package subscription

import (
	"iQuest/app/graphql/model"
	"sync"
)

type Service struct {
	Rooms map[string]*ChatRoom
	MU    sync.Mutex
}

type ChatRoom struct {
	Name    string
	Message chan *model.Message
}

var Server = Service{Rooms: map[string]*ChatRoom{}}
