package iyzico

type CheckoutFormInitRequest struct {
	Locale         Locale         `json:"locale"`
	ConversationID string         `json:"conversationId"`
	Price          string         `json:"price"`
	BasketID       string         `json:"basketId"`
	Buyer          Buyer          `json:"buyer"`
	BillingAddress BillingAddress `json:"billingAddress"`
	BasketItems    []BasketItems  `json:"basketItems"`
	CallbackURL    string         `json:"callbackUrl"`
	Currency       Currency       `json:"currency"`
	// EnabledInstallments uint           `json:"enabledinstallments"`
	PaidPrice string `json:"paidPrice"` // price + other value (discount, tax, ...)
}

type CheckoutFormSearchRequest struct {
	Locale         Locale `json:"locale"`
	ConversationID string `json:"conversationId"`
	Token          string `json:"token"`
}

type CallbackRequest struct {
	Token string `form:"token" json:"token"`
}

type RequestInfo struct {
	ApiKey    string `json:"-"`
	SecretKey string `json:"-"`
	URL       string `json:"-"`
}
