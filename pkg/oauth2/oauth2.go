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
	"golang.org/x/oauth2"
)

func NewServer(config *configs.Config) (oauth2.Config, *server.Server) {

	authServerURL := "http://localhost:8080"

	oauthConfig := oauth2.Config{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Scopes:       []string{"all"},
		RedirectURL:  authServerURL + "/oauth2/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth2/authorize",
			TokenURL: authServerURL + "/oauth2/token",
		},
	}

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// use mysql token store
	tokenStore := mysql.NewDefaultStore(
		mysql.NewConfig(config.Mysql.Dsn(true)),
	)

	defer tokenStore.Close()

	// token store
	manager.MapTokenStorage(tokenStore)

	// client store
	clientID := oauthConfig.ClientID
	ClientSecret := oauthConfig.ClientSecret

	clientStore := store.NewClientStore()
	clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: ClientSecret,
		Domain: authServerURL,
	})

	manager.MapClientStorage(clientStore)

	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("zerocmf2023"), jwt.SigningMethodHS512))

	// Initialize the oauth2 service
	srv := ginserver.InitServer(manager)
	ginserver.SetAllowGetAccessRequest(true)
	// ginserver.SetClientInfoHandler(server.ClientFormHandler)

	return oauthConfig, srv
}
