package scanner

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zorjak/tx-parser/storage"
)

type Scanner interface {
	ScanBlocks() error
	ScanOneBlock(blockNumber int) error
}

type scanner struct {
	url     string
	storage storage.Storage
}

func New(url string, storage storage.Storage) Scanner {
	return &scanner{
		url:     url,
		storage: storage,
	}
}

type BlockResult struct {
	Transactions []storage.Transaction `json:"transactions"`
}

type BlockResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  BlockResult `json:"result"`
}

func (s *scanner) ScanOneBlock(blockNumber int) error {

	hexBlockNumber := strconv.FormatInt(int64(blockNumber), 16)

	payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[ \"0x" + hexBlockNumber + "\",true]}")
	slog.Debug("payload", "payload", payload)
	req, err := http.NewRequest("POST", s.url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data BlockResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	slog.Debug("number of transactions in block", "blockNumber", blockNumber, "transactions", len(data.Result.Transactions))
	for _, tx := range data.Result.Transactions {
		s.storage.AddTransaction(tx)
	}

	return nil
}


func (s *scanner) ScanBlocks() error {
	for {
		lastBlock := s.storage.LastParsedBlock()
		lastBlock++
		err := s.ScanOneBlock(lastBlock)
		if err != nil {
			slog.Error("error occurred", "err", err)
			time.Sleep(2 * time.Second)
			continue
		}
		s.storage.SetLastParsedBlock(lastBlock)
	}
}