package settings

type CloudFlareConfig struct {
	BucketName      string
	AccountId       string
	AccessKeyId     string
	AccessKeySecret string
	DownloadFileTTL int    // 预签名下载链接的有效期，单位秒
	UploadFileTTL 	int    // 预签名上传链接的有效期，单位秒
	DomainAccess    string // 用于替换掉原来的域名，比如：https://cdn.cl.tv。这个域名的信息，需要在 R2 的具体一个桶里面的 Settings 去设置
}
