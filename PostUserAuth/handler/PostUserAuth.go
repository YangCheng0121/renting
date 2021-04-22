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
	"time"

	pb "renting/PostUserAuth/proto"
)

type PostUserAuth struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostUserAuth) PostUserAuth(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	// 打印被调用的函数
	logs.Info("---------------- 实名认证 PostUserAuth  api/v1.0/user/auth ----------------")

	// 创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/* 从session 中获取我们的 user_id */
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
	// 拼接key
	sessionIdUserId := req.Sessionid + "user_id"

	valueId, _ := bm.Get(context.TODO(), sessionIdUserId)
	logs.Info(valueId, reflect.TypeOf(valueId))

	id := int(valueId.([]uint8)[0])
	logs.Info(id, reflect.TypeOf(id))

	// 创建user对象
	user := models.User{
		Id:        id,
		Real_name: req.RealName,
		Id_card:   req.IdCard,
	}

	/* 更新user表中的 姓名 和 身份号 */
	o := orm.NewOrm()
	// 更新表
	_, err = o.Update(&user, "real_name", "id_card")
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	/* 更新我们的session中的user_id */
	_ = bm.Put(context.TODO(), sessionIdUserId, string(user.Id), time.Second*600)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostUserAuth) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.PostUserAuth_StreamStream) error {
	log.Infof("Received PostUserAuth.Stream request with count: %d", req.Count)

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
func (e *PostUserAuth) PingPong(ctx context.Context, stream pb.PostUserAuth_PingPongStream) error {
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
