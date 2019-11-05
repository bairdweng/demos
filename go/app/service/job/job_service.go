package job

import (
	"iQuest/app/constant"
	"iQuest/app/graphql/prisma"
	"iQuest/db"

	"github.com/adjust/rmq"
)

var QueueConn = rmq.OpenConnectionWithRedisClient(constant.QueueConnectionName, db.Redis())

type Service struct {
	Prisma *prisma.Client
}

func CloseQueueConn() bool {
	return QueueConn.Close()
}

func CloseAllQueueConn() error {
	return QueueConn.CloseAllQueuesInConnection()
}
