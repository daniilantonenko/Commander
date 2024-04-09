package server

import (
	"commander/config"
	"log"
	"net/http"
	"os/exec"
	"sync"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Server struct
type Server struct {
	cfg    *config.Config
	db     *gorm.DB
	mu     sync.Mutex
	cmdMap map[int]*exec.Cmd
}

// NewServer New Server constructor
func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{cfg: cfg, db: db, cmdMap: make(map[int]*exec.Cmd)}
}

func (s *Server) Run() error {

	router := mux.NewRouter()
	router.HandleFunc("/commands", s.createCommand).Methods("POST")
	router.HandleFunc("/commands", s.getCommands).Methods("GET")
	router.HandleFunc("/commands/{id}", s.getCommand).Methods("GET")
	router.HandleFunc("/commands/{id}", s.stopCommand).Methods("DELETE")

	log.Println("App started!")

	log.Fatal(http.ListenAndServe(":"+s.cfg.Port, router))

	return nil
}
