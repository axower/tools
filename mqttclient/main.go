package main

import (
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

var address string
var username string
var password string
var topic string
var qos uint

var logger *zap.Logger

func main() {
	initFlag()
	initLogger()

	over := make(chan struct{})

	clientID := genClientID()
	opts := mqtt.NewClientOptions()
	// 设置 Broker 地址
	opts.AddBroker(address)
	// 设置客户端ID
	opts.SetClientID(clientID)
	// 设置会话保持时长
	opts.SetKeepAlive(60 * time.Second)
	// 设置账号密码
	opts.SetUsername(username)
	opts.SetPassword(password)
	// 设置客户端连接时处理方法（初始化连接、自动重连均会出发）
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logger.Info("client is connect success")
	})
	// 设置客户端连接丢失处理方法
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logger.Warn("client is disconnect", zap.Error(err))
	})

	logger.Info("client info", zap.String("brokerAddress", address), zap.String("clientID", clientID), zap.String("subscribeTopic", topic), zap.Uint("Qos", qos))
	// 创建客户端
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Error("client connect error", zap.Error(token.Error()))
	}

	// 订阅
	token := client.Subscribe(topic, byte(qos), func(client mqtt.Client, message mqtt.Message) {
		logger.Info("handel subscribe", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))
	})
	if token.Wait() && token.Error() != nil {
		logger.Error("subscribe error", zap.Error(token.Error()))
	}

	<-over
}

func initFlag() {
	flag.StringVar(&address, "address", "tcp://127.0.0.1:1883", "mqtt broker address")
	flag.StringVar(&username, "u", "", "mqtt broker username")
	flag.StringVar(&password, "p", "", "mqtt broker password")
	flag.StringVar(&topic, "topic", "#", "subscribe topic")
	flag.UintVar(&qos, "qos", 0, "mqtt qos")
	flag.Parse()
}

func genClientID() string {
	return fmt.Sprintf("mc_%d", time.Now().Unix())
}

func initLogger() {
	en := zap.NewProductionEncoderConfig()
	en.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	en.ConsoleSeparator = " | "
	en.EncodeLevel = zapcore.CapitalLevelEncoder

	conf := zap.NewProductionConfig()
	conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	conf.EncoderConfig = en
	conf.Encoding = "console"
	var err error
	logger, err = conf.Build(zap.WithCaller(false))
	if err != nil {
		log.Fatal(err)
	}
}
