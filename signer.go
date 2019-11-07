package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

//ExecutePipeline execute piplines
func ExecutePipeline(jobs ...job) {

	in := make(chan interface{}, 10)
	out := make(chan interface{}, 10)

	wg := &sync.WaitGroup{}
	for i := range jobs {
		wg.Add(1)

		go worker(wg, jobs[i], in, out)

		in = out
		out = make(chan interface{}, 10)
	}

	wg.Wait()
}

func worker(wg *sync.WaitGroup, j job, in, out chan interface{}) {
	defer wg.Done()
	defer close(out)
	j(in, out)

}

//SingleHash single hash
func SingleHash(in, out chan interface{}) {

	wg := &sync.WaitGroup{}
	for i := range in {
		data := i.(int)
		str := fmt.Sprint(data)

		wg.Add(1)
		go singleHash(wg, out, str)
	}

	wg.Wait()
}

func singleHash(wg *sync.WaitGroup, out chan interface{}, data string) {
	defer wg.Done()
	ch1 := make(chan string)
	ch2 := make(chan string)

	go crc(ch1, data)

	go mdfive(ch2, data)
	cm := DataSignerCrc32(<-ch2)

	select {
	case v := <-ch1:
		rsp := v + "~" + cm
		fmt.Println(rsp)
		out <- rsp

	}
}

func crc(ch chan string, data string) {

	rsp := DataSignerCrc32(data)
	ch <- rsp
}

var mux sync.Mutex

func mdfive(ch chan string, data string) {
	mux.Lock()
	defer mux.Unlock()
	rsp := DataSignerMd5(data)
	ch <- rsp
}

//MultiHash multi hash
func MultiHash(in, out chan interface{}) {

	wg := &sync.WaitGroup{}
	for i := range in {
		data := i.(string)
		wg.Add(1)
		go multiHash(wg, out, data)
	}

	wg.Wait()
}

func multiHash(wg *sync.WaitGroup, out chan interface{}, data string) {
	defer wg.Done()
	chanels := make([]chan string, 6)
	for i := range chanels {
		chanels[i] = make(chan string)
		go crc(chanels[i], strconv.FormatInt(int64(i), 10)+data)
	}

	rsp := ""
	for i := range chanels {
		rsp += <-chanels[i]
	}

	out <- rsp
}

//CombineResults combine results
func CombineResults(in, out chan interface{}) {

	strs := make([]string, 0)

	for i := range in {
		str := i.(string)
		strs = append(strs, str)
	}

	sort.Strings(strs)
	rsp := ""
	if len(strs) > 0 {
		rsp = strs[0]
		for i := 1; i < len(strs); i++ {
			rsp += "_" + strs[i]
		}
	}

	out <- rsp
}
