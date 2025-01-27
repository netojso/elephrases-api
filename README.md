# Regras Funcionais e Não Funcionais para um Aplicativo de Flashcards

## Regras Funcionais

### 1. Gerenciamento de Flashcards
- Permitir a criação de flashcards com:
  - Frente: Texto, imagem, ou áudio para perguntas ou termos.
  - Verso: Texto, imagem, ou áudio para respostas ou definições.
- Permitir a edição e exclusão de flashcards existentes.
- Suporte para organização dos flashcards em decks temáticos.
- Funcionalidade de busca para encontrar flashcards ou decks específicos.

### 2. Modos de Estudo
#### 2.1. Modo Aprendizado
- Exibir flashcards sequencialmente para facilitar a memorização inicial.
- Permitir que o usuário marque cada flashcard como:
  - "Fácil": Indica que o usuário já domina o conteúdo.
  - "Médio": Indica que o conteúdo requer revisão moderada.
  - "Difícil": Indica que o conteúdo precisa de revisão frequente.
- Mostrar um resumo ao final com a classificação dos flashcards.

#### 2.2. Modo Repetição Espaciada
- Usar um algoritmo baseado em intervalos otimizados para reforço de memorização.
- Notificar o usuário quando for hora de revisar os flashcards.
- Ajustar automaticamente o intervalo com base no feedback do usuário.

#### 2.3. Modo Quiz
- Apresentar a frente do flashcard e solicitar que o usuário insira ou diga a resposta antes de revelar o verso.
- Avaliar automaticamente a resposta, se aplicável, ou permitir avaliação manual.
- Exibir um placar ao final com a porcentagem de acertos e erros.

#### 2.4. Modo Aleatório
- Apresentar os flashcards em ordem aleatória para testar a memorização.
- Incluir uma opção para revisar apenas os flashcards marcados como "difícil" ou "não conhecido".

### 3. Gerenciamento de Decks
- Criar, editar, excluir e duplicar decks.
- Importar e exportar decks em formatos como CSV ou JSON.
- Compartilhar decks com outros usuários por meio de links ou publicação na comunidade do aplicativo.
- Suporte para categorização de decks (ex.: idiomas, ciências, hobbies).

### 4. Avaliação e Feedback
- Permitir ao usuário marcar flashcards como:
  - "Conhecido": Revisão em intervalos maiores.
  - "Não conhecido": Revisão em intervalos menores.
- Exibir relatórios de desempenho:
  - Percentual de acertos e erros.
  - Tempo gasto estudando.
  - Estatísticas dos flashcards mais difíceis.

### 5. Personalização
- Opção de personalizar:
  - Temas e cores do aplicativo.
  - Ordem de exibição dos flashcards (ex.: aleatória ou ordenada).
  - Notificações de revisão.

### 6. Funcionalidades Comunitárias
- Criar uma seção de decks compartilhados, com:
  - Sistema de avaliação (ex.: estrelas ou curtidas).
  - Opção de seguir criadores de decks.
- Permitir comentários nos decks compartilhados.

### 7. Acessibilidade
- Suporte para idiomas variados.
- Compatibilidade com leitores de tela e ferramentas de acessibilidade.

### 8. Recursos Avançados
- Sugestões automáticas de criação de conteúdo com base em IA.
- Gamificação:
  - Pontos e conquistas por uso contínuo.
  - Níveis e rankings.
- Modo offline para estudo sem conexão à internet.

## Regras Não Funcionais

### 1. Desempenho
- O carregamento dos flashcards deve ocorrer em menos de 2 segundos.
- O sistema deve suportar decks com até 1.000 flashcards sem queda de desempenho.

### 2. Segurança
- Garantir que os dados dos usuários sejam armazenados de forma segura e criptografada.
- Implementar autenticação segura com suporte a autenticação de dois fatores (2FA).

### 3. Escalabilidade
- Suporte a até 1 milhão de usuários ativos simultaneamente sem degradação da performance.
- Arquitetura preparada para expansão de servidores e armazenamento.

### 4. Confiabilidade
- O aplicativo deve ter uptime de 99,9%.
- Backup automático dos dados dos usuários a cada 24 horas.

### 5. Usabilidade
- Interface intuitiva e fácil de usar para todas as idades.
- Suporte a dispositivos móveis e tablets, com layout responsivo.

### 6. Portabilidade
- Disponível para plataformas Android, iOS e Web.
- Sincronização de dados entre dispositivos em tempo real.

### 7. Manutenibilidade
- Código estruturado para fácil manutenção e evolução.
- Registro detalhado de logs de erro e eventos para diagnóstico.

### 8. Sustentabilidade
- Consumo eficiente de recursos para preservar a bateria de dispositivos móveis.
- Redução de uso de dados móveis para sincronização.

---

## Entidades da Aplicação

### 1. Usuário
- **Atributos:**
  - ID
  - Nome
  - E-mail
  - Senha (criptografada)
  - Data de criação
  - Configurações de personalização (tema, notificações, idioma)
- **Ações:**
  - Criar conta
  - Editar perfil
  - Alterar senha
  - Excluir conta

### 2. Flashcard
- **Atributos:**
  - ID
  - Frente (texto, imagem ou áudio)
  - Verso (texto, imagem ou áudio)
  - Data de criação
  - Última revisão
  - Status (fácil, médio, difícil)
- **Ações:**
  - Criar flashcard
  - Editar flashcard
  - Excluir flashcard

### 3. Deck
- **Atributos:**
  - ID
  - Nome
  - Descrição
  - Categoria
  - Lista de flashcards
  - Data de criação
  - Visibilidade (privado/público)
- **Ações:**
  - Criar deck
  - Editar deck
  - Compartilhar deck
  - Excluir deck

### 4. Comunidade
- **Atributos:**
  - ID do deck
  - Usuário criador
  - Avaliações (estrelas ou curtidas)
  - Comentários
  - Número de downloads
- **Ações:**
  - Publicar deck
  - Avaliar deck
  - Comentar em decks

### 5. Relatório
- **Atributos:**
  - ID do usuário
  - Estatísticas de desempenho (acertos, erros, tempo de estudo)
  - Flashcards mais difíceis
  - Histórico de progresso
- **Ações:**
  - Gerar relatório
  - Exibir resumo de desempenho

---
Essas entidades definem a estrutura básica da aplicação e as interações principais. Ajustes podem ser feitos conforme necessário.
