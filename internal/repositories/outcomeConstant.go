package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
)

func OutcomeConstantFind(reference_id string) (*types.OutcomeConstantItem, error) {
	var outcomeConstant *types.OutcomeConstantItem
	if err := database.DB.Table("outcome_constants").Where("reference_id =?", reference_id).First(&outcomeConstant).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return outcomeConstant, fmt.Errorf("OutcomeConstantFind: %v", err)
	}
	return outcomeConstant, nil
}
