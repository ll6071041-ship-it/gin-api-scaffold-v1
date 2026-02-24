package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

// 1. 初始化连接
func Init() {
	var err error
	// 这里的账号密码和端口对应你 docker-compose 里的配置
	// amqp://账号:密码@主机名:端口/虚拟主机 (这里用默认的 /)
	dsn := "amqp://admin:admin123@rabbitmq:5672/"

	Conn, err = amqp.Dial(dsn)
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}

	Channel, err = Conn.Channel()
	if err != nil {
		log.Fatalf("无法打开 Channel: %v", err)
	}

	// 声明一个队列（如果不存在会自动创建）
	_, err = Channel.QueueDeclare(
		"test_queue", // 队列名字
		true,         // durable: 是否持久化（重启不丢失）
		false,        // autoDelete: 用完是否自动删除
		false,        // exclusive: 是否排他（仅当前连接可用）
		false,        // noWait: 是否阻塞等待
		nil,          // arguments: 其他参数
	)
	if err != nil {
		log.Fatalf("声明队列失败: %v", err)
	}
	fmt.Println("🐰 RabbitMQ 初始化成功！")
}

// 2. 发送消息 (生产者)
func SendMessage(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Channel.PublishWithContext(ctx,
		"",           // exchange: 交换机（留空代表使用默认交换机）
		"test_queue", // routing key: 路由键（默认交换机模式下，这里写队列名）
		false,        // mandatory: 强制路由
		false,        // immediate: 立即投递
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message), // 把字符串转成字节发送
		})
	return err
}

// 3. 接收消息 (消费者)
func StartConsumer() {
	msgs, err := Channel.Consume(
		"test_queue", // 监听哪个队列
		"",           // consumer: 消费者名字（留空自动生成）
		true,         // auto-ack: 自动确认（收到消息就告诉MQ我搞定了）
		false,        // exclusive: 是否排他
		false,        // no-local:
		false,        // no-wait:
		nil,          // args:
	)
	if err != nil {
		log.Fatalf("注册消费者失败: %v", err)
	}

	// 开一个后台协程一直盯着这个队列
	go func() {
		for d := range msgs {
			log.Printf("📥 收到 MQ 消息: %s", d.Body)
			// 这里写你的业务逻辑，比如：去数据库扣库存！
		}
	}()
}
