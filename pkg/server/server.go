package server

import (
	"net/http"

	"github.com/onprem/muppet/pkg/api"
)

type Server struct{}

var _ api.ServerInterface = &Server{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) AddCommand(w http.ResponseWriter, r *http.Request) {}

func (s *Server) ListCommandQueue(w http.ResponseWriter, r *http.Request, host string) {}

func (s *Server) MarkCommandDone(w http.ResponseWriter, r *http.Request, host string) {}
