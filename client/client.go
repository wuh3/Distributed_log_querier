package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	//"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
	//"strings"
)

type Input []Config

type Config struct {
	SSHAdd    string `json:"sshAdd"`
	Port      string `json:"port"`
	Name      string `json:"name"`
	InputPath string `json:"inputPath"`
}

type receiveinfo struct {
	LogLines  string
	LineCount int
	Succ      bool
}

type Message struct {
	Pattern  string
	Filename string
}

var USER_DIR string = "/home/haozhew3"
var REMOTE_REPO_DIR string = "/cs425/mp1"
var input Input
var wg sync.WaitGroup

func loadConfig() Input {
	// Load config properties from JSON
	content, err := ioutil.ReadFile(USER_DIR + REMOTE_REPO_DIR + "/config.json")
	if err != nil {
		log.Fatal("Config File Not Found: ", err)
	}

	err = json.Unmarshal(content, &input)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return input
}

func communicateAddress(address string, pattern string, inputFile string, grepinfo chan receiveinfo) {
	defer wg.Done()

	connCh := make(chan *rpc.Client)
	errCh := make(chan error)

	// Start a goroutine to establish the connection
	go func() {
		conn, err := rpc.Dial("tcp", address)
		if err != nil {
			errCh <- err
			return
		}
		connCh <- conn
	}()

	select {
	case conn := <-connCh:
		defer conn.Close()
		var msg = Message{pattern, inputFile}
		fmt.Println("Input path: ", msg.Filename)
		var response string
		call := conn.Go("GrepService.HandleGrep", msg, &response, nil)

		select {
		case <-time.After(10 * time.Second): // Adjust the timeout duration as needed
			// Handle the timeout for the RPC call
			fmt.Println("RPC call timed out")
			grepinfo <- receiveinfo{"", 0, false}
		case <-call.Done:
			// The RPC call completed, check for errors
			if call.Error != nil {
				fmt.Println("Error:", call.Error)
				grepinfo <- receiveinfo{"", 0, false}
				return
			}

			if response == "" {
				fmt.Printf("no info matched")
				grepinfo <- receiveinfo{response, 0, true}
			} else {
				lineCount, err := strconv.Atoi(response)
				if err != nil {
					fmt.Println("Error: ", err)
				}
				grepinfo <- receiveinfo{response, lineCount, true}
			}
		}
	case err := <-errCh:
		fmt.Printf("Error connecting to %s: %v\n", address, err)
		grepinfo <- receiveinfo{"", 0, false}
	}
}

func writeFile(filename string, logline string, matchNum *int, successNum *int) {
	var err error
	var file *os.File
	_, err = os.Stat(filename)
	if err == nil { // if file exists
		err := os.Remove(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	file, err = os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	_, err = io.WriteString(file, logline)
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(file, "Matched Number: "+strconv.Itoa(*matchNum)+"\n")
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(file, "Number of successful log queries: "+strconv.Itoa(*successNum)+"\n")
	if err != nil {
		panic(err)
	}
	fmt.Println("Output result to " + filename)
}

func main() {
	defer wg.Wait()
	// Calculate run time
	t_init := time.Now()

	// Get grep command
	if len(os.Args) != 2 {
		fmt.Println("Format error! Should be: ./client <pattern>")
		return
	}

	var pattern string = os.Args[1]

	// Connect to servers
	input := loadConfig()

	addresses := []string{}
	for _, config := range input {
		var address string = config.SSHAdd + ":" + config.Port
		addresses = append(addresses, address)
	}
	grepChannel := make(chan receiveinfo)
	for i, address := range addresses {
		wg.Add(1)
		file := input[0].InputPath + "vm" + strconv.Itoa(i+1) + ".log"
		go communicateAddress(address, pattern, file, grepChannel)
	}
	var totalLineCount = 0
	var numOfSucc = 0
	var builder strings.Builder
	for i := 0; i < len(addresses); i++ {
		reclog := <-grepChannel

		if reclog.Succ {
			numOfSucc++
			totalLineCount += reclog.LineCount
		}
		fmt.Println("Current vm: ", i+1)
		fmt.Println("Current line count: ", reclog.LineCount)
		totalLineCount += reclog.LineCount
		builder.WriteString(reclog.LogLines)
		builder.WriteString("\n")
	}

	totalTime := time.Since(t_init)
	fmt.Println("Total line count: ", totalLineCount)
	fmt.Println("Total time for the grep function:", totalTime)
	fmt.Println("Total successful executions: ", numOfSucc)
	writeFile("output.log", builder.String(), &totalLineCount, &numOfSucc)
}
