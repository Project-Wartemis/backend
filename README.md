# Summary

Language: GO

* Core functionality
  * Provide communication to Engine for other components
  * Websockets for bots
    * parsing message
    * validation
    * relay to engine if needed
  * Websocket for frontend
    * Receive gamestate of each turn
    * Send signal to advance one turn (if the game was configured no not autoplay)
  * HTTP REST API for frontend
    * starting new game

* Links
  1. Backend <=> Engine : HTTP (2 way communication)
  2. Backend  <= Frontend : HTTP
  3. Backend <=> Frontend : Websockets
  4. Backend <=> Bot : Websockets
