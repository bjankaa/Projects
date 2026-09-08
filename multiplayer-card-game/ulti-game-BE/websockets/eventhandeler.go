package websockets

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"math/rand"

	"exmaple.com/ulti-restapi/database"
	"exmaple.com/ulti-restapi/models"
)

/* ----------------------------  Card Clicked Event  ------------------------------------- */

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

	var color string
	switch {
	case cardid >= 1 && cardid <= 8:
		color = "tok"
	case cardid >= 9 && cardid <= 16:
		color = "makk"
	case cardid >= 17 && cardid <= 24:
		color = "zold"
	case cardid >= 25 && cardid <= 32:
		color = "piros"
	default:
		color = ""
	}

	payload, err := json.Marshal(struct {
		CardId      int    `json:"cardid"`
		PlayerIndex int    `json:"playerindex"`
		Color       string `json:"color"`
	}{
		CardId:      cardid,
		PlayerIndex: pindex,
		Color:       color,
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

/* ----------------------------  Next Player Event  ------------------------------------- */

func SendNextPlayer(nextp int, game *Game) error {
	log.Printf("Sending next player: %d", nextp)
	tonextplayer := Event{
		Type:   "your_turn",
		GameId: game.id,
	}

	notYourTurn := Event{
		Type:   "not_your_turn",
		GameId: game.id,
	}

	for i, p := range game.players {
		if i == nextp {
			p.egress <- tonextplayer
		} else {
			p.egress <- notYourTurn
		}
	}

	return nil
}

/* ----------------------------  Round End Event  ------------------------------------- */

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

/* ----------------------------  Game End Event  ------------------------------------- */

func SendEndResult(g *Game) error {
	log.Printf("End of the game")

	bettingpool := g.bettingpool

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

	split := int(math.Ceil(float64(bettingpool) / float64(2)))

	if declarerWin {
		playerpoints = -split
		declarerpoints = bettingpool
	} else {
		playerpoints = split
		declarerpoints = -split
	}

	// saving game-------------------------------------------------own function
	for i, p := range g.players {
		points := 0
		if p == g.declarer {
			points = declarerpoints
		} else {
			points = playerpoints
		}
		if points != 0 && p.userID != 0 {

			u, err := models.GetUserByID(p.userID)
			if err != nil {
				log.Printf("Failed to fetch user %d for score update: %v", p.userID, err)
			} else {
				if err := u.AddToScore(points); err != nil {
					log.Printf("Failed to update score for user %d: %v", p.userID, err)
				}
			}
		} else if p.userID == 0 {
			log.Printf("Player index %d has no userID; skipping score update", i)
		}
	}

	var p1, p2, p3 int64
	var n1, n2, n3 string
	if len(g.players) == 3 {
		p1 = g.players[0].userID
		n1 = g.players[0].username
		p2 = g.players[1].userID
		n2 = g.players[1].username
		p3 = g.players[2].userID
		n3 = g.players[2].username
	}
	var declarerID int64
	var declarerName string
	if g.declarer != nil {
		declarerID = g.declarer.userID
		declarerName = g.declarer.username
	}

	insert := database.NormalizeQuery(database.Currentdb, "INSERT INTO finished_games(player1_id, player1_name, player2_id, player2_name, player3_id, player3_name, declarer_id, declarer_name, declarer_win, declarer_points, defenders_points) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)")
	_, derr := database.Database.Exec(
		insert,
		p1, n1, p2, n2, p3, n3, declarerID, declarerName, declarerWin, declarerpoints, playerpoints,
	)
	if derr != nil {
		log.Printf("Failed to insert finished game record: %v", derr)
	}
	// end saving game----------------------------------------------
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

/* ----------------------------  Game Exit Event  ------------------------------------- */

func GameExit(event Event, c *Client) error {

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

	err := SendGameInterrupted(game)
	if err != nil {
		return err
	}

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

/* ----------------------------  Play Again Event  ------------------------------------- */

func PlayAgain(event Event, c *Client) error {
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

	err := SendPlayAgainStatus(game, votes)
	if err != nil {
		return err
	}

	// Check if all players have voted
	if len(votes) == len(game.players) {
		allYes := true
		for _, vote := range votes {
			if !vote {
				allYes = false
				break
			}
		}

		if allYes {
			// restart the game
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

/* ----------------------------  Game Init Event  ------------------------------------- */

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
	if err != nil {
		log.Printf("error during gameinit %v", err)
		return err
	}
	log.Printf("event: id=%v cards=%v", id, cards)

	var game *Game
	for g := range c.manager.games {
		if g.id == id {
			game = g
			break
		}
	}

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

	game.StartBidding()

	event := Event{
		Type:   "game_start",
		GameId: game.id,
	}

	for _, p := range game.players {
		p.egress <- event
	}

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

	status, err := game.HandleBid(c, ev.Action, ev.Color)
	if err != nil {
		log.Printf("Bid error: %v", err)
		return err
	}

	err = SendBidResult(game, c, ev.Action, ev.Color, status)
	if err != nil {
		log.Printf("SendBidResult error: %v", err)
		return err
	}

	switch status {
	case "continue_bidding":
		log.Printf("Sending next bidder notification")
		err = SendNextBidder(game)
		if err != nil {
			log.Printf("SendNextBidder error: %v", err)
			return err
		}
	case "talon_exchange":
		err = SendTalonToDeclarer(game)
		if err != nil {
			return err
		}
	case "game_closed":
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

/* ----------------------------  Talon Event  ------------------------------------- */

func SendTalonToDeclarer(game *Game) error {
	game.Lock()
	declarer := game.declarer
	tallon := game.tallon

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

	declarer.egress <- event

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

func TalonExchange(event Event, c *Client) error {
	var ev TalonExchangeEvent
	if err := json.Unmarshal(event.Payload, &ev); err != nil {
		return errors.New("invalid talon_exchange payload")
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

	err := game.HandleTalonExchange(c, ev.DiscardCards)
	if err != nil {
		log.Printf("Talon exchange error: %v", err)
		return err
	}

	err = SendTalonExchangeComplete(game)
	if err != nil {
		return err
	}

	var declarerIndex int
	for i, player := range game.players {
		if player == game.declarer {
			declarerIndex = i
			break
		}
	}

	game.Lock()
	game.starterplayer = game.declarer
	game.Unlock()

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
