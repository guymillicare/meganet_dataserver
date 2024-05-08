package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
)

func SportsFindAll() ([]*types.SportItem, error) {
	var sports []*types.SportItem
	if err := database.DB.Table("sports").Find(&sports).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return sports, fmt.Errorf("SportsFindAll: %v", err)
	}
	return sports, nil
}

func SportsPreload() {
	database.SPORTS, _ = SportsFindAll()
}

func GetSportId(name string) (string, error) {
	for _, sport := range database.SPORTS {
		if sport.Slug == name {
			return sport.ReferenceId, nil
		}
	}
	return "0", fmt.Errorf("Sport not found")
}
