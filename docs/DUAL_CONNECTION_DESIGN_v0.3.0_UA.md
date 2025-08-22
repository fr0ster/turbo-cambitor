# Камбітор: Подвійне з'єднання (Основне + Сервісне) — Дизайн і Поведінка

Версія: v0.3.0 (2025-08-22)

Цей документ описує нову архітектуру камбітора з двома WebSocket-з'єднаннями: основним (для бізнес-функцій) і сервісним (для моніторингу здоров'я та керування життєвим циклом основного). Також наведено режими роботи (асинхронний і синхронний), API та зміни до тестів.

## Цілі
- Надійно обробляти падіння сервера і нестабільність мережі.
- Вчасно «ставити на паузу» основне з'єднання при відсутності Pong, і «знімати з паузи» при відновленні.
- Дати клієнту вибір між fail-fast і блокуючою (очікуючою) моделлю.
- Прибрати залежність тестів від штучних очікувань на рестарт сервера (waitServerListening) і перейти на здоров'я через Ping/Pong.

## Архітектура
- Сервісне з'єднання (healthConn):
  - Окремий постійний WebSocket до того ж endpoint.
  - Луп моніторингу: з інтервалом `healthInterval` шле Ping і чекає Pong до `healthTimeout`.
  - Debounce: `failThreshold` підряд → стан Unhealthy; `successThreshold` → назад у Healthy.
  - Самостійно перепідключається з backoff (мін/макс + джиттер), зберігаючи перевірки.
  - Публічне API:
    - `HealthStatus() enum {Connecting, Healthy, Unhealthy}`
    - `HealthPing(timeout) bool` — одноразовий пінг на вимогу
    - `WaitHealthy(ctx) error` — очікувати стан Healthy або таймаут
    - `EnableHealthMonitor(opts) / DisableHealthMonitor()`

- Основне з'єднання (mainConn):
  - Виконує бізнес-функції: `Call`, `Subscribe`, `Unsubscribe`, `ListOfSubscriptions`.
  - При Unhealthy: переходить у «paused» (рекомендовано — закрити сокет, щоб уникати напів-станів).
  - При Healthy: забезпечує підняття (Connect → WaitStarted → WaitReady), перевстановлює підписки/хендлери/таймаути.
  - Авто-reconnect підпорядковується монітору здоров'я (або вимикається).

## Режими роботи
- Асинхронний (fail-fast):
  - Публічні методи одразу перевіряють здоров'я/паузу; якщо не Healthy — повертають `ErrNotConnected`.
  - Клієнт сам викликає `WaitHealthy`/`HealthPing`/`HealthStatus` і ретраїть запити.

- Синхронний (wait):
  - Публічні методи перед I/O викликають `WaitHealthy(ctx)` з пер-кол таймаутом; блокуються до Healthy або таймауту.
  - Рекомендовано «один потік — один камбітор» для простої синхронізації.

## Публічне API (доповнення)
- `SetModeAsync()` / `SetModeSync(defaultTimeout)`
- `EnableHealthMonitor(opts{ healthInterval, healthTimeout, failThreshold, successThreshold, backoff })`
- `DisableHealthMonitor()`
- `HealthStatus() enum`
- `HealthPing(timeout time.Duration) bool`
- `WaitHealthy(ctx context.Context) error`
- `IsPaused() bool`
- Зворотна сумісність: `WaitForPong/WaitReady` переводяться на сервісне з'єднання (без gorilla у проді).

## Стан-машини
- Health: `Connecting → Healthy → Unhealthy → (reconnect) → Connecting → Healthy …`
- Main: `Active → Paused → Reconnecting → Active …`
- Тригери:
  - N фейлів Ping → Unhealthy → `pauseMain()` (close/stop лупи).
  - M успішних Pong → Healthy → `resumeMain()` (Connect → WaitStarted → WaitReady → reapply).

## Конкурентність
- Станові прапори — atomic; переходи — під невеликим м'ютексом/контрольною горутиною.
- Усі control/data-записи синхронізовані (writeMu у restler вже використовується).
- Життєвий цикл mainConn — під єдиним lifecycle-м'ютексом, щоб уникнути паралельних reconnect.

## Налаштування за замовчуванням
- `healthInterval`: 1s (тести), 5–10s (прод)
- `healthTimeout`: 500–1000ms (тести), 2–3s (прод)
- `failThreshold`: 2–3; `successThreshold`: 1–2
- Backoff сервісного reconnect: 200ms → 2s + джиттер
- Sync per-call timeout: 5s у тестах; конфігурується

## Зміни у тестах
- Видалити `waitServerListening` і локальний `waitForPong`.
- Авто-reconnect-сценарії:
  - Async: після «падіння з рестартом» — `Call` повертає `ErrNotConnected`, доки `HealthStatus != Healthy`; після `WaitHealthy` — `Call` успішний.
  - Sync: `Call` блокується до Healthy або таймауту; перевірити розблокування після рестарту.
- Readiness перевіряти через `HealthPing`/`WaitHealthy` (кабіторне API).

## Краєві випадки
- Якщо біржа обмежує кількість з'єднань — сервісний канал робимо опційним; документуємо ризики.
- Перевстановлення підписок після resume — автоматично (як зараз).
- Анти-флуд: обмежити частоту Ping.

---

# Дорожня карта реалізації

1. Типи і конфіги в `streamer.go`
   - Режими (Async/Sync), опції health, атоміки, м'ютекси, канали стану.
   - API: `HealthStatus`, `HealthPing`, `WaitHealthy`, `SetMode*`, `Enable/DisableHealthMonitor`.

2. Сервісний сокет
   - Створення/підтримка окремого WS.
   - Ping-луп з порогами (fail/success thresholds).
   - Reconnect із backoff + джиттер.
   - Публікація стану та керування переходами.

3. Інтеграція з основним
   - `pauseMain()` (close/stop), `resumeMain()` (Connect → WaitStarted → WaitReady → reapply конфіги/хендлери/підписки).
   - Перевірка стану у публічних методах (відмова або очікування залежно від режиму).

4. Переробка AutoReconnect
   - Підпорядкувати монітору здоров'я або вимкнути при активному health monitor.

5. Тести
   - Оновити сценарії (Async/Sync), прибрати `waitServerListening`/локальні Pong-хелпери.
   - Додати перевірки `HealthPing/WaitHealthy/HealthStatus`.

6. Документація
   - Опис режимів, конфігів, ноти міграції.

---

## Переваги та недоліки

[Аналог секції Pros/Cons/Security з англомовної версії]
