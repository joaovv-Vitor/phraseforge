# PhraseForge

PhraseForge e uma aplicacao de terminal escrita em Go que gera frases aleatorias por categoria.

As frases sao formadas a partir de um template e de partes como sujeito, verbo e complemento. Os dados sao carregados de um arquivo JSON.

## Requisitos

- Go 1.26.5 ou compativel

## Executar

Execute os comandos a partir da raiz do repositorio, pois o arquivo de dados usa o caminho relativo `data/phrases.json`.

```bash
go run ./cmd/phraseforge help
```

## Comandos

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

O template deve conter exatamente uma ocorrencia de cada placeholder aceito:

- `{subject}`
- `{verb}`
- `{complement}`

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
