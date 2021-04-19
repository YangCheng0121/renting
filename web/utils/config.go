package utils

import (
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"

	//使用了beego框架的配置文件读取模块
)

var (
	G_server_name  string //项目名称
	G_server_addr  string //服务器ip地址
	G_server_port  string //服务器端口
	G_redis_addr   string //redis ip地址
	G_redis_port   string //redis port端口
	G_redis_dbnum  string //redis db 编号
	G_mysql_addr   string //mysql ip 地址
	G_mysql_port   string //mysql 端口
	G_mysql_dbname string //mysql db name
	G_fastdfs_port   string //fastdfs 端口
	G_fastdfs_addr string //fastdfs ip
)

func InitConfig() {
	//从配置文件读取配置信息
	appconf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		logs.Debug(err)
		return
	}
	G_server_name, _ = appconf.String("appname")
	G_server_addr, _ = appconf.String("httpaddr")
	G_server_port, _ = appconf.String("httpport")
	G_redis_addr, _ = appconf.String("redisaddr")
	G_redis_port, _ = appconf.String("redisport")
	G_redis_dbnum, _ = appconf.String("redisdbnum")
	G_mysql_addr, _ = appconf.String("mysqladdr")
	G_mysql_port, _ = appconf.String("mysqlport")
	G_mysql_dbname, _ = appconf.String("mysqldbname")
	G_fastdfs_port, _ = appconf.String("fastdfsport")
	G_fastdfs_addr, _ = appconf.String("fastdfsaddr")
	return
}

func init() {
	InitConfig()
}