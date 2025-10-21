package utils

import (
	"log"
	"sync"
)

// Job represents a processing job
type Job struct {
	Type      string // "image", "video", "audio"
	FilePath  string
	UploadDir string
	FileName  string
}

// WorkerPool manages concurrent processing jobs
type WorkerPool struct {
	jobQueue    chan Job
	workerCount int
	wg          sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	pool := &WorkerPool{
		jobQueue:    make(chan Job, queueSize),
		workerCount: workerCount,
	}
	pool.start()
	return pool
}

// start initializes workers
func (p *WorkerPool) start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// worker processes jobs from the queue
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	for job := range p.jobQueue {
		log.Printf("[Worker %d] Processing %s: %s", id, job.Type, job.FileName)

		switch job.Type {
		case "image":
			_, err := ResizeImage(job.FilePath, job.UploadDir, job.FileName)
			if err != nil {
				log.Printf("[Worker %d] Image processing error: %v", id, err)
			} else {
				log.Printf("[Worker %d] Image processed successfully: %s", id, job.FileName)
			}

		case "video":
			_, err := ProcessVideo(job.FilePath, job.UploadDir, job.FileName)
			if err != nil {
				log.Printf("[Worker %d] Video processing error: %v", id, err)
			} else {
				log.Printf("[Worker %d] Video processed successfully: %s", id, job.FileName)
			}

		case "audio":
			_, err := ProcessAudio(job.FilePath, job.UploadDir, job.FileName)
			if err != nil {
				log.Printf("[Worker %d] Audio processing error: %v", id, err)
			} else {
				log.Printf("[Worker %d] Audio processed successfully: %s", id, job.FileName)
			}
		}
	}
}

// Submit adds a job to the queue
func (p *WorkerPool) Submit(job Job) {
	p.jobQueue <- job
}

// Shutdown gracefully stops the worker pool
func (p *WorkerPool) Shutdown() {
	close(p.jobQueue)
	p.wg.Wait()
}

// Global worker pool instance
var globalWorkerPool *WorkerPool
var once sync.Once

// GetWorkerPool returns the global worker pool instance (singleton)
func GetWorkerPool() *WorkerPool {
	once.Do(func() {
		// Initialize with 4 workers and queue size of 100
		// This means max 4 concurrent processing jobs, with up to 100 waiting in queue
		globalWorkerPool = NewWorkerPool(4, 100)
	})
	return globalWorkerPool
}
