package linuxdisk

import (
	"bufio"
	"fmt"
	"github.com/golangrustnode/log"

	"os"
	"strings"
)

type FstabEntry struct {
	Device     string
	MountPoint string
	FsType     string
	Options    string
	Dump       int
	Pass       int
}

func ParseFstab() ([]FstabEntry, error) {
	tabpath := "/etc/fstab"
	return parseFstab(tabpath)
}

func parseFstab(path string) ([]FstabEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []FstabEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 忽略空行和注释行
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// 按空格或制表符分割行
		fields := strings.Fields(line)

		// 每行应有至少6个字段
		if len(fields) < 6 {
			continue
		}

		// 解析 dump 和 pass 字段
		var dump, pass int
		fmt.Sscanf(fields[4], "%d", &dump)
		fmt.Sscanf(fields[5], "%d", &pass)

		entry := FstabEntry{
			Device:     fields[0],
			MountPoint: fields[1],
			FsType:     fields[2],
			Options:    fields[3],
			Dump:       dump,
			Pass:       pass,
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func GetBootDisk() (string, error) {
	entries, err := ParseFstab()
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if strings.TrimSpace(entry.MountPoint) != "/boot" {
			continue
		}
		target, err := os.Readlink(entry.Device)
		if err != nil {
			log.Error(err)
			return "", err
		}
		return target, nil
	}

	return "", fmt.Errorf("no boot disk found")
}

func IsBootDisk(diskname string) bool {
	bootdisk, err := GetBootDisk()
	if err != nil {
		log.Error(err)
		return false
	}
	return strings.Contains(bootdisk, diskname)
}
