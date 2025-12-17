package schema

import (
	schema2 "github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/accessTokenClient/posts/schema"
)

// BloggerPostUserInfosGetReq 表示 GET /blogger/v3/users/{userId}/blogs/{blogId}/posts/{postId} API 的请求参数
type BloggerPostUserInfosGetReq struct {
	BlogId      string  `json:"blogId"`                // 必需，博客的ID
	PostId      string  `json:"postId"`                // 必需，要获取的帖子的ID
	UserId      string  `json:"userId"`                // 必需，用户ID（"self"或用户个人资料标识符）
	MaxComments *string `json:"maxComments,omitempty"` // 可选，要为帖子检索的评论数量上限
}

type PostUserInfo struct {
	Kind          string `json:"kind"`          // 此实体的种类，始终为blogger#postPerUserInfo
	UserId        string `json:"userId"`        // 用户的ID
	BlogId        string `json:"blogId"`        // 帖子资源所属的博客的ID
	PostId        string `json:"postId"`        // 帖子资源的ID
	HasEditAccess bool   `json:"hasEditAccess"` // 如果用户对博文具有"作者"级别的访问权限，则为true
}

type PoserUserInfoResource struct {
	Kind         string               `json:"kind"` // 此实体的种类，始终为blogger#postUserInfo
	Post         schema2.PostResource `json:"post"` // "帖子"资源
	PostUserInfo PostUserInfo
}

// BloggerPostUserInfosGetRes 表示 GET /blogger/v3/users/{userId}/blogs/{blogId}/posts/{postId} API 的响应
type BloggerPostUserInfosGetRes struct {
	PoserUserInfoResource // 帖子用户的相关信息
}
