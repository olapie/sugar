package mobilex

import "sync"

type SafeFlags struct {
	mu    sync.Mutex
	flags int64
}

func (f *SafeFlags) TryLock(flag int64) bool {
	if f.flags&flag != 0 {
		return false
	}

	f.mu.Lock()
	if f.flags&flag != 0 {
		f.mu.Unlock()
		return false
	}
	f.flags |= flag
	f.mu.Unlock()
	return true
}

func (f *SafeFlags) Unlock(flag int64) {
	f.mu.Lock()
	f.flags &= ^flag
	f.mu.Unlock()
}

func (f *SafeFlags) IsOn(flag int64) bool {
	return f.flags&flag != 0
}
