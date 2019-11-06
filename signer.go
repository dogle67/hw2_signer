package main

import (
	"sort"
	"strconv"
	"sync"
)

//ExecutePipeline execute piplines
func ExecutePipeline(jobs ...job) {

}

//SingleHash single hash
func SingleHash(in, out chan interface{}) {
	for i := range in {

		data := i.(string)
		go singleHash(out, data)
	}
}

func singleHash(out chan interface{}, data string) {

	ch1 := make(chan string)
	ch2 := make(chan string)

	go crc(ch1, data)

	go mdfive(ch2, data)
	cm := DataSignerCrc32(<-ch2)

	select {
	case v := <-ch1:
		out <- v + "~" + cm

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

	for i := range in {
		data := i.(string)
		go multiHash(out, data)
	}
}

func multiHash(out chan interface{}, data string) {
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
