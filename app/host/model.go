package host

import "github.com/go-playground/validator/v10"

// 限制传参
var (
	validate = validator.New()
)

func (h *Host) Validate() error {
	return validate.Struct(h)
}

// 生成默认Host 防止空指针
func NewDefaultHost() *Host {
	return &Host{
		Resource:     &Resource{},
		DescribeHost: &DescribeHost{},
	}
}

type Vendor int

const (
	ALI_CLOUD Vendor = iota
	TX_CLOUD
)

type Host struct {
	*Resource
	*DescribeHost
	DescribeHash string `json:"describe_hash"`
	ResourceHash string `json:"resource_hash"`
}

type Resource struct {
	Id          string            `json:"id" validate:"required"`        //全局唯一id
	Vendor      Vendor            `json:"vendor" validate:"required"`    //厂商
	Region      string            `json:"regionId"validate:"required"`   //地域
	Zone        string            `json:"zone"`                          //可用区
	CreateAt    int64             `json:"create_at"validate:"required"`  //创建时间
	ExpireAt    int64             `json:"expire_at"`                     //过期时间
	Category    string            `json:"category"`                      //实例的类型
	Type        string            `json:"type"`                          //实例的规格
	InstanceId  string            `json:"instance_id"`                   //实例在云厂商内的ID
	Name        string            `json:"name"validate:"required"`       //实例名称
	Description string            `json:"description"`                   //实例描述
	Status      string            `json:"status"validate:"required"`     //状态
	Tags        map[string]string `json:"tags"`                          //标签
	UpdateAt    int64             `json:"update_at"`                     //更新时间
	SyncAt      int64             `json:"sync_at"`                       //同步时间
	SyncAccount string            `json:"syncAccount"`                   //同步的账户
	PublicIP    string            `json:"public_ip"`                     //公网ip
	PrivateIP   string            `json:"private_ip"validate:"required"` //内网ip
	PayType     string            `json:"pay_type"`                      //实例付费方式
}
type DescribeHost struct {
	CPU                     int    `json:"cpu"validate:"required"`     //核数
	Memory                  int    `json:"memory"validate:"required"`  //内存
	GPUAmount               int    `json:"gpu_amount"`                 //GPU数量
	GPUSpec                 string `json:"gpu_spec"`                   //GPU类型
	OSType                  string `json:"os_type"`                    //操作系统类型
	OSName                  string `json:"os_name"`                    //操作系统名称
	SerialNumber            string `json:"serial_number"`              //序列号
	ImageID                 string `json:"image_id"`                   //镜像ID
	InternetMaxBandwidthOut int    `json:"internet_max_bandwidth_out"` //公网出带宽最大值，单位Mbps
	InternetMaxBandwidthIn  int    `json:"internet_max_bandwidth_in"`  //公网入带宽最大值，单位Mbps
	KeyPairName             string `json:"key_pair_name"`              //密钥对名称
	SecurityGroups          string `json:"security_groups"`            //安全组采用逗号隔开
}

func NewSet() *Set {
	return &Set{
		Items: []*Host{},
	}
}
func (s *Set) ADD(itme *Host) {
	s.Items = append(s.Items, itme)
}

type Set struct {
	Total int64   `json:"total"`
	Items []*Host `json:"items"`
}
