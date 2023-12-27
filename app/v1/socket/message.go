package socket

import (
	"encoding/json"

	"github.com/anyshake/observer/publisher"
	"github.com/gorilla/websocket"
)

func (s *Socket) handleMessage(gp *publisher.Geophone, conn *websocket.Conn) error {
	data, err := json.Marshal(gp)
	if err != nil {
		conn.Close()
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		conn.Close()
		return err
	}

	return nil
}
