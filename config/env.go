package config

import helper "shopcore/pkg/helpers"

var (
	APP_PORT     = helper.GetEnv("APP_PORT", "3333")
	MONGO_DB_URI = helper.GetEnv("MONGO_DB_URI", "mongodb://username:password@localhost:9041/")
)
