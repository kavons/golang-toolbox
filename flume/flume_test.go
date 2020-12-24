package flume_test

import (
	"log"
	"os"
	"testing"

	"fmt"
	"github.com/gemalto/flume"
)

func TestExample(t *testing.T) {
	flume.Configure(flume.Config{
		Development:   true,
		DefaultLevel:  flume.DebugLevel,
		Encoding:      "json",
		EncoderConfig: flume.NewEncoderConfig(),
	})

	logFile, err := os.OpenFile("/var/log/beego/skype_alarm.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("skype alarm file can not be opened\n")
	}
	flume.SetOut(logFile)

	logTools := flume.New("tools")

	logTools.Info("Hello Tools!")
	logTools.Info("This entry has properties", "color", "red")
	logTools.Debug("This is a debug message")
	logTools.Error("This is an error message")
	logTools.Info("This message has a multiline value", "essay", `Four score and seven years ago
our fathers brought forth on this continent, a new nation, 
conceived in Liberty, and dedicated to the proposition that all men are created equal.`)

	child := logTools.With("child", "mayuxiang")
	child.Info("success")

	logModels := flume.New("models")
	logModels.Info("Hello Models!")
	child.Info("success")
}

func TestForTest(t *testing.T) {
	set := make(map[string]interface{})

	if set["name"] == nil {
		fmt.Println("nil")
	} else {
		fmt.Println("other")
	}
}
