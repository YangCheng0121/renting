package handler

import (
	"context"
	"encoding/json"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/core/logs"
	"github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"reflect"
	"renting/web/utils"

	pb "renting/GetSession/proto"
)

type GetSession struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetSession) GetSession(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- GET session    /api/v1.0/session ----------------")
	// 创建返回空间
	// 初始化的是否返回不存在
	rsp.Errno = utils.RECODE_SESSIONERR
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 获取前端的cookie
	logs.Info(req.Sessionid, reflect.TypeOf(req.Sessionid))
	redisConfigMap := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	logs.Info(redisConfigMap)
	redisConfig, _ := json.Marshal(redisConfigMap)
	logs.Info(string(redisConfig))

	//连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logs.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 拼接key
	sessionIdName := req.Sessionid + "name"
	// 从缓存中获取 session 那么使用唯一识别码 通过 key 查询用户名
	areasInfoValue, _ := bm.Get(context.TODO(), sessionIdName)
	// 查看返回数据类型
	logs.Info(reflect.TypeOf(areasInfoValue), areasInfoValue)

	// 通过redis方法进行转换
	name, err := redis.String(areasInfoValue, nil)
	if err != nil {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	// 查看返回数据类型
	logs.Info(name, reflect.TypeOf(name))

	// 获取到了session
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	rsp.Data = name

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetSession) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetSession_StreamStream) error {
	log.Infof("Received GetSession.Stream request with count: %d", req.Count)

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
func (e *GetSession) PingPong(ctx context.Context, stream pb.GetSession_PingPongStream) error {
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
