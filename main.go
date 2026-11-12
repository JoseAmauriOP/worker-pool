package main

import (
	"fmt"
	"sync"
	"time"
)

const numWorkers = 3
const numPedidos = 5

type Contador struct {
	mu    sync.Mutex
	total int
}

func (c *Contador) Incrementar() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.total++
}

func (c *Contador) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.total
}

func worker(id int, pedidos <-chan int, wg *sync.WaitGroup, c *Contador) {
	defer wg.Done()

	for pedido := range pedidos {
		fmt.Printf("Worker %d: Iniciando pedido %d\n", id, pedido)
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("Worker %d: Encerrando pedido %d\n", id, pedido)
		c.Incrementar()
	}

}

func main() {
	pedidos := make(chan int)
	var wg sync.WaitGroup
	var c Contador
	const numWorkers = 3
	const numPedidos = 5

	for id := 1; id <= numWorkers; id++ {
		wg.Add(1)
		go worker(id, pedidos, &wg, &c)
	}

	for pedido := 1; pedido <= numPedidos; pedido++ {
		pedidos <- pedido
	}

	close(pedidos)
	wg.Wait()
	fmt.Printf("Total: %d\n", c.Total())
}
