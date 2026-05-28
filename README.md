        # postgresql — Уровни изоляции: Read Committed, Repeatable Read, Serializable

        Homework-шаблон для урока **l2_isolation_levels** (Уровни изоляции: Read Committed, Repeatable Read, Serializable) на платформе Vibe Learn.

        ## Что делать

        Дано: testcontainers PG + таблица accounts. Реализуй на Go:
1) Функцию Transfer(fromID, toID, amount) с правильной concurrency: либо Serializable
   с retry, либо SELECT FOR UPDATE с детерминированным порядком.
2) Тест-нагрузку: 100 goroutine параллельно делают переводы, в конце сумма всех
   счетов должна остаться неизменной.
3) Worker pool, читающий очередь pending_jobs с FOR UPDATE SKIP LOCKED.
4) Замеры throughput на разных стратегиях (Read Committed без локов vs FOR UPDATE vs
   Serializable).
Тесты проверят сохранение invariants под нагрузкой и корректность retry-логики.

## Контекст (из transfer-задачи урока)

Тебе предложили спроектировать concurrency для трёх operations в банковском приложении:

(Op1) **Перевод между счетами** одного клиента: списать с одного, начислить на другой,
      сумма не должна уйти в минус.
(Op2) **Отчёт по счёту за месяц**: SELECT всех транзакций за период + остаток на начало
      и конец. Должна быть консистентная картинка.
(Op3) **Распределение баллов лояльности** ночным batch-job: для каждого клиента, у
      которого «выполнено условие X», начислить баллы. Должно отработать idempotent.
(Op4) **Резервирование автомобиля в каршеринге**: клиент жмёт «забронировать» — нужно
      атомарно проверить «не забронирован ли уже» и зарезервировать.
(Op5) **Воркеры обрабатывают очередь заявок** из таблицы `pending_jobs` — N воркеров
      параллельно, каждый берёт следующую необработанную заявку.

## Recap из урока

- **Read Committed (дефолт PG)** — snapshot на КАЖДУЮ команду; non-repeatable read и lost update возможны.
- **Repeatable Read = snapshot isolation в PG**: snapshot фиксируется на всю транзакцию, phantom read невозможен (строже ANSI).
- **Serializable** убирает write skew через SSI — но даёт 40001 'could not serialize', обязателен retry-loop в коде.
- **SELECT FOR UPDATE** — pessimistic lock на строку; SKIP LOCKED идеален для job queue, NOWAIT — для fail-fast.
- **Без retry-логики Serializable превращается в random failures.** Реализуй retry с exponential backoff для serialization (40001) и deadlocks (40P01).

        ## Как работать

        1. Платформа Vibe Learn создаёт копию этого репо в твоём GitHub-аккаунте по клику «Начать домашку» на странице урока (через GitHub `/generate`, codecrafters-pattern).
        2. Склонируй копию локально, реализуй TODO в `main.go`, прогони тесты, запушь.
        3. CI (`.github/workflows/ci.yml`) запускает `go vet` + `go test ./...` на каждый push. Платформа слушает результат через webhook от GitHub Actions и обновляет статус домашки на странице урока.

        ## Локальное окружение

        - Go 1.22+
        - Docker + docker-compose — `docker compose up -d` поднимает single-node PostgreSQL 16 на `localhost:5432` с healthcheck. DSN: `postgres://postgres:postgres@localhost:5432/postgres`. Переопределяется через env `DATABASE_URL`.

        ## Запуск

        ```bash
        # Поднять локальный PostgreSQL
        docker compose up -d

        # Прогнать тесты (интеграционный включается через PG_INTEGRATION=1)
        go test ./...
        PG_INTEGRATION=1 go test ./...

        # Запустить main (печатает marker; замени stub на реализацию)
        go run .
        ```

        ## Заметка автора

        Это baseline-шаблон, сгенерированный платформой. Бизнес-сущность задачи (что конкретно реализовать в `main.go`, какие тесты сделать строгими) расширяется по ходу итераций — параллельно с углублением теории урока.
