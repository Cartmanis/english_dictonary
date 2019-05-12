package provider_session

import (
	"container/list"
	"english_dictonary/app/cmd/session"
	"fmt"
	"sync"
	"time"
)

type Provider struct {
	sessions map[string]*list.Element
	list     *list.List
	lock     sync.Mutex
}

var prov = &Provider{list: list.New()}

type SessionStore struct {
	sid          string //session id
	timeAccessed time.Time
	value        map[interface{}]interface{} //session
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	return prov.SessionUpdate(st.sid)
}

func (st *SessionStore) Get(key interface{}) interface{} {
	if err := prov.SessionUpdate(st.sid); err != nil {
		fmt.Println("[ERROR] не удалось выполнить метод Get. Ошибка:", err)
		return nil
	}
	if v, ok := st.value[key]; ok {
		return v
	}
	return nil
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	return prov.SessionUpdate(st.sid)
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

func (p *Provider) SessionInit(sid string) (session.Session, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newSession := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := prov.list.PushBack(newSession)
	prov.sessions[sid] = element
	return newSession, nil
}

func (p *Provider) SessionRead(sid string) (session.Session, error) {
	//если есть, то возвращаем
	if elem, ok := p.sessions[sid]; ok {
		return elem.Value.(*SessionStore), nil
	}
	//если нет, то иницилизируем
	return prov.SessionInit(sid)
}

func (p *Provider) SessionDestroy(sid string) error {
	if element, ok := prov.sessions[sid]; ok {
		delete(prov.sessions, sid)
		prov.list.Remove(element)
	}
	return nil
}

func (p *Provider) SessionGC(maxLifeItem int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for {
		elem := prov.list.Back()
		if elem == nil {
			break
		}
		if (elem.Value.(*SessionStore).timeAccessed.Unix() + maxLifeItem) < time.Now().Unix() {
			prov.list.Remove(elem)
			delete(prov.sessions, elem.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (p *Provider) SessionUpdate(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if elem, ok := prov.sessions[sid]; ok {
		elem.Value.(*SessionStore).timeAccessed = time.Now()
		prov.list.MoveToFront(elem) //перемещаем элемент вначало
		return nil
	}
	return nil
}

func init() {
	prov.sessions = make(map[string]*list.Element, 0)
	if err := session.Register("memory", prov); err != nil {
		fmt.Println("[ERROR] не удалось создать сессию. Ошибка:", err)
	}
}
