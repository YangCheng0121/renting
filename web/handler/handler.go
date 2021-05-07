package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/asim/go-micro/v3"
	"github.com/beego/beego/v2/core/logs"
	"github.com/julienschmidt/httprouter"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	DELETESESSION "renting/DeleteSession/proto"
	GETAREA "renting/GetArea/proto"
	GETHOUSEINFO "renting/GetHouseInfo/proto"
	GETHOUSES "renting/GetHouses/proto"
	GETIMAGECD "renting/GetImageCd/proto"
	GETINDEX "renting/GetIndex/proto"
	GETSESSION "renting/GetSession/proto"
	GETSMSCD "renting/GetSmsCd/proto"
	GETUSERHOUSES "renting/GetUserHouses/proto"
	GETUSERINFO "renting/GetUserInfo/proto"
	GETUSERORDER "renting/GetUserOrder/proto"
	POSTAVATAR "renting/PostAvatar/proto"
	POSTHOUSES "renting/PostHouses/proto"
	POSTHOUSESIMAGE "renting/PostHousesImage/proto"
	POSTLOGIN "renting/PostLogin/proto"
	POSTORDERS "renting/PostOrders/proto"
	POSTRET "renting/PostRet/proto"
	POSTUSERAUTH "renting/PostUserAuth/proto"
	PUTCOMMENT "renting/PutComment/proto"
	PUTORDERS "renting/PutOrders/proto"
	PUTUSERINFO "renting/PutUserInfo/proto"
	"renting/web/models"
	"renting/web/utils"
)

// 获取地区
func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- 获取地区请求客户端 url : api/v1.0/areas ----------------")

	// 创建新的句柄
	service := micro.NewService()
	// 服务初始化
	service.Init()

	// 创建获取地区的服务并且返回句柄
	client := GETAREA.NewGetAreaService("go.micro.srv.GetArea", service.Client())

	// 调用函数并且获得返回数据
	rsp, err := client.GetArea(context.Background(), &GETAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	logs.Info("rsp:", rsp)

	// 创建返回类型的切片
	var areaList []models.Area

	// 循环读取服务返回的数据
	for _, value := range rsp.Data {
		tmp := models.Area{Id: int(value.Aid), Name: value.Aname, Houses: nil}
		areaList = append(areaList, tmp)
	}

	// 创建返回数据map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   areaList,
	}
	w.Header().Set("Content-Type", "application/json")

	// 将返回数据map发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
}

// 获取验证码图片
func GetImageCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Info("---------------- 获取图片验证码 url：api/v1.0/imagecode/:uuid ----------------")

	// 创建服务
	service := micro.NewService()
	// 服务初始化
	service.Init()

	// 连接服务
	client := GETIMAGECD.NewGetImageCdService("go.micro.srv.GetImageCd", service.Client())

	// 获取前端发送过来的唯一uuid
	logs.Info(ps.ByName("uuid"))
	// 通过句柄调用我们proto协议中准备好的函数
	// 第一个参数为默认,第二个参数 proto协议中准备好的请求包
	rsp, err := client.GetImageCd(context.Background(), &GETIMAGECD.Request{
		Uuid: ps.ByName("uuid"),
	})
	//判断函数调用是否成功
	if err != nil {
		logs.Info(err)
		http.Error(w, err.Error(), 502)
		return
	}

	// 处理前端发送过来的图片信息
	var img image.RGBA

	img.Stride = int(rsp.Stride)

	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	img.Rect.Max.X = int(rsp.Max.X)
	img.Rect.Max.Y = int(rsp.Max.Y)

	img.Pix = []uint8(rsp.Pix)

	var captchaImage captcha.Image
	captchaImage.RGBA = &img

	// 将图片发送给前端
	_ = png.Encode(w, captchaImage)
}

// 获取短信验证码
func GetSmsCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Info("---------------- 获取短信验证   api/v1.0/smscode/:id ----------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	// 获取前端发送过来的手机号码
	mobile := ps.ByName("mobile")
	logs.Info(mobile)

	//后端进行正则匹配
	//创建正则句柄
	myreg := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	//进行正则匹配
	bo := myreg.MatchString(mobile)

	//如果手机号错误则
	if bo == false {
		// we want to augment the response
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "手机号错误",
		}
		//设置返回数据格式
		w.Header().Set("Content-Type", "application/json")

		//将错误发送给前端
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		logs.Info("手机号错误返回")
		return
	}

	// 获取url携带的验证码 和key（uuid）
	logs.Info(r.URL.Query())
	// 获取url携带的参数
	text := r.URL.Query()["text"][0] // text=248484

	id := r.URL.Query()["id"][0] // id=9cd8faa9-5653-4f7c-b653-0a58a8a98c81

	// 调用服务
	client := GETSMSCD.NewGetSmsCdService("go.micro.srv.GetSmsCd", service.Client())
	rsp, err := client.GetSmsCd(context.Background(), &GETSMSCD.Request{
		Mobile: mobile,
		Id:     id,
		Text:   text,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)
		logs.Info(err)
		return
	}
	// 创建返回map
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	// 设置返回格式
	w.Header().Set("Content-Type", "application/json")

	// 将数据返回给前端
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logs.Info(err)
		return
	}
}

// 获取session
func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- 获取Session url：api/v1.0/session ----------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	// 创建句柄
	client := GETSESSION.NewGetSessionService("go.micro.srv.GetSession", service.Client())

	// 获取cookie
	userLogin, err := r.Cookie("userLogin")

	// 如果不存在就返回
	if err != nil {
		// 创建返回数据map
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	// 存在就发送数据给服务
	rsp, err := client.GetSession(context.Background(), &GETSESSION.Request{
		Sessionid: userLogin.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)
		logs.Info(err)
		return
	}

	// we want to augment the response
	//将获取到的用户名返回给前端
	data := make(map[string]string)
	data["name"] = rsp.Data
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	w.Header().Set("Content-Type", "application/json")

	// 将返回数据map发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
}

//注册请求
func PostRet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Info("---------------- 注册请求   /api/v1.0/users ----------------")
	/* 获取前端发送过来的json数据 */
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for key, value := range request {
		logs.Info(key, value, reflect.TypeOf(value))
	}

	// 由于前端每作所以后端进行下操作
	if request["mobile"] == "" || request["password"] == "" || request["sms_code"] == "" {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "信息有误请重新输入",
		}

		// 如果不存在直接给前端返回
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		logs.Info("有数据为空")
		return
	}

	// 创建服务
	service := micro.NewService()
	service.Init()

	// 连接服务将数据发送给注册服务进行注册
	client := POSTRET.NewPostRetService("go.micro.srv.PostRet", service.Client())

	rsp, err := client.PostRet(context.Background(), &POSTRET.Request{
		Mobile:   request["mobile"].(string),
		Password: request["password"].(string),
		SmsCode:  request["sms_code"].(string),
	})

	if err != nil {
		http.Error(w, err.Error(), 502)
		logs.Info(err)
		return
	}

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	// 读取cookie
	cookie, err := r.Cookie("userLogin")

	// 如果读取失败或者cookie的value中不存在则创建cookie
	if err != nil || "" == cookie.Value {
		cookie := http.Cookie{Name: "userLogin", Value: rsp.SessionID, Path: "/", MaxAge: 600}
		http.SetCookie(w, &cookie)
	}

	// 设置返回模式
	w.Header().Set("Content-Type", "application/json")

	// 将数据回发给前端
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logs.Info(err)
		return
	}

	return
}

// 登录
func PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- 登陆 api/v1.0/sessions ----------------")
	// 获取前端 post 请求发送的内容
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for key, value := range request {
		logs.Info(key, value, reflect.TypeOf(value))
	}

	// 判断账号密码是否为空
	if request["mobile"] == "" || request["password"] == "" {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "信息有误请重新输入",
		}
		w.Header().Set("Content-Type", "application/json")

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		logs.Info("有数据为空")
		return
	}

	// 创建连接
	service := micro.NewService()
	service.Init()

	client := POSTLOGIN.NewPostLoginService("go.micro.srv.PostLogin", service.Client())

	rsp, err := client.PostLogin(context.Background(), &POSTLOGIN.Request{
		Password: request["password"].(string),
		Mobile:   request["mobile"].(string),
	})

	if err != nil {
		http.Error(w, err.Error(), 502)
		logs.Info(err)
		return
	}

	cookie, err := r.Cookie("userLogin")
	if err != nil || "" == cookie.Value {
		cookie := http.Cookie{Name: "userLogin", Value: rsp.SessionID, Path: "/", MaxAge: 600}
		http.SetCookie(w, &cookie)
	}
	logs.Info(rsp.SessionID)
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logs.Info(err)
		return
	}
}

// 退出
func DeleteSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- DELETE  /api/v1.0/session DeleteSession() ------------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	client := DELETESESSION.NewDeleteSessionService("go.micro.srv.DeleteSession", service.Client())

	// 获取session
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	rsp, err := client.DeleteSession(context.Background(), &DELETESESSION.Request{
		Sessionid: userLogin.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)
		logs.Info(err)
		return
	}
	// 再次读取数据
	cookie, err := r.Cookie("userLogin")

	// 数据不为空则数据设置副的
	if err != nil || "" == cookie.Value {
		return
	} else {
		cookie := http.Cookie{Name: "userLogin", Path: "/", MaxAge: -1}
		http.SetCookie(w, &cookie)
	}

	// 返回数据
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logs.Info(err)
		return
	}

	return
}

// 获取用户信息
func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- GetUserInfo  获取用户信息   /api/v1.0/user ------------------")
	// 创建服务
	service := micro.NewService()
	service.Init()

	// 创建句柄
	client := GETUSERINFO.NewGetUserInfoService("go.micro.srv.GetUserInfo", service.Client())

	// 获取用户的登录信息
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	// 成功将信息发送给前端
	rsp, err := client.GetUserInfo(context.Background(), &GETUSERINFO.Request{
		Sessionid: userLogin.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)
		logs.Info(err)
		return
	}

	// 准备1个数据的map
	data := make(map[string]interface{})
	// 将信息发送给前端
	data["user_id"] = int(rsp.UserId)
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	// 设置格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logs.Info(err)
		return
	}
	return
}

// 上传用户头像 PostAvatar
func PostAvatar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- 上传用户头像 PostAvatar /api/v1.0/user/avatar ------------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	// 创建句柄
	client := POSTAVATAR.NewPostAvatarService("go.micro.srv.PostAvatar", service.Client())

	// 查看登录信息
	userLogin, err := r.Cookie("userLogin")

	// 如果没有登录就返回错误
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	// 接受前端发送过来的文集
	file, handler, err := r.FormFile("avatar")

	// 判断是否接受成功
	if err != nil {
		logs.Info("PostPpAvatar   c.GetFile(avatar) err", err)

		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}
	// 打印基本信息
	logs.Info(file, handler)
	logs.Info("文件大小", handler.Size)
	logs.Info("文件名", handler.Filename)

	// 二进制的空间用来存储文件
	fileBuffer := make([]byte, handler.Size)

	// 将文件读取到 fileBuffer 里
	_, err = file.Read(fileBuffer)
	if err != nil {
		logs.Info("PostUpAvatar   file.Read(fileBuffer) err", err)
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	// 调用函数传入数据
	rsp, err := client.PostAvatar(context.Background(), &POSTAVATAR.Request{
		Sessionid: userLogin.Value,
		Filename:  handler.Filename,
		Filesize:  handler.Size,
		Avatar:    fileBuffer,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		logs.Info(err)
		return
	}

	// 准备回传数据空间
	data := make(map[string]interface{})
	// url拼接回传数据
	data["avatar_url"] = rsp.AvatarUrl

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logs.Info(err)
		return
	}

	return
}

// 更新用户名 PutUserInfo
func PutUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- 更新用户名 PutUserInfo /api/v1.0/user/name ------------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	// 接受前端发送内容
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 调用服务
	client := PUTUSERINFO.NewPutUserInfoService("go.micro.srv.PutUserInfo", service.Client())

	// 获取用户登录信息
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	rsp, err := client.PutUserInfo(context.Background(), &PUTUSERINFO.Request{
		Sessionid: userLogin.Value,
		Username:  request["name"].(string),
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 接受回发数据
	data := make(map[string]interface{})
	data["name"] = rsp.Username

	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	// 返回前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

// 检查实名认证 GetUserAuth
func GetUserAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("----------------GetUserInfo  获取用户信息   /api/v1.0/user ------------------")

	// 初始化服务
	service := micro.NewService()
	service.Init()

	// 创建句柄
	client := GETUSERINFO.NewGetUserInfoService("go.micro.srv.GetUserInfo", service.Client())

	// 获取用户的登录信息
	userLogin, err := r.Cookie("userLogin")

	// 判断是否成功不成功就直接返回
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	// 成功就将信息发送给前端
	rsp, err := client.GetUserInfo(context.Background(), &GETUSERINFO.Request{
		Sessionid: userLogin.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		logs.Info(err)
		return
	}

	// 准备1个数据的map
	data := make(map[string]interface{})
	// 将信息发送给前端
	data["user_id"] = int(rsp.UserId)
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = rsp.AvatarUrl

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	// 设置格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logs.Info(err)
		return
	}
	return
}

// 实名认证 PostUserAuth
func PostUserAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- 实名认证 PostUserAuth  api/v1.0/user/auth ------------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	// 获取前端发送的数据
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	client := POSTUSERAUTH.NewPostUserAuthService("go.micro.srv.PostUserAuth", service.Client())

	// 获取cookie
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	rsp, err := client.PostUserAuth(context.Background(), &POSTUSERAUTH.Request{
		Sessionid: userLogin.Value,
		RealName:  request["real_name"].(string),
		IdCard:    request["id_card"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

// 获取当前用户所发布的资源 GetUserHouses
func GetUserHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- 获取当前用户所发布的房源 GetUserHouses /api/v1.0/user/houses ------------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	// call the backend service
	client := GETUSERHOUSES.NewGetUserHousesService("go.micro.srv.GetUserHouses", service.Client())

	// 获取cookie
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	rsp, err := client.GetUserHouses(context.TODO(), &GETUSERHOUSES.Request{
		Sessionid: userLogin.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var houseList []models.House
	_ = json.Unmarshal(rsp.Mix, &houseList)

	var houses []interface{}
	for _, houseInfo := range houseList {
		fmt.Printf("house.user = %+v\n", houseInfo.Id)
		fmt.Printf("house.area = %+v\n", houseInfo.Area)
		houses = append(houses, houseInfo.ToHouseInfo())
	}

	dataMap := make(map[string]interface{})
	dataMap["houses"] = houses

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   dataMap,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

// 发布房源信息
func PostHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- PostHouses 发布房源信息 /api/v1.0/houses ------------------")

	// 获取前端post请求发送的内容
	body, _ := ioutil.ReadAll(r.Body)

	// 获取cookie
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置回传格式
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	// 创建连接
	service := micro.NewService()
	service.Init()

	client := POSTHOUSES.NewPostHousesService("go.micro.srv.PostHouses", service.Client())

	rsp, err := client.PostHouses(context.Background(), &POSTHOUSES.Request{
		Sessionid: userLogin.Value,
		Max:       body,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		logs.Info(err)
		return
	}
	/* 得到插入房源信息表的 id */
	houseIdMap := make(map[string]interface{})
	houseIdMap["house_id"] = int(rsp.House_Id)

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   houseIdMap,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logs.Info(err)
		return
	}
}

// 发送房屋图片
func PostHousesImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Info("---------------- 发送房屋图片PostHousesImage  /api/v1.0/houses/:id/images ------------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	// call the backend service
	client := POSTHOUSESIMAGE.NewPostHousesImageService("go.micro.srv.PostHousesImage", service.Client())
	// 获取houseId
	houseId := ps.ByName("id")
	// 获取sessionId
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	file, handler, err := r.FormFile("house_image")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	logs.Info(file, handler)
	logs.Info("文件大小", handler.Size)
	logs.Info("文件名", handler.Filename)

	// 二进制的空间用来存储文件
	fileBuffer := make([]byte, handler.Size)
	// 将文件读取到 fileBuffer 里
	_, err = file.Read(fileBuffer)
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	rsp, err := client.PostHousesImage(context.Background(), &POSTHOUSESIMAGE.Request{
		Sessionid: userLogin.Value,
		Id:        houseId,
		Image:     fileBuffer,
		Filesize:  handler.Size,
		Filename:  handler.Filename,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 准备返回值
	data := make(map[string]interface{})
	data["url"] = utils.AddDomain2Url(rsp.Url)
	// 返回数据map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	// 回发数据
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

// 获取房源详细信息
func GetHouseInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Info("---------------- 获取房源详细信息 GetHouseInfo  api/v1.0/houses/:id ------------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	// call the backend service
	client := GETHOUSEINFO.NewGetHouseInfoService("go.micro.srv.GetHouseInfo", service.Client())

	id := ps.ByName("id")

	// 获取sessionId
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	rsp, err := client.GetHouseInfo(context.Background(), &GETHOUSEINFO.Request{
		Sessionid: userLogin.Value,
		Id:        id,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   nil,
	}

	if rsp.Userid > 0 {
		house := models.House{}
		_ = json.Unmarshal(rsp.Housedata, &house)
		dataMap := make(map[string]interface{})
		dataMap["user_id"] = int(rsp.Userid)
		dataMap["house"] = house.ToOneHouseDesc()
		response["data"] = dataMap
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
	return
}

// 获取首页轮播
func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- 获取首页轮播 url：api/v1.0/houses/index ----------------")
	service := micro.NewService()
	service.Init()

	client := GETINDEX.NewGetIndexService("go.micro.srv.GetIndex", service.Client())

	rsp, err := client.GetIndex(context.Background(), &GETINDEX.Request{})
	if err != nil {
		logs.Info(err)
		http.Error(w, err.Error(), 502)
		return
	}

	var data []interface{}
	_ = json.Unmarshal(rsp.Max, &data)

	// 创建反馈数据map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	w.Header().Set("Content-Type", "application/json")
	// 将返回数据map发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
}

// 搜索房屋
func GetHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- 搜索房屋 url：/api/v1.0/houses ----------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	// call the backend service
	client := GETHOUSES.NewGetHousesService("go.micro.srv.GetHouses", service.Client())

	// aid=5&sd=2017-11-12&ed=2017-11-30&sk=new&p=1
	aid := r.URL.Query()["aid"][0] // aid=5   		 地区编号
	sd := r.URL.Query()["sd"][0]   // sd=2017-11-1   开始世界
	ed := r.URL.Query()["ed"][0]   // ed=2017-11-3   结束世界
	sk := r.URL.Query()["sk"][0]   // sk=new    	 第三栏条件
	p := r.URL.Query()["p"][0]     // tp=1   		 页数

	rsp, err := client.GetHouses(context.Background(), &GETHOUSES.Request{
		Aid: aid,
		Sd:  sd,
		Ed:  ed,
		Sk:  sk,
		P:   p,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var housesL []interface{}
	_ = json.Unmarshal(rsp.Houses, &housesL)

	data := map[string]interface{}{}
	data["current_page"] = rsp.CurrentPage
	data["houses"] = housesL
	data["totalPage"] = rsp.TotalPage

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

// 发布订单
func PostOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- PostOrders  发布订单 /api/v1.0/orders ----------------")

	// 将 post 带过来的数据转化一下
	body, _ := ioutil.ReadAll(r.Body)

	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	// 创建服务
	service := micro.NewService()
	service.Init()

	// 调用服务
	client := POSTORDERS.NewPostOrdersService("go.micro.srv.PostOrders", service.Client())

	rsp, err := client.PostOrders(context.Background(), &POSTORDERS.Request{
		Sessionid: userLogin.Value,
		Body:      body,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	/* 得到插入信息房源表的 id */
	houseIdMap := make(map[string]interface{})
	houseIdMap["order_id"] = int(rsp.OrderId)

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   houseIdMap,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

// 获取订单
func GetUserOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- /api/v1.0/user/orders  GetUserOrder 获取订单 ----------------")

	// 创建服务
	service := micro.NewService()
	service.Init()

	// call the backend service
	client := GETUSERORDER.NewGetUserOrderService("go.micro.srv.GetUserOrder", service.Client())

	// 获取cookie
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
		}
		return
	}

	// 获取role
	role := r.URL.Query()["role"][0] // role

	rsp, err := client.GetUserOrder(context.Background(), &GETUSERORDER.Request{
		Sessionid: userLogin.Value,
		Role:      role,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var orderList []interface{}
	_ = json.Unmarshal(rsp.Orders, &orderList)

	data := map[string]interface{}{}
	data["orders"] = orderList

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	// encode the write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

// 房东同意/拒绝订单
func PutOrders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// decode the incoming request as json
	// 接收请求携带的数据
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
	//获取cookie
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 502)
			logs.Info(err)
			return
		}
		return
	}

	// 创建服务
	service := micro.NewService()
	service.Init()

	client := PUTORDERS.NewPutOrdersService("go.micro.srv.PutOrders", service.Client())

	rsp, err := client.PutOrders(context.Background(), &PUTORDERS.Request{
		Sessionid: userLogin.Value,
		Action:    request["action"].(string),
		Orderid:   ps.ByName("id"),
	})
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 504)
		return
	}
}

// 用户评价订单
func PutComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logs.Info("PutComment  用户评价 /api/v1.0/orders/:id/comment")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := micro.NewService()
	service.Init()
	// call the backend service
	client := PUTCOMMENT.NewPutCommentService("go.micro.srv.PutComment", service.Client())

	//获取cookie
	userLogin, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logs.Info(err)
			return
		}
		return
	}

	rsp, err := client.PutComment(context.TODO(), &PUTCOMMENT.Request{

		Sessionid: userLogin.Value,
		Comment:   request["comment"].(string),
		OrderId:   ps.ByName("id"),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}
