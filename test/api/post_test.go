package api

import (
	"fmt"
	"net/http"
	"testing"
)

func TestPost(t *testing.T) {

	prefix := "localhost:8080/api/v1/user"

	req, _ := http.NewRequest("GET", prefix+"/post", nil)

	fmt.Println("req", req.Body)

	// 获取岗位列表

	// 新增单个岗位

	// 修改单个岗位

	// 获取单个岗位

}
