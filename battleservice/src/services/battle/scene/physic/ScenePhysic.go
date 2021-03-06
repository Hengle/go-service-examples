// Package physic wraps ape physic engine.
package physic

// 场景物理层

import (
	"battleservice/src/services/base/ape"
	"github.com/cihub/seelog"
)

type ScenePhysic struct {
	engine    *ape.APEngine
	root      *ape.Group // 固定碰撞群
	blockAll  *ape.Group // 所有MapObject
	playerAll *ape.Group // 所有BallPlayer
	feedGroup *ape.Group // 所有BallFeed
}

type IPartical = ape.IAbstractParticle

func NewScenePhysic() *ScenePhysic {
	physic := &ScenePhysic{}
	physic.Init()
	return physic
}

func (this *ScenePhysic) Init() {
	this.engine = ape.NewAPEngine()
	this.engine.Init(float32(0.1))
	this.BuildGroups()
}

func (this *ScenePhysic) Tick() {
	this.engine.Step()
}

func (this *ScenePhysic) CreateBoard(size float32) {
	left := ape.NewRectangleParticle(-size/2, size/2, size, size*2)
	left.SetFixed(true)
	right := ape.NewRectangleParticle(size*3/2, size/2, size, size*2)
	right.SetFixed(true)
	top := ape.NewRectangleParticle(size/2, size*3/2, size*2, size)
	top.SetFixed(true)
	down := ape.NewRectangleParticle(size/2, -size/2, size*2, size)
	down.SetFixed(true)
	fiveColorStone := ape.NewCircleParticle(size, size, 1) //可能引擎问题，右上角会穿透，要堵一下
	fiveColorStone.SetFixed(true)
	this.root.AddParticle(left)
	this.root.AddParticle(right)
	this.root.AddParticle(top)
	this.root.AddParticle(down)
	this.root.AddParticle(fiveColorStone)
}

func (this *ScenePhysic) BuildGroups() {
	this.root = ape.NewGroup(false)
	this.engine.AddGroup(this.root)

	this.blockAll = ape.NewGroup(false)
	this.engine.AddGroup(this.blockAll)

	this.playerAll = ape.NewGroup(true)
	this.engine.AddGroup(this.playerAll)

	this.feedGroup = ape.NewGroup(false)
	this.engine.AddGroup(this.feedGroup)

	this.root.AddCollidable(this.playerAll)
	this.blockAll.AddCollidable(this.playerAll)
	this.feedGroup.AddCollidable(this.playerAll)
}

func (this *ScenePhysic) AddPlayer(player IPartical) {
	this.playerAll.AddParticle(player)
}

func (this *ScenePhysic) RemovePlayer(player IPartical) {
	if player != nil {
		if nil != this.playerAll {
			this.playerAll.RemoveParticle(player)
		}
	}
}

func (this *ScenePhysic) AddFeed(feed IPartical) {
	this.feedGroup.AddParticle(feed)
	if len(this.feedGroup.GetParticles()) > 1000 {
		seelog.Error("feed超过设定的最大值 ", len(this.feedGroup.GetParticles()))
	}
}

func (this *ScenePhysic) RemoveFeed(feed IPartical) {
	this.feedGroup.RemoveParticle(feed)
}

func (this *ScenePhysic) AddBlock(block IPartical) {
	this.blockAll.AddParticle(block)
}
