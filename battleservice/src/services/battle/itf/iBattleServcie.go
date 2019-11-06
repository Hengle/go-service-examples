package itf

import (
	"battleservice/src/services/battle/types"

	"github.com/giant-tech/go-service/framework/iserver"
)

// IBattleService battle服接口
type IBattleService interface {
	iserver.IServiceBase

	//验证token
	LookupToken(token types.Token) (bool, uint64, uint64)
}

// IScene 房间接口
type IScene interface {
	iserver.IEntityGroup
}
