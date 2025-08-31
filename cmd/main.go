package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MykolaSuperman/btrfs"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:           "btrfs-cli",
		Short:         "Tiny Btrfs helper (subvol + quota)",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// create
	var createPath string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new subvolume",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createPath == "" {
				return fmt.Errorf("create requires --path")
			}
			fmt.Printf("Creating subvolume at: %s\n", createPath)
			if err := btrfs.SubvolCreate(createPath); err != nil {
				return fmt.Errorf("create failed: %w", err)
			}
			fmt.Println("Subvolume created successfully")
			return nil
		},
	}
	createCmd.Flags().StringVar(&createPath, "path", "", "Path to create subvolume")

	// quota enable
	var quotaMount string
	quotaCmd := &cobra.Command{
		Use:   "quota",
		Short: "Enable quota on a mountpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			if quotaMount == "" {
				return fmt.Errorf("quota requires --mount")
			}
			fmt.Printf("Enabling quota on: %s\n", quotaMount)
			if err := btrfs.QuotaEnable(quotaMount); err != nil {
				return fmt.Errorf("quota enable failed: %w", err)
			}
			fmt.Println("Quota enabled successfully")
			return nil
		},
	}
	quotaCmd.Flags().StringVar(&quotaMount, "mount", "", "Mountpoint where quota should be enabled")

	// id
	var idPath string
	idCmd := &cobra.Command{
		Use:   "id",
		Short: "Get subvolume ID (treeid)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if idPath == "" {
				return fmt.Errorf("id requires --path")
			}
			fmt.Printf("Getting subvolume ID for: %s\n", idPath)
			id, err := btrfs.GetSubvolID(idPath)
			if err != nil {
				return fmt.Errorf("get id failed: %w", err)
			}
			fmt.Printf("Subvolume ID: %d\n", id)
			return nil
		},
	}
	idCmd.Flags().StringVar(&idPath, "path", "", "Path of subvolume to inspect")

	// limit
	var limitMount, limitSubvol string
	var limitBytes uint64
	limitCmd := &cobra.Command{
		Use:   "limit",
		Short: "Set quota limit for a subvolume (bytes)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if limitMount == "" || limitSubvol == "" || limitBytes == 0 {
				return fmt.Errorf("limit requires --mount, --subvol, and --bytes > 0")
			}
			fmt.Printf("Setting quota limit: %d bytes for %s on %s\n", limitBytes, limitSubvol, limitMount)
			if err := btrfs.QgroupLimit(limitMount, limitSubvol, limitBytes); err != nil {
				return fmt.Errorf("set limit failed: %w", err)
			}
			fmt.Println("Quota limit set successfully")
			return nil
		},
	}
	limitCmd.Flags().StringVar(&limitMount, "mount", "", "Mountpoint where quota is enabled")
	limitCmd.Flags().StringVar(&limitSubvol, "subvol", "", "Subvolume path to limit")
	limitCmd.Flags().Uint64Var(&limitBytes, "bytes", 0, "Quota size in bytes")

	// delete
	var deletePath string
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an existing subvolume",
		RunE: func(cmd *cobra.Command, args []string) error {
			if deletePath == "" {
				return fmt.Errorf("delete requires --path")
			}
			fmt.Printf("Deleting subvolume at: %s\n", deletePath)
			if err := btrfs.SubvolDelete(deletePath); err != nil {
				return fmt.Errorf("delete failed: %w", err)
			}
			fmt.Println("Subvolume deleted successfully")
			return nil
		},
	}
	deleteCmd.Flags().StringVar(&deletePath, "path", "", "Path of subvolume to delete")

	// wire up
	rootCmd.AddCommand(createCmd, quotaCmd, idCmd, limitCmd, deleteCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
