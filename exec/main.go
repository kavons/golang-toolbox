package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

func main() {
	// echo welcome
	var echoOut bytes.Buffer
	echoCmd := exec.Command("echo", "welcome")
	echoCmd.Stdout = &echoOut
	echoCmd.Run()
	fmt.Print(echoOut.String())

	// echo "select * from users" | tr "a-z" "A-Z"
	var trOut bytes.Buffer
	trCmd := exec.Command("tr", "a-z", "A-Z")
	trCmd.Stdin = strings.NewReader("select * from users")
	trCmd.Stdout = &trOut
	trCmd.Run()
	fmt.Println(trOut.String())

	// echo -n '{"Name": "Bob", "Age": 32}'
	cmd := exec.Command("echo", "-n", `{"Name": "Bob", "Age": 32}`)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	var person struct{
		Name string
		Age int
	}

	if err = json.NewDecoder(stdout).Decode(&person); err != nil {
		log.Fatal(err)
	}

	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(person)

	// dockerLogin
	dockerLogin("daocloud.io", "kavonm", "mjsc10230713")
}

func dockerLogin(registry, username, token string) {
	params := []string {
		"-c",
		fmt.Sprintf("docker login %s -u %s --password-stdin", registry, username),
	}
	cmd := exec.Command("bash", params...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, token+"\n")
	}()

	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	return
}