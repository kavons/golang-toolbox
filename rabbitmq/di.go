package rabbitmq

import (
	"go.uber.org/dig"
)

var container = dig.New()

func BuildContainer() *dig.Container {
	// amqp connection config
	container.Provide(NewAMQPConnectionConfig)

	// amqp information config
	container.Provide(NewAMQPInfoConfig)

	// message broker
	container.Provide(NewAMQPMessageBroker)

	// server
	container.Provide(NewAMQPServer)

	return container
}

func Invoke(i interface{}) error {
	return container.Invoke(i)
}
