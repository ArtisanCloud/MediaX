package blogger

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/blogUserInfos"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/blogs"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/comments"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/pageViews"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/pages"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/postUserInfos"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/posts"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/users"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/core"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type GoogleBloggerClient struct {
	//Logger             *logger.Logger
	//Cache              cache.ICache
	GoogleClient       *core.GoogleClient
	BloggerConfig      *config.GoogleBloggerConfig
	AccessTokenHandler *core.GoogleAccessTokenHandler

	// clients
	blogs         *blogs.BloggerBlogsClient
	blogUserInfos *blogUserInfos.BloggerBlogUserInfosClient
	comments      *comments.BloggerCommentsClient
	pages         *pages.BloggerPagesClient
	pageViews     *pageViews.BloggerPageViewsClient
	posts         *posts.BloggerPostsClient
	postUserInfos *postUserInfos.BloggerPostUserInfosClient
	users         *users.BloggerUsersClient
}

func NewGoogleBloggerClient(cfg *config.GoogleBloggerConfig, logger *logger.Logger, cache cache.ICache) (*GoogleBloggerClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = config.GoogleBloggerAPIUrl
	}
	c, err := core.NewGoogleClient(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	handler, err := core.NewGoogleAccessTokenHandler(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	// bind token handler to client
	c.TokenHandler = handler.AccessTokenHandler

	// override get custom token
	c.TokenHandler.GetCustomToken = cfg.GetOAuthToken

	return &GoogleBloggerClient{
		//Logger:             logger,
		//Cache:              cache,
		GoogleClient:       c,
		BloggerConfig:      cfg,
		AccessTokenHandler: handler,
	}, nil
}

func (client *GoogleBloggerClient) GetBlogsClient() *blogs.BloggerBlogsClient {
	if client.blogs == nil {
		client.blogs = blogs.NewClient(client.GoogleClient.BaseClient)
	}
	return client.blogs
}

func (client *GoogleBloggerClient) GetBlogUserInfosClient() *blogUserInfos.BloggerBlogUserInfosClient {
	if client.blogUserInfos == nil {
		client.blogUserInfos = blogUserInfos.NewClient(client.GoogleClient.BaseClient)
	}
	return client.blogUserInfos
}

func (client *GoogleBloggerClient) GetCommentsClient() *comments.BloggerCommentsClient {
	if client.comments == nil {
		client.comments = comments.NewClient(client.GoogleClient.BaseClient)
	}
	return client.comments
}

func (client *GoogleBloggerClient) GetPageViewsClient() *pageViews.BloggerPageViewsClient {
	if client.pageViews == nil {
		client.pageViews = pageViews.NewClient(client.GoogleClient.BaseClient)
	}
	return client.pageViews
}

func (client *GoogleBloggerClient) GetPagesClient() *pages.BloggerPagesClient {
	if client.pages == nil {
		client.pages = pages.NewClient(client.GoogleClient.BaseClient)
	}
	return client.pages
}

func (client *GoogleBloggerClient) GetPostUserInfosClient() *postUserInfos.BloggerPostUserInfosClient {
	if client.postUserInfos == nil {
		client.postUserInfos = postUserInfos.NewClient(client.GoogleClient.BaseClient)
	}
	return client.postUserInfos
}

func (client *GoogleBloggerClient) GetPostsClient() *posts.BloggerPostsClient {
	if client.posts == nil {
		client.posts = posts.NewClient(client.GoogleClient.BaseClient)
	}
	return client.posts
}

func (client *GoogleBloggerClient) GetUsersClient() *users.BloggerUsersClient {
	if client.users == nil {
		client.users = users.NewClient(client.GoogleClient.BaseClient)
	}
	return client.users
}
