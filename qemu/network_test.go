package qemu

import "testing"

func TestNetwork_UpdateNetwork(t *testing.T) {
	network := NewNetwork("br2_xxxx_network", "enp4s0f0")
	err := network.UpdateNetwork()
	if err != nil {
		t.Errorf("UpdateNetwork() failed, got %v", err)
	}
}
func TestNewNetworkCreate(t *testing.T) {
	network := NewNetwork("br2_xxxx_network1", "eth0")
	if network.NetName != "br2_xxxx_network" || network.Dev != "eth0" {
		t.Errorf("NewNetwork() failed, got %v", network)
	}
	NewLibvirt()
	err := network.CreateNetwork()
	if err != nil {
		t.Errorf("CreateNetwork() failed, got %v", err)
	}
}
func TestNewOpenwrt(t *testing.T) {
	openwrt := NewOpenwrt("ubuntu24.04", "default", 1)
	if openwrt.VmName != "ubuntu24.04" || openwrt.Source != "default" || openwrt.NetNum != 1 {
		t.Errorf("NewOpenwrt() failed, got %v", openwrt)
	}
}
func TestOpenwrt_UpdateNetwork(t *testing.T) {
	openwrt := NewOpenwrt("ubuntu24.04", "default", 1)
	NewLibvirt()
	err := openwrt.UpdateNetwork()
	if err != nil {
		t.Errorf("UpdateNetwork() failed, got %v", err)
	}
}
