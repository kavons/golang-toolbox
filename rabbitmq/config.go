package rabbitmq

import (
	"io/ioutil"
	"path/filepath"
	"runtime"

	"won-utils"

	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v2"
)

type AMQPConnectionConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func NewAMQPConnectionConfig() *AMQPConnectionConfig {
	return &AMQPConnectionConfig{
		Host:     beego.AppConfig.String("AMQPHost"),
		Port:     utils.String2Int(beego.AppConfig.String("AMQPPort")),
		User:     beego.AppConfig.String("AMQPUser"),
		Password: beego.AppConfig.String("AMQPPassword"),
	}
}

type AMQPExchangeInfo struct {
	Name       string     `yaml:"name"`
	Type       string     `yaml:"type"`
	Durable    bool       `yaml:"durable"`
	AutoDelete bool       `yaml:"auto_delete"`
	Internal   bool       `yaml:"internal"`
	NoWait     bool       `yaml:"no_wait"`
	Args       amqp.Table `yaml:"args"`
}

type AMQPQueueInfo struct {
	Name       string     `yaml:"name"`
	Durable    bool       `yaml:"durable"`
	AutoDelete bool       `yaml:"auto_delete"`
	Exclusive  bool       `yaml:"exclusive"`
	NoWait     bool       `yaml:"no_wait"`
	Args       amqp.Table `yaml:"args"`
}

type AMQPBindingInfo struct {
	QueueKey    string `yaml:"queue"`
	ExchangeKey string `yaml:"exchange"`
	CleanStart  bool   `yaml:"clean_start"`

	BindOptions    AMQPBindOptions    `yaml:"bind_options"`
	ConsumeOptions AMQPConsumeOptions `yaml:"consume_options"`
}

type AMQPChannelInfo struct {
	Prefetch int `yaml:"prefetch"`
}

type AMQPBindOptions struct {
	RoutingKey string     `yaml:"routing_key"`
	NoWait     bool       `yaml:"no_wait"`
	Args       amqp.Table `yaml:"args"`
}

type AMQPConsumeOptions struct {
	Tag       string     `yaml:"tag"`
	AutoAck   bool       `yaml:"auto_ack"`
	Exclusive bool       `yaml:"exclusive"`
	NoLocal   bool       `yaml:"no_local"`
	NoWait    bool       `yaml:"no_wait"`
	Args      amqp.Table `yaml:"args"`
}

type AMQPPublishOptions struct {
	RoutingKey string
	Mandatory  bool
	Immediate  bool
	Publishing amqp.Publishing
}

type AMQPInfoConfig struct {
	ExchangeInfoConfig map[string]AMQPExchangeInfo `yaml:"exchange"`
	QueueInfoConfig    map[string]AMQPQueueInfo    `yaml:"queue"`
	BindingInfoConfig  map[string]AMQPBindingInfo  `yaml:"binding"`
	ChannelInfoConfig  map[string]AMQPChannelInfo  `yaml:"channel"`
}

func NewAMQPInfoConfig() (*AMQPInfoConfig, error) {
	config := new(AMQPInfoConfig)

	_, file, _, _ := runtime.Caller(0)
	appPath, _ := filepath.Abs(filepath.Dir(file))
	configBytes, err := ioutil.ReadFile(appPath + "/amqp.yml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
