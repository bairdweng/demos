package graphql

import (
	"context"
	"iQuest/app/graphql/directive"
	"iQuest/app/graphql/prisma"
	"iQuest/app/graphql/resolver"
	"iQuest/app/graphql/schema"
	session "iQuest/app/model/user"
	"iQuest/config"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
)

type Service struct {
	Prisma *prisma.Client
}

var Server Service

// Handler GraphQL 处理
func Handler() gin.HandlerFunc {
	client := prisma.New(&prisma.Options{
		Endpoint: config.Viper.GetString("PRISMA_ENDPOINT"),
		Secret:   config.Viper.GetString("PRISMA_SECRET"),
	})

	Server.Prisma = client
	c := schema.Config{
		Resolvers: &resolver.Resolver{
			Prisma: client,
		},
	}
	c.Directives.IsAuthenticated = directive.IsAuthenticated

	//h := handler.GraphQL(schema.NewExecutableSchema(c), handler.WebsocketUpgrader(websocket.Upgrader{
	//	CheckOrigin: func(r *http.Request) bool {
	//		return true
	//	},
	//}))

	h := handler.GraphQL(schema.NewExecutableSchema(c))

	// 只需要通过Gin简单封装即可
	return func(c *gin.Context) {
		if config.Viper.GetBool("DEBUG") {
			debugUser(c)
		}

		ginCtx := context.WithValue(c.Request.Context(), "ginContext", c)
		c.Request = c.Request.WithContext(ginCtx)

		h.ServeHTTP(c.Writer, c.Request)
	}
}

func debugUser(c *gin.Context) *gin.Context {
	user := session.SessionUser{
		UserID:    44737,
		UserName:  "陈秋会",
		CompanyID: 10001489,
		OpenID:    "openid",
	}
	ctx := context.WithValue(c.Request.Context(), "user", user)
	c.Request = c.Request.WithContext(ctx)
	return c
}

// Playground GraphQL Playground
func Playground() gin.HandlerFunc {
	// 定义playground调用的接口地址对应gin路由
	h := handler.Playground("GraphQL", "/graphql")

	// 只需要通过Gin简单封装即可
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
