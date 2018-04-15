package config

import "time"

const SessionTimeout = 60 * time.Minute
const SessionKey = "sessionID"

const (
	WsCloseReasonUnknown = "Unknown"
	WsCloseFailureClientDisconnected = "ClientDisconnected"
	WsCloseReasonMessageSendFailure  = "MessageSendFailure"
	WsCloseReasonServerShutdown      = "ServerShutdown"
)
