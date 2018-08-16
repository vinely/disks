package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	disk "github.com/vinely/disks"
)

func main() {
	// diskinfo now only for windows
	log.SetOutput(ioutil.Discard)
	diskinfos := disk.DiskInfos("/version")
	fmt.Printf("Diskinfos:\n%+v\n", diskinfos)
	log.SetOutput(os.Stdout)

	// can check in windows and linux
	log.SetOutput(ioutil.Discard)
	vf := disk.CheckValidPathbyExist("/version")
	disk.CheckVolume(vf, disk.HandleSample)
	log.SetOutput(os.Stdout)

	// check device major letter for windows
	disks := disk.GetDisksAssociated()

	for _, disk := range disks {
		fmt.Print(disk.Index)
		fmt.Print("  -  ")
		for _, p := range disk.Partition {
			if p.Logical != nil {
				fmt.Print(p.Logical.DeviceID)
				break
			}
		}
		fmt.Println()
	}
}
