//go:build linux

package media

import (
	"fmt"
	"os/exec"
	"strings"
)

func New() (Controller, error) {
	c := &controller{}

	switch {
	case hasCmd("wpctl"):
		c.vol = &wpctl{}
	case hasCmd("pactl"):
		c.vol = &pactl{}
	case hasCmd("amixer"):
		c.vol = &amixer{}
	}

	if hasCmd("playerctl") {
		c.media = &playerctl{}
	}

	if c.vol == nil && c.media == nil {
		return nil, fmt.Errorf("no backends found (install wpctl/pactl/amixer and playerctl)")
	}

	return c, nil
}

func hasCmd(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func run(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}
