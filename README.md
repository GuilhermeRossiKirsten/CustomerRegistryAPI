# Customer Registry API

API HTTP em Go para cadastro e gestão de clientes (CRUD: criação, listagem paginada, busca por ID/documento e atualização de status). Persistência em PostgreSQL, documentação Swagger e deploy via Docker Compose ou Kubernetes.

> Todos os fluxos abaixo têm um atalho no `Makefile`. Rode `make help` para ver a lista completa.

## Pré-requisitos

| Ferramenta              | Versão mínima | Necessária para                 |
| ----------------------- | ------------- | ------------------------------- |
| Docker + Docker Compose | 24+           | Subir a stack completa          |
| Go                      | 1.26          | Rodar/compilar localmente       |
| `golang-migrate` CLI    | v4            | Executar migrations             |
| `swag` CLI              | v1.16+        | Regerar docs Swagger (opcional) |
| `make`                  | qualquer      | Atalhos de comando              |

Instalação das CLIs:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

## Como subir a aplicação

### Opção 1: Docker Compose (recomendada)

```bash
make up        # sobe Postgres + API (build se necessário)
make logs      # acompanha logs da API
make down      # derruba tudo e remove volumes
```

A API fica em `http://localhost:8080` e o Swagger UI em `http://localhost:8080/swagger/index.html`.

### Opção 2: Local (Go + Postgres em container)

```bash
docker compose up -d postgres   # só o banco
make run                        # API rodando no host
```

Variáveis de ambiente (com defaults):

| Variável       | Default                 | Descrição               |
| -------------- | ----------------------- | ----------------------- |
| `DB_HOST`      | `localhost`             | Host do Postgres        |
| `DB_PORT`      | `5432`                  | Porta                   |
| `DB_USER`      | `app`                   | Usuário                 |
| `DB_PASSWORD`  | `app`                   | Senha                   |
| `DB_NAME`      | `customers-registry`    | Database                |
| `APP_PORT`     | `8080`                  | Porta da API            |
| `SWAGGER_HOST` | `http://localhost:8080` | Host exibido no Swagger |

### Opção 3: Kubernetes

```bash
make docker-build       # builda a imagem customers-registry-api:0.1.0
make k8s-apply          # aplica k8s/manifest.yaml
make k8s-port-forward   # expõe a API em localhost:8080
make k8s-delete         # remove os recursos
```

O manifest contém apenas o `Deployment` da API (2 réplicas). O Postgres precisa ser provisionado à parte (ex.: Helm chart oficial). Variáveis customizáveis: `IMAGE_NAME`, `IMAGE_TAG`, `K8S_MANIFEST`.

## Como executar migrations

As migrations ficam em `migrations/` no formato do `golang-migrate`.

```bash
make migrate-up                          # aplica todas as pendentes
make migrate-down                        # reverte a última
make migrate-create name=add_phone       # cria nova migration
```

Customize o destino com `DATABASE_URL` se necessário:

```bash
DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=disable" make migrate-up
```

> Observação: o `docker-compose.yml` ainda não roda migrations automaticamente. Após `make up`, execute `make migrate-up` para criar a tabela `customers`.

## Como rodar os testes

```bash
make test           # roda toda a suíte
make test-race      # com race detector
make coverage       # gera coverage.out e abre HTML
make fuzz           # fuzz tests do service por 30s
make lint           # golangci-lint
```

## Endpoints

Base URL: `http://localhost:8080`

| Método | Path                             | Descrição                                           |
| ------ | -------------------------------- | --------------------------------------------------- |
| POST   | `/customers`                     | Cria cliente                                        |
| GET    | `/customers`                     | Lista (paginação via `limit`/`offset`)              |
| GET    | `/customers/{id}`                | Busca por UUID                                      |
| GET    | `/customers/document/{document}` | Busca por documento                                 |
| PATCH  | `/customers/{id}/status`         | Atualiza apenas o status                            |
| GET    | `/healthz`                       | Liveness (sempre 200 se vivo)                       |
| GET    | `/readyz`                        | Readiness (200 se DB respondeu, 503 caso contrário) |
| GET    | `/swagger/index.html`            | Swagger UI                                          |

### Exemplos com `curl`

**Criar cliente**

```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{
    "document": "DOC-00001",
    "name": "João Silva",
    "score": 850,
    "risk_level": "LOW",
    "income_range": "5000-8000",
    "status": "ACTIVE"
  }'
```

Resposta `201 Created`:

```json
{
  "id": "52f2b43e-8123-4b38-8284-dc66f2ca7748",
  "document": "DOC-00001",
  "name": "João Silva",
  "score": 850,
  "risk_level": "LOW",
  "income_range": "5000-8000",
  "status": "ACTIVE",
  "created_at": "2026-05-19T12:00:00Z",
  "updated_at": "2026-05-19T12:00:00Z"
}
```

**Listar com paginação**

```bash
curl "http://localhost:8080/customers?limit=10&offset=0"
```

**Buscar por ID**

```bash
curl http://localhost:8080/customers/52f2b43e-8123-4b38-8284-dc66f2ca7748
```

**Buscar por documento**

```bash
curl http://localhost:8080/customers/document/DOC-00001
```

**Atualizar status**

```bash
curl -X PATCH http://localhost:8080/customers/52f2b43e-8123-4b38-8284-dc66f2ca7748/status \
  -H "Content-Type: application/json" \
  -d '{"status":"UNDER_REVIEW"}'
```

**Health checks**

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

### Códigos de erro

| Código | Quando                                                                                         |
| ------ | ---------------------------------------------------------------------------------------------- |
| 400    | JSON inválido, campos obrigatórios faltando, score fora de 0–1000, risk_level/status inválidos |
| 404    | Cliente não encontrado por ID/documento                                                        |
| 409    | Documento duplicado                                                                            |
| 500    | Erro inesperado                                                                                |
| 503    | `/readyz` quando o DB não responde                                                             |

### Coleção pronta: `request.http`

Na raiz do projeto existe o arquivo **`request.http`** com todas as chamadas acima já parametrizadas, incluindo health checks, paginação, busca por documento, atualização de status e os principais cenários de erro (400, 404, 409).

Como usar:

- **VS Code**: instale a extensão [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) e clique em "Send Request" acima de cada bloco.

Ajuste as variáveis `@base_url` e `@customer_id` no topo do arquivo conforme o ambiente.

## Comandos do Makefile

| Alvo                                                                         | O que faz                                      |
| ---------------------------------------------------------------------------- | ---------------------------------------------- |
| `make help`                                                                  | Lista todos os alvos disponíveis               |
| `make up` / `make down` / `make logs`                                        | Ciclo do Docker Compose                        |
| `make build` / `make run` / `make tidy`                                      | Compilar, rodar e ajustar deps localmente      |
| `make test` / `make test-race` / `make fuzz` / `make coverage` / `make lint` | Suíte de qualidade                             |
| `make migrate-up` / `make migrate-down` / `make migrate-create name=...`     | Migrations                                     |
| `make swag`                                                                  | Regerar Swagger                                |
| `make docker-build`                                                          | Build da imagem `customers-registry-api:0.1.0` |
| `make k8s-apply` / `make k8s-delete` / `make k8s-port-forward`               | Kubernetes                                     |

## Decisões técnicas

- **Go 1.26 + `net/http` puro com camadas `handler → service → repository`**. O `ServeMux` do stdlib já tem method-routing e path params desde 1.22, então não precisa de framework. A separação em camadas mantém o handler só com HTTP, o service com regras de negócio e o repository isolando SQL, permitindo testar com mocks sem subir banco.
- **`pgx` via `database/sql` + migrations versionadas com `golang-migrate`**. `pgx` é o driver Postgres mais maduro; usá-lo pela interface `database/sql` preserva o pool padrão. Migrations fora do app garantem deploy controlado e rollback determinístico, com `CHECK`s no schema (`score`, `risk_level`, `status`) reforçando a validação como defesa em profundidade.
- **Dockerfile multi-stage com distroless + K8s com `securityContext` restritivo**. Imagem final ~15 MB (`gcr.io/distroless/static-debian12:nonroot`, binário estático com `-trimpath -ldflags="-s -w"`), rodando como UID 65532. No K8s: `runAsNonRoot`, `readOnlyRootFilesystem`, `capabilities: drop ALL`, `seccompProfile: RuntimeDefault`, rolling update com `maxUnavailable: 0`.

## O que seria melhorado com mais tempo

- **Graceful shutdown + migrations automáticas**. Capturar `SIGTERM` para drenar requests em andamento antes de fechar o pool, e rodar migrations num init container/job dedicado em vez de depender de `make migrate-up` manual.
- **Observabilidade completa**: logs estruturados (`slog`), métricas Prometheus em `/metrics` (latência, status, pool do DB) e tracing OpenTelemetry, já que hoje só há `log.Printf`.
- **Segurança e resiliência em produção**: secrets fora do manifest do K8s (SealedSecrets/External Secrets), autenticação + rate limiting na API e testes de integração com Postgres real via `testcontainers-go`.
