package socket

import (
	"github.com/gorilla/websocket"
	"sync"
	"github.com/astaxie/beego"
)

type Peer struct {
	Conn *websocket.Conn
	UID string
	wsLock *sync.Mutex
}


var lookUp = make(map[string]Peer)
var lookUpLock sync.RWMutex

func GetCopy() map[string]Peer {
	return lookUp
}

func RemoveLookUp(uid string) {
	lookUpLock.Lock()
	delete(lookUp, uid)
	lookUpLock.Unlock()
}

func GetLookUp(uid string) (Peer, bool) {
	lookUpLock.RLock()
	peer, present := lookUp[uid]
	lookUpLock.RUnlock()
	return peer, present
}

func SetLookUp(uid string, peer Peer) {
	lookUpLock.Lock()
	lookUp[uid] = peer
	lookUpLock.Unlock()
}

func SendData(uid string, data []byte) error {
	lookUpLock.RLock()
	peer := lookUp[uid];
	lookUpLock.RUnlock()
	beego.Info("peer=", peer)
	peer.wsLock.Lock()
	err := peer.Conn.WriteMessage(websocket.TextMessage, data)
	peer.wsLock.Unlock()
	if err != nil {
		beego.Info("error=", err)
	}
	return err
}

