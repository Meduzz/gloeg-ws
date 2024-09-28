package api

import "encoding/json"

type (
	WsEvent struct {
		Kind string          `json:"kind"`
		Flag *FlagEvent      `json:"flag,omitempty"`
		Log  json.RawMessage `json:"log,omitempty"`
	}

	FlagEvent struct {
		Kind  string `json:"kind"`
		Name  string `json:"name"`
		Value any    `json:"value"`
	}
)

const (
	EventKindLog         = "logs"
	EventKindFlagRemote  = "flag.remote"
	EventKindFlagLocal   = "flag.local"
	EventKindFlagRequest = "flag.request"
)
