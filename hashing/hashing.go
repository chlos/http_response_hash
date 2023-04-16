package hashing

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Hashing struct {
	ParallelLimit int
	URLs          []string
	hashCh        chan string
}

func NewHashing(limit int, urls []string) (*Hashing, error) {
	if limit < 1 {
		return nil, errors.New("parallel limit must be >= 1")
	}

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

	body, err := getHTTPBody(url)
	if err != nil {
		fmt.Printf("Hashing %s: %v\n", url, err)
		return
	}
	responseHash := getMD5Hash(body)

	h.hashCh <- responseHash
}

func getHTTPBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("running GET request: %w", err)
	}
	defer resp.Body.Close()

	var body bytes.Buffer
	_, err = io.Copy(&body, resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("reading response body: %w", err)
	}

	return body.Bytes(), nil
}

func getMD5Hash(body []byte) string {
	hasher := md5.New()
	hasher.Write(body)
	return hex.EncodeToString(hasher.Sum(nil))
}
