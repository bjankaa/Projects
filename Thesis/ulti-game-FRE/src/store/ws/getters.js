export default {
    socket(state){
        return state.socket;
    },

    isConnected(state){
        return state.isConnected;
    },

    gameID(state){
        return state.gameID;
    },
    cards(state){
        return state.cards;
    },
    yourturn(state){
        return state.yourturn;
    },
    roundcounter(state){
        return state.roundcounter;
    },
    idx(state){
        return state.idx;
    },
    yourPoints(state){
        return state.yourPoints;
    },
    bettingPool(state){
        return state.bettingPool;
    },
    roundWinner(state){
        return state.roundWinner;
    },
    playedCards(state){
        return state.playedCards;
    },
    gameStarted(state){
        return state.gameStarted;
    },
    playerNames(state){
        return state.playerNames;
    },
    gamePhase(state){
        return state.gamePhase;
    },
    currentBidder(state){
        return state.currentBidder;
    },
    isYourBid(state){
        return state.isYourBid;
    },
    lastBid(state){
        return state.lastBid;
    },
    gameColor(state){
        return state.gameColor;
    },
    declarerIndex(state){
        return state.declarerIndex;
    },
    talonCards(state){
        return state.talonCards;
    },
    gameResult(state){
        return state.gameResult;
    },
    playAgainVotes(state){
        return state.playAgainVotes;
    }
}