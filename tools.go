package gomini

import (
	"github.com/BurntSushi/toml"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tidwall/sjson"
	"github.com/bwmarrin/snowflake"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var cfg config
var node *snowflake.Node
type config struct {
	Dsn string
	signingKey string
}

func LoadCfg() {
	files, _ := ioutil.ReadDir("./configs")
	f := "./configs/" + files[0].Name()
	if _, err := toml.DecodeFile(f, &cfg); err != nil {
		log.Println("toml config:", f, err)
		return
	}
	initDB()
	node, _ = snowflake.NewNode(1)
}

type JsonClaims struct {
	Payload string `json:"payload"`
	jwt.RegisteredClaims
}

func JwtToken(js string) string {
	claims := JsonClaims{
		js,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "metactrl",
			Subject:   "traning_server",
			ID:        "1",
			Audience:  []string{"traning_client"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rs, err := token.SignedString([]byte(cfg.signingKey))
	if err != nil {
		log.Println(err)
	}
	return rs
}

func Req2json(r *http.Request) (string, string, error) {
	r.ParseForm()

	tokenString := r.Header.Get("Authorization")
	// tokenString := JwtToken(`{"uid":001,"uname":"ff"}`)
	header := ""
	token, err := jwt.ParseWithClaims(tokenString, &JsonClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.signingKey), nil
	})
	if claims, ok := token.Claims.(*JsonClaims); ok && token.Valid {
		log.Println(claims.Payload)
		header = claims.Payload
	} else {
		return "", "", err
	}

	body, _ := ioutil.ReadAll(r.Body)
	//    r.Body.Close()
	dat := strings.TrimSpace(string(body))
	// msg := "{}"
	for k, v := range r.Form {
		if strings.Index(k, "[]") > 1 {
			dat, _ = sjson.Set(dat, strings.Replace(k, "[]", "", 1), v)
		} else {
			dat, _ = sjson.Set(dat, k, v[0])
			if v[0] == "true" {
				dat, _ = sjson.Set(dat, k, true)
			}
			if v[0] == "false" {
				dat, _ = sjson.Set(dat, k, false)
			}
		}
	}
	if dat == "" {
		dat = "{}"
	}
	req, _ := sjson.SetRaw("", "req", dat)
	return header, req, nil
}

func Id() string {
	return node.Generate().String()
}