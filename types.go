package main

const (
	euro       = "EURO"
	usd        = "USD"
	deposit    = "deposit"
	withdrawal = "withdrawal"
)

type transactionPayload struct {
	UserID          string  `json:"user_id"`
	Currency        string  `json:"currency"`
	Amount          float32 `json:"amount"`
	TimePlaced      string  `json:"time_placed"`
	TransactionType string  `json:"type"`
}

type transactionModel struct {
	tableName         struct{} `pg:"transactions"`
	ID                int      `json:"id" pg:"id"`
	WalletID          int      `json:"wallet_id"`
	Amount            float32  `json:"amount"`
	TimePlaced        string   `json:"time_placed"`
	TransactionTypeID int      `json:"transaction_type_id"`
	Status            int      `json:"status"`
}

type walletModel struct {
	tableName  struct{} `pg:"wallets"`
	ID         int      `json:"id" pg:"id"`
	UserID     string   `json:"user_id"`
	CurrencyID float32  `json:"currency_id"`
	Balance    float32  `json:"balance"`
}

type currencyModel struct {
	tableName    struct{} `pg:"currency"`
	ID           int      `json:"id" pg:"id"`
	CurrencyName string   `json:"currency_name"`
}

type transactionTypeModel struct {
	tableName           struct{} `pg:"transaction_type"`
	ID                  int      `json:"id" pg:"id"`
	TransactionTypeName string   `json:"transaction_type_name"`
}
