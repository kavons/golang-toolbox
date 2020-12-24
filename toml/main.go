package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type server struct {
	IP string
	DC string
}

type ClientConfig struct {
	Name      string
	ProtoName string
	AddrList  []string
	EtcdAddrs []string
	Balancer  string
}

type ZRpcClientConfig struct {
	Clients []ClientConfig
}

type tomlConfig struct {
	Servers map[string]server
	Clients *ZRpcClientConfig
}

var (
	Conf *tomlConfig
)

func main() {
	_, err := toml.DecodeFile("./conf.toml", &Conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Conf)
}