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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		parser.Parse(data)
	}
}
