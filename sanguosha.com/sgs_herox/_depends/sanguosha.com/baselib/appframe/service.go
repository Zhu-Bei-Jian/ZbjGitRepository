package appframe

import (
	"errors"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/baselib/framework/worker"
	"sanguosha.com/baselib/framework/netframe"
)

// Service 服务通信接口.
// 一个 Service 对应同一类型的一组服务节点, Service 的具体实现负责消息路由和负载平衡.
type Service interface {
	// 服务类型
	Type() ServerType

	// 服务当前是否可用 (对一个节点来说, 就是没有退出也没有断线)
	Available() bool

	// SendMsg 给服务发送消息.
	SendMsg(msg proto.Message) error


	SendMsgWithSeqID(msg proto.Message,seqID int64) error

	// Request 对服务进行请求, 请求消息中必须有成员字段 Seqid int64.
	// 异步回调, 回调函数将会在 app 的 worker 中执行.
	// 返回值用于取消等待当前请求的响应.
	Request(msg proto.Message, cbk func(resp proto.Message, err error), timeout time.Duration) (cancel func())

	// RequestCall 对服务进行请求调用, 同步阻塞, 请求消息中必须有成员字段 Seqid int64.
	RequestCall(msg proto.Message, timeout time.Duration) (proto.Message, error)


	//优化版rpc
	ReqSugar(msg proto.Message,cbk interface{},duration time.Duration) (cancel func())
	CallSugar(req proto.Message, resp proto.Message,duration time.Duration)error
	

	// 转发来自用户会话的消息给该服务. traceId用于监控，可为空
	ForwardMsgFromSession(msg proto.Message,extend netframe.Server_Extend) error

	// 转发来自用户会话的原始消息给该服务. traceId用于监控，可为空
	ForwardRawMsgFromSession(msgid uint32, data []byte,extend netframe.Server_Extend) error
}

// LoadBalance 接口, 用于 Service 负载均衡实现.
type LoadBalance interface {
	Available() bool
	GetServerID(msg proto.Message) (uint32, bool)
	GetServerIDForwardMsgFromSession(session uint64, msg proto.Message) (uint32, bool)
	GetServerIDForwardRawMsgFromSession(session uint64, msgid uint32) (uint32, bool)
	OnServerEvent(svrid uint32, event netcluster.SvrEvent)
}

var (
	// ErrUnregisteredService 未注册的服务.
	ErrUnregisteredService = errors.New("ErrUnregisteredService")
	// ErrNoAvailableServer 没有可用的服务节点.
	ErrNoAvailableServer = errors.New("ErrNoAvailableServer")
)

// unregisteredService 未注册的服务, 用于简化错误处理.
type unregisteredService struct {
	app *Application
}

func (s unregisteredService) Type() ServerType {
	return 0
}
func (s unregisteredService) Available() bool {
	return false
}
func (s unregisteredService) SendMsg(msg proto.Message) error {
	return ErrUnregisteredService
}

func (s unregisteredService) SendMsgWithSeqID(msg proto.Message,seqId int64) error {
	return ErrUnregisteredService
}

func (s unregisteredService) SendUserMsg(userid uint64, msg proto.Message) error {
	return ErrUnregisteredService
}
func (s unregisteredService) ForwardUserMsg(userid uint64, msgid uint32, data []byte) error {
	return ErrUnregisteredService
}
func (s unregisteredService) Request(msg proto.Message, cbk func(resp proto.Message, err error), timeout time.Duration) func() {
	// 需要保证, 回调函数在 app 的 worker 中回调.
	return worker.AfterPost(s.app, 0, func() {
		cbk(nil, ErrUnregisteredService)
	})
}
func (s unregisteredService) RequestCall(msg proto.Message, timeout time.Duration) (proto.Message, error) {
	return nil, ErrUnregisteredService
}
func (s unregisteredService) ForwardMsgFromSession(msg proto.Message,extend netframe.Server_Extend) error {
	return ErrUnregisteredService
}
func (s unregisteredService) ForwardRawMsgFromSession(msgid uint32, data []byte,extend netframe.Server_Extend) error {
	return ErrUnregisteredService
}
func (s unregisteredService) ReqSugar(msg proto.Message, cbk interface{},duration time.Duration) (func()) {
	return nil
}

func (s unregisteredService) CallSugar(msg proto.Message, resp proto.Message,duration time.Duration) (error) {
	return ErrUnregisteredService
}
// 通用的服务基础结构.
type commonService struct {
	typ        ServerType
	app        *Application
	loadBlance LoadBalance
}

func (s *commonService) Type() ServerType {
	return s.typ
}
func (s *commonService) Available() bool {
	return s.loadBlance.Available()
}
func (s *commonService) SendMsg(msg proto.Message) error {
	svrid, ok := s.loadBlance.GetServerID(msg)
	if !ok {
		return ErrNoAvailableServer
	}
	return s.app.slave.SendServerMsg(msg,netframe.Server_Extend{ServerId: svrid})
}

func (s *commonService) SendMsgWithSeqID(msg proto.Message,seqId int64) error {
	svrid, ok := s.loadBlance.GetServerID(msg)
	if !ok {
		return ErrNoAvailableServer
	}
	return s.app.slave.SendServerMsg(msg,netframe.Server_Extend{ServerId: svrid,SeqId: seqId})
}

func (s *commonService) Request(msg proto.Message, cbk func(resp proto.Message, err error), timeout time.Duration) (cancel func()) {
	return s.app.reqc.Req(s.SendMsgWithSeqID, msg, cbk, timeout)
}
func (s *commonService) RequestCall(msg proto.Message, timeout time.Duration) (proto.Message, error) {
	return s.app.reqc.Call(s.SendMsgWithSeqID, msg, timeout)
}
func (s *commonService) ForwardMsgFromSession(msg proto.Message,extend netframe.Server_Extend) error {
	svrid, ok := s.loadBlance.GetServerIDForwardMsgFromSession(extend.SessionId, msg)
	if !ok {
		return ErrNoAvailableServer
	}
	extend.ServerId = svrid

	return s.app.slave.SendServerMsg(msg,extend)
}
func (s *commonService) ForwardRawMsgFromSession(msgid uint32, data []byte,extend netframe.Server_Extend) error {
	svrid, ok := s.loadBlance.GetServerIDForwardRawMsgFromSession(extend.SessionId, msgid)
	if !ok {
		return ErrNoAvailableServer
	}
	extend.ServerId = svrid

	return s.app.slave.SendServerBytes(msgid, data,extend)
}

func (s *commonService)ReqSugar(msg proto.Message,cbk interface{},duration time.Duration) (cancel func()){
	return s.app.reqc.ReqSugar(s.SendMsgWithSeqID,msg,cbk,duration)
}

func (s *commonService)CallSugar(msg proto.Message,resp proto.Message,duration time.Duration) (error){
	return s.app.reqc.CallSugar(s.SendMsgWithSeqID,msg,resp,duration)
}

// Application 和 GateApplication 都实现了这个接口.
type iGetAvailableServerIDs interface {
	GetAvailableServerIDs(typ ServerType) []uint32
}

// WithLoadBalanceSingleton 单节点服务负载均衡策略 (即无策略), 目的是统一抽象.
// 参数 app 可以是 *Application, 也可以是 *GateApplication 对象.
func WithLoadBalanceSingleton(app iGetAvailableServerIDs, svrtype ServerType) LoadBalance {
	return &singleton{app: app, typ: svrtype}
}

// WithLoadBalanceRandom 随机选择服务节点负载均衡策略.
// 参数 app 可以是 *Application, 也可以是 *GateApplication 对象.
func WithLoadBalanceRandom(app iGetAvailableServerIDs, svrtype ServerType) LoadBalance {
	return &random{app: app, typ: svrtype}
}

// 单节点服务负载均衡策略实现 (即无策略)
type singleton struct {
	app iGetAvailableServerIDs
	typ ServerType
	_id uint32
}

func (s *singleton) Type() ServerType {
	return s.typ
}
func (s *singleton) Available() bool {
	return len(s.app.GetAvailableServerIDs(s.typ)) > 0
}
func (s *singleton) GetServerID(msg proto.Message) (uint32, bool) {
	id := atomic.LoadUint32(&s._id)
	if id != 0 {
		return id, true
	}
	ids := s.app.GetAvailableServerIDs(s.typ)
	if len(ids) > 0 {
		id = ids[0]
		if atomic.CompareAndSwapUint32(&s._id, 0, id) {
			return id, true
		}
		return s.GetServerID(msg)
	}
	return 0, false
}
func (s *singleton) GetServerIDForwardMsgFromSession(session uint64, msg proto.Message) (uint32, bool) {
	return s.GetServerID(msg)
}
func (s *singleton) GetServerIDForwardRawMsgFromSession(session uint64, msgid uint32) (uint32, bool) {
	return s.GetServerID(nil)
}
func (s *singleton) OnServerEvent(svrid uint32, event netcluster.SvrEvent) {}

// 随机选择服务节点负载均衡策略实现.
type random struct {
	app iGetAvailableServerIDs
	typ ServerType
}

func (r *random) Available() bool {
	return len(r.app.GetAvailableServerIDs(r.typ)) > 0
}
func (r *random) GetServerID(msg proto.Message) (uint32, bool) {
	ids := r.app.GetAvailableServerIDs(r.typ)
	l := len(ids)
	if l == 0 {
		return 0, false
	}
	return ids[rand.Intn(l)], true
}
func (r *random) GetServerIDForwardMsgFromSession(session uint64, msg proto.Message) (uint32, bool) {
	return r.GetServerID(msg)
}
func (r *random) GetServerIDForwardRawMsgFromSession(session uint64, msgid uint32) (uint32, bool) {
	return r.GetServerID(nil)
}
func (r *random) OnServerEvent(svrid uint32, event netcluster.SvrEvent) {}
