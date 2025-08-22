package barista

import (
	"fmt"
	"strings"
)

// TODO validation
// TODO efficiency

type PacketParser interface {
	Parse(data []byte) (Packet, error)
}

type DefaultParser struct{}

func (p *DefaultParser) Parse(data []byte) (Packet, error) {
	parts := strings.Split(string(data), " ")
	packetType := parts[0]
	switch packetType {
	case "SUBSCRIBE":
		return p.parseSubscribe(data)
	case "PUBLISH":
		return p.parsePublish(data)
	default:
		return nil, fmt.Errorf("unsupported packet type: %s", packetType)
	}
}

func (p *DefaultParser) parseSubscribe(data []byte) (Packet, error) {
	parts := strings.Split(string(data), " ")
	topic := parts[1]
	return &SubscribePacket{Topic: topic}, nil
}

func (p *DefaultParser) parsePublish(data []byte) (Packet, error) {
	parts := strings.Split(string(data), " ")
	topic := parts[1]
	content := parts[2]
	return &PublishPacket{Topic: topic, Content: content}, nil
}
