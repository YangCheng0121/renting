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

	pb "renting/PutComment/proto"
)

type PutComment struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PutComment) PutComment(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- api/v1.0/orders  PutComment post success----------------")
	// 创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	// 1得到被评论的order_id
	// 获得用户id
	// 构建连接缓存的数据
	redisConfigMap := map[string]string{
		"key": utils.G_server_name,
		//"conn":"127.0.0.1:6379",
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

	// 得到订单id
	orderId, _ := strconv.Atoi(req.OrderId)
	// 获得参数

	comment := req.Comment
	// 检验评价信息是否合法 确保不为空
	if comment == "" {

		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 2根据order_id找到所关联的房源信息
	// 查询数据库，订单必须存在，订单状态必须为WAIT_COMMENT待评价状态
	order := models.OrderHouse{}
	o := orm.NewOrm()
	if err := o.QueryTable("order_house").Filter("id", orderId).Filter("status", models.ORDER_STATUS_WAIT_COMMENT).One(&order); err != nil {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 关联查询order订单所关联的user信息
	if _, err := o.LoadRelated(&order, "User"); err != nil {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 确保订单所关联的用户和该用户是同一个人
	if userId != order.User.Id {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 关联查询order订单所关联的House信息
	if _, err := o.LoadRelated(&order, "House"); err != nil {

		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)

		return nil
	}
	house := order.House
	// 3将房源信息的评论字段追加评论信息
	// 更新order的status为COMPLETE
	order.Status = models.ORDER_STATUS_COMPLETE
	order.Comment = comment

	// 将房屋订单成交量+1
	house.OrderCount++

	// 将order和house更新数据库
	if _, err := o.Update(&order, "status", "comment"); err != nil {
		logs.Error("update order status, comment error, err = ", err)
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	if _, err := o.Update(house, "order_count"); err != nil {
		logs.Error("update house order_count error, err = ", err)
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 将house_info_[house_id]的缓存key删除 （因为已经修改订单数量）
	houseInfoKey := "house_info_" + strconv.Itoa(house.Id)
	if err := bm.Delete(context.TODO(), houseInfoKey); err != nil {
		logs.Error("delete ", houseInfoKey, "error , err = ", err)
	}
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PutComment) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.PutComment_StreamStream) error {
	log.Infof("Received PutComment.Stream request with count: %d", req.Count)

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
func (e *PutComment) PingPong(ctx context.Context, stream pb.PutComment_PingPongStream) error {
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
