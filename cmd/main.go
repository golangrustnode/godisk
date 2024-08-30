package main

import (
	"github.com/golangrustnode/godisk/linuxdisk"
	"github.com/golangrustnode/log"
)

func main() {
	res, err := linuxdisk.ParseFstab()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(res)
	boot, err := linuxdisk.GetBootDisk()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(boot)
	devices, err := linuxdisk.ListBlockDevices()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(devices)
}
