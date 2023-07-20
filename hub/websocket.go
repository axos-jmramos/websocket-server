package hub

import "github.com/gorilla/websocket"

type Hub struct {
	connections map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[string]*websocket.Conn),
	}
}
func (h *Hub) Add(userID string, conn *websocket.Conn) {
	h.connections[userID] = conn
}

func (h *Hub) NotifyUser(userID string, message interface{}) {
	conn, ok := h.connections[userID]

	//check first if the user id is connected to us
	if ok {
		err := conn.WriteJSON(message)

		//if we had an error sending the message, we delete it in our connections
		if err != nil {
			delete(h.connections, userID)
		}
	}
}
