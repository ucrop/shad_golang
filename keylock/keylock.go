// +build !solution

package keylock

// KeyLock ...
type KeyLock struct {
	mt   chan struct{}
	used map[string]chan struct{}
}

func New() *KeyLock {
	res := &KeyLock{
		mt:   make(chan struct{}, 1),
		used: make(map[string]chan struct{}),
	}
	res.mt <- struct{}{}
	return res
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	for {
		select {
		case <-l.mt:
		case <-cancel:
			return true, nil
		}

		again := false
	Loop:
		for i, key := range keys {
			_, ok := l.used[key]
			if !ok {
				l.used[key] = make(chan struct{}, 1)
				l.used[key] <- struct{}{}
			}

			select {
			case <-l.used[key]:
			default:
				for j := 0; j < i; j++ {
					l.used[keys[j]] <- struct{}{}
				}
				l.mt <- struct{}{}
				again = true
				break Loop
			}
		}

		if again {
			continue
		}

		l.mt <- struct{}{}

		return false, func() {
			for _, key := range keys {
				l.used[key] <- struct{}{}
			}
		}
	}
}
