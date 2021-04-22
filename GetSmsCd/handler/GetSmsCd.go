package handler

import (
	"context"
	"encoding/json"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"math/rand"
	"reflect"
	"renting/web/models"
	"renting/web/utils"
	"time"

	pb "renting/GetSmsCd/proto"
)

type GetSmsCd struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetSmsCd) GetSmsCd(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- GET smscd  api/v1.0/smscode/:id ----------------")

	// 初始化返回正确的返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/* 验证uuid的缓存 */
	// 验证手机号
	o := orm.NewOrm()
	user := models.User{Mobile: req.Mobile}
	err := o.Read(&user)

	if err == nil {
		logs.Info("用户已存在")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info(err)

	// 连接redis数据库
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

	logs.Info(req.Id, reflect.TypeOf(req.Id))

	// 查询相关数据
	value, _ := bm.Get(context.TODO(), req.Id)
	if value == nil {
		logs.Info("获取到缓存数据查询失败", value)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	logs.Info(value, reflect.TypeOf(value))

	valueStr, _ := redis.String(value, nil)
	logs.Info(valueStr, reflect.TypeOf(valueStr))

	// 数据对比
	if req.Text != valueStr {
		logs.Info("图片验证码 错误")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	size := r.Intn(8999) + 10000
	logs.Info(size)

	err = bm.Put(context.TODO(), req.Mobile, size, time.Second*300)
	if err != nil {
		logs.Info("缓存出现问题")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetSmsCd) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetSmsCd_StreamStream) error {
	log.Infof("Received GetSmsCd.Stream request with count: %d", req.Count)

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
func (e *GetSmsCd) PingPong(ctx context.Context, stream pb.GetSmsCd_PingPongStream) error {
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
