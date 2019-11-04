package main

import (
	"iQuest/db"
	"iQuest/router"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	defer db.Close()
	g := gin.Default()
	router.Load(g)

	g.Run()
	// 貌似起到监听的作用
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
