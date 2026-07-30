# Tabyte — Endpoints da API Local v0.2

## Objetivo

Este documento define a **API HTTP local** do Tabyte sob a ótica de contrato técnico. O foco é descrever recursos, rotas, verbos HTTP, formatos de payload, convenções de resposta e organização dos endpoints expostos pelo servidor local da aplicação, sem transformar o documento em especificação completa de regras de negócio [1][2][3].

O cenário assumido é o de uma API servida apenas em ambiente local, por um processo iniciado via CLI, usando HTTP como contrato entre a interface web e o backend em Go [3].

## Convenções gerais

### Base URL

A API deve ser exposta localmente sob um prefixo versionado.

```text
http://127.0.0.1:{porta}/api/v1
```

O uso de prefixo versionado reduz quebra de contrato quando a API evoluir ao longo das versões do produto [1][2].

### Formato

- Requisições e respostas usam `application/json` por padrão.
- Endpoints de upload podem aceitar `multipart/form-data` quando necessário.
- Todas as respostas JSON devem usar envelope consistente.

### Envelope de resposta

#### Sucesso

```json
{
  "data": {},
  "meta": {
    "request_id": "req_123",
    "timestamp": "2026-07-30T11:30:00Z"
  },
  "error": null
}
```

#### Erro

```json
{
  "data": null,
  "meta": {
    "request_id": "req_123",
    "timestamp": "2026-07-30T11:30:00Z"
  },
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "payload inválido",
    "details": []
  }
}
```

### Semântica HTTP

A API deve respeitar a semântica básica dos métodos HTTP: `GET` para leitura, `POST` para criação ou processamento, `PUT` para substituição idempotente, `PATCH` para atualização parcial e `DELETE` para remoção. A semântica de segurança e idempotência desses métodos é parte do modelo HTTP definido em RFC 9110 e material técnico de referência [1][2].

## Recursos principais

A API local é organizada em seis grupos de recursos:

1. Sistema.
2. Configurações locais.
3. Sessões de análise.
4. Resultados estruturais.
5. Alertas.
6. Extensões futuras.

## 1. Sistema

### `GET /health`

Retorna estado básico do processo local.

#### Resposta esperada

```json
{
  "data": {
    "status": "ok",
    "app": "tabyte",
    "version": "0.1.0"
  },
  "meta": {
    "request_id": "req_1",
    "timestamp": "2026-07-30T11:30:00Z"
  },
  "error": null
}
```

### `GET /info`

Retorna informações técnicas da instância local, como versão, modo de execução e capacidades habilitadas.

## 2. Configurações locais

### `GET /settings`

Lista configurações efetivas relevantes ao runtime local.

### `PUT /settings/{key}`

Substitui o valor de uma configuração específica. O uso de `PUT` é adequado quando a intenção é sobrescrever o estado representado pelo recurso de maneira idempotente [1][2].

#### Exemplo de payload

```json
{
  "value": "postgres"
}
```

### `PATCH /settings`

Atualiza parcialmente um conjunto pequeno de configurações locais.

#### Exemplo de payload

```json
{
  "default_engine": "sqlserver",
  "theme": "dark"
}
```

## 3. Sessões de análise

### `POST /analysis-sessions`

Cria uma nova sessão de análise local. `POST` é o verbo mais adequado para criação de recurso ou submissão de processamento que produz novo estado [1][2].

#### Payload sugerido

```json
{
  "engine": "sqlserver",
  "source_name": "schema-produto.sql",
  "ddl_text": "CREATE TABLE Product (...)",
  "options": {
    "persist": true
  }
}
```

#### Resposta sugerida

```json
{
  "data": {
    "id": "as_01JXYZ",
    "engine": "sqlserver",
    "status": "created"
  },
  "meta": {
    "request_id": "req_2",
    "timestamp": "2026-07-30T11:30:00Z"
  },
  "error": null
}
```

### `GET /analysis-sessions`

Lista sessões persistidas localmente, com suporte futuro a filtros simples.

#### Parâmetros opcionais

- `engine`
- `limit`
- `offset`
- `q`

### `GET /analysis-sessions/{sessionId}`

Retorna os metadados principais de uma sessão específica.

### `DELETE /analysis-sessions/{sessionId}`

Remove uma sessão persistida. `DELETE` é apropriado para remoção de recurso identificado [1][2].

## 4. Resultados estruturais

### `GET /analysis-sessions/{sessionId}/summary`

Retorna o resumo estrutural agregado de uma sessão.

#### Resposta sugerida

```json
{
  "data": {
    "session_id": "as_01JXYZ",
    "engine": "sqlserver",
    "total_tables": 12,
    "estimated_total_bytes": 102400,
    "warning_count": 3,
    "error_count": 0
  },
  "meta": {
    "request_id": "req_3",
    "timestamp": "2026-07-30T11:30:00Z"
  },
  "error": null
}
```

### `GET /analysis-sessions/{sessionId}/tables`

Lista os resultados por tabela pertencentes à sessão.

### `GET /analysis-sessions/{sessionId}/tables/{tableId}`

Retorna o detalhamento de uma tabela persistida dentro da sessão.

### `GET /analysis-sessions/{sessionId}/tables/{tableId}/columns`

Lista o detalhamento persistido das colunas de uma tabela.

### `GET /analysis-sessions/{sessionId}/tables/{tableId}/indexes`

Lista índices associados à tabela persistida.

## 5. Alertas

### `GET /analysis-sessions/{sessionId}/warnings`

Lista alertas associados à sessão inteira.

### `GET /analysis-sessions/{sessionId}/tables/{tableId}/warnings`

Lista alertas associados a uma tabela específica.

## 6. Entrada por arquivo

### `POST /imports/sql`

Recebe um arquivo SQL por upload para criação de nova sessão. Como envolve submissão de conteúdo e potencial criação de recurso novo, `POST` permanece o verbo mais adequado [1][2].

#### Formato

`multipart/form-data`

#### Campos sugeridos

- `file`
- `engine`
- `source_name`
- `persist`

## 7. Extensões futuras

Esses endpoints não precisam existir no primeiro ciclo, mas a estrutura pode reservá-los conceitualmente.

### `GET /analysis-sessions/{sessionId}/insights`

Recurso opcional para listar insights auxiliares gerados por mecanismo externo.

### `POST /exports`

Recurso opcional para geração de exportações estruturadas da sessão.

## Status codes sugeridos

| Situação | Status |
|---|---|
| Leitura bem-sucedida | `200 OK` |
| Criação bem-sucedida | `201 Created` |
| Remoção sem corpo | `204 No Content` |
| Payload inválido | `400 Bad Request` |
| Recurso não encontrado | `404 Not Found` |
| Conflito lógico local | `409 Conflict` |
| Erro inesperado | `500 Internal Server Error` |

Essa convenção segue a semântica geral esperada para APIs HTTP orientadas a recurso [1][2].

## Convenções de nomenclatura

- URLs em kebab-case.
- Recursos no plural quando representam coleções.
- IDs como segmentos de rota.
- Sem verbos na URL, exceto quando o recurso representar de fato uma operação de sistema inevitavelmente orientada a ação.

## Ordem sugerida de implementação

Para um primeiro ciclo enxuto, a API pode começar com:

- `GET /health`
- `GET /info`
- `POST /analysis-sessions`
- `GET /analysis-sessions`
- `GET /analysis-sessions/{sessionId}`
- `GET /analysis-sessions/{sessionId}/summary`
- `GET /analysis-sessions/{sessionId}/tables`
- `GET /analysis-sessions/{sessionId}/warnings`
- `DELETE /analysis-sessions/{sessionId}`

Configurações, upload de arquivo, colunas, índices e extensões podem entrar depois, à medida que o modelo persistido amadurecer.

## Observação final

A principal revisão aplicada neste documento foi manter o foco em **contrato HTTP local**: recursos, rotas, verbos, payloads e convenções. Regras funcionais detalhadas, critérios analíticos e comportamento interno do domínio devem permanecer em documentos específicos de requisitos e aplicação, não na especificação de endpoints [1][3].