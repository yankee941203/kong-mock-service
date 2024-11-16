package model

import "gorm.io/gorm"

type ServiceInfo struct {
	ID               uint     `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	AvaliableVersion []string `json:"available_version"`
}

type ServiceEntity struct {
	gorm.Model
	Name        string
	Description string
	Images      []*ImageEntity `gorm:"many2many:service_images;"`
}

type ImageEntity struct {
	gorm.Model
	Name     string
	Version  string
	Services []*ServiceEntity `gorm:"many2many:service_images;"`
}
