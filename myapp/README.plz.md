```bash
.
├── adapter     # 适配器；
│   ├── db
│   │   └── db.go # 数据库配置信息文件，return *sql.DB and error，做数据迁移的目的数据库；
│   └── gorm
│       └── gorm.go # 数据库配置信息文件，return *gorm.DB and error，用于对上面的数据库进行CRUD等操作：
├── app
│   ├── app
│   │   ├── app.go    # app 结构体初始化(new)文件，包含logger日志，gorm.DB数据库操作和validator数据校验；
│   │   ├── bookHandler.go  # 主要文件，handler处理文件，包含List、Create、Update和Delete的先关操作；
│   │   ├── heathHandler.go   # 存活性和就绪性探测
│   │   └── indexHandler.go   # 首页根目录的处理文件
│   ├── handler
│   │   ├── handler.go      # Handler结构体实现了http.handler
│   │   └── logEntry.go     # 重写请求日志，来源：https://github.com/google/go-cloud/blob/master/server/requestlog/requestlog.go
│   └── router
│       ├── middleware      # 中间件
│       │   ├── content_type_json.go    # 中间件，判断content_type是否为json
│       │   └── content_type_json_test.go
│       └── router.go       # 主要文件，路由文件
├── cmd
│   ├── app
│   │   └── main.go         # myApp主入口文件，包含配置文件、数据库、logger的初始化和http的监听；
│   └── migrate
│       └── main.go         # 数据库迁移的主入口文件，使用goose进行数据库迁移。
├── config
│   └── config.go           # 配置文件，使用envdecdoe包从系统环境变量中获取变量
├── docker
│   ├── app
│   │   ├── Dockerfile      # 构建程序以及拷贝编译后的程序到开发环境
│   │   └── bin            
│   │       ├── init.sh     # docker容器主初始化文件，启动迁移，启动app主程序
│   │       └── wait-for-mysql.sh     # 等待mysql运行成功脚本
│   └── mariadb
│       └── Dockerfile      # 数据库
├── docker-compose.yml      # docker compose 文件
├── go.mod
├── go.sum
├── k8s                     # k8s相关文件
│   ├── app-configmap.yaml
│   ├── app-deployment.yaml
│   ├── app-secret.yaml
│   └── app-service.yaml
├── migrations
│   └── 20190805170000_create_books_table.sql     # 数据库迁移文件
├── model
│   └── book.go             # Book模型文件，外部展示的和内部的做处理
├── repository
│   └── book.go             # gorm的CRUD操作文件
└── util
    ├── logger
    │   ├── logger.go       # 设置日志级别以及日志类型
    │   └── logger_test.go  # 日志测试文件
    └── validator
        ├── validator.go    # 数据校验文件
        └── validator_test.go # 数据校验测试文件

22 directories, 28 files

```

