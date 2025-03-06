package websockethandler

import (
	"context"
	"log"
	"net/http"
	"postgre/database"
	"time"

	"github.com/gorilla/websocket"
)

func NewMyUpgrader() *MyUpgrader {
	return &MyUpgrader{
		&websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		make(chan []byte),
	}
}

type MyUpgrader struct {
	vc                 *websocket.Upgrader
	responsefromperson chan []byte
}

func (upgrader *MyUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return upgrader.vc.Upgrade(w, r, responseHeader)
}
func (upgrader *MyUpgrader) HandleWebSocketOnline(conn *websocket.Conn, userID int, ps *database.PersonRepository, ctx context.Context) {

	conn.SetPingHandler(func(data string) error {
		return conn.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(time.Second))
	})


	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	defer func() {
		conn.Close()
		ps.UpdateOnline(userID, false, ctx)
	}()


	done := make(chan struct{})
	errChan := make(chan error)


	go func() {
		defer close(done)
		for {
			messageType, _, err := conn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}

			switch messageType {
			case websocket.TextMessage:

			case websocket.PingMessage:
				if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case err := <-errChan:
			log.Printf("WebSocket error for user %d: %v", userID, err)
			return
		case <-ctx.Done():

			err := conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("Error sending close message: %v", err)
			}
			return
		case <-ticker.C:
			if err := ps.UpdateLastEntered(userID, time.Now(), "", ctx); err != nil {
				log.Printf("Error updating last entered: %v", err)
			}
		}
	}
}
