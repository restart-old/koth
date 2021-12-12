package koth

import (
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/dragonfly-on-steroids/area"
	"github.com/go-gl/mathgl/mgl64"
)

var koths []*KOTH

func Register(k *KOTH) {
	koths = append(koths, k)
}

type KOTH struct {
	world           *world.World
	captureArea     area.Area
	duration        time.Duration
	hMutex          sync.RWMutex
	h               KOTHHandler
	capturing       *player.Player
	shouldCaptureAt time.Time
	started         bool
}

func NewKOTH(world *world.World, captureArea area.Area, duration time.Duration) *KOTH {
	return &KOTH{
		world:       world,
		captureArea: captureArea,
		duration:    duration,
	}
}

func (k *KOTH) Handle(h KOTHHandler) {
	k.hMutex.Lock()
	defer k.hMutex.Unlock()
	if h == nil {
		h = NopHandler{}
	}
	k.h = h
}

func (k *KOTH) World() *world.World     { return k.world }
func (k *KOTH) CaptureArea() area.Area  { return k.captureArea }
func (k *KOTH) Duration() time.Duration { return k.duration }
func (k *KOTH) handler() KOTHHandler    { return k.h }
func (k *KOTH) Capturing() (*player.Player, bool) {
	return k.capturing, k.capturing != nil
}
func (k *KOTH) Start(src Source) {
	if !k.started {
		ctx := event.C()
		k.handler().HandleStart(ctx, src)
		ctx.Continue(func() {
			k.started = true

		})
	}
}
func (k *KOTH) Stop(src Source) {
	if k.started {
		ctx := event.C()
		k.handler().HandleStop(ctx, src)
		ctx.Continue(func() {
			k.started = false
		})
	}

}
func (k *KOTH) StartCapturing(p *player.Player) {
	if k.started {
		if k.capturing != p {
			ctx := event.C()
			k.handler().HandleStartCapturing(ctx, p)
			ctx.Continue(func() {
				k.capturing = p
				k.shouldCaptureAt = time.Now().Add(k.duration)
				time.AfterFunc(k.duration, k.captureFunc(p))
			})
		}
	}
}
func (k *KOTH) StopCapturing(p *player.Player) {
	if k.started {
		if k.capturing == p {
			ctx := event.C()
			k.handler().HandleStopCapturing(ctx, p)
			ctx.Continue(func() {
				k.capturing = nil
				k.shouldCaptureAt = time.Now().Add(43830 * time.Minute)
			})
		}
	}
}
func (k *KOTH) captureFunc(p *player.Player) func() {
	return func() {
		if k.capturing != nil && k.capturing == p {
			if k.shouldCaptureAt.Before(time.Now()) || k.shouldCaptureAt.Equal(time.Now()) {
				ctx := event.C()
				k.h.HandleCapture(ctx, p)
				ctx.Continue(func() {
					k.Stop(SourceCapture{winner: p})
				})
			}
		}
	}
}

func (k *KOTH) InCaptureArea(pos mgl64.Vec3) bool {
	return k.captureArea.Vec2Within(mgl64.Vec2{pos[0], pos[2]})
}
