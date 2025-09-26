package revenuecat

type Event struct {
	ID                string    `json:"id"`
	AppID             string    `json:"app_id"`
	AppUserID         string    `json:"app_user_id"`
	OriginalAppUserID string    `json:"original_app_user_id"`
	Aliases           []string  `json:"aliases"`
	Type              EventType `json:"type"`
	EventTimestampMs  int64     `json:"event_timestamp_ms"`
}

type SubscriberAttributes struct {
	Attributes map[string]struct {
		UpdatedAtMs int64  `json:"updated_at_ms"`
		Value       string `json:"value"`
	} `json:"-"`
}

type EventProperty struct {
	Event
	SubscriberAttributes      `json:"subscriber_attributes"`
	ProductID                 string       `json:"product_id"`
	EntitlementID             string       `json:"entitlement_id"`
	EntitlementIDs            []string     `json:"entitlement_ids"`
	PeriodType                PeriodType   `json:"period_type"`
	PurchasedAtMs             uint64       `json:"purchased_at_ms"`
	GracePeriodExpirationAtMs uint64       `json:"grace_period_expiration_at_ms"` // only BILLING_ISSUE events
	ExpirationAtMs            uint64       `json:"expiration_at_ms"`
	AutoResumeAtMs            uint64       `json:"auto_resume_at_ms"` // only SUBSCRIPTION_PAUSED events
	Store                     Store        `json:"store"`
	Environment               Environment  `json:"environment"`
	IsTrialConversion         bool         `json:"is_trial_conversion"` // only RENEWAL events
	CancelReason              CancelReason `json:"cancel_reason"`       // only CANCELLATION evets
	ExpirationReason          CancelReason `json:"expiration_reason"`   // only EXPIRATION events
	NewProductID              string       `json:"new_product_id"`      // only PRODUCT_CHANGE events
	PresentedOfferingID       string       `json:"presented_offering_id"`
	Price                     float32      `json:"price"`
	Currency                  string       `json:"currency"`
	PriceInPurchasedCurrency  float32      `json:"price_in_purchased_currency"`
	TaxPercentage             float32      `json:"tax_percentage"`
	CommissionPercentage      float32      `json:"commission_percentage"`
	TakehomePercentage        float32      `json:"takehome_percentage"`
	TransactionID             string       `json:"transaction_id"`
	OriginalTransactionID     string       `json:"original_transaction_id"`
	IsFamilyShare             bool         `json:"is_family_share"`
	TransferredFrom           []string     `json:"transferred_from"` // only TRANSFER events
	TransferredTo             []string     `json:"transferred_to"`   // only TRANSFER events
	CountryCode               string       `json:"country_code"`
	OfferCode                 string       `json:"offer_code"` // only SUBSCRIBER_ALIAS or TRANSFER
}

type EventRequest struct {
	ApiVersion    string        `json:"api_version"`
	EventProperty EventProperty `json:"event"`
}
