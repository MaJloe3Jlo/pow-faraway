package main

import (
	_ "embed"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"gitlab.com/avpetkun/pow-wow/internal/server"
)

//go:embed quotes.json
var jsonQuotes []byte

func main() {
	zerolog.DurationFieldUnit = time.Millisecond

	conf := server.NewConfig()

	log := zerolog.New(&zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().Timestamp().
		Logger()

	log.Debug().
		Str("listen_addr", conf.ListenAddr).
		Int("proof_token_size", conf.ProofTokenSize).
		Int("difficulty", int(conf.Difficulty)).
		Msg("server started")

	book, err := server.NewBook(jsonQuotes)
	check(log, err)

	server, err := server.StartServer(conf, log, book.ServeRequest)
	check(log, err)
	defer server.Close()

	waitForExit()
}

func waitForExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func check(log zerolog.Logger, err error) {
	if err != nil {
		log.Fatal().Err(err).CallerSkipFrame(1).Msg("start failed")
		panic(err)
	}
}
