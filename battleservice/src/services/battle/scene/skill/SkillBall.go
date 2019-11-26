package skill

// 球技能

import (
	b3core "battleservice/src/services/base/behavior3go/core"
	"battleservice/src/services/battle/scene/bll"
	"battleservice/src/services/battle/scene/interfaces"
	"battleservice/src/services/battle/scene/plr"
	"battleservice/src/services/battle/usercmd"
)

type SkillBall struct {
	blackboard *b3core.Blackboard
	bevTree    *b3core.BehaviorTree
	beginFrame uint32

	player *plr.ScenePlayer
	ball   *bll.BallSkill
}

func NewISkillBall(player *plr.ScenePlayer, ball *bll.BallSkill) interfaces.ISkillBall {
	return NewSkillBall(player, ball)
}

//新球
func NewSkillBall(player *plr.ScenePlayer, ball *bll.BallSkill) *SkillBall {
	skill := SkillBall{
		player: player,
		ball:   ball,
	}

	skill.blackboard = b3core.NewBlackboard()
	skill.blackboard.Set("castskill", float64(0), "", "")
	skill.blackboard.Set("player", player, "", "")
	skill.blackboard.Set("ballskill", &skill, "", "")

	return &skill
}

func (this *SkillBall) CastSkill(skillid uint32) {
	this.beginFrame = uint32(GetEndFrame(this.player, 1)) + 1
	this.bevTree = GetSkillBevTree(skillid)
}

func (this *SkillBall) Update() {
	if this.bevTree != nil && this.player.Frame() >= this.beginFrame {
		this.bevTree.Tick(nil, this.blackboard)
		skillid := uint32(this.blackboard.Get("castskill", "", "").(float64))
		if skillid == 0 {
			this.bevTree = nil
		}
	}
}

func (this *SkillBall) IsFinish() bool {
	return this.bevTree == nil
}

func (this *SkillBall) TryGetHit(player *plr.ScenePlayer) bool {
	if this.ball.GetBallType() == usercmd.BallType_SkillBomb {
		return true
	}
	return false
}

func (this *SkillBall) GetHit(player *plr.ScenePlayer) {
	skillid := uint32(this.blackboard.Get("castskill", "", "").(float64))
	if skillid == 0 {
		return
	}

	if this.ball.GetBallType() == usercmd.BallType_SkillBomb {
		angleVel := this.ball.GetPos()
		angleVel.SubS(player.GetPos())
		angleVel.Normalize()

		angleVel.MulS(0.34)
		this.ball.SetSpeed(&angleVel)
		if this.bevTree != nil {
			this.bevTree.Close(nil, this.blackboard)
		}
		this.bevTree = GetGetHitBevTree(skillid)
	}
}

func (this *SkillBall) GetBeginFrame() uint32 {
	return this.beginFrame
}
