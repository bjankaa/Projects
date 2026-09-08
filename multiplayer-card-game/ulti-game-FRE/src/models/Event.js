// need the gameid to identify/check if its the right game
// type to know what kind of event is this
// payload other information regarding the event

export class Event {
    constructor(type, id, payload) {
        this.type = type;
        this.id = id;
        this.payload = payload;
    }

    // handle event
    routeEvent() {
        if (this.type === undefined) {
            alert("no type field in the event");
        }

        switch (this.type) {
            case "card_clicked":
                console.log("card clicked response");
                return null;
            case "game_init":
                return {
                    gameId: this.id,
                    payload: this.payload,
                };
            case "card_played":
                return {
                    cardId: this.payload.cardid,
                    playerIndex: this.payload.playerindex
                };
            case "your_turn":
                return true;
            case "you_won":
                return {
                    points: this.payload.points
                };
            case "round_result":
                return {
                    winnerIndex: this.payload.winnerindex
                };
            case "player_joined":
                return {
                    playerNames: this.payload.playernames
                };
            case "game_start":
                return true;
            case "game_exit":
                return true;
            case "game_interrupted":
                return true;
            case "bid_result":
                return {
                    bidderIndex: this.payload.bidderindex,
                    action: this.payload.action,
                    color: this.payload.color,
                    status: this.payload.status,
                    bettingpool: this.payload.bettingpool
                };
            case "next_bidder":
                return {
                    bidderIndex: this.payload.bidderindex
                };
            case "talon_received":
                return {
                    talonCards: this.payload.taloncards
                };
            case "talon_exchange_phase":
                return true;
            case "talon_exchange_complete":
                return true;
            case "game_closed":
                return true;
            case "game_end":
                return {
                    isDeclarerWin: this.payload.isdeclarerwin,
                    points: this.payload.points
                };
            case "play_again_status":
                return {
                    votes: this.payload.votes
                };
            case "play_again_denied":
                return true;
            case "play_again_restart":
                return true;
            case "set_cards":
                return {
                    cards: this.payload.cards
                };
            default:
                alert("unsupported event type");
                return null;
        }
    }
}
