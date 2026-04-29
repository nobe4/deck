//go:build linux

package media

import (
	"fmt"
	"strconv"
	"strings"
)

// pactl implements volumeCtrl via PulseAudio.
// pactl get-sink-volume @DEFAULT_SINK@ -> "Volume: front-left: 32768 /  50% / ..."
// pactl get-sink-mute @DEFAULT_SINK@ -> "Mute: yes" or "Mute: no"
type pactl struct{}

func (p *pactl) getVolume() (int, error) {
	out, err := run("pactl", "get-sink-volume", "@DEFAULT_SINK@")
	if err != nil {
		return 0, err
	}
	_, after, ok := strings.Cut(out, "/ ")
	if !ok {
		return 0, fmt.Errorf("unexpected pactl output: %q", out)
	}
	pct, _, ok := strings.Cut(after, "%")
	if !ok {
		return 0, fmt.Errorf("unexpected pactl output: %q", out)
	}
	return strconv.Atoi(strings.TrimSpace(pct))
}

func (p *pactl) setVolume(level int) error {
	_, err := run("pactl", "set-sink-volume", "@DEFAULT_SINK@", fmt.Sprintf("%d%%", level))
	return err
}

func (p *pactl) mute() error {
	_, err := run("pactl", "set-sink-mute", "@DEFAULT_SINK@", "toggle")
	return err
}

func (p *pactl) isMuted() (bool, error) {
	out, err := run("pactl", "get-sink-mute", "@DEFAULT_SINK@")
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "yes"), nil
}
