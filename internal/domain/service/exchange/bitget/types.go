package bitget

// BaseResponse 基础响应结构
type BaseResponse[T any] struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    T      `json:"data"`
}

// CustomerListData 客户列表数据
type CustomerListData struct {
	CustomerInfoList []CustomerInfo `json:"customerInfoList"`
}

type CustomerInfo struct {
	Uid          string `json:"uid"`
	RegisterTime string `json:"registerTime"`
}

type CustomerTradeVolumeData struct {
	CustomerVolumeList []CustomerVolume `json:"customerVolumeList"`
}

type CustomerVolume struct {
	Uid    string `json:"uid"`
	Volume string `json:"volumn"`
	Time   string `json:"time"`
}

// IsSuccess 检查响应是否成功
func (r BaseResponse[T]) IsSuccess() bool {
	return r.Code == "00000"
}

// GetError 获取错误信息
func (r BaseResponse[T]) GetError() string {
	if r.IsSuccess() {
		return ""
	}
	return r.Message
}
