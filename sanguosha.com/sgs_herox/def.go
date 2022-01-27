package sgs_herox

import "sanguosha.com/baselib/appframe"

const AppID = 100

const (
	SvrTypeGate    appframe.ServerType = 2 // SvrTypeGate 网关服务
	SvrTypeAuth    appframe.ServerType = 3 // SvrTypeAuth 认证服务
	SvrTypeLobby   appframe.ServerType = 4 // SvrTypeLobby 大厅服务.
	SvrTypeEntity  appframe.ServerType = 5 // SvrTypeEntity 用户数据服务.
	SvrTypeGame    appframe.ServerType = 6 // SvrTypeGame 游戏服务.
	SvrTypeAI      appframe.ServerType = 9 // SvrTypeAI ai服务
	SvrTypeAdmin   appframe.ServerType = 11
	SvrTypeAccount appframe.ServerType = 15 //账号服，无状态，可开多个
	SvrTypeEnd     appframe.ServerType = 18 //用于标记结束，加服务时，要修改
)
