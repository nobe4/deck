//go:build linux

package media

// playerctl implements mediaCtrl via playerctl (MPRIS/D-Bus).
type playerctl struct{}

func (p *playerctl) playPause() error { _, err := run("playerctl", "play-pause"); return err }
func (p *playerctl) next() error      { _, err := run("playerctl", "next"); return err }
func (p *playerctl) previous() error  { _, err := run("playerctl", "previous"); return err }
