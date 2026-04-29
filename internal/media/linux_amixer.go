//go:build linux

package media

import (
	"fmt"
	"strings"
)

// amixer implements volumeCtrl via ALSA.
// amixer sget Master -> "... [50%] [on]" or "... [50%] [off]"
type amixer struct{}

func (a *amixer) getVolume() (int, error) {
	out, err := run("amixer", "sget", "Master")
	if err != nil {
		return 0, err
	}
	return parsePercentage(out, "[")
}

func (a *amixer) setVolume(level int) error {
	_, err := run("amixer", "sset", "Master", fmt.Sprintf("%d%%", level))
	return err
}

func (a *amixer) mute() error {
	_, err := run("amixer", "sset", "Master", "toggle")
	return err
}

func (a *amixer) isMuted() (bool, error) {
	out, err := run("amixer", "sget", "Master")
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "[off]"), nil
}
