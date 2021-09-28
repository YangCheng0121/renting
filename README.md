# 微服务实战Go Micro v3

## 系列文章
[点击跳转](https://zhuanlan.zhihu.com/p/368545133)

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
<img width="300" height="400" src="https://z3.ax1x.com/2021/06/08/2rVWX8.png"/>
<img width="300" height="400" src="https://z3.ax1x.com/2021/06/08/2rZFc6.png"/>
<img width="300" height="400" src="https://z3.ax1x.com/2021/06/08/2rZVBD.png"/>
