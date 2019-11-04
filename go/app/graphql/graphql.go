package graphql

import (
	"context"

	"iQuest/app/graphql/directive"

	// "iQuest/app/graphql/resolver"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"

	"iQuest/app/graphql/prisma"
	"iQuest/app/graphql/schema"
	"iQuest/config"
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
		// Resolvers: &resolver.Resolver{
		// 	Prisma: client,
		// },
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
			// debugUser(c)
		}

		ginCtx := context.WithValue(c.Request.Context(), "ginContext", c)
		c.Request = c.Request.WithContext(ginCtx)

		print("123123")
		h.ServeHTTP(c.Writer, c.Request)
	}
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
