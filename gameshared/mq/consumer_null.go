package mq

type nullConsumer struct {
}

func newNullConsumer(address []string) (*nullConsumer, error) {
	p := &nullConsumer{}
	return p, nil
}

func (p *nullConsumer) Close() {
}

// Sub 注册topic回调.
func (p *nullConsumer) Sub(topic string, callback func([]byte)) error {
	return nil
}
