package example

import "sync"

type SafetyMap struct {
	m 	map[string]string
	mutex   sync.RWMutex
}

func (sm *SafetyMap) Add(k string, v string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.m[k] = v
}

// LoadOrStore double-check
func (sm *SafetyMap) LoadOrStore(k, v string) (string, bool) {
	// 读锁可叠加
	sm.mutex.RLock()
	v, ok := sm.m[k]
	sm.mutex.RUnlock()  // 拿到读锁加写锁会panic
	if ok {
		return v, true
	}
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.m[k] = v
	return v, false
}
