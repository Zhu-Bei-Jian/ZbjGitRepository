package mq

type Config struct {
	Open    bool
	Type    string
	Address []string
}

type Msg struct {
	Topic string
	Data  []byte
}
