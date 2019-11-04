package router

import (
	"iQuest/app/graphql"

	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine) *gin.Engine {

	g.GET("/graphql", graphql.Handler())
	g.POST("/graphql", graphql.Handler())
	g.GET("/playground", graphql.Playground())

	return g
}
