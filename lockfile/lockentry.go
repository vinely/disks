package lockfile

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrTimeout = errors.New("Timeout for Lock!")
	lockFile   = LockFile{
		Path: filepath.Join(os.TempDir(), "mlock"),
	}
)

type TimeOutFunctions = func() bool

func SetLockFile(path string) {
	lockFile.Path = path
}

func TryLock(tof TimeOutFunctions) error {
	var err error
	for err = lockFile.Lock(); err != nil; {
		if tof == nil {
			tof = DefaultTimeoutFunction(10, 10)
		}
		if tof() {
			return ErrTimeout
		}
	}
	return nil
}

func TryRelease() {
	lockFile.Release()
}

func DefaultTimeoutFunction(interval time.Duration, count int) TimeOutFunctions {
	index := 0
	intval := interval
	cnt := count
	return func() bool {
		time.Sleep(intval)
		index++
		if cnt > 10 {
			return true
		} else {
			return false
		}
	}
}
