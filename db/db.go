package db

import (
	"errors"
	"fmt"
	"github.com/Altoros/cf-broker-boilerplate/model"
	"github.com/jinzhu/gorm"
	"os"
)

func New() (*gorm.DB, error) {

	var databaseUrl string
	if os.Getenv("DATABASE_URL") != "" {
		databaseUrl = os.Getenv("DATABASE_URL")
	} else {
		creds, err := LoadServiceCredentials("p-mysql")
		if err != nil {
			return errors.New("Can not load Mysql DB configuration from DATABASE_URL or VCAP_SERVICES env variables")
		}
		databaseUrl = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			creds.GetUsername(),
			creds.GetPassword(),
			creds.GetHost(),
			creds.GetPort(),
			creds.GetDBName())
	}

	db, err := gorm.Open("mysql", databaseUrl)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.ServiceInstance{}, &model.ServiceBinding{}).Error

	return db, err
}
