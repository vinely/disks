package disk

import (
	"strconv"
)

func DiskInfos(symbolFile string) []DiskInfo {
	diskinfos := GetDiskDrive()
	handle := func(blk BlkInfo) {
		mp := []byte(blk.Mount.MountPoint)
		no, _ := strconv.Atoi(FoundDiskNo(string(mp[0:2])))
		for index, _ := range diskinfos {
			if diskinfos[index].Disk.Index == uint32(no) {
				diskinfos[index].Installed = true
			}
		}
	}
	vf := CheckValidPathbyExist(symbolFile)
	CheckVolume(vf, handle)
	return diskinfos
}
