package main

import (
	"fmt"
	"sync"
	"time"
)

const numWorkers = 3
const numPedidos = 5

type Resultado struct {
	WorkerID int
	Pedido   int
}

func worker(id int, pedidos <-chan int, resultados chan<- Resultado, wg *sync.WaitGroup) {
	defer wg.Done()

	for pedido := range pedidos {
		fmt.Printf("Worker %d: Iniciando pedido %d\n", id, pedido)
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("Worker %d: Encerrando pedido %d\n", id, pedido)
		resultados <- Resultado{WorkerID: id, Pedido: pedido}
	}

}

func main() {
	pedidos := make(chan int)
	resultados := make(chan Resultado)
	var wg sync.WaitGroup
	total := 0
	processadosPorWorker := make(map[int]int)

	for id := 1; id <= numWorkers; id++ {
		wg.Add(1)
		go worker(id, pedidos, resultados, &wg)
	}

	go func() {
		defer close(pedidos)

		for pedido := 1; pedido <= numPedidos; pedido++ {
			pedidos <- pedido
		}
	}()

	go func() {
		wg.Wait()
		close(resultados)
	}()

	for resultado := range resultados {
		total++
		processadosPorWorker[resultado.WorkerID]++
	}

	fmt.Printf("Total: %d\n", total)
	for w := 1; w <= numWorkers; w++ {
		fmt.Printf("Worker %d: processou %d pedidos\n", w, processadosPorWorker[w])
	}
}
