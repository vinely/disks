package disk

import (
	"log"
	"os/exec"
	"time"

	lockfile "github.com/vinely/disks/lockfile"
	"golang.org/x/sys/windows"
)

const (
	// MaxVolumeLabelLength is the maximum number of characters in a volume label.
	MaxVolumeLabelLength = windows.MAX_PATH + 1

	// MaxVolumeNameLength is the maximum number of characters in a volume name.
	//
	//   \\?\Volume{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}\
	MaxVolumeNameLength = windows.MAX_PATH + 1 // 50?

	// MaxFileSystemNameLength is the maximum number of characters in a file
	// system name.
	MaxFileSystemNameLength = windows.MAX_PATH + 1

	MaximumComponentLength = 255 //for FAT.
)

var (
	TmpMountPoint = ""
	VolumeName    [MaxVolumeNameLength]uint16
	// check timeout for every volume
	CheckEachTimeout = time.Duration(5)
)

func run(cmd string) []byte {
	out, err := exec.Command("cmd", "/C", cmd).Output()
	if err != nil {
		log.Printf("%v\n", err.Error())
	}
	return out
}

func GetAvailableLetter() string {
	drivers, err := windows.GetLogicalDrives()
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	for i := uint32(2); i < 26; i++ {
		if drivers&(1<<i) == 0 {
			return string(i + 'A')
		}
	}
	for i := uint32(0); i < 2; i++ {
		if drivers&(1<<i) == 0 {
			return string(i + 'A')
		}
	}
	return ""
}

func checkMountPoint(verifyFunc VerifyFunction, handle HandleFunction) {
	var (
		volumeMountPoint [MaxFileSystemNameLength]uint16
		strMountPoint    string
		ret_len          uint32
	)
	log.Print(windows.UTF16ToString(VolumeName[:]))
	err := windows.GetVolumePathNamesForVolumeName(&VolumeName[0], &volumeMountPoint[0], MaxFileSystemNameLength, &ret_len)
	if err != nil {
		log.Printf("%s\n", err.Error())
	} else {
		strMountPoint = windows.UTF16ToString(volumeMountPoint[:])
		var blk BlkInfo
		volumeName := windows.UTF16ToString(VolumeName[:])
		// \\?\Volume{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}\
		//
		id := windows.UTF16ToString(VolumeName[11:47])
		if strMountPoint == "" {
			// didn't have a mountpoint
			if TmpMountPoint == "" {
				TmpMountPoint = GetAvailableLetter()
				if TmpMountPoint == "" {
					log.Panic("Don't have available mount point!")
				}
				TmpMountPoint += ":\\"
			}
			tmp, _ := windows.UTF16FromString(TmpMountPoint)
			lockfile.TryLock(nil)
			defer lockfile.TryRelease()
			err = windows.SetVolumeMountPoint(&tmp[0], &VolumeName[0])
			defer windows.DeleteVolumeMountPoint(&tmp[0])
			if err != nil {
				log.Printf("%s\n", err.Error())
				return
			}
			blk = BlkInfo{
				Device: &BlockDevice{
					OS:     "windows",
					Device: volumeName,
					ID:     id,
					Type:   "volume"},
				Mount: &MountPoint{
					MountPoint: TmpMountPoint},
			}
			if verifyFunc(blk) {
				handle(blk)
			}
		} else {
			blk = BlkInfo{
				Device: &BlockDevice{
					OS:     "windows",
					Device: volumeName,
					ID:     id,
					Type:   "volume"},
				Mount: &MountPoint{
					MountPoint: strMountPoint},
			}
			if verifyFunc(blk) {
				handle(blk)
			}
		}
		log.Printf("BlockDevice: %v\n", blk.Device)
		log.Printf("MountPoint: %v\n", blk.Mount)

	}
}

func CheckVolume(verifyFunc VerifyFunction, handle HandleFunction) {
	t := time.NewTimer(CheckEachTimeout * time.Second)
	defer t.Stop()
	stopflag := make(chan bool, 1)
	goWrap := func() {
		checkMountPoint(verifyFunc, handle)
		stopflag <- true
	}

	hvol, err := windows.FindFirstVolume(&VolumeName[0], MaxVolumeNameLength)
	if err != nil {
		log.Printf("%s\n", err.Error())
	}
	defer windows.FindVolumeClose(hvol)

	go goWrap()
	select {
	case <-t.C:
	case <-stopflag:
	}

	for {
		if err := windows.FindNextVolume(hvol, &VolumeName[0], MaxVolumeNameLength); err != nil {
			break
		}
		t.Reset(CheckEachTimeout * time.Second)
		go goWrap()
		select {
		case <-t.C:
		case <-stopflag:
		}
	}
}
