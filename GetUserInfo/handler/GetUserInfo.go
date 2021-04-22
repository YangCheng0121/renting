package handler

import (
	"context"
	"encoding/json"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"reflect"
	"renting/web/models"
	"renting/web/utils"

	pb "renting/GetUserInfo/proto"
)

type GetUserInfo struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetUserInfo) GetUserInfo(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- GET  /api/v1.0/user GetUserInfo() ------------------")
	// 打印sessionId
	logs.Info(req.Sessionid, reflect.TypeOf(req.Sessionid))
	// 错误码
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 构建连接缓存的数据
	redisConfigMap := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	logs.Info(redisConfigMap)

	redisConfig, _ := json.Marshal(redisConfigMap)

	// 连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logs.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 拼接用户信息缓存字段
	sessionIdUserId := req.Sessionid + "user_id"

	// 获取到当前登录用户的user_id
	valueId, _ := bm.Get(context.TODO(), sessionIdUserId)

	// 数据格式转换
	id := int(valueId.([]uint8)[0])
	logs.Info(id, reflect.TypeOf(id))

	// 创建user表
	user := models.User{Id: id}
	// 创建数据库orm句柄
	o := orm.NewOrm()
	// 查询表
	err = o.Read(&user)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 将查询到的数据依次赋值
	rsp.UserId = int64(user.Id)
	rsp.Name = user.Name
	rsp.Mobile = user.Mobile
	rsp.RealName = user.Real_name
	rsp.IdCard = user.Id_card

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetUserInfo) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetUserInfo_StreamStream) error {
	log.Infof("Received GetUserInfo.Stream request with count: %d", req.Count)

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
func (e *GetUserInfo) PingPong(ctx context.Context, stream pb.GetUserInfo_PingPongStream) error {
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
