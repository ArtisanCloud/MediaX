package schema

// QueryViolateTalent 数据结构体
type QueryViolateTalent struct {
	ViolateTalentsURL string `json:"ViolateTalentsURL"`
}

// DouYinTaskBoxQueryViolateTalentListRes 查询违规达人列表响应结果
type DouYinTaskBoxQueryViolateTalentListRes struct {
	LogId  string             `json:"log_id"`
	Data   QueryViolateTalent `json:"data"`
	ErrMsg string             `json:"err_msg"`
	ErrNo  int                `json:"err_no"`
}
