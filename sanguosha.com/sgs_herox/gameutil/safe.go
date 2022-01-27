package gameutil

import (
	"github.com/pkg/errors"
	"sanguosha.com/baselib/util"
	"time"
)

var (
	ErrDivZero = errors.New("ErrDivZero")
)

func SafeCall(f func()) {
	defer util.Recover()
	f()
}

func SafeCallAfter(t time.Duration, f func()) *time.Timer {
	return time.AfterFunc(t, func() {
		defer util.Recover()
		f()
	})
}

func SafeDivTimeDuration(a, b time.Duration) time.Duration {
	//if b == 0 {
	//	defer util.Recover()
	//	panic(ErrDivZero)
	//	return 0
	//}
	return a / b
}

func SafeDivInt(a, b int) int {
	//if b == 0 {
	//	defer util.Recover()
	//	panic(ErrDivZero)
	//	return 0
	//}
	return a / b
}

func SafeDivInt32(a, b int32) int32 {
	//if b == 0 {
	//	defer util.Recover()
	//	panic(ErrDivZero)
	//	return 0
	//}
	return a / b
}

func SafeDivInt64(a, b int64) int64 {
	//if b == 0 {
	//	//	defer util.Recover()
	//	//	panic(ErrDivZero)
	//	//	return 0
	//	//}
	return a / b
}

func SafeDivUint64(a, b uint64) uint64 {
	//if b == 0 {
	//	defer util.Recover()
	//	panic(ErrDivZero)
	//	return 0
	//}
	return a / b
}

func SafeDivFloat32(a, b float32) float32 {
	//if b == 0 {
	//	defer util.Recover()
	//	panic(ErrDivZero)
	//	return 0
	//}
	return a / b
}

func SafeDivFloat64(a, b float64) float64 {
	//if b == 0 {
	//	defer util.Recover()
	//	panic(ErrDivZero)
	//	return 0
	//}
	return a / b
}
