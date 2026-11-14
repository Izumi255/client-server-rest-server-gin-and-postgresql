# Практична робота № 5: Розробка REST-сервісу на базі фреймворку Gin та СУБД PostgreSQL

## Інформація про роботу
- **Дисципліна:** Розробка веб-застосунків / Backend-розробка
- **Тема:** Розробка REST-сервісу на базі фреймворку Gin та СУБД PostgreSQL
- **Проєкт:** Сервіс обліку технічних інцидентів та аварій (Incident Tracker API)
- **Стек технологій:** Go (Golang) 1.22+, Gin Web Framework, GORM, PostgreSQL, JWT (golang-jwt/jwt/v5), bcrypt (golang.org/x/crypto/bcrypt), Swagger UI (swaggo/gin-swagger), Docker.

---

## Зміст
1. [Опис проєкту](#опис-проєкту)
2. [Структура проєкту](#структура-проєкту)
3. [Вимоги та встановлення](#вимоги-та-встановлення)
4. [Запуск застосунку](#запуск-застосунку)
5. [Специфікація REST API](#специфікація-rest-api)
6. [Інтерактивна документація Swagger UI](#інтерактивна-документація-swagger-ui)
7. [Тестування API та приклади відповідей](#тестування-api-та-приклади-відповідей)
8. [Відповіді на контрольні питання (для захисту)](#відповіді-на-контрольні-питання-для-захисту)

---

## Опис проєкту

Проєкт реалізує повноцінний RESTful API сервіс для диспетчерських центрів та технічних служб підприємств з метою реєстрації, моніторингу та керування технічними інцидентами та аварійними подіями.

### Реалізований функціонал:
- **Архітектура на Gin:** Маршрутизація на базі Radix Tree за допомогою фреймворку `gin-gonic/gin`, групування ендпоінтів (`/api/v1`), обробка контексту запиту `gin.Context`.
- **Робота з СУБД PostgreSQL через GORM:**
  - Підключення через офіційний драйвер `gorm.io/driver/postgres`.
  - Автоматична міграція схем таблиць `incidents` та `users` через `AutoMigrate`.
  - Прив'язка інциденту до ідентифікатора інженера `user_id`.
- **Декларативна валідація вхідних даних (DTO):**
  - Використання тегів `binding` валідатора `go-playground/validator`.
  - Перевірка обов'язковості полів, мінімальної/максимальної довжини рядків та валідація формату email.
  - Обмеження значень полів фіксованим переліком (`oneof`): пріоритети (`low, medium, high, critical`), статуси (`open, in_progress, resolved, closed`).
- **Повний життєвий цикл інцидентів (CRUD):**
  - Отримання списку всіх інцидентів (`GET /api/v1/incidents`).
  - Отримання деталей конкретного інциденту за ID (`GET /api/v1/incidents/:id`) з коректною обробкою відсутності (код `404 Not Found`).
  - Створення нового інциденту від імені авторизованого користувача (`POST /api/v1/incidents`, код `201 Created`).
  - Оновлення статусу інциденту (`PATCH /api/v1/incidents/:id/status`).
  - Видалення інциденту (`DELETE /api/v1/incidents/:id`, код `204 No Content`).
- **Автентифікація та криптографічний захист (JWT + bcrypt):**
  - Безпечне гешування паролів за допомогою адаптивного алгоритму `bcrypt` (`golang.org/x/crypto/bcrypt`). Пароль ніколи не зберігається у відкритому вигляді.
  - Приховування хешу пароля від JSON-серіалізації тегом `json:"-"`.
  - Реєстрація користувачів з валідацією email та пароля (`POST /api/v1/auth/register`).
  - Авторизація користувача з перевіркою хешу у постійному часі (`CompareHashAndPassword`).
  - Генерація підписаних цифрових JWT-токенів (`golang-jwt/jwt/v5`) з терміном дії 24 години на базі алгоритму HMAC-SHA256 (`POST /api/v1/auth/login`).
- **Кастомний JWT Auth Middleware:**
  - Перевірка наявності заголовка `Authorization: Bearer <token>`.
  - Валідація цифрового підпису токена та терміну його дії.
  - Автоматичне блокування неавторизованих запитів зі статусом `401 Unauthorized` (`c.AbortWithStatusJSON`).
  - Витягування `userID` з токена та збереження в контексті запиту `c.Set("userID", claims.UserID)`.
  - Розподіл маршрутів на публічні (перегляд списку, реєстрація, вхід) та захищені (створення, зміна статусу, видалення).
- **Інтерактивна документація Swagger UI:**
  - Повний опис ендпоінтів та моделей за допомогою Swagger-анотацій.
  - Інтегрований вебінтерфейс Swagger UI (`/swagger/index.html`) з підтримкою авторизації через кнопку «Authorize» (схема `BearerAuth`).

---

## Структура проєкту

```text
incident-tracker/
├── database/
│   └── database.go          # Ініціалізація підключення до PostgreSQL та авто-міграція
├── docs/
│   └── docs.go              # Специфікація Swagger (OpenAPI 2.0)
├── dto/
│   ├── auth.go              # DTO для реєстрації та входу (RegisterDTO, LoginDTO)
│   └── incident.go          # DTO створення та оновлення статусу (CreateIncidentDTO, UpdateStatusDTO)
├── handlers/
│   ├── auth.go              # Обробники реєстрації та авторизації з випуском JWT
│   └── incident.go          # Обробники CRUD-операцій для сутності Incident
├── middleware/
│   └── auth.go              # Кастомний JWT Auth Middleware для захисту маршрутів
├── models/
│   ├── incident.go          # GORM-модель таблиці incidents (з прив'язкою user_id)
│   └── user.go              # GORM-модель таблиці users (хеш пароля приховано)
├── docker-compose.yml       # Швидкий запуск PostgreSQL у Docker-контейнері
├── requests.http            # Комплексний набір сценаріїв тестування (VS Code REST Client)
├── main.go                  # Точка входу, Swagger, захищені групи маршрутів та запуск
├── go.mod                   # Залежності проєкту
├── .gitignore               # Список виключень для Git
└── README.md                # Звіт, інструкція запуску та відповіді на контрольні питання
```

---

## Вимоги та встановлення

### Передумови
1. **Go SDK** версії 1.22 або новішої ([golang.org](https://go.dev/dl/)).
2. **PostgreSQL** 14+ (локально або через Docker).

### Встановлення залежностей
Перейдіть до каталогу проєкту та завантажте залежності:
```bash
go mod tidy
```

---

## Запуск застосунку

### 1. Запуск бази даних PostgreSQL
Якщо використовується Docker Desktop:
```bash
docker compose up -d
```
Або через пряму команду Docker:
```bash
docker run -d --name postgres-dev -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:alpine
```

### 2. Запуск бекенд-сервера
```bash
go run main.go
```

За замовчуванням сервер запускається на порті `8080`: `http://localhost:8080`.

---

## Специфікація REST API

| Метод | Маршрут | Доступ | Опис | Статус успіху | Помилка |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **GET** | `/swagger/*any` | Публічний | Інтерактивна документація Swagger UI | `200 OK` | — |
| **POST** | `/api/v1/auth/register` | Публічний | Реєстрація нового користувача | `201 Created` | `400 Bad Request` |
| **POST** | `/api/v1/auth/login` | Публічний | Вхід та генерація JWT токена | `200 OK` | `401 Unauthorized` |
| **GET** | `/api/v1/incidents` | Публічний | Отримання списку всіх інцидентів | `200 OK` | `500 Server Error` |
| **GET** | `/api/v1/incidents/:id` | Публічний | Отримання інциденту за ID | `200 OK` | `404 Not Found` |
| **POST** | `/api/v1/incidents` | **Bearer JWT** | Реєстрація інциденту (прив'язка до user_id) | `201 Created` | `400` / `401 Unauthorized` |
| **PATCH** | `/api/v1/incidents/:id/status`| **Bearer JWT** | Оновлення статусу інциденту | `200 OK` | `400` / `401` / `404` |
| **DELETE**| `/api/v1/incidents/:id` | **Bearer JWT** | Видалення запису про інцидент | `204 No Content`| `401` / `404` |

---

## Інтерактивна документація Swagger UI

Після запуску сервера відкрийте браузер за адресою:
👉 **`http://localhost:8080/swagger/index.html`**

### Як протестувати приватні ендпоінти у Swagger UI:
1. Виконайте запит `POST /api/v1/auth/login` та скопіюйте отриманий токен.
2. Вгорі сторінки Swagger UI натисніть зелену кнопку **«Authorize»**.
3. У полі значення введіть: `Bearer <ваш_jwt_токен>` та натисніть **Authorize**.
4. Тепер усі захищені маршрути (`POST /incidents`, `PATCH /incidents/{id}/status`, `DELETE /incidents/{id}`) успішно виконуватимуться прямо з браузера.

---

## Тестування API та приклади відповідей

### 1. Лог консолі під час успішного старту сервера:
```text
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
[GORM] Connected to PostgreSQL successfully (host=localhost:5432, db=postgres)
[GORM] Auto-migration completed: 'users', 'incidents' tables are up to date.
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.WrapHandler.func1 (3 handlers)
[GIN-debug] POST   /api/v1/auth/register     --> incident-tracker/handlers.RegisterHandler (3 handlers)
[GIN-debug] POST   /api/v1/auth/login        --> incident-tracker/handlers.LoginHandler (3 handlers)
[GIN-debug] GET    /api/v1/incidents         --> incident-tracker/handlers.GetIncidentsHandler (3 handlers)
[GIN-debug] GET    /api/v1/incidents/:id     --> incident-tracker/handlers.GetIncidentByIDHandler (3 handlers)
[GIN-debug] POST   /api/v1/incidents         --> incident-tracker/handlers.CreateIncidentHandler (4 handlers)
[GIN-debug] PATCH  /api/v1/incidents/:id/status --> incident-tracker/handlers.UpdateIncidentStatusHandler (4 handlers)
[GIN-debug] DELETE /api/v1/incidents/:id     --> incident-tracker/handlers.DeleteIncidentHandler (4 handlers)
[СЕРВЕР] Сервіс успішно запущено на http://localhost:8080
[СЕРВЕР] Swagger UI доступний за адресою: http://localhost:8080/swagger/index.html
```

---

### 2. Приклади HTTP-запитів та відповідей

#### 2.1. Авторизація та отримання токена JWT (POST /api/v1/auth/login)
**Запит:**
```http
POST http://localhost:8080/api/v1/auth/login HTTP/1.1
Content-Type: application/json

{
  "email": "engineer@company.internal",
  "password": "SuperSecretPassword123"
}
```

**Відповідь:**
```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsImV4cCI6MTczMTY4NDYwMH0.3J8..."
}
```

#### 2.2. Спроба доступу до захищеного маршруту БЕЗ токена (DELETE /api/v1/incidents/1)
**Запит:**
```http
DELETE http://localhost:8080/api/v1/incidents/1 HTTP/1.1
```

**Відповідь:**
```http
HTTP/1.1 401 Unauthorized
Content-Type: application/json; charset=utf-8

{
  "error": "Unauthorized access"
}
```

#### 2.3. Створення технічного інциденту З валідним токеном (POST /api/v1/incidents)
**Запит:**
```http
POST http://localhost:8080/api/v1/incidents HTTP/1.1
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "title": "Перегрів сервера Gateway-01",
  "description": "Температура процесора досягла 85C, спрацював термозахист",
  "priority": "critical"
}
```

**Відповідь:**
```http
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8

{
  "id": 1,
  "title": "Перегрів сервера Gateway-01",
  "description": "Температура процесора досягла 85C, спрацював термозахист",
  "priority": "critical",
  "status": "open",
  "user_id": 1,
  "created_at": "2026-11-14T15:31:00Z"
}
```

#### 2.4. Перевірка валідації DTO при некоректному пріоритеті
**Запит:**
```http
POST http://localhost:8080/api/v1/incidents HTTP/1.1
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "title": "Тестовий інцидент",
  "description": "Опис",
  "priority": "urgent"
}
```

**Відповідь:**
```http
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8

{
  "error": "Key: 'CreateIncidentDTO.Priority' Error:Field validation for 'Priority' failed on the 'oneof' tag"
}
```

#### 2.5. Оновлення статусу інциденту (PATCH /api/v1/incidents/1/status)
**Запит:**
```http
PATCH http://localhost:8080/api/v1/incidents/1/status HTTP/1.1
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "status": "in_progress"
}
```

**Відповідь:**
```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "id": 1,
  "title": "Перегрів сервера Gateway-01",
  "description": "Температура процесора досягла 85C, спрацював термозахист",
  "priority": "critical",
  "status": "in_progress",
  "user_id": 1,
  "created_at": "2026-11-14T15:31:00Z"
}
```

#### 2.6. Видалення інциденту З токеном (DELETE /api/v1/incidents/1)
**Запит:**
```http
DELETE http://localhost:8080/api/v1/incidents/1 HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Відповідь:**
```http
HTTP/1.1 204 No Content
```

---

## Відповіді на контрольні питання (для захисту)

### 1. У чому полягають головні архітектурні відмінності між фреймворком Gin та маршрутизатором Chi?
- **Gin** є повноцінним веб-фреймворком ("battery-included"). Він містить власну систему контексту (`*gin.Context`), вбудовану валідацію DTO (`validator.v10`), серіалізацію/десеріалізацію JSON/XML, пул пам'яті для зменшення навантаження на Garbage Collector та спеціалізоване дерево маршрутизації на базі Radix Tree. Обробники в Gin мають сигнатуру `func(*gin.Context)`.
- **Chi** — це мінімалістичний ідіоматичний HTTP-роутер, який повністю базується на стандартному інтерфейсі Go `net/http` (`http.Handler`, `http.HandlerFunc`) та стандартному `context.Context`. Chi не нав'язує власних абстракцій і надає переважно гнучку маршрутизацію та ланцюжки Middleware.

---

### 2. Яким чином сутність gin.Context оптимізує навантаження на Garbage Collector (GC) за допомогою пулу об’єктів?
- При обробці кожного вхідного HTTP-запиту необхідно зберігати параметри запиту, заголовки, посилання на ResponseWriter тощо. Якби на кожен запит створювався новий екземпляр `gin.Context` за допомогою динамічного виділення пам'яті в купі (`new(gin.Context)`), при високому навантаженні це спричиняло б часті зупинки програми збирачем сміття (GC pressure / GC pause).
- Фреймворк Gin використовує пул об'єктів стандартної бібліотеки **`sync.Pool`**. Коли запит надходить, екземпляр `gin.Context` береться з пулу готових об'єктів (`pool.Get()`), скидає попередній стан (`c.reset()`) та обслуговує поточний запит. Після завершення обробки контекст повертається назад у пул (`pool.Put(c)`), уникаючи нових алокацій пам'яті в heap.

---

### 3. Як працює зв’язування даних методом c.ShouldBindJSON та чим воно відрізняється від c.BindJSON у контексті обробки помилок?
- Обидва методи десеріалізують тіло вхідного HTTP-запиту формату JSON у передану структуру Go та запускають автоматичну валідацію полів за тегами `binding`.
- **Головна відмінність:**
  - **`c.BindJSON()`:** У разі помилки валідації чи синтаксису JSON автоматично перериває запит, записує у відповідь HTTP-статус `400 Bad Request` і встановлює заголовок відповіді. Розробник втрачає контроль над формуванням формату повідомлення про помилку.
  - **`c.ShouldBindJSON()`:** Не змінює статус відповіді самостійно, а просто повертає помилку `error`. Це дозволяє розробнику самостійно перевірити `if err != nil`, записати потрібний лог і повернути клієнту власний структурований JSON з деталями помилки через `c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})`.

---

### 4. Які переваги надає використання ORM-бібліотеки GORM порівняно з прямими запитами через драйвер бази даних database/sql?
- **Автоматичне відображення (Data Mapping):** GORM автоматично транслює рядки результатів SQL-запитів у типізовані структури Go та навпаки без необхідності ручного виклику `rows.Scan(&field1, &field2...)`.
- **Автоматична міграція схем (Auto-Migration):** Метод `db.AutoMigrate(&Model{})` аналізує структури Go та їхні теги `gorm:"..."` і автоматично створює або оновлює таблиці, типи колонок, індекси та зовнішні ключі в базі даних.
- **Захист від SQL Injection:** GORM за замовчуванням будує безпечні параметризовані SQL-запити.
- **Зручний CRUD та транзакції:** Надає лаконічний API (`db.Create()`, `db.First()`, `db.Save()`, `db.Delete()`, `db.Transaction()`) і підтримує Hooks (хуки життєвого циклу моделі, наприклад `BeforeCreate`).

---

### 5. Чому паролі користувачів заборонено гешувати простими геш-функціями (MD5, SHA-256) і чому алгоритм bcrypt вважається стандартом безпеки?
- **Прості геш-функції (MD5, SHA-1, SHA-256):** Проєктувалися для максимальної швидкодії (перевірка цілісності файлів, підпис повідомлень). Сучасні відеокарти (GPU) та ASIC-чіпи здатні перебирати мільярди таких гешів за секунду. Крім того, вони вразливі до атак за райдужними таблицями (Rainbow Tables).
- **Алгоритм bcrypt:**
  1. Є **повільною функцією деривації ключа** (key derivation function) на базі шифру Blowfish.
  2. Має вбудовану параметризацію складності (**Work Factor / Cost**). Складність можна збільшувати з розвитком потужностей процесорів, змушуючи обчислення одного гешу тривати, наприклад, 100-300 мс, що робить масовий перебір (brute-force) практично неможливим.
  3. Автоматично генерує та зберігає криптографічну **сіль (salt)** всередині самого рядка гешу, що унеможливлює використання Rainbow Tables.

---

### 6. З яких трьох частин складається структура токена JWT (Header, Payload, Signature) та як сервер перевіряє його цілісність без звернення до БД?
- Токен JWT складається з трьох частин, закодованих у форматі Base64URL і розділених крапками (`Header.Payload.Signature`):
  1. **Header (Заголовок):** Містить тип токена (`JWT`) та алгоритм підпису (наприклад, `HS256` — HMAC з SHA-256).
  2. **Payload (Корисне навантаження / Claims):** Містить передані дані (наприклад, `userID`, час створення `iat`, термін дії `exp`).
  3. **Signature (Цифровий підпис):** Обчислюється сервером за формулою:  
     `HMACSHA256(base64UrlEncode(Header) + "." + base64UrlEncode(Payload), secretKey)`.
- **Перевірка без звернення до БД:** Сервер бере отримані від клієнта `Header` та `Payload`, заново обчислює підпис за допомогою власного секретного ключа `secretKey`, який відомий лише серверу, і порівнює результат з переданим `Signature`. Якщо підписи ідентичні — дані не були підроблені третіми особами, і токен вважається валідним.

---

### 7. Як за допомогою методу c.AbortWithStatusJSON() у Middleware перервати виконання подальшого ланцюжка обробників?
- У фреймворку Gin виклик `c.Next()` передає керування наступному обробнику в ланцюжку (Middleware чи Handler).
- Якщо запит не пройшов перевірку (наприклад, відсутній або недійсний токен авторизації), виклик:
  ```go
  c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
  ```
  робить дві речі:
  1. Встановлює статус-код `401 Unauthorized` та відправляє клієнту JSON-відповідь.
  2. Встановлює внутрішній індекс ланцюжка обробників на максимальне значення, запобігаючи виконанню будь-яких наступних Middleware та цільового бізнес-хендлера.

---

### 8. Який принцип лежить в основі автоматичної генерації документації утилітою swag за допомогою спеціальних коментарів у коді?
- **Статичний аналіз коду (AST Parsing):** Утиліта `swag` не запускає сам додаток. Вона аналізує вихідний код програми на Go як текст, будуючи абстрактне синтаксичне дерево (Abstract Syntax Tree).
- Утиліта шукає спеціальні декларативні коментарі перед функцією `main()` та функціями-обробниками (наприклад: `@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`).
- На основі цих анотацій та структур DTO утиліта автоматично формує файли специфікації стандарту **OpenAPI / Swagger** (`docs.go`, `swagger.json`, `swagger.yaml`), які потім інтегруються у веб-інтерфейс Swagger UI через `ginSwagger.WrapHandler`.
