package handler

import (
	"context"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"io/ioutil"
	"os"
	"path"
	"renting/web/models"
	"renting/web/utils"
	"strconv"

	pb "renting/PostHousesImage/proto"
)

type PostHousesImage struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostHousesImage) PostHousesImage(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	// 打印被调用的函数
	logs.Info("---------------- 发送房屋图片PostHousesImage  /api/v1.0/houses/:id/images ----------------")

	// 初始化返回正确的返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/* 获取文件的后缀名 */ // dsnlkjfajadskfksda.sadsdasd.sdasd.jpg
	logs.Info("后缀名", path.Ext(req.Filename))
	fileType := "other"
	/* 存储文件到 fastdfs 当中并且获取 url */
	//.jpg
	fileExt := path.Ext(req.Filename)
	logs.Info(fileExt)

	if fileExt == ".jpg" || fileExt == ".png" || fileExt == ".gif" || fileExt == ".jpeg" {
		fileType = "img"
	}

	if fileType != "img" {
		logs.Info("不支持该文件上传")
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	err := ioutil.WriteFile(req.Filename, req.Image, 0666)
	if err != nil {
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// group1 group1/M00/00/00/wKgLg1t08pmANXH1AAaInSze-cQ589.jpg
	// 上传数据
	houseImageFile, err := models.UploadByFileName(req.Filename)
	if err != nil {
		logs.Info("文件读写错误", err)
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("houseImageFile:", houseImageFile)
	defer os.Remove(req.Filename)

	/* 从请求url中得到我们的house_id */
	houseId, _ := strconv.Atoi(req.Id)

	// 创建house对象
	house := models.House{Id: houseId}
	// 创建数据库句柄
	o := orm.NewOrm()
	err = o.Read(&house)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/* 判断index_image_url 是否为空 */
	if house.IndexImageUrl == "" {
		/* 空就把这张图片设置为主图片 */
		house.IndexImageUrl = houseImageFile.Path
	}

	/* 将该图片添加到 house 的全部图片当中 */
	houseImage := models.HouseImage{House: &house, Url: houseImageFile.Path}

	house.Images = append(house.Images, &houseImage)
	// 将图片对象插入表单之中
	_, err = o.Insert(&houseImage)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 对house表进行更新
	_, err = o.Update(&house)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/* 返回正确的数据回显给前端 */
	rsp.Url = houseImageFile.Path
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostHousesImage) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.PostHousesImage_StreamStream) error {
	log.Infof("Received PostHousesImage.Stream request with count: %d", req.Count)

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
func (e *PostHousesImage) PingPong(ctx context.Context, stream pb.PostHousesImage_PingPongStream) error {
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
