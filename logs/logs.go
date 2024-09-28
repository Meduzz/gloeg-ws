package logs

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"

	"github.com/Meduzz/gloegg-ws/transport"
	"github.com/Meduzz/gloegg/common"
	"github.com/Meduzz/gloegg/logging"
)

type (
	wsLogFactory struct {
		ws *transport.Transport
	}
)

func NewLogFactory(ws *transport.Transport) logging.HandlerFactory {
	return &wsLogFactory{ws}
}

func (f *wsLogFactory) Spawn(level slog.Leveler, tags common.Tags) slog.Handler {
	r, w := io.Pipe()

	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	})

	go f.listen(r)

	return handler
}

func (f *wsLogFactory) listen(reader io.Reader) {
	decoder := json.NewDecoder(reader)

	for decoder.More() {
		record := json.RawMessage{}
		err := decoder.Decode(&record)

		if err == nil {
			err = f.ws.SendLog(record)

			if err != nil && !errors.Is(err, io.EOF) {
				println("sending log threw error", err.Error())
			}
		}

		// TODO safe to ignore this error?
	}
}
