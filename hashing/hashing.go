package hashing

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
)

type Hashing struct {
	ParallelLimit int
	URLs          []string
	hashCh        chan string
}

func NewHashing(limit int, urls []string) (*Hashing, error) {
	return &Hashing{
		ParallelLimit: limit,
		URLs:          urls,
	}, nil
}

func (h *Hashing) Start() {
	// run the pool of goroutines
	var wg sync.WaitGroup
	wg.Add(len(h.URLs))
	waitPoolCh := make(chan struct{}, h.ParallelLimit) // limit the number of parallel goroutines
	h.hashCh = make(chan string, h.ParallelLimit)
	go func() {
		for _, url := range h.URLs {
			waitPoolCh <- struct{}{} // loop is blocked if we have reached max num of running goroutines
			go h.hashURL(url, waitPoolCh, &wg)
		}
	}()

	// wait for finishing
	go func() {
		wg.Wait()
		close(h.hashCh)
	}()
}

func (h *Hashing) Print() {
	for hash := range h.hashCh {
		fmt.Println(hash)
	}
}

func (h *Hashing) hashURL(url string, waitCh <-chan struct{}, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		<-waitCh
	}() // done, free up space in the pool waiting channel for the next goroutine

	// FIXME: do something
	responseHash := getMD5Hash(url)

	h.hashCh <- responseHash
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
