package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

var srv *server.Server
var clientStore *store.ClientStore

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var Persons = []Person{
	{
		ID:   1,
		Name: "Guilherme",
		Age:  22,
	},
	{
		ID:   2,
		Name: "Jorge",
		Age:  35,
	},
}

func Config() {
	manager := manage.NewDefaultManager()

	// armazenamento de token
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// armazenamento de cliente
	clientStore = store.NewClientStore()
	manager.MapClientStorage(clientStore)

	srv = server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	refreshTokenCfg := &manage.RefreshingConfig{
		AccessTokenExp:     time.Hour,
		RefreshTokenExp:    2 * time.Hour,
		IsGenerateRefresh:  true,
		IsResetRefreshTime: false,
	}

	manager.SetRefreshTokenCfg(refreshTokenCfg)
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
	c.JSON(200, gin.H{"clientId": clientId, "clientSecret": clientSecret})
}

func MiddlewareAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		_, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})

		}
	}

}

func HelloPerson(c *gin.Context) {
	c.JSON(200, gin.H{
		"Return": Persons,
	})
}

func PostPerson(c *gin.Context) {
	var Post Person
	err := c.ShouldBindJSON(&Post)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{
		"Return": Post,
	})
}

func main() {
	Config()
	r := gin.Default()

	v1 := r.Group("v1")
	v2 := r.Group("v2")
	v2.Use(MiddlewareAuth())

	v2.POST("/PostPerson", PostPerson).Use(MiddlewareAuth())
	v2.GET("/HelloPerson", HelloPerson).Use(MiddlewareAuth())

	v1.GET("/token", Token)
	v1.GET("/credentials", Credentials)
	r.Run(":9090")

}
