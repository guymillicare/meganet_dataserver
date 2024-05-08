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
	Sports            string `json:"sports"`
	Order             int32  `json:"order"`
}

type MarketConstantItem struct {
	Id          int32  `json:"id"`
	ReferenceId string `json:"reference_id"`
	Description string `json:"description"`
}

type OutcomeConstantItem struct {
	Id          int32  `json:"id"`
	ReferenceId string `json:"reference_id"`
	Name        string `json:"name"`
	Order       int32  `json:"order"`
}
