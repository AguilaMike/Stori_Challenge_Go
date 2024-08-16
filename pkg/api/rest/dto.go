package rest

type AccountDTO struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type TransactionDTO struct {
	AccounID string                            `json:"account_id"`
	Summary  TransactionSummaryDTO             `json:"summary"`
	Monthly  map[string]*TransactionMonthlyDTO `json:"monthly"`
}

type TransactionSummaryDTO struct {
	AverageCredit float64 `json:"average_credit"`
	AverageDebit  float64 `json:"average_debit"`
	CreditCount   int     `json:"credit_count"`
	DebitCount    int     `json:"debit_count"`
	TotalBalance  float64 `json:"total_balance"`
	TotalCount    int     `json:"total_count"`
	TotalCredit   float64 `json:"total_credit"`
	TotalDebit    float64 `json:"total_debit"`
}

type TransactionMonthlyDTO struct {
	Year          int                    `json:"year"`
	Month         int                    `json:"month"`
	Total         int                    `json:"total_transactions"`
	Balance       float64                `json:"balance"`
	AverageCredit float64                `json:"average_credit"`
	AverageDebit  float64                `json:"average_debit"`
	Transactions  []TransactionDetailDTO `json:"transactions"`
}

type TransactionDetailDTO struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
	InputDate string  `json:"input_date"`
}
