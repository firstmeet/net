package qemu

import (
	"errors"
	"fmt"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"libvirt.org/go/libvirt"
)

type NetworkInterface interface {
	UpdateNetwork() error
}
type Network struct {
	NetName string
	Dev     string
}
type Openwrt struct {
	VmName  string
	Source  string
	NetName string
	NetNum  int
}

func NewNetwork(netName, dev string) *Network {
	return &Network{
		NetName: netName,
		Dev:     dev,
	}
}
func NewOpenwrt(vmName, source string, netNum int) *Openwrt {
	return &Openwrt{
		VmName: vmName,
		Source: source,
		NetNum: netNum,
	}
}
func (n *Network) CreateNetwork() error {
	network := libvirtxml.Network{
		Name: n.NetName,
		Forward: &libvirtxml.NetworkForward{
			Mode: "bridge",
			Dev:  n.Dev,
		},
	}
	netXml, err := network.Marshal()
	if err != nil {
		return err
	}
	fmt.Println(netXml)
	_, err = Conn.NetworkDefineXML(netXml)
	if err != nil {
		return errors.New("NetworkDefineXML failed:" + err.Error())
	}
	net, err := Conn.LookupNetworkByName(n.NetName)
	if err != nil {
		return err
	}
	err = net.Create()
	if err != nil {
		return err
	}
	//set autostart
	err = net.SetAutostart(true)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) UpdateNetwork() error {
	net, err := Conn.LookupNetworkByName(n.NetName)
	if err != nil {
		return err
	}
	desc, err := net.GetXMLDesc(0)
	if err != nil {
		return err
	}
	domainInterface := libvirtxml.Network{}
	err = domainInterface.Unmarshal(desc)
	if err != nil {
		return err
	}
	domainInterface.Bridge.Name = n.Dev
	netXml, err := domainInterface.Marshal()
	if err != nil {
		return err
	}
	err = net.Update(libvirt.NETWORK_UPDATE_COMMAND_MODIFY, libvirt.NETWORK_SECTION_BRIDGE, -1, netXml, libvirt.NETWORK_UPDATE_AFFECT_LIVE)
	if err != nil {
		return err
	}
	//update config
	err = net.Update(libvirt.NETWORK_UPDATE_COMMAND_MODIFY, libvirt.NETWORK_SECTION_BRIDGE, -1, netXml, libvirt.NETWORK_UPDATE_AFFECT_CONFIG)
	if err != nil {
		return err
	}

	return nil
}

func (o *Openwrt) UpdateNetwork() error {
	dom, err := Conn.LookupDomainByName(o.VmName)
	if err != nil {
		return err
	}
	desc, err := dom.GetXMLDesc(0)
	if err != nil {
		return err
	}
	domainInterface := libvirtxml.Domain{}
	err = domainInterface.Unmarshal(desc)
	if err != nil {
		return err
	}
	//get all network interface
	interfaces := domainInterface.Devices.Interfaces
	if len(interfaces) == 0 {
		return errors.New("no network interface")
	}
	if o.NetNum > len(interfaces) {
		return errors.New("NetNum is out of range")
	}
	//first detach the network interface
	netXml, err := domainInterface.Devices.Interfaces[o.NetNum-1].Marshal()
	if err != nil {
		return err
	}
	err = dom.DetachDeviceFlags(netXml, libvirt.DOMAIN_DEVICE_MODIFY_LIVE)
	if err != nil {
		return err
	}
	err = dom.DetachDeviceFlags(netXml, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG)
	if err != nil {
		return err
	}
	interfaces[o.NetNum-1].Source.Network = &libvirtxml.DomainInterfaceSourceNetwork{
		Network: o.Source,
	}
	netXml, err = interfaces[o.NetNum-1].Marshal()
	if err != nil {
		return err
	}
	err = dom.AttachDeviceFlags(netXml, libvirt.DOMAIN_DEVICE_MODIFY_LIVE)
	if err != nil {
		return err
	}
	err = dom.AttachDeviceFlags(netXml, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG)
	if err != nil {
		return err
	}
	return nil
}
