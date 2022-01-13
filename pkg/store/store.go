package store

import "github.com/onprem/muppet/pkg/api"

type Status string

const (
	StatusPending  Status = "pending"
	StatusFinished Status = "finished"
	StatusErrored  Status = "errored"
)

type Store interface {
	AddCommand(host string, command api.Command) error
	GetCommandsWithStatus(host string, status Status) (api.Commands, error)
	UpdateStatus(host, uuid string, status Status) error
}

func GetPendingCommands(s Store, host string) (api.Commands, error) {
	return s.GetCommandsWithStatus(host, StatusPending)
}
