package main

type Packet interface {
	Type() string
}

type SubscribePacket struct {
	Topic string
}

func (p *SubscribePacket) Type() string {
	return "SUBSCRIBE"
}

type PublishPacket struct {
	Content string
	Topic   string
}

func (p *PublishPacket) Type() string {
	return "PUBLISH"
}
