package disk

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type VerifyFunction = func(blk BlkInfo) bool
type HandleFunction = func(blk BlkInfo)

type BlockDevice struct {
	OS     string // os type
	Device string // device name
	ID     string // id (windows guid, linux uuid)
	Type   string // filesystem type
}

type BlkInfo struct {
	Device *BlockDevice
	Mount  *MountPoint
}

func CheckValidPathbyExist(file string) VerifyFunction {
	var VeriFile = file
	return func(blk BlkInfo) bool {
		if blk.Mount == nil {
			log.Println("Didn't have a valid mountpoint!")
			return false
		}
		filepath := filepath.Join(blk.Mount.MountPoint, VeriFile)
		_, err := os.Stat(filepath)
		if err != nil {
			log.Println(err)
			return false
		} else {
			return true
		}
	}
}

func HandleSample(blk BlkInfo) {
	if blk.Mount == nil {
		log.Println("Didn't have a valid mountpoint!")
		return
	}
	fmt.Println("Valid in : " + blk.Mount.MountPoint)

}
