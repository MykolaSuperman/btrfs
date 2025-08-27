package main

import (
	"flag"
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

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	quotaCmd := flag.NewFlagSet("quota", flag.ExitOnError)
	idCmd := flag.NewFlagSet("id", flag.ExitOnError)
	limitCmd := flag.NewFlagSet("limit", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	switch command {
	case "create":
		createPath := createCmd.String("path", "", "Path to create subvolume")
		createCmd.Parse(os.Args[2:])
		if *createPath == "" {
			log.Fatal("create requires -path")
		}
		fmt.Printf("Creating subvolume at: %s\n", *createPath)
		if err := btrfs.SubvolCreate(*createPath); err != nil {
			log.Fatalf("Failed to create subvolume: %v", err)
		}
		fmt.Println("Subvolume created successfully")

	case "quota":
		quotaMount := quotaCmd.String("mount", "", "Mountpoint to enable quota")
		quotaCmd.Parse(os.Args[2:])
		if *quotaMount == "" {
			log.Fatal("quota requires -mount")
		}
		fmt.Printf("Enabling quota on: %s\n", *quotaMount)
		if err := btrfs.QuotaEnable(*quotaMount); err != nil {
			log.Fatalf("Failed to enable quota: %v", err)
		}
		fmt.Println("Quota enabled successfully")

	case "id":
		idPath := idCmd.String("path", "", "Path to get subvolume ID")
		idCmd.Parse(os.Args[2:])
		if *idPath == "" {
			log.Fatal("id requires -path")
		}
		fmt.Printf("Getting subvolume ID for: %s\n", *idPath)
		id, err := btrfs.GetSubvolID(*idPath)
		if err != nil {
			log.Fatalf("Failed to get subvolume ID: %v", err)
		}
		fmt.Printf("Subvolume ID: %d\n", id)

	case "limit":
		limitMount := limitCmd.String("mount", "", "Mountpoint for quota limit")
		limitSubvol := limitCmd.String("subvol", "", "Subvolume path for quota limit")
		limitBytes := limitCmd.Uint64("bytes", 0, "Bytes for quota limit")

		limitCmd.Parse(os.Args[2:])
		if *limitMount == "" || *limitSubvol == "" || *limitBytes == 0 {
			log.Fatal("limit requires -mount, -subvol, and -bytes")
		}
		fmt.Printf("Setting quota limit: %d bytes for %s on %s\n", *limitBytes, *limitSubvol, *limitMount)
		if err := btrfs.QgroupLimit(*limitMount, *limitSubvol, *limitBytes); err != nil {
			log.Fatalf("Failed to set quota limit: %v", err)
		}
		fmt.Println("Quota limit set successfully")

	case "delete":
		deletePath := deleteCmd.String("path", "", "Path to delete subvolume")
		deleteCmd.Parse(os.Args[2:])
		if *deletePath == "" {
			log.Fatal("delete requires -path")
		}
		fmt.Printf("Deleting subvolume at: %s\n", *deletePath)
		if err := btrfs.SubvolDelete(*deletePath); err != nil {
			log.Fatalf("Failed to delete subvolume: %v", err)
		}
		fmt.Println("Subvolume deleted successfully")

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
