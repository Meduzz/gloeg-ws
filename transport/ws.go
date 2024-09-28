package transport

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/Meduzz/gloegg-ws/api"
	"golang.org/x/net/websocket"
)

type (
	Transport struct {
		ws           *websocket.Conn
		config       *websocket.Config
		reconnecting bool
		handler      func(*api.FlagEvent)
	}
)

const (
	disconnected = "Gloegg WS disconnected"
	reconnecting = "Gloegg WS reconnecting"
	failed       = "Gloegg WS reconnecting failed"
	success      = "Gloegg WS reconnecting success"
)

func NewTransport(url, token string) (*Transport, error) {
	config, err := websocket.NewConfig(url, url)

	if err != nil {
		return nil, err
	}

	if token != "" {
		config.Header.Set("Authorization", token)
	}

	t := &Transport{
		config:       config,
		reconnecting: false,
	}

	err = t.connect()

	if err != nil {
		go t.Reconnect()
		return t, err
	}

	err = t.RequestFlags()

	if err != nil {
		return t, err
	}

	return t, nil
}

func (t *Transport) SendFlag(evt *api.FlagEvent) error {
	wsEvent := &api.WsEvent{}
	wsEvent.Kind = api.EventKindFlagLocal
	wsEvent.Flag = evt

	err := websocket.JSON.Send(t.ws, wsEvent)

	if err != nil && errors.Is(err, io.EOF) {
		// drop this flag and start reconnecting
		go t.Reconnect()
	}

	return err
}

func (t *Transport) OnFlagUpdated(handler func(*api.FlagEvent)) {
	t.handler = handler

	if t.ws == nil {
		return
	}

	for {
		wsEvent := &api.WsEvent{}
		err := websocket.JSON.Receive(t.ws, wsEvent)

		if err != nil {
			if errors.Is(err, io.EOF) {
				// socket is closed
				fakelog(disconnected)

				t.Reconnect()

				break
			}

			// we're kind of out of options at this point
			println(time.Now().Format("2006-01-02 15:04:05"), "serializing event threw error", err.Error())
		}

		if wsEvent.Kind == api.EventKindFlagRemote {
			handler(wsEvent.Flag)
		}
	}
}

func (t *Transport) RequestFlags() error {
	wsEvent := &api.WsEvent{}
	wsEvent.Kind = api.EventKindFlagRequest

	err := websocket.JSON.Send(t.ws, wsEvent)

	if err != nil && errors.Is(err, io.EOF) {
		go t.Reconnect()
	}

	return err
}

func (t *Transport) SendLog(log json.RawMessage) error {
	wsEvent := &api.WsEvent{}
	wsEvent.Kind = api.EventKindLog
	wsEvent.Log = log

	if t.ws == nil {
		fakelog(string(wsEvent.Log))
		go t.Reconnect()

		return nil
	}

	err := websocket.JSON.Send(t.ws, wsEvent)

	if err != nil && errors.Is(err, io.EOF) {
		// dump the log in stdout
		fakelog(string(wsEvent.Log))
		go t.Reconnect()
	}

	return err
}

func (t *Transport) connect() error {
	ws, err := websocket.DialConfig(t.config)

	if err != nil {
		return err
	}

	t.ws = ws

	return nil
}

func (t *Transport) Reconnect() error {
	if t.reconnecting {
		return nil
	}

	t.reconnecting = true

	fakelog(reconnecting)

	defer func() { t.reconnecting = false }()

	err := t.connect()

	if err != nil {
		fakelog(failed)
		ticker := time.NewTicker(5 * time.Second)

		for range ticker.C {
			fakelog(reconnecting)
			err = t.connect()

			if err == nil {
				ticker.Stop()
				break
			} else {
				fakelog(failed)
			}
		}
	}

	if t.handler != nil {
		go t.OnFlagUpdated(t.handler)
	}

	fakelog(success)

	// re-request flags just in case
	t.RequestFlags()

	return nil
}

func fakelog(msg string) {
	println(time.Now().Format("2006-01-02 15:04:05"), msg)
}
