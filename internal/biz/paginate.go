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

func ParsePaginate(current *int, pageSize *int) (int, int) {
	// 默认值

	_current := 1
	_pageSize := 10

	if current != nil {
		_current = *current
	}

	if pageSize != nil {
		_pageSize = *pageSize
	}

	return _current, _pageSize

}
