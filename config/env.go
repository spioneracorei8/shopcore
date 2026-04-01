package config

import helper "shopcore/pkg/helpers"

var (
	APP_PORT     = helper.GetEnv("APP_PORT", "")
	MONGO_DB_URI = helper.GetEnv("MONGO_DB_URI", "")
)
