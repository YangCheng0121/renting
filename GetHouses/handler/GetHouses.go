package handler

import (
	"context"
	"encoding/json"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"renting/web/models"
	"renting/web/utils"
	"strconv"

	pb "renting/GetHouses/proto"
)

type GetHouses struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetHouses) GetHouses(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("----------------  api/v1.0/houses  GETHouseData GET success ----------------")

	// 创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 获取url上的参数信息
	// api/v1.0/houses?aid=1&sd=&ed=&sk=price-inc&p=1
	var aid int // 地区
	aid, _ = strconv.Atoi(req.Aid)
	var sd string // 起始时间
	sd = req.Sd
	var ed string // 结束时间
	ed = req.Ed
	var sk string // 第三栏的信息
	sk = req.Sk
	var page int // 页
	page, _ = strconv.Atoi(req.P)
	logs.Info(aid, sd, ed, sk, page)

	/* 返回json */
	var houses []models.House
	// 创建orm句柄
	o := orm.NewOrm()
	// 设置查找的表
	qs := o.QueryTable("house")
	// 根据查询条件来查找内容
	// 查找传入地区的所有房屋
	num, err := qs.Filter("area_id", aid).All(&houses)
	if err != nil {
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 计算一下所有房屋/一页现实的数量
	totalPage := int(num)/models.HouseListPageCapacity + 1
	housePage := 1

	var houseList []interface{}
	for _, house := range houses {
		_, _ = o.LoadRelated(&house, "Area")
		_, _ = o.LoadRelated(&house, "User")
		_, _ = o.LoadRelated(&house, "Images")
		_, _ = o.LoadRelated(&house, "Facilities")
		houseList = append(houseList, house.ToOneHouseDesc())
	}
	rsp.TotalPage = int64(totalPage)
	rsp.CurrentPage = int64(housePage)
	rsp.Houses, _ = json.Marshal(houseList)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetHouses) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetHouses_StreamStream) error {
	log.Infof("Received GetHouses.Stream request with count: %d", req.Count)

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
func (e *GetHouses) PingPong(ctx context.Context, stream pb.GetHouses_PingPongStream) error {
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
