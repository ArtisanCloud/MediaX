package publish

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	response2 "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/publish/request"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/publish/response"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type OfficialAccountPublishClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *OfficialAccountPublishClient {
	return &OfficialAccountPublishClient{
		BaseClient: c,
	}
}

// DraftAdd 新建草稿
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Add_draft.html
func (comp *OfficialAccountPublishClient) DraftAdd(ctx context.Context, data *request.DraftAddReq) (*response.DraftAddRes, error) {
	result := &response.DraftAddRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/draft/add", nil, data, nil, result)
	return result, err
}

// DraftGet 获取草稿
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Get_draft.html
func (comp *OfficialAccountPublishClient) DraftGet(ctx context.Context, mediaID string) (*response.DraftGetRes, error) {
	result := &response.DraftGetRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/draft/get", nil, &object.HashMap{
		"media_id": mediaID,
	}, nil, result)

	return result, err
}

// DraftDelete 删除草稿
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Delete_draft.html
func (comp *OfficialAccountPublishClient) DraftDelete(ctx context.Context, mediaID string) (*response2.OfficialAccountRes, error) {
	result := &response2.OfficialAccountRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/draft/delete", nil, &object.HashMap{
		"media_id": mediaID,
	}, nil, result)

	return result, err
}

// DraftUpdate 修改草稿
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Update_draft.html
func (comp *OfficialAccountPublishClient) DraftUpdate(ctx context.Context, data *request.DraftUpdateReq) (*response2.OfficialAccountRes, error) {
	result := &response2.OfficialAccountRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/draft/update", nil, data, nil, result)

	return result, err
}

// DraftCount 获取草稿总数
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Count_drafts.html
func (comp *OfficialAccountPublishClient) DraftCount(ctx context.Context) (*response.DraftCountRes, error) {
	result := &response.DraftCountRes{}

	_, err := comp.BaseClient.HttpGet(ctx, "cgi-bin/draft/count", &object.StringMap{}, nil, result)

	return result, err
}

// DraftBatchGet 获取草稿列表
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Count_drafts.html
func (comp *OfficialAccountPublishClient) DraftBatchGet(ctx context.Context, data *request.BatchGetReq) (*response.BatchGetRes, error) {
	result := &response.BatchGetRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/draft/batchget", nil, data, nil, result)

	return result, err
}

// DraftSwitch MP端开关
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Temporary_MP_Switch.html
func (comp *OfficialAccountPublishClient) DraftSwitch(ctx context.Context) (*response2.OfficialAccountRes, error) {
	result := &response2.OfficialAccountRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/draft/switch", nil, &object.HashMap{}, nil, result)

	return result, err
}

// DraftCheckSwitch M检查P端开关
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Temporary_MP_Switch.html
func (comp *OfficialAccountPublishClient) DraftCheckSwitch(ctx context.Context) (*response.CheckSwitchRes, error) {
	result := &response.CheckSwitchRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/draft/switch", nil, &object.HashMap{}, &object.StringMap{
		"checkonly": "1",
	}, result)

	return result, err
}

// PublishSubmit 发布接口
// https://developers.weixin.qq.com/doc/offiaccount/Publish/Publish.html
func (comp *OfficialAccountPublishClient) PublishSubmit(ctx context.Context, mediaID string) (*response.PublishSubmitRes, error) {
	result := &response.PublishSubmitRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/freepublish/submit", nil, &object.HashMap{
		"media_id": mediaID,
	}, nil, result)

	return result, err
}

// PublishGet 发布状态轮询接口
// https://developers.weixin.qq.com/doc/offiaccount/Publish/Get_status.html
func (comp *OfficialAccountPublishClient) PublishGet(ctx context.Context, publishID uint64) (*response.PublishGetRes, error) {
	result := &response.PublishGetRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/freepublish/get", nil, &object.HashMap{
		"publish_id": publishID,
	}, nil, result)

	return result, err
}

// PublishDelete 删除发布
// https://developers.weixin.qq.com/doc/offiaccount/Publish/Delete_posts.html
func (comp *OfficialAccountPublishClient) PublishDelete(ctx context.Context, articleID string, index int) (*response2.OfficialAccountRes, error) {
	result := &response2.OfficialAccountRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/freepublish/delete", nil, &object.HashMap{
		"article_id": articleID,
		"index":      index,
	}, nil, result)

	return result, err
}

// PublishGetArticle 通过 article_id 获取已发布文章
// https://developers.weixin.qq.com/doc/offiaccount/Publish/Get_article_from_id.html
func (comp *OfficialAccountPublishClient) PublishGetArticle(ctx context.Context, articleID string) (*response.PublishGetArticleRes, error) {
	result := &response.PublishGetArticleRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/freepublish/getarticle", nil, &object.HashMap{
		"article_id": articleID,
	}, nil, result)

	return result, err
}

// PublishBatchGet 获取成功发布列表
// https://developers.weixin.qq.com/doc/offiaccount/Publish/Get_publication_records.html
func (comp *OfficialAccountPublishClient) PublishBatchGet(ctx context.Context, data *request.BatchGetReq) (*response.BatchGetRes, error) {
	result := &response.BatchGetRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "cgi-bin/freepublish/batchget", nil, data, nil, result)

	return result, err
}
