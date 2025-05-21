package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangPromoteCreativityCreateNoteReq 创建笔记创意请求
type JuGuangPromoteCreativityCreateNoteReq struct {
	AdvertiserID             int64    `json:"advertiser_id"`                        // 广告主ID
	UnitID                   int64    `json:"unit_id"`                              // 单元ID
	CreativityName           string   `json:"creativity_name"`                      // 创意名称
	NoteID                   string   `json:"note_id"`                              // 笔记ID
	ClickURLs                []string `json:"click_urls,omitempty"`                 // 点击链接
	ExpoURLs                 []string `json:"expo_urls,omitempty"`                  // 曝光链接
	CustomMask               int      `json:"custom_mask,omitempty"`                // 自选封面
	CustomTitle              int      `json:"custom_title,omitempty"`               // 是否自提标题
	TitleFills               []string `json:"title_fills,omitempty"`                // 自提标题
	MaskGen                  int      `json:"mask_gen,omitempty"`                   // 是否开启自动优化封面
	TitleGen                 int      `json:"title_gen,omitempty"`                  // 是否开启自动优化标题
	ConversionType           int      `json:"conversion_type"`                      // 组件类型
	JumpURL                  string   `json:"jump_url,omitempty"`                   // 落地页/外链url
	LandingPageType          int      `json:"landing_page_type,omitempty"`          // 落地页链接类型
	BarContent               string   `json:"bar_content,omitempty"`                // 按钮文案内容
	ConversionComponentTypes []int64  `json:"conversion_component_types,omitempty"` // 组件位置
	Comment                  string   `json:"comment,omitempty"`                    // 置顶评论文案
	AppCompIcon              string   `json:"app_comp_icon,omitempty"`              // 唤端下，商品主图的地址
	FallBackJumpURL          string   `json:"fall_back_jump_url,omitempty"`         // 唤端下，兜底链接
	QualInfo                 QualInfo `json:"qual_info"`                            // 资质信息
	MiniProgramPath          string   `json:"mini_program_path,omitempty"`          // 小程序组件链接
	PrimaryTitle             string   `json:"primary_title,omitempty"`              // 主标题
	IOSDownloadLink          string   `json:"ios_download_link,omitempty"`          // ios兜底下载链接
	AndroidDownloadLink      string   `json:"android_download_link,omitempty"`      // 安卓兜底下载链接
}

// QualInfo 资质信息
type QualInfo struct {
	ApplyID           string `json:"apply_id"`                       // 申请id
	ProductQualIDList []int  `json:"product_qual_id_list,omitempty"` // 产品资质id
	BrandQualIDList   []int  `json:"brand_qual_id_list,omitempty"`   // 品牌资质id
}

// JuGuangPromoteCreativityCreateNoteRes 创建笔记创意响应
type JuGuangPromoteCreativityCreateNoteRes struct {
	response.RedBookAccessTokenRes
	Data CreateNoteData `json:"data"` // 返回数据
}

// Data 返回数据
type CreateNoteData struct {
	CreativityID int64 `json:"creativity_id"` // 创意ID
}
