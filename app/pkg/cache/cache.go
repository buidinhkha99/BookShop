package cache

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ctx = context.Background()

func InsertData(key, data string) error {
	address := viper.GetString("redis.address")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	err := rdb.Set(ctx, key, data, 0).Err()
	if err != nil {
		return nil
	} else {
		return err
	}

}
func ServeJQueryWithRemoteCache(w http.ResponseWriter, key string) string {
	address := viper.GetString("redis.address")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	val2, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Info("No data in remote cache")
		return ""
	} else if err != nil {
		log.Info("Data has not been created in remote cache")
		return ""
	} else {
		log.Info("There are data in remote cache")
		return val2
	}
}
