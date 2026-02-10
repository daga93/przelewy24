package przelewy24

// NotificationBody is the body of webhook notification that you get
// after success payment.
type NotificationBody struct {
	MerchantId   int      `json:"merchantId"`
	PosId        int      `json:"posId"`
	SessionId    string   `json:"sessionId"`
	Amount       int      `json:"amount"`
	OriginAmount int      `json:"originAmount"`
	Currency     Currency `json:"currency"`
	OrderId      int      `json:"orderId"`
	MethodId     int      `json:"methodId"`
	Statement    string   `json:"statement"`
	Sign         string   `json:"sign"`
}

// RefundNotificationBody is the body of webhook notification that you get
// after success refund.
type RefundNotificationBody struct {
	OrderId     int      `json:"orderId"`
	SessionId   string   `json:"sessionId"`
	MerchantId  int      `json:"merchantId"`
	RequestId   string   `json:"requestId"`
	RefundsUuid string   `json:"refundsUuid"`
	Amount      int      `json:"amount"`
	Currency    Currency `json:"currency"`
	Timestamp   int      `json:"timestamp"`
	Status      int      `json:"status"`
	Sign        string   `json:"sign"`
}

type Currency string

type Country string

type Language string

type Channel int

type CartParameter struct {
	SellerID       string
	SellerCategory string
	Name           string
	Description    string
	Quantity       int
	Price          int
	Number         string
}
type Cart []CartParameter

type ShippingType int

type Shipping struct {
	ShippingType ShippingType `json:"type"`
	Address      string
	Zip          string
	City         string
	Country      string
}

type PSU struct {
	IP        string
	userAgent string
}

type AdditionalData struct {
	Shipping Shipping
	PSU      PSU
}

type verifyTransactionData struct {
	MerchantId int      `json:"merchantId"`
	PosId      int      `json:"posId"`
	SessionId  string   `json:"sessionId"`
	Amount     int      `json:"amount"`
	Currency   Currency `json:"currency"`
	OrderId    int      `json:"orderId"`
}

type verifyTransactionRequestBody struct {
	MerchantId int      `json:"merchantId"`
	PosId      int      `json:"posId"`
	SessionId  string   `json:"sessionId"`
	Amount     int      `json:"amount"`
	Currency   Currency `json:"currency"`
	OrderId    int      `json:"orderId"`
	Sign       string   `json:"sign"`
}

type registerTransactionRequestBody struct {
	MerchantId       int             `json:"merchantId"`
	PosId            int             `json:"posId"`
	SessionId        string          `json:"sessionId"`
	Amount           int             `json:"amount"`
	Currency         Currency        `json:"currency"`
	Description      string          `json:"description"`
	Email            string          `json:"email"`
	Client           string          `json:"client,omitempty"`
	Address          string          `json:"address,omitempty"`
	Zip              string          `json:"zip,omitempty"`
	City             string          `json:"city,omitempty"`
	Country          Country         `json:"country"`
	Phone            string          `json:"phone,omitempty"`
	Language         Language        `json:"language"`
	Method           int             `json:"method,omitempty"`
	UrlReturn        string          `json:"urlReturn"`
	UrlStatus        string          `json:"urlStatus,omitempty"`
	UrlNotify        string          `json:"urlNotify,omitempty"`
	TimeLimit        int             `json:"timeLimit,omitempty"` // In minutes
	Channel          Channel         `json:"channel,omitempty"`
	WaitForResult    bool            `json:"waitForResult,omitempty"`
	RegulationAccept bool            `json:"regulationAccept,omitempty"`
	Shipping         int             `json:"shipping,omitempty"`
	TransferLabel    string          `json:"transferLabel,omitempty"`
	Sign             string          `json:"sign,omitempty"`
	Encoding         string          `json:"encoding,omitempty"`
	MethodRefId      string          `json:"methodRefId,omitempty"`
	Cart             Cart            `json:"cart,omitempty"`
	Additional       *AdditionalData `json:"additional,omitempty"`
}

type registerTransactionResponseBody struct {
	ErrorData    string                       `json:"error,omitempty"`
	Code         int                          `json:"code,omitempty"`
	ResponseCode int                          `json:"responseCode,omitempty"`
	Data         registerTransactionDataField `json:"Data,omitempty"`
}
type verifyTransactionResponseBody struct {
	ErrorData    string                     `json:"error,omitempty"`
	Code         int                        `json:"code,omitempty"`
	ResponseCode int                        `json:"responseCode,omitempty"`
	Data         verifyTransactionDataField `json:"Data,omitempty"`
}
type registerTransactionDataField struct {
	Token string `json:"token"`
}
type verifyTransactionDataField struct {
	Status string `json:"status"`
}

type refund struct {
	OrderId     int    `json:"orderId"`
	SessionId   string `json:"sessionId"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type refundRequestBody struct {
	RequestId   string   `json:"requestId"`
	Refunds     []refund `json:"refunds"`
	RefundsUuid string   `json:"refundsUuid"`
	UrlStatus   string   `json:"urlStatus"`
}

type GetTransactionDetailsResponseBody struct {
	Data         TransactionDetails `json:"data"`
	ResponseCode int                `json:"responseCode"`
}

type TransactionDetails struct {
	Statement         string   `json:"statement"`
	OrderId           int      `json:"orderId"`
	SessionId         string   `json:"sessionId"`
	Status            int      `json:"status"`
	Amount            int      `json:"amount"`
	Currency          Currency `json:"currency"`
	Date              string   `json:"date"`
	DateOfTransaction string   `json:"dateOfTransaction"`
	ClientEmail       string   `json:"clientEmail"`
	AccountMD5        string   `json:"accountMD5"`
	PaymentMethod     int      `json:"paymentMethod"`
	Description       string   `json:"description"`
	ClientName        string   `json:"clientName"`
	ClientAddress     string   `json:"clientAddress"`
	ClientCity        string   `json:"clientCity"`
	ClientPostcode    string   `json:"clientPostcode"`
	BatchId           int      `json:"batchId"`
	Fee               string   `json:"fee"`
}

type RefundResponseBody struct {
	Data         []RefundResponseData `json:"data"`
	ResponseCode int                  `json:"responseCode"`
}

type RefundResponseData struct {
	OrderId     int    `json:"orderId"`
	SessionId   string `json:"sessionId"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	Message     string `json:"message"`
}
