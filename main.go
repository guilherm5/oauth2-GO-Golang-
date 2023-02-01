package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

var srv *server.Server
var clientStore *store.ClientStore

func Config() {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	// armazenamento de token
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// armazenamento de cliente
	clientStore = store.NewClientStore()
	manager.MapClientStorage(clientStore)

	srv = server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

}

func Token(c *gin.Context) {
	srv.HandleTokenRequest(c.Writer, c.Request)
}

func Credentials(c *gin.Context) {
	clientId := uuid.New().String()
	clientSecret := uuid.New().String()
	err := clientStore.Set(clientId, &models.Client{
		ID:     clientId,
		Secret: clientSecret,
		Domain: "http://localhost:9094",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	c.Header("Content-Type", "application/json")
	c.JSON(200, gin.H{"CLIENT_ID": clientId, "CLIENT_SECRET": clientSecret})
}

func Protedcted(c *gin.Context) {
	_, err := srv.ValidationBearerToken(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.String(200, "Hello World")
}

func main() {
	Config()
	r := gin.Default()

	r.GET("/token", Token)

	r.GET("/credentials", Credentials)

	r.GET("/protected", Protedcted)

	r.Run(":9090")
}
