package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"
	"strings"
)

func CreateSportEvent(prematch *proto.Prematch) (*types.SportEventItem, error) {
	newSportEvent := &types.SportEventItem{}
	sportEvent, _ := SportEventFindByReferenceId(prematch.Id)
	if sportEvent == nil {
		newSportEvent.ProviderId = 1
		newSportEvent.ReferenceId = prematch.Id
		newSportEvent.SportId, _ = GetSportRefId(prematch.Sport)
		if prematch.League != "" {
			country := strings.Split(prematch.League, " - ")[0]
			tournament := strings.Split(prematch.League, " - ")[0]
			if len(strings.Split(prematch.League, " - ")) > 1 {
				tournament = strings.Split(prematch.League, " - ")[1]
			}
			newSportEvent.CountryId, _ = GetCountryId(country)
			newSportEvent.TournamentId, _ = GetTournamentId(tournament)
		}
		newSportEvent.Name = prematch.HomeTeam + " vs " + prematch.AwayTeam
		newSportEvent.StartAt = prematch.StartDate
		newSportEvent.Status = prematch.Status
		if err := database.DB.Table("sport_events").Create(&newSportEvent).Error; err != nil {
			return newSportEvent, fmt.Errorf("CreateSportEvent: %v", err)
		}
		return newSportEvent, nil
	}
	return sportEvent, nil
}

func SportEventFindByReferenceId(id string) (*types.SportEventItem, error) {
	var sportEvent *types.SportEventItem
	if err := database.DB.Table("sport_events").Where("reference_id =?", id).First(&sportEvent).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return sportEvent, fmt.Errorf("FindById: %v", err)
	}
	return sportEvent, nil
}
