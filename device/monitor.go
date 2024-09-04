package device

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/gousb"
	"github.com/fsnotify/fsnotify"
	"github.com/sanjay7178/go-ukip/config"
	"github.com/sanjay7178/go-ukip/keystroke"
	"github.com/sanjay7178/go-ukip/logging"
)

type Monitor struct {
	cfg       *config.Config
	processor *keystroke.Processor
	context   *gousb.Context
	watcher   *fsnotify.Watcher
	devices   map[string]*gousb.Device
	mutex     sync.Mutex
	done      chan struct{}
}

func NewMonitor(cfg *config.Config) (*Monitor, error) {
	ctx := gousb.NewContext()
	processor, err := keystroke.NewProcessor(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create keystroke processor: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fs watcher: %v", err)
	}

	return &Monitor{
		cfg:       cfg,
		processor: processor,
		context:   ctx,
		watcher:   watcher,
		devices:   make(map[string]*gousb.Device),
		done:      make(chan struct{}),
	}, nil
}

func (m *Monitor) Start() error {
	if err := m.watcher.Add("/dev/input"); err != nil {
		return fmt.Errorf("failed to add /dev/input to watcher: %v", err)
	}

	go m.watchForDevices()
	go m.pollExistingDevices()

	return nil
}

func (m *Monitor) watchForDevices() {
	for {
		select {
		case event, ok := <-m.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				if strings.HasPrefix(filepath.Base(event.Name), "event") {
					go m.handleNewDevice(event.Name)
				}
			}
		case err, ok := <-m.watcher.Errors:
			if !ok {
				return
			}
			logging.Log.Errorf("Error watching for devices: %v", err)
		case <-m.done:
			return
		}
	}
}

func (m *Monitor) pollExistingDevices() {
	for {
		devices, err := ioutil.ReadDir("/dev/input")
		if err != nil {
			logging.Log.Errorf("Failed to read /dev/input: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, device := range devices {
			if strings.HasPrefix(device.Name(), "event") {
				devicePath := filepath.Join("/dev/input", device.Name())
				go m.handleNewDevice(devicePath)
			}
		}

		select {
		case <-time.After(30 * time.Second):
			// Poll every 30 seconds
		case <-m.done:
			return
		}
	}
}

func (m *Monitor) handleNewDevice(devicePath string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.devices[devicePath]; exists {
		return
	}

	dev, err := m.openDevice(devicePath)
	if err != nil {
		logging.Log.Errorf("Failed to open device %s: %v", devicePath, err)
		return
	}

	m.devices[devicePath] = dev
	go m.monitorDevice(devicePath, dev)
}

func (m *Monitor) openDevice(devicePath string) (*gousb.Device, error) {
	sysfsPath := fmt.Sprintf("/sys/class/input/%s/device", filepath.Base(devicePath))

	vendorID, err := readHexFile(filepath.Join(sysfsPath, "id_vendor"))
	if err != nil {
		return nil, fmt.Errorf("failed to read vendor ID: %v", err)
	}

	productID, err := readHexFile(filepath.Join(sysfsPath, "id_product"))
	if err != nil {
		return nil, fmt.Errorf("failed to read product ID: %v", err)
	}

	dev, err := m.context.OpenDeviceWithVIDPID(gousb.ID(vendorID), gousb.ID(productID))
	if err != nil {
		return nil, fmt.Errorf("failed to open USB device: %v", err)
	}

	return dev, nil
}

func (m *Monitor) monitorDevice(devicePath string, dev *gousb.Device) {
	defer func() {
		m.mutex.Lock()
		delete(m.devices, devicePath)
		m.mutex.Unlock()
		dev.Close()
	}()

	// Open the device file
	file, err := os.Open(devicePath)
	if err != nil {
		logging.Log.Errorf("Failed to open device file %s: %v", devicePath, err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 24) // Typical size for an input event

	for {
		select {
		case <-m.done:
			return
		default:
			n, err := file.Read(buffer)
			if err != nil {
				logging.Log.Errorf("Error reading from device %s: %v", devicePath, err)
				return
			}

			if n == 24 { // Full event read
				timestamp := time.Now()
				eventType := uint16(buffer[16]) | uint16(buffer[17])<<8
				code := uint16(buffer[18]) | uint16(buffer[19])<<8
				value := int32(buffer[20]) | int32(buffer[21])<<8 | int32(buffer[22])<<16 | int32(buffer[23])<<24

				if eventType == 1 && value == 1 { // Key press event
					m.processor.ProcessKeystroke(
						devicePath,
						dev.Desc.Product,
						fmt.Sprintf("%04x", dev.Desc.Vendor),
						fmt.Sprintf("%04x", dev.Desc.Product),
						rune(code),
						timestamp,
					)
				}
			}
		}
	}
}

func (m *Monitor) Stop() {
	close(m.done)
	m.watcher.Close()
	m.context.Close()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, dev := range m.devices {
		dev.Close()
	}
}

func readHexFile(path string) (uint16, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	trimmed := strings.TrimSpace(string(content))
	value, err := strconv.ParseUint(trimmed, 16, 16)
	if err != nil {
		return 0, err
	}

	return uint16(value), nil
}