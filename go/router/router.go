package router

import (
	"iQuest/app/graphql"
)

func Load(g *gin.Engine) *gin.Engine {
	g.GET("/graphql", graphql.Handler());
	return g
}
