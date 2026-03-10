package websockets

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	ingame     bool
	egress     chan Event
	username   string
	userID     int64
}
type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Manager, username string, userID int64) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		ingame:     false,
		username:   username,
		userID:     userID,
	}
}

/* ------------------------------------- Heart Beat ------------------------------------- */
// HEARTBEAT to keep connection alive a detect if it goes silent
var (
	//how much we wait before the connection drops
	pongWait = 10 * time.Second

	// how often we send pings
	pingInterval = (pongWait * 9) / 10
)

// Resetting timer
func (c *Client) pongHandler(pongMsg string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

/* ------------------------------------- READ ------------------------------------- */
func (c *Client) readMessages() {
	//clean up
	defer func() {
		c.manager.deleteClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		// if closed or unexpectedly closed
		if err != nil {
			//abnormal close
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading message: ", err)
			}
			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Println("error marshalling event:", err)
			break
		}
		log.Println(request)

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("error handeling message", err)
		}

	}
}

/* ------------------------------------- WRITE ------------------------------------- */
// we are reading the egress and sending back what we found on it

func (c *Client) writeMessages() {
	//clean up
	defer func() {
		c.manager.deleteClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				err := c.connection.WriteMessage(websocket.CloseMessage, nil)
				if err != nil {
					log.Println("connection close: ", err)
				}
				return
			}

			data, err := json.Marshal(message)

			if err != nil {
				log.Println(err)
				return
			}

			err = c.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("failed to send message:", err)
			}
			log.Println("message sent", data)

		case <-ticker.C:
			log.Println("ping")
			//send ping to client
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("write message err: ", err)
				return
			}

		}

	}
}

//--------------------------NOT USING IT-----------------------------------
//broadcast for all clients FOR NOW, FOR TEST
// we are getting the message and writing it to every egress
// for wsclients := range c.manager.clients {
// 	wsclients.egress <- payload
// }

// log.Println(messageType)
// log.Println(string(payload))
