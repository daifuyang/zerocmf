package data

import (
	"context"
	"time"
	"zerocmf/internal/biz"
)

type SmsRepo struct {
	data *Data
}

// 缓存验证码

func (repo *SmsRepo) CacheCode(smsCtx *biz.Sms) error {
	code := smsCtx.Code
	ctx := context.Background()
	err := repo.data.rdb.Set(ctx, smsCtx.Account+"_"+smsCtx.SceneCode, code, 3*time.Minute).Err()
	return err
}

func NewSmsRepo(data *Data) biz.SmsRepo {
	return &SmsRepo{
		data: data,
	}
}
