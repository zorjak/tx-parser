package parser_test

import (
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/zorjak/tx-parser/parser"
	mockStorage "github.com/zorjak/tx-parser/storage/mock"
    storage "github.com/zorjak/tx-parser/storage"
)

func TestGetTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockStorage.NewMockStorage(ctrl)
	mockParser := parser.New(mockStorage)

	// Mock behavior for IsAddressObserved
	mockStorage.EXPECT().IsAddressObserved(gomock.Any()).Return(true).Times(1)

	// Mock behavior for TransactionsOfAddress
	mockHashes := []string{"hash1", "hash2"}
	mockStorage.EXPECT().TransactionsOfAddress(gomock.Any()).Return(mockHashes).Times(1)

	// Mock behavior for Transaction
	mockStorage.EXPECT().Transaction("hash1").Return(storage.Transaction{Hash: "hash1"}, nil).Times(1)
	mockStorage.EXPECT().Transaction("hash2").Return(storage.Transaction{Hash: "hash2"}, nil).Times(1)

	transactions := mockParser.GetTransactions("address")

	if len(transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(transactions))
	}

	if transactions[0].Hash != "hash1" {
		t.Errorf("Unexpected transaction: %v", transactions[0])
	}

	if transactions[1].Hash != "hash2" {
		t.Errorf("Unexpected transaction: %v", transactions[1])
	}
}