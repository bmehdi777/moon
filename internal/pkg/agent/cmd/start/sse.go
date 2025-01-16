package start

type SseMessage struct {
	Data string
}

func (m *SseMessage) String() string {
	raw := "data: " + m.Data

	return raw
}
