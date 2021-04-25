package handler

import (
	"context"
	"encoding/json"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"renting/web/models"
	"renting/web/utils"
	"time"

	pb "renting/GetIndex/proto"
)

type GetIndex struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetIndex) GetIndex(ctx context.Context, req *pb.Request, rsp *pb.Response) error {

	// 创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	var data []interface{}
	// 1 从缓存服务器中请求 "home_page_data" 字段,如果有值就直接返回
	// 先从缓存中获取房屋数据,将缓存数据返回前端即可
	redisConfigMap := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}

	redisConfig, _ := json.Marshal(redisConfigMap)
	logs.Info(string(redisConfig))

	cacheConn, err := cache.NewCache("redis", string(redisConfig))
	//cache_conn, err := cache.NewCache("redis", `{"key":"ilhome","conn":"127.0.0.1:6379","dbNum":"0"} `)
	if err != nil {
		logs.Debug("connect cache error", err)
	}

	housePageKey := "home_page_data"
	housePageValue, _ := cacheConn.Get(context.TODO(), housePageKey)

	if housePageValue != nil {
		logs.Debug("======= get house page info  from CACHE!!! ========")
		// 直接将二进制发送给客户端
		rsp.Max = housePageValue.([]byte)
	}
	var houses []models.House

	// 2 如果缓存没有，需要从数据库中查到房屋列表
	o := orm.NewOrm()
	if _, err := o.QueryTable("house").Limit(models.HOME_PAGE_MAX_HOUSES).All(&houses); err == nil {
		for _, house := range houses {
			_, _ = o.LoadRelated(&house, "Area")
			_, _ = o.LoadRelated(&house, "User")
			_, _ = o.LoadRelated(&house, "Images")
			_, _ = o.LoadRelated(&house, "Facilities")
			data = append(data, house.ToOneHouseDesc())
		}
	}
	logs.Info(data, houses)

	// 将data存入缓存数据
	housePageValue, _ = json.Marshal(data)
	_ = cacheConn.Put(context.TODO(), housePageKey, housePageValue, 3600*time.Second)

	rsp.Max = housePageValue.([]byte)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetIndex) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetIndex_StreamStream) error {
	log.Infof("Received GetIndex.Stream request with count: %d", req.Count)

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
func (e *GetIndex) PingPong(ctx context.Context, stream pb.GetIndex_PingPongStream) error {
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
