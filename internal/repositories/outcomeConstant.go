package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
	"sync"
)

var outcomeConstantCache sync.Map

func OutcomeConstantFind(reference_id string) (*types.OutcomeConstantItem, error) {
	// Check if the item is in the cache
	if value, ok := outcomeConstantCache.Load(reference_id); ok {
		return value.(*types.OutcomeConstantItem), nil
	}

	// If not in the cache, query the database
	var outcomeConstant *types.OutcomeConstantItem
	if err := database.DB.Table("outcome_constants").Where("reference_id =?", reference_id).First(&outcomeConstant).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return outcomeConstant, fmt.Errorf("OutcomeConstantFind: %v", err)
	}

	// Store the result in the cache
	outcomeConstantCache.Store(reference_id, outcomeConstant)
	return outcomeConstant, nil
}
