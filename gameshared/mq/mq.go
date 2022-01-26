package mq

type Options struct {
	aecb func(msg *Msg, err error)
}

func GetDefaultOptions() Options {
	return Options{
		aecb: nil,
	}
}

type Option func(*Options) error

func PublishAsyncErrHandler(cb func(msg *Msg, err error)) Option {
	return func(o *Options) error {
		o.aecb = cb
		return nil
	}
}
