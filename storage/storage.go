package storage

import (
	"errors"
	"log/slog"
	"strings"
)

type AccessList struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

type Transaction struct {
	Type                 string       `json:"type"`
	BlockHash            string       `json:"blockHash"`
	BlockNumber          string       `json:"blockNumber"`
	From                 string       `json:"from"`
	Gas                  string       `json:"gas"`
	Hash                 string       `json:"hash"`
	Input                string       `json:"input"`
	Nonce                string       `json:"nonce"`
	To                   string       `json:"to"`
	TransactionIndex     string       `json:"transactionIndex"`
	Value                string       `json:"value"`
	V                    string       `json:"v"`
	R                    string       `json:"r"`
	S                    string       `json:"s"`
	GasPrice             string       `json:"gasPrice"`
	MaxFeePerGas         string       `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string       `json:"maxPriorityFeePerGas"`
	ChainId              string       `json:"chainId"`
	AccessList           []AccessList `json:"accessList"`
	YParity              string       `json:"yParity"`
}

type Storage interface {
	LastParsedBlock() int
	SetLastParsedBlock(block int)
	AddAddressToObserver(address string) bool
	IsAddressObserved(address string) bool
	AddTransactionToAddress(address string, transaction string) error
	TransactionsOfAddress(address string) []string
	AddTransaction(transaction Transaction)
	Transaction(hash string) (Transaction, error)
}

type storage struct {
	currentBlock          int
	observedAddresses     map[string]bool
	transactionsOfAddress map[string][]string
	transactions          map[string]Transaction
}

func New(startingBlock int) Storage {
	return &storage{
		currentBlock:          startingBlock,
		observedAddresses:     make(map[string]bool),
		transactionsOfAddress: make(map[string][]string),
		transactions:          make(map[string]Transaction),
	}
}

// SetLastParsedBlock sets the last parsed block
func (s *storage) SetLastParsedBlock(block int) {
	s.currentBlock = block
}

// LastParsedBlock returns the last parsed block
func (s *storage) LastParsedBlock() int {
	return s.currentBlock
}

// AddAddressToObserver adds an address to the observer list
func (s *storage) AddAddressToObserver(address string) bool {
	slog.Debug("adding address to observer", "address", address)

	s.observedAddresses[strings.ToLower(address)] = true
	return true
}

// IsAddressObserved checks if an address is observed
func (s *storage) IsAddressObserved(address string) bool {
	_, ok := s.observedAddresses[strings.ToLower(address)]
	return ok
}

// TransactionsOfAddress returns the transactions hash
func (s *storage) TransactionsOfAddress(address string) []string {
	transactions, ok := s.transactionsOfAddress[strings.ToLower(address)]
	slog.Debug("TransactionsOfAddress", "address", address, "transactionsNumber", len(transactions))
	if !ok {
		return nil
	}
	return transactions
}

func (s *storage) AddTransactionToAddress(address string, transaction string) error {
	slog.Debug("adding transaction to address", "address", address, "transaction", transaction)
	s.transactionsOfAddress[strings.ToLower(address)] = append(s.transactionsOfAddress[strings.ToLower(address)], transaction)
	return nil
}

func (s *storage) AddTransaction(transaction Transaction) {
	slog.Debug("adding transaction", "transaction", transaction)
	s.transactions[strings.ToLower(transaction.Hash)] = transaction
}

func (s *storage) Transaction(hash string) (Transaction, error) {
	transaction, ok := s.transactions[strings.ToLower(hash)]
	if !ok {
		return Transaction{}, errors.New("transaction not found")
	}
	return transaction, nil
}
