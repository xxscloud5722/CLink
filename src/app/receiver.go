package app

type Receiver interface {
	// receiver 消息大小和超时配置, 接收消息并解码.
	receiver(size, timeout int) (*[]LogMessage, error)
}
