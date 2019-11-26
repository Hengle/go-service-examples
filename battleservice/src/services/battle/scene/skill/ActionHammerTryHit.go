package skill

import (
	b3 "battleservice/src/services/base/behavior3go"
	b3config "battleservice/src/services/base/behavior3go/config"
	b3core "battleservice/src/services/base/behavior3go/core"
	"battleservice/src/services/battle/scene/plr"

	_ "github.com/cihub/seelog"
)

type ActionHammerTryHit struct {
	b3core.Action
	scale  float64
	gethit uint32
}

func (this *ActionHammerTryHit) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.scale = setting.GetProperty("scale")
	this.gethit = uint32(setting.GetPropertyAsInt("gethit"))
}

func (this *ActionHammerTryHit) OnTick(tick *b3core.Tick) b3.Status {
	ballskill := tick.Blackboard.Get("ballskill", "", "").(*SkillBall).ball
	player := tick.Blackboard.Get("player", "", "").(*plr.ScenePlayer)
	hits := tick.Blackboard.Get("hits", "", "").(map[uint64]int)
	scene := player.GetScene()

	attckRect := ballskill.GetRect()
	attckRect.SetRadius(float64(ballskill.GetRadius()) + 0.5)
	//cells := scene.GetAreaCells(attckRect)

	scene.TravsalPlayers(func(other *plr.ScenePlayer) {
		if other.GetId() == player.GetId() {
			return
		}
		if _, ok := hits[other.GetID()]; ok {
			return
		}
		if BallSkillAttack(tick, player, ballskill, this.scale, &other.BallPlayer) {
			hits[other.GetID()] = 1
			x, _, z := ballskill.GetPos()
			other.Skill.GetHit2(float64(x), float64(z), this.gethit)
		}
	})

	//TODO wei: cell
	// for _, cell := range cells {
	// 	for _, feed := range cell.Feeds {
	// 		distance := feed.GetPosV().SqrMagnitudeTo(ballskill.GetPosV())
	// 		tmp := feed.GetRadius() + ballskill.GetRadius()
	// 		if distance <= tmp*tmp {
	// 			return b3.FAILURE
	// 		}
	// 	}
	// }

	return b3.SUCCESS
}
