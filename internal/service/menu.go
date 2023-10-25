package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"zerocmf/internal/biz"

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
func mustLoad(configFile string, menu *[]biz.SysMenu) {
	// 解析配置项
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic("读取菜单文件失败：" + err.Error())
	}

	if err := json.Unmarshal(data, &menu); err != nil {
		panic("解析菜单文件失败：" + err.Error())
	}
}

func (s *menu) ImportMenu() {

	var menu []biz.SysMenu

	currentDir, err := os.Getwd()
	if err != nil {
		panic("无法获取当前工作目录!")
	}

	mustLoad(currentDir+"/static/menu.json", &menu)

	var parentId int64 = 0

	ctx := context.Background()

	s.recursionMenu(ctx, menu, parentId, "")

}

// 递归菜单
func (s *menu) recursionMenu(ctx context.Context, menu []biz.SysMenu, parentId int64, perms string) {
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
			s.recursionMenu(ctx, v.Children, nextParentId, newPerms)
		}
	}
}
