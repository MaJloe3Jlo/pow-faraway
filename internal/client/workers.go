package client

import (
	"context"
	"io/ioutil"
	"net"
	"time"

	"github.com/rs/zerolog"
	"gitlab.com/avpetkun/pow-wow/internal/pow"
)

func StartFetchWorkers(ctx context.Context, conf *Config, log zerolog.Logger) {
	creationPause := conf.Timeout / time.Duration(conf.FetchWorkers)
	for i := 0; i < conf.FetchWorkers; i++ {
		go runFetchWorker(ctx, conf, log, i)
		time.Sleep(creationPause)
	}
}

func runFetchWorker(ctx context.Context, conf *Config, log zerolog.Logger, workerID int) {
	log = log.With().Int("worker_id", workerID).Logger()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !fetchQuote(log, conf.ServerAddr) {
				time.Sleep(conf.Timeout)
			}
		}
	}
}

func fetchQuote(log zerolog.Logger, serverAddr string) bool {
	log.Trace().Msg("connect")
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Err(err).Msg("failed to connect")
		return false
	}

	log.Debug().Msg("process handshake")

	calcDifficulty, calcDuration, err := pow.Establish(conn)
	if err != nil {
		conn.Close()
		log.Err(err).Msg("error occured")
		return false
	}

	response, err := ioutil.ReadAll(conn)
	conn.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to read from connection")
		return false
	}
	if len(response) == 0 {
		log.Warn().Msg("empty response")
		return false
	}

	log.Info().
		Bytes("response", response).
		Int("difficulty", int(calcDifficulty)).
		Dur("proof_calc_duration", calcDuration).
		Msg("received response")
	return true
}
