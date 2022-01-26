package fsm

import (
	"fmt"
	"strings"
)

// 代码修改自 https://github.com/looplab/fsm
// 1. 修改了组建状态机的方式 (修改为通过 State 和 Transition 来组建, 而不是通过 Event 来组建, 这太反人类了...)
// 2. 增加了全局跳转功能, "*".
// 3. 增加了一些错误处理
// 4. 删除了一些不需要的功能 (比如锁 和 所谓的异步)
// 5. 支持了自跳转和内部事件.
// 6. 增加了并行状态的支持.

// FSM ...
type FSM struct {
	current     string
	transitions map[eKey]string
	callbacks   map[cKey]Callback
	transit     bool
	// 内部事件, 特定状态的下内部处理事件.
	internalEvents map[eKey]Callback

	OnError func(error)
}

// Callback ...
type Callback func(*Event)

// Callbacks ...
type Callbacks map[string]Callback

// State ...
type State struct {
	Transitions    map[string]string
	InternalEvents map[string]Callback
	OnEnter        Callback
	OnLeave        Callback
}

// States 状态组, key 为状态名.
type States map[string]State

// Transitions 状态之间的转移, key 是事件, value 是目标状态.
type Transitions map[string]string

// NewFSM 通过 State 和 Transition 来的方式来组建状态机, 更方便管理, 并可以通过 "*" 标识符支持全局跳转.
func NewFSM(initial string, states map[string]State, otherCallbacks map[string]Callback) *FSM {
	f := &FSM{
		current:        initial,
		transitions:    make(map[eKey]string),
		callbacks:      make(map[cKey]Callback),
		internalEvents: make(map[eKey]Callback),
	}

	allStates := states
	allEvents := make(map[string]bool)
	for sName, s := range states {
		if sName == "" {
			panic("state name is empty")
		}
		for eName, dstName := range s.Transitions {
			if eName == "" || eName == "*" {
				panic(fmt.Errorf("state %s transition event name invalid", sName))
			}
			if dstName == "*" {
				panic(fmt.Errorf("state %s transition dst name invalid", sName))
			}
			f.transitions[eKey{eName, sName}] = dstName
			allEvents[eName] = true
		}
		for eName, cbk := range s.InternalEvents {
			if eName == "" || eName == "*" {
				panic(fmt.Errorf("state %s internal event name invalid", sName))
			}
			if cbk == nil {
				panic(fmt.Errorf("state %s internal event cbk is nil", sName))
			}
			f.transitions[eKey{eName, sName}] = ""
			f.internalEvents[eKey{eName, sName}] = cbk
		}
		if s.OnEnter != nil {
			f.callbacks[cKey{sName, cbkEnterState}] = s.OnEnter
		}
		if s.OnLeave != nil {
			f.callbacks[cKey{sName, cbkLeaveState}] = s.OnLeave
		}
	}

	for name, c := range otherCallbacks {
		if c == nil {
			panic(fmt.Errorf("%s cbk is nil", name))
		}

		var target string
		var cbkType int

		switch {
		case strings.HasPrefix(name, "before_"):
			target = strings.TrimPrefix(name, "before_")
			if target == "*" {
				cbkType = cbkBeforeEvent
			} else if _, ok := allEvents[target]; ok {
				cbkType = cbkBeforeEvent
			}
		case strings.HasPrefix(name, "leave_"):
			target = strings.TrimPrefix(name, "leave_")
			if target == "*" {
				cbkType = cbkLeaveState
			} else if _, ok := allStates[target]; ok {
				cbkType = cbkLeaveState
			}
		case strings.HasPrefix(name, "enter_"):
			target = strings.TrimPrefix(name, "enter_")
			if target == "*" {
				cbkType = cbkEnterState
			} else if _, ok := allStates[target]; ok {
				cbkType = cbkEnterState
			}
		case strings.HasPrefix(name, "after_"):
			target = strings.TrimPrefix(name, "after_")
			if target == "*" {
				cbkType = cbkAfterEvent
			} else if _, ok := allEvents[target]; ok {
				cbkType = cbkAfterEvent
			}
		default:
			target = name
			if _, ok := allStates[target]; ok {
				cbkType = cbkEnterState
			} else if _, ok := allEvents[target]; ok {
				cbkType = cbkAfterEvent
			}
		}

		if cbkType == cbkInvalid {
			panic(fmt.Errorf("invalid callback type %s", name))
		}
		if _, exist := f.callbacks[cKey{target, cbkType}]; exist {
			panic(fmt.Errorf("callback %s %d already exist", target, cbkType))
		}

		f.callbacks[cKey{target, cbkType}] = c
	}

	return f
}

// Current ...
func (f *FSM) Current() string {
	return f.current
}

// Is ...
func (f *FSM) Is(state string) bool {
	return state == f.current
}

// Can ...
func (f *FSM) Can(event string) bool {
	if f.transit {
		return false
	}
	_, ok := f.transitions[eKey{event, f.current}]
	if !ok {
		_, ok = f.transitions[eKey{event, "*"}]
	}
	return ok
}

// Event ...
func (f *FSM) Event(event string, args ...interface{}) error {
	dst, ok := f.transitions[eKey{event, f.current}]
	if !ok {
		// 对全局跳转进行支持.
		dst, ok = f.transitions[eKey{event, "*"}]
		if !ok {
			for ekey := range f.transitions {
				if ekey.event == event {
					err := &InvalidEventError{event, f.current}
					if f.OnError != nil {
						f.OnError(err)
					}
					return err
				}
			}
			err := &UnknownEventError{event}
			if f.OnError != nil {
				f.OnError(err)
			}
			return err
		}
	}

	e := &Event{f, event, f.current, dst, nil, args, false}

	if dst == "" {
		// 内部事件处理.
		if cbk, ok := f.internalEvents[eKey{event, f.current}]; ok {
			cbk(e)
			return nil
		}
		if cbk, ok := f.internalEvents[eKey{event, "*"}]; ok {
			cbk(e)
			return nil
		}
		return &InternalEventError{Event: event, State: f.current}
	}

	if f.transit {
		err := &InTransitionError{Event: event}
		if f.OnError != nil {
			f.OnError(err)
		}
		return err
	}

	f.transit = true

	err := f.beforeEventCallbacks(e)
	if err != nil {
		f.transit = false
		if f.OnError != nil {
			f.OnError(err)
		}
		return err
	}

	err = f.leaveStateCallbacks(e)
	if err != nil {
		f.transit = false
		if f.OnError != nil {
			f.OnError(err)
		}
		return err
	}

	f.current = dst

	f.enterStateCallbacks(e)

	f.afterEventCallbacks(e)

	f.transit = false
	return e.Err
}

func (f *FSM) beforeEventCallbacks(e *Event) error {
	if fn, ok := f.callbacks[cKey{e.Event, cbkBeforeEvent}]; ok {
		fn(e)
		if e.canceled {
			return &CanceledError{e.Err}
		}
	}
	if fn, ok := f.callbacks[cKey{"*", cbkBeforeEvent}]; ok {
		fn(e)
		if e.canceled {
			return &CanceledError{e.Err}
		}
	}
	return nil
}

func (f *FSM) leaveStateCallbacks(e *Event) error {
	if fn, ok := f.callbacks[cKey{f.current, cbkLeaveState}]; ok {
		fn(e)
		if e.canceled {
			return &CanceledError{e.Err}
		}
	}
	if fn, ok := f.callbacks[cKey{"*", cbkLeaveState}]; ok {
		fn(e)
		if e.canceled {
			return &CanceledError{e.Err}
		}
	}
	return nil
}

func (f *FSM) enterStateCallbacks(e *Event) {
	if fn, ok := f.callbacks[cKey{f.current, cbkEnterState}]; ok {
		fn(e)
	}
	if fn, ok := f.callbacks[cKey{"*", cbkEnterState}]; ok {
		fn(e)
	}
}

func (f *FSM) afterEventCallbacks(e *Event) {
	if fn, ok := f.callbacks[cKey{e.Event, cbkAfterEvent}]; ok {
		fn(e)
	}
	if fn, ok := f.callbacks[cKey{"*", cbkAfterEvent}]; ok {
		fn(e)
	}
}

const (
	cbkInvalid int = iota
	cbkBeforeEvent
	cbkLeaveState
	cbkEnterState
	cbkAfterEvent
)

type cKey struct {
	target       string
	callbackType int
}

type eKey struct {
	event string
	src   string
}

// StatesGroup ...
type StatesGroup map[string]States

const (
	// EventParallelStart ...
	EventParallelStart = "parallel_start"
	// EventParallelFinish ...
	EventParallelFinish = "parallel_finish"
	// EventParallelTerminate ...
	EventParallelTerminate = "parallel_terminate"
	// EventParallelDone ...
	EventParallelDone = "parallel_done"
)

// NewParallelState 并行子状态.
func NewParallelState(finishEvent string, statesGroup map[string]States) State {
	var parentFSM *FSM

	internalEvents := Callbacks{}
	ms := make(map[string]*FSM, len(statesGroup))
	for name, states := range statesGroup {
		// 这个名字即是子状态机的名字, 也是该子状态机的初始化状态.
		name := name

		// 错误检查.
		if initState, ok := states[name]; ok {
			_, exist := initState.Transitions["parallel_start"]
			if !exist {
				_, exist = initState.InternalEvents["parallel_start"]
				if !exist {
					panic(fmt.Errorf("parallel state group %s init state not handle event parallel_start", name))
				}
			}
		} else {
			panic(fmt.Errorf("parallel state group %s no init state", name))
		}

		for k, state := range states {
			if state.InternalEvents == nil {
				state.InternalEvents = Callbacks{}
			}
			state.InternalEvents["parallel_done"] = func(e *Event) {
				// 报告给父状态机该子状态机完成了工作可以结束了.
				parentFSM.Event("parallel_done", name)
			}
			states[k] = state
		}

		// 全局跳转, 用于终止.
		states["*"] = State{
			Transitions: Transitions{
				"parallel_terminate": name,
				"parallel_finish":    name,
			},
		}

		m := NewFSM(name, states, nil)

		handled := map[string]bool{}
		for _, state := range states {
			for k := range state.InternalEvents {
				if k == "parallel_done" {
					continue
				}
				if !handled[k] {
					handled[k] = true
					if f, exist := internalEvents[k]; !exist {
						internalEvents[k] = func(e *Event) {
							m.Event(e.Event, e.Args...)
						}
					} else {
						internalEvents[k] = func(e *Event) {
							f(e)
							m.Event(e.Event, e.Args...)
						}
					}
				}
			}
		}

		ms[name] = m
	}

	counter := map[string]bool{}
	internalEvents["parallel_done"] = func(e *Event) {
		// 结束子状态机的工作, 并回到初始化状态.
		name := e.Args[0].(string)
		ms[name].Event("parallel_finish")

		// 如果全部完成工作则发送 finishEvent.
		counter[name] = true
		if len(counter) == len(ms) {
			e.FSM.Event(finishEvent)
		}
	}

	parallelState := State{
		Transitions:    Transitions{},
		InternalEvents: internalEvents,
		OnEnter: func(e *Event) {
			parentFSM = e.FSM
			for _, m := range ms {
				m.Event("parallel_start", e)
			}
		},
		OnLeave: func(e *Event) {
			for k, m := range ms {
				if !counter[k] {
					m.Event("parallel_terminate", e)
				}
			}
			for k := range counter {
				delete(counter, k)
			}
			parentFSM = nil
		},
	}

	return parallelState
}
