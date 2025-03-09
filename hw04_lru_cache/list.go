package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	MoveToBack(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	node := &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}

	if l.len == 0 {
		l.back = node
	} else {
		l.front.Prev = node
	}
	l.front = node

	l.len++
	return node
}

func (l *list) PushBack(v interface{}) *ListItem {
	node := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}

	if l.len == 0 {
		l.front = node
	} else {
		l.back.Next = node
	}
	l.back = node

	l.len++
	return node
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if l.front == i {
		l.front = i.Next
	}
	if l.back == i {
		l.back = i.Prev
	}

	l.len--
	if l.len == 0 {
		l.front = nil
		l.back = nil
	}
}

func (l *list) MoveToFront(i *ListItem) {
	/*
		Утечка в памяти, создастся i` в таком случае
		l.Remove(i)
		l.PushFront(i.Value)
	*/
	if l.front == i {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	i.Prev = nil
	i.Next = l.front
	if l.front != nil {
		l.front.Prev = i
	}

	l.front = i
}

func (l *list) MoveToBack(i *ListItem) {
	/*
		Утечка в памяти, создастся i` в таком случае
		l.Remove(i)
		l.PushBack(i.Value)
	*/
	if l.back == i {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i.Next = nil
	i.Prev = l.back
	if l.back != nil {
		l.back.Next = i
	}
	l.back = i
}

func NewList() List {
	return new(list)
}
