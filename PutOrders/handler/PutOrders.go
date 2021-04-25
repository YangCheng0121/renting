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
	"strconv"

	pb "renting/PutOrders/proto"
)

type PutOrders struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PutOrders) PutOrders(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- api/v1.0/orders  PutOrders post success ----------------")

	// 创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 1通过session得到当前的user_id
	//构建连接缓存的数据
	redisConfigMap := map[string]string{
		"key": utils.G_server_name,
		//"conn":"127.0.0.1:6379",
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
	sessionIdUserId := req.Sessionid + "user_id"

	valueId, _ := bm.Get(context.TODO(), sessionIdUserId)
	if valueId == nil {
		logs.Info("获取登录缓存失败", err)
		rsp.Errno = utils.RECODE_SESSIONERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	userId := int(valueId.([]uint8)[0])
	logs.Info(userId, reflect.TypeOf(userId))

	// 2通过url参数得到当前订单id
	orderId, _ := strconv.Atoi(req.Orderid)

	// 3通过客户端请求的json数据得到action参数
	logs.Info(req.Action)
	// 得到请求指令
	action := req.Action

	// 5查找订单，找到该订单并确定当前订单状态是wait_accept
	o := orm.NewOrm()
	order := models.OrderHouse{}
	err = o.QueryTable("order_house").Filter("id", orderId).Filter("status", models.ORDER_STATUS_WAIT_ACCEPT).One(&order)
	if err != nil {
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	if _, err := o.LoadRelated(&order, "House"); err != nil {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	house := order.House
	// 6检验该订单的user_id 是否是当前用户的user_id
	// 返回错误json
	if house.User.Id != userId {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = "订单用户不匹配，操作无效"
		return nil
	}
	// 7action为accept
	if action == "accept" {
		// 如果是接受订单，将订单状态变成待评价状态
		order.Status = models.ORDER_STATUS_WAIT_COMMENT
		// 8更换订单状态为status为reject
		reason := req.Action
		// 添加评论
		order.Comment = reason
		logs.Debug("action = reject!, reason is ", reason)
	}

	// 更新该数据到数据库中
	if _, err := o.Update(&order); err != nil {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PutOrders) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.PutOrders_StreamStream) error {
	log.Infof("Received PutOrders.Stream request with count: %d", req.Count)

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
func (e *PutOrders) PingPong(ctx context.Context, stream pb.PutOrders_PingPongStream) error {
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
