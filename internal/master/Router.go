package master

import (
	"flag"
	"fmt"
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/Project-Wartemis/pw-backend/internal/wrapper"
)

type Router struct {
	router *mux.Router
}

func NewRouter() *Router {
	return &Router {
		router: mux.NewRouter(),
	}
}

func (this *Router) Start(port int) {
	endpoint := fmt.Sprintf("0.0.0.0:%d", port)
	address := flag.String("addr", endpoint, "http service address")
	log.Infof("Starting http listener on port %d", port)
	err := http.ListenAndServe(*address, this.router)
	if err != nil {
		log.Error("Could not start http listener")
		log.Panic(err)
	}
}

func (this *Router) Initialise(lobbyWrapper *wrapper.LobbyWrapper, roomWrapper *wrapper.RoomWrapper) {
	this.router.HandleFunc("/lobby",         lobbyWrapper.GetLobby).Methods("GET")
	this.router.HandleFunc("/socket/{room}", roomWrapper.NewConnection)
	this.router.HandleFunc("/socket",        lobbyWrapper.NewConnection)
	this.router.HandleFunc("/*",             NotFoundHandler)
}

func NotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
}
