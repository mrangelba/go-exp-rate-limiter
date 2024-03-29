# Desafio Rate Limiter - Pós Go Expert

## Como executar

### Executar o docker-compose
Para executar o docker compose 

```sh
make run
```
Ou
```sh
docker-compose -f deployments/docker-compose.yml up
```

### Buildar imagem docker
```sh
make build-pkg
```

### Executar stress-test
Execute cada comando em uma janela separada do terminal
```sh
make run

docker run --network=host --rm mrangelba/go-exp-stress-test --url http://localhost:8080 --concurrency 50 --requests 5000
```

### Verificar cobertura de testes
```sh
make test-coverage
```

## Arquivo de configuração `.env`
Por exemplo, para definir todas as opções disponíveis no .env use:

|Variável|Descrição|
|-|-|
|`CACHE`|`redis` ou `inmemory`|
|`REDIS_HOST`|Nome do host do servidor redis|
|`REDIS_DB`|Banco de dados do redis|
|`REDIS_PASSWORD`|Senha do banco de dados redis|
|`REDIS_PORT`|Porta do banco de dados redis|
|`RATE_LIMIT_DEFAULT_REQUESTS`|Número máximo padrão de requisições permitidas.                 |
|`RATE_LIMIT_DEFAULT_EVERY`|Intervalo de tempo padrão (em segundos) para o limite de requisições. |
|`RATE_LIMIT_IP_0`|Endereço IP específico (ex.: 192.168.65.1) para aplicar limites de requisição. |
|`RATE_LIMIT_IP_0_REQUESTS`|Número máximo de requisições permitidas para o IP especificado. |
|`RATE_LIMIT_IP_0_EVERY`|Intervalo de tempo (em segundos) para o limite de requisições para o IP especificado. |
|`RATE_LIMIT_TOKEN_0`|Token de acesso específico (ex.: token_1) para aplicar limites de requisição. |
|`RATE_LIMIT_TOKEN_0_REQUESTS`|Número máximo de requisições permitidas para o token especificado. |
|`RATE_LIMIT_TOKEN_0_EVERY`|Intervalo de tempo (em segundos) para o limite de requisições para o token especificado. |
|`RATE_LIMIT_TOKEN_1`|Outro token de acesso específico (ex.: token_2) para aplicar limites de requisição. |
|`RATE_LIMIT_TOKEN_1_REQUESTS`|Número máximo de requisições permitidas para o segundo token especificado. |
|`RATE_LIMIT_TOKEN_1_EVERY`|Intervalo de tempo (em segundos) para o limite de requisições para o segundo token especificado. |

### Exemplo

Crie um arquivo na raiz do seu projeto chamado `.env`.

```env
CACHE=redis

REDIS_HOST=redis
REDIS_DB=0
REDIS_PASSWORD=
REDIS_PORT=6379

RATE_LIMIT_DEFAULT_REQUESTS=20
RATE_LIMIT_DEFAULT_EVERY=60

RATE_LIMIT_IP_0=192.168.65.1
RATE_LIMIT_IP_0_REQUESTS=10
RATE_LIMIT_IP_0_EVERY=60

RATE_LIMIT_TOKEN_0=token_1
RATE_LIMIT_TOKEN_0_REQUESTS=5
RATE_LIMIT_TOKEN_0_EVERY=60

RATE_LIMIT_TOKEN_1=token_2
RATE_LIMIT_TOKEN_1_REQUESTS=5
RATE_LIMIT_TOKEN_1_EVERY=30
```

## Features
O middleware verifica se o limite de requisições foi atingido para o token ou IP específico. Se o limite for excedido, uma resposta HTTP 429 (Too Many Requests) será retornada.

Header da resposta

|Header|Descrição|
|-|-|
|`Ratelimit-Limit`|Limite total de requests|
|`Ratelimit-Remaining`|Limite de requests restante|
|`Ratelimit-Reset`|Tempo para reiniciar|


## Adicionando o middleware ao seu router

### Exemplo de uso com NET/HTTP

```go
package main

import (
	"log"
	"net/http"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/usecases"
	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config"
	"github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/http/middlewares"
	"github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/strategies"
)

func main() {
	log.Println("Starting server...")

	config := config.GetConfig()

	uc := usecases.NewRateLimitUseCase(
		config,
		strategies.GetCacheStrategy(config.Cache),
	)

	rateLimit := middlewares.NewRateLimiter(uc)

	http.Handle("/",
		rateLimit.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, Go Expert!"))
		})),
	)

	http.ListenAndServe(":8080", nil)
}

```

### Exemplo de uso com GO-CHI

```go
package main

import (
	"log"
	"net/http"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/usecases"
	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config"
	"github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/http/middlewares"
	"github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/strategies"

	"github.com/go-chi/chi/v5"
)

func main() {
	log.Println("Starting server...")

	config := config.GetConfig()

	uc := usecases.NewRateLimitUseCase(
		config,
		strategies.GetCacheStrategy(config.Cache),
	)

	rateLimit := middlewares.NewRateLimiter(uc)

	router := chi.NewRouter()
	router.Use(rateLimit.Handler)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Go Expert!"))
	})

	http.ListenAndServe(":8080", router)
}
```

### Saída do stress-test

```sh
----------------------------------------------------------------------
Relatório de execução
----------------------------------------------------------------------
Tempo total gasto na execução: 4.453823378s
Quantidade total de requests realizados: 5000
----------------------------------------------------------------------
Quantidade de requests com status HTTP 200: 20
----------------------------------------------------------------------
Quantidade de requests com status HTTP 429: 4980
----------------------------------------------------------------------
```

### Saída da cobertura de testes
```sh
coverage: 100.0% of statements in ./...
ok      github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/utils  1.069s  coverage: 100.0% of statements in ./...
go tool cover -html coverage/coverage.out -o coverage/coverage.html
go tool cover -func coverage/coverage.out
github.com/mrangelba/go-exp-rate-limiter/internal/domain/entities/rate_limiter.go:13:                           MarshalBinary           100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/domain/entities/rate_limiter.go:17:                           UnmarshalBinary         100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/domain/usecases/rate_limit_usecase.go:26:                     NewRateLimitUseCase     100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/domain/usecases/rate_limit_usecase.go:33:                     VerifyLimit             100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/domain/usecases/rate_limit_usecase.go:82:                     validateCacheLimit      100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/domain/usecases/rate_limit_usecase.go:102:                    GetHttpHeaders          100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/drivers/cache/redis/client.go:15:                             GetClient               100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/drivers/cache/redis/client.go:23:                             connectRedis            100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config/config.go:17:                                  GetConfig               100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config/config.go:33:                                  String                  100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config/rate_limiter/rate_limiter_config.go:33:        GetRateLimiterConfig    100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config/redis/redis_config.go:14:                      GetRedisConfig          100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/http/middlewares/rate_limit.go:17:             NewRateLimiter          100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/http/middlewares/rate_limit.go:24:             Handler                 100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/http/middlewares/rate_limit.go:46:             checkLimitAddHeaders    100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/http/middlewares/rate_limit.go:65:             getIPs                  100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/strategies/rate_limit_in_memory.go:16:         NewRateLimitInMemory    100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/strategies/rate_limit_in_memory.go:22:         Set                     100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/strategies/rate_limit_in_memory.go:27:         Get                     100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/strategies/rate_limit_redis.go:20:             NewRateLimitRedis       100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/strategies/rate_limit_redis.go:26:             Set                     100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/strategies/rate_limit_redis.go:37:             Get                     100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/strategies/rate_limit_strategy.go:9:           GetCacheStrategy        100.0%
github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/utils/gin.go:9:                                MiddlewareToGin         100.0%
total:                                                                                                          (statements)            100.0%
```

## Requisitos do Desafio

:white_check_mark: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.

:white_check_mark: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
API_KEY: <TOKEN>

:white_check_mark: O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web

:white_check_mark: O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.

:white_check_mark: O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.

:white_check_mark: As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.

:white_check_mark: Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.

:white_check_mark: O sistema deve responder adequadamente quando o limite é excedido:
- Código HTTP: 429
- Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame

:white_check_mark: Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. 

:white_check_mark: Você pode utilizar docker-compose para subir o Redis.
Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.

:white_check_mark: A lógica do limiter deve estar separada do middleware.