# Tabyte — Guia de Correção das Estimativas de Armazenamento para PostgreSQL e SQL Server

## Objetivo

Este documento orienta a correção do cálculo de `EstimatedBytes` por coluna no Tabyte para os módulos **PostgreSQL** e **SQL Server**, com base nas particularidades oficiais de armazenamento de tipos de dados. O objetivo é ajudar a Cursor a aplicar mudanças de código coerentes, sem continuar reutilizando a mesma heurística simplificada para engines diferentes [1][2].

O escopo deste guia é o cálculo de tamanho por tipo de dado em nível de coluna, além das limitações conhecidas dessa abordagem. O documento também aponta onde a implementação atual está errada, o que deve mudar imediatamente e o que deve ficar para uma segunda etapa em nível de linha e página [1][3][2].

## Problema atual

Hoje os módulos `postgres` e `sqlserver` usam praticamente a mesma função de estimativa, o que cria distorções importantes. Isso acontece porque PostgreSQL e SQL Server não compartilham o mesmo modelo físico para `varchar`, `numeric`, `date`, tipos de data/hora e armazenamento variável [1][3][2].

Os principais problemas da implementação atual são:

- tratar `varchar` como `n bytes fixos` nas duas engines;
- tratar `numeric` com uma fórmula genérica reaproveitada;
- assumir `date = 4` no SQL Server;
- reduzir `timestamp` a uma única regra, ignorando diferenças de tipos concretos no SQL Server;
- usar heurísticas de `text` sem distinguir armazenamento variável e sem documentar isso.

## Princípio central

A correção não deve tentar produzir exatidão absoluta só com metadado de schema, porque isso não é possível para tipos variáveis sem hipóteses sobre conteúdo real. O objetivo correto é implementar uma estimativa **documentada, reproduzível e aderente ao storage model da engine**, distinguindo entre:

- tipos de tamanho fixo;
- tipos de tamanho variável;
- tipos cujo custo depende do valor armazenado, não apenas da assinatura do schema [3][2].

## SQL Server

## Regras oficiais relevantes

### Inteiros

No SQL Server, `smallint` ocupa 2 bytes, `int` ocupa 4 bytes e `bigint` ocupa 8 bytes [4].

### decimal / numeric

No SQL Server, `decimal` e `numeric` são equivalentes, e o tamanho de armazenamento depende da **precisão**, com a seguinte tabela oficial [5][1]:

| Precisão | Bytes |
|---|---:|
| 1 a 9 | 5 |
| 10 a 19 | 9 |
| 20 a 28 | 13 |
| 29 a 38 | 17 |

A implementação atual está errada ao usar uma heurística própria para `numeric`, porque o SQL Server já fornece uma regra objetiva por faixa de precisão [5][1].

### char / varchar

No SQL Server, `char(n)` usa `n` bytes, enquanto `varchar(n)` usa o comprimento real armazenado em bytes **mais 2 bytes** [3]. Isso significa que `varchar(100)` não consome sempre 100 bytes; 100 é limite, não consumo fixo [3].

Portanto:

- `char(n)` pode ser estimado como `n`;
- `varchar(n)` precisa distinguir entre **limite declarado** e **tamanho médio assumido** [3].

### date

No SQL Server, `date` ocupa 3 bytes, não 4 [6].

### datetime2

No SQL Server, `datetime2` varia com a precisão fracionária. A documentação da Microsoft informa que `datetime2` ocupa de 6 a 8 bytes conforme a precisão [7].

Regra prática:

| Precisão fracionária | Bytes |
|---|---:|
| 0 a 2 | 6 |
| 3 a 4 | 7 |
| 5 a 7 | 8 |

Se o sistema normaliza tudo para `timestamp`, essa informação está sendo perdida [7].

### float / real

No SQL Server, `float` depende do parâmetro `n`: de 1 a 24 bits usa 4 bytes, de 25 a 53 bits usa 8 bytes [8]. Se o domínio já normaliza `float` como equivalente de dupla precisão, 8 bytes é aceitável como padrão documentado; caso contrário, será necessário carregar também a precisão binária [8].

### uniqueidentifier

`uniqueidentifier` ocupa 16 bytes [6][9].

### bit

O SQL Server agrupa colunas `bit` em nível de linha, e não como 1 byte independente por coluna. Como a implementação atual estima coluna isolada, `1 byte` pode continuar como aproximação local, mas a documentação do código deve deixar explícito que o ajuste correto ocorre no estimador de linha [6][10].

## Regra recomendada para SQL Server no Tabyte

### O que mudar agora

1. Corrigir `date` para 3 bytes [6].
2. Trocar `numeric` por tabela oficial de precisão [5][1].
3. Tratar `char` e `varchar` de forma diferente [3].
4. Preparar o domínio para diferenciar `datetime`, `datetime2`, `smalldatetime` e similares [7][6].
5. Documentar que `boolean` em SQL Server é aproximação para `bit`, não cálculo final físico [10].

### Implementação mínima sugerida

```go
func EstimateColumn(col domain.Column) domain.Column {
	var n int64

	switch col.NormalizedType {
	case "smallint":
		n = 2
	case "int":
		n = 4
	case "bigint":
		n = 8
	case "boolean":
		n = 1
	case "uuid":
		n = 16
	case "date":
		n = 3
	case "datetime":
		n = 8
	case "datetime2":
		n = estimateSQLServerDateTime2(col.Scale)
	case "float":
		n = 8
	case "char":
		if col.Length != nil {
			n = int64(*col.Length)
		}
	case "varchar":
		avg := assumedLength(col, 64)
		n = int64(avg) + 2
	case "text":
		n = 256 + 2
	case "numeric":
		n = estimateSQLServerNumeric(col.Precision)
	default:
		return col
	}

	col.EstimatedBytes = &n
	return col
}

func estimateSQLServerNumeric(precision *int) int64 {
	p := 18
	if precision != nil && *precision > 0 {
		p = *precision
	}
	switch {
	case p <= 9:
		return 5
	case p <= 19:
		return 9
	case p <= 28:
		return 13
	default:
		return 17
	}
}

func estimateSQLServerDateTime2(scale *int) int64 {
	s := 7
	if scale != nil && *scale >= 0 {
		s = *scale
	}
	switch {
	case s <= 2:
		return 6
	case s <= 4:
		return 7
	default:
		return 8
	}
}
```

## PostgreSQL

## Regras oficiais relevantes

### Inteiros

No PostgreSQL, `smallint` ocupa 2 bytes, `integer` ocupa 4 bytes e `bigint` ocupa 8 bytes [2].

### boolean

No PostgreSQL, `boolean` ocupa 1 byte [2].

### uuid

No PostgreSQL, `uuid` ocupa 16 bytes [11].

### date e timestamp

No PostgreSQL, `date` ocupa 4 bytes e `timestamp` ocupa 8 bytes [12].

### char / varchar / text

No PostgreSQL, `character(n)`, `varchar(n)` e `text` são **tipos de comprimento variável**. A documentação informa que o requisito de armazenamento é o tamanho real da string em bytes, mais 1 byte se a string tiver menos de 127 bytes, ou 4 bytes caso contrário [13][14].

Isso significa que `varchar(100)` não deve ser estimado como 100 bytes fixos. O cálculo precisa usar um comprimento assumido para o valor armazenado, não o limite declarado da coluna [13].

### numeric

No PostgreSQL, `numeric` tem armazenamento variável: **dois bytes para cada grupo de quatro dígitos decimais**, mais **3 a 8 bytes de overhead** [2]. A própria documentação enfatiza que o armazenamento físico depende do número real de dígitos armazenados, não simplesmente do modificador `numeric(p,s)` [2].

Isso significa que a precisão declarada ajuda a construir uma estimativa conservadora, mas não deve ser tratada como custo fixo equivalente ao SQL Server [2].

## Regra recomendada para PostgreSQL no Tabyte

### O que mudar agora

1. Não reutilizar a heurística de `numeric` do SQL Server [2][1].
2. Não usar `varchar(n) == n bytes` [13].
3. Tratar `text` como varlena, com tamanho assumido + overhead [13][14].
4. Manter `date = 4`, `timestamp = 8`, `uuid = 16`, `boolean = 1` [12][11][2].

### Implementação mínima sugerida

```go
func EstimateColumn(col domain.Column) domain.Column {
	var n int64

	switch col.NormalizedType {
	case "smallint":
		n = 2
	case "int":
		n = 4
	case "bigint":
		n = 8
	case "boolean":
		n = 1
	case "uuid":
		n = 16
	case "date":
		n = 4
	case "timestamp":
		n = 8
	case "float":
		n = 8
	case "char":
		if col.Length != nil {
			n = pgVarlenaSize(int64(*col.Length))
		}
	case "varchar":
		avg := assumedLength(col, 64)
		n = pgVarlenaSize(int64(avg))
	case "text":
		n = pgVarlenaSize(256)
	case "numeric":
		n = estimatePostgresNumeric(col.Precision)
	default:
		return col
	}

	col.EstimatedBytes = &n
	return col
}

func pgVarlenaSize(dataBytes int64) int64 {
	if dataBytes < 127 {
		return dataBytes + 1
	}
	return dataBytes + 4
}

func estimatePostgresNumeric(precision *int) int64 {
	p := 18
	if precision != nil && *precision > 0 {
		p = *precision
	}
	groups := (p + 3) / 4
	return int64(groups*2 + 4)
}
```

## Mudança importante no domínio

Para melhorar a qualidade das estimativas, o Tabyte não deve depender só de `Length`, `Precision` e `Scale`. Para tipos variáveis, a engine precisa distinguir **limite declarado** de **tamanho médio assumido** [3][13].

A recomendação é incluir no domínio algo como:

```go
type Column struct {
	NormalizedType   string
	Length           *int
	Precision        *int
	Scale            *int
	AssumedAvgLength *int
	EstimatedBytes   *int64
}
```

E a função auxiliar:

```go
func assumedLength(col domain.Column, fallback int) int {
	if col.AssumedAvgLength != nil && *col.AssumedAvgLength > 0 {
		return *col.AssumedAvgLength
	}
	if col.Length != nil && *col.Length > 0 {
		if *col.Length < fallback {
			return *col.Length
		}
		return *col.Length / 2
	}
	return fallback
}
```

Essa mudança evita inflar todas as colunas variáveis até o tamanho máximo declarado, que é exatamente o erro mais visível no modelo atual [3][13].

## Limitação inevitável da estimativa por coluna

Mesmo com essas correções, estimativa por coluna continua incompleta para representar o armazenamento real final. Isso acontece porque parte do custo físico depende de estrutura de linha, overhead de registro, alinhamento, null bitmap, agrupamento de bits, headers internos e efeitos de armazenamento fora da linha ou varlena [10][2].

Portanto, a implementação correta deve ser tratada em duas etapas:

1. **Etapa 1**: corrigir o custo base por tipo em nível de coluna.
2. **Etapa 2**: criar estimadores de **linha** e depois de **tabela** por engine.

## Recomendação de roadmap técnico

### Fase imediata

- corrigir SQL Server `date`, `numeric`, `varchar` e `datetime2` [6][5][3][7];
- corrigir PostgreSQL `varchar`, `text` e `numeric` [13][2];
- introduzir `AssumedAvgLength` no domínio [3][13].

### Fase seguinte

- criar `EstimateRow(table)` por engine;
- calcular overhead fixo e variável fora da função por coluna;
- ajustar agrupamento de colunas `bit` no SQL Server [10];
- separar `text`/LOB/out-of-row/TOAST em tratamento explícito quando o produto evoluir [15][2].

## Instrução direta para a Cursor

Aplicar as seguintes mudanças no código:

1. Substituir a heurística atual de `numeric` do SQL Server por tabela oficial por precisão [5][1].
2. Corrigir `date` do SQL Server para 3 bytes [6].
3. Introduzir distinção entre `char` e `varchar` no SQL Server, usando `+2` bytes para `varchar` [3].
4. Parar de tratar `varchar` e `text` do PostgreSQL como tamanho fixo máximo; usar tamanho assumido + overhead varlena [13][14].
5. Substituir a heurística de `numeric` do PostgreSQL por fórmula baseada em grupos de 4 dígitos + overhead [2].
6. Adicionar `AssumedAvgLength` ao modelo de coluna e priorizar esse valor quando disponível [3][13].
7. Preparar uma segunda etapa para estimativa em nível de linha, sem tentar resolver row overhead dentro de `EstimateColumn` [10][2].

## Critério de aceite

A mudança pode ser considerada correta quando:

- SQL Server não usar mais a mesma regra de `numeric` do PostgreSQL [5][2];
- `date` do SQL Server retornar 3 bytes [6];
- PostgreSQL deixar de usar `varchar(n) = n` como regra fixa [13];
- o código documentar claramente o que é cálculo por coluna e o que depende de cálculo por linha [10][2].