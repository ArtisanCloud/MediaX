package schema

type AgencyQueryBillLink struct {
	BillLink     string `json:"bill_link"`
	LiveBillLink string `json:"live_bill_link"`
}

type DouYinTaskBoxAgencyQueryBillLinkRes struct {
	ErrNo  int64               `json:"err_no"`
	ErrMsg string              `json:"err_msg"`
	LogId  string              `json:"log_id"`
	Data   AgencyQueryBillLink `json:"data"`
}
