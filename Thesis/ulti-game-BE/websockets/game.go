package websockets

import (
	"errors"
	"log"
	"math/rand"
	"sync"
)

type Game struct {
	id      int
	players []*Client
	cards   CardList
	isFull  bool

	turncounter        int
	starterplayer      *Client
	roundcards         CardList
	cardsPlayedInRound int // tracks how many cards have been played in current round
	points             []int
	bettingpool        int
	gamecolor          string
	declarer           *Client
	tallon             CardList

	// Bidding phase fields
	biddingActive       bool
	currentBidder       int // index of player whose turn it is to bid
	startingBidderIndex int // tracks which player starts bidding (rotates each game)
	consecutivePasses   int
	passesAfterDeclarer int
	gamePhase           string // "bidding", "talon_exchange", "playing"

	// Play again voting
	playAgainVotes map[int]bool // player index -> vote (true/false)

	sync.RWMutex
}

// Color strength values for comparison
var colorStrength = map[string]int{
	"tok":   1,
	"makk":  2,
	"zold":  3,
	"piros": 4,
}

type GameList map[*Game]bool

func NewGame(id int) *Game {
	g := &Game{
		id:                  id,
		players:             []*Client{},
		cards:               DealCards(),
		isFull:              false,
		turncounter:         0,
		starterplayer:       nil,
		roundcards:          make(CardList, 3),
		cardsPlayedInRound:  0,
		points:              []int{0, 0, 0},
		bettingpool:         0,
		gamecolor:           "",
		declarer:            nil,
		biddingActive:       false,
		currentBidder:       0,
		startingBidderIndex: 0,
		consecutivePasses:   0,
		passesAfterDeclarer: 0,
		gamePhase:           "bidding",
		playAgainVotes:      make(map[int]bool),
	}
	// Set tallon to last 2 cards
	g.tallon = g.cards[30:32]
	return g
}

func (g *Game) addPlayer(client *Client) (int, error) {
	g.Lock()
	defer g.Unlock()

	if len(g.players) >= 3 {
		return -1, errors.New("too manny players")
	}

	if len(g.players) == 0 {
		g.starterplayer = client
	}

	g.players = append(g.players, client)

	if len(g.players) >= 3 {
		g.isFull = true
	}

	return len(g.players) - 1, nil
}

/* ------------------------------------- Card operations ------------------------------------- */

// CARDS randomize list order by players
func DealCards() CardList {
	cards := make(CardList, len(CardDeck))
	copy(cards, CardDeck)

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}

// From the shuffled cards assign cards to a new player
func (g *Game) SetHands(c *Client) ([]int, error) {
	g.Lock()
	defer g.Unlock()

	count := len(g.players)

	if count < 1 || count > 3 {
		err := errors.New("there is a problem with the number of players")
		return nil, err
	}

	start := (count - 1) * 10
	end := start + 10

	ids := make([]int, 10)
	for i := range ids {
		ids[i] = g.cards[start+i].ID
	}

	for i := start; i < end; i++ {
		g.cards[i].Hand = c
	}

	if count == 3 {
		g.isFull = true
	}

	return ids, nil
}

// checking if the card that was played is in the hands of the player
func (g *Game) playedCardCheck(cardId int, c *Client) (Card, error) {

	for i := range g.cards {
		if g.cards[i].ID == cardId && g.cards[i].Hand == c {
			g.cards[i].Hand = nil
			return g.cards[i], nil
		}
	}

	return Card{}, errors.New("the card is not in the hand of the player wrong move")
}

/* ------------------------------------- Playing ------------------------------------- */

func (g *Game) RoundHandler(c *Client, cardId int) (int, *Client, string) {
	g.Lock()
	defer g.Unlock()

	err := g.currentPlayerCheck(c)

	if err != nil {
		log.Println("current player errror: ", err)
		return -1, nil, ""
	}

	card, err := g.playedCardCheck(cardId, c)

	if err != nil {
		log.Println("Played card error", err)
		return -1, nil, ""
	}

	for i := range g.players {
		if g.players[i] == c {
			g.roundcards[i] = card
			break
		}
	}

	g.cardsPlayedInRound++
	g.turncounter++
	if g.turncounter >= 30 {
		_, err := g.EvaluateRound()
		if err != nil {
			log.Printf("Something went wrong during evaluation")
			return -1, nil, ""
		}
		return -1, nil, "end"
	}

	nextplayer, err := g.NextPlayer()

	if err != nil {
		log.Println("couldn't idetify next player")
		return -1, nil, ""
	}

	if g.cardsPlayedInRound >= len(g.players) {
		roundwinner, err := g.EvaluateRound()
		if err != nil {
			log.Printf("Something went wrong during evaluation")
			return -1, nil, ""
		}
		g.starterplayer = roundwinner
		g.roundcards = make(CardList, 3)
		g.cardsPlayedInRound = 0

		// Find the round winner's index to return as next player
		var winnerIndex int
		for i := range g.players {
			if g.players[i] == roundwinner {
				winnerIndex = i
				break
			}
		}

		return winnerIndex, roundwinner, ""
	}

	return nextplayer, nil, ""

}

// ---------------------------------------------------------------------------------------------
func (g *Game) NextPlayer() (int, error) {

	index := -1

	for i := range g.players {
		if g.players[i] == g.starterplayer {
			index = i
			break
		}
	}
	if index == -1 {
		return -1, errors.New("there is no starter player")
	}

	nextidx := (g.cardsPlayedInRound + index) % len(g.players)

	//nextp := g.players[nextidx]

	if g.players[nextidx] == nil {
		return -1, errors.New("no next player")
	}

	return nextidx, nil
}

// ---------------------------------------------------------------------------------------------
func (g *Game) currentPlayerCheck(c *Client) error {

	index := -1

	for i := range g.players {
		if g.players[i] == g.starterplayer {
			index = i
			break
		}
	}
	if index == -1 {
		log.Printf("DEBUG: starterplayer is nil or not found")
		return errors.New("there is no starter player")
	}

	offset := g.cardsPlayedInRound
	currentp := (index + offset) % len(g.players)

	log.Printf("DEBUG: starter index=%d, cardsPlayedInRound=%d, offset=%d, currentp=%d, clickingPlayer=%s, expectedPlayer=%s",
		index, g.cardsPlayedInRound, offset, currentp, c.username, g.players[currentp].username)

	if g.players[currentp] != c {
		return errors.New("there is no player like this")
	}
	return nil
}

/* ------------------------------------- Evaluate Round  ------------------------------------- */
func (g *Game) EvaluateRound() (*Client, error) {

	// Find the starter player's index
	starterIndex := -1
	for i := range g.players {
		if g.players[i] == g.starterplayer {
			starterIndex = i
			break
		}
	}
	if starterIndex == -1 {
		return nil, errors.New("starter player not found")
	}

	startcolor := g.roundcards[starterIndex].Color
	highestcardidx := g.roundcards[starterIndex].ID
	highestcardplayer := g.starterplayer
	hasGameColor := false

	for i := range g.roundcards {
		if g.roundcards[i].Color == g.gamecolor && startcolor != g.gamecolor {
			// This is a trump card (gamecolor)
			if !hasGameColor || g.roundcards[i].ID > highestcardidx {
				hasGameColor = true
				highestcardidx = g.roundcards[i].ID
				highestcardplayer = g.players[i]
			}
		} else if !hasGameColor && g.roundcards[i].Color == startcolor && g.roundcards[i].ID > highestcardidx {
			// Only update if no gamecolor has been played yet
			highestcardidx = g.roundcards[i].ID
			highestcardplayer = g.players[i]
		}
	}

	for i := range g.players {
		if g.players[i] == highestcardplayer {
			g.points[i]++
		}
	}

	return highestcardplayer, nil
}

/* ------------------------------------- Bidding Phase ------------------------------------- */

// StartBidding initializes the bidding phase
func (g *Game) StartBidding() {
	g.Lock()
	defer g.Unlock()

	g.biddingActive = true
	g.currentBidder = g.startingBidderIndex
	g.consecutivePasses = 0
	g.passesAfterDeclarer = 0
	g.gamePhase = "bidding"
	g.bettingpool = 0
}

// HandleBid processes a bid from a player
// action can be "pass" or "declare"
// color is the chosen color if declaring (tok, makk, zold, piros)
func (g *Game) HandleBid(c *Client, action string, color string) (string, error) {
	g.Lock()
	defer g.Unlock()

	// Verify it's this player's turn to bid
	if g.players[g.currentBidder] != c {
		return "", errors.New("not your turn to bid")
	}

	if action == "pass" {
		g.consecutivePasses++

		// If declarer exists, count passes after declarer
		if g.declarer != nil {
			g.passesAfterDeclarer++

			// If 2 consecutive passes after declarer, move to talon phase
			if g.passesAfterDeclarer >= 2 {
				g.gamePhase = "talon_exchange"
				g.biddingActive = false
				return "talon_exchange", nil
			}
		} else {
			// No declarer yet, check if 6 consecutive passes
			if g.consecutivePasses >= 6 {
				return "game_closed", errors.New("no declarer after 6 passes")
			}
		}

		// Move to next bidder
		g.currentBidder = (g.currentBidder + 1) % len(g.players)
		return "continue_bidding", nil
	}

	if action == "declare" {
		// Validate color
		if _, exists := colorStrength[color]; !exists {
			return "", errors.New("invalid color choice")
		}

		// Check if this is a valid one-up
		if g.declarer != nil {
			currentStrength := colorStrength[g.gamecolor]
			newStrength := colorStrength[color]

			// Can only one-up with stronger color OR same piros
			if newStrength < currentStrength {
				return "", errors.New("must choose stronger color to one-up")
			}
			if newStrength == currentStrength && color != "piros" {
				return "", errors.New("can only choose same color if it's piros")
			}
		}

		// Set new declarer
		g.declarer = c
		g.gamecolor = color
		g.bettingpool++
		g.consecutivePasses = 0
		g.passesAfterDeclarer = 0

		// Move to next bidder
		g.currentBidder = (g.currentBidder + 1) % len(g.players)
		return "continue_bidding", nil
	}

	return "", errors.New("invalid action")
}

// HandleTalonExchange processes talon exchange from declarer
func (g *Game) HandleTalonExchange(c *Client, discardCards []int) error {
	g.Lock()
	defer g.Unlock()

	// Only declarer can exchange
	if c != g.declarer {
		return errors.New("only declarer can exchange talon")
	}

	if g.gamePhase != "talon_exchange" {
		return errors.New("not in talon exchange phase")
	}

	// Validate discard count (1 or 2 cards)
	if len(discardCards) < 1 || len(discardCards) > 2 {
		return errors.New("must discard 1 or 2 cards")
	}

	// Verify cards belong to declarer and remove them from hand
	for _, cardID := range discardCards {
		found := false
		for i := range g.cards {
			if g.cards[i].ID == cardID && g.cards[i].Hand == c {
				// Remove from hand (these cards go to talon)
				g.cards[i].Hand = nil
				found = true
				break
			}
		}
		if !found {
			return errors.New("card not in your hand")
		}
	}

	// Talon cards were already added to declarer's hand in SendTalonToDeclarer
	// No need to add them again here

	// Move to playing phase
	g.gamePhase = "playing"
	return nil
}

/* ------------------------------------- Helper Functions ------------------------------------- */

func (g *Game) GetCurrentBidder() *Client {
	g.RLock()
	defer g.RUnlock()
	return g.players[g.currentBidder]
}

func (g *Game) GetGamePhase() string {
	g.RLock()
	defer g.RUnlock()
	return g.gamePhase
}
