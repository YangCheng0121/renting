package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"reflect"
	pb "renting/PostRet/proto"
	"renting/web/models"
	"renting/web/utils"
	"strconv"
	"time"
)

type PostRet struct{}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostRet) PostRet(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- POST userreg    /api/v1.0/users ----------------")

	// 初始化错误码
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 构建连接缓存的数据
	redisConfigMap := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	redisConfig, _ := json.Marshal(redisConfigMap)

	//连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logs.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 查询相关数据
	value, _ := bm.Get(context.TODO(), req.Mobile)
	if value == nil {
		logs.Info("获取到缓存数据查询失败", value)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	logs.Info(value, reflect.TypeOf(value))

	// 进行解码
	var info interface{}
	_ = json.Unmarshal(value.([]byte), &info)
	logs.Info(info, reflect.TypeOf(info))

	// 类型转换
	s := int(info.(float64))
	logs.Info(s, reflect.TypeOf(s))
	s1, err := strconv.Atoi(req.SmsCode)

	if s1 != s {
		logs.Info("短信验证码错误")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	user := models.User{}
	user.Name = req.Mobile
	// 密码正常情况下 md5 sha256 sm9  存入数据库的是你加密后的编码不是明文存入
	user.Password_hash = GetMd5String(req.Password)
	//user.Password_hash = req.Password
	user.Mobile = req.Mobile

	// 创建数据库剧本
	o := orm.NewOrm()
	// 插入数据库
	id, err := o.Insert(&user)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("id", id)

	// 生成sessionID 保证唯一性
	h := GetMd5String(req.Mobile + req.Password)
	// 返回给客户端session
	rsp.SessionID = h

	// 拼接key sessionid + name
	_ = bm.Put(context.TODO(), h+"name", string(user.Mobile), time.Second*3600)
	// 拼接key sessionid + user_id
	_ = bm.Put(context.TODO(), h+"user_id", string(user.Id), time.Second*3600)
	// 拼接key sessionid + mobile
	_ = bm.Put(context.TODO(), h+"mobile", string(user.Mobile), time.Second*3600)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostRet) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.PostRet_StreamStream) error {
	log.Infof("Received PostRet.Stream request with count: %d", req.Count)

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
func (e *PostRet) PingPong(ctx context.Context, stream pb.PostRet_PingPongStream) error {
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
