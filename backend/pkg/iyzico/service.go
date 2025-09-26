package iyzico

import (
	"encoding/json"
)

type IService interface {
	CheckoutFormInit(rq CheckoutFormInitRequest, rqInfo RequestInfo) (CheckoutFormInitResponse, error)
	CheckoutFormSearch(rq CheckoutFormSearchRequest, rqInfo RequestInfo) (CheckoutFormSearchResponse, error)
}

type service struct {
}

func newService() IService {
	return &service{}
}

func (s *service) CheckoutFormInit(req CheckoutFormInitRequest, reqInfo RequestInfo) (CheckoutFormInitResponse, error) {
	randStr := RandomString(NUMERIC, 6)
	hashedStr := GenerateHash(reqInfo.ApiKey + randStr + reqInfo.SecretKey + GenerateRequestString(req))
	reqStr := "IYZWS" + " " + reqInfo.ApiKey + ":" + hashedStr

	headers := GetHTTPHeaders(reqStr, randStr)

	res, err := HTTP(POST, reqInfo.URL, &headers, req)
	if err != nil {
		return CheckoutFormInitResponse{}, err
	}

	var body CheckoutFormInitResponse

	err = json.Unmarshal(res, &body)
	if err != nil {
		return CheckoutFormInitResponse{}, err
	}

	return body, nil
}

func (s *service) CheckoutFormSearch(req CheckoutFormSearchRequest, reqInfo RequestInfo) (CheckoutFormSearchResponse, error) {
	randStr := RandomString(NUMERIC, 6)
	hashedStr := GenerateHash(reqInfo.ApiKey + randStr + reqInfo.SecretKey + GenerateRequestString(req))
	reqStr := "IYZWS" + " " + reqInfo.ApiKey + ":" + hashedStr

	headers := GetHTTPHeaders(reqStr, randStr)

	res, err := HTTP(POST, reqInfo.URL, &headers, req)
	if err != nil {
		return CheckoutFormSearchResponse{}, err
	}

	var body CheckoutFormSearchResponse

	err = json.Unmarshal(res, &body)
	if err != nil {
		return CheckoutFormSearchResponse{}, err
	}

	return body, nil
}
