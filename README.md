# go-ukip
cross platform runtime protection from usb keystroke injection , DNS assignment over DHCP spoof attacks via BadUSB 
# Under Development 

### Project Structure
```bash
ukip/
│
├── cmd/
│   └── ukip/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── device/
│   │   └── monitor.go
│   ├── keystroke/
│   │   └── processor.go
│   ├── allowlist/
│   │   └── allowlist.go
│   └── logging/
│       └── logging.go
│
├── configs/
│   ├── ukip.service
│   ├── allowlist.txt
│   └── keycodes.json
│
├── scripts/
│   └── install.sh
│
├── go.mod
├── go.sum
└── README.md
```

# UKIP Configuration Files Metadata

## ukip.service

### Purpose
This file is a systemd service unit configuration file for the UKIP (USB Keystroke Injection Protection) service. It defines how and when the UKIP service should be started, stopped, and managed by the systemd init system.

### Project Location
```
ukip/configs/ukip.service
```

### System Installation Location
```
/etc/systemd/system/ukip.service
```

### Content Format
Text file in systemd unit file format. It typically includes sections like [Unit], [Service], and [Install].

### Example Content
```ini
[Unit]
Description=UKIP
Requires=systemd-udevd.service
After=systemd-udevd.service

[Service]
ExecStart=/usr/local/bin/ukip
Restart=always
User=root

[Install]
WantedBy=multi-user.target
```

### Notes
- The service is set to start after the udev service, which is necessary for USB device detection.
- It's configured to always restart if it stops, ensuring continuous protection.
- The service runs as root, which is typically necessary for USB device monitoring.

## allowlist.txt

### Purpose
This file contains a list of allowed USB devices and their permitted keystrokes. It's used by UKIP to determine which devices should be exempt from keystroke injection protection or which specific keystrokes should be allowed for certain devices.

### Project Location
```
ukip/configs/allowlist.txt
```

### System Installation Location
```
/etc/ukip/allowlist
```

### Content Format
Text file with each line representing a device and its allowed keystrokes. The format is:

```
<product ID in hex>:<vendor ID in hex> <allowed characters, comma separated>
```

### Example Content
```
# Yubikey example
0x0010:0x1050 c,b,d,e,f,g,h,i,j,k,l,n,r,t,u,v

# Allow all characters for a specific device
0x1234:0x5678 any

# Block all characters for a specific device
0x9ABC:0xDEF0 none
```

### Notes
- Lines starting with '#' are treated as comments.
- The 'any' keyword allows all characters for a device.
- The 'none' keyword blocks all characters for a device.
- If a device is not listed, it's treated as if 'none' was specified.
- The file should be readable by the UKIP service (typically root-owned with 644 permissions).
```

This metadata provides a comprehensive overview of these two critical configuration files for the UKIP system. It includes their purpose, locations both within the project and on the installed system, content format, and example content. This information is crucial for developers working on UKIP, system administrators installing or managing UKIP, and anyone trying to understand or modify UKIP's configuration.