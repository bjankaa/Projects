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
        state.roundColor = '';
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
        if (state.roundWinner === -1) {
            state.roundWinner = 0;
        }
    },

    setPlayerNames(state, payload) {
        state.playerNames = payload;
    },

    cardPlayed(state, payload) {
        console.log(`Card ${payload.cardId} played by player ${payload.playerIndex}`);
        // If this is the first card, set the round color
        if (state.playedCards.length === 0 && payload.color) {
            state.roundColor = payload.color;
        }
        state.playedCards.push({
            cardId: payload.cardId,
            playerIndex: payload.playerIndex,
            color: payload.color
        });
    },

    roundWon(state, payload) {
        console.log(`You won the round with ${payload.points} points!`);
        state.yourPoints = payload.points;
    },

    setRoundResult(state, payload) {
        console.log(`Round won by player ${payload.winnerIndex}`);
        state.roundWinner = payload.winnerIndex;
        state.roundcounter = state.roundcounter + 1;
        state.playedCards = [];
        state.roundColor = '';
    },

    playerJoined(state, payload) {
        state.playerNames = payload.playerNames;
    },

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
        state.roundColor = '';
        state.gameStarted = false;
        state.gamePhase = 'waiting';
        state.currentBidder = -1;
        state.isYourBid = false;
        state.lastBid = null;
        state.gameColor = '';
        state.declarerIndex = -1;
        state.talonCards = [];
    },

    gameInterrupted(state) {
        // This mutation exists just to trigger component watchers
        // The actual message is handled in the component
    },

    bidResult(state, payload) {
        console.log('bidResult mutation:', payload);
        state.lastBid = {
            bidderIndex: payload.bidderIndex,
            action: payload.action,
            color: payload.color,
            status: payload.status
        };
        
        if (payload.bettingpool !== undefined) {
            state.bettingPool = payload.bettingpool;
        }
        
        if (payload.action === 'declare' && payload.color) {
            state.gameColor = payload.color;
            state.declarerIndex = payload.bidderIndex;
        }
    },

    nextBidder(state, payload) {
        console.log('nextBidder mutation:', payload, 'myIdx:', state.idx);
        state.currentBidder = payload.bidderIndex;
        state.isYourBid = (payload.bidderIndex === state.idx);
        console.log('isYourBid now:', state.isYourBid);
    },


    talonReceived(state, payload) {
        console.log('talonReceived mutation:', payload);
        state.talonCards = payload.talonCards;
        state.cards = [...state.cards, ...payload.talonCards];
        state.gamePhase = 'talon_exchange';
    },


    talonExchangePhase(state) {
        state.gamePhase = 'talon_exchange';
    },

    talonExchangeComplete(state) {
        console.log('Talon exchange complete - cards in hand:', state.cards.length, 'cards:', state.cards);
        state.gamePhase = 'playing';
        state.talonCards = [];

        state.yourturn = false;
    },

    gameClosed(state) {
        state.gamePhase = 'closed';
    },

    gameEnd(state, payload) {
        console.log('gameEnd mutation:', payload);
        state.gamePhase = 'ended';
        state.gameResult = {
            isDeclarerWin: payload.isDeclarerWin,
            points: payload.points
        };
    },

    updatePlayAgainVotes(state, votes) {
        state.playAgainVotes = votes;
    },

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