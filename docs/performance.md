# Performance — fila de estudo

Medições das consultas usadas por `flashcardService.Study()` com `EXPLAIN (ANALYZE, BUFFERS)`.

## Ambiente

- PostgreSQL 15.19 (container `postgres:15`), Apple Silicon, cache aquecido (duas execuções descartadas antes da medida).
- **50.000 cartões em 250 decks** (200 por deck). Cerca de 30% são novos (`last_review_at IS NULL`). Os outros 70% foram revisados nos últimos 30 dias, com `next_review_at` entre -10 e +30 dias.
- Deck alvo: 66 cartões novos, 26 vencidos e 0 revisados hoje.
- **Antes:** schema na versão 6 (sem índices além da PK). **Depois:** versão 7 (`000007_indexes_and_constraints`), seguida de `ANALYZE flashcards`.

<details>
<summary>Seed usado</summary>

```sql
INSERT INTO decks (id, name, category, visibility, created_at)
SELECT uuid_generate_v4(), 'deck ' || g, 'perf', 'private', now() FROM generate_series(1, 250) g;

-- 200 cartões por deck: ~30% novos, ~70% já revisados nos últimos 30 dias
-- (uma parte hoje), com próxima revisão entre -10 e +30 dias.
INSERT INTO flashcards (id, deck_id, front, back, created_at, last_review_at, next_review_at, state, interval, ease_factor)
SELECT uuid_generate_v4(), c.deck_id, 'front ' || c.g, 'back ' || c.g, now() - interval '60 days',
       CASE WHEN c.r < 0.3 THEN NULL ELSE now() - (c.r2 * interval '30 days') END,
       CASE WHEN c.r < 0.3 THEN NULL ELSE now() + ((c.r3 * 40 - 10) * interval '1 day') END,
       CASE WHEN c.r < 0.3 THEN 'new' WHEN c.r < 0.45 THEN 'learning' ELSE 'review' END,
       0, 2.5
FROM (
    SELECT d.id AS deck_id, g, random() AS r, random() AS r2, random() AS r3
    FROM decks d CROSS JOIN generate_series(1, 200) g
) c;
ANALYZE decks; ANALYZE flashcards;
```

</details>

## Resumo

| Consulta | Antes | Depois | Buffers antes → depois |
|---|---|---|---|
| Revisados hoje (`deck_id = ? AND last_review_at > ?`) | 2.422 ms (Seq Scan) | 0.033 ms (Index) | 839 → 6 |
| Novos (`deck_id = ? AND last_review_at IS NULL LIMIT 20`) | 0.448 ms (Seq Scan) | 0.093 ms (Index) | 217 → 23 |
| Vencidos (`deck_id = ? AND next_review_at <= ?`) | 1.948 ms (Seq Scan) | 0.100 ms (Index) | 839 → 29 |

Antes da 000007, as três consultas faziam **Seq Scan na tabela inteira** (50 mil linhas e 839 páginas) para devolver dezenas de cartões de um único deck. Com os índices compostos `(deck_id, last_review_at)` e `(deck_id, next_review_at)`, o Postgres vai direto ao intervalo do deck. O custo passa a depender do tamanho do deck, não do total de cartões do sistema. A Q2 já era relativamente barata antes porque o `LIMIT 20` interrompe o scan cedo, mas mesmo assim lia 217 páginas.

## Planos — antes (v6)

```
--- Q1 revisados hoje
                                                         QUERY PLAN                                                         
----------------------------------------------------------------------------------------------------------------------------
 Seq Scan on flashcards  (cost=0.00..1839.00 rows=1 width=611) (actual time=2.400..2.401 rows=0 loops=1)
   Filter: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (last_review_at > date_trunc('day'::text, now())))
   Rows Removed by Filter: 50000
   Buffers: shared hit=839
 Planning:
   Buffers: shared hit=97
 Planning Time: 0.149 ms
 Execution Time: 2.422 ms
(8 rows)
--- Q2 novos (LIMIT 20)
                                                   QUERY PLAN                                                    
-----------------------------------------------------------------------------------------------------------------
 Limit  (cost=0.00..488.00 rows=20 width=611) (actual time=0.046..0.440 rows=20 loops=1)
   Buffers: shared hit=217
   ->  Seq Scan on flashcards  (cost=0.00..1464.00 rows=60 width=611) (actual time=0.045..0.438 rows=20 loops=1)
         Filter: ((last_review_at IS NULL) AND (deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid))
         Rows Removed by Filter: 12982
         Buffers: shared hit=217
 Planning Time: 0.043 ms
 Execution Time: 0.448 ms
(8 rows)
--- Q3 vencidos
                                                QUERY PLAN                                                 
-----------------------------------------------------------------------------------------------------------
 Seq Scan on flashcards  (cost=0.00..1714.00 rows=35 width=611) (actual time=0.057..1.943 rows=26 loops=1)
   Filter: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (next_review_at <= now()))
   Rows Removed by Filter: 49974
   Buffers: shared hit=839
 Planning:
   Buffers: shared hit=3
 Planning Time: 0.027 ms
 Execution Time: 1.948 ms
(8 rows)
```

## Planos — depois (v7)

```
--- Q1 revisados hoje
                                                                  QUERY PLAN                                                                  
----------------------------------------------------------------------------------------------------------------------------------------------
 Index Scan using idx_flashcards_deck_last_review on flashcards  (cost=0.42..8.44 rows=1 width=611) (actual time=0.020..0.020 rows=0 loops=1)
   Index Cond: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (last_review_at > date_trunc('day'::text, now())))
   Buffers: shared hit=6
 Planning:
   Buffers: shared hit=137
 Planning Time: 0.181 ms
 Execution Time: 0.033 ms
(7 rows)
--- Q2 novos (LIMIT 20)
                                                                   QUERY PLAN                                                                   
------------------------------------------------------------------------------------------------------------------------------------------------
 Limit  (cost=5.03..67.36 rows=20 width=611) (actual time=0.015..0.082 rows=20 loops=1)
   Buffers: shared hit=23
   ->  Bitmap Heap Scan on flashcards  (cost=5.03..192.03 rows=60 width=611) (actual time=0.014..0.080 rows=20 loops=1)
         Recheck Cond: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (last_review_at IS NULL))
         Heap Blocks: exact=20
         Buffers: shared hit=23
         ->  Bitmap Index Scan on idx_flashcards_deck_last_review  (cost=0.00..5.01 rows=60 width=0) (actual time=0.005..0.005 rows=66 loops=1)
               Index Cond: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (last_review_at IS NULL))
               Buffers: shared hit=3
 Planning Time: 0.025 ms
 Execution Time: 0.093 ms
(11 rows)
--- Q3 vencidos
                                                                QUERY PLAN                                                                
------------------------------------------------------------------------------------------------------------------------------------------
 Bitmap Heap Scan on flashcards  (cost=4.77..120.83 rows=34 width=611) (actual time=0.009..0.096 rows=26 loops=1)
   Recheck Cond: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (next_review_at <= now()))
   Heap Blocks: exact=26
   Buffers: shared hit=29
   ->  Bitmap Index Scan on idx_flashcards_deck_next_review  (cost=0.00..4.76 rows=34 width=0) (actual time=0.006..0.006 rows=26 loops=1)
         Index Cond: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (next_review_at <= now()))
         Buffers: shared hit=3
 Planning:
   Buffers: shared hit=3
 Planning Time: 0.018 ms
 Execution Time: 0.100 ms
(11 rows)
```

## Fase 5 — fila em uma única query (CTE)

`FindStudyQueue` substitui as três consultas acima por uma única query com CTEs (`reviewed_today`, `new_cards` e `due_cards`). Mesmo dataset e mesmo deck, schema v8, limite diário de 20.

| | 3 consultas (v7) | 1 query com CTE |
|---|---|---|
| Idas ao banco | 3 | 1 |
| Tempo de execução no servidor (soma) | 0.226 ms | 0.451 ms |
| Buffers | 58 | 113 |
| Índices usados | os dois | os dois (inclusive *Index Only Scan* na contagem) |

No servidor, a query única é um pouco mais cara, e o motivo é intencional: os cartões novos agora saem **ordenados por `created_at`** (ordem determinística). Para escolher os 20 mais antigos, o Postgres lê os 66 novos do deck. Antes, o `LIMIT 20` sem `ORDER BY` parava nos primeiros 20 que encontrasse, em ordem arbitrária. O ganho real está em outros dois pontos:

- **Uma ida ao banco em vez de três.** Com o banco em outra máquina, cada ida custa da ordem de 0,5 a 2 ms de rede, valor bem maior que o tempo de execução no servidor.
- **Snapshot consistente.** A contagem de revisados hoje e as listas de novos e vencidos vêm do mesmo instante. Com três consultas separadas, uma revisão concorrente podia cair entre elas.

A query também corrige um bug: com 20 ou mais revisões no dia, `dailyLimit - revisadosHoje` ficava `<= 0`, o repositório ignorava o `LIMIT` e devolvia **todos** os cartões novos. Agora o limite é `GREATEST(limite - revisados, 0)`.

```
                                                                                QUERY PLAN                                                                                 
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------
 Sort  (cost=320.77..320.87 rows=40 width=615) (actual time=0.374..0.375 rows=46 loops=1)
   Sort Key: (0), f.next_review_at, f.created_at, f.id
   Sort Method: quicksort  Memory: 31kB
   Buffers: shared hit=113
   ->  Append  (cost=198.27..319.71 rows=40 width=615) (actual time=0.266..0.352 rows=46 loops=1)
         Buffers: shared hit=104
         ->  Limit  (cost=198.27..198.28 rows=6 width=615) (actual time=0.266..0.267 rows=20 loops=1)
               Buffers: shared hit=75
               InitPlan 1 (returns $0)
                 ->  Aggregate  (cost=4.44..4.45 rows=1 width=8) (actual time=0.020..0.020 rows=1 loops=1)
                       Buffers: shared hit=6
                       ->  Index Only Scan using idx_flashcards_deck_last_review on flashcards  (cost=0.42..4.44 rows=1 width=0) (actual time=0.019..0.019 rows=0 loops=1)
                             Index Cond: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (last_review_at >= date_trunc('day'::text, now())))
                             Heap Fetches: 0
                             Buffers: shared hit=6
               ->  Sort  (cost=193.80..193.95 rows=60 width=615) (actual time=0.244..0.244 rows=20 loops=1)
                     Sort Key: f.created_at, f.id
                     Sort Method: top-N heapsort  Memory: 29kB
                     Buffers: shared hit=69
                     ->  Bitmap Heap Scan on flashcards f  (cost=5.03..192.03 rows=60 width=615) (actual time=0.017..0.226 rows=66 loops=1)
                           Recheck Cond: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (last_review_at IS NULL))
                           Heap Blocks: exact=66
                           Buffers: shared hit=69
                           ->  Bitmap Index Scan on idx_flashcards_deck_last_review  (cost=0.00..5.01 rows=60 width=0) (actual time=0.008..0.008 rows=66 loops=1)
                                 Index Cond: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (last_review_at IS NULL))
                                 Buffers: shared hit=3
         ->  Bitmap Heap Scan on flashcards f_1  (cost=4.77..120.83 rows=34 width=615) (actual time=0.011..0.083 rows=26 loops=1)
               Recheck Cond: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (next_review_at <= now()))
               Heap Blocks: exact=26
               Buffers: shared hit=29
               ->  Bitmap Index Scan on idx_flashcards_deck_next_review  (cost=0.00..4.76 rows=34 width=0) (actual time=0.008..0.008 rows=26 loops=1)
                     Index Cond: ((deck_id = '4df442a9-0417-48f8-b166-4f8182612bb4'::uuid) AND (next_review_at <= now()))
                     Buffers: shared hit=3
 Planning:
   Buffers: shared hit=161
 Planning Time: 0.323 ms
 Execution Time: 0.451 ms
(37 rows)
```
