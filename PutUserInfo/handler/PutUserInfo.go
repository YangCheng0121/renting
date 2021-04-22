package handler

import (
	"context"
	"encoding/json"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"reflect"
	"renting/web/models"
	"renting/web/utils"
	"time"

	log "github.com/micro/micro/v3/service/logger"
	pb "renting/PutUserInfo/proto"
)

type PutUserInfo struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PutUserInfo) PutUserInfo(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	// 打印被调用的函数
	logs.Info("---------------- PUT  /api/v1.0/user/name PutUserInfo() ------------------")

	// 创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/* 得到用户发送过来的name */
	logs.Info(rsp.Username)

	/* 从从sessionId 获取当前的 userId */
	// 连接redis
	redisConfigMap := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	logs.Info(redisConfigMap)
	redisConfig, _ := json.Marshal(redisConfigMap)
	logs.Info(string(redisConfig))

	// 连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logs.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 拼接key
	sessionIdUserId := req.Sessionid + "user_id"
	// 获取userId
	valueId, _ := bm.Get(context.TODO(), sessionIdUserId)
	logs.Info(valueId, reflect.TypeOf(valueId))

	id := int(valueId.([]uint8)[0])
	logs.Info(id, reflect.TypeOf(id))

	// 创建对象
	user := models.User{Id: id, Name: req.Username}
	/* 更新对应 user_id 的 name 字段的内容 */
	// 创建数据库句柄
	o := orm.NewOrm()
	// 更新
	_, err = o.Update(&user, "name")
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/* 更新session user_id */
	sessionIdName := req.Sessionid + "name"
	_ = bm.Put(context.TODO(), sessionIdUserId, string(user.Id), time.Second*600)
	/* 更新session name */
	_ = bm.Put(context.TODO(), sessionIdName, string(user.Name), time.Second*600)

	/* 成功返回数据 */
	rsp.Username = user.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PutUserInfo) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.PutUserInfo_StreamStream) error {
	log.Infof("Received PutUserInfo.Stream request with count: %d", req.Count)

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
func (e *PutUserInfo) PingPong(ctx context.Context, stream pb.PutUserInfo_PingPongStream) error {
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
