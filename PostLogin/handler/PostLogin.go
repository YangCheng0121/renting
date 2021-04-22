package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"renting/web/models"
	"renting/web/utils"
	"time"

	pb "renting/PostLogin/proto"
)

type PostLogin struct{}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostLogin) PostLogin(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- 登陆 api/v1.0/sessions ----------------")

	// 返回给前端的map结构体
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 查询数据库
	var user models.User
	o := orm.NewOrm()

	// select * from user
	// 创建查询句柄
	qs := o.QueryTable("user")
	// 查询符合的数据
	err := qs.Filter("mobile", req.Mobile).One(&user)
	if err != nil {
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)

		return nil
	}

	// 判断密码是否正确
	if req.Password != user.Password_hash {
		rsp.Errno = utils.RECODE_PWDERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 编写redis缓存数据库信息
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

	// 生成sessionID
	h := GetMd5String(req.Mobile + req.Password)
	rsp.SessionID = h
	logs.Info(h)

	// 拼接key sessionid + name
	_ = bm.Put(context.TODO(), h+"name", string(user.Name), time.Second*3600)
	// 拼接key sessionid + user_id
	_ = bm.Put(context.TODO(), h+"user_id", string(user.Id), time.Second*3600)
	// 拼接key sessionid + mobile
	_ = bm.Put(context.TODO(), h+"mobile", string(user.Mobile), time.Second*3600)

	// 成功返回数据
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostLogin) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.PostLogin_StreamStream) error {
	log.Infof("Received pb.Stream request with count: %d", req.Count)

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
func (e *PostLogin) PingPong(ctx context.Context, stream pb.PostLogin_PingPongStream) error {
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
