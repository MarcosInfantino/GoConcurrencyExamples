package main

import (
	"sync"
)

const THREAD_COUNT = 3
const VALUE_TO_ADD_PER_THREAD = 10000

type Counter struct {
	mutex sync.Mutex
	value int
}

func (c *Counter) getValue() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.value
}

func (c *Counter) increment() int {
	c.mutex.Lock()
	c.value++
	c.mutex.Unlock()
	return c.value
}

func main() {
	counter := Counter{value: 0}
	wg := sync.WaitGroup{}
	wg.Add(THREAD_COUNT)

	for i := 0; i < THREAD_COUNT; i++ {
		go func() {
			for i := 0; i < VALUE_TO_ADD_PER_THREAD*2; i++ {
				if i%2 == 0 {
					counter.increment()
				} else {
					println("El valor del contador es: ", counter.getValue())
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	println("El valor FINAL del contador es: ", counter.getValue())
	// Es necesario que este ultimo acceso al contador sea a traves del monitor?
	// A esta altura de la ejecucion, se está accediendo al recurso para lecturas y escrituras
	// de forma concurrente?
}
