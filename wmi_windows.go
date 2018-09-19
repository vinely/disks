package disk

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/StackExchange/wmi"
)

type Win32_DiskDrive struct {
	Index         uint32
	InterfaceType string
	Model         string
	Size          uint64
}

// type Win32_LogicalDiskToPartition struct {
// 	EndingAddress   uint64
// 	StartingAddress uint64
// 	Antecedent      Win32_DiskPartition
// 	Dependent       Win32_LogicalDisk
// }

type Win32_DiskPartition struct {
	Index     uint32
	DiskIndex uint32
	DeviceID  string
	Name      string
	Type      string
	Size      uint64
}

type Win32_LogicalDisk struct {
	DeviceID   string
	VolumeName string
	DriveType  uint32
	FileSystem string
	Name       string
	Size       uint64
	FreeSpace  uint64
}

type Win32_Partition struct {
	Win32_DiskPartition
	Logical *Win32_Logical
}

type Win32_Disk struct {
	Win32_DiskDrive
	Partition PartMap
}
type Win32_Logical struct {
	Win32_LogicalDisk
	AssignedPart *Win32_Partition
}

type Disks map[uint32]Win32_Disk
type PartMap map[uint32]Win32_Partition
type LogicalMap map[string]Win32_Logical
type LetterDiskPartMap map[string]LetterToDiskPartition

type LetterToDiskPartition struct {
	Letter       string
	Isassociated bool
	DiskIndex    uint32
	PartIndex    uint32
}

func GetDisks() Disks {
	disks := GetDiskDrive()
	diskList := make(Disks)
	for i, _ := range disks {
		diskList[disks[i].Index] = Win32_Disk{Win32_DiskDrive: disks[i], Partition: make(PartMap)}
	}
	parts := GetPartition()
	for _, p := range parts {
		if _, ok := diskList[p.DiskIndex]; !ok {
			continue
		} else {
			diskList[p.DiskIndex].Partition[p.Index] = Win32_Partition{Win32_DiskPartition: p}
		}
	}
	return diskList
}

func GetDisksAssociated() Disks {
	disks := GetDisks()
	ldmap := GetLogicalMap()
	disks.AssociateToLogical(ldmap)
	return disks
}

func (disk Disks) AssociateToLogical(mp map[string]Win32_Logical) {
	AssociateDiskToLogical(disk, mp)
}

func GetLogicalMap() LogicalMap {
	mp := make(LogicalMap)
	logicals := GetLogicalDisk()
	for i, _ := range logicals {
		symbol := string([]byte(logicals[i].DeviceID)[0:1])
		mp[symbol] = Win32_Logical{Win32_LogicalDisk: logicals[i]}
	}
	return mp
}
func GetLogicalMapAssociated() LogicalMap {
	disks := GetDisks()
	mp := GetLogicalMap()
	mp.AssociateToDisks(disks)
	return mp
}

func (mp LogicalMap) AssociateToDisks(disk Disks) {
	AssociateDiskToLogical(disk, mp)
}

func GetLetterToDiskPartition() LetterDiskPartMap {
	cmdStr := fmt.Sprintf(`wmic logicaldisk assoc /assocclass:Win32_LogicalDiskToPartition`)
	cmd := exec.Command("cmd", "/C", cmdStr)
	bytes, err := cmd.Output()
	if err != nil {
		log.Println(err.Error())
	}
	lines := strings.Split(string(bytes), "\r\n")

	var found *LetterToDiskPartition
	var tmpu64 uint64
	found = nil

	maplist := make(map[string]LetterToDiskPartition)
	ldreg := regexp.MustCompile(`Win32_LogicalDisk.DeviceID="(?P<b>[A-Za-z]):"`)
	for _, line := range lines {
		if found != nil {
			if ldstrs := ldreg.FindAllStringSubmatch(line, 1); len(ldstrs) != 0 {
				found = &LetterToDiskPartition{Letter: ldstrs[0][1]}
				continue
			}
			reg := regexp.MustCompile(`Win32_DiskPartition.DeviceID="Disk #(?P<b>[\d]*), Partition #(?P<c>[\d]*)"`)
			strs := reg.FindAllStringSubmatch(line, 1)
			if len(strs) != 0 {

				tmpu64, err = strconv.ParseUint(strs[0][1], 10, 32)
				if err != nil {
					found.Isassociated = false
					continue
				}
				found.DiskIndex = uint32(tmpu64)
				tmpu64, err = strconv.ParseUint(strs[0][2], 10, 32)
				if err != nil {
					found.Isassociated = false
					continue
				}
				found.PartIndex = uint32(tmpu64)
				found.Isassociated = true
				//fmt.Println(strs[0][1] + "   " + strs[0][2])
			} else {
				found.Isassociated = false
				found.DiskIndex = 0
				found.PartIndex = 0
				//fmt.Println("null")
			}
			maplist[found.Letter] = *found
			found = nil
		} else {
			ldstrs := ldreg.FindAllStringSubmatch(line, 1)
			if len(ldstrs) != 0 {
				found = &LetterToDiskPartition{Letter: ldstrs[0][1]}
				//fmt.Print(ldstrs[0][1] + "  -  ")
			}
		}
	}
	return maplist
}

func AssociateDiskToLogical(disk Disks, mp map[string]Win32_Logical) {
	ml := GetLetterToDiskPartition()
	for k, v := range ml {
		if !v.Isassociated {
			continue
		} else {
			d := v.DiskIndex
			p := v.PartIndex
			wp := disk[d].Partition[p]
			wl := mp[k]
			wl.AssignedPart = &wp
			mp[k] = wl
			wp.Logical = &wl
			disk[d].Partition[p] = wp
		}
	}

}

func GetDiskDrive() []Win32_DiskDrive {
	var dst []Win32_DiskDrive
	err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst)
	if err != nil {
		log.Fatalf("getpartition: %s", err)
	}
	return dst
}

func GetPartition() []Win32_DiskPartition {
	var dst []Win32_DiskPartition
	err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst)
	if err != nil {
		log.Fatalf("getpartition: %s", err)
	}
	return dst
}

func GetLogicalDisk() []Win32_LogicalDisk {
	var dst []Win32_LogicalDisk
	err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst)
	if err != nil {
		log.Fatalf("getpartition: %s", err)
	}
	return dst
}

func FoundDiskNoFromMountPoint(letter string) (uint32, bool) {
	l := []byte(letter)
	lmp := GetLetterToDiskPartition()
	if ldp, ok := lmp[string(l[0:1])]; ok {
		return ldp.DiskIndex, true
	}
	return 0, false
}
