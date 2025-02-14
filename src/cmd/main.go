package main

import (
	"fmt"
	"log"

	// "os"
	"net"
	// "time"
	"bufio"
	buffer "delob/internal/buffer"
	p "delob/internal/processor"
	"errors"
	"regexp"
	// "strings"
	// "log"
	// "io"
	// write "baobab/internal/write"
	// read "baobab/internal/read"
)

func main() {

	// PLAN
	// 1. integration tests
	// 2. refactor tokenizer -> token scan -> parser
	// 3. add tcp
	// 4. add SCRAM
	// 5. add pipeline
	// 6. dockerize

	// need to handle transaction -> optimistic locking?
	bufferManager, err := buffer.NewBufferManager()
	if err != nil {
		return
	}

	processor := p.NewProcessor(&bufferManager)
	errInit := processor.Initialize()
	if errInit != nil {
		fmt.Println(errInit)
		return
	}

	processor.Execute("ADD PLAYER 'Tomek';")
	processor.Execute("ADD PLAYER 'Justyna';")
	processor.Execute("SET WIN FOR 'Tomek' AND LOSE FOR 'Justyna';")
	processor.Execute("SET WIN FOR 'Tomek' AND LOSE FOR 'Justyna';")
	processor.Execute("SET WIN FOR 'Justyna' AND LOSE FOR 'Tomek';")

	result, err := processor.Execute("SELECT Players ORDER BY Elo DESC;")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

	// port, err := newPort(os.Args)
	// if err != nil {
	// 	LogError(err)
	// }

	// expr1 := "CREATE CATALOG Testowy (Id int, Name text, Value float)"
	// core.HandleCall(expr1)

	// expr2 := "INSERT INTO Testowy (Id=1, Name=Tomasz, Value=69)"
	// core.HandleCall(expr2)

	// // filemanager.Write()
	// entityId, _ := hasher.HashId("1")

	// memory_access.Read("Testowy", entityId)

	// write.WriteTwo("hello")
	// read.Read()
	// fmt.Println(read.Read())

	// startTcpServer(port)
}

type Port string

func (p Port) listenFormatted() string {
	return fmt.Sprintf(":%s", p)
}

func newPort(args []string) (Port, error) {
	if len(args) < 2 {
		return Port(""), errors.New("There are not enough arguments provided.")
	}

	if match, err := regexp.MatchString("^[1-9][0-9][0-9][0-9]$", args[1]); err != nil || (err == nil && !match) {
		return Port(""), errors.New("Given port number is incorrect (it should be 4 digit number)")
	}
	return Port(args[1]), nil
}

func startTcpServer(port Port) {
	l, err := net.Listen("tcp4", port.listenFormatted())
	if err != nil {
		LogError(err)
		return
	}

	LogInfo("Started listening on port: " + string(port))
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			LogError(err)
			return
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	LogInfo(fmt.Sprintf("Serving %s", c.RemoteAddr().String()))
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		LogInfo(netData)

		c.Write([]byte("+OK\r\n"))
	}
	c.Close()
}

const (
	InfoLevel  string = "INFO"
	ErrorLevel        = "ERROR"
)

func LogInfo(msg string) {
	log.Print(universalLogFormat(msg, InfoLevel))
}

func LogError(err error) {
	log.Print(universalLogFormat(err.Error(), ErrorLevel))
}

func universalLogFormat(msg string, level string) string {
	return fmt.Sprintf("[%s] %s",
		level,
		msg)
}
