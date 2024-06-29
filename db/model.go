package db

import "gorm.io/gorm"

type Router struct {
	gorm.Model
	UserId    uint
	NetworkId uint
	VmName    string
}
type Network struct {
	gorm.Model
	NetName string
	Dev     string
	UserId  uint
}
