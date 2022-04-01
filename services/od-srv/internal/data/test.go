package data

import (
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://root:zxcvbnm@172.19.0.4:5672")
	if err != nil {
		panic(err)
	}
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	err = channel.Publish("dead_exchange", "dead", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("774f3cff-8b46-4f50-b0c5-3439a4380c45"),
	})
	if err != nil {
		fmt.Println(err)
	}
}
