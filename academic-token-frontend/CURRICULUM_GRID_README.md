# Visualização de Grade Curricular

Esta funcionalidade permite visualizar a grade curricular organizada por períodos/semestres, com visualização de pré-requisitos e status de conclusão das disciplinas.

## Como Usar

### 1. Instalar Dependências (Opcional - para React Query)

Se quiser usar o React Query para melhor gerenciamento de estado:

```bash
npm install @tanstack/react-query @tanstack/react-query-devtools
```

Depois, descomente as linhas do QueryProvider em `app/layout.tsx`.

### 2. Acessar a Visualização

1. Navegue para a página de Curriculum: `/curriculum`
2. Clique no botão "🎯 Grid View (New)" 
3. Ou acesse diretamente: `/curriculum/view/[courseId]`

### 3. Estrutura de Dados

A visualização busca dados de três endpoints principais:

- `/curriculum/tree/{courseId}` - Estrutura do currículo com semestres
- `/subjects/course/{courseId}` - Detalhes das disciplinas
- `/prerequisites/course/{courseId}` - Pré-requisitos (ainda não implementado no backend)

### 4. Funcionalidades

- **Visualização por Período**: Disciplinas organizadas em colunas por semestre
- **Status Visual**: 
  - Verde: Disciplina concluída
  - Cinza com cadeado: Bloqueada por pré-requisitos
  - Branco: Disponível para matrícula
- **Detalhes ao Clicar**: Modal com informações completas da disciplina
- **Pré-requisitos**: Lista visual dos pré-requisitos de cada disciplina

### 5. Customização

Para adicionar disciplinas concluídas pelo aluno, edite o array `completedSubjects` em:
`app/curriculum/view/[courseId]/page.tsx`

### 6. Visualização Alternativa - Fluxo de Pré-requisitos

Se quiser usar a visualização D3.js (opcional):

```bash
npm install d3 @types/d3
```

Depois importe e use o componente `PrerequisiteFlow` na página.

## Estrutura de Arquivos Criados

```
app/
├── types/
│   └── curriculum.types.ts         # Tipos TypeScript
├── hooks/
│   └── useCurriculumData.ts       # Hook para buscar dados
├── components/
│   ├── CurriculumGrid/
│   │   ├── CurriculumGrid.tsx     # Grade principal
│   │   ├── SubjectCard.tsx        # Card de disciplina
│   │   └── CurriculumGrid.css     # Estilos
│   ├── SubjectDetailsModal/
│   │   ├── SubjectDetailsModal.tsx # Modal de detalhes
│   │   └── SubjectDetailsModal.css # Estilos do modal
│   └── PrerequisiteFlow/          # (Opcional)
│       └── PrerequisiteFlow.tsx    # Visualização D3
└── curriculum/
    └── view/
        └── [courseId]/
            ├── page.tsx            # Página principal
            └── page.css            # Estilos da página
```

## Próximos Passos

1. Implementar endpoint de pré-requisitos no backend
2. Integrar com dados reais do estudante logado
3. Adicionar funcionalidade de matrícula ao clicar em disciplina disponível
4. Implementar filtros e busca na grade
