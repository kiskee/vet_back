package websocket

const (
	EventRequestSearching  = "request.searching"
	EventRequestAssigned   = "request.assigned"
	EventRequestCancelled  = "request.cancelled"
	EventRequestStatus     = "request.status_update"
	EventVetLocationUpdate = "vet.location_update"
	EventVetIncomingReq    = "vet.incoming_request"
	EventVetAcceptReq      = "vet.accept_request"
	EventVetRejectReq      = "vet.reject_request"
	EventSubscribe         = "subscribe"
	EventUnsubscribe       = "unsubscribe"
)

type Event struct {
	Event     string      `json:"event"`
	RequestID string      `json:"request_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type SubscribePayload struct {
	RequestID string `json:"request_id"`
}

type IncomingRequestData struct {
	RequestID string  `json:"request_id"`
	ServiceID string  `json:"service_id"`
	ClientLat float64 `json:"client_lat"`
	ClientLng float64 `json:"client_lng"`
}

type AssignedData struct {
	VetID            string `json:"vet_id"`
	VetName          string `json:"vet_name"`
	EstimatedArrival string `json:"estimated_arrival"`
}

type LocationData struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
