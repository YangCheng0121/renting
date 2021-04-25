package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	pb "renting/PostAvatar/proto"
	"renting/web/models"
	"renting/web/utils"
)

type PostAvatar struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostAvatar) PostAvatar(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- 上传用户头像 PostAvatar /api/v1.0/user/avatar ----------------")

	// 初始化返回正确的返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 检查下是否正常
	logs.Info(len(req.Avatar), req.Filesize)

	/* 获取文件后缀名 */ //dsnlkjfajadskfksda.sadsdasd.sdasd.jpg
	logs.Info("后缀名", path.Ext(req.Filename))

	fileType := "other"
	/* 存储文件到 fastdfs 当中并且获取 url */
	//.jpg
	fileExt := path.Ext(req.Filename)
	logs.Info(fileExt)

	if fileExt == ".jpg" || fileExt == ".png" || fileExt == ".gif" || fileExt == ".jpeg" {
		fileType = "img"
	}

	if fileType != "img" {
		return errors.New("不支持文件上传")
	}

	err := ioutil.WriteFile(req.Filename, req.Avatar, 0666)
	if err != nil {
		logs.Info("文件读写错误", err)
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// group1 group1/M00/00/00/wKgLg1t08pmANXH1AAaInSze-cQ589.jpg
	// 上传数据
	avatar, err := models.UploadByFileName(req.Filename)
	if err != nil {
		return err
	}
	logs.Info("avatar:", avatar)
	defer os.Remove(req.Filename)

	/* 通过session 获取我们当前现在用户的user_id */
	redisConfigMap := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	logs.Info(redisConfigMap)

	redisConfig, _ := json.Marshal(redisConfigMap)
	logs.Info(string(redisConfig))

	// 连接redis 数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logs.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 拼接key
	sessionIdUserId := req.Sessionid + "user_id"

	// 获得当前用户的userId
	valueId, _ := bm.Get(context.TODO(), sessionIdUserId)
	if valueId == nil {
		logs.Info("获取登录缓存失败", err)
		rsp.Errno = utils.RECODE_SESSIONERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	id := int(valueId.([]uint8)[0])
	logs.Info(id, reflect.TypeOf(id))

	// 创建表对象
	user := models.User{Id: id, AvatarUrl: utils.AddDomain2Url(avatar.Path)}
	/* 将当前 fastdfs-url 存储到我们当前用户的表中 */
	o := orm.NewOrm()
	// 将图片的地址存入表中
	_, err = o.Update(&user, "avatar_url")
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	// 回传图片地址
	rsp.AvatarUrl = utils.AddDomain2Url(avatar.Path)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostAvatar) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.PostAvatar_StreamStream) error {
	log.Infof("Received PostAvatar.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&pb.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *PostAvatar) PingPong(ctx context.Context, stream pb.PostAvatar_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&pb.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
