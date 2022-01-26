package mq

type nullProducer struct {
}

func newNullProducer(address []string) (*nullProducer, error) {
	return &nullProducer{}, nil
}

func (p *nullProducer) PublishAsync(m *Msg) {

}

func (p *nullProducer) Close() {

}
