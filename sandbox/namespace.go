package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

//noinspection GoUnusedExportedFunction
func InitNamespace(newRoot string) error {
	_, _ = os.Stderr.WriteString(fmt.Sprintf("InitNamespace(%s) starting...\n", newRoot))

	if err := pivotRoot(newRoot); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("pivotRoot(%s) failed, err: %s\n", newRoot, err.Error()))
		return err
	}

	if err := syscall.Sethostname([]byte("justice")); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("syscall.Sethostname failed, err: %s\n", err.Error()))
		return err
	}

	_, _ = os.Stderr.WriteString(fmt.Sprintf("InitNamespace(%s) done\n", newRoot))
	return nil
}

func pivotRoot(newRoot string) error {
	putOld := filepath.Join(newRoot, "/.pivot_root")

	// bind mount new_root to itself - this is a slight hack needed to satisfy requirement (2)
	//
	// The following restrictions apply to new_root and put_old:
	// 1.  They must be directories.
	// 2.  new_root and put_old must not be on the same filesystem as the current root.
	// 3.  put_old must be underneath new_root, that is, adding a nonzero
	//     number of /.. to the string pointed to by put_old must yield the same directory as new_root.
	// 4.  No other filesystem may be mounted on put_old.
	if err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("syscall.Mount(%s, %s, \"\", syscall.MS_BIND|syscall.MS_REC, \"\") failed\n", newRoot, newRoot))
		return err
	}

	// create put_old directory
	if err := os.MkdirAll(putOld, 0700); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("os.MkdirAll(%s, 0700) failed\n", putOld))
		return err
	}

	// call pivotRoot
	if err := syscall.PivotRoot(newRoot, putOld); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("syscall.PivotRoot(%s, %s) failed\n", newRoot, putOld))
		return err
	}

	// Note that this also applies to the calling process: pivotRoot() may
	// or may not affect its current working directory.  It is therefore
	// recommended to call chdir("/") immediately after pivotRoot().
	if err := os.Chdir("/"); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("os.Chdir(\"/\") failed\n"))
		return err
	}

	// umount put_old, which now lives at /.pivot_root
	putOld = "/.pivot_root"
	if err := syscall.Unmount(putOld, syscall.MNT_DETACH); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("syscall.Unmount(%s, syscall.MNT_DETACH) failed\n", putOld))
		return err
	}

	// remove put_old
	if err := os.RemoveAll(putOld); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("os.RemoveAll(%s) failed\n", putOld))
		return err
	}

	return nil
}
