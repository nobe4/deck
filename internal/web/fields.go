package web

import (
	"encoding/json"
	"fmt"
)

func (s *Server) get(field string) (any, error) {
	switch field {
	case "volume":
		return s.ctrl.Volume()
	case "mute":
		return s.ctrl.IsMuted()
	default:
		return nil, fmt.Errorf("unknown field: %s", field)
	}
}

func (s *Server) set(field string, body []byte) error {
	switch field {
	case "volume":
		var req struct {
			Volume int `json:"volume"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			return err
		}
		return s.ctrl.SetVolume(req.Volume)

	case "playpause":
		return s.ctrl.PlayPause()

	case "next":
		return s.ctrl.Next()

	case "previous":
		return s.ctrl.Previous()

	case "mute":
		return s.ctrl.Mute()

	default:
		return fmt.Errorf("unknown field: %s", field)
	}
}
