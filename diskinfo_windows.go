package disk

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	InstalledFile = "\\version"
	DiskInfos     []DiskInfo
)

func HandlInstalled(blk BlkInfo) {
	mp := []byte(blk.Mount.MountPoint)
	no, _ := strconv.Atoi(FoundDiskNo(string(mp[0:2])))
	for index, _ := range DiskInfos {
		if DiskInfos[index].Disk.Index == uint32(no) {
			DiskInfos[index].Installed = true
		}
	}
}

func VerifyInstalled(blk BlkInfo) bool {
	filepath := filepath.Join(blk.Mount.MountPoint, InstalledFile)
	_, err := os.Stat(filepath)
	if err != nil {
		return false
	} else {
		return true
	}
}

func DiskInfoPrint() {
	DiskInfos = GetDiskDrive()
	CheckVolume(VerifyInstalled, HandlInstalled)
	log.Printf("%+v\n", DiskInfos)
}
