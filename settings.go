package main

const (
	MongoDBHost    = "localhost"
	MongoDb        = "canaryhaiku"
	HmacSecret     = "secret"
	SessionsSecret = "secret"
	SessionName    = "haikuSession"
	Port           = ":3000"
)

var hmacSampleSecret = []byte(HmacSecret)
