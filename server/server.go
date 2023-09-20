package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os/exec"
	"strings"
)

type Input []Config

type Config struct {
	SSHAdd    string `json:"sshAdd"`
	Port      string `json:"port"`
	Name      string `json:"name"`
	InputPath string `json:"inputPath"`
}

type GrepService struct{}

type Message struct {
	Pattern  string
	Filename string
}

var USER_DIR string = "/home/haozhew3"
var REMOTE_REPO_DIR string = "/cs425/mp1"
var input Input

func loadConfig() Input {
	// Load config properties from JSON
	content, err := ioutil.ReadFile(USER_DIR + REMOTE_REPO_DIR + "/config.json")
	if err != nil {
		log.Fatal("File Not Found: ", err)
	}

	err = json.Unmarshal(content, &input)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	for _, c := range input {
		fmt.Println("addr: " + c.SSHAdd)
		fmt.Println("port: " + c.Port)
		fmt.Println("vm: " + c.Name)
		fmt.Println("path: " + c.InputPath)
	}
	return input
}

func (serv *GrepService) HandleGrep(msg Message, response *string) error {
	// The receiver method to execute Grep command
	command := "grep"
	path := msg.Filename
	fmt.Println("grep command: " + command + msg.Pattern + path)

	cmd := exec.Command("grep", "-Ec", msg.Pattern, path)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Grep Error:", err)
		return err
	}

	result := strings.TrimSpace(string(output))
	*response = result

	return nil
}

func main() {
	input := loadConfig()
	// Register rpc connection
	grep := new(GrepService)
	err := rpc.RegisterName("GrepService", grep)
	if err != nil {
		log.Fatal("Error: rpc registration failed!", err)
	}

	// Start listening
	var servAddr string = ":" + input[0].Port
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		log.Fatal("Error: Incorrect TCP address!", err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	fmt.Println("Server starts listening...")
	if err != nil {
		log.Fatal("Error: RPC Registration Failed!", err)
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error: connection failed!", err)
			continue
		}
		go rpc.ServeConn(connection)
	}
}
