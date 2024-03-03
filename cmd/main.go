package main

import (
	"flag"

	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zorjak/tx-parser/cmd/handlers"
	"github.com/zorjak/tx-parser/parser"
	"github.com/zorjak/tx-parser/scanner"
	"github.com/zorjak/tx-parser/storage"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		slog.Error("error occurred", "err", err)
	}
}

func run() error {
	startingBlock := flag.Int("starting_block", 0, "starting block number")
	debug := flag.Bool("debug", false, "Set logs to debug level")
	url := flag.String("url", "https://cloudflare-eth.com", "url of the ethereum node (rpc)")
	flag.Parse()

	setLogger(*debug)

	slog.Info("startingBlock", "startingBlock", *startingBlock)
	//storage := storage.New(19350586)
	storage := storage.New(*startingBlock)
	scanner := scanner.New(*url, storage)

	parser := parser.New(storage)

	router := mux.NewRouter()
	router.HandleFunc("/api/current-block", handlers.CurrentBlockHandler(parser)).Methods("GET")
	router.HandleFunc("/api/subscribe", handlers.SubscribeHandler(parser)).Methods("POST")
	router.HandleFunc("/api/transactions/{address}", handlers.GetTransactionsHandler(parser)).Methods("GET")

	group := new(errgroup.Group)
	group.Go(func() error {
		slog.Info("Server started at port", "port", 8080)
		return http.ListenAndServe(":8080", router)
	})

	group.Go(func() error {
		return scanner.ScanBlocks()
	})

	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func setLogger(debug bool) {
	// Default slog.Level is Info (0)
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}).WithAttrs([]slog.Attr{}))
	slog.SetDefault(logger)
}
