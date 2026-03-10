import mutations from "./mutations";
import actions from "./actions";
import getters from "./getters";

export default{
    namespaced: true,
    state() {
        return {
            socket: null,
            isConnected: false,
            gameID: 0,
            cards: [],
            yourturn: false,
            roundcounter: 0,
            idx:-1,
            yourPoints: 0,
            bettingPool: 0,
            roundWinner: -1,
            playedCards: [],
            gameStarted: false,
            playerNames: {}, // Map of playerIndex -> playerName
            // Bidding phase
            gamePhase: 'waiting', // 'waiting', 'bidding', 'talon_exchange', 'playing', 'closed', 'ended'
            currentBidder: -1,
            isYourBid: false,
            lastBid: null,
            gameColor: '',
            declarerIndex: -1,
            talonCards: [],
            gameResult: null, // { isDeclarerWin, points }
            playAgainVotes: {}, // Map of playerIndex -> vote (true/false)
        };
    },
    getters,
    actions,
    mutations
}