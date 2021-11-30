# qp
<p align="center"><img src="./.github/qp.png" width="600" /></p>

## What is qp?
qp is a queue-processor designed with DX in mind. It basically accepts a queue, accepts your business logic then runs your business logic against all the incoming messages without you having to worry about ack'ing messages, backing off or pushing metrics to your metrics server.

## Why not use the sqs SDK directly?
We can, the issue is that you want the provider to be abstract so that your code runs against various things. Also as a developers, we don't really like handling boring bits.

## More about qp
Imagine a business requirement, where you need to process incoming data, and if it fails, you want to backoff exponentially, and on top of that, you want to try the same message again later if it fails.
Implementation of this becomes a little messy while handling queue and backing off (your queue provider probably handles retries and dead letter queue).

qp basically let's you write 1 callback which receives your message, and on your callback you write your business logic, and return error IF there's an error. qp will take care of the following:

- backing off exponentially
- retrying when failure
- graceful sigint handling
- reporting metrices

## Catch
As the time of this writing, there's only 1 queue provider which is `SNS` BUT we'll probably come up with others very soon, but it's not that difficult to write a driver, it's an interface with few methods.

It also does exponential back off which might not be what you want, we'll probably allow this to be changed later on in the future.

## Code Examples

```go
package main
import (
	"github.com/honestbank/qp/queue"
	"github.com/honestbank/qp"
)
func main() {
	q, _ := queue.NewSQSQueue("your-queue-name")
	procesor := qp.NewJob(q)
	procesor.SetWorker(processMessage)
	procesor.OnResult(handlers.ReportToPrometheus()) // catch with this is that you need a prometheus push gateway (or implement your own)
	procesor.Start()
}
```
