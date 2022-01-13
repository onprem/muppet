package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/onprem/muppet/pkg/api"
)

func TestServer(t *testing.T) {
	handler := api.Handler(NewServer())

	var data api.Commands

	t.Run("add command", func(t *testing.T) {
		req, err := api.NewAddCommandRequest(
			"localhost",
			api.AddCommandJSONRequestBody{Host: "host001", ShellCommand: "apt update"},
		)
		testutil.Ok(t, err)

		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		testutil.Equals(t, http.StatusOK, recorder.Code)

		req, err = api.NewListCommandQueueRequest("localhost", "host001")
		testutil.Ok(t, err)

		recorder = httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		testutil.Equals(t, http.StatusOK, recorder.Code)

		resp, err := api.ParseListCommandQueueResponse(recorder.Result())
		testutil.Ok(t, err)

		data = *resp.JSON200
		testutil.Equals(t, 1, len(data))
		testutil.Equals(t, "apt update", data[0].ShellCommand)
	})

	t.Run("mark a command done", func(t *testing.T) {
		req, err := api.NewMarkCommandDoneRequest(
			"localhost",
			"host001", api.MarkCommandDoneJSONRequestBody{Uuid: data[0].Uuid, ExitStatus: 1},
		)
		testutil.Ok(t, err)

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		testutil.Equals(t, http.StatusOK, recorder.Code)

		req, err = api.NewListCommandQueueRequest("localhost", "host001")
		testutil.Ok(t, err)

		recorder = httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		testutil.Equals(t, http.StatusOK, recorder.Code)

		resp, err := api.ParseListCommandQueueResponse(recorder.Result())
		testutil.Ok(t, err)

		data = *resp.JSON200
		testutil.Equals(t, 0, len(data))
	})
}
