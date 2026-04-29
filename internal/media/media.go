// Package media provides an interface for controlling system media playback and volume.
package media

// Controller controls system media playback and volume.
type Controller interface {
	PlayPause() error
	Next() error
	Previous() error
	Mute() error
	IsMuted() (bool, error)
	SetVolume(level int) error
	Volume() (int, error)
}

type volumeCtrl interface {
	getVolume() (int, error)
	setVolume(int) error
	mute() error
	isMuted() (bool, error)
}

type mediaCtrl interface {
	playPause() error
	next() error
	previous() error
}

type controller struct {
	vol   volumeCtrl
	media mediaCtrl
}

func (c *controller) Volume() (int, error)   { return c.vol.getVolume() }
func (c *controller) SetVolume(v int) error  { return c.vol.setVolume(v) }
func (c *controller) Mute() error            { return c.vol.mute() }
func (c *controller) IsMuted() (bool, error) { return c.vol.isMuted() }
func (c *controller) PlayPause() error       { return c.media.playPause() }
func (c *controller) Next() error            { return c.media.next() }
func (c *controller) Previous() error        { return c.media.previous() }
