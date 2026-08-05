# PhraseForge Roadmap

## Concluido

- CLI para listar categorias e gerar frases aleatorias.
- Dados JSON, templates e partes opcionais.
- API HTTP com health check, categorias e geracao em lote.
- Configuracao por flags e variaveis de ambiente.
- Timeouts, graceful shutdown e respostas de erro JSON.
- Imagem Docker para a API.

## Persistencia Atual

SQLite sera a primeira persistencia do projeto.

Motivos:

- Banco SQL embutido, sem processo de servidor adicional.
- Permite praticar `database/sql`, SQL, migrations, constraints e transacoes.
- O modelo relacional podera ser migrado para PostgreSQL quando houver uma necessidade real de operacao multi-instancia ou maior concorrencia de escrita.

### Proximas Etapas

1. Criar e validar o schema SQLite inicial.
2. Conectar a aplicacao com `database/sql`.
3. Implementar repositorio SQLite para listar categorias.
4. Migrar o carregamento da API de JSON para SQLite.
5. Persistir templates, partes, favoritos e historico de geracao.

## Evolucoes Futuras

- Avaliar PostgreSQL quando SQLite deixar de atender aos requisitos operacionais.
- Adicionar migrations automatizadas.
- Adicionar logs estruturados, CI e metricas.
- Avaliar concorrencia para lotes maiores somente apos medir beneficio.
- Avaliar arquitetura distribuida somente se o monolito apresentar um problema concreto que a separacao resolva.

## Validar Migrations Localmente

Aplicar a migration:

```bash
sqlite3 /tmp/phraseforge.db < db/migrations/000001_initial_schema.up.sql
```

Ao usar chaves estrangeiras, habilite a verificacao na conexao SQLite:

```sql
PRAGMA foreign_keys = ON;
```

Reverter a migration:

```bash
sqlite3 /tmp/phraseforge.db < db/migrations/000001_initial_schema.down.sql
```
