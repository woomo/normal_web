package util

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"normal_web/database"
	"normal_web/util/logger"
	"time"
)

const (
	TOKEN_PREFIX = "dual_token_"
	TOKEN_EXPIRE = 7 * 24 * time.Hour //一次登录7天有效
)

// 把<refreshToken, authToken>写入redis
func SetToken(refreshToken, authToken string) {
	client := database.GetRedisClient()
	cxt := context.Background()
	if err := client.Set(cxt, TOKEN_PREFIX+refreshToken, authToken, TOKEN_EXPIRE).Err(); err != nil { //7天之后就拿不到authToken了
		logger.Error(fmt.Sprintf("write token pair(%s, %s) to redis failed: %s", refreshToken, authToken, err))
	}
}

// 根据refreshToken获取authToken
func GetToken(refreshToken string) (authToken string) {
	client := database.GetRedisClient()
	cxt := context.Background()
	var err error
	if authToken, err = client.Get(cxt, TOKEN_PREFIX+refreshToken).Result(); err != nil {
		if err != redis.Nil {
			logger.Error(fmt.Sprintf("get auth token of refresh token %s failed: %s", refreshToken, err))
		}
	}
	return
}
