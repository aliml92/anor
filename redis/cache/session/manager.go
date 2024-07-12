package session

import (
	"github.com/alexedwards/scs/v2"
	"github.com/aliml92/anor/config"
	"github.com/redis/go-redis/v9"
)

type Manager struct {
	Auth  *scs.SessionManager
	Guest *scs.SessionManager
}

func NewManager(cfg config.SessionConfig, client *redis.Client) *Manager {
	as := scs.New()
	as.Store = NewRedisStore(client).WithPrefix("keys:session:user:authenticated:")
	as.Codec = MessagePackCodec{}
	as.Lifetime = cfg.AuthLifetime
	as.Cookie.Name = cfg.AuthCookieName
	//auth.Cookie.Name = "__anor_ust" // anor, authenticated user's session token

	gs := scs.New()
	gs.Store = NewRedisStore(client).WithPrefix("keys:session:user:guest:")
	gs.Codec = MessagePackCodec{}
	gs.Lifetime = cfg.GuestLifetime
	gs.Cookie.Name = cfg.GuestCookieName
	//guest.Cookie.Name = "__anor_gst" // anor, guest's session token

	return &Manager{
		Auth:  as,
		Guest: gs,
	}
}
