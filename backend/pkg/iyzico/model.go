package iyzico

type ItemType string

const (
	PHYSICAL ItemType = "PHYSICAL"
	VIRTUAL  ItemType = "VIRTUAL"
)

type Locale string

const (
	TR Locale = "tr"
	EN Locale = "en"
)

type Status string

const (
	SUCCESS Status = "success"
	FAILURE Status = "failure"
)

type PaymentStatus string

const (
	PSUCCESS         PaymentStatus = "SUCCESS"
	PFAILURE         PaymentStatus = "FAILURE"
	INIT_THREEDS     PaymentStatus = "INIT_THREEDS"
	CALLBACK_THREEDS PaymentStatus = "CALLBACK_THREEDS"
	BKM_POS_SELECTED PaymentStatus = "BKM_POS_SELECTED"
	CALLBACK_PECCO   PaymentStatus = "CALLBACK_PECCO"
)

type CardAssociation string

const (
	VISA             CardAssociation = "VISA"
	MASTER_CARD      CardAssociation = "MASTER_CARD"
	AMERICAN_EXPRESS CardAssociation = "AMERICAN_EXPRESS"
	TROY             CardAssociation = "TROY"
)

type CardType string

const (
	CREDIT_CARD  CardType = "CREDIT_CARD"
	DEBIT_CARD   CardType = "DEBIT_CARD"
	PREPAID_CARD CardType = "PREPAID_CARD"
)

type FraudStatus int8

const (
	APPROVED FraudStatus = 1
	WAITED   FraudStatus = 0
	DENIED   FraudStatus = -1
)

type TransactionStatus int8

const (
	TAPPROVED TransactionStatus = 2
	TWAITED   TransactionStatus = 1
	TDENIED   TransactionStatus = -1
	CHECKING  TransactionStatus = 0
)

type Currency string

const (
	TRY Currency = "TRY"
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	IRR Currency = "IRR"
	RUB Currency = "RUB"
	CNY Currency = "CNY"
	JPY Currency = "JPY"
	AUD Currency = "AUD"
)

type Category string

const (
	ELECTRONIC Category = "ELECTRONIC"
)

type Buyer struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Surname             string `json:"surname"`
	IdentityNumber      string `json:"identityNumber"`
	Email               string `json:"email"`
	GSMNumber           string `json:"gsmNumber"`
	RegistrationAddress string `json:"registrationAddress"`
	City                string `json:"city"`
	Country             string `json:"country"`
	IP                  string `json:"ip"`
}

type BillingAddress struct {
	Address     string `json:"address"`
	ContactName string `json:"contactName"`
	City        string `json:"city"`
	Country     string `json:"country"`
}

type BasketItems struct {
	ID        string   `json:"id"`
	Price     string   `json:"price"`
	Name      string   `json:"name"`
	Category1 Category `json:"category1"`
	ItemType  ItemType `json:"itemType"`
}

type ItemTransactions struct {
	PaymentTransactionID         string            `json:"paymentTransactionId"`
	ItemID                       string            `json:"itemId"`
	Price                        float32           `json:"price"`
	PaidPrice                    float32           `json:"paidPrice"`
	TransactionStatus            TransactionStatus `json:"transactionStatus"`
	BlockageRate                 float32           `json:"blockageRate"`
	BlockageRateAmountMerchant   float32           `json:"blockageRateAmountMerchant"`
	BlockageResolvedDate         string            `json:"blockageResolvedDate"`
	IyziCommissionFee            float32           `json:"iyziCommissionFee"`
	MerchantCommissionRate       float32           `json:"merchantCommissionRate"`
	MerchantCommissionRateAmount float32           `json:"merchantCommissionRateAmount"`
	ConvertedPayout              ConvertedPayout   `json:"convertedPayout"`
}

type ConvertedPayout struct {
	PaidPrice                  float32 `json:"paidPrice"`
	IyziCommissionFee          float32 `json:"iyziCommissionFee"`
	IyziCommissionRateAmount   float32 `json:"iyziCommissionRateAmount"`
	BlockageRateAmountMerchant float32 `json:"blockageRateAmountMerchant"`
	MerchantPayoutAmount       float32 `json:"merchantPayoutAmount"`
	IyziConversationRate       float32 `json:"iyziConversationRate"`
	IyziConversationRateAmount float32 `json:"iyziConversationRateAmount"`
	Currency                   string  `json:"currency"`
	MDStatus                   string  `json:"mdStatus"`
}
