package koth

import (
	"math"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
)

type KOTHHandler struct {
	player.NopHandler
	p *player.Player
}

func NewKothHandler(p *player.Player) *KOTHHandler {
	return &KOTHHandler{
		p: p,
	}
}

func (*KOTHHandler) Name() string { return "KothHandler" }

func actuallyMoved(old, new mgl64.Vec3) bool {
	return old.X() != new.X() || old.Z() != new.Z()
}
func (k *KOTHHandler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	if actuallyMoved(k.p.Position(), newPos) {
		k.p.SendTip(math.Round(newPos[0]), math.Round(newPos[2]))
		for _, koth := range koths {
			if _, ok := IsStarted(koth); ok {
				if InCaptureArea(newPos, koth) {
					StartCapturing(k.p, koth)
				} else {
					StopCapturing(k.p, koth)
				}
			}
		}
	}
}
