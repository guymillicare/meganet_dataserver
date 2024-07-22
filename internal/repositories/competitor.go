package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"

	"github.com/go-redis/redis"
	"gorm.io/gorm/clause"
)

// CreateCompetitorsBatch inserts competitors into the database in a single batch.
func CreateCompetitorsBatch(competitors []*types.CompetitorItem) error {
	nonEmptyCompetitors := []*types.CompetitorItem{}
	for _, competitor := range competitors {
		if competitor != nil && competitor.Name != "" {
			nonEmptyCompetitors = append(nonEmptyCompetitors, competitor)
		}
	}

	if len(nonEmptyCompetitors) > 0 {
		err := database.DB.Table("competitors").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "reference_id"}},
			DoNothing: true,
		}).Create(&nonEmptyCompetitors).Error

		if err != nil {
			return fmt.Errorf("CompetitorCreateBatch: %v", err)
		}
	}
	ctx := context.Background()
	for _, competitor := range nonEmptyCompetitors {
		saveCompetitorToRedis(ctx, competitor)
	}

	return nil
}

func CompetitorFindByName(name string) (*types.CompetitorItem, error) {
	var competitor *types.CompetitorItem
	if err := database.DB.Table("competitors").Where("name =?", name).First(&competitor).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return competitor, fmt.Errorf("CompetitorFindOne: %v", err)
	}
	return competitor, nil
}

// PrepareCompetitors collects competitors from a prematch
func PrepareCompetitors(BaseURL string, APIKey string, prematch *proto.Prematch) []*types.CompetitorItem {
	sport, err := GetSportFromRedis(prematch.Sport)
	if err != nil {
		fmt.Println("GetSportFromRedis:", err)
		return nil
	}
	if prematch.HomeTeamInfo == nil || prematch.AwayTeamInfo == nil {
		return nil
	}

	url := fmt.Sprintf(
		"%s/api/v2/teams?key=%s&id=%s&id=%s&include_logos=true",
		BaseURL,
		APIKey,
		prematch.HomeTeamInfo.Id,
		prematch.AwayTeamInfo.Id,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var teamInfoResponse types.TeamInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&teamInfoResponse); err != nil {
		return nil
	}

	// teams := []*proto.TeamInfo{prematch.HomeTeamInfo, prematch.AwayTeamInfo}
	competitors := make([]*types.CompetitorItem, 2)

	for i, teamInfo := range teamInfoResponse.Data {
		competitors[i] = &types.CompetitorItem{
			Name:        teamInfo.TeamName,
			ReferenceId: teamInfo.Id,
			SportId:     sport.ReferenceId,
			Logo:        teamInfo.Logo,
		}
	}

	return competitors
}

func CompetitorsFindAll() ([]*types.CompetitorItem, error) {
	var competitors []*types.CompetitorItem
	if err := database.DB.Table("competitors").Where("data_feed='huge_data'").Find(&competitors).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return competitors, fmt.Errorf("CompetitorsFindAll: %v", err)
	}

	return competitors, nil

}

func saveCompetitorToRedis(ctx context.Context, competitor *types.CompetitorItem) error {
	competitorJSON, err := json.Marshal(competitor)
	if err != nil {
		return fmt.Errorf("saveCompetitorToRedis: error marshaling competitor: %v", err)
	}

	key := fmt.Sprintf("competitor:%s", competitor.ReferenceId)

	err = database.RedisDB.Set(ctx, key, competitorJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("saveCompetitorToRedis: error saving competitor to Redis: %v", err)
	}

	return nil
}

func CompetitorsPreload() {
	competitors, _ := CompetitorsFindAll()
	ctx := context.Background()
	for _, competitor := range competitors {
		if err := saveCompetitorToRedis(ctx, competitor); err != nil {
			fmt.Printf("CompetitorsPreload: error saving sportEvent to Redis: %v\n", err)
			continue
		}
	}
}

func GetCompetitorFromRedis(refId string) (*types.CompetitorItem, error) {
	// Construct the key for the competitor item
	ctx := context.Background()
	key := fmt.Sprintf("competitor:%s", refId)

	// Retrieve the competitor item from Redis
	compeitorJSON, err := database.RedisDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// competitor item not found in Redis
			return nil, nil
		}
		return nil, fmt.Errorf("GetCompetitorFromRedis: error fetching competitor item from Redis: %v", err)
	}

	// Deserialize the competitor item from JSON
	var competitor types.CompetitorItem
	if err := json.Unmarshal(compeitorJSON, &competitor); err != nil {
		return nil, fmt.Errorf("GetCompetitorFromRedis: error unmarshaling competitor item: %v", err)
	}

	return &competitor, nil
}
