# Lika Queue

Пакет включает в себя реализацию компонента для работы с очередями, а так же простые, но эффективные инструменты, которые
помогут просто и изящно построить многопоточные приложения.

# Очереди (Queue)

### Базовое использование
Если все сообщения прочитаны из очереди, далее вы будете получать nil

```go
package main

import (
    "fmt"
    queue "github.com/lika-queue"
)

func main() {
    q := queue.NewQueue(10)
    go q.Publish("my message 1")
    go q.Publish("my message 2")
    
    for {
        m := q.Consume()
    
        if m == nil {
            break
        }
    
        fmt.Printf("%+v\n", m)
    }
}
```

Если необходимо использовать очередь с блокировкой в ожидании новых сообщений, то можете для объекта Queue указать аттрибут
`Infinitie:true` и установить интервал времени в миллисикундах, в течении которого очередь будет перепроверять наличие новых
сообщений `Duration:100` или же создать очередь вызвав конструктор `q := queue.NewInfiniteQueue(100)`

```go
package main

import (
    "fmt"
    queue "github.com/lika-queue"
)

func main() {
    q := queue.NewQueue(10)
    q.Infinite = true
    
    go pushMessages(q, 1000)
    
    for {
        m := q.Consume()
    
        fmt.Printf("%+v\n", m)
    }
}

func pushMessages(q *queue.Queue, count int) {
    i := 0
    
    for i < count {
        q.Publish("rr")
        i++
    }
}
```

### Использование очередей в высококонкурентной среде
Очереди отлично подходят для работы в высококонкурентной среде, когда происходит параллельные процессы на чтение и запись.

```go
package main

import (
    "fmt"
    queue "github.com/lika-queue"
)

func main() {
    q := queue.NewInfiniteQueue(300)
    
    go pushMessages(q, 1000)
    go pushMessages(q, 1000)
    
    go readMessages(q)
    go readMessages(q)
}

func pushMessages(q *queue.Queue, count int) {
    i := 0
    
    for i < count {
        q.Publish("rr")
        i++
    }
}

func readMessages(q *queue.Queue) {
    for {
        m := q.Consume()
        
        // my code
    }
}
```

# Worker

Зачастую использование очередей предполагает под собой чтение данных из нескольких параллельных потоков и выполнением идентичных
операций в каждом из них, для этого в пакете есть готовые инструменты для упрощения данной задачи.

Компонент Worker предназначен для создания динамического количества потоков в которых запускается выполнение переданной вами функции,
а так же предоставляет возможность синхронизировать их работу через стандартные библиотеки Golang.

В пакете есть две реализации воркеров, стандартный воркер `Worker`, который никак не привязан к очередям и можем использоваться
для различных целей, где нужен пулл горутин с одновременной синхранизацией, а так же `QueueWorker` который помогает читать данные из очереди 
и одновременно работает с базовым воркером

### Classic worker implementation
```go
package main

import (
    "fmt"
    queue "github.com/lika-queue"
    "time"
)

func main() {
    w := queue.Worker{
        Threads: 2, // threads count
        SyncThreads: true, // if set false, worker will not wait for threads to finish
        Callable: printCounter,
    }
    w.Run()
}

func printCounter() {
    i := 4
    
    for i > 0 {
        fmt.Printf("%+v  \n", i)
        i--
        time.Sleep(1000 * time.Millisecond)
    }
}
```

### Queue based worker implementation

Обратите внимание на то, что функцию вы теперь передаете компоненту `QueueWorker`, и функция должна удовлетворять интерфейс принимая
на вход ссылку объекта `Message`

```go
package main

import (
    "fmt"
    queue "github.com/lika-queue"
    "time"
)

func main() {
    q := queue.NewQueue()
    
    pushMessages(q, 10000)
    
    w := queue.QueueWorker{
        Queue: q,
        Callable: printMessage, // NOTE: your callback function you shoud set to QueueWorker, not to Worker
        Worker: &queue.Worker{
            Threads:     5,
            SyncThreads: true,
        },
    }
    
    w.Run()
}

func printMessage(m *queue.Message)  {
    fmt.Printf("%+v\n", m)
}

```

А теперь давайте посмотрим на тот же самый код, но без использования воркеров

```go
package main

import (
    "fmt"
    queue "github.com/lika-queue"
    "sync"
    "time"
)

func main() {
    q := queue.NewQueue()
    
    pushMessages(q, 10000)
    gorutinesCount := 5
    
    group := &sync.WaitGroup{}
    group.Add(gorutinesCount)
    
    i := 0
    for i < gorutinesCount {
        go consumeMessages(q, group)
        i++
    }

    group.Wait()
}

func consumeMessages(q *queue.Queue, g *sync.WaitGroup) {
    defer g.Done()
    
    for {
        m := q.Consume()
    
        if m == nil {
            break
        }
    
        fmt.Printf("%+v\n", m)
    }
}
```