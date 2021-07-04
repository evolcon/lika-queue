# Lika Queue

Пакет включает в себя реализацию компонента для работы с очередями, а так же простые, но эффективные инструменты, которые
помогут просто и изящно построить многопоточные приложения.

# Базовый компонент

Базовый компонент очередей используется для хранения и управления несколькими брокерами сообщений, а так же является промежуточным
звеном между пользователем и конкретным брокером.

Это позволяет в свою очередь иметь возможность использовать разные типы очередей для разных целей, и при необходимости 
заменять одну на другую, не задумываясь о деталях реализации.

```go

package main

import (
	"fmt"
	queue "github.com/lika_queue"
	MemoryBroker "github.com/lika_queue/brokers/memory"
)

func main() {
	var queueComponent queue.QueueInterface
	var broker queue.BrokerInterface

	queueComponent = queue.New()
	broker = MemoryBroker.New(10000)

	queueComponent.Add("main", broker)

	myMessage1 := "my string"
	myMessage2 := map[string]interface{} {
		"key1": "Vladimir",
		"key2": map[string]string {
			"name": "Dzhamshud",
			"surname": "Moskvich",
		},
	}

	_ = queueComponent.Publish("testQueue", myMessage1, nil)
	_ = queueComponent.Publish("testQueue", myMessage2, nil)

	for {
		message, _ := queueComponent.Consume("testQueue", nil)

		if message != nil {
			fmt.Println(message.GetData())
		} else {
			break
		}
	}
}


```


### Использование очередей в высококонкурентной среде
Очереди отлично подходят для работы в высококонкурентной среде, когда происходит параллельные процессы на чтение и запись.

```go

package main

import (
	"fmt"
	queue "github.com/lika_queue"
	MemoryBroker "github.com/lika_queue/brokers/memory"
	"time"
)

func main() {
	queueComponent := queue.New()
	queueComponent.Add("mem", MemoryBroker.New(10000))

	go publishMessages(queueComponent)

	go consumeMessages(queueComponent, 1)
	go consumeMessages(queueComponent, 2)

	go publishMessages(queueComponent)

	time.Sleep(20000 * time.Millisecond)
}

func publishMessages(queueComponent queue.QueueInterface) {
	i := 0

	for i < 1000 {
		_ = queueComponent.Publish("testQueue", i, nil)
		i++
	}
}

func consumeMessages(queueComponent queue.QueueInterface, thread int) {
	for {
		message, _ := queueComponent.Consume("testQueue", nil)

		if message != nil {
			fmt.Println(fmt.Sprintf("%s %d", message.GetData(), thread))
		} else {
			break
		}

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

[подробнее о воркерах](worker.md)

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
		Threads: 2, // Количество потоков в которых необходимо запустить вашу функцию
		SyncThreads: true, // Если false, воркер не будет ожидать окончания работы горутин
		Callable: printCounter,
	}
	w.Run()
}

func printCounter() error {
	i := 4

	for i > 0 {
		fmt.Printf("%+v  \n", i)
		i--
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}
```

### Queue based worker implementation

Обратите внимание на то, что функцию вы теперь передаете компоненту `QueueWorker`, и функция должна удовлетворять интерфейс принимая
на вход ссылку объекта `Message`

```go
package main

import (
	"fmt"
	queue "github.com/lika_queue"
	MemoryBroker "github.com/lika_queue/brokers/memory"
)

func main() {
	queueComponent := queue.New()
	queueComponent.Add("mem", MemoryBroker.New(10000))

	go publishMessages(queueComponent)

	var worker queue.QueueWorker

	worker = queue.QueueWorker {
		IsInfinite: true,
		QueueName: "testQueue",
		Broker: queueComponent, // or you can use any broker by calling queueComponent.Broker("mySpecialBroker")
		Worker: &queue.Worker {
			Threads: 5,
			SyncThreads: true,
		},
		Callable: myFunc,
	}

	worker.Run()
}

func myFunc(message queue.MessageInterface)  {
	fmt.Println(message.GetData())
}

func publishMessages(queueComponent queue.QueueInterface) {
	i := 0

	for i < 1000 {
		_ = queueComponent.Publish("testQueue", i, nil)
		i++
	}
}

```