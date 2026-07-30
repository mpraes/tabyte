# Tabyte — Requisitos Funcionais v0.2

## Visão do produto

O **Tabyte** é uma ferramenta open source para estimativa de armazenamento e análise estrutural de impacto em performance a partir de DDLs relacionais. A primeira linha de suporte deve contemplar **SQL Server** e **PostgreSQL**, respeitando as diferenças documentadas entre tipos de dados, estrutura de tabela e indexação de cada engine [1][2][3].

A solução deverá ser executada localmente por meio de uma **CLI cross-platform** que inicializa um servidor HTTP local e abre a interface web no navegador padrão do usuário. Esse modelo mantém paridade operacional entre Windows PowerShell e Ubuntu/Linux, além de preservar simplicidade de distribuição e uso local-first [4][5][6].

## Objetivo funcional

O sistema deverá receber DDLs, interpretar tabelas e objetos relevantes do schema, aplicar regras de estimativa específicas da engine selecionada e apresentar resultados compreensíveis ao usuário. Para SQL Server, a documentação oficial trata a estimativa como composição da estrutura base com índices adicionais, enquanto o PostgreSQL documenta índices como recurso de performance com custo adicional de manutenção e armazenamento [1][3][7].

## Escopo da primeira versão

A v0.2 inicial do produto deverá cobrir análise de DDL por texto, seleção de engine, cálculo estimado por coluna, linha, tabela e schema, alertas estruturais básicos, projeção de crescimento simples e persistência local opcional para histórico e configurações. O sistema não deverá depender de banco servidor externo e poderá usar SQLite local apenas como mecanismo interno leve quando a persistência estiver ativada [8][9].

Permanecem fora do escopo inicial funcionalidades como benchmark real, leitura de planos de execução, conexão obrigatória a bancos remotos, múltiplos SGBDs além de SQL Server e PostgreSQL e tuning automatizado avançado [1][3].

## Atores

### Usuário analista

Pessoa que cola ou edita um DDL e deseja compreender impacto estimado de modelagem em armazenamento e sinais estruturais de impacto em performance.

### Usuário arquiteto ou DBA

Pessoa que deseja comparar decisões de schema, revisar tipos de dados, estimar crescimento e avaliar implicações estruturais de índices, largura de linha e variação de modelagem [3][7].

### Usuário CLI

Pessoa que executa a aplicação localmente por terminal, iniciando o servidor HTTP local e acessando a UI via navegador. Esse ator deve conseguir iniciar o produto em Windows PowerShell ou Linux com um fluxo uniforme [4][5].

## Requisitos funcionais

### RF-01 — Iniciar aplicação por CLI

O sistema deverá disponibilizar um comando principal para inicializar a aplicação localmente, por exemplo `tabyte serve`. Ao iniciar com sucesso, o sistema deverá informar a URL local ativa e, quando configurado, abrir automaticamente o navegador padrão do usuário [4][5].

### RF-02 — Operar em localhost

O servidor HTTP da aplicação deverá operar por padrão apenas em `127.0.0.1` ou `localhost`, evitando exposição de rede externa desnecessária. A configuração padrão deverá privilegiar segurança e uso local [6].

### RF-03 — Informar DDL por texto

O sistema deverá permitir que o usuário cole ou edite manualmente um DDL em uma área de texto para análise. O conteúdo informado deverá ser preservado na sessão corrente até nova limpeza ou substituição manual.

### RF-04 — Selecionar engine de banco de dados

O sistema deverá permitir que o usuário escolha explicitamente a engine de análise entre **SQL Server** e **PostgreSQL**. As regras de parsing, mapeamento de tipos e cálculo deverão ser ajustadas conforme a engine selecionada [1][2].

### RF-05 — Validar entrada mínima

O sistema deverá validar se há conteúdo suficiente para análise antes do processamento. Em caso de DDL vazio, incompleto ou semanticamente insuficiente, o sistema deverá apresentar mensagem de erro clara.

### RF-06 — Interpretar estruturas de tabela

O sistema deverá identificar no DDL, ao menos, tabelas, colunas, tipos de dados, nullability, chave primária e índices suportados pelo parser da versão inicial. Quando encontrar estruturas não suportadas, o sistema deverá sinalizar a limitação sem interromper toda a análise, sempre que possível [1][2][3].

### RF-07 — Normalizar tipos de dados por engine

O sistema deverá converter a representação textual do tipo encontrado no DDL para um modelo interno normalizado, preservando comprimento, precisão, escala e demais atributos relevantes para cálculo. Essa normalização é necessária porque SQL Server e PostgreSQL definem tipos com semânticas distintas [1][2].

### RF-08 — Estimar tamanho por coluna

O sistema deverá calcular uma estimativa de armazenamento por coluna, considerando tipo normalizado, parâmetros associados e regras da engine selecionada. O resultado por coluna deverá ser apresentado ao usuário em bytes estimados ou unidade derivada [1][2].

### RF-09 — Estimar tamanho por linha

O sistema deverá calcular o tamanho estimado de cada linha da tabela, incluindo payload das colunas e overhead modelado pela engine. No caso de SQL Server, a composição deverá refletir a estimativa da estrutura base antes da soma de índices adicionais [1].

### RF-10 — Estimar volume por tabela

O sistema deverá calcular o volume estimado por tabela a partir do tamanho de linha estimado multiplicado pela quantidade de registros informada pelo usuário. O sistema deverá aceitar valor padrão configurável para cardinalidade quando o usuário ainda não tiver informado esse dado [10].

### RF-11 — Estimar volume total do schema

O sistema deverá consolidar o volume estimado total do schema analisado, somando as tabelas identificadas e demais componentes considerados pela engine. O total deverá ser exibido em unidade legível, como KB, MB ou GB.

### RF-12 — Permitir informar quantidade de linhas por tabela

O sistema deverá permitir que o usuário informe ou altere a quantidade de linhas esperada para cada tabela individualmente. O recálculo deverá ocorrer de forma imediata ou mediante ação explícita de atualização, conforme a UX adotada [10].

### RF-13 — Projetar crescimento

O sistema deverá permitir projetar crescimento por tabela com base em taxa simples informada pelo usuário, como crescimento por hora, dia ou mês. O objetivo é apoiar cenários iniciais de volumetria, e não modelagem temporal avançada [10].

### RF-14 — Exibir detalhamento por tabela

O sistema deverá apresentar, para cada tabela, nome, número de colunas interpretadas, bytes por linha estimados, quantidade de linhas adotada, volume total estimado e lista das colunas com seus tamanhos individuais.

### RF-15 — Exibir memória de cálculo

O sistema deverá mostrar como o valor final foi obtido, discriminando payload de colunas, overhead de linha e contribuição estimada de índices quando aplicável. Essa memória de cálculo é parte funcional central do produto [1].

### RF-16 — Interpretar índices suportados

O sistema deverá reconhecer, no escopo da versão inicial, índices definidos no DDL ou elementos suficientes para inferência básica de estruturas indexadas, como chaves primárias e índices explícitos. No PostgreSQL, os índices devem ser tratados como estruturas que aceleram determinadas buscas, mas adicionam overhead [3][7].

### RF-17 — Estimar impacto estrutural de índices

O sistema deverá fornecer estimativa básica do impacto estrutural de índices sobre armazenamento e custo potencial de manutenção. No SQL Server, a estimativa deverá seguir a separação entre estrutura base e índices não clusterizados; no PostgreSQL, os índices deverão ser tratados como custo adicional associado ao desempenho de leitura [1][3].

### RF-18 — Emitir alertas estruturais de modelagem

O sistema deverá gerar alertas baseados em regras para situações como colunas potencialmente largas, tipos excessivamente genéricos, uso de precisão desnecessária, excesso de colunas variáveis e estruturas propensas a aumento de volume ou custo de indexação [1][2][3].

### RF-19 — Emitir sinais de impacto em performance

O sistema deverá informar sinais estruturais que possam afetar leitura, escrita, ordenação ou manutenção de índices, sem afirmar desempenho real da carga de trabalho. Esse requisito está alinhado ao fato de que índices podem beneficiar filtros e ordenações, mas também introduzem custo adicional [3][7][11].

### RF-20 — Exportar resultados

O sistema deverá permitir exportar os resultados da análise em formato estruturado, preferencialmente JSON e CSV, para reaproveitamento em documentação, comparação ou automação externa.

### RF-21 — Reprocessar análise após ajuste de parâmetros

O sistema deverá permitir ao usuário alterar engine, volume de linhas ou outros parâmetros suportados e recalcular a análise sem necessidade de recarregar a aplicação inteira.

### RF-22 — Persistir configurações locais opcionalmente

O sistema deverá permitir persistência local opcional de configurações de uso, preferências de execução e histórico resumido de análises por meio de SQLite embutido. Essa persistência não deverá ser obrigatória para funcionamento básico [8][9].

### RF-23 — Operar sem integração externa obrigatória

O sistema deverá ser plenamente utilizável sem dependência de APIs externas. O núcleo da versão inicial deverá permanecer disponível localmente para parsing, estimativa e alertas baseados em regra.

### RF-24 — Preparar ponto de extensão para IA opcional

O sistema deverá prever um ponto de extensão funcional para geração futura de insights textuais por IA. Esse recurso, quando existir, deverá atuar sobre o resultado da análise e não substituir o motor de cálculo [3].

## Regras de negócio funcionais

### RN-01 — Engine define a semântica da estimativa

Toda análise deverá ser conduzida segundo a engine selecionada pelo usuário. O mesmo DDL textual poderá produzir resultados diferentes sob SQL Server e PostgreSQL devido às diferenças documentadas entre tipos, índices e procedimentos de estimativa [1][2][3].

### RN-02 — Cálculo não equivale a medição real

Os resultados apresentados pelo sistema deverão ser tratados como estimativas e não como medição exata de armazenamento físico em instância real. A interface deverá comunicar essa natureza de forma explícita [1][3].

### RN-03 — Índices representam benefício e custo

Sempre que índices forem considerados na análise, o sistema deverá comunicar que podem melhorar determinadas operações de leitura, busca ou ordenação, mas também trazem custo adicional de armazenamento e manutenção [3][11].

### RN-04 — Falhas parciais não invalidam todo o resultado

Se parte do DDL não puder ser interpretada, o sistema deverá tentar analisar o restante do conteúdo sempre que isso não comprometer a consistência global. As limitações deverão ser apresentadas ao usuário.

## Casos de uso principais

### UC-01 — Iniciar a ferramenta localmente

O usuário executa `tabyte serve`. O sistema inicializa o servidor local, informa a URL e abre a interface web no navegador padrão [4][5].

### UC-02 — Analisar DDL de uma tabela

O usuário cola o DDL de uma tabela, seleciona a engine e executa a análise. O sistema interpreta a estrutura, calcula tamanhos estimados e apresenta o detalhamento da tabela.

### UC-03 — Analisar múltiplas tabelas

O usuário informa um conjunto de `CREATE TABLE` e objetos relacionados. O sistema identifica as tabelas individualmente e mostra consolidado por tabela e por schema.

### UC-04 — Projetar volumetria

O usuário informa quantidade de linhas por tabela e ajusta crescimento esperado. O sistema recalcula os totais e apresenta nova estimativa global [10].

### UC-05 — Revisar alertas estruturais

O usuário navega pelos alertas de modelagem e sinais de impacto em performance para entender quais decisões do schema merecem revisão [3][7].

### UC-06 — Exportar resultado

Após a análise, o usuário exporta os dados em formato estruturado para uso em documentação técnica, avaliação arquitetural ou comparação entre cenários.

## Priorização funcional para a primeira entrega

Os requisitos mais críticos para a primeira entrega são RF-01 a RF-12, RF-14, RF-15, RF-18, RF-19 e RF-23. Esse subconjunto já permite entregar uma aplicação local utilizável, coerente com a proposta do produto e alinhada às diferenças reais entre SQL Server e PostgreSQL em termos de armazenamento, tipos e índices [1][2][3].

Os requisitos de persistência opcional e futura integração com IA podem entrar de forma incremental, desde que a primeira release já entregue cálculo explicável, operação local simples e alertas estruturalmente defensáveis [8][9].