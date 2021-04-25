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

	pb "renting/GetUserOrder/proto"
)

type GetUserOrder struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetUserOrder) GetUserOrder(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- GET  /api/v1.0/user GetUserInfo() ------------------")

	// 创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 根据session得到当前用户的user_id
	// 构建连接缓存的数据
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

	// 得到用户角色
	logs.Info(req.Role)

	o := orm.NewOrm()
	var orders []models.OrderHouse
	var orderList []interface{} // 存放预订单切片

	if "landlord" == req.Role {
		// 角色为房东
		// 现在找到自己已经发布了哪些房子
		var landLordHouses []models.House
		_, _ = o.QueryTable("house").Filter("user__id", userId).All(&landLordHouses)
		var housesIds []int
		for _, house := range landLordHouses {
			housesIds = append(housesIds, house.Id)
		}
		// 在从订单中找到房屋id为自己房源的id
		_, _ = o.QueryTable("order_house").Filter("house_id__in", housesIds).OrderBy("ctime").All(&orders)
	} else {
		// 角色为租客
		_, err := o.QueryTable("order_house").Filter("user__id", userId).OrderBy("ctime").All(&orders)
		if err != nil {
			logs.Info(err)
		}
	}
	// 循环将数据放到切片中
	for _, order := range orders {
		_, _ = o.LoadRelated(&order, "User")
		_, _ = o.LoadRelated(&order, "House")
		orderList = append(orderList, order.ToOrderInfo())
	}
	rsp.Orders, _ = json.Marshal(orderList)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetUserOrder) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetUserOrder_StreamStream) error {
	log.Infof("Received GetUserOrder.Stream request with count: %d", req.Count)

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
func (e *GetUserOrder) PingPong(ctx context.Context, stream pb.GetUserOrder_PingPongStream) error {
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
