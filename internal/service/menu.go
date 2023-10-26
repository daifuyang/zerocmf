package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"zerocmf/internal/biz"
	"zerocmf/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type menu struct {
	*Context
}

func NewMenu(c *Context) *menu {
	return &menu{
		Context: c,
	}
}

// 解析配置文件
func mustLoad(configFile string, menu *[]*biz.SysMenu) {
	// 解析配置项
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic("读取菜单文件失败：" + err.Error())
	}

	if err := json.Unmarshal(data, &menu); err != nil {
		panic("解析菜单文件失败：" + err.Error())
	}
}

// 配置文件导入菜单
func (s *menu) ImportMenu() {

	var menu []*biz.SysMenu

	currentDir, err := os.Getwd()
	if err != nil {
		panic("无法获取当前工作目录!")
	}

	mustLoad(currentDir+"/static/menu.json", &menu)

	var parentId int64 = 0

	ctx := context.Background()

	s.recursionImportMenu(ctx, menu, parentId, "")

}

// 递归导入菜单
func (s *menu) recursionImportMenu(ctx context.Context, menu []*biz.SysMenu, parentId int64, perms string) {
	for k, v := range menu {

		newPerms := v.Perms
		if perms != "" {
			newPerms = perms + "." + v.Perms
		}

		// 查询菜单是否存在
		one, err := s.mc.FindOneByMenuName(ctx, v.MenuName)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic("导入菜单失败：" + err.Error())
		}

		localOne := &biz.SysMenu{
			MenuName: v.MenuName,
			ParentID: parentId,
			Path:     v.Path,
			OrderNum: k + 1,
			MenuType: v.MenuType,
			Perms:    newPerms,
			CreateId: 1,
		}

		// 只能一条一条的插入
		if one != nil {
			localOne.MenuID = one.MenuID
			localOne.Status = one.Status
			localOne.OrderNum = one.OrderNum
			localOne.Visible = one.Visible
			localOne.Icon = one.Icon
			localOne.CreatedAt = one.CreatedAt
			localOne.Remark = one.Remark
			s.mc.Update(ctx, one)
		} else {
			s.mc.Insert(ctx, localOne)
		}
		nextParentId := localOne.MenuID
		if v.Children != nil {
			s.recursionImportMenu(ctx, v.Children, nextParentId, newPerms)
		}
	}
}

// 递归显示树菜单
func recursionMenu(menu []*biz.SysMenu, parentId int64) []*biz.SysMenu {
	var result = make([]*biz.SysMenu, 0)
	for _, v := range menu {
		if v.ParentID == parentId {
			children := recursionMenu(menu, v.MenuID)
			if len(children) > 0 {
				v.Children = children
			}
			result = append(result, v)
		}
	}
	return result
	// return menu
}

// 查看全部菜单树
func (s *menu) Tree(c *gin.Context) {
	sysMenu, err := s.mc.Find(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	menus := recursionMenu(sysMenu, 0)
	fmt.Println("menus", menus)
	response.Success(c, "获取成功！", menus)
}

// 添加一条菜单
func (s *menu) Add(c *gin.Context) {
	var req biz.SysMenu
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	err := s.mc.Insert(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "添加成功！", req)

}

// 查看一条菜单
func (s *menu) Show(c *gin.Context) {

}

// 更新一条菜单
func (s *menu) Update(c *gin.Context) {

}

// 新增和跟新统一逻辑
func (s *menu) Save(c *gin.Context) {

}

// 删除一条菜单
func (s *menu) Delete(c *gin.Context) {

}

// 批量删除菜单
func (s *menu) DeleteBatch(c *gin.Context) {
}
