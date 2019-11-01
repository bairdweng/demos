package main

import (
	"github.com/gin-gonic/gin"
	"iQuest/db"
	"iQuest/router"
)

func main() {
	defer db.Close()
	println("123123k12k31k23")
	g := gin.Default()

	router.Load(g)

}
