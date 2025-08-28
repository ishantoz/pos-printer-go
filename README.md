# 🖨️ POS Printer Go

A high-performance Go service for printing barcodes and receipts to POS printers via USB connection. Built with modern Go practices, featuring a RESTful API, SQLite job queue, and cross-platform USB device support.

## ✨ Features

- **🖨️ USB Printer Support**: Direct communication with POS printers via USB
- **📊 Job Queue System**: SQLite-based job management with background processing
- **🔒 Secure API**: HTTPS server with configurable certificates
- **📱 RESTful API**: Simple HTTP endpoints for printing operations
- **🔄 Background Workers**: Asynchronous job processing with configurable workers
- **📏 Flexible Barcode Printing**: Customizable size, direction, and label gaps
- **🌐 Cross-Platform**: Works on Windows, macOS, and Linux
- **⚡ High Performance**: Built with Go for optimal performance

## 🚀 Quick Start

### Prerequisites

- **Go 1.24.4 or higher**
- **USB POS Printer** (Thermal, Label, or Receipt printer)
- **Printer's VID (Vendor ID) and PID (Product ID)**

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/ishantoz/pos-printer-go.git
   cd pos-printer-go
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Build the application**
   ```bash
   go build -o pos-printer ./cmd/pos-printer
   ```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# Server Configuration
POS_PRINTER_ENDPOINT=:5000
POS_PRINTER_SERVER_CERT_PATH=./certs/cert.pem
POS_PRINTER_SERVER_KEY_PATH=./certs/cert.key

# Database Configuration
POS_PRINTER_DB_SQLITE_PATH=./data/db/pos-printer.sqlite.db
POS_PRINTER_DB_MIGRATE=1

# Printer Configuration
POS_PRINTER_MAX_BARCODE_PRINT_COUNT=1000
POS_PRINTER_MAX_BARCODE_DATA_LENGTH=100
POS_PRINTER_MAX_TOP_TEXT_LENGTH=50

# Worker Configuration
POS_PRINTER_MAX_JOB_ATTEMPTS=3
POS_PRINTER_BARCODE_WORKER_COUNT=3
```

### SSL Certificates

For HTTPS support, place your SSL certificates in the `certs/` directory:
- `cert.pem` - SSL certificate
- `cert.key` - Private key

## 🖥️ Operating System Setup

### Windows Setup

1. **Install USB Driver (Required)**
   - Download Zadig from the `usb-driver/` folder
   - Extract `zadig-2.9.zip`
   - Run Zadig as Administrator
   - Connect your POS printer via USB
   - Select your printer device
   - Install the WinUSB driver

2. **Find Printer VID/PID**
   - Open Device Manager
   - Look for your printer under "Universal Serial Bus controllers"
   - Right-click → Properties → Details → Hardware Ids
   - Note the VID and PID values (e.g., `USB\VID_6EF0&PID_6550`)

3. **Run the Service**
   ```cmd
   pos-printer.exe
   ```

### macOS Setup

1. **Install Go (if not already installed)**
   ```bash
   brew install go
   ```

2. **Find Printer VID/PID**
   ```bash
   system_profiler SPUSBDataType | grep -A 10 -B 5 "Your Printer Name"
   ```
   Look for `Product ID: 0xXXXX` and `Vendor ID: 0xXXXX`

3. **Run the Service**
   ```bash
   ./pos-printer
   ```

### Linux Setup

1. **Install Dependencies**
   ```bash
   # Ubuntu/Debian
   sudo apt-get update
   sudo apt-get install libusb-1.0-0-dev

   # CentOS/RHEL/Fedora
   sudo yum install libusb1-devel
   # or
   sudo dnf install libusb1-devel
   ```

2. **Find Printer VID/PID**
   ```bash
   lsusb
   ```
   Look for your printer and note the VID:PID format (e.g., `6ef0:6550`)

3. **Set USB Permissions (Optional)**
   ```bash
   # Create udev rule for persistent access
   sudo nano /etc/udev/rules.d/99-pos-printer.rules
   ```
   
   Add this line (replace VID and PID with your values):
   ```
   SUBSYSTEM=="usb", ATTRS{idVendor}=="6ef0", ATTRS{idProduct}=="6550", MODE="0666"
   ```
   
   Reload rules:
   ```bash
   sudo udevadm control --reload-rules
   sudo udevadm trigger
   ```

4. **Run the Service**
   ```bash
   ./pos-printer
   ```

## 📡 API Usage

### Health Check
```bash
curl -k https://localhost:5000/health
```

### Print Barcode
```bash
curl -k -X POST https://localhost:5000/barcode/print \
  -H "Content-Type: application/json" \
  -d '{
    "vid": "0x6ef0",
    "pid": "0x6550",
    "sizeX": 55,
    "sizeY": 45,
    "direction": 0,
    "topText": "100tk",
    "barcodeData": "AX2B2CL21LL2",
    "printCount": 2,
    "labelGap": {
      "length": 0,
      "offset": 0
    }
  }'
```

### Check Job Status
```bash
curl -k https://localhost:5000/barcode/job/{jobId}
```

## 📋 Request Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `vid` | string | Vendor ID (hex format: 0x6EF0) | Required |
| `pid` | string | Product ID (hex format: 0x6550) | Required |
| `sizeX` | int | Barcode width in mm | Required |
| `sizeY` | int | Barcode height in mm | Required |
| `direction` | int | Print direction (0=horizontal, 1=vertical) | 0 |
| `topText` | string | Text above barcode | "" |
| `barcodeData` | string | Barcode content | Required |
| `printCount` | int | Number of copies to print | 1 |
| `labelGap` | object | Label gap configuration | Auto-detect |

### Label Gap Configuration
```json
{
  "length": 0,    // Gap length in mm (0 = auto-detect)
  "offset": 0     // Gap offset in mm
}
```

## 🏗️ Project Structure

```
pos-printer-go/
├── cmd/                    # Application entry points
│   ├── pos-printer/       # Main service
│   └── escpos-test/       # ESC/POS testing utility
├── internal/               # Internal packages
│   ├── api/               # HTTP API handlers
│   ├── config/            # Configuration management
│   ├── db/                # Database operations
│   ├── helper/            # Utility functions
│   ├── job/               # Job processing system
│   ├── lib/               # External library wrappers
│   ├── model/             # Data models
│   └── printer/           # Printer communication
├── assets/                 # Static assets
├── certs/                  # SSL certificates
├── usb-driver/            # Windows USB drivers
├── go.mod                 # Go module file
└── README.md              # This file
```

## 🔍 Finding Your Printer's VID/PID

### Windows
1. Connect printer via USB
2. Open Device Manager
3. Expand "Universal Serial Bus controllers"
4. Find your printer device
5. Right-click → Properties → Details → Hardware Ids
6. Look for format: `USB\VID_XXXX&PID_XXXX`

### macOS
```bash
system_profiler SPUSBDataType | grep -A 10 -B 5 "Printer Name"
```

### Linux
```bash
lsusb | grep "Printer Name"
# or
dmesg | grep -i "usb.*printer"
```

## 🚨 Troubleshooting

### Common Issues

1. **Printer not found**
   - Verify USB connection
   - Check VID/PID values
   - Ensure proper driver installation (Windows)
   - Check USB permissions (Linux)

2. **Permission denied (Linux)**
   - Run with sudo or set up udev rules
   - Check USB device permissions

3. **Driver issues (Windows)**
   - Use Zadig to install WinUSB driver
   - Run Zadig as Administrator
   - Restart after driver installation

4. **SSL certificate errors**
   - Check certificate paths in `.env`
   - Ensure certificates are valid
   - Use `-k` flag with curl for self-signed certs

### Debug Mode

Enable debug logging by setting environment variables:
```bash
export POS_PRINTER_DEBUG=1
```

## 🧪 Testing

### Test ESC/POS Commands
```bash
go run ./cmd/escpos-test/main.go
```

### API Testing
Use the included `client.http` file with REST Client extensions in VS Code or similar tools.

## 📦 Building for Distribution

### Windows
```bash
GOOS=windows GOARCH=amd64 go build -o pos-printer.exe ./cmd/pos-printer
```

### macOS
```bash
GOOS=darwin GOARCH=amd64 go build -o pos-printer ./cmd/pos-printer
```

### Linux
```bash
GOOS=linux GOARCH=amd64 go build -o pos-printer ./cmd/pos-printer
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/ishantoz/pos-printer-go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ishantoz/pos-printer-go/discussions)
- **Wiki**: [Project Wiki](https://github.com/ishantoz/pos-printer-go/wiki)

## 🙏 Acknowledgments

- [gousb](https://github.com/google/gousb) - USB device communication
- [Echo](https://echo.labstack.com/) - Web framework
- [SQLite](https://www.sqlite.org/) - Database engine

---

**Made with ❤️ using Go**
