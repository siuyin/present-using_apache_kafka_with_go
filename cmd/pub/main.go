package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/siuyin/dflt"
)

func main() {
	fmt.Println("kafka producer")

	p := prod(cfg())             // p is a kafka producer
	ch := make(chan kafka.Event) // message delivery channel

	topic := dflt.EnvString("TOPIC_PREFIX", "vwjkx5cb-") + dflt.EnvString("TOPIC", "test")

	publishTime(p, topic, ch)
	printDelvMsg(ch)

	select {} // wait forever
}

func cfg() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"metadata.broker.list": dflt.EnvString("BROKERS", "sulky-01.srvs.cloudkafka.com:9094,sulky-02.srvs.cloudkafka.com:9094,sulky-03.srvs.cloudkafka.com:9094"),

		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "SCRAM-SHA-256",
		"sasl.username":     dflt.EnvString("USERNAME", "vwjkx5cb"),
		"sasl.password":     dflt.EnvString("KAFKA_PASSWORD", "someS3cRet"),
	}
}

func prod(cfg *kafka.ConfigMap) *kafka.Producer {
	p, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return p
}

func publishTime(p *kafka.Producer, topic string, ch chan kafka.Event) {
	go func() {
		for {
			err := p.Produce(msg(topic), ch)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(100 * time.Millisecond)
		}

	}()

}

func msg(topic string) *kafka.Message {
	value := time.Now().Format("15:04:05.000")
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}
}

func printDelvMsg(ch chan kafka.Event) {
	go func() {
		for {
			e := <-ch
			m := e.(*kafka.Message)

			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				return
			}

			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
	}()
}
