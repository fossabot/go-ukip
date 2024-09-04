package allowlist

import (
	"bufio"
	"os"
	"strings"
)

type Allowlist struct {
	devices map[string][]rune
}

func Load() (*Allowlist, error) {
	file, err := os.Open("/etc/ukip/allowlist")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	al := &Allowlist{devices: make(map[string][]rune)}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		deviceID := parts[0]
		if parts[1] == "any" {
			al.devices[deviceID] = nil // nil means any character is allowed
		} else if parts[1] != "none" {
			al.devices[deviceID] = []rune(parts[1])
		}
	}

	return al, scanner.Err()
}

func (al *Allowlist) IsAllowed(deviceID string, keystroke rune) bool {
	allowed, ok := al.devices[deviceID]
	if !ok {
		return false
	}
	if allowed == nil {
		return true // Any character is allowed
	}
	for _, r := range allowed {
		if r == keystroke {
			return true
		}
	}
	return false
}