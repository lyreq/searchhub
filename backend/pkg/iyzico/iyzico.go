package iyzico

import (
	"fmt"
	"strings"
)

type Iyzico interface {
	CheckoutFormInit(rq CheckoutFormInitRequest) (CheckoutFormInitResponse, error)
	CheckoutFormSearch(rq CheckoutFormSearchRequest) (CheckoutFormSearchResponse, error)
	CreateCheckoutFormRequest(price float64, itemID string, itemName string, buyer Buyer) CheckoutFormInitRequest

	IsSandbox() bool
}

type client struct {
	apiKey      string
	secretKey   string
	baseUrl     string
	callbackURL string
	service     IService
}

func New(apiKey, secretKey, baseUrl string, callbackURL string) *client {
	return &client{
		apiKey:      apiKey,
		secretKey:   secretKey,
		baseUrl:     baseUrl,
		callbackURL: callbackURL,
		service:     newService(),
	}
}

func (i *client) CheckoutFormInit(rq CheckoutFormInitRequest) (CheckoutFormInitResponse, error) {
	return i.service.CheckoutFormInit(rq, RequestInfo{
		ApiKey:    i.apiKey,
		SecretKey: i.secretKey,
		URL:       i.baseUrl + "/payment/iyzipos/checkoutform/initialize/auth/ecom",
	})
}

func (i *client) CheckoutFormSearch(rq CheckoutFormSearchRequest) (CheckoutFormSearchResponse, error) {
	return i.service.CheckoutFormSearch(rq, RequestInfo{
		ApiKey:    i.apiKey,
		SecretKey: i.secretKey,
		URL:       i.baseUrl + "/payment/iyzipos/checkoutform/auth/ecom/detail",
	})
}

func (s *client) CreateCheckoutFormRequest(price float64, itemID string, itemName string, buyer Buyer) CheckoutFormInitRequest {
	kdvPrice := price * 1.2

	itemPrice := fmt.Sprintf("%.2f", price)
	itemPrice, _ = strings.CutSuffix(itemPrice, "0")

	paidPrice := fmt.Sprintf("%.2f", kdvPrice)
	paidPrice, _ = strings.CutSuffix(paidPrice, "0")

	return CheckoutFormInitRequest{
		Locale:         TR,
		ConversationID: RandomString(ALPHANUMERIC, 6),
		Price:          itemPrice,
		BasketID:       RandomString(ALPHANUMERIC, 6),
		Buyer:          buyer,
		BillingAddress: BillingAddress{
			Address:     buyer.RegistrationAddress,
			ContactName: buyer.Name + " " + buyer.Surname,
			City:        buyer.City,
			Country:     buyer.Country,
		},
		BasketItems: []BasketItems{
			{
				ID:        itemID,
				Price:     itemPrice,
				Name:      itemName,
				Category1: ELECTRONIC,
				ItemType:  VIRTUAL,
			},
		},
		CallbackURL: s.callbackURL,
		Currency:    TRY,
		PaidPrice:   paidPrice,
	}
}

func (i *client) IsSandbox() bool {
	return i.baseUrl == "https://sandbox-api.iyzipay.com"
}
