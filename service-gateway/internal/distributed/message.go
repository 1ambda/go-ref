package distributed

const SingleKeyWsConnectionCount = "service-gateway/single/count-websocket-connection"
const SingleKeyBrowserHistoryCount = "service-gateway/single/count-browser-history"
const SingleKeyLeaderName = "service-gateway/single/leader-name"

const RangeKeyPrefixWebSocket = "service-gateway/range/websocket/gateway-"

type Message struct {
	key   string
	value string
}

func NewBrowserHistoryCountMessage(count string) *Message {
	return &Message{key: SingleKeyBrowserHistoryCount, value: count}
}
