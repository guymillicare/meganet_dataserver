package requests

type EventsList struct {
	EventIdList []string `json:"event_id_list"`
}

type MatchesRequest struct {
	Status      []string `json:"status,omitempty"`
	Sports      []int32  `json:"sports,omitempty"`
	Countries   []int32  `json:"countries,omitempty"`
	Tournaments []int32  `json:"tournaments,omitempty"`
	StartDate   int64    `json:"start_date,omitempty"`
	EndDate     int64    `json:"end_date,omitempty"`
}
