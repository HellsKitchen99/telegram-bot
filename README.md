# 📊 Telegram Финансовый Бот

Телеграм-бот на Go для учета личных расходов. Сохраняет траты, показывает статистику и управляется через 6 удобных кнопок.

## 🔧 Возможности

- ✅ Добавление расходов (пример: `155.5 Еда Обед`)
- 🗑 Удаление по ID
- 📅 Вывод всех трат
- 📂 Группировка по категориям
- 📎 Инлайн-кнопки (`Start`, `Help`, `Add`, `Delete`, `Expenses (All time)`, `Expenses (By categories)`)

## 🛠 Технологии

- Язык: **Go**
- База данных: **PostgreSQL**
- Telegram API: [go-telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api)
- Тестирование: **Go testing**, **mock-бот**

## 📁 Структура проекта

```
.
├── main.go                # запуск Telegram-бота
├── handlers.go           # обработка команд и кнопок
├── expenses.go           # бизнес-логика для работы с расходами
├── db.go                 # подключение и операции с БД
├── models/               # структуры данных
├── mock/                 # фейковый бот для тестов
├── .env.example          # пример .env с переменными
├── go.mod / go.sum       # зависимости
└── *_test.go             # тесты
```

## ⚙️ Установка и запуск

1. **Создай БД в PostgreSQL** и таблицу:

```sql
CREATE TABLE expensesTG (
  id SERIAL PRIMARY KEY,
  amount DOUBLE PRECISION NOT NULL,
  category TEXT NOT NULL,
  comment TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);
```

2. **Создай `.env` на основе `.env.example`** и добавь свой токен:

```env
Token=your_telegram_token_here
```

3. **Установи зависимости и запусти бота**:

```bash
go mod tidy
go run main.go
```

## 🧪 Тесты

В проекте написаны модульные тесты на:

- Добавление трат
- Удаление трат
- Получение трат по категориям и за всё время
- Хендлеры и Telegram API (с мок-ботом)

Запуск тестов:
```bash
go test ./...
```

## 💬 Примеры команд

- Ввод: `120.50 Еда обед`
- Ответ:  
  ```
  Трата добавлена  
  ID траты - 1  
  Сумма - 120.5  
  Категория - Еда  
  Комментарий - обед  
  Время траты - 2025-06-02 13:45:00
  ```