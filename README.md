# go-ukip

cross platform runtime protection from usb keystroke injection , DNS assignment over DHCP spoof attacks via BadUSB

# Under Development

### Installation

1. Clone the repository:

   ```
   git clone https://github.com/sanjay7178/go-ukip.git
   cd go-ukip
   ```

2. Build the project:

   ```
   go build
   ```

3. After building your Go project, run the installation script:

   ```
   sudo ./install.sh
   ```

   This script will copy the binary to `/usr/local/bin/ukip`, copy the configuration files to `/etc/ukip/`, and set up the systemd service.

4. You can then start the service with:

   ```
   sudo systemctl start ukip
   ```

   And check its status with:

   ```
   sudo systemctl status ukip
   ```

5. To ensure the service starts on boot, enable it:

   ```
   sudo systemctl enable ukip
   ```

Remember, since UKIP needs to access USB devices, it typically needs to run as root. That's why the service file specifies `User=root` and why the installation script needs to be run with sudo.

Also, make sure that your Go binary is built for the correct architecture of your system. If you're building on the same system where you'll run UKIP, this shouldn't be an issue.

If you make changes to the UKIP binary in the future, you'll need to copy it to `/usr/local/bin/` again and restart the service:

```
sudo cp ukip /usr/local/bin/ukip
sudo systemctl restart ukip
```

This process ensures that your UKIP binary is in the correct location, the service file points to this location, and the service is set up to run automatically on system boot.

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


### Contributing Instructions

The `gousb` package requires `libusb-1.0` to be installed on your system. Let's address this issue step by step:

1. First, we need to install the required dependencies. On a Ubuntu/Debian system, you can do this with the following commands:

```

sudo apt update
sudo apt install libusb-1.0-0-dev pkg-config

```

2. After installing these packages, try building your project again:

```

go build

```

If you still encounter issues, here are a few more things to check:

3. Ensure that the `PKG_CONFIG_PATH` environment variable includes the directory containing `libusb-1.0.pc`. You can check this with:

```

pkg-config --list-all | grep libusb

```

If it doesn't show up, you might need to set the `PKG_CONFIG_PATH`. The exact path can vary, but it's often something like:

```

export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/lib/x86_64-linux-gnu/pkgconfig

```

4. If you're still having issues, you might need to install `gcc` if it's not already on your system:

```

sudo apt install build-essential

```

5. After making these changes, clean your Go module cache and try building again:

```

go clean -modcache
go build

```

If you continue to face issues, please provide the full error message you're getting after trying these steps. It would also be helpful to know the output of:

```

go version
uname -a

```

This will give us more information about your Go version and system, which can help in troubleshooting.

Remember, working with USB devices often requires elevated privileges. When you run your UKIP application, you might need to use `sudo` or set up appropriate udev rules to allow non-root access to USB devices.
```
