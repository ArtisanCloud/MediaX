package schema

import (
	"time"
)

// Video 结构体
type Video struct {
	Kind                 string                  `json:"kind"`
	Etag                 string                  `json:"etag"`
	ID                   string                  `json:"id"`
	Snippet              Snippet                 `json:"snippet"`
	ContentDetails       ContentDetails          `json:"contentDetails"`
	Status               Status                  `json:"status"`
	Statistics           Statistics              `json:"statistics"`
	PaidProductPlacement PaidProductPlacement    `json:"paidProductPlacementDetails"`
	Player               Player                  `json:"player"`
	TopicDetails         TopicDetails            `json:"topicDetails"`
	RecordingDetails     RecordingDetails        `json:"recordingDetails"`
	FileDetails          FileDetails             `json:"fileDetails"`
	ProcessingDetails    ProcessingDetails       `json:"processingDetails"`
	Suggestions          Suggestions             `json:"suggestions"`
	LiveStreamingDetails LiveStreamingDetails    `json:"liveStreamingDetails"`
	Localizations        map[string]Localization `json:"localizations"`
}

// Snippet 结构体
type Snippet struct {
	PublishedAt      time.Time            `json:"publishedAt"`
	ChannelID        string               `json:"channelId"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	Thumbnails       map[string]Thumbnail `json:"thumbnails"`
	ChannelTitle     string               `json:"channelTitle"`
	Tags             []string             `json:"tags"`
	CategoryID       string               `json:"categoryId"`
	LiveBroadcast    string               `json:"liveBroadcastContent"`
	DefaultLanguage  string               `json:"defaultLanguage"`
	Localized        Localization         `json:"localized"`
	DefaultAudioLang string               `json:"defaultAudioLanguage"`
}

// Thumbnail 结构体
type Thumbnail struct {
	URL    string `json:"url"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
}

// Localization 结构体
type Localization struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// ContentDetails 结构体
type ContentDetails struct {
	Duration           string            `json:"duration"`
	Dimension          string            `json:"dimension"`
	Definition         string            `json:"definition"`
	Caption            string            `json:"caption"`
	LicensedContent    bool              `json:"licensedContent"`
	RegionRestriction  RegionRestriction `json:"regionRestriction"`
	ContentRating      ContentRating     `json:"contentRating"`
	Projection         string            `json:"projection"`
	HasCustomThumbnail bool              `json:"hasCustomThumbnail"`
}

// RegionRestriction 结构体
type RegionRestriction struct {
	Allowed []string `json:"allowed"`
	Blocked []string `json:"blocked"`
}

// ContentRating 结构体
type ContentRating struct {
	ACBRating  string `json:"acbRating"`
	MPAARating string `json:"mpaaRating"`
	YT_Rating  string `json:"ytRating"`
}

// Status 结构体
type Status struct {
	UploadStatus            string    `json:"uploadStatus"`
	FailureReason           string    `json:"failureReason"`
	RejectionReason         string    `json:"rejectionReason"`
	PrivacyStatus           string    `json:"privacyStatus"`
	PublishAt               time.Time `json:"publishAt"`
	License                 string    `json:"license"`
	Embeddable              bool      `json:"embeddable"`
	PublicStatsViewable     bool      `json:"publicStatsViewable"`
	MadeForKids             bool      `json:"madeForKids"`
	SelfDeclaredMadeForKids bool      `json:"selfDeclaredMadeForKids"`
	ContainsSyntheticMedia  bool      `json:"containsSyntheticMedia"`
}

// Statistics 结构体
type Statistics struct {
	ViewCount     string `json:"viewCount"`
	LikeCount     string `json:"likeCount"`
	DislikeCount  string `json:"dislikeCount"`
	FavoriteCount string `json:"favoriteCount"`
	CommentCount  string `json:"commentCount"`
}

// PaidProductPlacement 结构体
type PaidProductPlacement struct {
	HasPaidProductPlacement bool `json:"hasPaidProductPlacement"`
}

// Player 结构体
type Player struct {
	EmbedHtml   string `json:"embedHtml"`
	EmbedHeight int64  `json:"embedHeight"`
	EmbedWidth  int64  `json:"embedWidth"`
}

// TopicDetails 结构体
type TopicDetails struct {
	TopicIds         []string `json:"topicIds"`
	RelevantTopicIds []string `json:"relevantTopicIds"`
	TopicCategories  []string `json:"topicCategories"`
}

// RecordingDetails 结构体
type RecordingDetails struct {
	RecordingDate time.Time `json:"recordingDate"`
}

// FileDetails 结构体
type FileDetails struct {
	FileName     string        `json:"fileName"`
	FileSize     uint64        `json:"fileSize"`
	FileType     string        `json:"fileType"`
	Container    string        `json:"container"`
	VideoStreams []VideoStream `json:"videoStreams"`
	AudioStreams []AudioStream `json:"audioStreams"`
	DurationMs   uint64        `json:"durationMs"`
	BitrateBps   uint64        `json:"bitrateBps"`
	CreationTime string        `json:"creationTime"`
}

// VideoStream 结构体
type VideoStream struct {
	WidthPixels  uint    `json:"widthPixels"`
	HeightPixels uint    `json:"heightPixels"`
	FrameRateFps float64 `json:"frameRateFps"`
	AspectRatio  float64 `json:"aspectRatio"`
	Codec        string  `json:"codec"`
	BitrateBps   uint64  `json:"bitrateBps"`
	Rotation     string  `json:"rotation"`
	Vendor       string  `json:"vendor"`
}

// AudioStream 结构体
type AudioStream struct {
	ChannelCount uint   `json:"channelCount"`
	Codec        string `json:"codec"`
	BitrateBps   uint64 `json:"bitrateBps"`
	Vendor       string `json:"vendor"`
}

// ProcessingDetails 结构体
type ProcessingDetails struct {
	ProcessingStatus          string             `json:"processingStatus"`
	ProcessingProgress        ProcessingProgress `json:"processingProgress"`
	ProcessingFailureReason   string             `json:"processingFailureReason"`
	FileDetailsAvailability   string             `json:"fileDetailsAvailability"`
	ProcessingIssuesAvailable string             `json:"processingIssuesAvailability"`
}

// ProcessingProgress 结构体
type ProcessingProgress struct {
	PartsTotal     uint64 `json:"partsTotal"`
	PartsProcessed uint64 `json:"partsProcessed"`
	TimeLeftMs     uint64 `json:"timeLeftMs"`
}

// Suggestions 结构体
type Suggestions struct {
	ProcessingErrors   []string        `json:"processingErrors"`
	ProcessingWarnings []string        `json:"processingWarnings"`
	ProcessingHints    []string        `json:"processingHints"`
	TagSuggestions     []TagSuggestion `json:"tagSuggestions"`
	EditorSuggestions  []string        `json:"editorSuggestions"`
}

// TagSuggestion 结构体
type TagSuggestion struct {
	Tag               string   `json:"tag"`
	CategoryRestricts []string `json:"categoryRestricts"`
}

// LiveStreamingDetails 结构体
type LiveStreamingDetails struct {
	ActualStartTime    time.Time `json:"actualStartTime"`
	ActualEndTime      time.Time `json:"actualEndTime"`
	ScheduledStartTime time.Time `json:"scheduledStartTime"`
	ScheduledEndTime   time.Time `json:"scheduledEndTime"`
	ConcurrentViewers  uint64    `json:"concurrentViewers"`
	ActiveLiveChatID   string    `json:"activeLiveChatId"`
}
