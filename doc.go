/*
Package gearman provides a thread-safe Gearman client

Example

Here's an example program that submits a job to Gearman and listens for events from that job:

	package main

	import(
		"github.com/Clever/gearman"
		"github.com/Clever/gearman/job"
	)

	func main() {
		client, err := gearman.NewClient("tcp4", "localhost:4730")
		if err != nil {
			panic(err)
		}

		j, err := client.Submit("reverse", []byte("hello world!"))
		if err != nil {
			panic(err)
		}
		// Warnings are impossible for this worker, so we don't need to listen for warnings.
		// If they were and we didn't listen for warnings, we would block parsing from the server.
		for data := range j.Data() {
			println(data) // !dlrow olleh
		}
		println(j.State()) // job.State.Completed
	}
*/
package gearman
