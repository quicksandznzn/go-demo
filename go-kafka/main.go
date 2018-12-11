// Created by quicksandzn@gmail.com on 2018/7/25
package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"
	"sync"
)

type Config struct {
	Kafka struct {
		Address []string `yaml:"addrs"`
		Topic   string   `yaml:"topic"`
	}
}

const ConfigPath = "./config/conf.yml"

func main() {
	config := Config{}
	buffer, err := ioutil.ReadFile(ConfigPath)
	failOnError(err, "read config error")
	err = yaml.Unmarshal(buffer, &config)
	failOnError(err, "yml convert error")

	//producer(config)
	consumer(config)

}

func producer(config Config) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	//kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewSyncProducer(config.Kafka.Address, kafkaConfig)

	failOnError(err, "kafka producer error")
	defer producer.Close()

	value := "hello kafka"
	msg := &sarama.ProducerMessage{
		Topic: config.Kafka.Topic,
		Value: sarama.ByteEncoder(value),
		Key:   sarama.StringEncoder("key"),
	}

	for {
		// 生产消息
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			failOnError(err, "Send Message Fail")
		}
		fmt.Printf("Partion = %d, offset = %d\n", partition, offset)
	}
}

var (
	wg sync.WaitGroup
)

func consumer(config Config) {
	consumer, err := sarama.NewConsumer(config.Kafka.Address, nil)

	if err != nil {
		panic(err)
	}

	partitionList, err := consumer.Partitions(config.Kafka.Topic)

	if err != nil {
		panic(err)
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(config.Kafka.Topic, int32(partition), sarama.OffsetOldest)

		failOnError(err,"kafka consumer error")
		defer pc.AsyncClose()

		wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
		wg.Wait()
		consumer.Close()
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(err)
	}
}
