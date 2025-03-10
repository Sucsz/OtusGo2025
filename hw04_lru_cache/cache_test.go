package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// -.cache_test.go:12:6: Function 'TestCache' is too long (186 > 150) (funlen)
//
//nolint:funlen
func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a", 1)
		val, ok := c.Get("a")
		require.True(t, ok)
		require.Equal(t, 1, val)

		c.Set("b", 2)
		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		c.Set("c", 3)
		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)

		c.Clear()
		val, ok = c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("b")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("c")
		require.False(t, ok)
		require.Nil(t, val)

		// Проверка, что можно снова пользоваться кэшем после очистки
		c.Set("d", 4)
		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 4, val)

		c.Set("e", 5)
		val, ok = c.Get("e")
		require.True(t, ok)
		require.Equal(t, 5, val)

		c.Set("f", 6)
		val, ok = c.Get("f")
		require.True(t, ok)
		require.Equal(t, 6, val)
	})

	t.Run("CapacityElementEjection", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		val, ok := c.Get("a")
		require.True(t, ok)
		require.Equal(t, 1, val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)

		// вытолкнули а
		c.Set("d", 4)
		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 4, val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("LastUsedElementEjection", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		// Получаем значения
		val, ok := c.Get("a")
		require.True(t, ok)
		require.Equal(t, 1, val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = c.Get("a")
		require.True(t, ok)
		require.Equal(t, 1, val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		// Меняем значения
		c.Set("a", 10)
		val, ok = c.Get("a")
		require.True(t, ok)
		require.Equal(t, 10, val)

		c.Set("b", 20)
		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 20, val)

		c.Set("c", 30)
		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 30, val)

		// Получаем значение
		val, ok = c.Get("a")
		require.True(t, ok)
		require.Equal(t, 10, val)

		// Добавляем 4е значение
		c.Set("d", 40)
		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 40, val)

		// Смотрим на 3 последних упоминания = d, a, c -> вытолкнули b
		val, ok = c.Get("b")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()

	t.Run("TestRoutineSet", func(t *testing.T) {
		capas := 10
		c := NewCache(capas)

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			for i := 0; i < capas/2; i++ {
				c.Set(Key(strconv.Itoa(i)), i)
			}
		}()

		go func() {
			defer wg.Done()
			for i := capas / 2; i < capas; i++ {
				c.Set(Key(strconv.Itoa(i)), i)
			}
		}()

		wg.Wait()
		for i := 0; i < capas; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.True(t, ok)
			require.Equal(t, i, val)
		}
	})
}
