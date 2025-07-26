package structure

type Message struct {
	payload    []byte
	structured any
}

func FromBytes(data []byte) *Message {
	return &Message{
		payload:    data,
		structured: bytesToStructured(data),
	}
}

func bytesToStructured(data []byte) any {

	return nil
}
