package lockfile

import (
	"errors"
	"fmt"
	"os"
)

var ErrLocked = errors.New("Locked.")

type LockFile struct {
	Path string
}

func (l *LockFile) Lock() error {
	file, err := os.OpenFile(l.Path, os.O_CREATE|os.O_RDWR, os.ModeTemporary|0640)
	if err == nil {
		var pid int
		if _, err = fmt.Fscanf(file, "%d\n", &pid); err == nil {
			if pid != os.Getpid() {
				if ProcessRunning(pid) {
					file.Close()
					return ErrLocked
				}
			}
		}

		file.Seek(0, 0)
		file.Truncate(0)
		if _, err := fmt.Fprintf(file, "%d\n", os.Getpid()); err == nil {
			return nil
		} else {
			file.Close()
			return err
		}
	} else {
		return err
	}
}

func (l *LockFile) Release() {
	if file, err := os.OpenFile(l.Path, os.O_CREATE|os.O_RDWR, os.ModeTemporary|0640); err == nil {
		var pid int
		if _, err = fmt.Fscanf(file, "%d\n", &pid); err == nil {
			if pid != os.Getpid() {
				if ProcessRunning(pid) {
					file.Close()
					return
				}
			}
		}
		file.Close()
	}
	os.Remove(l.Path)
	l = nil
}

func (l *LockFile) Dispose() {
	if file, err := os.OpenFile(l.Path, os.O_CREATE|os.O_RDWR, os.ModeTemporary|0640); err == nil {
		var pid int
		_, err = fmt.Fscanf(file, "%d\n", &pid)
		if err == nil {
			p, err := os.FindProcess(pid)
			if err == nil {
				p.Release()
				p.Kill()
			}
		}
		file.Close()
	}
	os.Remove(l.Path)
}
