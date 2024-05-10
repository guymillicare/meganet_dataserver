package types

import "time"

type SportItem struct {
	Id          int32     `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Slug        string    `json:"slug"`
	ReferenceId string    `json:"reference_id"`
	Order       int32     `json:"order"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"`
}

type CountryItem struct {
	Id          int32     `json:"id"`
	Name        string    `json:"name"`
	Abbr        string    `json:"abbr"`
	ReferenceId string    `json:"reference_id"`
	Order       int32     `json:"order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"`
}

type TournamentItem struct {
	Id          int32     `json:"id"`
	ReferenceId string    `json:"reference_id"`
	SportId     string    `json:"sport_id"`
	CountryId   string    `json:"country_id"`
	Name        string    `json:"name"`
	Abbr        string    `json:"abbr"`
	Order       int32     `json:"order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"`
}

type CompetitorItem struct {
	Id          int32     `json:"id"`
	ReferenceId string    `json:"reference_id"`
	CountryId   int32     `json:"country_id"`
	Name        string    `json:"name"`
	Abbr        string    `json:"abbr"`
	SportId     string    `json:"sport_id`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"`
}

type SportEventItem struct {
	Id           int32     `json:"id"`
	ProviderId   int32     `json:"provider_id"`
	ReferenceId  string    `json:"reference_id"`
	SportId      string    `json:"sport_id"`
	CountryId    string    `json:"country_id"`
	TournamentId string    `json:"tournament"`
	Name         string    `json:"name"`
	StartAt      string    `json:"start_at"`
	Status       string    `json:"status"`
	Active       int32     `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"-"`
}

type MarketOutcomeItem struct {
	Id                int32  `json:"id"`
	MarketRefId       string `json:"market_ref_id"`
	MarketDescription string `json:"market_description"`
	OutcomeRefId      string `json:"outcome_ref_id"`
	OutcomeName       string `json:"outcome_name"`
	SportRefId        string `json:"sport_ref_id"`
	Order             int32  `json:"order"`
}

type MarketConstantItem struct {
	Id          int32  `json:"id"`
	ReferenceId string `json:"reference_id"`
	Description string `json:"description"`
}

type OutcomeConstantItem struct {
	Id          int    `json:"id"`
	ReferenceId string `json:"reference_id"`
	Name        string `json:"name"`
	Order       int    `json:"order"`
}

type SportMarketGroupItem struct {
	Id         int32  `json:"id"`
	SportId    int32  `json:"sport_id"`
	MarketId   int32  `json:"market_id"`
	SportName  string `json:"sport_name"`
	MarketName string `json:"market_name"`
}

type OutcomeItem struct {
	Id          int32     `json:"id"`
	ReferenceId string    `json:"reference_id"`
	EventId     int32     `json:"event_id"`
	MarketId    int32     `json:"market_id"`
	Name        string    `json:"name"`
	Odds        float64   `json:"odds"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OddsStream struct {
	Data []struct {
		BetName         string  `json:"bet_name"`
		BetPoints       float64 `json:"bet_points"`
		BetPrice        float64 `json:"bet_price"`
		BetType         string  `json:"bet_type"`
		GameId          string  `json:"game_id"`
		Id              string  `json:"id"`
		IsLive          bool    `json:"is_live"`
		IsMain          bool    `json:"is_main"`
		League          string  `json:"league"`
		PlayerId        string  `json:"player_id"`
		Selection       string  `json:"selection"`
		SelectionLine   string  `json:"selection_line"`
		SelectionPoints float64 `json:"selection_points"`
		Sport           string  `json:"sport"`
		Sportsbook      string  `json:"sportsbook"`
		Timestamp       float64 `json:"timestamp"`
	} `json:"data"`
	EntryId string `json:"entry_id"`
	Type    string `json:"type"`
}
