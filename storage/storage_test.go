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

func TestAddTransactionToAddress(t *testing.T) {
    storage := New(5)
	// Test TransactionsOfAddress and AddTransactionToAddress
	transactionHash := "0x456def"
    address := "0x123abc"
	storage.AddTransactionToAddress(address, transactionHash)
	transactions := storage.TransactionsOfAddress(address)
	if len(transactions) != 1 || transactions[0] != transactionHash {
		t.Errorf("Expected transaction %s to be added to address %s transactions", transactionHash, address)
	}

	// Test AddTransaction and Transaction
	transaction := Transaction{
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

	// Test Transaction when transaction is not found
	_, err = storage.Transaction("0xnotfound")
	if err == nil {
		t.Errorf("Expected error when retrieving non-existent transaction")
	}
}