package stack

import (
	"errors"
	"sync"
)

type Item interface{}

type ItemStack struct {
	items []Item //栈的元素

	stackLen int //栈的长度

	lock sync.RWMutex //锁
}

// 创建栈
func New() *ItemStack {
	return &ItemStack{items: make([]Item, 0), stackLen: 0}
}

// 入栈
func (s *ItemStack) Push(t Item) Item {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = append(s.items, t)

	s.stackLen++

	return t
}

// 出栈
func (s *ItemStack) Pop() (Item, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.stackLen <= 0 {
		return nil, errors.New("STACK IS NULL")
	}

	item := s.items[s.stackLen-1]
	s.items = s.items[:s.stackLen-1]
	s.stackLen--

	return item, nil
}

//获取栈的长度
func (s *ItemStack) GetLen() int {
	return s.stackLen
}
