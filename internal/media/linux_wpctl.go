//go:build linux

package media

import (
	"fmt"
	"strconv"
	"strings"
)

// wpctl implements volumeCtrl via WirePlumber/PipeWire.
// wpctl get-volume @DEFAULT_AUDIO_SINK@ -> "Volume: 0.50" or "Volume: 0.50 [MUTED]"
type wpctl struct{}

func (w *wpctl) getVolume() (int, error) {
	out, err := run("wpctl", "get-volume", "@DEFAULT_AUDIO_SINK@")
	if err != nil {
		return 0, err
	}
	_, val, ok := strings.Cut(out, " ")
	if !ok {
		return 0, fmt.Errorf("unexpected wpctl output: %q", out)
	}
	// val may contain "[MUTED]" suffix, take first field
	val, _, _ = strings.Cut(val, " ")
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, fmt.Errorf("parse wpctl volume: %w", err)
	}
	return int(f * 100), nil
}

func (w *wpctl) setVolume(level int) error {
	_, err := run("wpctl", "set-volume", "@DEFAULT_AUDIO_SINK@", fmt.Sprintf("%d%%", level))
	return err
}

func (w *wpctl) mute() error {
	_, err := run("wpctl", "set-mute", "@DEFAULT_AUDIO_SINK@", "toggle")
	return err
}

func (w *wpctl) isMuted() (bool, error) {
	out, err := run("wpctl", "get-volume", "@DEFAULT_AUDIO_SINK@")
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "[MUTED]"), nil
}
