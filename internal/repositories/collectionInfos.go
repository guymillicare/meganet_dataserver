package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"

	"github.com/go-redis/redis"
)

func CollectionInfosFindAll() ([]*types.CollectionInfoItem, error) {
	var collectionInfos []*types.CollectionInfoItem
	if err := database.DB.Table("collection_infos").Find(&collectionInfos).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return collectionInfos, fmt.Errorf("CollectionInfosFindAll: %v", err)
	}
	return collectionInfos, nil
}

func saveCollectionInfoToRedis(ctx context.Context, collectionInfo *types.CollectionInfoItem) error {
	collectionInfoJSON, err := json.Marshal(collectionInfo)
	if err != nil {
		return fmt.Errorf("saveCollectionInfoToRedis: error marshaling collectionInfo: %v", err)
	}

	key := fmt.Sprintf("collectionInfo:%d", collectionInfo.Id)
	err = database.RedisDB.Set(ctx, key, collectionInfoJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("saveCollectionInfoToRedis: error saving collectionInfo to Redis: %v", err)
	}

	return nil
}

func CollectionInfosPreload() {
	collectionInfos, _ := CollectionInfosFindAll()
	ctx := context.Background()
	for _, collectionInfo := range collectionInfos {
		if err := saveCollectionInfoToRedis(ctx, collectionInfo); err != nil {
			fmt.Printf("saveCollectionInfoToRedis: error saving collectionInfo to Redis: %v\n", err)
			continue
		}
	}
}

func GetCollectionInfoFromRedis(id int32) (*types.CollectionInfoItem, error) {
	// Construct the key for the collectionInfo item
	ctx := context.Background()
	key := fmt.Sprintf("collectionInfo:%d", id)

	// Retrieve the collectionInfo item from Redis
	collectionInfoJSON, err := database.RedisDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// collectionInfo item not found in Redis
			return nil, nil
		}
		return nil, fmt.Errorf("GetCollectionInfoFromRedis: error fetching collectionInfo item from Redis: %v", err)
	}

	// Deserialize the collectionInfo item from JSON
	var collectionInfo types.CollectionInfoItem
	if err := json.Unmarshal(collectionInfoJSON, &collectionInfo); err != nil {
		return nil, fmt.Errorf("GetCollectionInfoFromRedis: error unmarshaling collectionInfo item: %v", err)
	}

	return &collectionInfo, nil
}
