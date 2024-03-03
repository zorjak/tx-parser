package parser

import (
	"log/slog"

	"github.com/zorjak/tx-parser/scanner"
	"github.com/zorjak/tx-parser/storage"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []storage.Transaction
}

type parser struct {
	storage storage.Storage
	scanner scanner.Scanner
	url     string
}

func New(storage storage.Storage, scanner scanner.Scanner, url string) Parser {
	return &parser{
		storage: storage,
		scanner: scanner,
		url:     url,
	}
}

func (p *parser) GetCurrentBlock() int {
	return p.storage.LastParsedBlock()
}

func (p *parser) Subscribe(address string) bool {
	return p.storage.AddAddressToObserver(address)
}

func (p *parser) GetTransactions(address string) []storage.Transaction {
	observable := p.storage.IsAddressObserved(address)
	slog.Debug("is address observable", "observable", observable)
	if !observable {
		return nil
	}

	hashes := p.storage.TransactionsOfAddress(address)
	slog.Debug("GetTransactions: number of", "transactions", len(hashes))

	var transactions []storage.Transaction
	for _, hash := range hashes {
		transaction, err := p.storage.Transaction(hash)
		if err != nil {
			return nil
		}
		transactions = append(transactions, transaction)
	}

	return transactions
}
