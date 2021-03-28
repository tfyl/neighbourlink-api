package types

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)
// class/struct definition for WebSocketMapStruct
type WebSocketMapStruct struct {
	m map[*SyncCon]int  // [*SyncCon] to AccountID
	s sync.RWMutex
}

type SyncCon struct {
	*websocket.Conn
	sync.Mutex
}
// method of WebSocketMapStruct
// adds connection to map
func (webmap *WebSocketMapStruct) AddConn (accountID int,websocket *websocket.Conn){
	// locks mutex to stop concurrent processing errors on the map
	webmap.s.Lock()
	// add new entry for map (dictionary)
	webmap.m[&SyncCon{websocket,sync.Mutex{}}] = accountID
	// unlocks mutex
	webmap.s.Unlock()
}

func (webmap *WebSocketMapStruct) SendAll (post Post) error {
	// locks mutex to stop concurrent processing errors on the map
	webmap.s.RLock()

	// convert post to byte json data
	postByte,err := json.Marshal(post)
	if err != nil{
		return err
	}
	for conn,_ := range webmap.m {
		conn.Lock()
		// sends byte encoded post via websocket
		err := conn.WriteMessage(1,postByte)
		conn.Unlock( )
		if err != nil{
			fmt.Println("Websocket error:",err)
		}
	}

	// unlocks mutex
	webmap.s.RUnlock()
	return nil

}

// instantiates class/struct
func NewWebSocketMap() (WebSocketMap WebSocketMapStruct) {
	// initialises map for the map attribute
	WebSocketMap.m = make(map[*SyncCon]int)
	return WebSocketMap
}