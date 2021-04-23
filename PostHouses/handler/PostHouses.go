package handler

import (
	"context"
	"encoding/json"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"reflect"
	"renting/web/models"
	"renting/web/utils"
	"strconv"

	log "github.com/micro/micro/v3/service/logger"

	pb "renting/PostHouses/proto"
)

type PostHouses struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostHouses) PostHouses(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	// 打印被调用的函数
	logs.Info("---------------- PostHouses 发布房源信息 /api/v1.0/houses ----------------")

	// 创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	var requestMap = make(map[string]interface{})
	_ = json.Unmarshal(req.Max, &requestMap)
	for key, value := range requestMap {
		logs.Info(key, value)
	}

	house := models.House{}

	/* 插入房源信息 */
	// "title":"上奥世纪中心"
	house.Title = requestMap["title"].(string)
	// "price":"666",
	price, _ := strconv.Atoi(requestMap["price"].(string))
	house.Price = price * 100
	//	"address":"西三旗桥东建材城1号",
	house.Address = requestMap["address"].(string)
	//	"room_count":"2",
	house.Room_count, _ = strconv.Atoi(requestMap["room_count"].(string))
	//	"acreage":"60",
	house.Acreage, _ = strconv.Atoi(requestMap["acreage"].(string))
	//	"unit":"2室1厅",
	house.Unit = requestMap["unit"].(string)
	//	"capacity":"3",
	house.Capacity, _ = strconv.Atoi(requestMap["capacity"].(string))
	//	"beds":"双人床2张",
	house.Beds = requestMap["beds"].(string)
	//	"deposit":"200",
	deposit, _ := strconv.Atoi(requestMap["deposit"].(string))
	house.Deposit = deposit * 100
	//	"min_days":"3",
	house.Min_days, _ = strconv.Atoi(requestMap["min_days"].(string))
	//	"max_days":"0",
	house.Max_days, _ = strconv.Atoi(requestMap["max_days"].(string))

	// 设施
	// "facility":["1","2","3","7","12","14","16","17","18","21","22"]
	var facility []*models.Facility
	for _, fId := range requestMap["facility"].([]interface{}) {
		fid, _ := strconv.Atoi(fId.(string))
		fac := &models.Facility{Id: fid}
		facility = append(facility, fac)
	}

	//	"area_id":"5"，地区
	areaId, _ := strconv.Atoi(requestMap["area_id"].(string))
	area := models.Area{Id: areaId}
	house.Area = &area

	// 获的 userId
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
	sessionIdUserId := req.Sessionid + "user_id"

	valueId, _ := bm.Get(context.Background(), sessionIdUserId)
	if valueId == nil {
		logs.Info("获取缓存失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	userId := int(valueId.([]uint8)[0])
	logs.Info(userId, reflect.TypeOf(userId))

	// 添加user信息
	user := models.User{Id: userId}
	house.User = &user

	// 创建数据库句柄
	o := orm.NewOrm()
	houseId, err := o.Insert(&house)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info(houseId, reflect.TypeOf(houseId), house.Id)

	/* 插入房源与设施信息的多对多表中 */
	m2m := o.QueryM2M(&house, "Facilities")
	num, err := m2m.Add(facility)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	if num == 0 {
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	rsp.House_Id = int64(house.Id)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostHouses) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.PostHouses_StreamStream) error {
	log.Infof("Received PostHouses.Stream request with count: %d", req.Count)

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
func (e *PostHouses) PingPong(ctx context.Context, stream pb.PostHouses_PingPongStream) error {
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
