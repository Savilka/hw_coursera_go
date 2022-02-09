package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

func ExecutePipeline(jobs ...job) {

	wg := &sync.WaitGroup{}

	in := make(chan interface{})
	for _, job := range jobs {

		wg.Add(1)

		out := make(chan interface{})
		go jobMaker(wg, job, in, out)
		in = out
	}

	wg.Wait()
}

func jobMaker(wg *sync.WaitGroup, job job, in, out chan interface{}) {

	job(in, out)
	defer wg.Done()
	defer close(out)
}

func SingleHash(in, out chan interface{}) {

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for data := range in {

		wg.Add(1)

		go func(data interface{}, out chan interface{}, wg *sync.WaitGroup, mu *sync.Mutex) {

			defer wg.Done()
			dataString := strconv.Itoa(data.(int))
			mu.Lock()
			md5 := DataSignerMd5(dataString)
			mu.Unlock()

			crc32chan := make(chan string)
			go func(data string, crc32chan chan string) {
				crc32chan <- DataSignerCrc32(data)
			}(dataString, crc32chan)

			crc32Md5chan := make(chan string)
			go func(md5 string, crc32Md5chan chan string) {
				crc32Md5chan <- DataSignerCrc32(md5)
			}(md5, crc32Md5chan)

			result := <-crc32chan + "~" + <-crc32Md5chan
			out <- result
		}(data, out, wg, mu)
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for data := range in {

		wg.Add(1)

		go func(data interface{}, out chan interface{}, wg *sync.WaitGroup) {

			defer wg.Done()

			wgCrc32 := &sync.WaitGroup{}
			mu := &sync.Mutex{}

			multiHash := make([]string, 6)

			for i := 0; i < 6; i++ {
				wgCrc32.Add(1)
				data := strconv.Itoa(i) + data.(string)
				go func(wgCrc32 *sync.WaitGroup, mu *sync.Mutex, data string, multiHash []string, i int) {
					defer wgCrc32.Done()
					data = DataSignerCrc32(data)
					mu.Lock()
					multiHash[i] = data
					mu.Unlock()
				}(wgCrc32, mu, data, multiHash, i)

			}

			wgCrc32.Wait()
			result := strings.Join(multiHash, "")
			out <- result

		}(data, out, wg)
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var combineResult []string
	for data := range in {
		data := fmt.Sprint(data)
		combineResult = append(combineResult, data)
	}
	sort.Strings(combineResult)
	out <- strings.Join(combineResult, "_")
}
