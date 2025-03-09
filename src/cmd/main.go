package main

import (
	buffer "delob/internal/buffer"
	"delob/internal/interfaces"
	p "delob/internal/processor"
	"delob/internal/utils/logger"
	"fmt"
)

func main() {
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
