package main

import (
	"fmt"

	// "os"

	// "time"

	buffer "delob/internal/buffer"
	"delob/internal/interfaces"
	p "delob/internal/processor"
	"delob/internal/utils/logger"
	// "strings"
	// "log"
	// "io"
	// write "baobab/internal/write"
	// read "baobab/internal/read"
)

func main() {

	// PLAN
	// OK 1. integration tests
	// OK 2. refactor tokenizer -> token scan -> parser
	// OK 3. add tcp
	// 4. add SCRAM (user, credentials, store, handshake, session store)
	// OK 5. add pipeline
	// OK 6. dockerizecd
	// 7. backup
	// 8. move tests to higher level
	// FIX: order by elo asc; works but should throw since it needs to be Elo not elo
	// FIX: parse -> should not be able to use key more than once

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

	tcpServer := interfaces.NewTcpServer(5678)

	logger.Info("", "delob is up and running!")
	tcpServer.Start(processor.Execute)
}
