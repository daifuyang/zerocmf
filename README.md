# 配置

## 环境变量

# 快速上手

## 创建领域对象

```
type testRepo struct {
	data *Data
}

// NewTestRepo .
func NewTestRepo(data *Data) biz.TestRepo {
	return &testRepo{
		data: data,
	}
}
```
