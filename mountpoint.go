package disk

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	workDir = os.TempDir()
)

type MountPoint struct {
	MountPoint string
}

var ErrOther = errors.New("Unknown Error.")

func SetMountPath(path string) {
	workDir = path
}

func GetMountPoint(name string) (*MountPoint, error) {
	var (
		err error
		mp  = MountPoint{}
	)
	if name == "" {
		mp.MountPoint = filepath.Join(workDir, "/tmp")
	} else {
		mp.MountPoint = filepath.Join(workDir, name)
	}
	err = os.MkdirAll(mp.MountPoint, 0640)
	if err != nil {
		return nil, err
	}
	return &mp, nil
}

func (mp *MountPoint) Release() {
	os.Remove(mp.MountPoint)
}
