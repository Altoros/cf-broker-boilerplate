package model

import (
	"github.com/jinzhu/gorm"
)

type ServiceInstance struct {
	gorm.Model
	InstanceId string `gorm:"not null"`
}

type ServiceBinding struct {
	gorm.Model
	BindingId  string `gorm:"not null"`
	InstanceId string `gorm:"not null"`
}
