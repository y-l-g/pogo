package supervisor

//go:generate go run ../../cmd/genproto/main.go

const (
	PktTypeData     = 0x00
	PktTypeError    = 0x01
	PktTypeFatal    = 0x02
	PktTypeHello    = 0x03
	PktTypeShm      = 0x04
	PktTypeShutdown = 0x09

	MaxPayloadSize = 16 * 1024 * 1024
)
