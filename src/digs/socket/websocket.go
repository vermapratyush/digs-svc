package socket

import (
	"github.com/gorilla/websocket"
	"sync"
	"digs/logger"
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
	defer lookUpLock.Unlock()

	delete(lookUp, uid)
}

func GetLookUp(uid string) (Peer, bool) {
	lookUpLock.RLock()
	defer lookUpLock.RUnlock()

	peer, present := lookUp[uid]
	return peer, present
}

func SetLookUp(uid string, peer Peer) {
	lookUpLock.Lock()
	defer lookUpLock.Unlock()
	lookUp[uid] = peer
}

func SendData(uid string, data []byte) error {
	lookUpLock.RLock()
	defer lookUpLock.RUnlock()

	peer := lookUp[uid];

	peer.wsLock.Lock()
	defer peer.wsLock.Unlock()

	err := peer.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logger.Error("SOCKET|UnableToWriteToSocket|UID=", uid, "|Error=%v", err)
	}
	return err
}

