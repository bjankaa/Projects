export default {
    setSocketConnection(state, payload) {
        state.socket = payload.socket;
        state.isConnected = payload.isConnected;
    },
    setGameId(state, payload) {
        state.gameID = payload;
    },
    setCards(state, payload){
        state.cards = payload;
    },

    deleteCard(state,payload){
        state.cards = state.cards.filter(id => id !== payload);
    },

    resetWebSocket(state) {
        state.socket = null;
        state.isConnected = false;
        state.gameID = 0;
        state.cards = [];
        state.yourturn = false;
        state.roundcounter = 0;
        state.idx = -1;
        state.yourPoints = 0;
        state.bettingPool = 0;
        state.roundWinner = -1;
        state.playedCards = [];
        state.gameStarted = false;
        state.playerNames = {};
        state.gamePhase = 'waiting';
        state.currentBidder = -1;
        state.isYourBid = false;
        state.lastBid = null;
        state.gameColor = '';
        state.declarerIndex = -1;
        state.talonCards = [];
        state.gameResult = null;
        state.playAgainVotes = {};
    },
    //setting turn to true or false
    setYourTurn(state, payload){
        state.yourturn = payload;
    },

    setRoundCounter(state, payload){
        state.roundcounter = payload;
    },
    setIndex(state, payload){
        state.idx = payload;
        // First round starter is player 0
        if (state.roundWinner === -1) {
            state.roundWinner = 0;
        }
    },

    setPlayerNames(state, payload) {
        state.playerNames = payload;
    },

    // Handle card played by another player
    cardPlayed(state, payload) {
        console.log(`Card ${payload.cardId} played by player ${payload.playerIndex}`);
        state.playedCards.push({
            cardId: payload.cardId,
            playerIndex: payload.playerIndex
        });
    },

    // Handle round won by current player
    roundWon(state, payload) {
        console.log(`You won the round with ${payload.points} points!`);
        state.yourPoints = payload.points;
    },

    // Handle round result (winner announcement)
    setRoundResult(state, payload) {
        console.log(`Round won by player ${payload.winnerIndex}`);
        state.roundWinner = payload.winnerIndex;
        state.roundcounter = state.roundcounter + 1;
        state.playedCards = [];
    },

    // Handle player joined
    playerJoined(state, payload) {
        state.playerNames = payload.playerNames;
    },

    // Handle game start
    gameStart(state) {
        state.gameStarted = true;
        state.gamePhase = 'bidding';
    },

    // Reset game data but keep socket connection
    resetGameData(state) {
        state.gameID = 0;
        state.cards = [];
        state.yourturn = false;
        state.roundcounter = 0;
        state.idx = -1;
        state.yourPoints = 0;
        state.bettingPool = 0;
        state.roundWinner = -1;
        state.playedCards = [];
        state.gameStarted = false;
        state.gamePhase = 'waiting';
        state.currentBidder = -1;
        state.isYourBid = false;
        state.lastBid = null;
        state.gameColor = '';
        state.declarerIndex = -1;
        state.talonCards = [];
    },

    // Handle game interrupted (triggers component watcher)
    gameInterrupted(state) {
        // This mutation exists just to trigger component watchers
        // The actual message is handled in the component
    },

    // Handle bid result
    bidResult(state, payload) {
        console.log('bidResult mutation:', payload);
        state.lastBid = {
            bidderIndex: payload.bidderIndex,
            action: payload.action,
            color: payload.color,
            status: payload.status
        };
        
        // Update betting pool if provided
        if (payload.bettingpool !== undefined) {
            state.bettingPool = payload.bettingpool;
        }
        
        // Store game color if someone declared
        if (payload.action === 'declare' && payload.color) {
            state.gameColor = payload.color;
            state.declarerIndex = payload.bidderIndex;
        }
    },

    // Handle next bidder
    nextBidder(state, payload) {
        console.log('nextBidder mutation:', payload, 'myIdx:', state.idx);
        state.currentBidder = payload.bidderIndex;
        state.isYourBid = (payload.bidderIndex === state.idx);
        console.log('isYourBid now:', state.isYourBid);
    },

    // Handle talon received
    talonReceived(state, payload) {
        console.log('talonReceived mutation:', payload);
        state.talonCards = payload.talonCards;
        // Add talon cards to player's hand
        state.cards = [...state.cards, ...payload.talonCards];
        // Set game phase to talon_exchange for the declarer
        state.gamePhase = 'talon_exchange';
    },

    // Handle talon exchange phase
    talonExchangePhase(state) {
        state.gamePhase = 'talon_exchange';
    },

    // Handle talon exchange complete
    talonExchangeComplete(state) {
        console.log('Talon exchange complete - cards in hand:', state.cards.length, 'cards:', state.cards);
        state.gamePhase = 'playing';
        state.talonCards = [];
    },

    // Handle game closed
    gameClosed(state) {
        state.gamePhase = 'closed';
    },

    // Handle game end
    gameEnd(state, payload) {
        console.log('gameEnd mutation:', payload);
        state.gamePhase = 'ended';
        state.gameResult = {
            isDeclarerWin: payload.isDeclarerWin,
            points: payload.points
        };
    },

    // Update play again votes
    updatePlayAgainVotes(state, votes) {
        state.playAgainVotes = votes;
    },

    // Reset for play again (keep connection, reset game state)
    resetForPlayAgain(state) {
        state.cards = [];
        state.yourturn = false;
        state.roundcounter = 0;
        state.yourPoints = 0;
        state.roundWinner = -1;
        state.playedCards = [];
        state.gamePhase = 'waiting';
        state.currentBidder = -1;
        state.isYourBid = false;
        state.lastBid = null;
        state.gameColor = '';
        state.declarerIndex = -1;
        state.talonCards = [];
        state.gameResult = null;
        state.playAgainVotes = {};
    }

}