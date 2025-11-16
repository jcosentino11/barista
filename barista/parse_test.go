package barista

import (
	"testing"
)

func BenchmarkParse(b *testing.B) {
	parser := NewDefaultParser()
	packet := PublishPacket{
		Topic:   "some-topic",
		Content: "example content",
	}
	data := packet.Bytes()

	for b.Loop() {
		parser.Parse(data)
	}
}
