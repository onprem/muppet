package store

import "github.com/onprem/muppet/pkg/api"

type Command struct {
	Cmd        api.Command
	ExitStatus uint
}

type Store interface {
	AddCommand(host string, command api.Command) error
	GetPendingCommands(host string) (api.Commands, error)
	GetDoneCommands(host string) ([]Command, error)
	MarkDone(host, uuid string, exitStatus uint) error
}
