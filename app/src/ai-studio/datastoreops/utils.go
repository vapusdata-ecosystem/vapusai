package dmstores

import (
	"context"

	utils "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/utils"
)

func (ds *DMStore) CacheFilter(ctx context.Context, action, key string, value ...string) (interface{}, error) {
	switch action {
	case utils.LIST:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	case utils.ADD:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.ADD", key, value[0]).Result()
	case utils.EXISTS:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	case utils.COUNT:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.CARD", key).Result()
	case utils.MADD:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.MADD", key, value).Result()
	case utils.DEL:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.DEL", key, value[0]).Result()
	default:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	}
}
