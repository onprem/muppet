package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"

	"github.com/onprem/muppet/pkg/api"
	"github.com/onprem/muppet/pkg/store"
)

type Server struct {
	store  store.Store
	logger log.Logger
}

var _ api.ServerInterface = &Server{}

func NewServer(s store.Store, logger log.Logger) *Server {
	return &Server{
		store:  s,
		logger: logger,
	}
}

func (s *Server) AddCommand(w http.ResponseWriter, r *http.Request) {
	logger := log.With(s.logger, "handler", "add command")

	var data api.AddCommandJSONBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		level.Debug(logger).Log("msg", "parsing request body", "err", err)

		return
	}
	defer r.Body.Close()

	cmd := api.Command{ShellCommand: data.ShellCommand, Uuid: uuid.NewString()}
	err := s.store.AddCommand(data.Host, cmd)
	if err != nil {
		http.Error(w, "error adding command to queue", http.StatusInternalServerError)
		level.Error(logger).Log("msg", "adding command to queue", "err", err)

		return
	}

	level.Info(s.logger).Log("msg", "added command to queue", "host", data.Host, "cmd", cmd.ShellCommand, "uuid", cmd.Uuid)

	w.Write([]byte("command added to queue"))
}

func (s *Server) ListCommandQueue(w http.ResponseWriter, r *http.Request, host string) {
	logger := log.With(s.logger, "handler", "list command queue")

	cmds, err := s.store.GetPendingCommands(host)
	if err != nil {
		http.Error(w, "error fetching commands from queue", http.StatusInternalServerError)
		level.Error(logger).Log("msg", "fetching commands from queue", "err", err)

		return
	}

	w.Header().Set("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(cmds); err != nil {
		http.Error(w, "error encoding data", http.StatusInternalServerError)
		level.Error(logger).Log("msg", "encoding commands", "err", err)

		return
	}
}

func (s *Server) MarkCommandDone(w http.ResponseWriter, r *http.Request, host string) {
	logger := log.With(s.logger, "handler", "mark command done")

	var data api.MarkCommandDoneJSONBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		level.Debug(logger).Log("msg", "parsing request body", "err", err)

		return
	}
	defer r.Body.Close()

	err := s.store.MarkDone(host, data.Uuid, uint(data.ExitStatus))
	if err != nil {
		http.Error(w, "failed marking command as done", http.StatusInternalServerError)
		level.Error(logger).Log("msg", "marking command as done", "uuid", data.Uuid, "err", err)

		return
	}

	level.Info(s.logger).Log("msg", "marked command as done", "host", host, "uuid", data.Uuid, "exitcode", data.ExitStatus)

	w.Write([]byte("command marked as done"))
}
