package websockets

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"exmaple.com/ulti-restapi/models"
	"exmaple.com/ulti-restapi/utility"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// take an http request and upgrade it to a websocket connection
var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     CheckOrigin,
	}
)

// to stop other websites opening sockets to this backend
// CORS(Cross-Origin Resource Sharing)
var allowedOrigins = map[string]bool{
	"https://localhost:5173": true,
	"https://localhost:8080": true,
	"https://localhost:8443": true,
}

func CheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	return allowedOrigins[origin]
}

/* ------------------------------------- Manager ------------------------------------- */
type Manager struct {
	clients ClientList
	sync.RWMutex

	handlers map[string]EventHandler
	games    GameList
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
		games:    make(GameList),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) WSHandler(context *gin.Context) {
	log.Println("new connection")

	// Get token from query parameter
	token := context.Query("token")
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	// Verify token and get user info
	userID, err := utility.VerifyToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
		return
	}

	conn, err := websocketUpgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m, user.Name, user.ID)

	m.addClient(client)

	// Start client processes
	go client.readMessages()
	go client.writeMessages()

}

/* ------------------------------------- Client Handeling ------------------------------------- */

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

func (m *Manager) deleteClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	_, ok := m.clients[client]
	if ok {
		client.connection.Close()
		delete(m.clients, client)
	}

}

/* ------------------------------------- Game Handeling ------------------------------------- */

func (m *Manager) addGame(game *Game) {
	m.Lock()
	defer m.Unlock()
	m.games[game] = true
}

func (m *Manager) deleteGame(game *Game) {
	m.Lock()
	defer m.Unlock()
	delete(m.games, game)
	log.Printf("Game %d deleted", game.id)
}

/* ------------------------------------- Event Handeling ------------------------------------- */
type EventHandler func(event Event, c *Client) error

func (m *Manager) setupEventHandlers() {
	m.handlers[EventCardClicked] = CardClicked
	m.handlers[EventGameInit] = SendGameInit
	m.handlers[EventGameExit] = GameExit
	m.handlers[EventBidAction] = BidAction
	m.handlers[EventTalonExchange] = TalonExchange
	m.handlers[EventPlayAgain] = PlayAgain
}

// called in readmessage function to recognise the event
func (m *Manager) routeEvent(event Event, c *Client) error {
	handler, ok := m.handlers[event.Type]
	if ok {
		err := handler(event, c)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}
