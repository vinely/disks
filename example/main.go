package main

import (
	"fmt"

	disk "github.com/vinely/disks"
)

func main() {
	// diskinfo now only for windows
	diskinfos := disk.DiskInfos("/version")
	fmt.Printf("Diskinfos:\n%+v\n", diskinfos)

	// can check in windows and linux
	vf := disk.CheckValidPathbyExist("/version")
	disk.CheckVolume(vf, disk.HandleLs)
}
