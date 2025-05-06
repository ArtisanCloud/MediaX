package schema

// DouYinContentActivityQueryUserStatusReq 查询用户状态请求参数
type DouYinContentActivityQueryUserStatusReq struct {
	ActivityId   int64   `json:"activity_id"`    // 活动ID
	TargetOpenId string  `json:"target_open_id"` // 目标用户OpenId
	AlliedId     string  `json:"allied_id"`      // 联盟ID
	TaskIdList   []int64 `json:"task_id_list"`   // 任务ID列表
}

// TaskCompleteStatusMap 任务完成状态映射
type TaskCompleteStatusMap map[string]bool

// PostingNumStageMap 发布数量阶段映射
type PostingNumStageMap map[string]bool

// DiggStageMap 点赞阶段映射
type DiggStageMap map[string]bool

// CommentStageMap 评论阶段映射
type CommentStageMap map[string]bool

// PlayVideoStageMap 视频播放阶段映射
type PlayVideoStageMap map[string]bool

// MapCompleteInfo 映射完成信息
type MapCompleteInfo struct {
	PostingNumStageMap PostingNumStageMap `json:"posting_num_stage_map"` // 发布数量阶段映射
	DiggStageMap       DiggStageMap       `json:"digg_stage_map"`        // 点赞阶段映射
	CommentStageMap    CommentStageMap    `json:"comment_stage_map"`     // 评论阶段映射
	PlayVideoStageMap  PlayVideoStageMap  `json:"play_video_stage_map"`  // 视频播放阶段映射
}

// PostingVideoBindCompleteInfoMap 发布视频绑定完成信息映射
type PostingVideoBindCompleteInfoMap map[string]MapCompleteInfo

// DouYinContentActivityQueryUserStatusRes 查询用户状态响应结构
type DouYinContentActivityQueryUserStatusRes struct {
	TaskCompleteStatusMap           TaskCompleteStatusMap           `json:"task_complete_status_map"`              // 任务完成状态映射
	PostingVideoBindCompleteInfoMap PostingVideoBindCompleteInfoMap `json:"posting_video_bind_compelete_info_map"` // 发布视频绑定完成信息映射
}
