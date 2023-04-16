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

var (
	// getHTTPBody defines a function which should be used for getting HTTP body (production/test mock/etc).
	getHTTPBody = _getHTTPBody
)

// urlHash is a struct where the result of each URL hashing is stored.
type urlHash struct {
	// url is the current website's URL.
	url string
	// hash is a MD5 hash of the current URL's HTTP response body.
	hash string
}

// Hashing makes http requests and prints the address of the request along with the MD5 hash of the response.
type Hashing struct {
	// ParallelLimit is a value for the number of parallel requests.
	ParallelLimit int
	// URLs is a slice of URLs of websites which should be hashed.
	URLs []string
	// hashCh is a channel for urlHash results.
	hashCh chan urlHash
}

// NewHashing is a constructor of Hashing.
// It also checks the parallel limit and returns error in case of a wrong limit.
func NewHashing(limit int, urls []string) (*Hashing, error) {
	if limit < 1 {
		return nil, errors.New("parallel limit must be >= 1")
	}

	return &Hashing{
		ParallelLimit: limit,
		URLs:          urls,
	}, nil
}

// Start makes HTTP requests and calculate responses' HD5 hashes. It sends the result to the hashCh channel.
func (h *Hashing) Start() {
	// run the pool of goroutines
	var wg sync.WaitGroup
	wg.Add(len(h.URLs))
	waitPoolCh := make(chan struct{}, h.ParallelLimit) // limit the number of parallel goroutines
	h.hashCh = make(chan urlHash, h.ParallelLimit)
	go func() {
		for _, url := range h.URLs {
			waitPoolCh <- struct{}{} // loop is blocked if we have reached max num of running goroutines
			go h.hashHTTPBody(url, waitPoolCh, &wg)
		}
	}()

	// wait for finishing
	go func() {
		wg.Wait()
		close(h.hashCh)
	}()
}

// Print prints the address of the request along with the MD5 hash of the response.
func (h *Hashing) Print() {
	for item := range h.hashCh {
		fmt.Println(item.url, item.hash)
	}
}

// hashHTTPBody makes a HTTP request, calculates the response's hash and sends it to hashCh channel.
func (h *Hashing) hashHTTPBody(url string, waitCh <-chan struct{}, wg *sync.WaitGroup) {
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

	h.hashCh <- urlHash{url: url, hash: responseHash}
}

// _getHTTPBody makes a HTTP request using the specified URL and returns the response's body.
func _getHTTPBody(url string) ([]byte, error) {
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

// getMD5Hash calculates and returns the MD5 hash of a body.
func getMD5Hash(body []byte) string {
	hasher := md5.New()
	hasher.Write(body)
	return hex.EncodeToString(hasher.Sum(nil))
}
