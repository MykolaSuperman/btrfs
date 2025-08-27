package btrfs

/*
#include <linux/btrfs.h>
#include <sys/ioctl.h>
#include <linux/ioctl.h>
#include <string.h>
#include <errno.h>
#include <stdlib.h>

static inline int btrfs_subvol_create_v2_c(int dirfd, const char* name, int *perrno) {
    struct btrfs_ioctl_vol_args_v2 args;
    memset(&args, 0, sizeof(args));
    args.fd = dirfd;
    size_t nlen = strlen(name);
    if (nlen >= BTRFS_VOL_NAME_MAX) nlen = BTRFS_VOL_NAME_MAX - 1;
    memcpy(args.name, name, nlen);
    args.name[nlen] = '\0';
    int ret = ioctl(dirfd, BTRFS_IOC_SUBVOL_CREATE_V2, &args);
    if (ret < 0) { if (perrno) *perrno = errno; return ret; }
    return 0;
}

static inline int btrfs_quota_enable_c(int fd, int *perrno) {
    struct btrfs_ioctl_quota_ctl_args ctl;
    memset(&ctl, 0, sizeof(ctl));
    ctl.cmd = BTRFS_QUOTA_CTL_ENABLE;
    int ret = ioctl(fd, BTRFS_IOC_QUOTA_CTL, &ctl);
    if (ret < 0) { if (perrno) *perrno = errno; return ret; }
    return 0;
}

static inline int btrfs_get_subvol_id_c(int fd, unsigned long long* treeid, int *perrno) {
    struct btrfs_ioctl_get_subvol_info_args info;
    memset(&info, 0, sizeof(info));
    int ret = ioctl(fd, BTRFS_IOC_GET_SUBVOL_INFO, &info);
    if (ret < 0) { if (perrno) *perrno = errno; return ret; }
    if (treeid) *treeid = info.treeid;
    return 0;
}

static inline int btrfs_qgroup_limit_c(int fd, unsigned long long qgroupid, unsigned long long max_rfer, int *perrno) {
    struct btrfs_ioctl_qgroup_limit_args args;
    memset(&args, 0, sizeof(args));
    args.qgroupid = qgroupid;
    args.lim.flags = BTRFS_QGROUP_LIMIT_MAX_RFER;
    args.lim.max_rfer = max_rfer;
    int ret = ioctl(fd, BTRFS_IOC_QGROUP_LIMIT, &args);
    if (ret < 0) { if (perrno) *perrno = errno; return ret; }
    return 0;
}

static inline int btrfs_subvol_delete_c(int dirfd, const char* name, int *perrno) {
    struct btrfs_ioctl_vol_args args;
    memset(&args, 0, sizeof(args));
    size_t nlen = strlen(name);
    if (nlen >= BTRFS_VOL_NAME_MAX) nlen = BTRFS_VOL_NAME_MAX - 1;
    memcpy(args.name, name, nlen);
    args.name[nlen] = '\0';
    int ret = ioctl(dirfd, BTRFS_IOC_SNAP_DESTROY, &args);
    if (ret < 0) { if (perrno) *perrno = errno; return ret; }
    return 0;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

// SubvolCreate creates a new BTRFS subvolume at the specified path
func SubvolCreate(path string) error {
	dir, name := filepath.Split(filepath.Clean(path))

	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var cerr C.int
	ret := C.btrfs_subvol_create_v2_c(C.int(d.Fd()), cname, &cerr)
	if ret != 0 {
		return syscall.Errno(cerr)
	}
	return nil
}

// SubvolDelete deletes an existing BTRFS subvolume at the specified path
func SubvolDelete(path string) error {
	dir, name := filepath.Split(filepath.Clean(path))

	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var cerr C.int
	ret := C.btrfs_subvol_delete_c(C.int(d.Fd()), cname, &cerr)
	if ret != 0 {
		return syscall.Errno(cerr)
	}
	return nil
}

// QuotaEnable enables quota support on the BTRFS filesystem at the specified mountpoint
func QuotaEnable(mountpoint string) error {
	mp, err := os.Open(mountpoint)
	if err != nil {
		return err
	}
	defer mp.Close()

	var cerr C.int
	ret := C.btrfs_quota_enable_c(C.int(mp.Fd()), &cerr)
	if ret != 0 {
		return fmt.Errorf("quota enable: %w", syscall.Errno(cerr))
	}
	return nil
}

// GetSubvolID retrieves the subvolume ID for the specified path
func GetSubvolID(path string) (uint64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var treeid C.ulonglong
	var cerr C.int
	ret := C.btrfs_get_subvol_id_c(C.int(f.Fd()), &treeid, &cerr)
	if ret != 0 {
		return 0, fmt.Errorf("get subvol info: %w", syscall.Errno(cerr))
	}
	return uint64(treeid), nil
}

// QgroupLimit sets a quota limit for a subvolume at the specified path
func QgroupLimit(mountpoint, subvolPath string, maxBytes uint64) error {
	// ensure quota is enabled on the fs that contains mountpoint
	if err := QuotaEnable(mountpoint); err != nil && !errors.Is(err, syscall.EEXIST) {
		// enabling quota twice may error; ignore benign failure if already enabled
	}

	id, err := GetSubvolID(subvolPath)
	if err != nil {
		return err
	}

	mp, err := os.Open(mountpoint)
	if err != nil {
		return err
	}
	defer mp.Close()

	var cerr C.int
	ret := C.btrfs_qgroup_limit_c(C.int(mp.Fd()), C.ulonglong(id), C.ulonglong(maxBytes), &cerr)
	if ret != 0 {
		return fmt.Errorf("qgroup limit: %w", syscall.Errno(cerr))
	}
	return nil
}

// ioctl is a helper function to perform ioctl operations
func ioctl(fd uintptr, request, arg uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, request, arg)
	if errno != 0 {
		return errno
	}
	return nil
}
