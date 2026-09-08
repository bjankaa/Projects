package websockets

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	GameId  int             `json:"id"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventCardClicked   = "card_clicked"
	EventGameInit      = "game_init"
	EventPlayerJoined  = "player_joined"
	EventGameStart     = "game_start"
	EventGameExit      = "game_exit"
	EventBidAction     = "bid_action"
	EventTalonExchange = "talon_exchange"
	EventPlayAgain     = "play_again"
)

// json blob
type CardClickedEvent struct {
	CardId int `json:"cardid"`
}

type BidActionEvent struct {
	Action string `json:"action"` // "pass" or "declare"
	Color  string `json:"color"`  // "tok", "makk", "zold", "piros" (only if declaring)
}

type TalonExchangeEvent struct {
	DiscardCards []int `json:"discardcards"` // 1 or 2 card IDs to discard
}

type PlayAgainEvent struct {
	Vote bool `json:"vote"` // true to play again, false to exit
}

// type GameInitEvent struct {
// 	Cards CardList `json:"cards"`
// }
