# ayg-ishouru-app-service
> author：Allen,Sean
---
#### 介绍
应用层服务，逻辑处理层，GraphQL

#### 克隆项目
```bash
git clone git@gitee.com:aiyuangong/ayg-ishouru-app-service.git
```


#### 安装教程
进入根目录
```bash
cp .env.example .env
go mod tidy
go build main.go
```

#### GraphQL教程
* 工具gqlgen https://gqlgen.com/getting-started/
* 配置文件: app/graphql/gqlgen.yml
* schema文件: app/graphql/schema/schema.graphql
* 进去配置文件所在目录执行gqlgen即可生成model和resolver
* 已集成gin，把router/router.go注释的graphql路由打开即可访问
```bash
go get -u github.com/99designs/gqlgen
cd app/graphql
gqlgen
go run github.com/99designs/gqlgen (only window平台)
```

#### 使用说明
> 项目启动前，请先修改env相关配置

运行
```bash
go run main.go
```

编译
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go -o main
```


