package reporter

import "errors"

type Reporter interface {
	Send(content string)
}

type Type int

const (
	DingDing Type = 1
)

func New(typ Type, token string) (Reporter, error) {
	switch typ {
	case DingDing:
		return newDingDingReporter(token), nil
	default:
		return nil, errors.New("not support type")
	}
}
