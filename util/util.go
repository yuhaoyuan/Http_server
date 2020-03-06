package util

import (
	"context"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/yuhaoyuan/Http_server/config"
	"log"
)

// 七牛云CDN上传
func UploadQiniu(filePath, fileName string) string {
	bucket := config.BaseConf.Bucket
	accessKey:= config.BaseConf.AccessKEY
	secretKey:= config.BaseConf.SecretKEY
	putPolicy := storage.PutPolicy{
		Scope:               bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, fileName, filePath, &putExtra)
	if err != nil {
		log.Println(err)
		return ""
	}
	return config.BaseConf.CDNUrl + fileName
}