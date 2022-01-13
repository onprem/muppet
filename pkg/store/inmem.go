package store

import (
	"errors"

	"github.com/onprem/muppet/pkg/api"
)

type commands struct {
	pending map[string]Command
	done    map[string]Command
}

type InMemStore struct {
	data map[string]commands
}

var _ Store = &InMemStore{}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		data: make(map[string]commands),
	}
}

func (im *InMemStore) AddCommand(host string, cmd api.Command) error {
	if host == "" {
		return errors.New("empty host is not allowed")
	}

	if cmd.Uuid == "" {
		return errors.New("empty UUID is not allowed")
	}

	if _, ok := im.data[host]; !ok {
		im.data[host] = commands{
			pending: make(map[string]Command),
			done:    make(map[string]Command),
		}
	}

	im.data[host].pending[cmd.Uuid] = Command{
		Cmd: cmd,
	}

	return nil
}

func (im *InMemStore) GetPendingCommands(host string) (api.Commands, error) {
	cmds := make([]api.Command, 0, len(im.data[host].pending))

	for _, v := range im.data[host].pending {
		cmds = append(cmds, v.Cmd)
	}

	return cmds, nil
}

func (im *InMemStore) GetDoneCommands(host string) ([]Command, error) {
	cmds := make([]Command, 0, len(im.data[host].done))

	for _, v := range im.data[host].done {
		cmds = append(cmds, v)
	}

	return cmds, nil
}

func (im *InMemStore) MarkDone(host, uuid string, exitStatus uint) error {
	cmds, ok := im.data[host]
	if !ok {
		return errors.New("invalid host")
	}

	cmd, ok := cmds.pending[uuid]
	if !ok {
		return errors.New("invalid uuid or already marked done")
	}

	cmd.ExitStatus = exitStatus
	im.data[host].done[uuid] = cmd

	delete(im.data[host].pending, uuid)

	return nil
}
