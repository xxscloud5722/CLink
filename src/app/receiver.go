package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/fatih/color"
	"github.com/samber/lo"
	"log"
	"os"
	"os/signal"
	"time"
)

type Receiver interface {
	// receiver 消息大小和超时配置, 接收消息并解码.
	receiver(size int, timeout int64, handler MessageHandler) error
}
type MessageHandler interface {
	Handler(message []*LogMessage) error
}
type KafkaReceiver struct {
	server   string   // server 服务器地址.
	port     int      // port 端口号.
	groupId  string   // groupId 分组ID.
	topic    []string // topic 主题.
	username string   // username SASL 账号.
	password string   // password SASL 密码.
	debug    bool     // debug 是否启用调试.
}

func NewKafka(config *Config) (*KafkaReceiver, error) {
	return &KafkaReceiver{
		server:   config.Kafka.Server,
		port:     config.Kafka.Port,
		groupId:  config.Kafka.Consumer.GroupId,
		topic:    config.Kafka.Consumer.Topic,
		username: config.Kafka.SASL.Username,
		password: config.Kafka.SASL.Password,
		debug:    config.debug,
	}, nil
}

func (kafka *KafkaReceiver) receiver(size int, timeout int64, handler MessageHandler) error {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	kafkaConfig.Net.SASL.Enable = true
	kafkaConfig.Net.SASL.User = kafka.username
	kafkaConfig.Net.SASL.Password = kafka.password
	kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	brokerList := []string{fmt.Sprintf("%s:%d", kafka.server, kafka.port)}
	consumer, err := sarama.NewConsumerGroup(brokerList, kafka.groupId, kafkaConfig)
	if err != nil {
		return err
	}
	defer func() {
		_ = consumer.Close()
	}()
	log.Println("[Receiver] Connect to `Kafka` Success ...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	consumerHandler := &LogConsumerHandler{
		batchSize: size,
		timeout:   timeout,
		handler:   handler,
		debug:     kafka.debug,
	}
	err = consumer.Consume(ctx, kafka.topic, consumerHandler)
	if err != nil {
		return err
	}
	// 等待程序结束.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	if err = consumer.Close(); err != nil {
		return errors.New(fmt.Sprintf("Error closing consumer: %v", err))
	}
	return nil
}

type LogConsumerHandler struct {
	batchSize int
	timeout   int64
	handler   MessageHandler
	debug     bool
}

func (h *LogConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *LogConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *LogConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batch := make([]*sarama.ConsumerMessage, 0, h.batchSize)
	timeout := time.NewTimer(time.Duration(h.timeout) * time.Millisecond)

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			if h.debug {
				log.Println(fmt.Sprintf("[Debug]: %s", string(message.Value)))
			}
			batch = append(batch, message)
			// 数据条数出发事件
			if len(batch) >= h.batchSize {
				err := h.handler.Handler(lo.Map(batch, func(item *sarama.ConsumerMessage, index int) *LogMessage {
					return &LogMessage{Body: item.Value}
				}))
				if err != nil {
					// TODO 错误日志
					color.Red("[Kafka] Consumption failure: %v", err)
				} else {
					for _, item := range batch {
						session.MarkMessage(item, "")
					}
					batch = batch[:0]
					timeout.Reset(time.Duration(h.timeout) * time.Millisecond)
				}
			}
		case <-timeout.C:
			// 计时器触发事件
			if len(batch) > 0 {
				err := h.handler.Handler(lo.Map(batch, func(item *sarama.ConsumerMessage, index int) *LogMessage {
					return &LogMessage{Body: item.Value}
				}))
				if err != nil {
					// TODO 错误日志
					color.Red("[Kafka] Consumption failure: %v", err)
				} else {
					for _, item := range batch {
						session.MarkMessage(item, "")
					}
					batch = batch[:0]
				}
			}
			timeout.Reset(2 * time.Second)
		}
	}
}
