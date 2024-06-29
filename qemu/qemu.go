package qemu

import "libvirt.org/go/libvirt"

var Conn *libvirt.Connect

func NewLibvirt() {
	var err error
	Conn, err = libvirt.NewConnect("qemu:///system")
	if err != nil {
		panic(err)
	}
}
