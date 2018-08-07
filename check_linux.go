// +build  linux

package disk

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	lockfile "github.com/vinely/disks/lockfile"
)

func run(cmd string) []byte {
	out, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		log.Printf("%v\n", err.Error())
	}
	return out
}

func CheckVolume(verifyFunc VerifyFunction, handle HandleFunction) {
	if verifyFunc == nil || handle == nil {
		log.Fatal("Verify Function or Handle can't be none!\n")
		return
	}
	lines := strings.Split(string(run("blkid")), "\n")
	//log.Printf("%v, %d\n", lines, len(lines))
	for _, line := range lines {
		reg := regexp.MustCompile(`^(?P<device>.*):[ ]+UUID="(?P<uuid>[^"]*)"[ ]+TYPE="(?P<type>[^"]*)"`)

		values := reg.FindAllStringSubmatch(line, 1)
		if len(values) > 0 && len(values[0]) == 4 {
			var device = BlockDevice{}
			device.OS = "linux"
			device.Device = values[0][1]
			device.ID = values[0][2]
			device.Type = values[0][3]

			mp, err := GetMountPoint(device.ID)
			defer mp.Release()
			if err != nil {
				log.Printf("GetMountPoint(%s)\n", device.ID)
				log.Fatal(err)
				return
			}

			lockfile.TryLock(nil)
			defer lockfile.TryRelease()
			err = syscall.Mount(device.Device, mp.MountPoint, device.Type, 0, "")

			if err != nil {
				log.Printf("Mount(\"%s\", \"%s\", \"%s\", 0, \"\")\n", device.Device, mp.MountPoint, device.Type)
				log.Fatal(err)
				return
			}

			blk := BlkInfo{
				Device: &device,
				Mount:  mp}
			if verifyFunc(blk) {
				handle(blk)
			}

			err = syscall.Unmount(mp.MountPoint, syscall.MNT_DETACH)
			if err != nil {
				log.Printf("Unmount(\"%s\", \"%v\")\n", device.Device, syscall.MNT_DETACH)
				log.Fatal(err)
			}
		}
	}
}
