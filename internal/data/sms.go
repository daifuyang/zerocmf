package data

import (
	"context"
	"errors"
	"strconv"
	"time"
	"zerocmf/internal/biz"

	"github.com/go-redis/redis/v8"
)

type SmsRepo struct {
	data *Data
}

// 缓存短信验证码

func (repo *SmsRepo) CacheCode(smsCtx *biz.Sms) error {

	// 获取是否存在已发送的短信验证码
	config := repo.data.config
	code := smsCtx.Code
	ctx := context.Background()

	// 使用客户端IP地址作为键，设置每小时最多发送5次的计数器，并设置过期时间
	countKey := smsCtx.Account + "_" + smsCtx.ClientIP

	// 发送频率
	key := smsCtx.Account + "_" + smsCtx.SceneCode
	_, err := repo.data.rdb.Get(ctx, key).Result()
	if err != redis.Nil {
		return errors.New("发送过于频繁，请稍后再试！")
	}

	remaining, err := repo.data.rdb.Incr(ctx, countKey).Result()
	if err != nil {
		return err
	}

	if remaining > 5 {
		return errors.New("发送频率超限，请稍后再试！")
	}

	err = repo.data.rdb.Set(ctx, key, code, time.Duration(config.Sms.ExpiresTime)*time.Minute).Err()
	return err
}

// 验证短信验证码

func (repo *SmsRepo) VerifyCode(smsCtx *biz.Sms) error {
	code := smsCtx.Code // 用户输入的验证码
	ctx := context.Background()
	key := smsCtx.Account + "_" + smsCtx.SceneCode
	storeCode, err := repo.data.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return errors.New("验证码不存在或已失效！")
	}
	if storeCode != strconv.Itoa(int(code)) {
		return errors.New("验证码错误！")
	}
	// 验证成功后清除消费验证码
	return repo.data.rdb.Del(ctx, key).Err()
}

func NewSmsRepo(data *Data) biz.SmsRepo {
	return &SmsRepo{
		data: data,
	}
}
