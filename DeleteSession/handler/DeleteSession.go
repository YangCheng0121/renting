package handler

import (
	"context"
	"encoding/json"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/core/logs"
	log "github.com/micro/micro/v3/service/logger"
	pb "renting/DeleteSession/proto"
	"renting/web/utils"
)

type DeleteSession struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *DeleteSession) DeleteSession(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- DELETE session    /api/v1.0/session ----------------")

	// 创建返回空间
	// 初始化的是否返回不存在
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 获取连接缓存的数据
	redisConfigMap := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	logs.Info(redisConfigMap)
	redisConfig, _ := json.Marshal(redisConfigMap)

	//连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logs.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	sessionIdName := req.Sessionid + "name"
	sessionIdUserId := req.Sessionid + "user_id"
	sessionIdMobile := req.Sessionid + "mobile"

	// 从缓存中获取session 那么使用唯一识别码
	_ = bm.Delete(context.TODO(), sessionIdName)
	_ = bm.Delete(context.TODO(), sessionIdUserId)
	_ = bm.Delete(context.TODO(), sessionIdMobile)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *DeleteSession) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.DeleteSession_StreamStream) error {
	log.Infof("Received DeleteSession.Stream request with count: %d", req.Count)

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
func (e *DeleteSession) PingPong(ctx context.Context, stream pb.DeleteSession_PingPongStream) error {
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
