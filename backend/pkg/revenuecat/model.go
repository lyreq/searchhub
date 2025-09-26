package revenuecat

type EventType string

const (
	EventType_InitialPurchase EventType = "INITIAL_PURCHASE"
	EventType_Cancellation    EventType = "CANCELLATION"
	EventType_ProductChange   EventType = "PRODUCT_CHANGE"
	EventType_Test            EventType = "TEST"
)

type PeriodType string

const (
	PeriodType_Trial       PeriodType = "TRIAL"
	PeriodType_Intro       PeriodType = "INTRO"
	PeriodType_Normal      PeriodType = "NORMAL"
	PeriodType_Promotional PeriodType = "PROMOTIONAL"
	PeriodType_Prepaid     PeriodType = "PREPAID"
)

type Store string

const (
	Store_Amazon      Store = "AMAZON"
	Store_Appstore    Store = "APPSTORE"
	Store_MacAppStore Store = "MAC_APP_STORE"
	Store_PlayStore   Store = "PLAY_STORE"
	Store_Promotional Store = "PROMOTIONAL"
	Store_Stripe      Store = "STRIPE"
)

type Environment string

const (
	Environment_Sandbox    Environment = "SANDBOX"
	Environment_Production Environment = "PRODUCTION"
)

type CancelReason string

const (
	CancelReason_Unsubscribe        CancelReason = "UNSUBSCRIBE"
	CancelReason_BillingError       CancelReason = "BILLING_ERROR"
	CancelReason_DeveloperInitiated CancelReason = "DEVELOPER_INITIATED"
	CancelReason_PriceIncrease      CancelReason = "PRICE_INCREASE"
	CancelReason_CustomerSupport    CancelReason = "CUSTOMER_SUPPORT"
	CancelReason_Unknown            CancelReason = "UNKNOWN"
)

type Currency string

const (
	Currency_TRY Currency = "TRY"
	Currency_USD Currency = "USD"
)
