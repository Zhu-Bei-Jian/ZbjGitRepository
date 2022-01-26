package userdb

import (
	"github.com/golang/protobuf/proto"
)

//工作原理：
//当玩家上线时，从db加载此模块数据（依据dbFieldName),使用init接口初始化模块数据
//当模块有数据更新需要存入db时，调用setDirty()，玩家线程会在一定时间后检查是否有模块数据要落地(isDirty())，调用toProtoMessage()接口获取最新模块数据存入db
type Component interface {
	dbFieldName() string           //DB字段的名称，必需实现
	init(*User, []byte) error      //数据导入，初始化，必需实现
	toProtoMessage() proto.Message //数据导出，必需实现

	onLogin()  //登录触发事件
	onLogout() //登出触发事件
	onClock0() //0点事件
	onMinute() //每分钟会调用一次，可选

	//模块内嵌BaseComp即可，不用实现
	isDirty() bool //是否有数据变更
	setDirty()     //设置有数据变更
	clearDirty()   //清dirty
}

type BaseComp struct {
	dirty bool
}

func (p *BaseComp) isDirty() bool {
	return p.dirty
}

func (p *BaseComp) setDirty() {
	p.dirty = true
}

func (p *BaseComp) clearDirty() {
	p.dirty = false
}

func (p *BaseComp) onLogin() {

}

func (p *BaseComp) onLogout() {

}

func (p *BaseComp) onClock0() {

}

func (p *BaseComp) onMinute() {

}
