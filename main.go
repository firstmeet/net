package main

import (
	"awesomeProject/qemu"
	"awesomeProject/services"
	"fmt"
)

type Request struct {
	UserId   uint
	RouterId uint
	WanType  uint
	Dev      string
}

const (
	WanTypeDev     = 1
	WanTypeVirtual = 2
)

func main() {
	qemu.NewLibvirt()
	var req Request
	req.UserId = 1
	req.RouterId = 1
	req.WanType = 1
	req.Dev = "eth0"
	router, err := services.GetRouter(req.UserId, req.RouterId)
	if err != nil {
		fmt.Println(err)
	}
	switch req.WanType {
	case WanTypeDev:
		getNetwork, err := services.GetNetwork(router.UserId, router.NetworkId)
		if nil != err {
			fmt.Println(err)
			return
		}
		network := qemu.NewNetwork(getNetwork.NetName, req.Dev)
		err = network.UpdateNetwork()
		if err != nil {
			fmt.Println(err)
			return
		}
		// update openwrt network
		err = qemu.NewOpenwrt(router.VmName, network.NetName, 1).UpdateNetwork()
		if err != nil {
			fmt.Println(err)
			return
		}
	case WanTypeVirtual:
		// update openwrt network
		err = qemu.NewOpenwrt(router.VmName, req.Dev, 1).UpdateNetwork()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
