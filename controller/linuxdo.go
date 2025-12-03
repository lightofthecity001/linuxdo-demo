package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ClientId       string
	ClientSecret   string
	CallbackUrl    string
	DisableUserUrl string
}

var config EnvConfig

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("错误：无法加载 .env 文件")
	}
	config.ClientId = os.Getenv("client_id")
	config.ClientSecret = os.Getenv("client_secret")
	config.CallbackUrl = os.Getenv("callback_url")
	config.DisableUserUrl = os.Getenv("disable_user_url")
}

type LinuxdoUser struct {
	Id             int         `json:"id"`
	Sub            string      `json:"sub"`
	Username       string      `json:"username"`
	Login          string      `json:"login"`
	Name           string      `json:"name"`
	Email          string      `json:"email"`
	AvatarTemplate string      `json:"avatar_template"`
	AvatarUrl      string      `json:"avatar_url"`
	Active         bool        `json:"active"`
	TrustLevel     int         `json:"trust_level"`
	Silenced       bool        `json:"silenced"`
	ExternalIds    interface{} `json:"external_ids"`
	ApiKey         string      `json:"api_key"`
}

var loginUser LinuxdoUser

var authorizeUrl string = "https://connect.linux.do/oauth2/authorize?response_type=code&client_id=%s&state=%s"
var tokenUrl string = "https://connect.linux.do/oauth2/token"
var userUrl string = "https://connect.linux.do/api/user"

func CallAuthorize(c *gin.Context) {
	log.Println("111111111111111111", c.Request.URL.Query())
	targetUrl := fmt.Sprintf(authorizeUrl, config.ClientId, c.Request.URL.Query().Get("state"))
	c.Redirect(http.StatusFound, targetUrl)
}

func Callback(c *gin.Context) {
	log.Println("22222222222222", c.Request.URL.Query())
	code := c.Request.URL.Query().Get("code")
	state := c.Request.URL.Query().Get("state")
	c.Redirect(http.StatusFound, config.CallbackUrl+"?code="+code+"&state="+state)
}

func GetToken(c *gin.Context) {
	code := c.PostForm("code")
	log.Println("33333333333333333333", code)
	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", config.ClientId)
	formData.Set("client_secret", config.ClientSecret)
	resp, err := http.PostForm(tokenUrl, formData)
	if err != nil {
		log.Fatal("调用token接口失败")
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Println("服务端返回的详细错误:", string(bodyBytes))
	c.Data(resp.StatusCode, "application/json", bodyBytes)
}

func GetUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	log.Println("authHeader", authHeader)
	client := &http.Client{}
	req, err := http.NewRequest("GET", userUrl, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", authHeader)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Println("44444444444", string(body))
	json.Unmarshal(body, &loginUser)
	log.Println("55555555555", loginUser.Name)
	if IsUserBanned(strconv.Itoa(loginUser.Id)) {
		loginUser.Id = 0
	}
	c.JSON(http.StatusOK, loginUser)
}
