package types

import (
	"net/http"
	"time"
)

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
	Logo        string    `json:"logo"`
	SportId     string    `json:"sport_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"`
}

type SportEventItem struct {
	Id             int32     `json:"id"`
	ProviderId     int32     `json:"provider_id"`
	ReferenceId    string    `json:"reference_id"`
	SportId        string    `json:"sport_id"`
	CountryId      string    `json:"country_id"`
	TournamentId   string    `json:"tournament"`
	Name           string    `json:"name"`
	StartAt        string    `json:"start_at"`
	Status         string    `json:"status"`
	Active         int32     `json:"active"`
	HomeScore      int32     `json:"home_score"`
	AwayScore      int32     `json:"away_score"`
	StatsperformId string    `json:"statsperform_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"-"`
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
	Active      bool      `json:"active"`
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

type GameScoreStream struct {
	Data struct {
		GameID string `json:"game_id"`
		Score  struct {
			Clock                    string `json:"clock"`
			ScoreAwayPeriod1         int    `json:"score_away_period_1"`
			ScoreAwayPeriod1Tiebreak *int   `json:"score_away_period_1_tiebreak,omitempty"`
			ScoreAwayPeriod2         int    `json:"score_away_period_2"`
			ScoreAwayTotal           int    `json:"score_away_total"`
			ScoreHomePeriod1         int    `json:"score_home_period_1"`
			ScoreHomePeriod1Tiebreak *int   `json:"score_home_period_1_tiebreak,omitempty"`
			ScoreHomePeriod2         int    `json:"score_home_period_2"`
			ScoreHomeTotal           int    `json:"score_home_total"`
		} `json:"score"`
	} `json:"data"`
	EntryId string `json:"entry_id"`
}

type GenericResponse struct{}

func (rd *GenericResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type OutcomeListResponse struct {
	GenericResponse
	SportEvent *SportEventItem `json:"sportEvent"`
	Outcome    []*OutcomeItem  `json:"outcome"`
}

type SportEventListResponse struct {
	GenericResponse
	SportEventList []*SportEventOddsItem `json:"sportEventList"`
}

type SportEventsListResponse struct {
	GenericResponse
	SportEventsList []*SportEventsOddsItem `json:"sportEventsList"`
}

type TeamInfoResponse struct {
	Data []struct {
		Id               string `json:"id"`
		TeamName         string `json:"team_name"`
		TeamCity         string `json:"team_city"`
		TeamMascot       string `json:"team_mascot"`
		TeamAbbreviation string `json:"team_abbreviation"`
		Sport            string `json:"sport"`
		League           string `json:"league"`
		Logo             string `json:"logo"`
	} `json:"data"`
}

type GameScoreResponse struct {
	Data []struct {
		GameId         string  `json:"game_id"`
		ScoreHomeTotal int     `json:"score_home_total"`
		ScoreAwayTotal int     `json:"score_away_total"`
		Clock          *string `json:"clock,omitempty"`
		Sport          string  `json:"sport"`
		League         string  `json:"league"`
		// Period                 int     `json:"period"`
		Status                 string  `json:"status"`
		IsLive                 bool    `json:"is_live"`
		Duration               *string `json:"duration,omitempty"`
		AwayTeam               string  `json:"away_team"`
		HomeTeam               string  `json:"home_team"`
		StartDate              string  `json:"start_date"`
		Description            string  `json:"description"`
		CurrentOuts            *int    `json:"current_outs,omitempty"`
		CurrentBalls           *int    `json:"current_balls,omitempty"`
		CurrentStrikes         *int    `json:"current_strikes,omitempty"`
		ScoreAwayPeriod1       *int    `json:"score_away_period_1,omitempty"`
		ScoreAwayPeriod2       *int    `json:"score_away_period_2,omitempty"`
		ScoreAwayPeriod3       *int    `json:"score_away_period_3,omitempty"`
		ScoreAwayPeriod4       *int    `json:"score_away_period_4,omitempty"`
		ScoreAwayPeriod5       *int    `json:"score_away_period_5,omitempty"`
		ScoreAwayPeriod6       *int    `json:"score_away_period_6,omitempty"`
		ScoreAwayPeriod7       *int    `json:"score_away_period_7,omitempty"`
		ScoreAwayPeriod8       *int    `json:"score_away_period_8,omitempty"`
		ScoreAwayPeriod9       *int    `json:"score_away_period_9,omitempty"`
		ScoreHomePeriod1       *int    `json:"score_home_period_1,omitempty"`
		ScoreHomePeriod2       *int    `json:"score_home_period_2,omitempty"`
		ScoreHomePeriod3       *int    `json:"score_home_period_3,omitempty"`
		ScoreHomePeriod4       *int    `json:"score_home_period_4,omitempty"`
		ScoreHomePeriod5       *int    `json:"score_home_period_5,omitempty"`
		ScoreHomePeriod6       *int    `json:"score_home_period_6,omitempty"`
		ScoreHomePeriod7       *int    `json:"score_home_period_7,omitempty"`
		ScoreHomePeriod8       *int    `json:"score_home_period_8,omitempty"`
		ScoreHomePeriod9       *int    `json:"score_home_period_9,omitempty"`
		RunnerOnFirst          *string `json:"runner_on_first,omitempty"`
		DecisionMethod         *string `json:"decision_method,omitempty"`
		Decision               *string `json:"decision,omitempty"`
		Broadcast              *string `json:"broadcast,omitempty"`
		HomeStarter            *string `json:"home_starter,omitempty"`
		LastPlay               *string `json:"last_play,omitempty"`
		Weather                *string `json:"weather,omitempty"`
		HomeTeamCity           *string `json:"home_team_city,omitempty"`
		RunnerOnThird          *string `json:"runner_on_third,omitempty"`
		CurrentDownAndDistance *string `json:"current_down_and_distance,omitempty"`
		HomeTeamAbb            *string `json:"home_team_abb,omitempty"`
		AwayTeamAbb            *string `json:"away_team_abb,omitempty"`
		AwayStarter            *string `json:"away_starter,omitempty"`
		SeasonType             *string `json:"season_type,omitempty"`
		VenueName              *string `json:"venue_name,omitempty"`
		WeatherTemp            *string `json:"weather_temp,omitempty"`
		Attendance             *string `json:"attendance,omitempty"`
		CurrentBatterName      *string `json:"current_batter_name,omitempty"`
		CurrentPitcherName     *string `json:"current_pitcher_name,omitempty"`
		SeasonWeek             *string `json:"season_week,omitempty"`
		HomeTeamName           *string `json:"home_team_name,omitempty"`
		RunnerOnSecond         *string `json:"runner_on_second,omitempty"`
		VenueLocation          *string `json:"venue_location,omitempty"`
		WeatherTempHigh        *string `json:"weather_temp_high,omitempty"`
		AwayTeamCity           *string `json:"away_team_city,omitempty"`
		Capacity               *string `json:"capacity,omitempty"`
		SeasonYear             string  `json:"season_year"`
		AwayTeamName           *string `json:"away_team_name,omitempty"`
		CurrentPossession      *string `json:"current_possession,omitempty"`
	} `json:"data"`
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
}

type SportEventFullItem struct {
	Id             int32     `json:"id"`
	ReferenceId    string    `json:"reference_id"`
	Name           string    `json:"name"`
	StartAt        time.Time `json:"start_at"`
	Active         bool      `json:"active"`
	SportName      string    `json:"sport_name"`
	CountryName    string    `json:"country_name"`
	TournamentName string    `json:"tournament_name"`
	HomeScore      int32     `json:"home_score"`
	AwayScore      int32     `json:"away_score"`
}

type SportEventOddsItem struct {
	SportEvent *SportEventFullItem `json:"sportEvent"`
	Outcome    []*OutcomeItem      `json:"outcome"`
}

type SportEventsOddsItem struct {
	SportEvent *SportEventItem `json:"sportEvent"`
	Outcome    []*OutcomeItem  `json:"outcome"`
}
