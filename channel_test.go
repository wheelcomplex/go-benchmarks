package test

import (
	"runtime"
	"sync"
	"testing"

	"github.com/Workiva/go-datastructures/queue"
)

func BenchmarkChannelSPSC(b *testing.B) {
	ch := make(chan interface{}, 1)

	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			<-ch
		}
	}()

	for i := 0; i < b.N; i++ {
		ch <- `a`
	}
}

func BenchmarkChanRingBufferSPSC(b *testing.B) {
	q := queue.NewRingBuffer(1)

	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			q.Get()
		}
	}()

	for i := 0; i < b.N; i++ {
		q.Put(`a`)
	}
}

func BenchmarkChannelMPMC(b *testing.B) {
	ch := make(chan interface{}, 1)
	var wg sync.WaitGroup
	pn := runtime.GOMAXPROCS(-1)
	wg.Add(pn)

	b.ResetTimer()
	for i := 0; i < pn; i++ {
		go func() {
			for i := 0; i < b.N; i++ {
				<-ch
			}
		}()
	}

	for i := 0; i < pn; i++ {
		go func() {
			for i := 0; i < b.N; i++ {
				ch <- `a`
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkChanRingBufferMPMC(b *testing.B) {
	q := queue.NewRingBuffer(1)

	var wg sync.WaitGroup
	pn := runtime.GOMAXPROCS(-1)
	wg.Add(pn)

	b.ResetTimer()
	for i := 0; i < pn; i++ {
		go func() {
			go func() {
				for i := 0; i < b.N; i++ {
					q.Get()
				}
			}()

		}()
	}

	for i := 0; i < pn; i++ {
		go func() {
			for i := 0; i < b.N; i++ {
				q.Put(`a`)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
