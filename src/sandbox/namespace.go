package sandbox

import (
	"path/filepath"
	"syscall"
	"os"
)

func InitNamespace(newRoot string) error {
	os.Stderr.WriteString("InitNamespace starting...\n")

	if err := pivotRoot(newRoot); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("pivotRoot failed...\n")
		return err
	}

	if err := mountProc(newRoot); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("mountProc failed...\n")
		return err
	}

	if err := syscall.Sethostname([]byte("justice")); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("Sethostname failed...\n")
		return err
	}

	os.Stderr.WriteString("InitNamespace done \n")
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
		os.Stderr.WriteString("Mount newRoot failed...\n")
		return err
	}

	// create put_old directory
	if err := os.MkdirAll(putOld, 0700); err != nil {
		os.Stderr.WriteString("Mkdir putOld failed...\n")
		return err
	}

	// call pivotRoot
	if err := syscall.PivotRoot(newRoot, putOld); err != nil {
		os.Stderr.WriteString("PivotRoot failed...\n")
		return err
	}

	// Note that this also applies to the calling process: pivotRoot() may
	// or may not affect its current working directory.  It is therefore
	// recommended to call chdir("/") immediately after pivotRoot().
	if err := os.Chdir("/"); err != nil {
		os.Stderr.WriteString("Chdir / failed...\n")
		return err
	}

	// umount put_old, which now lives at /.pivot_root
	putOld = "/.pivot_root"
	if err := syscall.Unmount(putOld, syscall.MNT_DETACH); err != nil {
		os.Stderr.WriteString("Unmount putOld failed...\n")
		return err
	}

	// remove put_old
	if err := os.RemoveAll(putOld); err != nil {
		os.Stderr.WriteString("Remove putOld failed...\n")
		return err
	}

	return nil
}

func mountProc(newRoot string) error {
	target := filepath.Join(newRoot, "/proc")

	if err := os.MkdirAll(target, 0755); err != nil {
		os.Stderr.WriteString("Mkdir target failed...\n")
		return err
	}

	if err := syscall.Mount("proc", target, "proc", uintptr(0), ""); err != nil {
		os.Stderr.WriteString("Mount proc failed...\n")
		return err
	}

	return nil
}
