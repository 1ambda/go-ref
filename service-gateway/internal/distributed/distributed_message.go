package distributed

const SingleKeyWsConnectionCount = "gateway/single/count-websocket-connection"
const SingleKeyBrowserHistoryCount = "gateway/single/count-browser-history"
const SingleKeyLeaderName = "gateway/single/leader-name"

const RangeKeyPrefixWebSocket = "gateway/range/websocket/gateway-"

type Message struct {
	key   string
	value string
}

func NewBrowserHistoryCountMessage(count string) *Message {
	return &Message{key: SingleKeyBrowserHistoryCount, value: count}
}
