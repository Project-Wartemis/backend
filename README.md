# project-artemis

## Summary

    Basically Warlight with an API allow Bots to play each other:
        -> Player starts bot
        -> Bot registers with the backend (with websocket)
        -> Player starts game in Frontend/http request
        -> Bot receives message from backend that websocket is available for game 
        -> Bot sends/recieves moves and gamestates
        -> ...
        -> Game over
        -> Bot is nottified and disconnects


## 4 Components

### Game Engine
Language: GO

    Game features:
        Possible player moves:
            - Deploy troops
            - Move troops

    Core functionality:
        - Start new game:
            * Load initial game state from file
        - Validate moves
        - Execute player moves
        - Return game results

### Backend
Language: GO

    Core Functionality:
        - Communicate with Game engine
        - Websockets for bots:
            * parsing message
            * validation
        - Websocket for frontend:
            * Receive message from frontend:
                * next turn              
        - Minimal REST API:
            * starting new game

### Frontend
Language: Angular

    Core Functionality:
        - View all connected bots
        - View all ongoing games
        - Start new game and select bots
        - Send 'next turn' message
        - Show currect game state (via graph library)

### Demo bot
Language: TBD

  Core Functionality:
    Very simple bot that:
      * makes valid moves
      * doesn't crash
      * disconnects at end of game
    
