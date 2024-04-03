package app

type Transmitter interface {
	// transmitter 发送器, 推送日志消息到目的地.
	transmitter(*[]LogMessage) error
}
