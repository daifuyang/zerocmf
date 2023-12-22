package data

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"zerocmf/internal/biz"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type roleRepo struct {
	data *Data
}

func NewRoleRepo(data *Data) biz.RoleRepo {
	return &roleRepo{data: data}
}

var (
	roleCachePrefix     = "cache:role:id:"
	rolePermCachePrefix = "cache:rolePerm:id:"
)

// 查询全部
func (repo *roleRepo) Find(ctx context.Context, listQuery *biz.SysRoleListQuery) (data interface{}, err error) {

	var total int64 = 0
	query := []string{"deleted_at is null"}
	queryArgs := make([]interface{}, 0)

	if strings.TrimSpace(listQuery.RoleName) != "" {
		query = append(query, "role_name like ?")
		queryArgs = append(queryArgs, "%"+listQuery.RoleName+"%")
	}

	if listQuery.Status != nil {
		query = append(query, "status = ?")
		queryArgs = append(queryArgs, *listQuery.Status)
	}

	queryStr := strings.Join(query, " and ")

	var roles []*biz.SysRole

	current, pageSize := biz.ParsePaginate(listQuery.Current, listQuery.PageSize)

	offset := (current - 1) * pageSize

	if pageSize == 0 {
		tx := repo.data.db.Where(queryStr, queryArgs...).Find(&roles)
		if tx.Error != nil {
			err = tx.Error
			return
		}
		data = roles
		return
	}

	tx := repo.data.db.Model(&biz.SysRole{}).Where("deleted_at is null").Count(&total)
	if tx.Error != nil {
		err = tx.Error
		return
	}

	tx = repo.data.db.Where(queryStr, queryArgs...).Offset(offset).Limit(pageSize).Find(&roles)
	if tx.Error != nil {
		err = tx.Error
		return
	}

	data = &biz.Paginate{
		Total:    total,
		Current:  current,
		PageSize: pageSize,
		Data:     roles,
	}
	return
}

// 根据角色id获取单个角色
func (repo *roleRepo) FindOne(ctx context.Context, id int64) (*biz.SysRole, error) {
	var sysRole *biz.SysRole
	key := fmt.Sprintf("%s%v", roleCachePrefix, id)

	err := repo.data.RGet(ctx, &sysRole, key)
	if err != redis.Nil {
		return sysRole, nil
	}

	tx := repo.data.db.Where("role_id = ? AND deleted_at is null", id).First(&sysRole)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// 将结果存入redis
	err = repo.data.RSet(ctx, key, sysRole)
	if err != nil {
		return nil, err
	}

	return sysRole, nil
}

// 根基角色id获取授权通过的权限列表
func (repo *roleRepo) FindPermissions(ctx context.Context, id int64) ([]*int, error) {

	key := fmt.Sprintf("%s%v", rolePermCachePrefix, id)

	var permissionIds []*int

	err := repo.data.RGet(ctx, &permissionIds, key)
	if err != redis.Nil {
		return permissionIds, nil
	}

	e := repo.data.e
	roleId := strconv.FormatInt(id, 10)
	permissions := e.GetPermissionsForUser(roleId)
	for _, p := range permissions {
		menuId, err := strconv.Atoi(p[1])
		if err != nil {
			return nil, err
		}
		permissionIds = append(permissionIds, &menuId) // 在 Casbin 中权限的位置可能不同，根据你的策略文件进行调整
	}

	// 将结果存入redis
	err = repo.data.RSet(ctx, key, permissionIds)
	if err != nil {
		return nil, err
	}

	return permissionIds, nil
}

// 插入一条数据
func (repo *roleRepo) Insert(ctx context.Context, role *biz.SysRole) (err error) {

	e := repo.data.e
	db := repo.data.db
	return db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&role).Error; err != nil {
			return err
		}

		roleId := strconv.FormatInt(role.RoleID, 10)
		menuIds := make([]string, len(role.MenuIds))
		for i, v := range role.MenuIds {
			menuIds[i] = strconv.Itoa(*v)
		}
		_, err = e.AddPermissionForUser(roleId, menuIds...)
		if err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

// 辅助函数：检查切片中是否包含某个值
func contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

// 更新一条数据
func (repo *roleRepo) Update(ctx context.Context, role *biz.SysRole) (err error) {

	roleId := role.RoleID
	roleIdStr := strconv.FormatInt(roleId, 10)

	roleKey := fmt.Sprintf("%s%v", roleCachePrefix, roleId)
	repo.data.rdb.Del(ctx, roleKey)

	permKey := fmt.Sprintf("%s%v", rolePermCachePrefix, roleId)
	repo.data.rdb.Del(ctx, permKey)

	e := repo.data.e
	db := repo.data.db
	return db.Transaction(func(tx *gorm.DB) error {
		// redis重置

		if err := tx.Save(&role).Error; err != nil {
			return err
		}

		// 传入的新规则
		menuIds := role.MenuIds

		// 转成字符串
		newAccess := make([]string, len(menuIds))
		for i, v := range menuIds {
			newAccess[i] = strconv.Itoa(*v)
		}

		//获取原来的数据

		permissions := e.GetPermissionsForUser(roleIdStr)

		existAccess := make([]string, len(permissions))

		for i, permission := range permissions {
			existAccess[i] = permission[1]
		}

		// 如果不存在，则全部添加（首次添加）
		if len(existAccess) == 0 {
			for _, access := range newAccess {
				e.AddPermissionForUser(roleIdStr, access)
			}
		} else if len(newAccess) == 0 {
			// 清空
			for _, access := range existAccess {
				e.DeletePermissionForUser(roleIdStr, access)
			}
		} else {
			// 先找出需要删除的值
			for _, exist := range existAccess {
				if !contains(newAccess, exist) {
					e.DeletePermissionForUser(roleIdStr, exist)
				}
			}

			// 再找出需要添加的值
			for _, new := range newAccess {
				if !contains(existAccess, new) {
					e.AddPermissionForUser(roleIdStr, new)
				}
			}
		}

		return nil
	})
}

// 删除一条数据
func (repo *roleRepo) Delete(ctx context.Context, id int64) error {

	roleKey := fmt.Sprintf("%s%v", roleCachePrefix, id)
	repo.data.rdb.Del(ctx, roleKey)

	permKey := fmt.Sprintf("%s%v", rolePermCachePrefix, id)
	repo.data.rdb.Del(ctx, permKey)

	return repo.data.db.Model(&biz.SysRole{}).Where("role_id = ?", id).Update("deleted_at", time.Now()).Error
}
