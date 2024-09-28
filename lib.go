package gloegws

import (
	"fmt"

	"github.com/Meduzz/gloegg"
	"github.com/Meduzz/gloegg-ws/flags"
	"github.com/Meduzz/gloegg-ws/logs"
	"github.com/Meduzz/gloegg-ws/transport"
	"github.com/Meduzz/helper/utilz"
)

func Setup(token, project, service string) error {
	if project != "" {
		gloegg.AddMeta("@project", project)
	} else {
		return fmt.Errorf("no project specified")
	}

	if service != "" {
		gloegg.AddMeta("@service", service)
	} else {
		return fmt.Errorf("no service specified")
	}

	url := utilz.Env("GLOEGG_URL", "http://localhost:8080")

	ws, err := transport.NewTransport(url, token)

	if ws == nil {
		// TODO this one is problematic (but should be very rare)
		return err
	}

	flags.SetupFlagHandler(ws)
	factory := logs.NewLogFactory(ws)
	gloegg.SetHandlerFactory(factory)

	return err
}
