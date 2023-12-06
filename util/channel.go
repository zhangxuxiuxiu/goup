package util

import (
	"context"
	"sync"
)

func CAny[T any](c <-chan T, pred func(T) bool, done chan struct{}) bool {
	for v := range c {
		if pred(v) {
			if done != nil {
				close(done)
			}
			return true
		}
	}
	return false
}

func CAll[T any](c <-chan T, pred func(T) bool, done chan struct{}) bool {
	return !CAny(c, Not(pred), done)
}

func CEqual[T comparable](t1, t2 <-chan T, d1, d2 chan struct{}) bool {
	defer func() {
		if d1 != nil {
			close(d1)
		}
		if d2 != nil {
			close(d2)
		}
	}()

	for s := range t1 {
		d, isOpen := <-t2
		if !isOpen || s != d {
			return false
		}
	}
	_, isOpen := <-t2
	return !isOpen
}

func CBatch[T any](c <-chan T, batchSize int) <-chan chan T {
	batches := make(chan chan T)

	go func() {
		size := 0
		var batch chan T
		for v := range c {
			size++
			if size == 1 {
				batch = make(chan T, 0)
				batches <- batch
			}
			batch <- v
			if size == batchSize {
				size = 0
				close(batch)
			}
		}
		if size != 0 {
			close(batch)
		}
		close(batches)
	}()

	return batches
}

func CUnique[T comparable](t1 <-chan T) <-chan T {
	t2 := make(chan T)

	go func() {
		first := true
		var last T
		for v := range t1 {
			if first {
				t2 <- v
				first = false
				last = v
			} else if v != last {
				t2 <- v
				last = v
			}
		}
		close(t2)
	}()

	return t2
}

// CFilter , Map, Reduce
func CFilter[T any](t1 <-chan T, pred func(T) bool) <-chan T {
	t2 := make(chan T)

	go func() {
		for v := range t1 {
			if pred(v) {
				t2 <- v
			}
		}
		close(t2)
	}()

	return t2
}

func CMap[T, U any](t <-chan T, fn func(T) U) <-chan U {
	u := make(chan U)

	go func() {
		for v := range t {
			u <- fn(v)
		}
		close(u)
	}()

	return u
}

func CFlatMap[T, U any](src <-chan T, fn func(T) []U) <-chan U {
	dest := make(chan U)
	go func() {
		for v := range src {
			if us := fn(v); len(us) != 0 {
				for _, d := range us {
					dest <- d
				}
			}
		}
		close(dest)
	}()
	return dest
}

// CRange params : end (,start (,step))
func CRange(args ...int) <-chan int {
	c := make(chan int, 1)

	end, start, step := 0, 0, 1
	switch len(args) {
	case 1:
		end = args[0]
	case 2:
		end = args[0]
		start = args[1]
	case 3:
		end = args[0]
		start = args[1]
		step = args[2]
	default:
		close(c)
		return c
	}

	go func() {
		for ; start < end; start = start + step {
			c <- start
		}
		close(c)
	}()

	return c
}

func WhenAll[T any](arrF []func() T, fn func(T) bool) <-chan []T {
	fnGen := func(_ context.Context, f func() T) func() T {
		return f
	}
	return GWhenAll(context.Background(), arrF, fnGen, fn)
}

func WhenAny[T any](arrF []func() T, fn func(T) bool) <-chan T {
	fnGen := func(_ context.Context, f func() T) func() T {
		return f
	}
	return GWhenAny(context.Background(), arrF, fnGen, fn)
}

func CWhenAll[T any](ctxt context.Context, arrFn []func(context.Context) T, fn func(T) bool) <-chan []T {
	fnGen := func(ctx context.Context, f func(context.Context) T) func() T {
		return func() T {
			return f(ctx)
		}
	}
	return GWhenAll(ctxt, arrFn, fnGen, fn)
}

func CWhenAny[T any](ctxt context.Context, arrFn []func(context.Context) T, pred func(T) bool) <-chan T {
	fnGen := func(ctx context.Context, f func(context.Context) T) func() T {
		return func() T {
			return f(ctx)
		}
	}
	return GWhenAny(ctxt, arrFn, fnGen, pred)
}

func GWhenAll[T, F any](ctxt context.Context, arrFn []F, fnGen func(context.Context, F) func() T, pred func(T) bool) <-chan []T {
	ctxt, cancel := context.WithCancel(ctxt)

	mapFn := func(f F) func() T {
		return fnGen(ctxt, f)
	}
	// construct result array
	consumer := func(c <-chan T, done chan struct{}) <-chan []T {
		ret := make(chan []T)
		go func() {
			defer close(ret)
			arr := make([]T, 0, len(arrFn))
			for v := range c {
				if !pred(v) {
					close(done)
					cancel()
					return
				}
				arr = append(arr, v)
			}
			ret <- arr
		}()
		return ret
	}

	return whenImpl(Map(arrFn, mapFn), consumer)
}

func GWhenAny[T, F any](ctxt context.Context, arrFn []F, fnGen func(context.Context, F) func() T, fn func(T) bool) <-chan T {
	ctxt, cancel := context.WithCancel(ctxt)

	mapFn := func(f F) func() T {
		return fnGen(ctxt, f)
	}
	// construct result array
	consumer := func(c <-chan T, done chan struct{}) <-chan T {
		ret := make(chan T)
		go func() {
			defer close(ret)
			for v := range c {
				if fn(v) {
					close(done)
					cancel()
					ret <- v
					return
				}
			}
		}()

		return ret
	}
	return whenImpl(Map(arrFn, mapFn), consumer)
}

// whenImpl ref:https://go101.org/article/channel-closing.html
func whenImpl[T, U any](arrFn []func() T, consumer func(<-chan T, chan struct{}) <-chan U) <-chan U {
	done := make(chan struct{})
	genProducer := func(f func() T) func(chan<- T) {
		return func(c chan<- T) {
			select {
			case <-done:
			case c <- f():
			}
		}
	}
	pipe := Parallelize(Map(arrFn, genProducer))
	return consumer(pipe, done)
}

// Parallelize https://go.dev/blog/pipelines
func Parallelize[T any](arrF []func(chan<- T)) <-chan T {
	c := make(chan T)
	wg := sync.WaitGroup{}
	wg.Add(len(arrF))

	for _, f := range arrF { // multiple senders
		go func(fn func(chan<- T)) {
			fn(c)
			wg.Done()
		}(f)
	}

	go func() {
		wg.Wait() // all senders finished
		close(c)  // no send to $c anymore, so it's safe to close it
	}()

	return c
}

func Async[T any](fn func() T) chan T {
	c := make(chan T)
	go func() {
		c <- fn()
		close(c)
	}()
	return c
}

func ToChan[T any](arr []T) chan T {
	c := make(chan T)
	go func() {
		for _, v := range arr {
			c <- v
		}
		close(c)
	}()
	return c
}
