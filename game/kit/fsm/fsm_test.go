package fsm

import (
	"fmt"
	"testing"
)

// TestFSM ...
func TestFSM(t *testing.T) {
	sm := NewFSM(
		"end",
		States{
			"end": State{
				Transitions: Transitions{
					"start": "start",
				},
				OnEnter: func(e *Event) {
					fmt.Println("game end")
				},
				OnLeave: func(e *Event) {
					fmt.Println("leave game state end")
				},
			},

			"start": State{
				InternalEvents: Callbacks{
					"internal": func(e *Event) {
						fmt.Println("start internal event")
					},
				},
				Transitions: Transitions{
					"play": "playing",
				},
				OnEnter: func(e *Event) {
					fmt.Println("game start")
				},
				OnLeave: func(e *Event) {
					fmt.Println("leave game state start")
				},
			},

			"playing": func() State {
				statesGroup := StatesGroup{
					"s1_init": States{
						"s1_init": State{
							Transitions: Transitions{
								EventParallelStart: "s1",
							},
						},
						"s1": State{
							InternalEvents: Callbacks{
								"s1done": func(e *Event) {
									fmt.Println("s1 do s1done")
									e.FSM.Event(EventParallelDone)
								},
							},
							OnEnter: func(e *Event) {
								fmt.Println("s1 onenter")
							},
							OnLeave: func(e *Event) {
								fmt.Println("s1 onleave")
							},
						},
					},
					"s2_init": States{
						"s2_init": State{
							Transitions: Transitions{
								EventParallelStart: "s2",
							},
						},
						"s2": State{
							InternalEvents: Callbacks{
								"s2done": func(e *Event) {
									fmt.Println("s2 do s2done")
									e.FSM.Event(EventParallelDone)
								},
							},
							OnEnter: func(e *Event) {
								fmt.Println("s2 onenter")
							},
							OnLeave: func(e *Event) {
								fmt.Println("s2 onleave")
							},
						},
					},
				}

				parallel := NewParallelState("finish", statesGroup)
				parallel.Transitions["finish"] = "after_playing"
				onEnter := parallel.OnEnter
				parallel.OnEnter = func(e *Event) {
					fmt.Println("parallel playing onenter")
					onEnter(e)
				}
				onLeave := parallel.OnLeave
				parallel.OnLeave = func(e *Event) {
					onLeave(e)
					fmt.Println("parallel playing onleave")
				}
				return parallel
			}(),

			"after_playing": State{},

			// 全局跳转.
			"*": State{
				Transitions: Transitions{
					"end": "end",
				},
				OnEnter: func(e *Event) {
					//fmt.Println("enter", e.Src, "->", e.Dst)
				},
				OnLeave: func(e *Event) {
					//fmt.Println("leave", e.Src, "->", e.Dst)
				},
			},
		},
		nil,
	)

	err := sm.Event("start")
	if err != nil {
		t.Error(err)
	}
	if !sm.Is("start") {
		t.Fail()
	}
	err = sm.Event("internal")
	if err != nil {
		t.Error(err)
	}
	if !sm.Is("start") {
		t.Fail()
	}

	err = sm.Event("play")
	if err != nil {
		t.Error(err)
	}
	if !sm.Is("playing") {
		t.Fail()
	}
	sm.Event("s1done")
	sm.Event("s2done")
	if !sm.Is("after_playing") {
		t.Fail()
	}

	err = sm.Event("end")
	if err != nil {
		t.Error(err)
	}
	if !sm.Is("end") {
		t.Fail()
	}
}
