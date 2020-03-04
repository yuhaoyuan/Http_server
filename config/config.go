package config

// BaseConf 参数
var BaseConf = BaseConfig{}

// BaseConfig config 结构体
type BaseConfig struct {
	Addr      string
	AccessKEY string
	SecretKEY string
	Bucket    string
	CdnUrl string
	LogName string
}

// BaseConfInit 初始化环境变量
func BaseConfInit() {
	//BaseConf.Addr = os.Getenv("ADDR")
	//BaseConf.AccessKEY = os.Getenv("ACCESSKEY")
	//BaseConf.SecretKEY = os.Getenv("SECRETKEY")
	//BaseConf.Bucket = os.Getenv("BUCKET")
	//BaseConf.CdnUrl = os.Getenv("CDNURL")
	//BaseConf.LogName = os.Getenv("HSLOGNAME")
	BaseConf.Addr = "127.0.0.1:8009"
	BaseConf.AccessKEY = "4w-q4XHzGx3eV_a1aMPog9hu44MLVpoLWrNv8rGH"
	BaseConf.SecretKEY = "qAiFkbYFUEkiaIJjZawLvfaIp2K3P-CbpzcNLkKo"
	BaseConf.Bucket = "yuhaoyuan"
	BaseConf.CdnUrl = "http://q6gy4v9f7.bkt.clouddn.com/"
	BaseConf.LogName = "yhy_http_server_log"
}