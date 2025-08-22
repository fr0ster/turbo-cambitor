# Дорожня карта: Подвійне з'єднання (Основне + Сервісне)

Цей документ деталізує етапи впровадження, API, задачі, критерії приймання, тест-план і ризики для архітектури з двома WebSocket-з'єднаннями в камбіторі.

## 0) Обсяг змін (Scope)
- Додати сервісне з'єднання (healthConn) для моніторингу стану через Ping/Pong.
- Керувати основним з'єднанням (mainConn) через стани: пауза/відновлення.
- Додати режими роботи: Async (fail-fast) і Sync (wait).
- Оновити публічний API та тести, прибравши залежність від `waitServerListening` і локальних Pong-хелперів.

## 1) Публічні API (сигнатури)
- Конфіг/режими:
  - `func (sw *StreamWrapper) SetModeAsync() *StreamWrapper`
  - `func (sw *StreamWrapper) SetModeSync(defaultTimeout time.Duration) *StreamWrapper`
  - `type HealthOptions struct { Interval, Timeout, BackoffMin, BackoffMax time.Duration; FailThreshold, SuccessThreshold int; Jitter float64 }`
  - `func (sw *StreamWrapper) EnableHealthMonitor(opts HealthOptions) *StreamWrapper`
  - `func (sw *StreamWrapper) DisableHealthMonitor() *StreamWrapper`
- Перевірка стану:
  - `type HealthState int` + `const (HealthConnecting HealthState = iota; HealthHealthy; HealthUnhealthy)`
  - `func (sw *StreamWrapper) HealthStatus() HealthState`
  - `func (sw *StreamWrapper) HealthPing(timeout time.Duration) bool`
  - `func (sw *StreamWrapper) WaitHealthy(ctx context.Context) error`
  - Back-compat: `WaitForPong(timeout)`/`WaitReady(ctx, interval)` маршрутизуються через сервісний канал.

### Тести на рівні функцій (покрокові сценарії)

Нижче — конкретні сценарії тестування для кожної нової/оновленої функції API.

1. `SetModeAsync()`
  - Передумова: HealthMonitor увімкнено, сервер у стані Unhealthy (імітувати пропуск Pong N разів через мок).
  - Дія: викликати `sw.SetModeAsync()`, далі `sw.Call(LIST_SUBSCRIPTIONS)`.
  - Очікування: негайна помилка `ErrNotConnected`; після відновлення здоров'я (мок починає відповідати Pong) — повторний `Call` успішний.

2. `SetModeSync(defaultTimeout)`
  - Передумова: Unhealthy.
  - Дія: `sw.SetModeSync(2*s)`, викликати `sw.Call(...)` у окремій горутині; через 1s відновити здоров'я (мок повертає Pong/рестарт).
  - Очікування: виклик блокується < 2s і завершується успішно; якщо не відновлювати — по закінченню таймауту помилка `context deadline exceeded`.

3. `EnableHealthMonitor(opts)`
  - Передумова: вимкнений монітор; встановити `opts.Interval=200ms`, `Timeout=150ms`, `FailThreshold=2`, `SuccessThreshold=1`.
  - Дія: `EnableHealthMonitor(opts)`; змусити мок 2 рази не відповідати на Ping (імітувати пропуски Pong), потім відповісти.
  - Очікування: стани послідовно `Connecting→Healthy→Unhealthy` (після двох пропусків) → `Healthy` (після успішного Pong); перевірити через `HealthStatus()` і часові межі.

4. `DisableHealthMonitor()`
  - Передумова: Monitor увімкнений і шле Ping.
  - Дія: `DisableHealthMonitor()`; захопити часове вікно 2*Interval.
  - Очікування: нових Ping немає (перевіряється опосередковано — стан не змінюється на Unhealthy при симуляції пропусків; також у мок-сервері можна підрахувати кількість отриманих Ping).

5. `HealthStatus()`
  - Передумова: тільки но піднято healthConn.
  - Дія/Очікування: спершу `Connecting`; після першого Pong — `Healthy`; при послідовних провалах — `Unhealthy`.

6. `HealthPing(timeout)`
  - Сценарій 1: сервер живий — повертає `true` < timeout.
  - Сценарій 2: сервер не відповідає — повертає `false` ≈ timeout.
  - Сценарій 3: під час реконекту — залежно від мок-налаштувань, переконатися що не панікує та повертає коректне булеве значення.

7. `WaitHealthy(ctx)`
  - Сценарій 1: server Healthy — повертає `nil` швидко (< 100ms при малих інтервалах).
  - Сценарій 2: server Unhealthy — з `ctx` з таймаутом 500ms повертає `ctx.Err()` після таймауту.
  - Сценарій 3: Unhealthy → Healthy протягом вікна — повертає `nil` до дедлайну.

8. `IsPaused()` (непрямий тест pause/resume)
  - Передумова: активні підписки на mainConn.
  - Дія: перевести health у `Unhealthy` (імітувати пропуски Pong ≥ FailThreshold).
  - Очікування: `IsPaused()==true`; mainConn закритий; `Subscribe/Call` у Async режимі дає `ErrNotConnected`, у Sync — блокується до `Healthy`.
  - Після відновлення: `IsPaused()==false`; автоматичне відновлення підписок; отримання повідомлень від мок-а.

9. Back-compat: `WaitForPong(timeout)` / `WaitReady(ctx, interval)`
  - Перевірити, що реалізація використовує сервісне з'єднання і поводиться ідентично: повертає `true` коли мок відповідає Pong, `false` при таймауті.

10. Координація з `EnableAutoReconnect()`
  - Передумова: HealthMonitor увімкнено; автоперепідключення також увімкнено.
  - Дія: викликати "падіння+рестарт" у мок-і.
  - Очікування: відсутність подвійних реконектів/ресабскрайбів; один цикл pause→resume; отримання даних після resume.

## 2) Внутрішні структури/поля
- У `StreamWrapper` додати:
  - `healthConn web_socket.WebSocketCommonInterface`
  - `healthState atomic.Value` (HealthState)
  - `paused atomic.Bool` (керує доступом до mainConn)
  - `mode enum` (Async/Sync) + `defaultSyncTimeout time.Duration`
  - `lifecycleMu sync.Mutex` (життєвий цикл main)
  - `healthMu sync.Mutex` (життєвий цикл health)
  - контексти/скасування для лупів, канали подій переходів
- Реюз існуючих: реєстрація хендлерів, timeouts, message logger, ресабскрайб підписок.

## 3) Етапи впровадження

### Етап A: Каркас health API (без лупа)
- Додати типи, поля, методи `HealthStatus`, `HealthPing`, `WaitHealthy` (через тимчасовий Ping контролем).
- `EnableHealthMonitor` поки лише зберігає опції; без фонового монітору.
- Тести: перевести на `HealthPing/WaitHealthy` без видалення старих хелперів (позначити deprecated).

### Етап B: Сервісний сокет і монітор
- Створити/підтримувати окремий `healthConn` через factory.
- Запустити ping-луп з порогами (`FailThreshold/SuccessThreshold`).
- Реалізувати reconnect з backoff + джиттер.
- Публікація переходів стану.
- Метрики/логи: RTT, кількість фейлів/відновлень, рівень backoff.

### Етап C: Pause/Resume основного
- `pauseMain()`: атомарно виставити `paused=true`, коректно закрити mainConn, зупинити лупи.
- `resumeMain()`: підняти mainConn (Open → WaitStarted → WaitReady), перевстановити таймаути/хендлери/логгер, ресабскрайб підписки.
- Серіалізація через `lifecycleMu`, ідемпотентність.

### Етап D: Режими Async/Sync
- У публічних методах (Call/Subscribe/Unsubscribe/ListOfSubscriptions):
  - Async: якщо `paused` або `HealthStatus!=Healthy` → `ErrNotConnected`.
  - Sync: `WaitHealthy(ctx)` з пер-кол таймаутом перед I/O.
- Налаштування таймаутів і дефолтів.

### Етап E: Очищення тестів і сумісність
- Закоментувати `waitServerListening` і локальний `waitForPong` у `streamer_test.go`.
- Переписати сценарії reconnect на використання `HealthPing/WaitHealthy/HealthStatus`.
- Додати тести на флапінг (чередування Healthy/Unhealthy), таймаути в Sync-режимі, деградацію healthConn.

### Етап F: Документація/Приклади
- Оновити README та додати code snippets використання для Async/Sync.
- Залишити нотатки міграції і поведінку при конфліктах із AutoReconnect.

## 4) Критерії приймання (Acceptance)
- Всі юніт/інтеграційні тести проходять без флаків, у т.ч. під дебагером.
- Відсутні паніки/рейс-кондишени (`-race` чисто).
- Без gorilla-залежностей у прод-коді камбітора (використовуємо restler API).
- Реконект/пауза/резюм працюють детерміновано; підписки відновлюються автоматично.
- Пінг-частота в безпечних межах; немає flood на сервер.

## 5) Тест-план (матриця)
- HealthConn:
  - Успішний Ping/Pong; таймаут Pong; N підряд фейлів → Unhealthy; успішні → Healthy.
  - Reconnect із backoff; jitter; збереження моніторингу після відновлення.
- MainConn:
  - Pause на Unhealthy (всі API fail-fast або блокуються залежно від режиму).
  - Resume на Healthy: Connect → WaitStarted → WaitReady → resubscribe; отримання даних.
- Режими:
  - Async: негайний `ErrNotConnected` під час Unhealthy; успіх після `WaitHealthy`.
  - Sync: блокування до Healthy/таймаут; перевірка часу очікування.
- Флапінг/стрес:
  - Часта зміна станів; відсутність deadlock/панік; відсутність «штормів» ресабскрайбу.
- Безпека:
  - Не логуються секрети; ліміти на пінги; дотримання таймаутів; структуровані логи.

## 6) Ризики і пом'якшення
- Ліміти на кількість з'єднань — зробити сервісний канал опційним; попередження в логах.
- Надто часті пінги — верхні/нижні межі конфігів; дефолти консервативні.
- Двоєдна логіка реконекту — пріоритет health monitor; AutoReconnect вимикається або делегує.
- Розходження станів — серіалізація переходів, ідемпотентні pause/resume.

## 7) Рол-аут
- Фаза 1 (A+B): API каркас + healthConn з монітором за фіче-флагом.
- Фаза 2 (C+D): Pause/Resume + режими; тестове покриття, вимкнення старих хелперів у тестах.
- Фаза 3 (E+F): Прибирання, документація, приклади, валідація в інтеграційному середовищі.
