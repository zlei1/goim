package comet

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/zlei1/goim/proto"
	"net/http"
	"strconv"
	"time"
)

func (cmt *Comet) InitWebsocket() {
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		serveWs(cmt, writer, request)
	})
	log.Infof("comet websocket bind %s", cmt.C.Websocket.Bind)
	err := http.ListenAndServe(cmt.C.Websocket.Bind, nil)
	if err != nil {
		panic(err)
	}
}

func serveWs(cmt *Comet, writer http.ResponseWriter, request *http.Request) {
	var upGrader = websocket.Upgrader{
		ReadBufferSize: cmt.C.Base.ReadBufferSize,
		WriteBufferSize: cmt.C.Base.WriteBufferSize,
	}
	upGrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upGrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Errorf("comet websocket upgrader err: %s", err.Error())
		return
	}
	var ch *Channel
	ch = NewChannel(cmt.C.Base.BroadcastSize)
	ch.conn = conn

	go writePump(cmt, ch)
	go readPump(cmt, ch)
}

func readPump(cmt *Comet, ch *Channel) {
	defer func() {
		if ch.Room == nil || ch.UserId == 0 {
			ch.conn.Close()
			return
		}
		disconnectReq := new(proto.DisconnectReq)
		disconnectReq.RoomId = ch.Room.Id
		disconnectReq.UserId = ch.UserId
		cmt.assignBucket(strconv.FormatInt(ch.UserId, 10)).DeleteChannel(ch)
		if err := cmt.operator.Disconnect(disconnectReq); err != nil {
			log.Warnf("comet disconnect err: %s", err.Error())
		}
		ch.conn.Close()
	}()

	ch.conn.SetReadLimit(cmt.C.Base.MaxMessageSize)
	ch.conn.SetReadDeadline(time.Now().Add(cmt.C.Base.PongWait))
	ch.conn.SetPongHandler(func(string) error {
		ch.conn.SetReadDeadline(time.Now().Add(cmt.C.Base.PongWait))
		return nil
	})

	for {
		_, message, err := ch.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("readPump ReadMessage err: %s", err.Error())
				return
			}
		}
		if message == nil {
			return
		}
		var connectReq *proto.ConnectReq
		if err := json.Unmarshal(message, &connectReq); err != nil {
			log.Errorf("comet readPump message struct %+v", connectReq)
		}
		if connectReq.UserId == 0 {
			log.Warnf("userId empty")
			return
		}
		connectReq.ServerId = cmt.C.Base.ServerId
		err = cmt.operator.Connect(connectReq)
		if err != nil {
			return
		}
		bucket := cmt.assignBucket(strconv.FormatInt(connectReq.UserId, 10))
		err = bucket.PutChannel(connectReq.UserId, connectReq.RoomId, ch)
		if err != nil {
			log.Errorf("bucket putchannel err: %s", err.Error())
			ch.conn.Close()
		}
		log.Infof("message: %s", message)
	}
}

func writePump(cmt *Comet, ch *Channel) {
	ticker := time.NewTicker(cmt.C.Base.PingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-ch.broadcast:
			// write data dead time , like http timeout , default 10s
			ch.conn.SetWriteDeadline(time.Now().Add(cmt.C.Base.WriteWait))
			if !ok {
				log.Warn("writePump SetWriteDeadline not ok")
				ch.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ch.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Warn("WritePump ch.conn.NextWriter err: %s", err.Error())
				return
			}
			w.Write(message.Body)
			log.Infof("message write body: %s", message.Body)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// heartbeat，if ping error will exit and close current websocket conn
			ch.conn.SetWriteDeadline(time.Now().Add(cmt.C.Base.WriteWait))
			if err := ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}