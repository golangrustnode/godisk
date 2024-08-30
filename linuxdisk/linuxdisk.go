package linuxdisk

import (
	"github.com/golangrustnode/gotrimreader/trimreader"
	"github.com/golangrustnode/log"
	"os"
)

type PhysicalDisk struct {
	SectorNum  uint64
	SectorSize uint64
	Rotate     uint64
	IsBootDisk bool
}

func GetPhysicalDiskInfo() {

}

func ListBlockDevices() (map[string]PhysicalDisk, error) {
	dir_path := "/sys/class/block/"
	f, err := os.Open(dir_path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	files, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}
	var block_devices = make(map[string]PhysicalDisk)
	for _, file := range files {
		if !IsPhysicalBlockDeivce(dir_path + file.Name() + "/device") {
			continue
		}
		pdisk, err := GetBlockInfo(dir_path + file.Name())
		if err != nil {
			log.Error(err)
			continue
		}
		pdisk.IsBootDisk = IsBootDisk(file.Name())
		block_devices[file.Name()] = pdisk
	}
	return block_devices, nil
}

func IsPhysicalBlockDeivce(device_path string) bool {
	log.Debug("Checking if ", device_path, " is a physical block device")
	_, err := os.Stat(device_path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func GetBlockInfo(device_path string) (PhysicalDisk, error) {
	log.Debug("Getting block info for ", device_path)
	var disk PhysicalDisk
	var err error
	disk.IsBootDisk = false
	disk.SectorNum, err = trimreader.ReadAsUint64(device_path + "/size")
	if err != nil {
		log.Error("Error reading size from ", device_path, " Error: ", err)
	}
	disk.SectorSize, err = trimreader.ReadAsUint64(device_path + "/queue/hw_sector_size")
	if err != nil {
		log.Error("Error reading hw_sector_size from ", device_path, " Error: ", err)
	}
	disk.Rotate, err = trimreader.ReadAsUint64(device_path + "/queue/rotational")
	if err != nil {
		log.Error("Error reading rotational from ", device_path, " Error: ", err)
		return disk, err
	}
	return disk, nil
}
