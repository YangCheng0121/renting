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
	"time"

	pb "renting/PostOrders/proto"
)

type PostOrders struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostOrders) PostOrders(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- api/v1.0/orders  PostOrders  发布订单 ----------------")

	// 创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 1根据session得到当前用户的 user_id
	// 构建连接缓存数据
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

	//拼接key
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

	// 2得到你用户请求的json数据并校验合法性
	// 获取得到用户请求Response数据的name
	var requestMap = make(map[string]interface{})
	err = json.Unmarshal(req.Body, &requestMap)

	if err != nil {
		rsp.Errno = utils.RECODE_REQERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info(requestMap)

	// 校验合法性
	// 用户参数做合法判断
	if requestMap["house_id"] == "" || requestMap["start_date"] == "" || requestMap["end_data"] == "" {
		rsp.Errno = utils.RECODE_REQERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 3确定end_date 在 start_date之后
	// 格式化日期时间
	startDateTime, _ := time.Parse("2006-01-02 15:04:05", requestMap["start_date"].(string))
	endDateTime, _ := time.Parse("2006-01-02 15:04:05", requestMap["end_date"].(string))

	//4得到一共入住的天数
	logs.Info(startDateTime, endDateTime)
	days := endDateTime.Sub(startDateTime).Hours()/24 + 1
	logs.Info(days)

	// 5根据order_id 得到关联的房源信息
	houseId, _ := strconv.Atoi(requestMap["house_id"].(string))

	// 房屋对象
	house := models.House{Id: houseId}
	o := orm.NewOrm()
	if err := o.Read(&house); err != nil {
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	_, _ = o.LoadRelated(&house, "user")

	// 6确保当前的user_id不是房源信息所关联的user_id
	if userId == house.User.Id {
		rsp.Errno = utils.RECODE_ROLEERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 7确保用户选择的房租未预定，日期没有冲突
	if endDateTime.Before(startDateTime) {
		rsp.Errno = utils.RECODE_ROLEERR
		rsp.Errmsg = "结束时间在开始时间之前"
		return nil
	}

	// 7.1添加征信步骤
	// 8封装order订单
	amount := days * float64(house.Price)
	order := models.OrderHouse{}
	order.House = &house
	user := models.User{Id: userId}
	order.User = &user
	order.BeginDate = startDateTime
	order.EndDate = endDateTime
	order.Days = int(days)
	order.HousePrice = house.Price
	order.Amount = int(amount)
	order.Status = models.ORDER_STATUS_WAIT_ACCEPT
	// 征信
	order.Credit = false

	logs.Info(order)
	// 9将订单信息入库表中
	if _, err := o.Insert(&order); err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 10返回order_id
	_ = bm.Put(context.TODO(), sessionIdUserId, string(userId), time.Second*7200)
	rsp.OrderId = int64(order.Id)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostOrders) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.PostOrders_StreamStream) error {
	log.Infof("Received PostOrders.Stream request with count: %d", req.Count)

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
func (e *PostOrders) PingPong(ctx context.Context, stream pb.PostOrders_PingPongStream) error {
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
