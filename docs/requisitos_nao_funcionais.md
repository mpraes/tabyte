# Tabyte — Requisitos Não Funcionais v0.3

## Objetivo

Este documento define **apenas** os requisitos não funcionais do **Tabyte**, isto é, os atributos de qualidade e restrições técnicas do produto, sem descrever capacidades funcionais do sistema. A organização adotada segue a lógica de atributos de qualidade de software amplamente usada em modelos como ISO/IEC 25010, que estrutura requisitos de qualidade em características como desempenho, compatibilidade, usabilidade, confiabilidade, segurança, manutenibilidade e portabilidade [1][2][3].

O contexto assumido é o de uma aplicação **local-first**, distribuída como ferramenta **CLI cross-platform** com interface acessada via navegador em `localhost`, implementada em **Go**, com assets embutidos no binário e persistência opcional por **SQLite** para uso interno [4][5][6][7].

## Escopo

Este documento não define o que o Tabyte faz do ponto de vista de negócio, mas como ele deve se comportar em termos de qualidade, restrições operacionais, mantenibilidade e implantação. Assim, itens como “analisar DDL”, “estimar volumetria” ou “expor endpoint específico” não pertencem a este documento, exceto quando impactam diretamente um atributo de qualidade [1][8].

## Estrutura de qualidade

Para manter consistência, os requisitos foram agrupados em categorias de qualidade inspiradas em ISO/IEC 25010. Esse padrão é útil porque oferece terminologia consistente para especificar e avaliar requisitos não funcionais de produto de software [1][2].

## Desempenho e eficiência

### RNF-01 — Tempo de inicialização

A aplicação deve iniciar em tempo curto em ambiente local comum, incluindo carregamento de configuração, inicialização do servidor local e disponibilização da interface web. O tempo de boot deve ser suficientemente baixo para preservar a percepção de ferramenta utilitária de uso imediato [2][9].

### RNF-02 — Tempo de resposta interativo

As operações interativas da aplicação devem responder de forma fluida para schemas pequenos e médios. A experiência do usuário não deve transmitir sensação de travamento durante operações normais de uso local [2][9].

### RNF-03 — Eficiência de recursos

A aplicação deve operar com consumo moderado de CPU e memória, de modo compatível com estações de trabalho comuns e notebooks de desenvolvimento. O produto não deve pressupor infraestrutura robusta para uso local [2][9].

### RNF-04 — Baixa dependência operacional

A execução do sistema deve exigir o mínimo possível de dependências externas no ambiente do usuário final. Sempre que viável, o produto deve ser entregue de forma autossuficiente, reduzindo etapas manuais de preparação [4][6].

## Compatibilidade e portabilidade

### RNF-05 — Paridade cross-platform

O produto deve manter comportamento operacional equivalente entre Windows e Ubuntu/Linux nos fluxos principais de execução local. Diferenças específicas de sistema operacional devem ser minimizadas e, quando inevitáveis, claramente documentadas [2][9].

### RNF-06 — Instalação simples

O processo de instalação ou disponibilização do binário deve ser simples o suficiente para um usuário técnico comum, sem exigir configuração extensa do ambiente. A experiência de entrada no produto deve ser curta, previsível e bem documentada [2][9].

### RNF-07 — Portabilidade de build

O projeto deve poder ser compilado e empacotado para Windows e Linux sem pipeline excessivamente complexa. A stack técnica deve favorecer portabilidade prática do artefato final [2][7].

### RNF-08 — Ausência de banco servidor externo

O sistema não deve depender de banco de dados servidor externo para sua operação normal. Quando houver persistência local, ela deve ser embutida e transparente para o usuário, preferencialmente em formato de arquivo de aplicação [5].

## Usabilidade técnica

### RNF-09 — Operabilidade clara

A aplicação deve ser fácil de operar por um usuário técnico, com fluxo de inicialização, uso e encerramento compreensível. Mensagens de estado, erro e sucesso devem ser curtas, objetivas e úteis [2][9].

### RNF-10 — Consistência da interface

A interface web local deve manter consistência visual, terminológica e estrutural entre telas, painéis e mensagens. O sistema não deve apresentar variações arbitrárias de nomenclatura ou comportamento em interações semelhantes [1][2].

### RNF-11 — Prevenção de erro de uso

A experiência deve reduzir erros evitáveis de operação por meio de validações, mensagens claras e feedback imediato quando uma entrada estiver incompleta, inválida ou incompatível com o contexto atual [2][9].

### RNF-12 — Acessibilidade básica

A interface web deve atender a princípios básicos de acessibilidade, incluindo navegação por teclado, foco visível, contraste suficiente e marcação semântica adequada. Mesmo sendo uma ferramenta técnica, o produto não deve ignorar acessibilidade fundamental [1][3].

## Confiabilidade e resiliência

### RNF-13 — Comportamento determinístico

Sob as mesmas condições de entrada, configuração e versão, o sistema deve apresentar comportamento consistente e resultados reproduzíveis. Esse requisito é essencial para confiança do usuário e repetibilidade de uso [1][9].

### RNF-14 — Tolerância a falhas parciais

A aplicação deve degradar de forma controlada diante de erros localizados, evitando interromper desnecessariamente toda a execução quando uma falha parcial puder ser isolada. Em caso de falha, o sistema deve privilegiar mensagem clara em vez de encerramento abrupto [2][9].

### RNF-15 — Recuperabilidade operacional

Após erro de execução, encerramento inesperado ou falha local não crítica, a aplicação deve poder ser reiniciada com baixo esforço operacional. Quando houver persistência local, ela não deve se tornar ponto frequente de corrupção ou recuperação manual [2][5].

### RNF-16 — Observabilidade local

A aplicação deve gerar logs suficientes para diagnóstico local de problemas, sem depender de infraestrutura externa de observabilidade. Os registros devem ser compreensíveis para manutenção e depuração em ambiente local [2][8].

## Segurança e privacidade

### RNF-17 — Bind local restrito

Por padrão, a aplicação deve aceitar conexões apenas em `localhost` ou `127.0.0.1`, evitando exposição indevida a interfaces de rede externas. Qualquer flexibilização desse comportamento deve ser explícita e não padrão [6].

### RNF-18 — Privacidade do conteúdo analisado

O conteúdo informado pelo usuário deve permanecer local por padrão. Nenhum dado submetido à aplicação deve ser transmitido a serviços externos sem ação explícita ou configuração consciente do usuário [5].

### RNF-19 — Persistência mínima necessária

Quando a persistência local estiver ativa, o sistema deve armazenar apenas o necessário para a experiência proposta, evitando retenção excessiva de dados sem valor claro. A política de persistência deve favorecer simplicidade, previsibilidade e controle local [5].

### RNF-20 — Transparência de integração externa

Se houver integração futura com serviços externos, como provedores de IA, o sistema deve comunicar de forma explícita a existência dessa interação e seus efeitos sobre privacidade e tráfego de dados [1][8].

## Manutenibilidade

### RNF-21 — Modularidade interna

A solução deve ser organizada em módulos internos com responsabilidades claras, reduzindo acoplamento entre interface, runtime CLI, domínio, persistência e integrações. A modularidade é requisito central para evolução segura do produto [2][9].

### RNF-22 — Testabilidade

A arquitetura deve permitir testes automatizados em múltiplos níveis, incluindo unidade, integração e smoke tests do runtime local. Componentes centrais não devem depender de acoplamentos que dificultem verificação automatizada [2][8].

### RNF-23 — Modificabilidade

O sistema deve permitir evolução incremental com baixo custo de mudança em módulos isolados. Alterações em persistência, UI ou integrações futuras não devem exigir reestruturação do núcleo inteiro [2][9].

### RNF-24 — Analisabilidade

A base de código deve favorecer entendimento técnico por contribuidores, com estrutura previsível, nomes consistentes e separação clara entre decisões de domínio e detalhes de infraestrutura [2][8].

### RNF-25 — Rastreabilidade técnica

Decisões de implementação relevantes, especialmente sobre runtime local, empacotamento, persistência e atributos de qualidade, devem poder ser rastreadas em documentação do projeto. Isso reduz ambiguidade e melhora governança técnica do repositório [8].

## Restrições tecnológicas

### RNF-26 — Linguagem principal

A implementação principal do produto deve usar **Go**, visando simplicidade de distribuição, portabilidade prática e servidor HTTP local nativo [6].

### RNF-27 — Assets embutidos no binário

A interface web estática deve ser empacotada preferencialmente dentro do binário usando o mecanismo `embed` da linguagem Go, reduzindo dependências externas de distribuição [4].

### RNF-28 — Persistência opcional via SQLite

Quando persistência local for necessária, o mecanismo preferencial deve ser **SQLite**, por sua adequação como formato de arquivo de aplicação local e simplicidade operacional [5].

### RNF-29 — Preferência por stack sem CGO quando viável

A solução deve preferir componentes que reduzam dependências de build e runtime complexas, especialmente em cenários de distribuição cross-platform. Quando possível, a stack deve evitar exigências de CGO para simplificar empacotamento [7].

## Implantação e operação local

### RNF-30 — Runtime único por CLI

O modelo operacional principal do produto deve ser a execução por CLI, com inicialização do servidor local e acesso à interface pelo navegador. Esse padrão deve ser o mesmo em Windows e Linux para reduzir variação operacional [6][2].

### RNF-31 — Ausência de serviços auxiliares obrigatórios

A instalação padrão não deve exigir processos complementares, agentes residentes ou infraestrutura local paralela para que o produto funcione [4][5].

### RNF-32 — Encerramento limpo

A aplicação deve poder ser encerrada sem deixar processos órfãos, portas presas ou arquivos temporários desnecessários. O ciclo de vida local do processo deve ser simples e previsível [2][9].

## Priorização

Para a primeira entrega, os requisitos mais críticos são RNF-04, RNF-05, RNF-06, RNF-08, RNF-09, RNF-13, RNF-17, RNF-21, RNF-22, RNF-26, RNF-27, RNF-28 e RNF-30. Esse conjunto protege o núcleo de qualidade esperado: simplicidade operacional, execução local previsível, modularidade e portabilidade real do produto [1][2][5].

## Observação final

A principal correção aplicada nesta revisão foi remover formulações que descreviam comportamento funcional do sistema e manter apenas atributos de qualidade, restrições tecnológicas e expectativas de operação do produto. Isso torna o documento mais aderente ao conceito de requisito não funcional em engenharia de software [1][3][8].