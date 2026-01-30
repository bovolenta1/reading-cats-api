# Copilot Instructions — Reading Cats API

Este documento guia o desenvolvimento e manutenção da **Reading Cats API**, uma aplicação serverless em Go com arquitetura em camadas (Clean Architecture + DDD).

---

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Arquitetura](#arquitetura)
- [Estrutura de Pastas](#estrutura-de-pastas)
- [Padrões & Convenções](#padrões--convenções)
- [Fluxo de Requisição](#fluxo-de-requisição)
- [Autenticação & Autorização](#autenticação--autorização)
- [Adicionando Novos Endpoints](#adicionando-novos-endpoints)
- [Guia de Implementação por Camada](#guia-de-implementação-por-camada)
- [Tratamento de Erros](#tratamento-de-erros)
- [Testes & Desenvolvimento Local](#testes--desenvolvimento-local)
- [Deployment](#deployment)

---

## 🎯 Visão Geral

**Propósito:** Backend para a aplicação "Reading Cats", que rastreia progresso de leitura do usuário.

**Stack:**
- **Runtime:** Go 1.21+
- **Deployment:** AWS Lambda + API Gateway (HTTP API)
- **Infrastructure as Code:** AWS SAM
- **Database:** PostgreSQL (Neon)
- **Migrations:** golang-migrate
- **Autenticação:** AWS Cognito (JWT)

**Endpoints Atuais:**
```
GET  /v1/me                   → Usuário autenticado (me)
POST /v1/reading/logs         → Registrar leitura diária
GET  /v1/reading/progress     → Progresso de leitura
PUT  /v1/reading/goal         → Alterar meta de leitura
POST /v1/groups               → Criar novo grupo
```

---

## 🏗️ Arquitetura

A aplicação segue **Clean Architecture** com **Domain-Driven Design**:

```
┌─────────────────────────────────────────────────┐
│         HTTP API (API Gateway)                   │
└────────────────┬────────────────────────────────┘
                 │
         ┌───────▼────────┐
         │ Presentation   │  (handlers, routing, I/O)
         └───────┬────────┘
                 │
         ┌───────▼────────┐
         │ Application    │  (use cases, DTOs)
         └───────┬────────┘
                 │
         ┌───────▼────────┐
         │ Domain         │  (entidades, lógica, VOs)
         └───────┬────────┘
                 │
         ┌───────▼────────┐
         │ Infrastructure │  (repositories, DB)
         └────────────────┘
```

### Princípios:
- **Domain Layer:** Totalmente independente de frameworks; contém a lógica pura de negócio
- **Application Layer:** Orquestra os use cases; não conhece detalhes de HTTP ou DB
- **Infrastructure Layer:** Implementação concreta (Postgres, etc.)
- **Presentation Layer:** Controllers, handlers, validação de entrada

---

## 📁 Estrutura de Pastas

```
reading-cats-api/
├── main.go                          # Entrada Lambda, inicialização de dependências
├── go.mod & go.sum                  # Dependências Go
├── template.yaml                    # AWS SAM template
├── Makefile                         # Build, run, migrations
├── env.json & env.local             # Variáveis de ambiente
│
├── internal/
│   ├── domain/                      # DOMAIN LAYER (lógica pura)
│   │   ├── user/
│   │   │   ├── user.go              # Entidade User
│   │   │   ├── value_objects.go     # CognitoSub, Email, DisplayName, etc.
│   │   │   ├── errors.go            # Erros de domínio do usuário
│   │   │
│   │   ├── reading/
│   │   │   ├── progress.go          # Lógica de progresso de leitura
│   │   │   ├── value_objects.go     # LocalDate, Pages, etc.
│   │   │   └── errors.go            # Erros de domínio de leitura
│   │   │
│   │   └── group/
│   │       ├── group.go             # Entidade Group
│   │       ├── value_objects.go     # GroupName, IconID, Visibility
│   │       └── errors.go            # Erros de domínio de grupo
│   │
│   ├── application/                 # APPLICATION LAYER (use cases)
│   │   ├── user/
│   │   │   ├── ensure_me.go         # UseCase: encontrar ou criar usuário
│   │   │   ├── dto.go               # DTOs (Input/Output)
│   │   │   └── repository.go        # Interface do repositório
│   │   │
│   │   ├── reading/
│   │   │   ├── register_reading.go  # UseCase: registrar leitura
│   │   │   ├── get_reading_progress.go
│   │   │   ├── change_goal.go       # UseCase: alterar meta
│   │   │   ├── dto.go               # DTOs (Input/Output)
│   │   │   └── repository.go
│   │   │
│   │   └── group/
│   │       ├── create_group.go      # UseCase: criar grupo
│   │       ├── dto.go               # DTOs (Input/Output)
│   │       └── repository.go        # Interface do repositório
│   │├── reading/
│   │   │   └── postgres_repository.go # Implementação do repo de reading
│   │   └── group/
│   │       └── postgres_repository.go # Implementação do repo de grupo
│   │   ├── db/
│   │   │   └── postgres.go          # Pool de conexões Postgres
│   │   ├── user/
│   │   │   └── postgres_repository.go # Implementação do repo de usuário
│   │   └── reading/
│   │       └── postgres_repository.go # Implementação do repo de reading
│   │
│   ├── presentation/                # PRESENTATION LAYER (HTTP)
│   │   └── httpapi/
│   │       ├── router.go            # Roteamento principal
│   │       ├── auth_claims.go       # Extração de JWT
│   │       ├── me_input.go
│   │       ├── register_reading_handler.go
│   │       ├── register_reading_input.go
│   │       ├── get_reading_progress_handler.go
│   │       ├── get_reading_progress_input.go
│   │       ├── change_goal_handler.go
│   │       ├── change_goal_input.go
│   │       ├── create_group_handler.go      # Handler: POST /v1/groups
│   │       └── create_group_inputading_input.go
│   │       ├── get_reading_progress_input.go
│   │       └── auth_claims.go
│   │
│   └── config/
│   ├── 000002_create_reading_tables.down.sql
│   ├── 000003_add_valid_from_to_reading_goal.up.sql
│   ├── 000003_add_valid_from_to_reading_goal.down.sql
│   ├── 000004_switch_to_uuid_keys.up.sql
│   ├── 000004_switch_to_uuid_keys.down.sql
│   ├── 000005_rename_reading_day_to_user_checkins.up.sql
│   ├── 000005_rename_reading_day_to_user_checkins.down.sql
│   ├── 000006_create_groups_schema.up.sql
│   ├── 000006_create_groups_schema.down.sql
│   ├── 000007_alter_group_seasons_table.up.sql
│   └── 000007_alter_group_seasons_table # Variáveis de ambiente
│
├── migrations/                      # SQL migrations (golang-migrate)
│   ├── 000001_create_user.up.sql
│   ├── 000001_create_user.down.sql
│   ├── 000002_create_reading_tables.up.sql
│   └── 000002_create_reading_tables.down.sql
│
└── README.md                        # Setup & instruções locais
```

---

## 🎨 Padrões & Convenções

### 1. **Nomenclatura de Pacotes**
- Pacotes = `lowercase` sem underscores
- Tipos = `PascalCase`
- Interfaces = `PascalCase` (sufixo `er` quando verbo, ex: `Reader`, `Repository`)
- Variáveis privadas = `camelCase`
- Constantes = `PascalCase` (ou SCREAMING_SNAKE_CASE para grupos)

### 2. **Value Objects (VO)**
Valores que representam conceitos do domínio com validação:

```go
// Exemplo: CognitoSub é um VO
type CognitoSub string

func NewCognitoSub(v string) (CognitoSub, error) {
    v = strings.TrimSpace(v)
    if v == "" {
        return "", ErrInvalidCognitoSub
    }
    return CognitoSub(v), nil
}
```

**Regra:** Sempre usar construtores `New*` com validação; nunca criar direto.

### 3. **Entidades**
Objetos com identidade única e comportamento:

```go
type User struct {
    ID            string        // UUID
    CognitoSub    CognitoSub    // VO
    Email         Email         // VO
    DisplayName   DisplayName   // VO
    AvatarURL     AvatarURL     // VO
    ProfileSource ProfileSource // enum
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### 4. **Interfaces de Repositório**
Definidas na camada de aplicação (Application Layer), implementadas na camada de infra:

```go
// application/user/repository.go
type Repository interface {
    FindByCognitoSub(ctx context.Context, sub CognitoSub) (*User, error)
    Insert(ctx context.Context, u *User) error
}

// infra/user/postgres_repository.go
type PostgresRepository struct {
    pool *pgxpool.Pool
}

func (r *PostgresRepository) FindByCognitoSub(...) (*domain.User, error) {
    // Implementação
}
```

### 5. **Use Cases**
Orquestram a lógica de aplicação:

```go
type EnsureMeUseCase struct {
    repo Repository  // Depende de abstração (interface)
}

func (uc *EnsureMeUseCase) Execute(ctx context.Context, in Input) (Output, error) {
    // Lógica de negócio
}
```

**Padrão Constructor + Dependency Injection:**
```go
func NewEnsureMeUseCase(repo Repository) *EnsureMeUseCase {
    return &EnsureMeUseCase{repo: repo}
}
```

### 6. **DTOs (Data Transfer Objects)**
Usados apenas na camada de Presentation e Application (nunca Domain):

```go
type MeDTO struct {
    ID          string `json:"id"`
    CognitoSub  string `json:"cognito_sub"`
    Email       string `json:"email"`
    DisplayName string `json:"display_name"`
    AvatarURL   string `json:"avatar_url"`
    Source      string `json:"source"`
}
```

### 7. **Handlers HTTP**
Conversam com o mundo externo (HTTP):

```go
type MeHandler struct {
    uc *app.EnsureMeUseCase
}

func (h *MeHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
    // 1. Parse input
    in, err := BuildEnsureMeInput(event)
    
    // 2. Chamar use case
    result, err := h.uc.Execute(ctx, in)
    
    // 3. Retornar resposta HTTP
    return JSON(http.StatusOK, result), nil
}
```

---

## 🔄 Fluxo de Requisição

### Exemplo: GET /v1/me

```
1. AWS API Gateway
   ↓
2. Lambda Handler (main.go)
   → router.Route(ctx, event)
   ↓
3. Router.Route()
   → Identifica GET /v1/me
   → Chama r.me.Handle(ctx, event)
   ↓
4. MeHandler.Handle()
   → BuildEnsureMeInput(event)  [extrai JWT]
   → h.uc.Execute(ctx, input)
   ↓
5. EnsureMeUseCase.Execute()
   → uc.repo.FindByCognitoSub(ctx, sub)
   → Se não existe, cria novo (domain.NewFromIDP)
   → uc.repo.Insert(ctx, &user)
   → Retorna MeDTO
   ↓
6. PostgresRepository
   → Query ao Postgres
   ↓
7. MeHandler
   → JSON(http.StatusOK, meDTO)
   ↓
8. AWS API Gateway
   → Retorna response HTTP
```

---

## 🔐 Autenticação & Autorização

### Extração do Token (auth_claims.go)

**Dois modos:**

#### 1. **PROD (AWS API Gateway Authorizer)**
```go
if event.RequestContext.Authorizer != nil && event.RequestContext.Authorizer.JWT != nil {
    // Claims já validadas pelo API Gateway
    claims := event.RequestContext.Authorizer.JWT.Claims
    // Usar direto: claims["sub"], claims["email"], etc.
}
```

#### 2. **DEV Local (SAM)**
```go
if os.Getenv("AWS_SAM_LOCAL") == "true" {
    token := bearerToken(event.Headers)  // Extrai "Bearer TOKEN"
    payload, ok := decodeJwtPayload(token)  // Decodifica base64
    // Validar payload
}
```

### Funções-chave:

- **`bearerToken(headers)`** → Extrai token do header `Authorization: Bearer <token>`
- **`decodeJwtPayload(token)`** → Decodifica JWT (sem validação de assinatura, apenas local)
- **`BuildEnsureMeInput(event)`** → Monta o input do use case com claims validadas
- **`buildIDPClaims()`** → Cria Value Objects de claims com validação

### Validação:
- `CognitoSub`: **obrigatório**, senão → 401 Unauthorized
- `Email`, `Name`, `Picture`: **opcionais**; se inválidos, ignorados

---

## ➕ Adicionando Novos Endpoints

### Passo-a-passo:

#### 1. **Domain Layer** (se necessário)
Criar entidade ou value object:

```go
// internal/domain/book/book.go
type Book struct {
    ID        string
    Title     string
    Author    string
    CreatedAt time.Time
}

// internal/domain/book/value_objects.go
type Title string
func NewTitle(v string) (Title, error) {
    v = strings.TrimSpace(v)
    if v == "" || len([]rune(v)) > 200 {
        return "", ErrInvalidTitle
    }
    return Title(v), nil
}
```

#### 2. **Application Layer**
Criar use case + DTO:

```go
// internal/application/book/repository.go
type Repository interface {
    Insert(ctx context.Context, b *domain.Book) error
    FindByID(ctx context.Context, id string) (*domain.Book, error)
}

// internal/application/book/register_book.go
type RegisterBookUseCase struct {
    repo Repository
}

type RegisterBookInput struct {
    Title  string
    Author string
}

type RegisterBookOutput struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

func (uc *RegisterBookUseCase) Execute(ctx context.Context, in RegisterBookInput) (RegisterBookOutput, error) {
    // Validar entrada
    title, err := domain.NewTitle(in.Title)
    if err != nil {
        return RegisterBookOutput{}, err
    }
    
    // Criar entidade
    book := domain.Book{
        ID:    uuid.NewString(),
        Title: title,
        Author: in.Author,
        CreatedAt: time.Now().UTC(),
    }
    
    // Persistir
    if err := uc.repo.Insert(ctx, &book); err != nil {
        return RegisterBookOutput{}, err
    }
    
    // Retornar DTO
    return RegisterBookOutput{
        ID: book.ID,
        Title: string(book.Title),
        Author: book.Author,
    }, nil
}
```

#### 3. **Infrastructure Layer**
Implementar repositório:

```go
// internal/infra/book/postgres_repository.go
type PostgresRepository struct {
    pool *pgxpool.Pool
}

func (r *PostgresRepository) Insert(ctx context.Context, b *domain.Book) error {
    _, err := r.pool.Exec(ctx,
        `INSERT INTO books (id, title, author, created_at) VALUES ($1, $2, $3, $4)`,
        b.ID, b.Title, b.Author, b.CreatedAt,
    )
    return err
}
```

#### 4. **Presentation Layer**
Criar handler:

```go
// internal/presentation/httpapi/register_book_handler.go
type RegisterBookHandler struct {
    uc *app.RegisterBookUseCase
}

func NewRegisterBookHandler(uc *app.RegisterBookUseCase) *RegisterBookHandler {
    return &RegisterBookHandler{uc: uc}
}

func (h *RegisterBookHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
    var in app.RegisterBookInput
    if err := json.Unmarshal([]byte(event.Body), &in); err != nil {
        return Error(event, http.StatusBadRequest, "invalid request body"), nil
    }
    
    out, err := h.uc.Execute(ctx, in)
    if err != nil {
        return Error(event, http.StatusInternalServerError, err.Error()), nil
    }
    
    return JSON(http.StatusCreated, out), nil
}
```

#### 5. **Atualizar Router**
```go
// internal/presentation/httpapi/router.go
type Router struct {
    me              *MeHandler
    registerReading *RegisterReadingHandler
    getReadingProgress *GetReadingProgressHandler
    registerBook    *RegisterBookHandler  // ← NOVO
}

func (r *Router) Route(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
    // ... endpoints existentes ...
    
    if event.RequestContext.HTTP.Method == http.MethodPost && event.RawPath == "/v1/books" {
        return r.registerBook.Handle(ctx, event)
    }
    
    return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNotFound}, nil
}
```

#### 6. **main.go**
Adicionar instância e injeção:

```go
// main.go
func init() {
    // ... código existente ...
    
    // book
    bookRepo := infraBook.NewPostgresRepository(pool)
    registerBookUC := appBook.NewRegisterBookUseCase(bookRepo)
    registerBookHandler := httpapi.NewRegisterBookHandler(registerBookUC)
    
    router = httpapi.NewRouter(meHandler, registerReadingHandler, getReadingProgressHandler, registerBookHandler)
}
```

#### 7. **Migrations** (se DB schema mudar)
```bash
make migrate-create NAME=add_books_table
# Editar migrations/000003_add_books_table.up.sql
# Editar migrations/000003_add_books_table.down.sql
```

---

## 📚 Guia de Implementação por Camada

### Domain Layer (internal/domain/)

**Responsabilidades:**
- Entidades com identidade e ciclo de vida
- Value Objects com validação
- Lógica pura de negócio
- Interfaces (quando necessário)
- Erros customizados

**Regras:**
- ✅ Use `type MyType string` ou `struct` sem métodos complexos
- ✅ Valide tudo via construtores `New*()`
- ✅ Sem imports de `infra`, `application`, `presentation`
- ❌ Nunca use `package main`, `log.Fatal`, `fmt.Println`
- ❌ Nunca acesse DB ou HTTP direto

**Exemplo:**
```go
// domain/reading/progress.go
type ReadingProgress struct {
    UserID      string
    TotalPages  int
    StreakDays  int
}

func NewReadingProgress(userID string, pages int) (ReadingProgress, error) {
    if userID == "" {
        return ReadingProgress{}, ErrInvalidUserID
    }
    if pages < 0 {
        return ReadingProgress{}, ErrNegativePages
    }
    return ReadingProgress{UserID: userID, TotalPages: pages}, nil
}
```

### Application Layer (internal/application/)

**Responsabilidades:**
- Orquestração de use cases
- DTOs (Input/Output)
- Interfaces de repositório
- Lógica de negócio que coordena entidades

**Regras:**
- ✅ Use cases = `*UseCase` com método `Execute(ctx, input)`
- ✅ Dependa de abstrações (interfaces)
- ✅ Valide entrada e trate erros
- ✅ Retorne DTOs (nunca entidades domain direto)
- ❌ Não acesse HTTP ou logs diretos

**Exemplo:**
```go
// application/reading/register_reading.go
type RegisterReadingUseCase struct {
    repo Repository
    tz   string  // timezone para cálculo de dias
}

func NewRegisterReadingUseCase(repo Repository, timezone string) *RegisterReadingUseCase {
    return &RegisterReadingUseCase{repo: repo, tz: timezone}
}

type RegisterReadingInput struct {
    UserID string `json:"user_id"`
    Pages  int    `json:"pages"`
}

type RegisterReadingOutput struct {
    Date       string `json:"date"`
    Pages      int    `json:"pages"`
    StreakDays int    `json:"streak_days"`
}

func (uc *RegisterReadingUseCase) Execute(ctx context.Context, in RegisterReadingInput) (RegisterReadingOutput, error) {
    // 1. Validar
    if in.Pages < 0 {
        return RegisterReadingOutput{}, errors.New("pages must be >= 0")
    }
    
    // 2. Usar repositório
    day, err := uc.repo.AddPages(ctx, nil, in.UserID, today, in.Pages)
    if err != nil {
        return RegisterReadingOutput{}, err
    }
    
    // 3. Retornar DTO
    return RegisterReadingOutput{
        Date: day.Date.String(),
        Pages: day.Pages,
        StreakDays: day.StreakDays,
    }, nil
}
```

### Infrastructure Layer (internal/infra/)

**Responsabilidades:**
- Implementação de repositórios (interfaces da application layer)
- Conexão com DB, cache, APIs externas
- Conversão entre domain models e DB models

**Regras:**
- ✅ Implemente interfaces definidas em `application/`
- ✅ Use `*pgxpool.Pool` para queries
- ✅ Trate erros de DB e converta para domain errors se necessário
- ✅ Mantenha a lógica DB-específica isolada aqui
- ❌ Nunca exponha detalhes SQL na application layer

**Exemplo:**
```go
// infra/reading/postgres_repository.go
type PostgresRepository struct {
    pool *pgxpool.Pool
}

func (r *PostgresRepository) AddPages(ctx context.Context, tx pgx.Tx, subID string, date domain.LocalDate, delta int) (application.DayRow, error) {
    var pages, streak int
    
    err := tx.QueryRow(ctx,
        `UPDATE reading_days SET pages = pages + $1, updated_at = now()
         WHERE user_id = $2 AND date = $3
         RETURNING pages, streak_days`,
        delta, subID, date,
    ).Scan(&pages, &streak)
    
    if err != nil {
        return application.DayRow{}, err
    }
    
    return application.DayRow{
        Date: date,
        Pages: pages,
        StreakDays: streak,
    }, nil
}
```

### Presentation Layer (internal/presentation/httpapi/)

**Responsabilidades:**
- Handlers HTTP (parse input, chamar use case, retornar response)
- Roteamento
- Extração de autenticação
- Transformação HTTP ↔ application DTOs

**Regras:**
- ✅ Handlers = `*Handler` com método `Handle(ctx, event)`
- ✅ Sempre extraia e valide autenticação primeiro
- ✅ Retorne `events.APIGatewayV2HTTPResponse`
- ✅ Use `Error()` para erros e `JSON()` para sucesso
- ❌ Nunca acesse DB diretamente
- ❌ Nunca valide ou processe lógica de domínio aqui

#### Padrão: BuildInput + Claims

**O `BuildInput` sempre:**
1. Extrai claims com `ExtractClaims(event)`
2. Parseia o body JSON
3. Retorna o DTO da application layer (que já contém as claims)

**Exemplo correto** (`register_reading_input.go`):
```go
type registerReadingBody struct {
	Pages int `json:"pages"`
}

func BuildRegisterReadingInput(event events.APIGatewayV2HTTPRequest) (app.RegisterReadingInput, error) {
	// Extract Claims from event
	claims, err := ExtractClaims(event)
	if err != nil {
		return app.RegisterReadingInput{}, err
	}

	// Parse body
	var body registerReadingBody
	if err := json.Unmarshal([]byte(event.Body), &body); err != nil {
		return app.RegisterReadingInput{}, errors.New("invalid request body")
	}

	pagesVO, err := readingDomain.NewPages(body.Pages)
	if err != nil {
		return app.RegisterReadingInput{}, err
	}

	return app.RegisterReadingInput{
		Claims: claims,
		Pages:  pagesVO,
	}, nil
}
```

**E o handler fica simples:**
```go
func (h *RegisterReadingHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// 1. Parse input (claims já incluídas)
	in, err := BuildRegisterReadingInput(event)
	if err != nil {
		return Error(event, http.StatusBadRequest, err.Error()), nil
	}

	// 2. Chamar use case
	out, err := h.uc.Execute(ctx, in)
	if err != nil {
		return Error(event, http.StatusInternalServerError, err.Error()), nil
	}

	// 3. Retornar
	return JSON(http.StatusOK, out), nil
}
```

**Exemplo:**
```go
// presentation/httpapi/register_reading_handler.go
type RegisterReadingHandler struct {
    uc *app.RegisterReadingUseCase
}

func (h *RegisterReadingHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
    // 1. Autenticação
    authInput, err := BuildEnsureMeInput(event)
    if err != nil {
        return Error(event, http.StatusUnauthorized, "unauthorized"), nil
    }
    
    // 2. Parse body
    var in app.RegisterReadingInput
    if err := json.Unmarshal([]byte(event.Body), &in); err != nil {
        return Error(event, http.StatusBadRequest, "invalid json"), nil
    }
    
    // 3. Chamar use case
    out, err := h.uc.Execute(ctx, in)
    if err != nil {

        return Error(event, http.StatusInternalServerError, err.Error()), nil
    }
    
    // 4. Retornar
    return JSON(http.StatusOK, out), nil
}
```

---

## ⚠️ Tratamento de Erros

### Padrão de Erros

#### 1. **Domain Errors** (errors.go em cada domain package)
```go
// domain/user/errors.go
var (
    ErrInvalidCognitoSub = errors.New("invalid cognito sub")
    ErrInvalidEmail      = errors.New("invalid email format")
    ErrInvalidName       = errors.New("name too long")
)
```

#### 2. **Application Errors** (tratados no use case)
```go
// application/user/ensure_me.go
if existing == nil {
    u := domain.NewFromIDP(in.Claims)
    if err := uc.repo.Insert(ctx, &u); err != nil {
        return MeDTO{}, fmt.Errorf("failed to insert user: %w", err)
    }
}
```

#### 3. **HTTP Errors** (mapeados em handlers)
```go
// Unauthorized
if err == ErrUnauthorized {
    return Error(event, http.StatusUnauthorized, "unauthorized"), nil
}

// Bad Request
if strings.Contains(err.Error(), "invalid") {
    return Error(event, http.StatusBadRequest, err.Error()), nil
}

// Server Error
return Error(event, http.StatusInternalServerError, err.Error()), nil
```

### Logs
Use `log.Printf()` via `Error()` helper:

```go
// response.go
func Error(event events.APIGatewayV2HTTPRequest, status int, msg string) events.APIGatewayV2HTTPResponse {
    reqID := event.RequestContext.RequestID
    log.Printf("[httpapi] Error reqId=%s status=%d msg=%s", reqID, status, msg)
    return JSON(status, map[string]string{"error": msg})
}
```

---

## 🧪 Testes & Desenvolvimento Local

### Desenvolvimento Local (SAM)

#### Setup:
```bash
# Build
make build

# Start local server
make run
# Acessar: http://localhost:3000
```

#### Testar endpoints:
```bash
# GET /v1/me (com JWT no header)
curl -X GET http://localhost:3000/v1/me \
  -H "Authorization: Bearer <token>"

# POST /v1/reading/logs
curl -X POST http://localhost:3000/v1/reading/logs \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"pages": 20}'

# GET /v1/reading/progress
curl -X GET http://localhost:3000/v1/reading/progress \
  -H "Authorization: Bearer <token>"

# PUT /v1/reading/goal
curl -X PUT http://localhost:3000/v1/reading/goal \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"pages": 100}'

# POST /v1/groups (criar novo grupo)
curl -X POST http://localhost:3000/v1/groups \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "React Lovers", "icon_id": "books"}'
```

### Variáveis de Ambiente
Criar `env.local` (não commitar):
```
DATABASE_URL=postgres://user:password@localhost:5432/reading_cats?sslmode=require
AWS_SAM_LOCAL=true
```

### Migrations Locais
```bash
# Criar nova migration
make migrate-create NAME=my_feature

# Rodar migrations UP
make migrate-up

# Rodar migrations DOWN
make migrate-down
```

---

## 🚀 Deployment

### Build & Deploy (GitHub Actions)

Ao fazer push para `main`:

1. **Build:** `sam build --use-container`
2. **Deploy:** `sam deploy` (via GitHub Secrets)
3. **Migrations:** Executar `golang-migrate` before/after deploy

### Checklist Pre-Deploy:
- [ ] Código passes linting (`go fmt`, `go vet`)
- [ ] Migrations estão criadas e testadas
- [ ] Variáveis de ambiente definidas em AWS Secrets Manager
- [ ] Testes locais passando
- [ ] Template.yaml atualizado com novos handlers

### Environment Variables (AWS):
```
DATABASE_URL=<neon-postgres-url>
TIMEZONE=America/Sao_Paulo
```

---

## ✅ Exemplo Prático: POST /v1/groups

### Estrutura Criada:

**Domain Layer** (`internal/domain/group/`)
- `group.go` → Entidade `Group` com construtor `New()`
- `value_objects.go` → VOs: `GroupName`, `IconID`, `Visibility` (com validação)
- `errors.go` → Erros: `ErrInvalidGroupName`, `ErrInvalidIconID`

**Application Layer** (`internal/application/group/`)
- `dto.go` → `CreateGroupInput` (contém `Claims` + dados do corpo), `CreateGroupOutput`
- `create_group.go` → `CreateGroupUseCase` que:
  - Busca o usuário pelo `CognitoSub` usando `userRepo.FindByCognitoSub(in.Claims.Sub)`
  - Obtém o `user.ID` (UUID) para usar como `CreatedByUserID`
  - Valida entrada (name, icon_id)
  - Cria entidade `Group` com visibility = INVITE_ONLY
  - Insere no DB
  - Adiciona creator como ADMIN em `group_members`
  - Recebe injeção de dependência: `groupRepo` + `userRepo`
- `repository.go` → Interface com `Insert()` e `AddMember()`

**Infrastructure Layer** (`internal/infra/group/`)
- `postgres_repository.go` → Implementação com queries diretos ao Postgres

**Presentation Layer** (`internal/presentation/httpapi/`)
- `create_group_input.go` → `BuildCreateGroupInput()` que:
  - Extrai claims com `ExtractClaims()`
  - Parseia body JSON
  - Retorna `app.CreateGroupInput` (já com claims)
- `create_group_handler.go` → Handler simples que recebe input pronto
- `router.go` → Rota `POST /v1/groups` adicionada
- `main.go` → Dependency injection com pool → repo → UC → handler

### Request:
```bash
POST /v1/groups
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "React Lovers",
  "icon_id": "books"
}
```

### Response (201 Created):
```json
{
  "id": "uuid-xxx",
  "name": "React Lovers",
  "icon_id": "books",
  "visibility": "INVITE_ONLY",
  "created_by_user_id": "uuid-yyy",
  "created_at": "2026-01-29T10:00:00Z",
  "updated_at": "2026-01-29T10:00:00Z"
}
```

---

## 📝 Resumo: Do Simples ao Complexo

### Fluxo de Desenvolvimento Típico:

1. **Identifique a Feature** → Ex: "Listar livros do usuário"

2. **Domain First:**
   - Crie `domain/book/book.go` (entidade)
   - Crie `domain/book/value_objects.go` (Title, ISBN, etc.)
   - Defina `domain/book/errors.go` (erros específicos)

3. **Application:**

   - Defina `application/book/repository.go` (interface)
   - Implemente `application/book/list_books.go` (use case)
   - Defina `application/book/dto.go` (Input/Output)

4. **Infrastructure:**
   - Implemente `infra/book/postgres_repository.go`
   - Escreva migration SQL em `migrations/000X_add_books.up.sql`

5. **Presentation:**
   - Crie `presentation/httpapi/list_books_handler.go`
   - Atualize `presentation/httpapi/router.go`

6. **Integração:**
   - Atualize `main.go` (dependency injection)
   - Teste localmente com `make run`
   - Deploy com `sam deploy`

---

## 🎓 Referências Rápidas

### Imports Comuns:
```go
import (
    "context"
    "errors"
    "json"
    "time"
    
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/aws/aws-lambda-go/events"
)
```

### Tipos Lambda:
```go
events.APIGatewayV2HTTPRequest   // Input
events.APIGatewayV2HTTPResponse  // Output
```

### Response Helpers:
```go
JSON(status int, body any) → success response
Error(event, status, msg) → error response com log
```

### Database Context:
```go
// Com transação
uc.repo.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
    // Lógica dentro da transação
    return nil
})

// Sem transação (for reads)
```

---

**Última atualização:** Janeiro 2026  
**Mantido por:** Reading Cats Team
