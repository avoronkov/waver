package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	_ "github.com/tliron/commonlog/simple"
)

const lsName = "Waver Language Server"

var version string = "0.0.1"
var handler protocol.Handler

func main() {
	logfile, err := os.Create("waver-lsp.log")
	if err != nil {
		log.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(logfile, nil))
	slog.SetDefault(logger)

	commonlog.Configure(2, nil)

	serverImpl := NewServer()

	handler = protocol.Handler{
		Initialize:             serverImpl.Initialize,
		Shutdown:               serverImpl.Shutdown,
		TextDocumentCompletion: serverImpl.TextDocumentCompletion,
		TextDocumentDidChange:  serverImpl.TextDocumentDidChange,
	}

	server := server.NewServer(&handler, lsName, true)

	server.RunStdio()
}
