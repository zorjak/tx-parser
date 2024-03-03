package scanner

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	storage "github.com/zorjak/tx-parser/storage"
	mockStorage "github.com/zorjak/tx-parser/storage/mock"
	"go.uber.org/mock/gomock"
)

func TestScanOneBlock(t *testing.T) {
	// Prepare a mock HTTP server to return the desired response.
	expectedBlockNumber := 123
	expectedResult := BlockResult{
		Transactions: []storage.Transaction{
			{To: "toAddress1", From: "fromAddress1", Hash: "hash1"},
			{To: "toAddress2", From: "fromAddress2", Hash: "hash2"},
		},
	}
	mockResponse := BlockResponse{
		JsonRpc: "2.0",
		Result:  expectedResult,
	}
	mockResponseBody, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request payload is as expected
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		expectedPayload := `{"id":1,"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":[ "0x7b",true]}`
		if string(body) != expectedPayload {
			t.Errorf("Unexpected request payload. Got: %s, Want: %s", body, expectedPayload)
		}

		// Respond with the mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBody)
	}))
	defer server.Close()

	store := mockStorage.NewMockStorage(gomock.NewController(t))
	// Create a new scanner with the mock storage and mock server URL
	s := &scanner{
		url:     server.URL,
		storage: store,
	}

	store.EXPECT().AddTransaction(storage.Transaction{To: "toAddress1", From: "fromAddress1", Hash: "hash1"}).Return()
	store.EXPECT().AddTransaction(storage.Transaction{To: "toAddress2", From: "fromAddress2", Hash: "hash2"}).Return()

	// Call the method to be tested
	err := s.ScanOneBlock(expectedBlockNumber)
	if err != nil {
		t.Fatalf("ScanOneBlock returned an error: %v", err)
	}
}
