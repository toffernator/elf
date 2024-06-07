package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ListenAddr string `envconfig:"ELF_LISTEN_ADDR" default:":4000"`
	Db         Db
	Auth       Auth
	OAuth      OAuth
	Auth0      Auth0
}

type Auth struct {
	SessionCookieName    string `envconfig:"ELF_SESSION_COOKIE_NAME" default:"session"`
	SessionCookieUserKey string `envconfig:"ELF_SESSION_COOKIE_USER_KEY" default:"user"`
}

type OAuth struct {
	StateLength     int    `envconfig:"ELF_OAUTH_STATE_LENGTH" default:"16"`
	StateCookieName string `envconfig:"ELF_OAUTH_STATE_COOKIE_NAME" default:"oauthstate"`
}

type Auth0 struct {
	Domain                      string `envconfig:"ELF_AUTH0_DOMAIN"`
	ClientId                    string `envconfig:"ELF_AUTH0_CLIENT_ID"`
	ClientSecret                string `envconfig:"ELF_AUTH0_CLIENT_SECRET"`
	LoginCallbackUrl            string `envconfig:"ELF_AUTH0_LOGIN_CALLBACK_URL"`
	LogoutCallbackUrl           string `envconfig:"ELF_AUTH0_LOGOUT_CALLBACK_URL"`
	SessionCookieAccessTokenKey string `envconfig:"ELF_AUTH0_SESSION_COOKIE_ACCESS_TOKEN_KEY"`
}

type Db struct {
	Name string `envconfig:"ELF_DB_NAME" default:"elf.db"`
}

func loadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func MustLoadConfig() *Config {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
