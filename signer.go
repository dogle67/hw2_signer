package main

import "sync"

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

	go crc(ch1, data)

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

}

//CombineResults combine results
func CombineResults(in, out chan interface{}) {

}
