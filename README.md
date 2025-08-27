# BTRFS Go Package

A Go package for working with BTRFS (B-tree file system) operations including subvolume management and quota controls.

## Features

- **Subvolume Creation**: Create new BTRFS subvolumes
- **Quota Management**: Enable and configure quota limits
- **Subvolume Information**: Retrieve subvolume IDs and metadata
- **CGO Integration**: Direct Linux kernel interface via ioctl calls

## Requirements

- Linux kernel with BTRFS support
- Go 1.22 or later
- CGO enabled (required for Linux system calls)

## Installation

```bash
go get github.com/msolianko/btrfs
```

## Usage

### Basic Operations

```go
package main

import (
    "github.com/msolianko/btrfs"
)

func main() {
    // Create a new subvolume
    err := btrfs.SubvolCreate("/mnt/btrfs/new_subvol")
    if err != nil {
        panic(err)
    }

    // Enable quota on mountpoint
    err = btrfs.QuotaEnable("/mnt/btrfs")
    if err != nil {
        panic(err)
    }

    // Get subvolume ID
    id, err := btrfs.GetSubvolID("/mnt/btrfs/new_subvol")
    if err != nil {
        panic(err)
    }

    // Set quota limit (100MB)
    err = btrfs.QgroupLimit("/mnt/btrfs", "/mnt/btrfs/new_subvol", 100*1024*1024)
    if err != nil {
        panic(err)
    }
}
```

### Command Line Interface

The package includes a CLI tool for testing and demonstration:

```bash
# Create a subvolume
go run cmd/main.go create /mnt/btrfs/test_subvol

# Enable quota on mountpoint
go run cmd/main.go quota /mnt/btrfs

# Get subvolume ID
go run cmd/main.go id /mnt/btrfs/test_subvol

# Set quota limit (100MB)
go run cmd/main.go limit /mnt/btrfs /mnt/btrfs/test_subvol 104857600
```

## API Reference

### Functions

#### `SubvolCreate(path string) error`

Creates a new BTRFS subvolume at the specified path.

#### `QuotaEnable(mountpoint string) error`

Enables quota support on the BTRFS filesystem at the specified mountpoint.

#### `GetSubvolID(path string) (uint64, error)`

Retrieves the subvolume ID for the specified path.

#### `QgroupLimit(mountpoint, subvolPath string, maxBytes uint64) error`

Sets a quota limit for a subvolume at the specified path.

## Building

Since this package uses CGO and Linux-specific headers, it must be built in a Linux environment:

```bash
# Build the package
go build

# Build the CLI tool
go build -o btrfs-cli cmd/main.go
```

## Docker/Container Usage

For development and testing in containers:

```dockerfile
FROM golang:1.22-bullseye

# Install BTRFS tools and development headers
RUN apt-get update && apt-get install -y \
    btrfs-tools \
    libbtrfs-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

# Build the package
RUN go build
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Notes

- This package requires root privileges for most operations
- BTRFS must be compiled into the kernel or available as a module
- The package uses unsafe.Pointer for C struct manipulation
- All operations are synchronous and may block on I/O
