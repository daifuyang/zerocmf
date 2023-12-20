# 快速入门

## 快速创建一个模块

- 创建路由

修改 `internal/server/api/v1/router.go` ,创建您新增模块的路由，推荐使用 restful 风格。示例：岗位管理

```go
// 模块化路由
 v1 := router.Group("/api/v1")
 {
    // 在用户管理下面添加岗位管理，方便以后微服务化
    user := v1.Group("/user")
    {
        // 部门管理
        user.GET("/post", service.NewPost(svcCtx).List)
        user.GET("/post/:id", service.NewPost(svcCtx).Show)
        user.POST("/post", service.NewPost(svcCtx).Add)
        user.POST("/post/:id", service.NewPost(svcCtx).Update)
        user.DELETE("/post/:id", service.NewPost(svcCtx).Delete)
    }
}
```

- 创建 service 服务层，示例：`internal/service/post.go`

```go
package service

import "github.com/gin-gonic/gin"

type post struct {
 *Context
}

func NewPost(c *Context) *post {
 return &post{
  Context: c,
 }
}

// 获取所有岗位
func (s *post) List(c *gin.Context) {

}

// 获取当个岗位
func (s *post) Show(c *gin.Context) {

}

// 新增单个岗位
func (s *post) Add(c *gin.Context) {

}

// 更新单个岗位
func (s *post) Update(c *gin.Context) {

}

// 【解耦】统一保存新增和更新逻辑
func (s *post) Save(c *gin.Context) {

}

// 删除单个岗位
func (s *post) Delete(c *gin.Context) {

}
```

- 创建 biz 业务逻辑层，示例：`internal/biz/post.go`

```go
package biz

import (
 "context"

 "gorm.io/gorm"
)

// SysPost 表示 sys_post 表的数据模型
type SysPost struct {
 PostID    int64      `gorm:"column:post_id;primaryKey;comment:岗位ID" json:"postId"`
 PostCode  string     `gorm:"column:post_code;size:64;unique;not null;comment:岗位编码" json:"postCode"`
 PostName  string     `gorm:"column:post_name;size:50;not null;comment:岗位名称" json:"postName"`
 ListOrder int        `gorm:"column:list_order;not null;comment:显示顺序" json:"listOrder"`
 Status    string     `gorm:"column:status;size:1;not null;default 1;comment:状态（1正常;0停用）" json:"status"`
 CreateId  int64      `gorm:"column:create_id;default:0;comment:创建者" json:"createId"`
 CreatedAt LocalTime  `gorm:"column:created_at;autoCreateTime;index;comment:创建时间" json:"createdAt"`
 UpdateId  int64      `gorm:"column:update_id;default:0;comment:更新者" json:"updateId"`
 UpdatedAt LocalTime  `gorm:"column:updated_at;autoUpdateTime;index;comment:更新时间" json:"updatedAt"`
 DeletedAt *LocalTime `gorm:"column:deleted_at;default:null;index;comment:删除时间" json:"deletedAt"`
 Remark    string     `gorm:"comment:备注;size:500;type:varchar(500)" json:"remark"`
}

// 列表筛选条件

type SysPostListQuery struct {
 PostCode string `form:"postCode"`
 PostName string `form:"postName"`
 Status   *int   `form:"status"`
 PaginateQuery
}

// TableName 指定表名
func (SysPost) TableName() string {
 return "sys_post"
}

// 数据库迁移
func (biz *SysPost) AutoMigrate(db *gorm.DB) error {
 err := db.AutoMigrate(&biz)
 if err != nil {
  return err
 }
 return nil
}

// 定义repo实体接口 （依赖倒置原则）
type PostRepo interface {
 Find(ctx context.Context, listQuery *SysPostListQuery) (interface{}, error) // 查看全部
 FindOne(ctx context.Context, id int64) (*SysPost, error)                    // 查询一条
 Insert(ctx context.Context, menu *SysPost) (err error)                      // 插入一条
 Update(ctx context.Context, menu *SysPost) (err error)                      // 更新一条
 Delete(ctx context.Context, id int64) error                                 // 删除一条
}

// 定义业务用例
type Postusecase struct {
 repo PostRepo
}

// 使用wire进行依赖注入
func NewPostusecase(post PostRepo) *Postusecase {
 return &Postusecase{
  repo: post,
 }
}

// 获取列表
func (biz *Postusecase) Find(ctx context.Context, listQuery *SysPostListQuery) (interface{}, error) {
 return biz.repo.Find(ctx, listQuery)
}

// 获取一条数据
func (biz *Postusecase) FindOne(ctx context.Context, id int64) (*SysPost, error) {
 return biz.repo.FindOne(ctx, id)
}

// 新增一条数据
func (biz *Postusecase) Insert(ctx context.Context, post *SysPost) error {
 return biz.repo.Insert(ctx, post)
}

// 更新一条数据
func (biz *Postusecase) Update(ctx context.Context, post *SysPost) error {
 return biz.repo.Update(ctx, post)
}

// 删除一条数据
func (biz *Postusecase) Delete(ctx context.Context, id int64) (*SysPost, error) {
 one, err := biz.repo.FindOne(ctx, id)
 if err != nil {
  return nil, err
 }
 err = biz.repo.Delete(ctx, id)
 if err != nil {
  return nil, err
 }
 return one, nil
}
```

- 修改 `internal/biz/migrate.go` 进行数据库迁移，自动创建表字段

```go
package biz

import (
 "zerocmf/configs"

 "gorm.io/gorm"
)

// 数据库迁移
func AutoMigrate(db *gorm.DB, config *configs.Config) error {
 err = new(SysPost).AutoMigrate(db)
 if err != nil {
  return err
 }
 return nil
}

```

- 修改 `internal/biz/biz.go` 绑定 usecase 到 wire

```go
package biz

import "github.com/google/wire"

// ...代表其他用例，这里省略

var ProviderSet = wire.NewSet(NewDemoUsecase, ... , NewPostusecase) // 加入NewPostusecase到wire中

```

- 新增 `internal/biz/post.go` 定义业务数据访问层

```go
package data

import (
 "context"
 "zerocmf/internal/biz"
)

type PostRepo struct {
 data *Data
}

// Find implements biz.PostRepo.
func (*PostRepo) Find(ctx context.Context, listQuery *biz.SysPostListQuery) (interface{}, error) {
 panic("unimplemented")
}

// FindOne implements biz.PostRepo.
func (*PostRepo) FindOne(ctx context.Context, id int64) (*biz.SysPost, error) {
 panic("unimplemented")
}

// Insert implements biz.PostRepo.
func (*PostRepo) Insert(ctx context.Context, menu *biz.SysPost) (err error) {
 panic("unimplemented")
}

// Update implements biz.PostRepo.
func (*PostRepo) Update(ctx context.Context, menu *biz.SysPost) (err error) {
 panic("unimplemented")
}

// Delete implements biz.PostRepo.
func (*PostRepo) Delete(ctx context.Context, id int64) error {
 panic("unimplemented")
}

func NewPostRepo(data *Data) biz.PostRepo {
 return &PostRepo{data: data}
}

```

- 修改 `internal/data/data.go` 绑定 repo 到 wire

```go
package data

import (
    //...省略其他代码
   "github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewData, ..., NewPostRepo) // 加入NewPostRepo到wire中

//...省略其他代码

```

- 修改 `internal/service/service.go` 绑定 usecase 到 service

```go
package service

import (
 "zerocmf/internal/biz"
 "github.com/google/wire"
 //...省略其他代码
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewContext)

// 用来组装上下文依赖
type Context struct {
//...省略其他代码
demouc *biz.Demousecase
postuc *biz.Postusecase
}

func NewContext(..., demouc *biz.Demousecase ,postpc *biz.Postusecase) *Context {

//...省略其他代码
return &Context{
    demouc: demouc,
    postuc: postpc,
  }
}
```

终端生成wire依赖 

```shell
cd cmd/server
wire .

```
