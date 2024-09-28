package flags

import (
	"github.com/Meduzz/gloegg-ws/api"
	"github.com/Meduzz/gloegg-ws/transport"
	"github.com/Meduzz/gloegg/toggles"
)

type (
	toggleConfig struct {
		ws *transport.Transport
	}
)

func SetupFlagHandler(ws *transport.Transport) {
	config := &toggleConfig{
		ws: ws,
	}

	toggles.Subscribe(config.localFlagHander)
	go ws.OnFlagUpdated(config.onRemoteFlag)
}

func (c *toggleConfig) localFlagHander(toggle *toggles.UpdatedToggle) {
	event := &api.FlagEvent{}
	event.Kind = toggle.Kind
	event.Name = toggle.Name
	event.Value = toggle.Value

	err := c.ws.SendFlag(event)

	if err != nil {
		// TODO error handling
	}
}

func (c *toggleConfig) onRemoteFlag(event *api.FlagEvent) {
	switch event.Kind {
	case toggles.KindBool:
		value, ok := event.Value.(bool)

		if ok {
			toggles.SetBoolToggle(event.Name, value)
		}
	case toggles.KindFloat32:
		value, ok := event.Value.(float32)

		if ok {
			toggles.SetFloat32Toggle(event.Name, value)
		}
	case toggles.KindFloat64:
		value, ok := event.Value.(float64)

		if ok {
			toggles.SetFloat64Toggle(event.Name, value)
		}
	case toggles.KindInt:
		value, ok := event.Value.(int)

		if ok {
			toggles.SetIntToggle(event.Name, value)
		}
	case toggles.KindInt64:
		value, ok := event.Value.(int64)

		if ok {
			toggles.SetInt64Toggle(event.Name, value)
		}
	case toggles.KindObject:
		value, ok := event.Value.(map[string]any)

		if ok {
			toggles.SetObjectToggle(event.Name, value)
		}
	case toggles.KindString:
		value, ok := event.Value.(string)

		if ok {
			toggles.SetStringToggle(event.Name, value)
		}
	default:
	}
}
