//go:build darwin

package media

/*
#cgo LDFLAGS: -framework Cocoa -framework CoreAudio -framework AudioToolbox
#include "bridge_darwin.h"
*/
import "C"

import "fmt"

type darwin struct{}

func New() (Controller, error) {
	d := &darwin{}
	return &controller{vol: d, media: d}, nil
}

func (d *darwin) getVolume() (int, error) {
	vol := int(C.GetVolume())
	if vol < 0 {
		return 0, fmt.Errorf("get volume failed")
	}
	return vol, nil
}

func (d *darwin) setVolume(level int) error {
	if C.SetVolume(C.int(level)) != 0 {
		return fmt.Errorf("set volume failed")
	}
	return nil
}

func (d *darwin) playPause() error { C.MediaKey(16); return nil }
func (d *darwin) next() error      { C.MediaKey(17); return nil }
func (d *darwin) previous() error  { C.MediaKey(18); return nil }

func (d *darwin) mute() error {
	muted := int(C.GetMute())
	if muted < 0 {
		return fmt.Errorf("get mute failed")
	}
	if C.SetMute(C.int(1-muted)) != 0 {
		return fmt.Errorf("set mute failed")
	}
	return nil
}

func (d *darwin) isMuted() (bool, error) {
	muted := int(C.GetMute())
	if muted < 0 {
		return false, fmt.Errorf("get mute failed")
	}
	return muted == 1, nil
}
