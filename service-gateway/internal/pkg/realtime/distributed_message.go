package realtime

const KeyWsConnectionStat = "gateway/stat/wsConnectionCount"
const KeyTotalAccessStat = "gateway/stat/totalAccessCount"
const KeyLeaderNameStat = "gateway/stat/leaderName"

type DistributedMessage struct {
	key string
	value string
}

func NewTotalAccessCountMessage(count string) *DistributedMessage {
	return &DistributedMessage{key: KeyTotalAccessStat, value: count}
}
