package keystroke

import (
	"container/ring"
	"time"

	"github.com/sanjay7178/go-ukip/allowlist"
	"github.com/sanjay7178/go-ukip/config"
	"github.com/sanjay7178/go-ukip/logging"
)

type Keystroke struct {
	Char      rune
	Timestamp time.Time
}

type Processor struct {
	cfg       *config.Config
	allowlist *allowlist.Allowlist
	devices   map[string]*DeviceData
}

type DeviceData struct {
	Keystrokes *ring.Ring
	Product    string
	VendorID   string
	ProductID  string
}

func NewProcessor(cfg *config.Config) (*Processor, error) {
	al, err := allowlist.Load()
	if err != nil {
		return nil, err
	}

	return &Processor{
		cfg:       cfg,
		allowlist: al,
		devices:   make(map[string]*DeviceData),
	}, nil
}

func (p *Processor) ProcessKeystroke(devicePath, product, vendorID, productID string, char rune, timestamp time.Time) {
	p.addKeystroke(devicePath, product, vendorID, productID, char, timestamp)
	p.checkForAttack(devicePath)
}

func (p *Processor) addKeystroke(devicePath, product, vendorID, productID string, char rune, timestamp time.Time) {
	if _, ok := p.devices[devicePath]; !ok {
		p.devices[devicePath] = &DeviceData{
			Keystrokes: ring.New(p.cfg.KeystrokeWindow),
			Product:    product,
			VendorID:   vendorID,
			ProductID:  productID,
		}
	}

	p.devices[devicePath].Keystrokes.Value = Keystroke{Char: char, Timestamp: timestamp}
	p.devices[devicePath].Keystrokes = p.devices[devicePath].Keystrokes.Next()
}

func (p *Processor) checkForAttack(devicePath string) {
	device := p.devices[devicePath]
	if device.Keystrokes.Len() < p.cfg.KeystrokeWindow {
		return
	}

	attackCounter := 0
	var prev Keystroke

	device.Keystrokes.Do(func(v interface{}) {
		if v == nil {
			return
		}

		current := v.(Keystroke)
		if !prev.Timestamp.IsZero() {
			timeDiff := current.Timestamp.Sub(prev.Timestamp)
			if timeDiff <= time.Duration(p.cfg.AbnormalTyping)*time.Microsecond {
				attackCounter++
			}
		}
		prev = current
	})

	if attackCounter == p.cfg.KeystrokeWindow-1 {
		p.handlePotentialAttack(devicePath)
	}
}

func (p *Processor) handlePotentialAttack(devicePath string) {
	device := p.devices[devicePath]
	deviceID := device.VendorID + ":" + device.ProductID

	// Check if all typed characters are in the allowlist
	allAllowed := true
	device.Keystrokes.Do(func(v interface{}) {
		if v == nil {
			return
		}
		keystroke := v.(Keystroke)
		if !p.allowlist.IsAllowed(deviceID, keystroke.Char) {
			allAllowed = false
		}
	})

	if allAllowed {
		return
	}

	if p.cfg.RunMode == "MONITOR" {
		p.monitorMode(devicePath)
	} else if p.cfg.RunMode == "HARDENING" {
		p.hardeningMode(devicePath)
	}
}

func (p *Processor) monitorMode(devicePath string) {
	device := p.devices[devicePath]
	logging.Log.Warningf(
		"[UKIP] The device %s with the vendor id %s and the product id %s would have been blocked. "+
			"The causing timings are: %v",
		device.Product, device.VendorID, device.ProductID, getTimings(device.Keystrokes),
	)
}

func (p *Processor) hardeningMode(devicePath string) {
	device := p.devices[devicePath]
	logging.Log.Warningf(
		"[UKIP] The device %s with the vendor id %s and the product id %s was blocked. "+
			"The causing timings are: %v",
		device.Product, device.VendorID, device.ProductID, getTimings(device.Keystrokes),
	)

	// Here you would implement the logic to unbind the device driver
	// This would typically involve interacting with the Linux sysfs
	// For example:
	// err := unbindDeviceDriver(devicePath)
	// if err != nil {
	//     logging.Log.Errorf("Failed to unbind device driver: %v", err)
	// }
}

func getTimings(keystrokes *ring.Ring) []time.Duration {
	var timings []time.Duration
	var prev Keystroke

	keystrokes.Do(func(v interface{}) {
		if v == nil {
			return
		}

		current := v.(Keystroke)
		if !prev.Timestamp.IsZero() {
			timings = append(timings, current.Timestamp.Sub(prev.Timestamp))
		}
		prev = current
	})

	return timings
}

// func unbindDeviceDriver(devicePath string) error {
//     // Implementation would go here
//     // This would involve writing to /sys/bus/usb/drivers/usb/unbind
//     return nil
// }
