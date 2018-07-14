package translate

import (
	"net/rpc"

	log "github.com/Sirupsen/logrus"
)

type ComfyWorker struct {
	client *rpc.Client
}

func newComfyWorker() ComfyWorker {
	w := ComfyWorker{}

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal("Failed to connect with translation service:", err)
	}

	w.client = client

	return w
}

func (w ComfyWorker) Process(payload interface{}) interface{} {
	return comfyTranslate(w.client, payload.(translateRequest))
}

// BlockUntilReady is called before each job is processed and must block the
// calling goroutine until the Worker is ready to process the next job.
func (w ComfyWorker) BlockUntilReady() {

}

// Interrupt is called when a job is cancelled. The worker is responsible
// for unblocking the Process implementation.
func (w ComfyWorker) Interrupt() {

}

// Terminate is called when a Worker is removed from the processing pool
// and is responsible for cleaning up any held resources.
func (w ComfyWorker) Terminate() {
	w.client.Close()
}
