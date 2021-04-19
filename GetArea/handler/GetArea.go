package handler

import (
	"context"
	"encoding/json"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	pb "renting/GetArea/proto"
	"renting/web/models"
	"renting/web/utils"
	"time"
)

type GetArea struct {
}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetArea) GetArea(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- GetArea    api/v1.0/areas ----------------")

	// 初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 连接redis创建句柄
	redisConfigMap := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}

	// 确定连接信息
	logs.Info(redisConfigMap)
	// 将map转化为json
	redisConfig, _ := json.Marshal(redisConfigMap)
	// 连接redis
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logs.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*1获取缓存数据*/
	areasInfoValue, err := bm.Get(context.TODO(), "areas_info")
	if areasInfoValue != nil {
		logs.Info("获取到缓存发送给前端")

		// 用来存放解码的json
		var areaInfo []map[string]interface{}
		// 解码
		err = json.Unmarshal(areasInfoValue.([]byte), &areaInfo)

		// 进行循环赋值
		for _, value := range areaInfo {
			// 创建对于数据类型并进行赋值
			areaAddress := pb.Response_Address{Aid: int32(value["aid"].(float64)), Aname: value["aname"].(string)}
			// 递增到切片
			rsp.Data = append(rsp.Data, &areaAddress)
		}
		return nil
	}
	logs.Info("没有拿到缓存")

	/*2 如果没有缓存我们就从mysql 里进行查询 */

	//orm的操作创建orm句柄
	o := orm.NewOrm()

	// 接受地区信息的切片
	var areas []models.Area
	// 创建查询条件
	qs := o.QueryTable("area")
	// 查询全部地区
	num, err := qs.All(&areas)
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

	logs.Info("写入缓存")

	/*3 获取数据写入缓存 */

	// 将查询到的数据编码成json格式
	areaInfoStr, _ := json.Marshal(areas)

	// Put(key string, val interface{}, timeout time.Duration) error
	// 存入缓存中
	err = bm.Put(context.TODO(), "areas_info", areaInfoStr, time.Second*3600)
	if err != nil {
		logs.Info("数据库中查出数据信息存入缓存中失误", err)
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 返回地区信息
	for _, value := range areas {
		area := pb.Response_Address{Aid: int32(value.Id), Aname: string(value.Name)}
		rsp.Data = append(rsp.Data, &area)
	}
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetArea) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetArea_StreamStream) error {
	log.Logf("Received Example.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&pb.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}
	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetArea) PingPong(ctx context.Context, stream pb.GetArea_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&pb.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
