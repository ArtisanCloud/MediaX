package kernel

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel/response"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/contract"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/helper"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type AccessToken struct {
	HttpHelper *helper.RequestHelper

	RequestMethod      string
	EndpointToGetToken string
	QueryName          string
	Token              *object.HashMap
	TokenKey           string
	CacheTokenKey      string
	CachePrefix        string

	GetCredentials func() *object.StringMap
	GetEndpoint    func() (string, error)

	SetCustomToken func(token *response.AccessTokenRes) interface{}
	GetCustomToken func(key string, refresh bool) object.HashMap

	GetMiddlewareOfLog func(l *logger.Logger) contract.RequestMiddleware
}
