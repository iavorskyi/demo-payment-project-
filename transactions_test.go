package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/go-pg/pg/v10"
)

func TestTransactions(t *testing.T) {
	handlers := Handlers()
	server := httptest.NewServer(handlers)
	defer server.Close()

	////////////////////////////////////////////// Prerequisites
	transactionCreatePath := "/transactions/create"
	db := pg.Connect(&pg.Options{
		Addr:     dbConnectionString,
		User:     dbUser,
		Password: dbPassword,
		Database: dbName,
	})
	defer db.Close()

	e := httpexpect.New(t, server.URL)

	testWalletEURO := walletModel{
		UserID:     "test",
		CurrencyID: 1,
		Balance:    1000.00,
	}
	testWalletUSD := walletModel{
		UserID:     "test",
		CurrencyID: 2,
		Balance:    1000.00,
	}

	depositUSDPayload := transactionPayload{
		UserID:          "test",
		Currency:        "USD",
		Amount:          100,
		TimePlaced:      "24-JAN-20 10:27:44",
		TransactionType: "deposit",
	}
	depositEUROPayload := transactionPayload{
		UserID:          "test",
		Currency:        "EURO",
		Amount:          100,
		TimePlaced:      "24-JAN-20 10:27:44",
		TransactionType: "deposit",
	}
	withdrawalUSDPayload := transactionPayload{
		UserID:          "test",
		Currency:        "USD",
		Amount:          100,
		TimePlaced:      "24-JAN-20 10:27:44",
		TransactionType: "withdrawal",
	}
	withdrawalEUROPayload := transactionPayload{
		UserID:          "test",
		Currency:        "EURO",
		Amount:          100,
		TimePlaced:      "24-JAN-20 10:27:44",
		TransactionType: "withdrawal",
	}
	withdrawalEUROPayloadStatus2 := transactionPayload{
		UserID:          "test",
		Currency:        "EURO",
		Amount:          100000,
		TimePlaced:      "24-JAN-20 10:27:44",
		TransactionType: "withdrawal",
	}
	withdrawalUSDPayloadStatus2 := transactionPayload{
		UserID:          "test",
		Currency:        "USD",
		Amount:          100000,
		TimePlaced:      "24-JAN-20 10:27:44",
		TransactionType: "withdrawal",
	}
	_, _ = db.Model(&testWalletEURO).Returning("*").Insert()
	_, _ = db.Model(&testWalletUSD).Returning("*").Insert()

	///////////////////////////// Positive tests

	// deposit to test EURO wallet
	expectedNewBalanceEURO := depositEUROPayload.Amount + testWalletEURO.Balance
	_ = e.POST(transactionCreatePath).
		WithJSON(depositEUROPayload).Expect().Status(http.StatusOK).Body().Equal(fmt.Sprintf("\"Success. Your current balance: %.3f\"", expectedNewBalanceEURO))

	// deposit to test USD wallet
	expectedNewBalanceUSD := depositUSDPayload.Amount + testWalletUSD.Balance
	_ = e.POST(transactionCreatePath).
		WithJSON(depositUSDPayload).Expect().Status(http.StatusOK).Body().Equal(fmt.Sprintf("\"Success. Your current balance: %.3f\"", expectedNewBalanceUSD))

	// withdraw from test EURO wallet (enough money)
	expectedNewBalanceEURO = expectedNewBalanceEURO - depositEUROPayload.Amount
	_ = e.POST(transactionCreatePath).
		WithJSON(withdrawalEUROPayload).Expect().Status(http.StatusOK).Body().Equal(fmt.Sprintf("\"Success. Your current balance: %.3f\"", expectedNewBalanceEURO))

	// withdraw from test USD wallet (enough money)
	expectedNewBalanceUSD = expectedNewBalanceUSD - depositUSDPayload.Amount
	_ = e.POST(transactionCreatePath).
		WithJSON(withdrawalUSDPayload).Expect().Status(http.StatusOK).Body().Equal(fmt.Sprintf("\"Success. Your current balance: %.3f\"", expectedNewBalanceUSD))

	// withdraw from test EURO wallet (not enough money)
	_ = e.POST(transactionCreatePath).
		WithJSON(withdrawalEUROPayloadStatus2).Expect().Status(http.StatusOK).Body().Equal("\"Failed. There is not enough of money in the wallet\"")

	// withdraw from test USD wallet (not enough money)
	_ = e.POST(transactionCreatePath).
		WithJSON(withdrawalUSDPayloadStatus2).Expect().Status(http.StatusOK).Body().Equal("\"Failed. There is not enough of money in the wallet\"")

	///////////////////////////// Clean up

	// delete test wallet
	var transactionModel transactionModel
	_, err := db.Model(&transactionModel).Where("wallet_id = ?", testWalletEURO.ID).Delete()
	if err != nil {
		log.Println("Failed to delete transaction", err)
	}
	_, err = db.Model(&transactionModel).Where("wallet_id = ?", testWalletUSD.ID).Delete()
	if err != nil {
		log.Println("Failed to delete transaction", err)
	}

	// delete all transactions created by testing
	_, err = db.Model(&testWalletEURO).Where("user_id = ?", testWalletEURO.UserID).Delete()
	if err != nil {
		log.Println("Failed to delete wallet", err)
	}
}
