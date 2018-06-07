/*
@Time : 2018/6/7 16:24 
@Author : yinsaki
@File : gopool
*/
package lib

import (
	"sync"
)

type Job func()

type worker struct {
	workerPool chan *worker
	jobQueue chan Job
	stop chan struct{}
}

type dispatcher struct {
	workerPool chan*worker
	jobQueue chan Job
	stop chan struct{}
}

type  Pool struct {
	dispatch *dispatcher
	wg sync.WaitGroup
	enableWaitForAll bool
}

func (this *Pool)AddJob(job Job) {
	if this.enableWaitForAll {
		this.wg.Add(1)
	}

	this.dispatch.jobQueue <- func() {
		job()
		if this.enableWaitForAll {
			this.wg.Done()
		}
	}
}

func (this *Pool)StopAll() {
	this.dispatch.stop <- struct{}{} // 阻塞，等dispatch协程调度worker stop后
	<- this.dispatch.stop
}

func (this *Pool)Start() *Pool {
	for i:= 0; i < cap(this.dispatch.workerPool); i ++ {
		worker := newWorker(this.dispatch.workerPool)
		go worker.start()
	}

	go this.dispatch.dispatch()
	return this
}

func newDispatcher(workerPool chan* worker, jobQueue chan Job) *dispatcher {
	return &dispatcher{workerPool:workerPool, jobQueue:jobQueue, stop:make(chan struct{})}
}

func NewPool(workNum, jobNum int) *Pool {
	workers := make(chan *worker,workNum)
	jobs := make(chan Job, jobNum)

	pool := &Pool{
		dispatch:newDispatcher(workers, jobs),
		enableWaitForAll:false,
	}
	return pool
}

func (dis *dispatcher)dispatch() {
	for {
		select {
		case job:= <- dis.jobQueue:
			worker := <- dis.workerPool
			worker.jobQueue <- job
		case <-dis.stop:
			for i:= 0; i < cap(dis.workerPool);i++ {
				worker := <- dis.workerPool
				worker.stop <- struct{}{}
				<- worker.stop
			}

			dis.stop <- struct{}{}
			return
		}
	}
}
func newWorker(workerPool chan* worker) *worker {
	return &worker{
		workerPool:workerPool,
		jobQueue:make(chan Job),
		stop:make(chan struct{}),
	}
}

func (this *worker)start() {
	for {
		this.workerPool <- this // 这里可以确定哪个协程拿到了消息
		select {
		case job:= <-this.jobQueue:
			job()
		case <- this.stop:
			this.stop <- struct{}{}
			return
		}
	}
}

func (this *Pool)WaitForAll() {
	if this.enableWaitForAll {
		this.wg.Wait()
	}
}

func (this *Pool)EnableWaitForAll(enable bool) *Pool {
	this.enableWaitForAll = enable
	return this
}