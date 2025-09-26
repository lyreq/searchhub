package iyzico

type CheckoutFormInitResponse struct {
	ResponseStatus
	CheckoutFormContent string `json:"checkoutFormContent"`
	PaymentPageURL      string `json:"paymentPageUrl"`
	Token               string `json:"token"`
	TokenExpireTime     int    `json:"tokenExpireTime"`
	Locale              Locale `json:"locale"`
	ConversationID      string `json:"conversationId"`
}

type CheckoutFormSearchResponse struct {
	ResponseStatus
	Token                        string             `json:"token"`
	CallbackURL                  string             `json:"callbackUrl"`
	PaymentStatus                PaymentStatus      `json:"paymentStatus"`
	Locale                       Locale             `json:"locale"`
	ConversationID               string             `json:"conversationId"`
	PaymentID                    string             `json:"paymentId"`
	Price                        float32            `json:"price"`
	PaidPrice                    float32            `json:"paidPrice"` // price + other value (discount, tax, ...)
	Currency                     Currency           `json:"currency"`
	Installment                  uint8              `json:"installment"`
	BasketID                     string             `json:"basketId"`
	BinNumber                    string             `json:"binNumber"`
	LastFourDigits               string             `json:"lastFourDigits"`
	CardAssociation              CardAssociation    `json:"cardAssociation"`
	CardFamily                   string             `json:"cardFamily"`
	CardType                     CardType           `json:"cardType"`
	FraudStatus                  FraudStatus        `json:"fraudStatus"`
	IyziCommissionFee            float32            `json:"iyziCommissionFee"`
	IyziCommissionRateAmount     float32            `json:"iyziCommissionRateAmount"`
	MerchantCommissionRate       float32            `json:"merchantCommissionRate"`
	MerchantCommissionRateAmount float32            `json:"merchantCommissionRateAmount"`
	ItemTransactions             []ItemTransactions `json:"itemTransactions"`
}

type ResponseStatus struct {
	Status       Status `json:"status"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	ErrorGroup   string `json:"errorGroup"`
	SystemTime   int64  `json:"systemTime"`
}

func (f *ResponseStatus) IsSuccess() bool {
	return f.Status == SUCCESS
}

func (f *CheckoutFormSearchResponse) PaymentIsSuccess() bool {
	return f.PaymentStatus == PSUCCESS
}
