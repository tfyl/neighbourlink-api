package types

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type WebSocketMapStruct struct {
	m map[*SyncCon]int  // [*SyncCon] to AccountID
	s sync.RWMutex
}

type SyncCon struct {
	*websocket.Conn
	sync.Mutex
}

func (webmap *WebSocketMapStruct) AddConn (accountID int,websocket *websocket.Conn){
	webmap.s.Lock()
	webmap.m[&SyncCon{websocket,sync.Mutex{}}] = accountID
	webmap.s.Unlock()
}

func (webmap *WebSocketMapStruct) SendAll (post Post) error {
	webmap.s.RLock()

	postByte,err := json.Marshal(post)
	if err != nil{
		return err
	}
	for conn,_ := range webmap.m {
		conn.Lock()
		err := conn.WriteMessage(1,postByte)
		conn.Unlock( )
		if err != nil{
			fmt.Println("Websocket error:",err)
		}
	}

	webmap.s.RUnlock()
	return nil

}

func NewWebSocketMap() (WebSocketMap WebSocketMapStruct) {
	WebSocketMap.m = make(map[*SyncCon]int)
	return WebSocketMap
}