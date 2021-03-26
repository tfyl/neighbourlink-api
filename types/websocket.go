package types

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type WebSocketMapStruct struct {
	m map[int]*SyncCon  // [accountID] to SyncCon
	s sync.RWMutex
}

type SyncCon struct {
	*websocket.Conn
	sync.Mutex
}

func (webmap *WebSocketMapStruct) AddConn (accountID int,websocket *websocket.Conn){
	webmap.s.Lock()
	webmap.m[accountID] = &SyncCon{websocket,sync.Mutex{}}
	webmap.s.Unlock()
}

func (webmap *WebSocketMapStruct) DeleteConn (accountID int){
	webmap.s.Lock()
	delete(webmap.m,accountID)
	webmap.s.Unlock()
}

func (webmap *WebSocketMapStruct) SendAll (post Post) error {
	webmap.s.RLock()

	postByte,err := json.Marshal(post)
	if err != nil{
		return err
	}
	for _,conn := range webmap.m {
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
	WebSocketMap.m = make(map[int]*SyncCon)
	return WebSocketMap
}