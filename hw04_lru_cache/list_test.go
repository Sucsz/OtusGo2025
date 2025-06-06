package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// list_test.go:9:6: Function 'TestList' is too long (187 > 150) (funlen)
//
//nolint:funlen
func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("TestLen", func(t *testing.T) {
		l := NewList()
		require.Equal(t, 0, l.Len())

		l.PushFront("10")
		require.Equal(t, 1, l.Len())
		l.PushBack(10)
		require.Equal(t, 2, l.Len())

		for _, v := range [...]int{10, 20, 30, 70, 80} {
			l.PushBack(v)
		}
		require.Equal(t, 7, l.Len())

		l.Remove(l.Front())
		require.Equal(t, 6, l.Len())

		lLen := l.Len()
		for i := 0; i < lLen; i++ {
			l.Remove(l.Back())
		}
		require.Equal(t, 0, l.Len())
	})

	t.Run("TestFront", func(t *testing.T) {
		l := NewList()
		require.Nil(t, l.Front())

		l.PushFront(10)
		require.Equal(t, 10, l.Front().Value)

		l.PushFront("15")
		require.Equal(t, "15", l.Front().Value)

		l.Remove(l.Front())
		require.Equal(t, 10, l.Front().Value)
	})

	t.Run("TestBack", func(t *testing.T) {
		l := NewList()
		require.Nil(t, l.Back())

		l.PushFront(10)
		require.Equal(t, 10, l.Back().Value)

		l.PushBack("15")
		require.Equal(t, "15", l.Back().Value)

		l.Remove(l.Back())
		require.Equal(t, 10, l.Front().Value)
	})

	t.Run("PushFront", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 1, l.Len())

		l.PushFront(20)
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 2, l.Len())
		require.Equal(t, 20, l.Front().Value)

		// Связи
		require.Nil(t, l.Front().Prev)
		require.Equal(t, 10, l.Front().Next.Value)
		require.Equal(t, 20, l.Front().Next.Prev.Value)

		l.PushFront("abc")
		require.Equal(t, "abc", l.Front().Value)
		require.Equal(t, 3, l.Len())
	})

	t.Run("PushBack", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 1, l.Len())

		l.PushBack(20)
		require.Equal(t, 20, l.Back().Value)
		require.Equal(t, 2, l.Len())

		// Связи
		require.Nil(t, l.Back().Next)
		require.Equal(t, 10, l.Back().Prev.Value)
		require.Equal(t, 20, l.Front().Next.Value)

		l.PushBack("abc")
		require.Equal(t, "abc", l.Back().Value)
		require.Equal(t, 3, l.Len())
	})

	t.Run("TestRemove", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)
		l.Remove(l.Front())
		require.Equal(t, 0, l.Len())

		l.PushFront(3)
		l.PushFront(2)
		l.PushFront(1)
		l.Remove(l.Front().Next)
		require.Equal(t, 2, l.Len())

		// Связи
		require.Equal(t, 3, l.Front().Next.Value)
		require.Equal(t, 1, l.Front().Next.Prev.Value)
	})

	t.Run("MoveToFront", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)
		l.MoveToFront(l.Front())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 1, l.Len())

		l.PushBack(20)
		l.MoveToFront(l.Back())
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 2, l.Len())

		l.PushBack(30)
		l.MoveToFront(l.Front().Next)
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 3, l.Len())

		// Связи
		require.Equal(t, 10, l.Front().Next.Prev.Value)
		require.Equal(t, 30, l.Front().Next.Next.Value)
	})

	t.Run("MoveToBack", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)
		l.MoveToBack(l.Front())
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 1, l.Len())

		l.PushFront(20)
		l.MoveToBack(l.Front())
		require.Equal(t, 20, l.Back().Value)
		require.Equal(t, 2, l.Len())

		l.PushFront(30)
		l.MoveToBack(l.Front().Next)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 3, l.Len())

		// Связи
		require.Equal(t, 20, l.Front().Next.Value)
		require.Equal(t, 10, l.Front().Next.Next.Value)
	})
}
