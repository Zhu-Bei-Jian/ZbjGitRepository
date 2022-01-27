package reporter

import (
	"testing"
	"time"
)

func TestDingDingReporter_Send(t *testing.T) {
	token := "f676d3735196c9d2f90cd023e6bae726e5d1859162b340d0b0751961fdb7aac7"
	p := newDingDingReporter(token)
	p.Send("dingding robot test,hi")
	time.Sleep(time.Second * 5)
}
