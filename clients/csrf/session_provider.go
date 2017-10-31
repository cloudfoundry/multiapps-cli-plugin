package csrf

type SessionProvider interface {
	GetSession() error
}
