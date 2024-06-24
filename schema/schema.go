package schema

import "fmt"

type OktaToken struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int32  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

type OktaError struct {
	Id      string   `json:"errorId"`
	Code    string   `json:"errorCode"`
	Summary string   `json:"errorSummary"`
	Link    string   `json:"errorLink"`
	Causes  []string `json:"errorCauses"`
}

// Error implements error.
func (e *OktaError) Error() string {
	return fmt.Sprintf("ErrorCode: %s, ErrorSummary: %s, ErrorLink: %s, ErrorId: %s, ErrorCauses: %v",
		e.Code, e.Summary, e.Link, e.Id, e.Causes)
}

type AtlasPayload[T any] struct {
	ExchangeId      int               `json:"exchangeId"`
	PublishTimeInMs int64             `json:"publishTimeInMs"`
	ReceiveTimeInMs int64             `json:"receiveTimeInMs"`
	Messages        []AtlasMessage[T] `json:"messages"`
}

type AtlasMessage[T any] struct {
	IsJson      bool `json:"isJson"`
	JsonMessage T    `json:"jsonMessage"`
}

type RfidHubScan struct {
	TagHexEpc      string                 `json:"tagHexEpc"`
	Timestamp      int64                  `json:"timestamp"`
	TrackingNumber int64                  `json:"trackingNumber"`
	IsFedexTag     bool                   `json:"isFedexTag"`
	Center         string                 `json:"center"`
	ExtraData      map[string]interface{} `json:"extra_data"`
}
