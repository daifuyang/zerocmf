package oauth2

import (
	"zerocmf/configs"

	ginserver "github.com/go-oauth2/gin-server"
	"github.com/go-oauth2/mysql/v4"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
)

func NewServer(config *configs.Config, tokenStore *mysql.Store) *server.Server {

	oauth2Conf := config.Oauth2

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MapTokenStorage(tokenStore)

	// client store
	clientID := oauth2Conf.ClientID
	ClientSecret := oauth2Conf.ClientSecret

	clientStore := store.NewClientStore()
	clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: ClientSecret,
		Domain: oauth2Conf.AuthServerURL,
	})

	manager.MapClientStorage(clientStore)

	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("zerocmf2023"), jwt.SigningMethodHS512))

	// Initialize the oauth2 service
	srv := ginserver.InitServer(manager)
	ginserver.SetAllowGetAccessRequest(true)
	// ginserver.SetClientInfoHandler(server.ClientFormHandler)

	return srv
}
