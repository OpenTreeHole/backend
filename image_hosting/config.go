package main

//
//import (
//	"github.com/caarlos0/env/v9"
//	"net/url"
//	"sync"
//)
//
//var Config struct {
//	ApiStyle      string `env:"API_STYLE" envDefault:"chevereto"`
//	ProxyType     string `env:"PROXY_TYPE" envDefault:"lskypro"`
//	ProxyUrl      string `env:"PROXY_URL,required"`
//	ProxyUsername string `env:"PROXY_USERNAME"`
//	ProxyPassword string `env:"PROXY_PASSWORD,required"`
//	ProxyEmail    string `env:"PROXY_EMAIL"`
//	HostRewrite   string `env:"HOST_REWRITE"`
//	Token         string `env:"TOKEN"`
//}
//
//var Token struct {
//	sync.RWMutex
//	data string
//}
//
//var ProxyUrlData *url.URL
//
//func GetToken() string {
//	Token.RLock()
//	defer Token.RUnlock()
//	return Token.data
//}
//
//func init() {
//	err := env.Parse(&Config)
//	if err != nil {
//		panic(err)
//	}
//	ProxyUrlData, err = url.ParseRequestURI(Config.ProxyUrl)
//	if err != nil {
//		panic(err)
//	}
//	if Config.Token == "" {
//		Config.Token, err = LskyRefreshToken()
//		if err != nil {
//			panic(err)
//		}
//	}
//	Token.data = Config.Token
//	log.Printf("%+v", Token.data)
//}
