package main

import (
	disk "github.com/vinely/disks"
)

func main() {
	disk.DiskInfoPrint()

	vf := disk.CheckValidPathbyExist("version")
	disk.CheckVolume(vf, disk.HandleLs)
}
