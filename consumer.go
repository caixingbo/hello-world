package main

import (
	"fmt"
	"log"
	"bytes"
	"time"
	"github.com/streadway/amqp"
	"os"
)

const (
	//AMQP URI
	uri          =  "amqp://cxb:dtct2017@101.37.243.50:5672/common"
	exchangeName =  "ha-vhost-exchange-topic"
	//Exchange type - direct|fanout|topic|x-custom
	exchangeType = "topic"
	bindingKey = ""
	//Durable AMQP queue name
	queueName    = "ha-vhost-queue"//"test-idoall-info"//
)

//如果存在错误，则输出

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main(){
	//调用消息接收者
	consumerExchange(uri, exchangeName,exchangeType,queueName,bindingKey)
}

//接收者方法
func consumerExchange(amqpURI string, exchange string, exchangeType string, queue string, key string){
	//建立连接
	log.Printf("dialing %q", amqpURI)
	connection, err := amqp.Dial(amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	//创建一个Channel
	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	log.Printf("got queue, declaring %q", queue)

	//创建一个queue
	err = channel.ExchangeDeclare(
		exchange, // name
		exchangeType,
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Exchange Declare:")

	//创建一个queue
	q, err := channel.QueueDeclare(
		queueName, // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive 当Consumer关闭连接时，这个queue要被deleted
		false,   // no-wait
		nil,     // arguments
	)


	failOnError(err, "Failed to declare a queue")

	//每次只取一条消息
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")


	//绑定到exchange
	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, exchange, s)
		//绑定到exchange
		err = channel.QueueBind(
			q.Name, // name of the queue
			s,        // bindingKey
			exchange,   // sourceExchange
			false,      // noWait
			nil,        // arguments
		);
		failOnError(err, "Failed to bind a queue")
	}


	log.Printf("Queue bound to Exchange, starting Consume")
	//订阅消息
	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//创建一个channel
	forever := make(chan bool)

	//调用gorountine
	go func() {
		num := 0
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			//dot_count := bytes.Count(d.Body, []byte("."))
			//t := time.Duration(dot_count)
			//time.Sleep( t * time.Second)
			//log.Printf("Done",num)
			//if num%1000 == 0 {
			//	log.Printf("Done",num)
			//}
			num++
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")


	conErr := make(chan *amqp.Error)
	connection.NotifyClose(conErr)
	go func() {
		select {
		case err:=<-conErr :
			println("reconnect %v",err)
			channel.Recover(true)
		}

	}()

	//没有写入数据，一直等待读，阻塞当前线程，目的是让线程不退出
	<-forever

	//test push
}
