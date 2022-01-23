package atoms

import (
	"fmt"
)

type Atom[T comparable] struct {
	val       T
			// key, a, oldState, newState
	watchers  []func(string, Atom[T], T, T)
	validator func(T) bool
}

func New[T comparable](init T) Atom[T] {
	return Atom[T]{
		val: init,
	}
}

func (a *Atom[T]) Deref() T {
	return a.val
}

func (a *Atom[T]) notify(x T) {
	for _, w := range a.watchers {
		w(x)
	}
}

func (a *Atom[T]) Swap(f func(T)T) T {
	// Todo: thread safe this shit
	new := f(a.val)
	a.val = new
	go a.notify(new)
	return new
}

func (a *Atom[T]) Reset(newVal T) {
	a.val = newVal
	go a.notify(newVal)
}

func (a *Atom[T]) CompareAndSet(oldVal, newVal T) bool {
	if a.val==oldVal {
		a.val = newVal
		go a.notify(newVal)
		return true
	}
	return false
}

func (a *Atom[T]) AddWatch(f func(T)bool)
