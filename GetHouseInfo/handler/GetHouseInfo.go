package handler

import (
	"context"
	"encoding/json"
	"fmt"
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

	pb "renting/GetHouseInfo/proto"
)

type GetHouseInfo struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetHouseInfo) GetHouseInfo(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- 获取房源详细信息 GetHouseInfo  api/v1.0/houses/:id ----------------")

	// 创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/* 从session 中获取我们的user_id的字段 得到当前用户 id */
	/* 通过 session 获取我们当前登陆用户的 user_id */
	// 构建连接缓存的数据
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
	sessionIdUserId := req.Sessionid + "user_id"

	valueId, err := bm.Get(context.TODO(), sessionIdUserId)
	if valueId == nil {
		logs.Info("获取登录缓存失败", err)
		rsp.Errno = utils.RECODE_SESSIONERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	id := int(valueId.([]uint8)[0])
	logs.Info(id, reflect.TypeOf(id))

	/* 从请求中的url获取房源id */
	houseId, _ := strconv.Atoi(req.Id)

	/* 从缓存数据库中获取当前房屋的数据 */
	houseInfoKey := fmt.Sprintf("house_info_%s", houseId)
	houseInfoValue, _ := bm.Get(context.TODO(), houseInfoKey)
	if houseInfoValue != nil {
		rsp.Userid = int64(id)
		rsp.Housedata = houseInfoValue.([]byte)
	}

	/* 查询当前数据库得到当前的house详细信息 */
	// 创建数据对象
	house := models.House{Id: houseId}
	// 创建数据库句柄
	o := orm.NewOrm()
	_ = o.Read(&house)
	/* 关联查询 area user images fac等表 */
	_, _ = o.LoadRelated(&house, "Area")
	_, _ = o.LoadRelated(&house, "User")
	_, _ = o.LoadRelated(&house, "Images")
	_, _ = o.LoadRelated(&house, "Facilities")

	/* 将查询到的结果存储到缓存当中 */
	houseMix, err := json.Marshal(house)
	_ = bm.Put(context.TODO(), houseInfoKey, houseMix, time.Second*3600)

	/* 返回正确数据给前端 */
	rsp.Userid = int64(id)
	rsp.Housedata = houseMix
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetHouseInfo) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetHouseInfo_StreamStream) error {
	log.Infof("Received GetHouseInfo.Stream request with count: %d", req.Count)

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
func (e *GetHouseInfo) PingPong(ctx context.Context, stream pb.GetHouseInfo_PingPongStream) error {
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
