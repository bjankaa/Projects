import { Event } from '../../models/Event.js';

export default{

    async connect(context) {
        console.log("Im here");

        if(context.getters.socket && context.getters.isConnected){
            console.log("Already connected");
            return;
        }

        try {
            // Get token from auth store
            const token = context.rootGetters['auth/token'];
            if (!token) {
                throw new Error('No authentication token available');
            }

            const socket = new WebSocket (`ws://localhost:3000/game?token=${token}`);

            socket.onopen = () => {
                const payload = {
                    socket: socket,
                    isConnected: true,
                };
                context.commit('setSocketConnection', payload);
                console.log('WebSocket connected successfully');
            };


        // handeling messages coming from the backend
        socket.onmessage = (ev) => {
            try {
            const eventdata = JSON.parse(ev.data);

            const event = new Event(eventdata.type, eventdata.id, eventdata.payload);

            // Validate gameId for all events except game_init
            const currentGameId = context.getters.gameID;
            if (event.type !== 'game_init' && currentGameId !== 0 && event.id !== currentGameId) {
                console.error(`Game ID mismatch! Expected: ${currentGameId}, Received: ${event.id}`);
                return;
            }

            const returned_data = event.routeEvent();
            
            switch(event.type) {
                case 'game_init':
                    if(returned_data != null) {
                        context.dispatch("gameInit", returned_data);
                    }
                    break;
                
                case 'card_played':
                    if(returned_data != null) {
                        console.log(`Player ${returned_data.playerIndex} played card ${returned_data.cardId}`);
                        context.commit('cardPlayed', returned_data);
                    }
                    break;
                
                case 'your_turn':
                    console.log('It\'s your turn!');
                    context.commit('setYourTurn', true);
                    break;
                
                case 'you_won':
                    if(returned_data != null) {
                        console.log(`You won the round! Points: ${returned_data.points}`);
                        context.commit('roundWon', returned_data);
                    }
                    break;
                
                case 'round_result':
                    if(returned_data != null) {
                        console.log(`Round won by player ${returned_data.winnerIndex}`);
                        // Delay clearing the cards so players can see the result
                        setTimeout(() => {
                            context.commit('setRoundResult', returned_data);
                        }, 2000); // 2 second delay
                    }
                    break;
                
                case 'player_joined':
                    if(returned_data != null) {
                        console.log(`Player joined: ${Object.keys(returned_data.playerNames).length}/3`);
                        context.commit('playerJoined', returned_data);
                    }
                    break;
                
                case 'game_start':
                    console.log('All players joined! Game starting...');
                    context.commit('gameStart');
                    break;
                
                case 'game_interrupted':
                    console.log('Game interrupted');
                    context.commit('resetGameData');
                    context.commit('gameInterrupted');
                    break;
                
                case 'bid_result':
                    if(returned_data != null) {
                        console.log(`Bid result: Player ${returned_data.bidderIndex} ${returned_data.action} ${returned_data.color || ''}`);
                        context.commit('bidResult', returned_data);
                    }
                    break;
                
                case 'next_bidder':
                    if(returned_data != null) {
                        console.log(`Next bidder: Player ${returned_data.bidderIndex}`);
                        context.commit('nextBidder', returned_data);
                    }
                    break;
                
                case 'talon_received':
                    if(returned_data != null) {
                        console.log('Talon received:', returned_data.talonCards);
                        context.commit('talonReceived', returned_data);
                    }
                    break;
                
                case 'talon_exchange_phase':
                    console.log('Talon exchange phase started');
                    context.commit('talonExchangePhase');
                    break;
                
                case 'talon_exchange_complete':
                    console.log('Talon exchange complete');
                    context.commit('talonExchangeComplete');
                    break;
                
                case 'game_closed':
                    console.log('Game closed - no declarer');
                    context.commit('resetGameData');
                    context.commit('gameClosed');
                    break;
                
                case 'game_end':
                    if(returned_data != null) {
                        console.log('Game ended:', returned_data);
                        context.commit('gameEnd', returned_data);
                    }
                    break;
                
                case 'play_again_status':
                    if(returned_data != null) {
                        console.log('Play again votes:', returned_data.votes);
                        context.commit('updatePlayAgainVotes', returned_data.votes);
                    }
                    break;
                
                case 'play_again_denied':
                    console.log('Play again denied, redirecting to game hub');
                    context.commit('resetGameData');
                    context.commit('resetWebSocket');
                    break;
                
                case 'play_again_restart':
                    console.log('Play again restart, resetting game');
                    context.commit('resetForPlayAgain');
                    break;
                
                case 'set_cards':
                    if(returned_data != null) {
                        console.log('Received new cards:', returned_data.cards);
                        context.commit('setCards', returned_data.cards);
                    }
                    break;
                
                default:
                    console.log('Unknown event type:', event.type);
            }
            } catch (error) {
                console.error('Error processing WebSocket message:', error);
            }
        };

        socket.onerror = (err) => {
            console.error("WebSocket error: ", err);
            context.commit('resetWebSocket');
        };

        socket.onclose = () => {
            context.commit('resetWebSocket');
            console.log('WebSocket connection closed');
        };

        } catch (error) {
            console.error('Failed to establish WebSocket connection:', error);
            context.commit('resetWebSocket');
        }
    },

    closeSocket(context) {
        const socket = context.getters.socket;

        if (socket && socket.readyState !== WebSocket.CLOSED) {
            socket.close();
        }
    },

    exitGameAndLogout(context) {
        const gameID = context.getters.gameID;
        const socket = context.getters.socket;
        
        // If in a game, send exit event
        if (gameID !== 0 && socket && socket.readyState === WebSocket.OPEN) {
            const event = new Event('game_exit', gameID, {});
            try {
                socket.send(JSON.stringify(event));
            } catch (error) {
                console.error('Error sending exit event:', error);
            }
        }
        
        // Close WebSocket connection
        if (socket && socket.readyState !== WebSocket.CLOSED) {
            socket.close();
        }
        
        // Reset all WebSocket state
        context.commit('resetWebSocket');
    },

    sendEvent(context, eventdata){
        const event = new Event(eventdata.type, eventdata.id, eventdata.payload);
        console.log(event);
        const wsocket = context.getters.socket;
        
        if (!wsocket || wsocket.readyState !== WebSocket.OPEN) {
            console.error('WebSocket is not connected. Cannot send event.');
            return;
        }
        
        try {
            wsocket.send(JSON.stringify(event));
        } catch (error) {
            console.error('Error sending WebSocket message:', error);
        }
    },

    gameInit(context, data){
        if (data.gameId == 0){
            console.log("No gameid was returned");
            const socket = context.getters.socket;
            if (socket) {
                socket.close();
            }
            return
        }
        context.commit("setGameId", data.gameId);
        console.log("GAME: ", data.gameId);

        context.commit("setCards", data.payload.cards);

        context.commit("setIndex", data.payload.index);
        
        if (data.payload.playerNames) {
            context.commit("setPlayerNames", data.payload.playerNames);
        }        
    },

    sendBidAction(context, { action, color }) {
        const gameID = context.getters.gameID;
        const event = new Event('bid_action', gameID, { action, color });
        
        const wsocket = context.getters.socket;
        if (!wsocket || wsocket.readyState !== WebSocket.OPEN) {
            console.error('WebSocket is not connected. Cannot send bid action.');
            return;
        }
        
        try {
            wsocket.send(JSON.stringify(event));
        } catch (error) {
            console.error('Error sending bid action:', error);
        }
    },

    sendTalonExchange(context, discardCards) {
        const gameID = context.getters.gameID;
        const event = new Event('talon_exchange', gameID, { discardcards: discardCards });
        
        const wsocket = context.getters.socket;
        if (!wsocket || wsocket.readyState !== WebSocket.OPEN) {
            console.error('WebSocket is not connected. Cannot send talon exchange.');
            return;
        }
        
        // Remove discarded cards from local state immediately
        discardCards.forEach(cardId => {
            context.commit('deleteCard', cardId);
        });
        
        try {
            wsocket.send(JSON.stringify(event));
        } catch (error) {
            console.error('Error sending talon exchange:', error);
        }
    }

}