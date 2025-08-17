# Lo

REST API на Go для управления задачами с асинхронной обработкой и логированием.

## Возможности
- **Эндпоинты**:
    - `POST /tasks` — создать задачу (`{"title": "string", "description": "string"}`).
    - `GET /tasks` — получить список задач (фильтр по `status` через query).
    - `GET /tasks/{id}` — получить задачу по ID.
- **Статусы задач**: `Pending`, `Completed`, `Failed`, `Error`.
- **Обработка**: Задачи выполняются 1–2 минуты. 20% шанс `Failed`, до 2 ретраев, затем `Error`.
- **Логирование**: Все изменения статусов логируются в stdout через канал.
- **Graceful Shutdown**: При остановке (Ctrl+C) выводятся все задачи (ID, Title, Status, Retries).

## Сборка и запуск
1. **Склонируйте репозиторий**:
   ```bash
   git clone https://github.com/SharpDenin/lo
   cd lo
   ```
2. **Инициализируйте зависимости**:
   ```bash
   go mod tidy
   ```
3. **Соберите приложение**:
   ```bash
   go build -o lo ./cmd/lo/main.go
   ```
4. **Запустите сервер**:
   ```bash
   go run ./cmd/lo/main.go
   ```

## Тестирование
Используйте `curl` для проверки.

### 1. Создание задачи
```bash
curl -X POST -H "Content-Type: application/json" -d '{"title":"TestTask","description":"Test"}' http://localhost:8080/tasks
```
**Ответ** (HTTP 201):
```json
{"data":{"id":1,"title":"TestTask","description":"Test","status":"Pending","retries":0},"meta":null}
```

### 2. Создание множества задач
```bash
for i in {1..10}; do curl -X POST -H "Content-Type: application/json" -d "{\"title\":\"Task$i\",\"description\":\"Desc$i\"}" http://localhost:8080/tasks; done
```
Все задачи создадутся.

### 3. Получение списка задач
```bash
curl http://localhost:8080/tasks
```
**Ответ** (HTTP 200):
```json
{"data":[{"id":1,"title":"TestTask","description":"Test","status":"Pending","retries":0}, ...],"meta":null}
```
Фильтрация по статусу:
```bash
curl http://localhost:8080/tasks?status=Pending
```

### 4. Получение задачи по ID
```bash
curl http://localhost:8080/tasks/1
```
**Ответ** (HTTP 200):
```json
{"data":{"id":1,"title":"TestTask","description":"Test","status":"Pending","retries":0},"meta":null}
```

### 5. Проверка асинхронной обработки
Подождите 1–2 минуты после создания задачи:
```bash
curl http://localhost:8080/tasks/1
```
Статус изменится на `Completed` (80%) или `Failed`/`Error` (20% после ретраев).

### 6. Проверка graceful shutdown
1. Создайте задачи.
2. Нажмите `Ctrl+C`.
3. Проверьте логи:
   ```
   Инициируется завершение работы сервера...
   Состояние задач на момент завершения:
   Task ID: 1, Title: TestTask, Status: Completed, Retries: 0
   ...
   Сервер успешно завершен
   ```
