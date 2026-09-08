package websockets

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
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
	"http://localhost:5173": true,
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

/* ---------  Card Clicked Event  ------------------------------------- */

func CardClicked(event Event, c *Client) error {
	/*
			            type: "card_clicked",
		                id: this.gameid,
		                payload: { cardId },
	*/
	var ev CardClickedEvent
	if err := json.Unmarshal(event.Payload, &ev); err != nil {
		return errors.New("invalid card_clicked payload")
	}

	var (
		nextplayer  int
		roundwinner *Client
		eof         string
		game        *Game
	)
	for g := range c.manager.games {
		if g.id == event.GameId {
			game = g
			nextplayer, roundwinner, eof = game.RoundHandler(c, ev.CardId)
			break
		}
	}

	if game == nil {
		return errors.New("mistake in the game id, no game found")

	}

	// Always send the card to all players first
	if nextplayer != -1 || eof == "end" {
		err := SendCardClickToPlayers(ev.CardId, c, game)
		if err != nil {
			return errors.New("cards cant be sent to other players")
		}
	}

	if eof == "end" {
		err := SendEndResult(game)
		if err != nil {
			return errors.New("couldnt finalize the game ")
		}
		return nil
	}

	if nextplayer != -1 {
		err := SendNextPlayer(nextplayer, game)
		if err != nil {
			return errors.New("next player wasnt sent")
		}
		if roundwinner != nil {
			err := SendEvaluation(roundwinner, game)
			if err != nil {
				return errors.New("something went wrong during round evaluation")
			}
		}
		return nil
	}

	return errors.New("there was a problem during the turn")
}

func SendCardClickToPlayers(cardid int, client *Client, game *Game) error {

	var pindex int
	for i := range game.players {
		if game.players[i] == client {
			pindex = i
			break
		}
	}

	payload, err := json.Marshal(struct {
		CardId      int `json:"cardid"`
		PlayerIndex int `json:"playerindex"`
	}{
		CardId:      cardid,
		PlayerIndex: pindex,
	})

	if err != nil {
		return errors.New("cardclicked send marhals went wrong")
	}

	response := Event{
		Type:    "card_played",
		GameId:  game.id,
		Payload: payload,
	}

	for _, p := range game.players {
		if p == client {
			continue
		}
		p.egress <- response
	}

	return nil
}

func SendNextPlayer(nextp int, game *Game) error {
	log.Printf("Sending next player: %d", nextp)
	tonextplayer := Event{
		Type:   "your_turn",
		GameId: game.id,
	}

	for i, p := range game.players {
		if i == nextp {
			p.egress <- tonextplayer
		}

	}

	return nil
}

func SendEvaluation(roundW *Client, game *Game) error {
	log.Printf("Sending evaluation")

	var (
		points      int
		winnerindex int
	)
	for i, p := range game.players {
		if p == roundW {
			points = game.points[i]
			winnerindex = i
		}
	}

	winnerres, err := json.Marshal(struct {
		Points int `json:"points"`
	}{
		Points: points,
	})

	if err != nil {
		return errors.New("winner response marhals went wrong")
	}
	winner := Event{
		Type:    "you_won",
		GameId:  game.id,
		Payload: winnerres,
	}

	playerres, err := json.Marshal(struct {
		WinnerIndex int `json:"winnerindex"`
	}{
		WinnerIndex: winnerindex,
	})

	if err != nil {
		return errors.New("player response marhals went wrong")
	}
	player := Event{
		Type:    "round_result",
		GameId:  game.id,
		Payload: playerres,
	}

	for _, p := range game.players {
		if p == roundW {
			p.egress <- winner
		}
		p.egress <- player
	}
	return nil
}

func SendEndResult(g *Game) error {
	log.Printf("End of the game")

	bettingpool := g.bettingpool

	// Find declarer's points (wins)
	var declarerPoints int
	for i, p := range g.players {
		if p == g.declarer {
			declarerPoints = g.points[i]
			break
		}
	}

	// Declarer wins if they have 6 or more wins
	declarerWin := declarerPoints >= 6

	var (
		declarerpoints int
		playerpoints   int
	)

	if declarerWin {
		declarerpoints = bettingpool
		playerpoints = -1
	} else {
		declarerpoints = -bettingpool / (len(g.players) - 1)
		playerpoints = bettingpool / (len(g.players) - 1)
	}

	declarerres, err := json.Marshal(struct {
		IsDeclarerWin bool `json:"isdeclarerwin"`
		Points        int  `json:"points"`
	}{
		IsDeclarerWin: declarerWin,
		Points:        declarerpoints,
	})

	if err != nil {
		return errors.New("player response marhals went wrong")
	}

	declarer := Event{
		Type:    "game_end",
		GameId:  g.id,
		Payload: declarerres,
	}

	playerres, err := json.Marshal(struct {
		IsDeclarerWin bool `json:"isdeclarerwin"`
		Points        int  `json:"points"`
	}{
		IsDeclarerWin: declarerWin,
		Points:        playerpoints,
	})

	if err != nil {
		return errors.New("player response marhals went wrong")
	}
	player := Event{
		Type:    "game_end",
		GameId:  g.id,
		Payload: playerres,
	}

	for _, p := range g.players {
		if p == g.declarer {
			p.egress <- declarer
		} else {
			p.egress <- player
		}
	}

	return nil
}

/* ---------  Game Exit Event  ------------------------------------- */

func GameExit(event Event, c *Client) error {
	log.Printf("Player requesting to exit game %d", event.GameId)

	var game *Game
	for g := range c.manager.games {
		if g.id == event.GameId {
			game = g
			break
		}
	}

	if game == nil {
		return errors.New("game not found")
	}

	// Send game_interrupted event to all players
	err := SendGameInterrupted(game)
	if err != nil {
		return err
	}

	// Delete the game
	c.manager.deleteGame(game)

	return nil
}

func SendGameInterrupted(game *Game) error {
	log.Printf("Sending game interrupted to all players")

	event := Event{
		Type:   "game_interrupted",
		GameId: game.id,
	}

	for _, p := range game.players {
		p.egress <- event
	}

	return nil
}

/* ---------  Play Again Event  ------------------------------------- */

func PlayAgain(event Event, c *Client) error {
	log.Printf("Player %s requesting play again for game %d", c.username, event.GameId)

	var ev PlayAgainEvent
	if err := json.Unmarshal(event.Payload, &ev); err != nil {
		return errors.New("invalid play_again payload")
	}

	var game *Game
	for g := range c.manager.games {
		if g.id == event.GameId {
			game = g
			break
		}
	}

	if game == nil {
		return errors.New("game not found")
	}

	// Find player index
	playerIdx := -1
	for i, p := range game.players {
		if p == c {
			playerIdx = i
			break
		}
	}

	if playerIdx == -1 {
		return errors.New("player not in game")
	}

	game.Lock()
	game.playAgainVotes[playerIdx] = ev.Vote
	votes := make(map[int]bool)
	for k, v := range game.playAgainVotes {
		votes[k] = v
	}
	game.Unlock()

	log.Printf("Player %d voted %v for play again. Total votes: %d/%d", playerIdx, ev.Vote, len(votes), len(game.players))

	// Send vote status to all players
	err := SendPlayAgainStatus(game, votes)
	if err != nil {
		return err
	}

	// Check if all players have voted
	if len(votes) == len(game.players) {
		// Check if all voted yes
		allYes := true
		for _, vote := range votes {
			if !vote {
				allYes = false
				break
			}
		}

		if allYes {
			// All players want to play again - restart the game
			log.Printf("All players voted yes, restarting game %d", game.id)
			err := SendPlayAgainRestart(game)
			if err != nil {
				return err
			}

			// Reset game state
			game.Lock()
			// Rotate the starting bidder
			game.startingBidderIndex = (game.startingBidderIndex + 1) % len(game.players)
			game.cards = DealCards()
			game.tallon = game.cards[30:32]
			game.turncounter = 0
			game.starterplayer = game.players[game.startingBidderIndex]
			game.roundcards = make(CardList, 0)
			game.points = []int{0, 0, 0}
			game.bettingpool = 0
			game.gamecolor = ""
			game.declarer = nil
			game.biddingActive = false
			game.currentBidder = game.startingBidderIndex
			game.consecutivePasses = 0
			game.passesAfterDeclarer = 0
			game.gamePhase = "bidding"
			game.playAgainVotes = make(map[int]bool)
			game.Unlock()

			// Deal cards to all players
			for _, player := range game.players {
				cards, err := game.SetHands(player)
				if err != nil {
					log.Printf("Error dealing cards to player: %v", err)
					return err
				}

				// Send cards to player
				payload, err := json.Marshal(struct {
					Cards []int `json:"cards"`
				}{
					Cards: cards,
				})

				if err != nil {
					log.Printf("Error marshaling cards: %v", err)
					return err
				}

				event := Event{
					Type:    "set_cards",
					GameId:  game.id,
					Payload: payload,
				}

				player.egress <- event
			}

			// Start new game
			err = SendGameStart(game)
			if err != nil {
				return err
			}

		} else {
			// At least one player voted no - send denial and delete game
			log.Printf("At least one player voted no, ending game %d", game.id)
			err := SendPlayAgainDenied(game)
			if err != nil {
				return err
			}
			c.manager.deleteGame(game)
		}
	}

	return nil
}

func SendPlayAgainStatus(game *Game, votes map[int]bool) error {
	payload, err := json.Marshal(struct {
		Votes map[int]bool `json:"votes"`
	}{
		Votes: votes,
	})

	if err != nil {
		return errors.New("failed to marshal play_again_status payload")
	}

	event := Event{
		Type:    "play_again_status",
		GameId:  game.id,
		Payload: payload,
	}

	for _, p := range game.players {
		p.egress <- event
	}

	return nil
}

func SendPlayAgainDenied(game *Game) error {
	event := Event{
		Type:   "play_again_denied",
		GameId: game.id,
	}

	for _, p := range game.players {
		p.egress <- event
	}

	return nil
}

func SendPlayAgainRestart(game *Game) error {
	event := Event{
		Type:   "play_again_restart",
		GameId: game.id,
	}

	for _, p := range game.players {
		p.egress <- event
	}

	return nil
}

/* ---------  Game Initialization Event  ------------------------------------- */

func GameInit(c *Client) (int, []int, int, error) {
	log.Printf("GAMEINIT")

	for i := range c.manager.games {
		if i.isFull {
			continue
		}

		idx, err := i.addPlayer(c)
		if err != nil {
			return 0, []int{}, 0, err
		}
		id := i.id
		cards, err := i.SetHands(c)

		if err != nil {
			return id, cards, idx, err
		}

		return id, cards, idx, nil
	}

	id := 0
	for {
		id = generateID()
		if !c.manager.gameIdExist(id) {
			break
		}
	}

	game := NewGame(id)
	c.manager.addGame(game)

	idx, err := game.addPlayer(c)
	if err != nil {
		return id, []int{}, idx, err
	}
	cards, err := game.SetHands(c)

	if err != nil {
		return id, cards, idx, err
	}

	return id, cards, idx, nil
}

func SendGameInit(event Event, c *Client) error {
	id, cards, index, err := GameInit(c)
	//error
	if err != nil {
		log.Printf("error during gameinit %v", err)
		return err
	}
	log.Printf("event: id=%v cards=%v", id, cards)

	// Find the game to get player names
	var game *Game
	for g := range c.manager.games {
		if g.id == id {
			game = g
			break
		}
	}

	// Build player names map
	playerNames := make(map[int]string)
	if game != nil {
		for i, player := range game.players {
			playerNames[i] = player.username
		}
	}

	payload, err := json.Marshal(struct {
		Cards       []int          `json:"cards"`
		Index       int            `json:"index"`
		PlayerNames map[int]string `json:"playerNames"`
	}{
		Cards:       cards,
		Index:       index,
		PlayerNames: playerNames,
	})

	if err != nil {
		log.Printf("error during converting cards to payload %v", err)
		return err
	}

	response := Event{
		Type:    "game_init",
		GameId:  id,
		Payload: payload,
	}

	c.egress <- response

	// Find the game and notify all players about player count
	if game == nil {
		for g := range c.manager.games {
			if g.id == id {
				game = g
				break
			}
		}
	}

	if game != nil {
		err := SendPlayerJoined(game)
		if err != nil {
			log.Printf("error sending player joined: %v", err)
		}

		// If game is full, start the game
		if game.isFull {
			err := SendGameStart(game)
			if err != nil {
				log.Printf("error sending game start: %v", err)
			}

			// Send your_turn to the first player
			err = SendNextPlayer(0, game)
			if err != nil {
				log.Printf("error sending first turn: %v", err)
			}
		}
	}

	return nil

}

func SendPlayerJoined(game *Game) error {

	playerNames := make(map[int]string)
	for i, player := range game.players {
		playerNames[i] = player.username
	}

	payload, err := json.Marshal(struct {
		PlayerNames map[int]string `json:"playernames"`
	}{
		PlayerNames: playerNames,
	})

	if err != nil {
		return err
	}

	event := Event{
		Type:    "player_joined",
		GameId:  game.id,
		Payload: payload,
	}

	for _, p := range game.players {
		p.egress <- event
	}

	return nil
}

func SendGameStart(game *Game) error {
	log.Printf("Game starting with %d players", len(game.players))

	// Initialize bidding phase
	game.StartBidding()

	event := Event{
		Type:   "game_start",
		GameId: game.id,
	}

	for _, p := range game.players {
		p.egress <- event
	}

	// Send first bidder notification
	err := SendNextBidder(game)
	if err != nil {
		return err
	}

	return nil
}

func generateID() int {
	id := 0
	for i := 0; i < 10; i++ {
		digit := rand.Intn(10)
		id = id*10 + digit
	}
	return id
}

func (m *Manager) gameIdExist(id int) bool {
	for game := range m.games {
		if game.id == id {
			return true
		}
	}
	return false
}

/* ------------------------------------- Bidding Phase ------------------------------------- */

func BidAction(event Event, c *Client) error {
	var ev BidActionEvent
	if err := json.Unmarshal(event.Payload, &ev); err != nil {
		return errors.New("invalid bid_action payload")
	}

	log.Printf("BidAction: Player %s, Action: %s, Color: %s", c.username, ev.Action, ev.Color)

	// Find the game
	var game *Game
	for g := range c.manager.games {
		if g.id == event.GameId {
			game = g
			break
		}
	}

	if game == nil {
		return errors.New("game not found")
	}

	// Process the bid
	status, err := game.HandleBid(c, ev.Action, ev.Color)
	if err != nil {
		log.Printf("Bid error: %v", err)
		return err
	}

	log.Printf("Bid processed. Status: %s", status)

	// Send bid result to all players
	err = SendBidResult(game, c, ev.Action, ev.Color, status)
	if err != nil {
		log.Printf("SendBidResult error: %v", err)
		return err
	}

	// Handle different statuses
	switch status {
	case "continue_bidding":
		// Send next bidder notification
		log.Printf("Sending next bidder notification")
		err = SendNextBidder(game)
		if err != nil {
			log.Printf("SendNextBidder error: %v", err)
			return err
		}
	case "talon_exchange":
		// Send talon to declarer
		err = SendTalonToDeclarer(game)
		if err != nil {
			return err
		}
	case "game_closed":
		// Game closed, no declarer - end game
		err = SendGameClosed(game)
		if err != nil {
			return err
		}
		c.manager.deleteGame(game)
	default:
		return errors.New("unknown bidding status")

	}

	return nil
}

func SendBidResult(game *Game, bidder *Client, action string, color string, status string) error {
	var bidderIndex int
	for i, player := range game.players {
		if player == bidder {
			bidderIndex = i
			break
		}
	}

	payload, err := json.Marshal(struct {
		BidderIndex int    `json:"bidderindex"`
		Action      string `json:"action"`
		Color       string `json:"color"`
		Status      string `json:"status"`
		BettingPool int    `json:"bettingpool"`
	}{
		BidderIndex: bidderIndex,
		Action:      action,
		Color:       color,
		Status:      status,
		BettingPool: game.bettingpool,
	})

	if err != nil {
		return errors.New("bid_result marshal error")
	}

	event := Event{
		Type:    "bid_result",
		GameId:  game.id,
		Payload: payload,
	}

	for _, p := range game.players {
		p.egress <- event
	}

	return nil
}

func SendNextBidder(game *Game) error {
	currentBidder := game.GetCurrentBidder()
	var bidderIndex int
	for i, player := range game.players {
		if player == currentBidder {
			bidderIndex = i
			break
		}
	}

	log.Printf("SendNextBidder: Next bidder is player %d (%s)", bidderIndex, currentBidder.username)

	payload, err := json.Marshal(struct {
		BidderIndex int `json:"bidderindex"`
	}{
		BidderIndex: bidderIndex,
	})

	if err != nil {
		return errors.New("next_bidder marshal error")
	}

	event := Event{
		Type:    "next_bidder",
		GameId:  game.id,
		Payload: payload,
	}

	for _, p := range game.players {
		p.egress <- event
	}

	return nil
}

func SendTalonToDeclarer(game *Game) error {
	game.Lock()
	declarer := game.declarer
	tallon := game.tallon

	// Add talon cards to declarer's hand immediately
	for i := range tallon {
		tallon[i].Hand = declarer
	}
	game.Unlock()

	if declarer == nil {
		return errors.New("no declarer found")
	}

	// Get talon card IDs
	talonCardIds := make([]int, len(tallon))
	for i, card := range tallon {
		talonCardIds[i] = card.ID
	}

	payload, err := json.Marshal(struct {
		TalonCards []int `json:"taloncards"`
	}{
		TalonCards: talonCardIds,
	})

	if err != nil {
		return errors.New("talon marshal error")
	}

	event := Event{
		Type:    "talon_received",
		GameId:  game.id,
		Payload: payload,
	}

	// Send only to declarer
	declarer.egress <- event

	// Notify other players that talon exchange is happening
	notifyEvent := Event{
		Type:   "talon_exchange_phase",
		GameId: game.id,
	}

	for _, p := range game.players {
		if p != declarer {
			p.egress <- notifyEvent
		}
	}

	return nil
}

func SendGameClosed(game *Game) error {
	event := Event{
		Type:   "game_closed",
		GameId: game.id,
	}

	for _, p := range game.players {
		p.egress <- event
	}
	return nil
}

func TalonExchange(event Event, c *Client) error {
	var ev TalonExchangeEvent
	if err := json.Unmarshal(event.Payload, &ev); err != nil {
		return errors.New("invalid talon_exchange payload")
	}

	log.Printf("TalonExchange: Player %s discarding cards: %v", c.username, ev.DiscardCards)

	// Find the game
	var game *Game
	for g := range c.manager.games {
		if g.id == event.GameId {
			game = g
			break
		}
	}

	if game == nil {
		return errors.New("game not found")
	}

	// Process talon exchange
	err := game.HandleTalonExchange(c, ev.DiscardCards)
	if err != nil {
		log.Printf("Talon exchange error: %v", err)
		return err
	}

	log.Printf("Talon exchange successful, moving to playing phase")

	// Notify all players that talon exchange is complete
	err = SendTalonExchangeComplete(game)
	if err != nil {
		return err
	}

	// Find declarer's index
	var declarerIndex int
	for i, player := range game.players {
		if player == game.declarer {
			declarerIndex = i
			break
		}
	}

	// Set starter player to declarer
	game.Lock()
	game.starterplayer = game.declarer
	game.Unlock()

	// Start the playing phase - declarer starts
	log.Printf("Sending first turn to declarer (player %d)", declarerIndex)
	err = SendNextPlayer(declarerIndex, game)
	if err != nil {
		return err
	}

	return nil
}

func SendTalonExchangeComplete(game *Game) error {
	event := Event{
		Type:   "talon_exchange_complete",
		GameId: game.id,
	}

	for _, p := range game.players {
		p.egress <- event
	}

	return nil
}
