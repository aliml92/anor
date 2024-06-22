package session

import "github.com/alexedwards/scs/v2"

type Manager struct {
	Auth  *scs.SessionManager
	Guest *scs.SessionManager
}

func NewManager(auth *scs.SessionManager, guest *scs.SessionManager) *Manager {
	return &Manager{
		Auth:  auth,
		Guest: guest,
	}
}
