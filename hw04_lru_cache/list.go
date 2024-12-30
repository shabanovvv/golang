package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
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
	newFront := &ListItem{Value: v}

	if l.front == nil {
		l.front = newFront
		l.back = newFront
	} else {
		newFront.Next = l.front
		l.front.Prev = newFront
		newFront.Prev = nil
		l.front = newFront
	}

	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBack := &ListItem{Value: v}
	if l.back == nil {
		l.back = newBack
		l.front = newBack
	} else {
		newBack.Prev = l.back
		l.back.Next = newBack
		newBack.Next = nil
		l.back = newBack
	}
	l.len++

	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	if l.front == nil {
		l.back = nil
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.front {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
