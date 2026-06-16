# Gblog API

API REST para gerenciamento de artigos de um blog, construída com **Go**, **Gin**, **GORM** e **PostgreSQL**, seguindo os princípios de **Clean Architecture**.

---

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Arquitetura](#arquitetura)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Tecnologias](#tecnologias)
- [Pré-requisitos](#pré-requisitos)
- [Instalação](#instalação)
- [Variáveis de Ambiente](#variáveis-de-ambiente)
- [Execução](#execução)
- [Endpoints da API](#endpoints-da-api)
- [Modelos de Dados](#modelos-de-dados)
- [Regras de Negócio](#regras-de-negócio)
- [Documentação Swagger](#documentação-swagger)

---

## Visão Geral

O **Gblog** é uma API para criação, edição, publicação e exclusão de posts de blog. O sistema implementa um fluxo completo de ciclo de vida dos posts, desde a criação como rascunho até a publicação ou exclusão lógica.

### Funcionalidades

- **Criar posts** como rascunho
- **Editar** título, slug e conteúdo
- **Publicar** posts com validação de regras de negócio
- **Exclusão lógica** (soft delete) de posts
- **Documentação Swagger** automática

---

## Arquitetura

O projeto segue os princípios de **Clean Architecture**, separando responsabilidades em camadas distintas:

```
┌─────────────────────────────────────────────┐
│              HTTP Handlers (Delivery)        │
│         Camada de apresentação / API         │
├─────────────────────────────────────────────┤
│              Use Cases (Business)            │
│           Lógica de negócio                  │
├─────────────────────────────────────────────┤
│              Repository (Data)               │
│          Acesso a dados (GORM/PostgreSQL)    │
└─────────────────────────────────────────────┘
```

- **Entities**: Definem as regras de negócio e estruturas de dados
- **Use Cases**: Orquestram as operações de negócio
- **Repositories**: Abstraem o acesso ao banco de dados
- **Handlers**: Expõem a API HTTP

---

## Estrutura do Projeto

```
Gblog/
├── cmd/
│   ├── docs/                    # Documentação Swagger gerada
│   └── main.go                  # Ponto de entrada da aplicação
├── internal/
│   ├── blogPost/
│   │   ├── entity.go            # Entidade BlogPost
│   │   ├── dto.go               # Data Transfer Objects
│   │   ├── repository.go        # Interface e implementação do repositório
│   │   ├── usecase.go           # Casos de uso (Create, Publish, Update, Delete)
│   │   └── http/
│   │       └── handlers.go      # Handlers HTTP (rotas)
│   ├── infra/
│   │   └── database.go          # Conexão com banco e configuração
│   └── user/
│       ├── entity.go            # Entidade User (em desenvolvimento)
│       ├── dto.go               # DTOs do usuário
│       ├── repository.go        # Repositório do usuário
│       └── usecase.go           # Caso de uso do usuário
├── pkg/
│   └── blogStatus.go            # Enum de status do blog
├── .env                         # Variáveis de ambiente (não versionado)
├── .gitignore                   # Arquivos ignorados pelo Git
├── go.mod                       # Dependências do módulo Go
└── go.sum                       # Sumário de dependências
```

---

## Tecnologias

| Tecnologia | Versão | Descrição |
|------------|--------|-----------|
| Go | 1.26.3 | Linguagem de programação |
| Gin | v1.12.0 | Framework HTTP |
| GORM | v1.31.1 | ORM para Go |
| PostgreSQL | - | Banco de dados relacional |
| Swag | v1.16.6 | Geração de documentação Swagger |

---

## Pré-requisitos

- [Go](https://go.dev/dl/) 1.26.3 ou superior
- [PostgreSQL](https://www.postgresql.org/) instalado e configurado
- [Swag CLI](https://github.com/swaggo/swag) (para gerar documentação)

---

## Instalação

1. **Clone o repositório:**

   ```bash
   git clone https://github.com/seu-usuario/Gblog.git
   cd Gblog
   ```

2. **Instale as dependências:**

   ```bash
   go mod download
   ```

3. **Crie o arquivo `.env` na raiz do projeto:**

   ```bash
   cp .env.example .env
   ```

4. **Configure as variáveis de ambiente** (veja seção [Variáveis de Ambiente](#variáveis-de-ambiente)).

5. **Gere a documentação Swagger:**

   ```bash
   swag init -g cmd/main.go
   ```

---

## Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `API_PORT` | Porta do servidor HTTP | `8080` |
| `DB_HOST` | Host do PostgreSQL | `localhost` |
| `DB_USER` | Usuário do banco | `postgres` |
| `DB_PASSWORD` | Senha do banco | `senha123` |
| `DB_NAME` | Nome do banco de dados | `gblog` |
| `DB_PORT` | Porta do PostgreSQL | `5432` |

---

## Execução

```bash
# Executar a aplicação
go run cmd/main.go

# Ou compilar e executar
go build -o gblog cmd/main.go
./gblog
```

A API estará disponível em: `http://localhost:{API_PORT}`

---

## Endpoints da API

Base URL: `/api/v1`

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/posts` | Criar um novo post (rascunho) |
| `PUT` | `/posts/:id` | Atualizar título, slug e conteúdo |
| `PATCH` | `/posts/:id/publish` | Publicar um post |
| `DELETE` | `/posts/:id` | Excluir um post (soft delete) |

### Detalhes dos Endpoints

#### Criar Post

```http
POST /api/v1/posts
Content-Type: application/json

{
  "title": "Meu Primeiro Post",
  "slug": "meu-primeiro-post",
  "content": "Conteúdo do post aqui..."
}
```

**Resposta:** `201 Created`

#### Publicar Post

```http
PATCH /api/v1/posts/{id}/publish
```

**Resposta:** `200 OK`

#### Atualizar Post

```http
PUT /api/v1/posts/{id}
Content-Type: application/json

{
  "title": "Título Atualizado",
  "slug": "titulo-atualizado",
  "content": "Conteúdo atualizado..."
}
```

**Resposta:** `200 OK`

#### Excluir Post

```http
DELETE /api/v1/posts/{id}
```

**Resposta:** `200 OK`

---

## Modelos de Dados

### BlogPost

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `ID` | `uint` | Identificador único (auto-incremento) |
| `Title` | `string` | Título do post (até 255 caracteres) |
| `Content` | `string` | Conteúdo do post (texto longo) |
| `Status` | `BlogStatus` | Status do post: `draft`, `published`, `archived` |
| `Slug` | `string` | URL amigável (única) |
| `IsActive` | `bool` | Indica se o post está ativo (soft delete) |
| `CreatedAt` | `time.Time` | Data de criação |
| `UpdatedAt` | `time.Time` | Data da última atualização |
| `DeletedAt` | `time.Time` | Data da exclusão lógica |

### Status do Post

| Status | Descrição |
|--------|-----------|
| `draft` | Rascunho (status padrão na criação) |
| `published` | Publicado |
| `archived` | Arquivado |

---

## Regras de Negócio

### Criação
- O título é obrigatório
- O status inicial é sempre `draft`

### Publicação
- O post deve existir e estar ativo
- O post não pode já estar publicado
- O conteúdo deve ter no mínimo 10 caracteres

### Edição
- O post deve existir e estar ativo
- Não é possível editar posts inativos/deletados

### Exclusão
- A exclusão é lógica (soft delete)
- O campo `IsActive` é definido como `false`
- O GORM preenche automaticamente `DeletedAt`

---

## Documentação Swagger

A documentação Swagger é gerada automaticamente e está disponível em:

```
http://localhost:{API_PORT}/swagger/index.html
```

Para regenerar a documentação após alterações nos comentários `godoc`:

```bash
swag init -g cmd/main.go
```

---

## Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
