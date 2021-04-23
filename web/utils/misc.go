package utils

/* 将url加上 http://IP:PROT/  前缀 */
//http:// + 127.0.0.1 + ：+ 8080 + 请求
func AddDomain2Url(url string) (domainUrl string) {
	domainUrl = "http://" + G_fastdfs_addr + ":" + G_fastdfs_port + url
	return domainUrl
}
