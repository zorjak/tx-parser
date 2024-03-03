package storage

import (
	"reflect"
	"testing"
)

func TestLastParsedBlock(t *testing.T) {
	startingBlock := 1000
	storage := New(startingBlock)

	// Test LastParsedBlock
	if storage.LastParsedBlock() != startingBlock {
		t.Errorf("Expected LastParsedBlock to return %d, got %d", startingBlock, storage.LastParsedBlock())
	}

	// Test SetLastParsedBlock
	newBlock := 2000
	storage.SetLastParsedBlock(newBlock)
	if storage.LastParsedBlock() != newBlock {
		t.Errorf("Expected LastParsedBlock to return %d after SetLastParsedBlock, got %d", newBlock, storage.LastParsedBlock())
	}
}

func TestObservableAddress(t *testing.T) {
    storage := New(5)
	// Test AddAddressToObserver and IsAddressObserved
	address := "0x123abc"
	storage.AddAddressToObserver(address)
	if !storage.IsAddressObserved(address) {
		t.Errorf("Expected address %s to be observed", address)
	}
}

func TestAddTransaction(t *testing.T) {
    storage := New(5)

	transaction := Transaction{
		To: "0x123abc",
		From: "0x456def",
		Hash: "0x789ghi",
	}
	storage.AddTransaction(transaction)
	transactionRetrieved, err := storage.Transaction(transaction.Hash)
	if err != nil {
		t.Errorf("Error retrieving transaction: %v", err)
	}
	if !reflect.DeepEqual(transactionRetrieved, transaction) {
		t.Errorf("Expected retrieved transaction to match added transaction")
	}

	transactionsOfAddress := storage.TransactionsOfAddress(transaction.To)
	if len(transactionsOfAddress) != 1 {
		t.Errorf("Expected 1 transaction for address, got %d", len(transactionsOfAddress))
	}
}