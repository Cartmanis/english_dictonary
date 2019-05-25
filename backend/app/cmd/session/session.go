package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxLifeTime int64
}

//глобальный map, который хранит все сессии
var provides = make(map[string]Provider)

func Register(name string, p Provider) error {
	if p == nil {
		return fmt.Errorf("был передан неиницилизированный Provider")
	}
	if _, ok := provides[name]; ok {
		return fmt.Errorf("cессия с именем %v уже существует", name)
	}
	provides[name] = p
	return nil
}

func NewManager(providerName, cookieName string, maxLifeTime int64) (*Manager, error) {
	p, ok := provides[providerName]
	if !ok {
		return nil, fmt.Errorf("неизвестный провайдер с именем %v", providerName)
	}
	return &Manager{
		cookieName:  cookieName,
		lock:        sync.Mutex{},
		provider:    p,
		maxLifeTime: maxLifeTime,
	}, nil
}

//основная функция для API создания сессии

func (m *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (Session, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	cookie, err := r.Cookie(m.cookieName)
	//если по тем или иным причинам cookie не была найдена в запросе
	//просто устанавливаем новую cookie
	if err != nil || cookie.Value == "" {
		sid, err := createSessionId()
		if err != nil {
			return nil, err
		}
		session, err := m.provider.SessionInit(sid)
		if err != nil {
			return nil, err
		}
		cookie := &http.Cookie{Name: m.cookieName, Value: url.QueryEscape(sid),
			Path: "/", HttpOnly: true, MaxAge: int(m.maxLifeTime)}
		http.SetCookie(w, cookie)
		return session, nil
	}
	sid := url.QueryEscape(cookie.Value)
	return m.provider.SessionRead(sid)
}

func (m *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	if err := m.provider.SessionDestroy(cookie.Value); err != nil {
		fmt.Println("[ERROR] не удалось уничтожить сессию. Ошибка:", err)
		return
	}
	cookie = &http.Cookie{Name: m.cookieName, Path: "/", HttpOnly: true,
		Expires: time.Now(), MaxAge: -1}
	http.SetCookie(w, cookie)
}

func (m *Manager) GC() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.provider.SessionGC(m.maxLifeTime)
	time.AfterFunc(time.Duration(m.maxLifeTime)*time.Second,
		func() { m.GC() })

}

//приватные функции и методы

func createSessionId() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
