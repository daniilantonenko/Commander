package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Command
type Command struct {
	ID      int    `json:"id" gorm:"unique"`
	Command string `json:"command"`
	Output  string `json:"output,omitempty"`
	Status  string `json:"status,omitempty"`
}

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
func createCommand(w http.ResponseWriter, r *http.Request) {
	var cmd Command
	_ = json.NewDecoder(r.Body).Decode(&cmd)

	if cmd.Command == "" || cmd.Output != "" || cmd.Status != "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Incorrect form data", cmd)
		w.Write([]byte("incorrect form data"))
		return
	}

	err := db.Create(&cmd).Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		w.Write([]byte("error creating command"))
		return
	}

	json.NewEncoder(w).Encode(cmd)

	// Run the command in goroutine
	go func(command string, id int) {
		mu.Lock()
		cmdMap[id] = exec.Command("bash", "-c", command)
		cmd := cmdMap[id]
		mu.Unlock()

		output, _ := cmd.CombinedOutput()

		mu.Lock()
		delete(cmdMap, id)
		mu.Unlock()

		db.Model(&Command{}).Where("id = ?", id).Update("output", string(output))
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
func getCommands(w http.ResponseWriter, r *http.Request) {
	var commands []Command
	db.Find(&commands)
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
func getCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var cmd Command
	err := db.First(&cmd, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Record id not found:", id)
		w.Write([]byte("find command error"))
		return
	}

	mu.Lock()
	if execCmd, ok := cmdMap[cmd.ID]; ok {
		cmd.Status = "Running"
		cmdMap[cmd.ID] = execCmd
	}
	mu.Unlock()

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
func stopCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id int
	var err error
	id, err = strconv.Atoi(params["id"])

	if err != nil {
		log.Println("Error converting id:", err)
		return
	}

	mu.Lock()
	if execCmd, ok := cmdMap[id]; ok {
		err := execCmd.Process.Kill()
		if err != nil {
			log.Println("Error stopping command:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error stopping command"))
			mu.Unlock()
			return
		}
		delete(cmdMap, id)
	}
	mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
