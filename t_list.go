package redis

import (
	"container/list"
	"log"

	"github.com/pkg/errors"
)

var _ Item = (*List)(nil)

type List struct {
	inner *list.List
}

func NewListFromArr(arr []string) *List {
	list := NewList()
	list.LPush(arr...)
	return list
}

func NewList() *List {
	return &List{inner: list.New()}
}

func (l *List) Value() interface{} {
	return l.inner
}

func (l List) Type() uint64 {
	return ValueTypeList
}

func (l List) TypeFancy() string {
	return ValueTypeFancyList
}

func (l List) OnDelete(key string, db RedisDb) {
	log.Printf("Deleting %s with key %s from database ID %d\n", l.TypeFancy(), key, db.id)
}

// LLen returns number of elements.
func (l *List) LLen() int {
	return l.inner.Len()
}

// LPush returns the length of the list after the push operation.
func (l *List) LPush(values ...string) int {
	for _, v := range values {
		l.inner.PushFront(v)
	}
	return l.LLen()
}

// RPush returns the length of the list after the push operation.
func (l *List) RPush(values ...string) int {
	for _, v := range values {
		l.inner.PushBack(v)
	}
	return l.LLen()
}

// LInsert see redis doc
func (l *List) LInsert(isBefore bool, pivot, value string) int {
	for e := l.inner.Front(); e.Next() != nil; e = e.Next() {
		if getString(e) == pivot {
			if isBefore {
				l.inner.InsertBefore(value, e)
			} else {
				l.inner.InsertAfter(value, e)
			}
			return l.LLen()
		}
	}
	return -1
}

// RPop pops the front of the list and returns the value.
// Returns true if its valid, false otherwise.
func (l *List) LPop() (string, bool) {
	if e := l.inner.Front(); e == nil {
		return "", false
	} else {
		l.inner.Remove(e)
		return getString(e), true
	}
}

// RPop pops the back of the list and returns the value.
// Returns true if its valid, false otherwise.
func (l *List) RPop() (string, bool) {
	if e := l.inner.Back(); e == nil {
		return "", false
	} else {
		l.inner.Remove(e)
		return getString(e), true
	}
}

// LRem see redis doc
func (l *List) LRem(count int, value string) int {
	// count > 0: Remove elements equal to value moving from head to tail.
	// count < 0: Remove elements equal to value moving from tail to head.
	// count = 0: Remove all elements equal to value.
	var rem int
	if count >= 0 {
		for e := l.inner.Front(); e.Next() != nil; {
			if getString(e) == value {
				r := e
				e = e.Next()
				l.inner.Remove(r)
				rem++
				if count != 0 && rem == count {
					break
				}
			} else {
				e = e.Next()
			}
		}
	} else if count < 0 {
		count = abs(count)
		for e := l.inner.Back(); e.Prev() != nil; {
			if getString(e) == value {
				r := e
				e = e.Prev()
				l.inner.Remove(r)
				rem++
				if count != 0 && rem == count {
					break
				}
			} else {
				e = e.Prev()
			}
		}
	}
	return rem
}

// LSet see redis doc
func (l *List) LSet(index int, value string) error {
	e := atIndex(index, l.inner)
	if e == nil {
		return errors.New("index out of range")
	}
	e.Value = value
	return nil
}

// LIndex see redis doc
func (l *List) LIndex(index int) (string, bool) {
	e := atIndex(index, l.inner)
	if e == nil {
		return "", false
	}
	return getString(e), true
}

// LRange see redis doc
func (l *List) LRange(start int, end int) []string {
	values := make([]string, 0)
	// from index to index
	from, to := startEndIndexes(start, end, l.LLen())
	if from > to {
		return values
	}
	// get start element
	e := atIndex(from, l.inner)
	if e == nil { // shouldn't happen
		return values
	}
	// fill with values
	values = append(values, getString(e))
	for i := 0; i < to; i++ {
		e = e.Next()
		values = append(values, getString(e))
	}
	return values
}

// LTrim see redis docs - returns true if list is now emptied so the key can be deleted.
func (l *List) LTrim(start int, end int) bool {
	// from index to index
	from, to := startEndIndexes(start, end, l.LLen())
	if from > to {
		l.inner.Init()
		return true
	}
	// trim before
	if from > 0 {
		i := 0
		e := l.inner.Front()
		for e != nil && i < from {
			del := e
			e = e.Next()
			l.inner.Remove(del)
			i++
		}
	}
	// trim after
	if to < l.LLen() {
		i := l.LLen()
		e := l.inner.Back()
		for e != nil && i > to {
			del := e
			e = e.Prev()
			l.inner.Remove(del)
			i--
		}
	}
	return false
}

// TODO: For now we only store strings so this should be enough.
func (list *List) ForEachF(f func(a string)) {
	l := list.inner
	for e := l.Front(); e != nil; e = e.Next() {
		f(e.Value.(string))
	}
}

func startEndIndexes(start, end int, listLen int) (int, int) {
	if end > listLen-1 {
		end = listLen - 1
	}
	return toIndex(start, listLen), toIndex(end, listLen)
}

// atIndex finds element at given index or nil.
func atIndex(index int, list *list.List) *list.Element {
	index = toIndex(index, list.Len())
	e, i := list.Front(), 0
	for ; e.Next() != nil && i < index; i++ {
		if e.Next() == nil {
			return nil
		}
		e = e.Next()
	}
	return e
}

// Converts to real index.
//
// E.g. i=5, len=10 -> returns 5
//
// E.g. i=-1, len=10 -> returns 10
//
// E.g. i=-10, len=10 -> returns 0
//
// E.g. i=-3, len=10 -> returns 7
func toIndex(i int, len int) int {
	if i < 0 {
		if len+i > 0 {
			return len + i
		} else {
			return 0
		}
	}
	return i
}

// Value of a list element to string.
func getString(e *list.Element) string {
	v := e.Value.(string)
	return v
}

// Return positive x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
