package main

import (
	"fmt"
	"log"
	"os"

	"github.com/msolianko/btrfs"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/main.go <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  create <path>     - Create a new subvolume")
		fmt.Println("  quota <mountpoint> - Enable quota on mountpoint")
		fmt.Println("  id <path>         - Get subvolume ID")
		fmt.Println("  limit <mountpoint> <subvol_path> <bytes> - Set quota limit")
		fmt.Println("  delete <path>     - Delete an existing subvolume")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("create command requires a path argument")
		}
		path := os.Args[2]
		fmt.Printf("Creating subvolume at: %s\n", path)
		if err := btrfs.SubvolCreate(path); err != nil {
			log.Fatalf("Failed to create subvolume: %v", err)
		}
		fmt.Println("Subvolume created successfully")

	case "quota":
		if len(os.Args) < 3 {
			log.Fatal("quota command requires a mountpoint argument")
		}
		mountpoint := os.Args[2]
		fmt.Printf("Enabling quota on: %s\n", mountpoint)
		if err := btrfs.QuotaEnable(mountpoint); err != nil {
			log.Fatalf("Failed to enable quota: %v", err)
		}
		fmt.Println("Quota enabled successfully")

	case "id":
		if len(os.Args) < 3 {
			log.Fatal("id command requires a path argument")
		}
		path := os.Args[2]
		fmt.Printf("Getting subvolume ID for: %s\n", path)
		id, err := btrfs.GetSubvolID(path)
		if err != nil {
			log.Fatalf("Failed to get subvolume ID: %v", err)
		}
		fmt.Printf("Subvolume ID: %d\n", id)

	case "limit":
		if len(os.Args) < 5 {
			log.Fatal("limit command requires mountpoint, subvol_path, and bytes arguments")
		}
		mountpoint := os.Args[2]
		subvolPath := os.Args[3]
		bytes := os.Args[4]

		var maxBytes uint64
		if _, err := fmt.Sscanf(bytes, "%d", &maxBytes); err != nil {
			log.Fatalf("Invalid bytes value: %s", bytes)
		}

		fmt.Printf("Setting quota limit: %d bytes for %s on %s\n", maxBytes, subvolPath, mountpoint)
		if err := btrfs.QgroupLimit(mountpoint, subvolPath, maxBytes); err != nil {
			log.Fatalf("Failed to set quota limit: %v", err)
		}
		fmt.Println("Quota limit set successfully")

	case "delete":
		if len(os.Args) < 3 {
			log.Fatal("delete command requires a path argument")
		}
		path := os.Args[2]
		fmt.Printf("Deleting subvolume at: %s\n", path)
		if err := btrfs.SubvolDelete(path); err != nil {
			log.Fatalf("Failed to delete subvolume: %v", err)
		}
		fmt.Println("Subvolume deleted successfully")

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
