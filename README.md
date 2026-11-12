# Estudo de Worker Pool em Go

Projeto individual de estudo, feito para aprender na prática como funciona um worker pool e os conceitos de concorrência de **Go** que ele envolve: goroutines, channels, `sync.WaitGroup` e `sync.Mutex`. Os pedidos são simulados: o objetivo não é a regra de negócio, e sim observar como as goroutines dividem o trabalho e se coordenam.

## Como funciona

```text
                                 ┌─> worker 1 ─┐
goroutine        ┌──────────┐    │             │    ┌────────────┐
produtora ─────> │ pedidos  │ ───┼─> worker 2 ─┼──> │ resultados │ ───> main
(envia 1..5)     └──────────┘    │             │    └────────────┘     (conta total
                                 └─> worker 3 ─┘                        e por worker)
```

- **3 workers, 5 pedidos**: definidos pelas constantes `numWorkers` e `numPedidos`.
- **Distribuição automática**: cada pedido do channel `pedidos` é entregue a um único worker livre. A divisão muda a cada execução.
- **Processamento simulado**: cada pedido leva 200ms (`time.Sleep(200 * time.Millisecond)`), uma espera que existe só para deixar a execução simultânea visível na saída.
- **Encerramento coordenado**: a goroutine produtora fecha `pedidos` depois do último envio. Os workers saem do `for range`, e uma goroutine que espera o `sync.WaitGroup` fecha `resultados`.
- **Contagem sem trava**: só a `main` lê `resultados` e altera o total e o map `processadosPorWorker`.
- **Tempo total de cerca de 400ms**: os 5 pedidos são processados em duas rodadas paralelas, contra cerca de 1s em sequência. Essa diferença mostra na prática o ganho da concorrência.

## Decisões de projeto

- **Mesma contagem com mutex e com channel, para comparar**: a primeira versão (commit `e9ab795`) foi escrita para praticar `sync.Mutex`, com um `Contador` compartilhado entre os workers. A versão atual refaz a contagem com channel: os workers enviam um `Resultado` e só a `main` mexe nos dados. Como existe um único dono, não há o que proteger. Ter as duas versões permitiu comparar quando usar cada abordagem.
- **Channels unbuffered**: a `main` lê os resultados enquanto os workers trabalham, então não é preciso buffer para evitar deadlock. Usar um buffer só esconderia um erro de ordem das operações.
- **Envio dos pedidos em uma goroutine separada**: essa decisão veio de um deadlock encontrado durante o estudo. Quando a `main` enviava os pedidos e só depois lia os resultados, os workers travavam no envio e a `main` travava no envio do próximo pedido.
- **Direção dos channels na assinatura**: `worker` recebe `pedidos <-chan int` (só recebe) e `resultados chan<- Resultado` (só envia). Um uso no sentido errado não compila.

## Tecnologias

- Go 1.26.5
- Biblioteca padrão: `sync`, `time` e `fmt` (sem dependências externas)

## Como rodar

```bash
git clone https://github.com/JoseAmauriOP/worker-pool.git
cd worker-pool
go run .
# Windows (PowerShell ou CMD): os mesmos comandos
```

A saída mostra o início e o fim de cada pedido e, no final, a contagem. Rode mais de uma vez: a ordem e a distribuição entre os workers mudam a cada execução, e é isso que o projeto quer mostrar.

```text
Worker 3: Iniciando pedido 1
Worker 1: Iniciando pedido 2
Worker 2: Iniciando pedido 3
...
Total: 5
Worker 1: processou 1 pedidos
Worker 2: processou 2 pedidos
Worker 3: processou 2 pedidos
```

Para checar condições de corrida com o detector do Go, usado durante todo o estudo:

```bash
go run -race .
```

## Como o projeto está organizado

```text
worker-pool/
├── go.mod      # módulo worker-pool, Go 1.26.5
├── main.go     # tipo Resultado, função worker e a main que monta o pool
└── README.md
```

## O que aprendi

- `go worker(...)` inicia uma goroutine. Quando a `main` termina, o programa acaba sem esperar as goroutines.
- Um envio em channel unbuffered bloqueia até alguém receber. Enviar antes de existir um receptor causa `fatal error: all goroutines are asleep - deadlock!`.
- `close(pedidos)` encerra o `for pedido := range pedidos` dos workers. Sem ele, o `wg.Wait()` nunca retorna.
- `sync.WaitGroup` (`wg.Add(1)` antes do `go`, `defer wg.Done()` no worker, `wg.Wait()`) coordena o fim das goroutines. Ele é passado como ponteiro.
- `processados++` em várias goroutines perde incrementos. `go run -race` detecta o problema, e `sync.Mutex` com `Lock`/`Unlock` resolve.
- Use mutex para proteger um estado compartilhado e channel para levar dados de uma goroutine para outra.

Desenvolvido por [José Amauri Patrício](www.linkedin.com/in/jose-amauri-patricio)