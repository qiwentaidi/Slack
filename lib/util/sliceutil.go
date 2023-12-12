package util

import (
	"sync"
)

type (
	SafeSlice struct {
		sync.RWMutex
		items []Items
	}

	Items struct {
		num   int
		value any
	}
)

func (ss *SafeSlice) Append(item any) {
	ss.Lock()
	defer ss.Unlock()

	ss.items = append(ss.items, Items{value: item})
}

func (ss *SafeSlice) Len() int {
	ss.RLock()
	defer ss.RUnlock()

	return len(ss.items)
}

func (ss *SafeSlice) Key(item any) int {
	ss.RLock()
	defer ss.RUnlock()

	for i, v := range ss.items {
		if v.value == item {
			return i
		}
	}

	return -1
}

func (ss *SafeSlice) Get(index int) any {
	ss.RLock()
	defer ss.RUnlock()

	return ss.items[index].value
}

func (ss *SafeSlice) Update(index int, item any) {
	ss.Lock()
	defer ss.Unlock()

	ss.items[index].value = item
}

func (ss *SafeSlice) List() []any {
	ss.RLock()
	defer ss.RUnlock()

	r := []any{}
	for _, v := range ss.items {
		r = append(r, v.value)
	}

	return r
}

func (ss *SafeSlice) Iter() chan any {
	ss.Lock()

	out := make(chan any)

	go func() {
		defer close(out)
		defer ss.Unlock()

		for _, item := range ss.items {
			out <- item.value
		}
	}()

	return out
}

func (ss *SafeSlice) Num(item any) int {
	ss.RLock()
	defer ss.RUnlock()

	for _, v := range ss.items {
		if v.value == item {
			return v.num
		}
	}

	return 0
}

func (ss *SafeSlice) UpdateNum(item any, num int) {
	ss.Lock()
	defer ss.Unlock()

	for i, v := range ss.items {
		if v.value == item {
			ss.items[i].num += num
			return
		}
	}
}

func (ss *SafeSlice) ResetNum(item any) {
	ss.Lock()
	defer ss.Unlock()

	for i, v := range ss.items {
		if v.value == item {
			ss.items[i].num = 0
			return
		}
	}
}

func (ss *SafeSlice) SetNum(item any, num int) {
	ss.Lock()
	defer ss.Unlock()

	for i, v := range ss.items {
		if v.value == item {
			ss.items[i].num = num
			return
		}
	}
}
