package biz

type PaginateQuery struct {
	Current  *int `form:"current"`
	PageSize *int `form:"pageSize"`
}

type Paginate struct {
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
	Total    int64       `json:"total"`
	Data     interface{} `json:"data"`
}
