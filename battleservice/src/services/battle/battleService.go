package battle

import (
	"battleservice/src/services/battle/conf"
	"battleservice/src/services/battle/proc"
	"battleservice/src/services/battle/scene"
	"battleservice/src/services/battle/scene/plr"
	"battleservice/src/services/battle/token"
	"fmt"
	"math/rand"
	"time"

	"github.com/cihub/seelog"
	"github.com/spf13/viper"
	"github.com/giant-tech/go-service/base/net/server"
	"github.com/giant-tech/go-service/framework/service"
)

// BattleService 战斗服务
type BattleService struct {
	service.BaseService
	svr *server.Server
	token.TokenMgr
}

// OnInit 初始化
func (bs *BattleService) OnInit() error {
	seelog.Debug("BattleService.OnInit")

	rand.Seed(time.Now().Unix())

	bs.TokenMgr.Init()

	// 注册proto
	bs.RegProtoType("Player", &plr.ScenePlayer{}, false)

	// 全局配置
	if !conf.ConfigMgr_GetMe().Init() {
		seelog.Error("[启动] 读取全局配置失败")
		return fmt.Errorf("conf init failed")
	}

	// 绑定本地端口
	address := viper.GetString("room.listen")
	svr, err := server.New("tcp+kcp", address, 10000)
	if err != nil {
		seelog.Error(fmt.Sprintf("创建服务器失败: ", err))
		return fmt.Errorf("server.New")
	}

	bs.svr = svr

	// 添加MsgProc, 这样新连接创建时会注册处理函数
	bs.svr.AddMsgProc(&proc.ProcBattle{IBattleService: bs})

	if !scene.Init() {
		return fmt.Errorf("scene.Init failed")
	}
	if !conf.InitMapConfig() {
		seelog.Error("[启动]InitMapConfig fail! ")
		return fmt.Errorf("InitMapConfig failed")
	}

	seelog.Info("[启动] 完成初始化")

	return nil
}

// OnTick tick
func (bs *BattleService) OnTick() {

}

// OnDestroy 析构
func (bs *BattleService) OnDestroy() {
	seelog.Debug("BattleService.OnDestroy")

	if bs.svr != nil {
		bs.svr.Close()
	}

}
