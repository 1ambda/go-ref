package distributed

const SingleKeyWsConnectionCount = "gateway/single/count-websocket-connection"
const SingleKeyTotalAccessCount = "gateway/single/count-total-access"
const SingleKeyLeaderName = "gateway/single/leader-name"

const RangeKeyPrefixWebSocket = "gateway/range/websocket/gateway-"

type DistributedMessage struct {
	key   string
	value string
}

func NewTotalAccessCountMessage(count string) *DistributedMessage {
	return &DistributedMessage{key: SingleKeyTotalAccessCount, value: count}
}
