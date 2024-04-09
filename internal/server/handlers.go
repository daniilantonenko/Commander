package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"commander/internal/models"
)

// createCommand godoc
// @Summary      Create command
// @Description  create accounts
// @Tags         commands
// @Accept       json
// @Produce      json
// @Param        data    body	Command true "form data"
// @Success      200  {object}   Command
// @Failure      400  {string}  string    "bad request"
// @Failure      404  {string}  string    "bad request"
// @Router       /commands [post]
func (s *Server) createCommand(w http.ResponseWriter, r *http.Request) {
	var cmd models.Command
	_ = json.NewDecoder(r.Body).Decode(&cmd)

	if cmd.Command == "" || cmd.Output != "" || cmd.Status != "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Incorrect form data", cmd)
		w.Write([]byte("incorrect form data"))
		return
	}

	err := s.db.Create(&cmd).Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		w.Write([]byte("error creating command"))
		return
	}

	json.NewEncoder(w).Encode(cmd)

	// Run the command in goroutine
	go func(command string, id int) {
		s.mu.Lock()
		s.cmdMap[id] = exec.Command("bash", "-c", command)
		cmd := s.cmdMap[id]
		s.mu.Unlock()

		output, _ := cmd.CombinedOutput()

		s.mu.Lock()
		delete(s.cmdMap, id)
		s.mu.Unlock()

		s.db.Model(&models.Command{}).Where("id = ?", id).Update("output", string(output))
	}(cmd.Command, cmd.ID)
}

// getCommands godoc
// @Summary      List commands
// @Description  get accounts
// @Tags         commands
// @Accept       json
// @Produce      json
// @Success      200  {array}   Command
// @Router       /commands [get]
func (s *Server) getCommands(w http.ResponseWriter, r *http.Request) {
	var commands []models.Command
	s.db.Find(&commands)
	json.NewEncoder(w).Encode(commands)
}

// getCommand godoc
// @Summary      Show command
// @Description  get accounts
// @Tags         commands
// @Accept       json
// @Produce      json
// @Param        some_id    path	int  false  "id search by some_id"
// @Success      200  {object}   Command
// @Failure      404  {string}  string    "record not found"
// @Router       /commands/{some_id} [get]
func (s *Server) getCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var cmd models.Command
	err := s.db.First(&cmd, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Record id not found:", id)
		w.Write([]byte("find command error"))
		return
	}

	s.mu.Lock()
	if execCmd, ok := s.cmdMap[cmd.ID]; ok {
		cmd.Status = "Running"
		s.cmdMap[cmd.ID] = execCmd
	}
	s.mu.Unlock()

	json.NewEncoder(w).Encode(cmd)
}

// stopCommand godoc
// @Summary      Stop command
// @Description  stop command
// @Tags         commands
// @Accept       json
// @Produce      json
// @Param        some_id    path	int  false  "id search by some_id"
// @Success      200  {object}   Command
// @Failure      404  {string}  string    "record not found"
// @Router       /commands/{some_id} [delete]
func (s *Server) stopCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id int
	var err error
	id, err = strconv.Atoi(params["id"])

	if err != nil {
		log.Println("Error converting id:", err)
		return
	}

	s.mu.Lock()
	if execCmd, ok := s.cmdMap[id]; ok {
		err := execCmd.Process.Kill()
		if err != nil {
			log.Println("Error stopping command:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error stopping command"))
			s.mu.Unlock()
			return
		}
		delete(s.cmdMap, id)
	}
	s.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
