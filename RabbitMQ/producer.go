package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)


const (
	//AMQP URI
	uri          =  "amqp://cxb:dtct2017@101.37.243.50:5672/common"
	//exchangeName =  "ha-vhost-exchange-topic"
	exchangeName =  "ha-vhost-warn-topic"
	//Exchange type - direct|fanout|topic|x-custom
	exchangeType = "topic"
	bindingKey = ""
	//Durable AMQP queue name
	queueName    = "ha-vhost-dw"
)

//如果存在错误，则输出
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main(){
	bodyMsg := bodyFrom(os.Args)
	//调用发布消息函数
	publishExchange(uri, exchangeName, exchangeType,bodyMsg)
	log.Printf("published %dB OK", len(bodyMsg))
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello idoall.org ok"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}

//发布者的方法
func publishExchange(amqpURI string, exchange string, exchangeType string, body string){
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

	log.Printf("got Channel, declaring %q Exchange (%q)", exchangeType, exchange)

	//创建一个queue

	err = channel.ExchangeDeclare(
		exchange, // name
		exchangeType,
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	//创建一个queue
	_, err = channel.QueueDeclare(
		queueName, // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive 当Consumer关闭连接时，这个queue要被deleted
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue")

	log.Printf("declared queue, publishing %dB body (%q)", len(body), body)

	// Producer只能发送到exchange，它是不能直接发送到queue的。
	// 现在我们使用默认的exchange（名字是空字符）。这个默认的exchange允许我们发送给指定的queue。
	// routing_key就是指定的queue名字。
	for i:=0 ; i< 100; i++ {
		err = channel.Publish(
			exchange,     // exchange
			severityFrom(os.Args), // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing {
				Headers:         amqp.Table{},
				//DeliveryMode: amqp.Persistent,
				ContentType: "text/plain",
				ContentEncoding: "",
				//Body:        []byte(body+strconv.Itoa(i)),
				Body:[]byte("{\"MsgID\":\"20310\",\"Name\":\"zhangsan\"}"+body+strconv.Itoa(i)),
			})
		log.Printf(body+strconv.Itoa(i))
		time.Sleep(time.Microsecond*500000)
	}

	//没有写入数据，一直等待读，阻塞当前线程，目的是让线程不退出
	//forever := make(chan bool)
	//<-forever

	failOnError(err, "Failed to publish a message")
}
