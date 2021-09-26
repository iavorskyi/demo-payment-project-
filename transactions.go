package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
)

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var payload transactionPayload
	var status int
	var responseMsg string

	db := pg.Connect(&pg.Options{
		Addr:     dbConnectionString,
		User:     dbUser,
		Password: dbPassword,
		Database: dbName,
	})
	defer db.Close()

	// get the transaction from the body and unmarshal it in the structure
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read JSON: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(requestBody, &payload)
	if err != nil {
		log.Printf("Error unmarshaling: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get currency id from db by name
	var currency currencyModel
	err = db.Model(&currency).Where("currency_name = ?", payload.Currency).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to get currency: %s", err)
		return
	}

	// get transaction type id from db by name
	var transactionType transactionTypeModel
	err = db.Model(&transactionType).Where("transaction_type_name = ?", payload.TransactionType).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to get transaction_type: %s", err)
		return
	}
	// get ballance from db (use cerrency from payload)
	var wallet walletModel
	err = db.Model(&wallet).Where("user_id = ?", payload.UserID).Where("currency_id = ?", currency.ID).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to get wallet: %s", err)
		return
	}

	// if this is withdrawal, check wallet possibility to commit this transaction
	if payload.TransactionType == withdrawal {
		// check if ballance >= amount from payload; if true, status 1 (comitted), false - 2 (failed)
		if wallet.Balance >= payload.Amount {
			status = 1
			wallet.Balance = wallet.Balance - payload.Amount

			_, err = db.Model(&wallet).WherePK().Update()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("Failed to update wallet balance: %s", err)
				return
			}

			// write response message
			responseMsg = fmt.Sprintf("Success. Your current balance: %.2f", wallet.Balance)

		} else {
			status = 2
			// write response message
			responseMsg = "Failed. There is not enough of money in the wallet"
		}

	}
	if payload.TransactionType == deposit {
		wallet.Balance = wallet.Balance + payload.Amount
		_, err = db.Model(&wallet).Column("balance").WherePK().Update()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to update wallet balance: %s", err)
			return
		}
		status = 1
		// write response message
		responseMsg = fmt.Sprintf("Success. Your current balance: %.2f", wallet.Balance)
	}

	// create a transaction struct to insert in db
	transaction := transactionModel{
		WalletID:          wallet.ID,
		Amount:            payload.Amount,
		TimePlaced:        payload.TimePlaced,
		TransactionTypeID: transactionType.ID,
		Status:            status,
	}

	_, err = db.Model(&transaction).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to update wallet balance: %s", err)
		return
	}

	// send the response
	err = ren.JSON(w, http.StatusOK, responseMsg)
	if err != nil {
		log.Printf("JSON rendering failed: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func getBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	db := pg.Connect(&pg.Options{
		Addr:     dbConnectionString,
		User:     dbUser,
		Password: dbPassword,
		Database: dbName,
	})
	defer db.Close()

	// get all of the user wallets
	var wallets []walletModel
	err := db.Model(&wallets).Where("user_id = ?", userID).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to get wallet: %s", err)
		return
	}

	// make and fill in a result map for response
	responseMap := make(map[string]float32)
	for _, wallet := range wallets {
		var currency currencyModel

		err = db.Model(&currency).Where("id = ?", wallet.CurrencyID).Select()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to get currency: %s", err)
			return
		}
		responseMap[currency.CurrencyName] = wallet.Balance
	}
	// send the response
	err = ren.JSON(w, http.StatusOK, responseMap)
	if err != nil {
		log.Printf("JSON rendering failed: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
