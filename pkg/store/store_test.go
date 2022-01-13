package store

import (
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"

	"github.com/onprem/muppet/pkg/api"
)

func TestInMemStore(t *testing.T) {
	t.Parallel()

	im := NewInMemStore()
	testCommandStore(t, im)
}

func testCommandStore(t *testing.T, s Store) {
	var host1, host2 string = "host1", "host2"

	cmd1 := api.Command{ShellCommand: "apt update", Uuid: "abc"}
	cmd2 := api.Command{ShellCommand: "apt upgrade", Uuid: "def"}
	cmd3 := api.Command{ShellCommand: "pacman -Syu", Uuid: "ghi"}

	t.Run("empty host", func(t *testing.T) {
		err := s.AddCommand("", api.Command{ShellCommand: "apt update", Uuid: "abc"})
		testutil.NotOk(t, err)
	})

	t.Run("empty uuid", func(t *testing.T) {
		err := s.AddCommand(host1, api.Command{ShellCommand: "apt update", Uuid: ""})
		testutil.NotOk(t, err)
	})

	t.Run("add a command", func(t *testing.T) {
		err := s.AddCommand(host1, cmd1)
		testutil.Ok(t, err)

		cmds, err := s.GetPendingCommands(host1)
		testutil.Ok(t, err)

		testutil.Equals(t, 1, len(cmds))
		testutil.Equals(t, cmd1, cmds[0])
	})

	t.Run("add commands to multiple hosts", func(t *testing.T) {
		err := s.AddCommand(host1, cmd2)
		testutil.Ok(t, err)

		err = s.AddCommand(host2, cmd3)
		testutil.Ok(t, err)

		cmds, err := s.GetPendingCommands(host1)
		testutil.Ok(t, err)

		testutil.Equals(t, 2, len(cmds))
		testutil.Equals(t, api.Commands{cmd1, cmd2}, cmds)

		cmds, err = s.GetPendingCommands(host2)
		testutil.Ok(t, err)

		testutil.Equals(t, 1, len(cmds))
		testutil.Equals(t, api.Commands{cmd3}, cmds)
	})

	t.Run("update status", func(t *testing.T) {
		err := s.MarkDone(host1, cmd2.Uuid, 0)
		testutil.Ok(t, err)

		cmds, err := s.GetDoneCommands(host1)
		testutil.Ok(t, err)

		testutil.Equals(t, []Command{{Cmd: cmd2}}, cmds)
	})
}
