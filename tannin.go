package tannin

import (
	"unsafe"

	"github.com/floatdrop/lru"
	"github.com/floatdrop/slru"
)

const (
	DefaultWindowRatio = 0.1
	DefaultSLRUSplit   = 0.2
)

// Tannin is a thread-safe fixed size cache with windowed TinyLFU eviction policy.
type Tannin[K comparable, V any] struct {
	window  *lru.LRU[K, V]
	c       cm4
	w       int
	samples int
	main    *slru.SLRU[K, V]
	hf      func(unsafe.Pointer, uintptr) uintptr
}

func (T *Tannin[K, V]) hash(key K) uintptr {
	return T.hf(unsafe.Pointer(&key), 0)
}

// Evicted holds key/value pair that was evicted from cache.
type Evicted[K comparable, V any] struct {
	Key   K
	Value V
}

func (T *Tannin[K, V]) Get(key K) *V {
	T.c.add(T.hash(key))

	T.w++
	if T.w >= T.samples {
		T.c.reset()
		T.w = 0
	}

	if e := T.main.Get(key); e != nil {
		return e
	}

	if e := T.window.Get(key); e != nil {
		return e
	}

	return nil
}

func fromLruEvicted[K comparable, V any](e *lru.Evicted[K, V]) *Evicted[K, V] {
	if e == nil {
		return nil
	}

	return &Evicted[K, V]{
		e.Key,
		e.Value,
	}
}

func fromSlruEvicted[K comparable, V any](e *slru.Evicted[K, V]) *Evicted[K, V] {
	if e == nil {
		return nil
	}

	return &Evicted[K, V]{
		e.Key,
		e.Value,
	}
}

func (T *Tannin[K, V]) Set(key K, value V) *Evicted[K, V] {
	mcv := T.main.Victim(key)
	if mcv == nil {
		// If there is no victims in main cache - propagate incoming key to main cache
		return fromSlruEvicted(T.main.Set(key, value))
	}

	wcv := T.window.Set(key, value)
	if wcv == nil {
		return nil
	}

	if T.c.estimate(T.hash(wcv.Key)) > T.c.estimate(T.hash(*mcv)) {
		return fromSlruEvicted(T.main.Set(wcv.Key, wcv.Value))
	}

	return fromLruEvicted(wcv)
}

func (T *Tannin[K, V]) Len() int {
	return T.window.Len() + T.main.Len()
}

// Peek returns value for key (if key was in cache), but does not modify its recency.
func (T *Tannin[K, V]) Peek(key K) *V {
	if e := T.window.Peek(key); e != nil {
		return e
	}

	return T.main.Peek(key)
}

// Remove method removes entry associated with key and returns pointer to removed value (or nil if entry was not in cache).
func (T *Tannin[K, V]) Remove(key K) *V {
	if e := T.window.Remove(key); e != nil {
		return e
	}

	return T.main.Remove(key)
}

func NewParams[K comparable, V any](windowSize int, mainASize int, mainBSize int, samples int) *Tannin[K, V] {
	// https://mdlayher.com/blog/go-generics-draft-design-building-a-hashtable/#bonus-a-generic-hash-function
	var m interface{} = make(map[K]struct{})
	hf := (*mh)(*(*unsafe.Pointer)(unsafe.Pointer(&m))).hf

	return &Tannin[K, V]{
		window:  lru.New[K, V](windowSize),
		c:       *newCM4(windowSize + mainASize + mainBSize),
		main:    slru.NewParams[K, V](mainASize, mainBSize),
		hf:      hf,
		samples: samples,
	}
}

func New[K comparable, V any](size int, samples int) *Tannin[K, V] {
	windowSize := int(DefaultWindowRatio * float64(size))
	mainASize := int(DefaultSLRUSplit * float64(size-windowSize))
	return NewParams[K, V](
		windowSize,
		mainASize,
		size-windowSize-mainASize,
		samples,
	)
}

///////////////////////////
/// stolen from runtime ///
///////////////////////////

// mh is an inlined combination of runtime._type and runtime.maptype.
type mh struct {
	_  uintptr
	_  uintptr
	_  uint32
	_  uint8
	_  uint8
	_  uint8
	_  uint8
	_  func(unsafe.Pointer, unsafe.Pointer) bool
	_  *byte
	_  int32
	_  int32
	_  unsafe.Pointer
	_  unsafe.Pointer
	_  unsafe.Pointer
	hf func(unsafe.Pointer, uintptr) uintptr
}
