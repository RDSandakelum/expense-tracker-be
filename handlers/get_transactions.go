package handlers

type TransactionsListDto struct {
	ID          string `json:"id"`
	Merchant    string `json:"merchant"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	Date        string `json:"date"`
	Amount      string `json:"amount"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}
