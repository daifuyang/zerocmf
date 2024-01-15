package data

import (
	"context"
	"fmt"
	"strings"
	"time"
	"zerocmf/internal/biz"

	"github.com/go-redis/redis/v8"
)

type PostRepo struct {
	data *Data
}

var (
	postCachePrefix = "cache:post:id:"
)

// 获取岗位列表
func (repo *PostRepo) Find(ctx context.Context, listQuery *biz.PostListQuery) (data interface{}, err error) {

	var total int64 = 0
	query := []string{"deleted_at is null"}
	queryArgs := make([]interface{}, 0)

	// 筛选条件
	if strings.TrimSpace(listQuery.PostName) != "" {
		query = append(query, "post_name like ?")
		queryArgs = append(queryArgs, "%"+listQuery.PostName+"%")
	}

	if strings.TrimSpace(listQuery.PostCode) != "" {
		query = append(query, "post_code like ?")
		queryArgs = append(queryArgs, "%"+listQuery.PostCode+"%")
	}

	// 状态
	if listQuery.Status != nil {
		query = append(query, "status = ?")
		queryArgs = append(queryArgs, *listQuery.Status)
	}

	queryStr := strings.Join(query, " and ")

	var posts []*biz.Post

	current, pageSize := biz.ParsePaginate(listQuery.Current, listQuery.PageSize)

	offset := (current - 1) * pageSize

	if pageSize == 0 {
		tx := repo.data.db.Where(queryStr, queryArgs...).Find(&posts)
		if tx.Error != nil {
			err = tx.Error
			return
		}
		data = posts
		return
	}

	tx := repo.data.db.Model(&biz.Post{}).Where("deleted_at is null").Count(&total)
	if tx.Error != nil {
		err = tx.Error
		return
	}

	tx = repo.data.db.Where(queryStr, queryArgs...).Offset(offset).Limit(pageSize).Order("list_order").Find(&posts)
	if tx.Error != nil {
		err = tx.Error
		return
	}

	data = &biz.Paginate{
		Total:    total,
		Current:  current,
		PageSize: pageSize,
		Data:     posts,
	}
	return

}

// 根据条件查询一条岗位
func (repo *PostRepo) First(query interface{}, args ...interface{}) (*biz.Post, error) {
	var Post *biz.Post
	tx := repo.data.db.Where(query, args...).First(&Post)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return Post, nil
}

// 根据id获取一条岗位
func (repo *PostRepo) FindOne(ctx context.Context, id int64) (*biz.Post, error) {
	var Post *biz.Post

	key := fmt.Sprintf("%s%v", postCachePrefix, id)
	err := repo.data.RGet(ctx, &Post, key)
	if err != redis.Nil {
		return Post, nil
	}

	tx := repo.data.db.Where("post_id = ? AND deleted_at is null", id).First(&Post)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 将结果存入redis
	err = repo.data.RSet(ctx, key, Post)
	if err != nil {
		return nil, err
	}

	return Post, nil
}

// 插入单个岗位
func (repo *PostRepo) Insert(ctx context.Context, post *biz.Post) (err error) {
	tx := repo.data.db.Create(&post)
	return tx.Error
}

// 更新单个岗位
func (repo *PostRepo) Update(ctx context.Context, post *biz.Post) (err error) {

	key := fmt.Sprintf("%s%v", postCachePrefix, post.PostID)
	repo.data.rdb.Del(ctx, key)

	tx := repo.data.db.Save(&post)
	return tx.Error
}

// 删除单个岗位
func (repo *PostRepo) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s%v", postCachePrefix, id)
	repo.data.rdb.Del(ctx, key)
	return repo.data.db.Model(&biz.Post{}).Where("post_id = ?", id).Update("deleted_at", time.Now()).Error
}

func NewPostRepo(data *Data) biz.PostRepo {
	return &PostRepo{data: data}
}
