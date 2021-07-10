package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/siuyin/dflt"
)

func main() {
	fmt.Println("kafka consumer")

	c := consumer(cfg())

	topic := dflt.EnvString("TOPIC_PREFIX", "vwjkx5cb-") + dflt.EnvString("TOPIC", "test")
	subscribe(c, topic)

	for {
		msg := read(c)
		fmt.Printf("msg on %s: %s\n", msg.TopicPartition, string(msg.Value))
	}
}

func cfg() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"metadata.broker.list": dflt.EnvString("BROKERS", "sulky-01.srvs.cloudkafka.com:9094,sulky-02.srvs.cloudkafka.com:9094,sulky-03.srvs.cloudkafka.com:9094"),

		"group.id":             dflt.EnvString("CONS_GRP", "abc"),
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},

		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "SCRAM-SHA-256",
		"sasl.username":     dflt.EnvString("USERNAME", "vwjkx5cb"),
		"sasl.password":     dflt.EnvString("KAFKA_PASSWORD", "someS3cRet"),
	}
}

func consumer(cfg *kafka.ConfigMap) *kafka.Consumer {
	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s\n", err)
	}
	return c
}

func subscribe(c *kafka.Consumer, topic string) {
	if err := c.Subscribe(topic, nil); err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}
}

func read(c *kafka.Consumer) *kafka.Message {
	msg, err := c.ReadMessage(-1) // -1 means never timeout
	if err != nil {
		log.Printf("error: %v: %v", err, msg)
		return &kafka.Message{}
	}
	return msg
}

func unsub(c *kafka.Consumer) {
	c.Unsubscribe()
	log.Println("done")
}
