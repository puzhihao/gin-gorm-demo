package host

import "context"

type Service interface {
	// CreateHost 录入主机信息
	CreateHost(ctx context.Context, host *Host) (*Host, error)
	// QueryHost 查询主机列表
	QueryHost(ctx context.Context, req *QueryHostRequest) (*Set, error)
	// DescribeHost 查看主机详情
	DescribeHost(ctx context.Context, req *DescribeHostRequest) (*Host, error)
	// UpdateHost 修改主机信息
	UpdateHost(ctx context.Context, req *UpdateHostRequest) (*Host, error)
	// DeleteHost 删除主机
	DeleteHost(ctx context.Context, req *DeleteHostRequest) (*Host, error)
}

type QueryHostRequest struct {
	PageSize   int
	PageNumber int
	Keywords   string
}

type DescribeHostRequest struct {
	Id string
}

const (
	PATCH = 0
	PUT   = 1
)

type UpdateMode int
type UpdateHostRequest struct {
	*Resource
	*DescribeHost
	UpdateMode
}
type DeleteHostRequest struct {
	Id string
}

func NewPatchUpdateHostRequest() *UpdateHostRequest {
	return &UpdateHostRequest{
		UpdateMode:   PATCH,
		Resource:     &Resource{},
		DescribeHost: &DescribeHost{},
	}
}

func NewPutUpdateHostRequest() *UpdateHostRequest {
	return &UpdateHostRequest{
		UpdateMode:   PUT,
		Resource:     &Resource{},
		DescribeHost: &DescribeHost{},
	}
}
