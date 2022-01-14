package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/onprem/muppet/pkg/api"
	"github.com/onprem/muppet/pkg/store"
)

type Server struct {
	store store.Store
}

var _ api.ServerInterface = &Server{}

func NewServer(s store.Store) *Server {
	return &Server{
		store: s,
	}
}

func (s *Server) AddCommand(w http.ResponseWriter, r *http.Request) {
	var data api.AddCommandJSONBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	err := s.store.AddCommand(data.Host, api.Command{ShellCommand: data.ShellCommand, Uuid: uuid.NewString()})
	if err != nil {
		http.Error(w, "error adding command to queue", http.StatusInternalServerError)

		return
	}

	w.Write([]byte("command added to queue"))
}

func (s *Server) ListCommandQueue(w http.ResponseWriter, r *http.Request, host string) {
	cmds, err := s.store.GetPendingCommands(host)
	if err != nil {
		http.Error(w, "error fetching commands from queue", http.StatusInternalServerError)

		return
	}

	w.Header().Set("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(cmds); err != nil {
		http.Error(w, "error encoding data", http.StatusInternalServerError)

		return
	}
}

func (s *Server) MarkCommandDone(w http.ResponseWriter, r *http.Request, host string) {
	var data api.MarkCommandDoneJSONBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	err := s.store.MarkDone(host, data.Uuid, uint(data.ExitStatus))
	if err != nil {
		http.Error(w, "failed marking command as done", http.StatusInternalServerError)

		return
	}

	w.Write([]byte("command marked as done"))
}
