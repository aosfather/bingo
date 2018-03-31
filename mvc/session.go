package mvc

import (
	"sync"
	"net/http"
	"time"
	"io"
	"strconv"
	"encoding/base64"
	"crypto/rand"
	"github.com/go-redis/redis"
)

type SessionStore interface {
	Exist(id string) bool
	Create(id string)
	GetValue(id, key string, value interface{})
	SetValue(id, key string, value interface{})
	Touch(id string)
	Delete(id string)
}

type memoryStore struct {
	sessions map[string]map[string]interface{}
}

func (this *memoryStore)Init(){
	this.sessions=make(map[string]map[string]interface{})
}

func (this *memoryStore) Exist(id string) bool{
	s:=this.sessions[id]
	return s!=nil

}
func (this *memoryStore) Create(id string){
	s:=make(map[string]interface{})
	this.sessions[id]=s
}
func (this *memoryStore) GetValue(id, key string,value interface{}){
	v:=this.sessions[id]
	if v!=nil {
		value= v[key]
		return
	}
	value=nil
}
func (this *memoryStore) SetValue(id, key string, value interface{}){
	v:=this.sessions[id]
	if v!=nil {
		v[key]=value
	}
}
func (this *memoryStore) Touch(id string){

}
func (this *memoryStore) Delete(id string){
	delete(this.sessions,id)
}

/*会话*/
type HttpSession struct {
	mLock            sync.RWMutex //互斥(保证线程安全)
	mSessionID       string       //唯一id
	mNew             bool
	mStrore          SessionStore
	lastTimeAccessed time.Time
}

func (this *HttpSession) IsNew() bool {
	return this.mNew
}

func (this *HttpSession) Touch() {
	this.mStrore.Touch(this.mSessionID)

}

func (this *HttpSession) SetValue(key string, value interface{}) {
	if key != "" {
		this.mLock.Lock()
		defer this.mLock.Unlock()
		this.mStrore.SetValue(this.mSessionID, key, value)
	}
}

func (this *HttpSession) GetValue(key string,value interface{})  {
	this.mLock.RLock()
	defer this.mLock.RUnlock()
	//如果找不到，尝试者从store中获取，若也无法获取则返回nil
	this.mStrore.GetValue(this.mSessionID, key,value)
}

type sessionManager struct {
	cookieName   string                  //客户端cookie名称
	lock         sync.RWMutex            //互斥(保证线程安全)
	store        SessionStore            //session存储对象
	maxCacheTime int64                   //垃圾回收时间
	mMaxLifeTime int64                   //session的有效时间
	sessions     map[string]*HttpSession //保存session的指针[sessionID] = session
}

func (this *sessionManager) Init() {
	//初始化
	this.sessions=make(map[string]*HttpSession)


	if this.store==nil {

		m:=memoryStore{}
		m.Init()
		this.store=&m
	}
	//启动定期清理
	go this.gc()
}

func (this *sessionManager) getSession(w http.ResponseWriter, r *http.Request) *HttpSession {


	var cookie, err = r.Cookie(this.cookieName)

	if cookie != nil && err == nil {

		this.lock.Lock()
		defer this.lock.Unlock()

		sessionID := cookie.Value
		if session, ok := this.sessions[sessionID]; ok {
			session.lastTimeAccessed = time.Now() //判断合法性的同时，更新最后的访问时间
			return session
		} else if this.store.Exist(sessionID) { //从store中检查session，如果存在则加载
			session = &HttpSession{mSessionID: sessionID, lastTimeAccessed: time.Now(), mNew: false, mStrore: this.store}
			this.sessions[sessionID] = session
			return session

		}

	}
	return this.create(w)

}

func (this *sessionManager) create(w http.ResponseWriter) *HttpSession {
	newSessionID := newSessionID()
	session := &HttpSession{mSessionID: newSessionID, lastTimeAccessed: time.Now(), mNew: true, mStrore: this.store}
	this.sessions[newSessionID] = session
	//让浏览器cookie设置过期时间
	cookie := http.Cookie{Name: this.cookieName, Value: newSessionID, Path: "/", HttpOnly: true, MaxAge: int(this.mMaxLifeTime)}
	http.SetCookie(w, &cookie)
	return nil
}

func (this *sessionManager) deleteSession(id string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.sessions, id)
	this.store.Delete(id)

}

func (this *sessionManager) gc() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for sessionID, session := range this.sessions {
		//删除超过时限的session
		if session.lastTimeAccessed.Unix()+this.maxCacheTime < time.Now().Unix() {
			delete(this.sessions, sessionID)
		}
	}

	//定时回收
	time.AfterFunc(time.Duration(this.maxCacheTime)*time.Second, func() { this.gc() })

}

func newSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		nano := time.Now().UnixNano() //微秒
		return strconv.FormatInt(nano, 10)
	}
	return base64.URLEncoding.EncodeToString(b)
}


// redis session store
type RedisSessionStore struct {
	client *redis.Client
	prefix string
}




