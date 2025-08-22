package barista

type Packet interface {
	Type() string
	// TODO decouple serialization?
	Bytes() []byte
}

type SubscribePacket struct {
	Topic string
}

func (p *SubscribePacket) Type() string {
	return "SUBSCRIBE"
}

func (p *SubscribePacket) Bytes() []byte {
	return []byte(p.Type() + " " + p.Topic)
}

type PublishPacket struct {
	Content string
	Topic   string
}

func (p *PublishPacket) Type() string {
	return "PUBLISH"
}

func (p *PublishPacket) Bytes() []byte {
	return []byte(p.Type() + " " + p.Topic + " " + p.Content)
}
