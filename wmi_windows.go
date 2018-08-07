package disk

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/StackExchange/wmi"
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

type DiskInfo struct {
	Disk      Win32_DiskDrive
	Installed bool
}

type Win32_DiskDrive struct {
	Index         uint32
	InterfaceType string
	Model         string
	Size          uint64
}

func FoundDiskNo(letter string) string {
	cmdStr := fmt.Sprintf(`wmic logicaldisk assoc /assocclass:Win32_LogicalDiskToPartition`)
	cmd := exec.Command("cmd", "/C", cmdStr)
	bytes, err := cmd.Output()
	if err != nil {
		log.Println(err.Error())
	}
	lines := strings.Split(string(bytes), "\r\n")
	found := false

	for _, line := range lines {
		if found == true {

			reg := regexp.MustCompile(`Win32_DiskPartition.DeviceID="Disk #(?P<b>[\d]*), Partition #(?P<c>[\d]*)"`)
			strs := reg.FindAllStringSubmatch(line, 1)
			return strs[0][1]
		}
		if strings.Contains(line, fmt.Sprintf("Win32_LogicalDisk.DeviceID=\"%s\"", letter)) {
			found = true
		}
	}
	return ""
}

func GetDiskDrive() []DiskInfo {
	var dst []Win32_DiskDrive
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		log.Fatalf("getpartition: %s", err)
	}
	info := make([]DiskInfo, len(dst))
	for index, _ := range dst {
		info[index].Disk = dst[index]
		info[index].Installed = false
	}
	//log.Printf("%+v\n", info)
	return info
}
