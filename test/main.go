package main

import (
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/dragonfly-on-steroids/area"
	"github.com/dragonfly-on-steroids/claim"
	"github.com/dragonfly-on-steroids/koth"
	"github.com/dragonfly-on-steroids/moreHandlers"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sirupsen/logrus"
)

type testClaim struct {
	claim.NopClaim
}

func (testClaim) Area() area.Area { return area.NewArea(mgl64.Vec2{-10, -10}, mgl64.Vec2{30, 30}) }
func (testClaim) AllowBreakBlock(p *player.Player, pos cube.Pos, drops *[]item.Stack) bool {
	return false
}

type testKoth struct {
	koth.NopKOTH
}

func (*testKoth) Capture(p *player.Player) {
	p.Message("gg")
}
func (*testKoth) CaptureArea() area.Area  { return area.NewArea(mgl64.Vec2{10, 10}, mgl64.Vec2{15, 15}) }
func (*testKoth) Duration() time.Duration { return 5 * time.Second }

func main() {
	c := server.DefaultConfig()
	c.Players.SaveData = false
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel
	s := server.New(&c, log)
	s.Start()

	claim.Register(testClaim{})
	k := &testKoth{}
	koth.Register(k)
	koth.StartKOTH(k)
	for {
		p, err := s.Accept()
		if err != nil {
			return
		}
		claimH := claim.NewClaimHandler(p, claim.NewWilderness(s.World(), "entered wilderness", "left wilderness"))
		kothH := koth.NewKothHandler(p)
		h := moreHandlers.NewPlayerHandler(claimH, kothH)
		p.Handle(h)
	}
}
