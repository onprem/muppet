package store

import (
	"errors"

	"github.com/onprem/muppet/pkg/api"
)

type command struct {
	cmd    api.Command
	status Status
}

type commands map[string]command

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
		im.data[host] = make(commands)
	}

	im.data[host][cmd.Uuid] = command{
		cmd:    cmd,
		status: StatusPending,
	}

	return nil
}

func (im *InMemStore) GetCommandsWithStatus(host string, status Status) (api.Commands, error) {
	var cmds api.Commands

	for _, v := range im.data[host] {
		if v.status != status {
			continue
		}

		cmds = append(cmds, v.cmd)
	}

	return cmds, nil
}

func (im *InMemStore) UpdateStatus(host, uuid string, status Status) error {
	cmds, ok := im.data[host]
	if !ok {
		return errors.New("invalid host")
	}

	cmd, ok := cmds[uuid]
	if !ok {
		return errors.New("invalid uuid")
	}

	cmd.status = status
	im.data[host][uuid] = cmd

	return nil
}
