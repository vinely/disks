package disk

type DiskInfo struct {
	Disk      Win32_DiskDrive
	Installed bool
}

func DiskInfos(symbolFile string) []DiskInfo {
	disks := GetDiskDrive()
	diskinfos := make([]DiskInfo, len(disks))
	for index, _ := range disks {
		diskinfos[index].Disk = disks[index]
		diskinfos[index].Installed = false
	}
	handle := func(blk BlkInfo) {
		var no, ok = uint32(0), true
		if no, ok = FoundDiskNoFromMountPoint(blk.Mount.MountPoint); !ok {
			return
		}
		for index, _ := range diskinfos {
			if diskinfos[index].Disk.Index == no {
				diskinfos[index].Installed = true
			}
		}
	}
	vf := CheckValidPathbyExist(symbolFile)
	CheckVolume(vf, handle)
	return diskinfos
}
