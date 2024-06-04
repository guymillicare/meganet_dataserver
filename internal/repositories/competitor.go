package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"

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
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&nonEmptyCompetitors).Error

		if err != nil {
			return fmt.Errorf("CompetitorCreateBatch: %v", err)
		}
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
