package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
)

func CountriesFindAll() ([]*types.CountryItem, error) {
	var countries []*types.CountryItem
	if err := database.DB.Table("countries").Find(&countries).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return countries, fmt.Errorf("CountriesFindAll: %v", err)
	}
	return countries, nil
}

func CountriesPreload() {
	database.COUNTRIES, _ = CountriesFindAll()
}

func GetCountryId(name string) (string, error) {
	for _, country := range database.COUNTRIES {
		if country.Name == name {
			return country.ReferenceId, nil
		}
	}
	return "0", fmt.Errorf("Country not found")
}
