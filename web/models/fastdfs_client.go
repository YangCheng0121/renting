package models

import (
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

const (
	Prefix = "http://"
	Group  = "/group1/upload"
)

//domain:http: //127.0.0.1:3666 md5:6672e1cf7c2069422c97de9ef0e4d54b mtime:1.61898525e+09 path:/group1/default/20210421/14/07/5/d5025964.jpg retcode:0 retmsg: scene:default scenes:default size:51175 src:/group1/default/20210421/14/07/5/d5025964.jpg url:http://127.0.0.1:3666/group1/default/20210421/14/07/5/d5025964.jpg?name=d5025964.jpg&download=1

type File struct {
	Domain  string
	MD5     string
	Mtime   int
	Path    string
	RetCode int
	RetMsg  string
	Scene   string
	Scenes  string
	Size    int
	Src     string
	Url     string
}

// 通过文件名的方式进行上传
func UploadByFileName(fileExtName string) (file File, err error) {
	fastDfsConf, err := config.NewConfig("ini", "./conf/fastdfs.conf")
	if err != nil {
		logs.Error(err)
		return
	}

	fastDfsHost, _ := fastDfsConf.String("fast_dfs_host")
	fastDfsPort, _ := fastDfsConf.String("fast_dfs_port")

	url := Prefix + fastDfsHost + ":" + fastDfsPort + Group

	req := httplib.Post(url)
	req.PostFile("file", fileExtName) //注意不是全路径
	req.Param("output", "json")
	req.Param("scene", "")
	req.Param("path", "")
	_ = req.ToJSON(&file)

	return file, nil
}
