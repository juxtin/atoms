package atoms

import "sync"

type Atom[T comparable] struct {
	val       T
	watchers  []watcher[T]
	validator func(T) error
	m         *sync.Mutex
}

type watcher[T comparable] struct {
	key string
	// watcher = func(key string, a Atom[T], oldState, newState T)
	f func(string, *Atom[T], T, T)
}

func New[T comparable](init T) Atom[T] {
	var m sync.Mutex
	return Atom[T]{
		val: init,
		validator: func(_ T)error { return nil },
		m: &m,
	}
}

func (a *Atom[T]) Deref() T {
	return a.val
}

func (a *Atom[T]) notify(old, new T) {
	for _, w := range a.watchers {
		w.f(w.key, a, old, new)
	}
}

func (a *Atom[T]) Swap(f func(T)T) (T, error) {
	// TODO: Clojure's atoms don't lock, and ultimately neither should these.
	// For now though, it's a fine corner to cut.
	a.m.Lock()
	old := a.val
	new := f(a.val)
	if err := a.validator(new); err != nil {
		a.m.Unlock()
		return old, err
	}
	a.val = new
	a.m.Unlock()
	go a.notify(old, new)
	return new, nil
}

func (a *Atom[T]) Reset(newVal T) error {
	// No need to lock for validation
	if err := a.validator(newVal); err != nil {
		return err
	}
	a.m.Lock()
	oldVal := a.val
	a.val = newVal
	a.m.Unlock()
	go a.notify(oldVal, newVal)
	return nil
}

func (a *Atom[T]) CompareAndSet(oldVal, newVal T) bool {
	a.m.Lock()
	if a.val==oldVal {
		a.val = newVal
		a.m.Unlock()
		go a.notify(oldVal, newVal)
		return true
	}
	a.m.Unlock()
	return false
}

func (a *Atom[T]) AddWatch(key string, f func(string, *Atom[T], T, T)) {
	w := watcher[T]{
		key: key,
		f: f,
	}
	a.watchers = append(a.watchers, w)
}

func (a *Atom[T]) SetValidator(f func(T)error) {
	a.validator = f
}