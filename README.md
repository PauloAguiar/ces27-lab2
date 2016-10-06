# CES-27 - LAB 2 - Dynamo

Instruções para configuração do ambiente: [Setup](SETUP.md)  
Para saber mais sobre a utilização do Git, acessar: [Git Passo-a-Passo](GIT.md)  
Para saber mais sobre a entrega do trabalho, acessar: [Entrega](ENTREGA.md)  

## Referências
O paper no qual essa atividade é baseada pode ser encontrado em: [Dynamo: Amazon’s Highly Available Key-value Store](http://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf)

## Executando o Código

Durante a execução dessa atividade, será necessário que sejam inicializadas quatro instâncias do serviço proposto. O arquivo main da execução(server.go) se encontra na pasta raiz do projeto.

A primeira instância pode ser inicializada como mostrado a seguir:
```bash
$ go run server.go
```
> [SERVER] Running Dynamo Server with id 'localhost:3000'  
> [RING] Added node 'localhost:3000'(hostname: 'localhost:3000') to the local ring(hash: '795')  
> [SERVER] Listening on 'localhost:3000'  

As demais instâncias devem utilizar o seguinte comando (incrementando o parâmetro port para cade instância a mais):
```bash
$ go run server.go -port 3001 -ring localhost:3000
```
> [SERVER] Running Dynamo Server with id 'localhost:3001'  
> [RING] Added node 'localhost:3000'(hostname: 'localhost:3000') to the local ring(hash: '795')  
> [RING] Added node 'localhost:3001'(hostname: 'localhost:3001') to the local ring(hash: '165')  
> [RING] Reporting to host 'localhost:3000'  
> [SERVER] Listening on 'localhost:3001'  

Para evitar erros de sincronização, aguarde um pouco entre a inicialização de cada uma das instâncias.

O Resultado da execução é mostrado a seguir:
![Run All](doc/run-all.PNG?raw=true)

## Console
Para simplificar a interação com o servidor, é possível entrar com comandos durante a execuçao do serviço como mostrado a seguir:

```bash
ring
```
> [CONSOLE] HASH__|ID_______________|  
> [CONSOLE] '165'_|'localhost:3001'_|  
> [CONSOLE] '473'_|'localhost:3003'_|  
> [CONSOLE] '543'_|'localhost:3002'_|  
> [CONSOLE] '795'_|'localhost:3000'_|  

Lista de comandos disponíveis:
put *key* *value* - put KV(key ,value) pair on this instance only  
rput *key* *value* *quorum* - route a put request (the same as if an external client made a put request) using a write quorum  
get *key* - get a the value associated with the given key  
rget *key* *quorum* - route a get request (the same as if an external client made a get request) using a read quorum  
print - print all the data stored in this instance  
ring - print the topology of the ring structure  
down - set the server as down (unavailable)  
up - set the server as up (available)  


# Hash Consistente

Na primeira parte da atividade, vamos corrigir o funcionamento do hash consistente.

Abra o arquivo common/consistenthash/consistenthash.go:

A função a seguir não possui o código necessário para a correta distribuição dos dados entre os nós.

```go
// search will find the index of the node that is responsible for the range that
// includes the hashed value of key.
func (r *Ring) search(key string) int {
    /////////////////////////
    // YOUR CODE GOES HERE //
    /////////////////////////

    return 0
}
```

Ao sempre retornar zero, a função redireciona todas as operações de leitura e escrita para o mesmo nó.

Para verificar o problema, basta executar o seguinte comando:
```bash
rput alice hi 2
```
> [ROUTER] Routing Put of KV('alice', 'hi') with quorum Q('2').  
> [ROUTER] Trying 'localhost:3001' as coordinator.  

A primeira tentativa de roteamento da key "alice" deveria ser o nó com id "localhost:3000", mas por conta do erro na função search, sempre é escolhido o nó com id "localhost:3001" pois este é o primeiro da lista(índice 0).

```bash
ring
```
> [CONSOLE] HASH__|ID_______________|  
> [CONSOLE] '165'_|'localhost:3001'_|  
> [CONSOLE] '473'_|'localhost:3003'_|  
> [CONSOLE] '543'_|'localhost:3002'_|  
> [CONSOLE] '795'_|'localhost:3000'_|  

Ao implementar a função search corretamente, o comando acima deve resultar em:

```bash
rput alice hi 2
```
> [ROUTER] Routing Put of KV('alice', 'hi') with quorum Q('2').  
> [ROUTER] Trying 'localhost:3000' as coordinator.  

Pois o código abaixo:
```go
func hashId(key string) uint32 {
    return crc32.ChecksumIEEE([]byte(key)) % uint32(1000)
}
```
retorna o hash 735

# Versionamento usando Timestamps

Na segunda parte desta atividade, vamos implementar a utilização de comparação de timestamps para decidir entre diferentes versões de uma mesma chave, o que pode acontecer quando um ou mais servidores ficam indisponíveis.

Para verificar o problema, execute os seguintes comandos:

Adiciona uma chave cujo hash caia dentro do range do servidor 3000:

```bash
rput alice hi 2
```
> [ROUTER] Routing Put of KV('alice', 'hi') with quorum Q('2').  
> [ROUTER] Trying 'localhost:3000' as coordinator.  

O sistema de replicação vai então colocar cópias dos dados no servidores 3001 e 3003.

> 2016/10/06 09:00:29 [INTERNAL] Replicating 'alice' = 'hi' (timestamp: '1475755228')  

Agora deixe o servidor 3000 indisponível:

```bash
down
```
> [CONSOLE] Putting server DOWN.  
> [SERVER] Server Stopped  

Escreva um novo valor para a chave utilizada acima:

```bash
rput alice oi 2
```
> [ROUTER] Routing Put of KV('alice', 'oi') with quorum Q('2').  
> [ROUTER] Trying 'localhost:3000' as coordinator.  
> [ROUTER] Coordinator tryout failed. Error: dial tcp [::1]:3000: connectex: No connection (...)  
> [ROUTER] Trying 'localhost:3001' as coordinator.  
> [ROUTER] Coordinate succeded.  

Como o servidor 3000 está indísponível, o próximo coordenador é o seguinte a ele no anel, nesse caso o servidor 3001.

O valor foi atualizado em 3001 e 3003, mas não em 3000 pois este está indisponível:

> [COORDINATOR] Error on replication: dial tcp [::1]:3000: connectex: No connection could be made because the target machine actively refused it.  

Reative o servidor 3000:

```bash
up
```
> [CONSOLE] Putting server UP.  
> [SERVER] Listening on 'localhost:3000'  

Observe que o servidor 3000 possui a versão antiga da chave:

```bash
print
```
> [CONSOLE] KEY_____|VALUE_|TIMESTAMP_|  
> [CONSOLE] 'alice'_|'hi'__|0_________|  

Realize uma requisição Get pela chave:

```bash
rget alice 2
```
> [COORDINATOR] Coordinating voting of K('alice') with Quorum '2'  
> [CACHE] Getting Key 'alice' with Value 'hi' @ timestamp '0'  
> [COORDINATOR] Voting with quorum '2' succeded.  
> [COORDINATOR] Vote: hi  
> [COORDINATOR] Vote: oi  

E o resultado:

> [ROUTER] Coordinate succeded.  
> [CONSOLE] Rget result: 'alice' = 'hi'  

Esse erro (receber "hi" ao invés do último valor "oi") pois o serviço não utiliza a compara
ção de timestamps. Ao invés disso, ele sempre seleciona o primeiro elemento da lista de votos.

Para corrigir o problema, é necessário alterar duas coisas:

O arquivo dynamo/cache.go não armazena as timestamps recebidas durante a gravação:

```go
// Put a value to a key in the storage. This will handle concurrent put
// requests by locking the structure.
func (cache *Cache) Put(key string, value string, timestamp int64) {
    log.Printf("[CACHE] Putting Key '%v' with Value '%v' @ timestamp '%v'\n", key, value, timestamp)

    cache.Lock()
    cache.data[key] = value
    cache.Unlock()

    return
}
```

Você precisa alterar esse arquivo (não somente esse método) para armazenar e retornar corretamente a timestamp associada a um dado.

Quando a implementação estiver correta, o valor de timestamp retornado pelo cache deve ser o mesmo que o foi passado durante a gravação.

> [CACHE] Putting Key 'alice' with Value 'hi' @ timestamp '1475755228'  
> [CACHE] Getting Key 'alice' with Value 'hi' @ timestamp '1475755228'  

Por último, no arquivo dynamo/coordinator.go você deve implementar a função abaixo:

```go
// aggregateVotes will select the right value from the votes received.
func aggregateVotes(votes []*vote) (result string) {
    for _, vote := range votes {
        log.Printf("[COORDINATOR] Vote: %v\n", vote.value)
    }

    /////////////////////////
    // YOUR CODE GOES HERE //
    /////////////////////////

    result = votes[0].value
    return
}
````