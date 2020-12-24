package yaml

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"gopkg.in/yaml.v2"
)

type AMQPExchangeInfo struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type AMQPQueueInfo struct {
	Name    string `yaml:"name"`
	Durable bool   `yaml:"durable"`
}

type AMQPConfig struct {
	AMQPExchangeConfig map[string]AMQPExchangeInfo `yaml:"exchange"`
	AMQPQueueConfig    map[string]AMQPQueueInfo    `yaml:"queue"`
}

func TestPeatioAMQPConfig(t *testing.T) {
	configBytes, err := ioutil.ReadFile("./amqp.yml")
	if err != nil {
		log.Fatal(err)
	}

	config := AMQPConfig{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)
}
