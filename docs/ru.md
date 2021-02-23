# Lika Queue

Пакет включает в себя реализацию компонента для работы с очередями, а так же простые, но эффективные инструменты, которые
помогут просто и изящно построить многопоточные приложения.

Данный пакет отлично подойдет для работы в приложениях с высококонкурентной средой, так как удовлетворяет thread-safety концепции.

# Очереди (Queue)

```go
package main

import (
	"fmt"
	queue "github.com/lika-queue"
)

func main() {
	q := queue.NewQueue()
    q.Publish("my message")

    message := q.Consume()
    fmt.Printf("%+v\n", message)
    // &{Data:my message Retries:0 done:false queue:0xc00005e1b0 id:1}
}



```