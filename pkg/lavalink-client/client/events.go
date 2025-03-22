package client

type Message string

var (
	MessageReady        Message = "ready"
	MessagePlayerUpdate Message = "playerUpdate"
	MessageStats        Message = "stats"
	MessageRaw          Message = "event"
)

type MessageHandlers map[Message]MessageHandlerFunc

type MessageHandlerFunc func([]byte) error
