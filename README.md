# PhraseForge

PhraseForge e uma aplicacao de terminal escrita em Go que gera frases aleatorias por categoria.

As frases sao formadas a partir de um template e de partes como sujeito, verbo e complemento. Os dados sao carregados de um arquivo JSON.

Veja o [roadmap do projeto](ROADMAP.md) para as proximas fases e a decisao de persistencia.

## Requisitos

- Go 1.26.5 ou compativel

## Executar

Execute os comandos a partir da raiz do repositorio. Por padrao, a aplicacao usa o caminho relativo `data/phrases.json`.

```bash
go run ./cmd/phraseforge help
```

## Comandos

Formato geral:

```bash
phraseforge [--data-file PATH] <command>
```

Para usar outro arquivo JSON, informe a flag global antes do comando:

```bash
go run ./cmd/phraseforge --data-file outro-arquivo.json categories
```

Gerar uma frase usando a categoria padrao, `programming`:

```bash
go run ./cmd/phraseforge generate
```

Gerar uma frase de uma categoria especifica:

```bash
go run ./cmd/phraseforge generate --category study
```

Gerar varias frases:

```bash
go run ./cmd/phraseforge generate --category programming --count 3
```

Listar categorias disponiveis:

```bash
go run ./cmd/phraseforge categories
```

Exibir ajuda:

```bash
go run ./cmd/phraseforge help
```

## API HTTP

Inicie a API a partir da raiz do repositorio:

```bash
go run ./cmd/phraseforge-api
```

Por padrao, ela fica disponivel em `http://localhost:8080`. Use `Ctrl+C` para encerrar o servidor de forma controlada.

### Configuracao

Variaveis de ambiente disponiveis:

- `PHRASEFORGE_API_ADDR`: endereco em que a API escuta. O padrao e `:8080`.
- `PHRASEFORGE_DATA_FILE`: caminho do arquivo JSON de categorias. O padrao e `data/phrases.json`.

Exemplo com configuracao customizada:

```bash
PHRASEFORGE_API_ADDR=:9090 \
PHRASEFORGE_DATA_FILE=data/phrases.json \
go run ./cmd/phraseforge-api
```

### Rotas

Health check:

```bash
curl http://localhost:8080/health
```

```json
{"status":"ok"}
```

Listar categorias:

```bash
curl http://localhost:8080/categories
```

```json
{"categories":["programming","study"]}
```

Gerar uma frase da categoria padrao, `programming`:

```bash
curl http://localhost:8080/phrases/random
```

```json
{
  "category": "programming",
  "phrases": [
    "Codigo simples reduz problemas futuros."
  ]
}
```

Gerar frases de uma categoria especifica:

```bash
curl 'http://localhost:8080/phrases/random?category=study'
```

Gerar varias frases:

```bash
curl 'http://localhost:8080/phrases/random?category=study&count=3'
```

```json
{
  "category": "study",
  "phrases": [
    "Com foco, a pratica constante fortalece o aprendizado, passo a passo.",
    "A cada dia, a revisao diaria melhora o raciocinio logico, com consistencia.",
    "Com foco, a revisao diaria melhora o aprendizado, passo a passo."
  ]
}
```

### Parametros

`GET /phrases/random` aceita os seguintes query parameters:

- `category`: opcional. Sem o parametro, a API usa `programming`. Valor vazio retorna `400`; categoria inexistente retorna `404`.
- `count`: opcional. O padrao e `1`; informe um inteiro entre `1` e `10`. Valores vazios, repetidos ou fora desse intervalo retornam `400`.

### Erros

Erros da API usam JSON:

```json
{"error":"category not found"}
```

Status mais comuns:

- `400 Bad Request`: parametro invalido.
- `404 Not Found`: rota ou categoria inexistente.
- `405 Method Not Allowed`: metodo HTTP nao suportado. A resposta inclui `Allow: GET`.
- `500 Internal Server Error`: estado interno inconsistente ou falha inesperada de geracao.

## Docker

Crie a imagem da API:

```bash
docker build -t phraseforge-api .
```

Execute a API na porta padrao:

```bash
docker run --rm -p 8080:8080 phraseforge-api
```

Em outro terminal, verifique o health check:

```bash
curl http://localhost:8080/health
```

Para configurar outro endereco dentro do container:

```bash
docker run --rm \
  -p 9090:9090 \
  -e PHRASEFORGE_API_ADDR=:9090 \
  phraseforge-api
```

O container inclui `data/phrases.json`. Para usar outro arquivo, monte um volume e configure `PHRASEFORGE_DATA_FILE`:

```bash
docker run --rm \
  -p 8080:8080 \
  -v "$(pwd)/data/custom-phrases.json:/app/data/custom-phrases.json:ro" \
  -e PHRASEFORGE_DATA_FILE=/app/data/custom-phrases.json \
  phraseforge-api
```

## Dados

As categorias ficam em `data/phrases.json`:

```json
{
  "categories": [
    {
      "name": "programming",
      "template": "{subject} {verb} {complement}",
      "subjects": ["Codigo simples"],
      "verbs": ["reduz"],
      "complements": ["problemas futuros"]
    }
  ]
}
```

Cada categoria precisa ter nome, template e listas nao vazias de `subjects`, `verbs` e `complements`.

O template deve conter exatamente uma ocorrencia dos placeholders obrigatorios:

- `{subject}`
- `{verb}`
- `{complement}`

Os placeholders opcionais podem aparecer no maximo uma vez:

- `{introduction}`
- `{conclusion}`

Quando um placeholder opcional e usado, sua lista correspondente no JSON deve conter ao menos um valor.

## Desenvolvimento

Formatar os arquivos Go:

```bash
gofmt -w cmd/phraseforge/*.go internal/phrase/*.go internal/storage/*.go
```

Executar testes:

```bash
go test ./...
```

Executar analise estatica:

```bash
go vet ./...
```

## Estrutura

```text
cmd/phraseforge/  Interface de linha de comando
internal/phrase/  Dominio de categorias, templates e geracao
internal/storage/ Leitura e validacao do arquivo JSON
data/             Dados usados pela aplicacao
```
