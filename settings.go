package main

const (
	MongoDBHost = "localhost"
	MongoDb     = "canaryhaiku"
	HmacSecret  = "secret"
	Port        = ":3000"
)

var hmacSampleSecret = []byte(HmacSecret)
