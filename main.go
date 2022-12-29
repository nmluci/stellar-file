package main

import (
	"github.com/nmluci/stellar-file/cmd/webservice"
	"github.com/nmluci/stellar-file/internal/component"
	"github.com/nmluci/stellar-file/internal/config"
)

func main() {
	config.Init()
	conf := config.Get()

	logger := component.NewLogger(component.NewLoggerParams{
		ServiceName: conf.ServiceName,
		PrettyPrint: true,
	})

	webservice.Start(conf, logger)
}
