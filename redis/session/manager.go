package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/aliml92/anor/session"
	"log/slog"
	"net/http"
	"time"
)

const (
	defaultAuthLifetime  = 24 * time.Hour
	defaultGuestLifetime = 30 * 24 * time.Hour
)

type Manager struct {
	manager        *scs.SessionManager
	guestSkipPaths map[string][]string
	authLifetime   time.Duration
	guestLifetime  time.Duration
}

type Option func(*Manager)

func NewManager(options ...Option) *Manager {
	sm := &Manager{
		manager:        scs.New(),
		guestSkipPaths: make(map[string][]string),
		authLifetime:   defaultAuthLifetime,
		guestLifetime:  defaultGuestLifetime,
	}

	for _, option := range options {
		option(sm)
	}

	return sm
}

func WithCookieName(name string) Option {
	return func(sm *Manager) {
		sm.manager.Cookie.Name = name
	}
}

func WithGuestSkipPaths(paths map[string][]string) Option {
	return func(sm *Manager) {
		sm.guestSkipPaths = paths
	}
}

func WithAuthLifetime(d time.Duration) Option {
	return func(sm *Manager) {
		sm.manager.Lifetime = d
		sm.authLifetime = d
	}
}

func WithGuestLifetime(d time.Duration) Option {
	return func(sm *Manager) {
		sm.guestLifetime = d
	}
}

func WithStore(store scs.Store) Option {
	return func(sm *Manager) {
		sm.manager.Store = store
	}
}

func WithCodec(codec scs.Codec) Option {
	return func(sm *Manager) {
		sm.manager.Codec = codec
	}
}

func (m *Manager) LoadAndSave(next http.Handler) http.Handler {
	return m.manager.LoadAndSave(next)
}

func (m *Manager) LoadAndSaveGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if methods, ok := m.guestSkipPaths[r.URL.Path]; ok {
			for _, method := range methods {
				if method == r.Method {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		ctx := r.Context()

		isAuth := m.manager.GetBool(ctx, session.IsAuthKey)
		if !isAuth {
			// Set expiration for guest users
			m.manager.Put(ctx, session.IsAuthKey, false)
			m.manager.Lifetime = m.guestLifetime
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (m *Manager) LoadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		u := m.getUser(ctx)
		slog.Info(">>>>>> user loaded from session <<<<<",
			"ID", u.ID,
			"IsAuth", u.IsAuth,
			"CartID", u.CartID,
			"ShippingAddressID", u.ShippingAddressID,
			"BillingAddressID", u.BillingAddressID,
			"StripeConfirmationTokenID", u.StripeConfirmationTokenID,
		)
		ctx = context.WithValue(ctx, session.UserKey, u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Manager) GetExpiry(ctx context.Context) time.Time {
	return m.manager.Deadline(ctx)
}

func (m *Manager) getUser(ctx context.Context) *session.User {
	u := &session.User{}
	isAuth := m.manager.GetBool(ctx, session.IsAuthKey)

	if isAuth {
		u.ID = m.manager.GetInt64(ctx, session.UserIDKey)
		u.Firstname = m.manager.GetString(ctx, session.UserFirstnameKey)
		u.ShippingAddressID = m.manager.GetInt64(ctx, session.ShippingAddressIDKey)
		u.BillingAddressID = m.manager.GetInt64(ctx, session.BillingAddressIDKey)
		u.PaymentMethod = m.manager.GetString(ctx, session.PaymentMethodKey)
		u.StripeConfirmationTokenID = m.manager.GetString(ctx, session.StripeConfirmationTokenIDKey)

		// TODO: get the user from db?
		// or
		// TODO: or make sure to delete this session when deleting the user from db
	}

	u.IsAuth = isAuth
	u.CartID = m.manager.GetInt64(ctx, session.CartIDKey)

	return u
}

func (m *Manager) SetAuthUser(ctx context.Context, u session.User) error {
	if u.ID == 0 {
		return errors.New("userID must be greater than zero")
	}

	if err := m.manager.RenewToken(ctx); err != nil {
		return fmt.Errorf("renew token: %w", err)
	}
	m.manager.Put(ctx, session.IsAuthKey, true)
	m.manager.Put(ctx, session.UserIDKey, u.ID)
	m.manager.Put(ctx, session.UserFirstnameKey, u.Firstname)
	m.manager.Put(ctx, session.CartIDKey, u.CartID)

	// TODO: do these actions in other place
	m.manager.Lifetime = m.authLifetime

	return nil
}

func (m *Manager) Get(ctx context.Context, key string) interface{} {
	return m.manager.Get(ctx, key)
}

func (m *Manager) Put(ctx context.Context, key string, value interface{}) {
	m.manager.Put(ctx, key, value)
}

func (m *Manager) Remove(ctx context.Context, key string) {
	m.manager.Remove(ctx, key)
}

func (m *Manager) Destroy(ctx context.Context) error {
	return m.manager.Destroy(ctx)
}
