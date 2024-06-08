package config

import "github.com/kelseyhightower/envconfig"

const DEV = "development"
const PROD = "production"

type Config struct {
	Environment  string `envconfig:"ELF_ENVIRONMENT" required:"true"`
	ListenAddr   string `envconfig:"ELF_LISTEN_ADDR" default:":4000"`
	SecureCookie SecureCookie
	Session      Session
	Db           Db
	Auth         Auth
	OAuth        OAuth
	Auth0        Auth0
}

func (c Config) IsDevelop() bool {
	return c.Environment == DEV
}

func (c Config) IsProduction() bool {
	return c.Environment == PROD
}

type SecureCookie struct {
	HashKey  string `envconfig:"ELF_SECURECOOKIE_HASHKEY" required:"true"`
	BlockKey string `envconfig:"ELF_SECURECOOKIE_BLOCKKEY" required:"true"`
}

type Session struct {
	Secret string `envconfig:"ELF_SESSION_SECRET" required:"true"`
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
	Domain                      string `envconfig:"ELF_AUTH0_DOMAIN" required:"true"`
	ClientId                    string `envconfig:"ELF_AUTH0_CLIENT_ID" required:"true"`
	ClientSecret                string `envconfig:"ELF_AUTH0_CLIENT_SECRET" required:"true"`
	LoginCallbackUrl            string `envconfig:"ELF_AUTH0_LOGIN_CALLBACK_URL" required:"true"`
	LogoutCallbackUrl           string `envconfig:"ELF_AUTH0_LOGOUT_CALLBACK_URL" required:"true"`
	SessionCookieAccessTokenKey string `envconfig:"ELF_AUTH0_SESSION_COOKIE_ACCESS_TOKEN_KEY" required:"true"`
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
