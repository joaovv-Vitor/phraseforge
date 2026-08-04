# PhraseForge

PhraseForge e uma aplicacao de terminal escrita em Go que gera frases aleatorias por categoria.

As frases sao formadas a partir de um template e de partes como sujeito, verbo e complemento. Os dados sao carregados de um arquivo JSON.

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
