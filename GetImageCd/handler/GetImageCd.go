package handler

import (
	"context"
	"encoding/json"
	"github.com/afocus/captcha"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/core/logs"
	"image/color"
	pb "renting/GetImageCd/proto"
	"renting/web/utils"
	"time"
)

type GetImageCd struct{}

func (e *GetImageCd) GetImageCd(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- GET  /api/v1.0/imagecode/:uuid GetImage() ------------------")

	// 创建1个句柄
	cap := captcha.New()
	// 通过句柄调用 字体文件
	if err := cap.SetFont("comic.ttf"); err != nil {
		logs.Info("没有字体文件")
		panic(err.Error())
	}

	// 设置图片的大小
	cap.SetSize(91, 41)
	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	//SetFrontColor(colors ...color.Color)  这两个颜色设置的函数属于不定参函数
	cap.SetFrontColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{R: 255, A: 255}, color.RGBA{B: 255, A: 255}, color.RGBA{G: 153, A: 255})
	// 生成图片 返回图片和 字符串(图片内容的文本形式)
	img, str := cap.Create(4, captcha.NUM)

	logs.Info(str)

	b := *img      //解引用
	c := *(b.RGBA) //解引用
	//成功返回
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//图片信息
	rsp.Pix = []byte(c.Pix)
	rsp.Stride = int64(c.Stride)
	rsp.Max = &pb.Response_Point{X: int64(c.Rect.Max.X), Y: int64(c.Rect.Max.Y)}
	rsp.Min = &pb.Response_Point{X: int64(c.Rect.Min.X), Y: int64(c.Rect.Min.Y)}

	/* 将uuid与 随机数验证码对应的存储在redis缓存中 */
	// 初始化缓存全局变量的对象

	redisConfigMap := map[string]string{
		"key":   "renting",
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}

	logs.Info(redisConfigMap)

	redisConfig, _ := json.Marshal(redisConfigMap)

	// 连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logs.Info("GetImage()   cache.NewCache err ", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	// 验证码进行1个小时缓存
	_ = bm.Put(context.TODO(), req.Uuid, str, 300*time.Second)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetImageCd) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetImageCd_StreamStream) error {
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
func (e *GetImageCd) PingPong(ctx context.Context, stream pb.GetImageCd_PingPongStream) error {
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
