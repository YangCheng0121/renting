# 微服务实战Go Micro v3

## 系列文章
* [微服务实战Go Micro v3 系列（一）- 基础篇](https://cleverbamboo.github.io/2021/04/27/GO/微服务实战Go-Micro-v3-系列（一）-基础篇/#more)
* [微服务实战Go Micro v3 系列（二）- HelloWorld](https://cleverbamboo.github.io/2021/04/27/GO/微服务实战Go-Micro-v3-系列（二）-HelloWorld/#more)
* [微服务实战Go-Micro v3 系列（三）- 启动HTTP服务](https://cleverbamboo.github.io/2021/04/28/GO/微服务实战Go-Micro-v3-系列（三）-启动HTTP服务/#more)
* [微服务实战Go Micro v3 系列（四）- 事件驱动(Pub/Sub)](https://cleverbamboo.github.io/2021/05/12/GO/微服务实战Go-Micro-v3-系列（四）-事件驱动-Pub-Sub/#more)
* [微服务实战Go Micro v3 系列（五）- 注册和配置中心](https://cleverbamboo.github.io/2021/06/02/GO/%E5%BE%AE%E6%9C%8D%E5%8A%A1%E5%AE%9E%E6%88%98Go-Micro-v3-%E7%B3%BB%E5%88%97%EF%BC%88%E4%BA%94%EF%BC%89-%E6%B3%A8%E5%86%8C%E5%92%8C%E9%85%8D%E7%BD%AE%E4%B8%AD%E5%BF%83/#more)
* [微服务实战Go Micro v3 系列（六）- 综合篇（爱租房项目）](https://cleverbamboo.github.io/2021/06/08/GO/%E5%BE%AE%E6%9C%8D%E5%8A%A1%E5%AE%9E%E6%88%98Go-Micro-v3-%E7%B3%BB%E5%88%97%EF%BC%88%E5%85%AD%EF%BC%89-%E7%BB%BC%E5%90%88%E7%AF%87%EF%BC%88%E7%88%B1%E7%A7%9F%E6%88%BF%E9%A1%B9%E7%9B%AE%EF%BC%89/#more)

## 技术栈
golang + docker + consul + grpc + protobuf + beego + mysql + redis + fastDFS + nginx

## 目标功能
- [x] 功能模块
    - [x] 用户模块
        - [x] 注册
            - [x] 获取验证码图片服务
            - [x] 获取短信验证码服务
            - [x] 发送注册信息服务
        - [x] 登录
            - [x] 获取session信息服务
            - [x] 获取登录信息服务
        - [x] 退出
        - [x] 个人信息获取
            - [x] 获取用户基本信息服务
            - [x] 更新用户名服务
            - [x] 发送上传用户头像服务
        - [x] 实名认证
            - [x] 获取用户实名信息服务
            - [x] 发送用户实名认证信息服务
    - [x] 房屋模块
        - [x] 首页展示
            - [x] 获取首页轮播图服务
        - [x] 房屋详情
            - [x] 发布房屋详细信息的服务
            - [x] 上传房屋图片的服务
        - [x] 地区列表
        - [x] 房屋搜索
    - [x] 订单模块
        - [x] 订单确认
        - [x] 发布订单
        - [x] 查看订单信息
        - [x] 订单评论

## 项目文档
​	document文件夹下：

1. ​	整体架构图
2. ​	微服务框架图
3. ​	接口文档

## 项目布局
```
├── DeleteSession
│   ├── 退出登录时清除session
├── GetArea
│   ├── 获取地区信息服务
├── GetImageCd
│   ├── 获取验证码图片服务
├── GetSession
│   ├── 获取Session信息服务
├── GetSmscd
│   ├── 获取短信信息服务
├── GetUserHouses
│   ├── 获取用户已发布房屋的服务
├── GetUserInfo
│   ├── 获取用户详细信息的服务
├── IhomeWeb
│   ├── conf 项目配置文件
│   │   ├── app.conf
│   │   ├── data.sql
│   │   └── redis.conf
│   ├── handler
│   │   └── handler.go 配置路由
│   ├── html 项目静态文件
│   ├── main.go 主函数
│   ├── model 数据库模型
│   │   └── models.go
│   ├── plugin.go
│   ├── server.sh
│   └── utils 项目中用到的工具函数
│       ├── config.go
│       ├── error.go
│       └── misc.go
├── PostAvatar
│   ├──	发送（上传）用户头像服务
├── PostHouses
│   ├── 发送（发布）房源信息服务
├── PostHousesImage
│   ├── 发送（上传）房屋图片服务
├── PostLogin
│   ├── 发送登录服务消息
├── PostRet
│   ├── 发现注册信息服务
├── PostUserAuth
│   ├── 发送用户实名认证信息服务
├── PutUserInfo
│   ├── 发送用户信息
├── GetUserAuth
│   ├── 获取（检查）用户实名信息服务
├── PostHousesImage
│   ├── 发送（上传）房屋图片服务
├── GetHouseInfo
│   ├── 获取房屋详细信息服务
├── GetIndex
│   ├── 获取首页轮播图片服务
├── GetHouses
│   ├── 获取（搜索）房源服务
├── PostOrders
│   ├── 发送（发布）订单服务
├── GetUserOrder
│   ├── 获取房东/租户订单信息服务
├── PutOrders
│   ├── 更新房东同意/拒绝订单
├── PutComments
│   ├── 更新用户评价订单信息
└── README.md
```

## 部分效果图
<a style="display:inline-block">
<img width="300" height="400" src="https://z3.ax1x.com/2021/06/08/2rVWX8.png"/>
</a>
<a style="display:inline-block">
<img width="300" height="400" src="https://z3.ax1x.com/2021/06/08/2rZFc6.png"/>
</a>
<a style="display:inline-block">
<img width="300" height="400" src="https://z3.ax1x.com/2021/06/08/2rZVBD.png"/>
</a>