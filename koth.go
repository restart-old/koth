package koth

import (
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/dragonfly-on-steroids/area"
	"github.com/go-gl/mathgl/mgl64"
)

var koths []KOTH

func Register(k KOTH) {
	koths = append(koths, k)
}

type kothStatus struct {
	capturing *player.Player
	captureAt time.Time
}

func newKothStatus() *kothStatus {
	return &kothStatus{
		capturing: nil,
	}
}

func StopCapturing(p1 *player.Player, k KOTH) {
	if p, ok := BeingCaptured(k); ok && p == p1 {
		k.StopCapturing(p)
		SetCapturing(nil, k)
	}
}
func StartCapturing(p1 *player.Player, k KOTH) {
	if p, ok := BeingCaptured(k); !ok && p != p1 {
		k.StartCapturing(p)
		SetCapturing(p1, k)
		time.AfterFunc(k.Duration(), ShouldCaptureFunc(p1, k))
	}
}

func ShouldCaptureFunc(p *player.Player, k KOTH) func() {
	return func() {
		if status, ok := IsStarted(k); ok {
			if status.captureAt.Before(time.Now()) || status.captureAt.Equal(time.Now()) {
				StopKOTH(k)
				k.Capture(p)
			}
		}
	}
}

var running sync.Map

func StartKOTH(k KOTH) {
	running.Store(k, newKothStatus())
	k.Start()
}
func StopKOTH(k KOTH) {
	running.Delete(k)
	k.Stop()
}

type KOTH interface {
	Start()
	Stop()

	StartCapturing(p *player.Player)
	StopCapturing(p *player.Player)

	CaptureArea() area.Area
	Duration() time.Duration

	Capture(p *player.Player)
}

type NopKOTH struct{}

func (NopKOTH) Start()                          {}
func (NopKOTH) Stop()                           {}
func (NopKOTH) StartCapturing(p *player.Player) {}
func (NopKOTH) StopCapturing(p *player.Player)  {}
func (NopKOTH) CaptureArea() area.Area          { return area.Area{} }
func (NopKOTH) Duration() time.Duration         { return 0 }
func (NopKOTH) Capture(p *player.Player)        {}

func IsStarted(k KOTH) (*kothStatus, bool) {
	status, ok := running.Load(k)
	s, ok := status.(*kothStatus)
	return s, ok
}

func BeingCaptured(k KOTH) (*player.Player, bool) {
	if status, ok := IsStarted(k); ok {
		return status.capturing, status.capturing != nil
	}
	return nil, false
}

func InCaptureArea(pos mgl64.Vec3, k KOTH) bool {
	return k.CaptureArea().Vec2Within(mgl64.Vec2{pos[0], pos[2]})
}

func SetCapturing(p *player.Player, k KOTH) {
	if status, ok := IsStarted(k); ok {
		status.capturing = p
		status.captureAt = time.Now().Add(k.Duration())
	}
}
