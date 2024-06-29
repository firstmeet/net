package services

import "awesomeProject/db"

// get router from userid and routerid
func GetRouter(userId uint, routerId uint) (db.Router, error) {
	var router db.Router
	err := db.DB.Where("user_id = ? AND router_id = ?", userId, routerId).First(&router).Error
	return router, err
}

// get network from userid and networkid
func GetNetwork(userId uint, networkId uint) (db.Network, error) {
	var network db.Network
	err := db.DB.Where("user_id = ? AND network_id = ?", userId, networkId).First(&network).Error
	return network, err
}
