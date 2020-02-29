package config

import "os"

var BaseConf = BaseConfig{}

type BaseConfig struct {
	Addr      string
	AccessKEY string
	SecretKEY string
	Bucket    string
	CdnUrl string
}

func BaseConfInit() {
	BaseConf.Addr = os.Getenv("ADDR")
	BaseConf.AccessKEY = os.Getenv("ACCESSKEY")
	BaseConf.SecretKEY = os.Getenv("SECRETKEY")
	BaseConf.Bucket = os.Getenv("BUCKET")
	BaseConf.CdnUrl = os.Getenv("CDNURL")
}
