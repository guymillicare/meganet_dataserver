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
		sport, _ := GetSportFromRedis(prematch.Sport)
		newSportEvent.SportId = sport.ReferenceId
		if prematch.League != "" {
			country := strings.Split(prematch.League, " - ")[0]
			tournament := strings.Split(prematch.League, " - ")[0]
			if len(strings.Split(prematch.League, " - ")) > 1 {
				tournament = strings.Split(prematch.League, " - ")[1]
			}
			countryItem, _ := GetCountryFromRedis(country)
			newSportEvent.CountryId = countryItem.ReferenceId
			tournamentItem, _ := GetTournamentFromRedis(tournament)
			newSportEvent.TournamentId = tournamentItem.ReferenceId
		}
		newSportEvent.Name = prematch.HomeTeam + " vs " + prematch.AwayTeam
		newSportEvent.StartAt = prematch.StartDate
		newSportEvent.Status = prematch.Status
		if err := database.DB.Table("sport_events").Create(&newSportEvent).Error; err != nil {
			return newSportEvent, fmt.Errorf("CreateSportEvent: %v", err)
		}
		return newSportEvent, nil
	} else {
		sportEvent.StartAt = prematch.StartDate
		sportEvent.Status = prematch.Status
		if err := database.DB.Table("sport_events").Save(&sportEvent).Error; err != nil {
			return sportEvent, fmt.Errorf("UpdateSportEvent: %v", err)
		}
		return sportEvent, nil
	}
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
