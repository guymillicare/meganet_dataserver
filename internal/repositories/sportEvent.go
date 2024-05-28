package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

func CreateOrUpdateSportEvent(prematch *proto.Prematch) (*types.SportEventItem, error) {
	sportEvent, _ := GetSportEventFromRedis(prematch.Id)
	if sportEvent == nil {
		sportEvent = &types.SportEventItem{}
		sportEvent.ProviderId = 1
		sportEvent.ReferenceId = prematch.Id
		sport, _ := GetSportFromRedis(prematch.Sport)
		if sport == nil {
			fmt.Print("SPORT", prematch.Sport)
		}
		sportEvent.SportId = sport.ReferenceId
		if prematch.League != "" {
			country := strings.Split(prematch.League, " - ")[0]
			tournament := strings.Split(prematch.League, " - ")[0]
			if len(strings.Split(prematch.League, " - ")) > 1 {
				tournament = strings.Split(prematch.League, " - ")[1]
			}
			countryItem, _ := GetCountryFromRedis(country)
			if countryItem == nil {
				return sportEvent, fmt.Errorf("GetCountryFromRedis")
			}
			sportEvent.CountryId = countryItem.ReferenceId
			tournamentItem, _ := GetTournamentFromRedis(tournament)
			if tournamentItem == nil {
				return sportEvent, fmt.Errorf("GetCountryFromRedis")
			}
			sportEvent.TournamentId = tournamentItem.ReferenceId
		}
		sportEvent.Name = prematch.HomeTeam + " vs " + prematch.AwayTeam
		sportEvent.StartAt = prematch.StartDate
		sportEvent.Status = prematch.Status
		if err := database.DB.Table("sport_events").Create(&sportEvent).Error; err != nil {
			return sportEvent, fmt.Errorf("CreateSportEvent: %v", err)
		}
	} else {
		sportEvent.StartAt = prematch.StartDate
		sportEvent.Status = prematch.Status
		if err := database.DB.Table("sport_events").Save(&sportEvent).Error; err != nil {
			return sportEvent, fmt.Errorf("UpdateSportEvent: %v", err)
		}
	}

	// Save the updated or newly created sport event to Redis
	ctx := context.Background()
	if err := saveSportEventToRedis(ctx, sportEvent); err != nil {
		return sportEvent, fmt.Errorf("CreateOrUpdateSportEvent: error saving sport event to Redis: %v", err)
	}

	return sportEvent, nil
}

func UpdateSportEventStatus(event *types.SportEventItem) {
	if err := database.DB.Table("sport_events").Save(&event).Error; err != nil {
		fmt.Printf("UpdateSportEventStatus: %v\n", err)
	}
	ctx := context.Background()
	if err := saveSportEventToRedis(ctx, event); err != nil {
		fmt.Printf("CreateOrUpdateSportEvent: error saving sport event to Redis: %v\n", err)
	}
}

func SportEventsFindAll() ([]*types.SportEventItem, error) {
	var sportEvent []*types.SportEventItem
	if err := database.DB.Table("sport_events").Where("status!='Completed'").Find(&sportEvent).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return sportEvent, fmt.Errorf("SportEventsFindAll: %v", err)
	}

	return sportEvent, nil

}

func saveSportEventToRedis(ctx context.Context, sportEvent *types.SportEventItem) error {
	sportEventJSON, err := json.Marshal(sportEvent)
	if err != nil {
		return fmt.Errorf("saveSportEventToRedis: error marshaling sportEvent: %v", err)
	}

	key := fmt.Sprintf("sportEvent:%s", sportEvent.ReferenceId)

	// Define the expiration time as 90 days
	expiration := 90 * 24 * time.Hour

	err = database.RedisDB.Set(ctx, key, sportEventJSON, expiration).Err()
	if err != nil {
		return fmt.Errorf("saveSportEventToRedis: error saving sportEvent to Redis: %v", err)
	}

	return nil
}

func SportEventsPreload() {
	sportEvents, _ := SportEventsFindAll()
	ctx := context.Background()
	for _, sportEvent := range sportEvents {
		if err := saveSportEventToRedis(ctx, sportEvent); err != nil {
			fmt.Printf("saveSportEventToRedis: error saving sportEvent to Redis: %v\n", err)
			continue
		}
	}
}

func GetSportEventFromRedis(refId string) (*types.SportEventItem, error) {
	// Construct the key for the sportEvent item
	ctx := context.Background()
	key := fmt.Sprintf("sportEvent:%s", refId)

	// Retrieve the sportEvent item from Redis
	sportEventJSON, err := database.RedisDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// sportEvent item not found in Redis
			return nil, nil
		}
		return nil, fmt.Errorf("GetSportEventFromRedis: error fetching sportEvent item from Redis: %v", err)
	}

	// Deserialize the sportEvent item from JSON
	var sportEvent types.SportEventItem
	if err := json.Unmarshal(sportEventJSON, &sportEvent); err != nil {
		return nil, fmt.Errorf("GetSportEventFromRedis: error unmarshaling sportEvent item: %v", err)
	}

	return &sportEvent, nil
}

func SportEventsFindByFilters(systemId int32, providerId int, status string, sportId int32, countryId int32, tournamentId int32, offset int, limit int) ([]*types.SportEventFullItem, error) {
	var result []*types.SportEventFullItem
	query := database.DB.
		Table("sport_events").
		Select("sport_events.id as id",
			"sport_events.name as name",
			"sport_events.start_at as start_at",
			"sports.name as sport_name",
			"countries.name as country_name",
			"tournaments.name as tournament_name",
			"sport_events.home_score as home_score",
			"sport_events.away_score as away_score").
		Joins("LEFT JOIN sports ON sports.id = sport_events.sport_id").
		Joins("LEFT JOIN countries ON countries.id = sport_events.country_id").
		Joins("LEFT JOIN tournaments ON tournaments.id = sport_events.tournament_id").
		Joins("LEFT JOIN system_sports ON system_sports.sport_id = sports.id").
		Where("sport_events.provider_id=? AND sport_events.status = ?", providerId, status)
	if systemId > 0 {
		query = query.Where("system_sports.system_id=?", systemId)
	}
	if sportId > 0 {
		query = query.Where("sport_events.sport_id=?", sportId)
	}
	if countryId > 0 {
		query = query.Where("sport_events.country_id=?", countryId)
	}
	if tournamentId > 0 {
		query = query.Where("sport_events.tournament_id=?", tournamentId)
	}
	if err := query.Offset(offset).Limit(limit).Find(&result).Order("created_at").Error; err != nil {
		return nil, fmt.Errorf("SportEventsFindByFilters: %v", err)
	}
	return result, nil
}
