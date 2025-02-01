package main

import (
	"flag"
	"io"
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
	flag.Parse()

	var logOutput io.Writer = os.Stderr
	if debugFile {
		logfile, err := os.Create("waver-lsp.log")
		if err != nil {
			log.Fatal(err)
		}
		logOutput = logfile
	}
	logger := slog.New(slog.NewTextHandler(logOutput, nil))
	slog.SetDefault(logger)

	commonlog.Configure(2, nil)

	serverImpl := NewServer()

	handler = protocol.Handler{
		Initialize:             serverImpl.Initialize,
		Shutdown:               serverImpl.Shutdown,
		TextDocumentCompletion: serverImpl.TextDocumentCompletion,
		TextDocumentDidOpen:    serverImpl.TextDocumentDidOpen,
		TextDocumentDidChange:  serverImpl.TextDocumentDidChange,
		TextDocumentHover:      serverImpl.TextDocumentHover,
		TextDocumentDefinition: serverImpl.TextDocumentDefinition,
	}

	server := server.NewServer(&handler, lsName, true)

	server.RunStdio()
}
