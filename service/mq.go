package service

/**
积分事件队列的处理
1、基于MQ
2、

*/
import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"github.com/aosfather/bingo/utils"

)

type MQExchangType string

const (
	MQ_EXCHANGE_DIRECT MQExchangType = "direct"
	MQ_EXCHANGE_TOPIC  MQExchangType = "topic"
	MQ_EXCHANGE_HEADER MQExchangType = "headers"
	MQ_EXCHANGE_FANOUT MQExchangType = "fanout"
)

// Receiver 观察者模式需要的接口
// 观察者用于接收指定的queue到来的数据
type Receiver interface {
	QueueName() string     // 获取接收者需要监听的队列
	RouterKey() string     // 这个队列绑定的路由
	OnError(error)         // 处理遇到的错误，当RabbitMQ对象发生了错误，他需要告诉接收者处理错误
	OnReceive([]byte) bool // 处理收到的消息, 这里需要告知RabbitMQ对象消息是否处理成功
}

type RabbitMQConnect struct {
	connection *amqp.Connection
	logger utils.Log
}

func (this *RabbitMQConnect)SetLogger(l utils.Log){
	this.logger=l
}

func (this *RabbitMQConnect) Connect(mqurl string) {
	if this.connection != nil {
		return
	}
	var err error
	this.connection, err = amqp.Dial(mqurl)
	this.failOnErr(err, "failed to connect tp rabbitmq")

}

func (this *RabbitMQConnect)failOnErr(err error, text string) {
	if err != nil {
		this.logger.Error("%s,%s",text,err.Error())
		panic(err.Error())
	}
}
func (this *RabbitMQConnect) GetChannel() *amqp.Channel {
	channel, err := this.connection.Channel()
	this.failOnErr(err, "failed to open a channel")
	return channel
}

func (this *RabbitMQConnect) Close() {
	err := this.connection.Close()
	this.failOnErr(err, "close mq connection error!")
}

func (this *RabbitMQConnect) NewClient(exchangeName string, exchangeType MQExchangType)*RabbitMQ{
	// 这里可以根据自己的需要去定义
	client:= &RabbitMQ{
		exchangeName: exchangeName,
		exchangeType: exchangeType,
	}
	client.channel=this.GetChannel()
	return client
}

// RabbitMQ 用于管理和维护rabbitmq的对象
type RabbitMQ struct {
	wg sync.WaitGroup

	channel      *amqp.Channel
	exchangeName string        // exchange的名称
	exchangeType MQExchangType // exchange的类型
	receivers    []Receiver
	logger utils.Log
}


func (this *RabbitMQ) Start() {
	this.prepareExchange()
    go this.startListen()

}

func (this *RabbitMQ) Push(queque string, msg interface{}) {
	data, _ := json.Marshal(msg)
	theMsg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}
	this.channel.Publish(this.exchangeName, queque, false, false, theMsg)

}

func (this *RabbitMQ) startListen() {
	for _, receiver := range this.receivers {
		this.wg.Add(1)
		go this.listen(receiver) // 每个接收者单独启动一个goroutine用来初始化queue并接收消息
	}

	this.wg.Wait()

	this.logger.Error("所有处理queue的任务都意外退出了")
}

func (this *RabbitMQ) RegisterReceiver(r Receiver) {

	this.receivers = append(this.receivers, r)

}

func (this *RabbitMQ) prepareExchange() error {
	// 申明Exchange
	err := this.channel.ExchangeDeclare(
		this.exchangeName,         // exchange
		string(this.exchangeType), // type
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   // args
	)

	if nil != err {
		return err
	}

	return nil
}

// Listen 监听指定路由发来的消息
// 这里需要针对每一个接收者启动一个goroutine来执行listen
// 该方法负责从每一个接收者监听的队列中获取数据，并负责重试
func (this *RabbitMQ) listen(receiver Receiver) {
	defer this.wg.Done()

	// 这里获取每个接收者需要监听的队列和路由
	queueName := receiver.QueueName()
	routerKey := receiver.RouterKey()

	// 申明Queue
	_, err := this.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive(排他性队列)
		false,     // no-wait
		nil,       // arguments
	)
	if nil != err {
		// 当队列初始化失败的时候，需要告诉这个接收者相应的错误
		receiver.OnError(fmt.Errorf("初始化队列 %s 失败: %s", queueName, err.Error()))
	}

	// 将Queue绑定到Exchange上去
	err = this.channel.QueueBind(
		queueName,         // queue name
		routerKey,         // routing key
		this.exchangeName, // exchange
		false,             // no-wait
		nil,
	)
	if nil != err {
		receiver.OnError(fmt.Errorf("绑定队列 [%s - %s] 到交换机失败: %s", queueName, routerKey, err.Error()))
	}

	// 获取消费通道
	//	this.channel.Qos(1, 0, true) --这个导致channel被关闭// 确保rabbitmq会一个一个发消息
	msgs, err := this.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if nil != err {
		receiver.OnError(fmt.Errorf("获取队列 %s 的消费通道失败: %s", queueName, err.Error()))
	}

	forever := make(chan bool)

	// 使用callback消费数据
	for msg := range msgs {
		// 当接收者消息处理失败的时候，
		// 比如网络问题导致的数据库连接失败，redis连接失败等等这种
		// 通过重试可以成功的操作，那么这个时候是需要重试的
		// 直到数据处理成功后再返回，然后才会回复rabbitmq ack
		for !receiver.OnReceive(msg.Body) {
			this.logger.Error("receiver 数据处理失败，将要重试")
			time.Sleep(1 * time.Second)
		}

		// 确认收到本条消息, multiple必须为false
		msg.Ack(false)
	}

	<-forever
}
