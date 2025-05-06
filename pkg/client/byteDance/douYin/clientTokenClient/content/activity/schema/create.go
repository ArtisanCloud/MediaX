package schema

// LiveWatchEventInfo 直播观看事件信息
type LiveWatchEventInfo struct {
	GameId        string `json:"game_id"`        // 游戏ID
	WatchDuration int64  `json:"watch_duration"` // 观看时长
	AwemeId       string `json:"aweme_id"`       // 抖音ID
}

// ShortVideoCommentEventInfo 短视频评论事件信息
type ShortVideoCommentEventInfo struct {
	ItemId   string `json:"item_id"`   // 视频ID
	VideoUrl string `json:"video_url"` // 视频URL
}

// AccountFollowEventInfo 账号关注事件信息
type AccountFollowEventInfo struct {
	AwemeId string `json:"aweme_id"` // 抖音ID
}

// LiveCommentEventInfo 直播评论事件信息
type LiveCommentEventInfo struct {
	AwemeId string `json:"aweme_id"` // 抖音ID
}

// ShortVideoFinishPlayingEventInfo 短视频播放完成事件信息
type ShortVideoFinishPlayingEventInfo struct {
	ItemId   string `json:"item_id"`   // 视频ID
	VideoUrl string `json:"video_url"` // 视频URL
	GameId   string `json:"game_id"`   // 游戏ID
	AnchorId string `json:"anchor_id"` // 主播ID
}

// ShortVideoDiggEventInfo 短视频点赞事件信息
type ShortVideoDiggEventInfo struct {
	ItemId   string `json:"item_id"`   // 视频ID
	VideoUrl string `json:"video_url"` // 视频URL
}

// PostingVideEventInfo 发布视频事件信息
type PostingVideEventInfo struct {
	PlayVideoStage []int64 `json:"play_video_stage"` // 视频播放阶段
	NumStage       []int64 `json:"num_stage"`        // 数量阶段
	DiggStage      []int64 `json:"digg_stage"`       // 点赞阶段
	CommentStage   []int64 `json:"comment_stage"`    // 评论阶段
}

// LiveShareEventInfo 直播分享事件信息
type LiveShareEventInfo struct {
	AwemeId string `json:"aweme_id"` // 抖音ID
}

// ShortVideoShareEventInfo 短视频分享事件信息
type ShortVideoShareEventInfo struct {
	VideoUrl string `json:"video_url"` // 视频URL
	ItemId   string `json:"item_id"`   // 视频ID
}

// ShortVideoCollectionEventInfo 短视频收藏事件信息
type ShortVideoCollectionEventInfo struct {
	VideoUrl string `json:"video_url"` // 视频URL
	ItemId   string `json:"item_id"`   // 视频ID
}

// LiveDiggEventInfo 直播点赞事件信息
type LiveDiggEventInfo struct {
	AwemeId string `json:"aweme_id"` // 抖音ID
}

// TaskEventInfo 任务事件信息
type TaskEventInfo struct {
	LiveWatchEventInfo               LiveWatchEventInfo               `json:"live_watch_event_info"`                 // 直播观看事件信息
	ShortVideoCommentEventInfo       ShortVideoCommentEventInfo       `json:"short_video_comment_event_info"`        // 短视频评论事件信息
	AccountFollowEventInfo           AccountFollowEventInfo           `json:"account_follow_event_info"`             // 账号关注事件信息
	LiveCommentEventInfo             LiveCommentEventInfo             `json:"live_comment_event_info"`               // 直播评论事件信息
	ShortVideoFinishPlayingEventInfo ShortVideoFinishPlayingEventInfo `json:"short_video_finish_playing_event_info"` // 短视频播放完成事件信息
	ShortVideoDiggEventInfo          ShortVideoDiggEventInfo          `json:"short_video_digg_event_info"`           // 短视频点赞事件信息
	PostingVideEventInfo             PostingVideEventInfo             `json:"posting_vide_event_info"`               // 发布视频事件信息
	LiveShareEventInfo               LiveShareEventInfo               `json:"live_share_event_info"`                 // 直播分享事件信息
	ShortVideoShareEventInfo         ShortVideoShareEventInfo         `json:"short_video_share_event_info"`          // 短视频分享事件信息
	ShortVideoCollectionEventInfo    ShortVideoCollectionEventInfo    `json:"short_video_collection_event_info"`     // 短视频收藏事件信息
	LiveDiggEventInfo                LiveDiggEventInfo                `json:"live_digg_event_info"`                  // 直播点赞事件信息
}

// BusinessTaskInfo 业务任务信息
type BusinessTaskInfo struct {
	TaskName      string        `json:"task_name"`       // 任务名称
	TaskEventInfo TaskEventInfo `json:"task_event_info"` // 任务事件信息
	StartTime     int64         `json:"start_time"`      // 开始时间
	EndTime       int64         `json:"end_time"`        // 结束时间
	TaskId        int64         `json:"task_id"`         // 任务ID
}

// DouYinContentActivityCreateReq 创建活动请求参数
type DouYinContentActivityCreateReq struct {
	StartTime                  int64              `json:"start_time"`                     // 开始时间
	EndTime                    int64              `json:"end_time"`                       // 结束时间
	CreateBusinessTaskInfoList []BusinessTaskInfo `json:"create_business_task_info_list"` // 创建业务任务信息列表
	NotBcCheck                 bool               `json:"not_bc_check"`                   // 是否不进行BC检查
	ActivityName               string             `json:"activity_name"`                  // 活动名称
}

// DouYinContentActivityCreateRes 创建活动响应结构
type DouYinContentActivityCreateRes struct {
	ActivityId         int64   `json:"activity_id"`           // 活动ID
	BusinessTaskIdList []int64 `json:"business_task_id_list"` // 业务任务ID列表
}
