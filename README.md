# ginuser
使用gin框架实现用户注册登录功能



## 包含的技术栈
- golang
- gin
- gorm
- mysql
- redis缓存
- jwt认证
- jaeger链路追踪
- zap日志







## 目录结构
```aiignore
.
├── README.md
├── config
│   ├── config.go
│   ├── db.go
│   ├── logger.go
│   ├── redis.go
│   └── tracer.go
├── consts
│   └── track.go
├── controller
│   └── user.go
├── etc
│   ├── dev.yaml
│   └── prod.yaml
├── go.mod
├── go.sum
├── main.go
├── middleware
│   ├── auth.go
│   ├── cors.go
│   ├── log.go
│   ├── promMiddleware.go
│   └── request_trace_log.go
├── pkg
│   ├── e
│   │   ├── code.go
│   │   └── msg.go
│   ├── jwt
│   │   └── jwt.go
│   └── utils
│       ├── ctl
│       │   ├── ctl.go
│       │   └── user_info.go
│       └── track
│           └── track.go
├── repository
│   └── db
│       ├── dao
│       │   └── user.go
│       └── model
│           └── user.go
├── router
│   └── router.go
├── service
│   └── user.go
├── test
│   ├── etc
│   │   ├── dev.yaml
│   ├── ginuser_test.go
│   └── init.go
└── types
    └── user.go

```