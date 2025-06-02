package handlers

import (
	"database/sql"
	"fmt"
	"testing"
	"tgBotWithDataBase/mock"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func setUpTestDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password='C529@4H0OdEO%}Y' dbname=puvelka_test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД")
	}
	errPing := db.Ping()
	if errPing != nil {
		return nil, fmt.Errorf("не удалось пингануть БД")
	}
	_, errDelete := db.Exec(`TRUNCATE TABLE expensesTG RESTART IDENTITY;`)
	if errDelete != nil {
		return nil, fmt.Errorf("ошибка удаления полей из БД")
	}
	return db, nil
}

func InsertTests(db *sql.DB) error {
	_, err1 := db.Exec(`INSERT INTO expensesTG (amount, category, comment) VALUES (1, 'test', 'rrr')`)
	if err1 != nil {
		return fmt.Errorf("ошибка из 1")
	}

	_, err2 := db.Exec(`INSERT INTO expensesTG (amount, category, comment) VALUES (1, 'test', 'rrr')`)
	if err2 != nil {
		return fmt.Errorf("ошибка из 2")
	}

	_, err3 := db.Exec(`INSERT INTO expensesTG (amount, category, comment) VALUES (1, 'test', 'rrr')`)
	if err3 != nil {
		return fmt.Errorf("ошибка из 3")
	}
	return nil
}

func TestHandlerExpensesByCategories(t *testing.T) {
	dbTest, errDB := setUpTestDB()
	if errDB != nil {
		t.Errorf("%v", errDB)
	}

	err := InsertTests(dbTest)
	if err != nil {
		t.Errorf("%v", err)
	}

	bot := mock.NewFakeBot()
	update := telegram.Update{
		CallbackQuery: &telegram.CallbackQuery{
			From: &telegram.User{ID: 555},
			Message: &telegram.Message{
				Chat: &telegram.Chat{ID: 1001},
			},
		},
	}
	errNew := HandlerExpensesByCategories(bot, &update, dbTest)
	if errNew != nil {
		t.Errorf("%v", errNew)
	}
	if len(bot.SentMessages) != 1 {
		t.Errorf("ожидалась длина 1, а не %v", len(bot.SentMessages))
	}
	msg, ok := bot.SentMessages[0].(telegram.MessageConfig)
	if !ok {
		t.Error("не удалось преобразовать")
	}

	output := fmt.Sprintf("Траты по категориям\n%v - %v\n", "test", 3)
	fmt.Printf("%q\n", output)
	fmt.Printf("%q\n", msg.Text)
	if output != msg.Text {
		t.Errorf("ожидалось %v, а пришло %v", output, msg.Text)
	}
}

func TestHandlerAdd(t *testing.T) {
	dbTest, errDB := setUpTestDB()
	if errDB != nil {
		t.Errorf("%v", errDB)
	}
	bot := mock.NewFakeBot()
	update := telegram.Update{
		Message: &telegram.Message{
			Text: "120.56 метро проезд",
			Chat: &telegram.Chat{ID: 1001},
			From: &telegram.User{ID: 555},
		},
	}
	status := map[int64]string{555: "add"}
	id, timee, err := HandlerAdd(bot, &update, status, dbTest)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(bot.SentMessages) != 1 {
		t.Errorf("ожидалась длина 1, а не %v", len(bot.SentMessages))
	}
	msg, ok := bot.SentMessages[0].(telegram.MessageConfig)
	if !ok {
		t.Error("не удалось привести к MessageConfig")
	}
	expected := fmt.Sprintf("Трата добавлена\nID траты - %v\nСумма - 120.56\nКатегория - метро\nКомментарий - проезд\nВремя траты - %v", id, timee.Format("2006-01-02 15:04:05"))
	if msg.Text != expected {
		t.Errorf("ожидалось\n%v\nа пришло\n%v", expected, msg.Text)
	}
}

func InsertForDelete(db *sql.DB) error {
	_, err1 := db.Exec("INSERT INTO expensesTG (amount, category, comment) VALUES (111, 'food', 'i am daredevil')")
	if err1 != nil {
		return fmt.Errorf("ошибка 1 запроса")
	}
	_, err2 := db.Exec("INSERT INTO expensesTG (amount, category, comment) VALUES (222, 'drinks', 'not even God can stop me')")
	if err2 != nil {
		return fmt.Errorf("ошибка 2 запроса")
	}
	return nil
}

func TestHandlerDelete(t *testing.T) {
	dbTest, errDB := setUpTestDB()
	if errDB != nil {
		t.Errorf("%v", errDB)
	}

	errTest := InsertForDelete(dbTest)
	if errTest != nil {
		t.Errorf("%v", errTest)
	}

	bot := mock.NewFakeBot()
	update := telegram.Update{
		Message: &telegram.Message{
			Text: "1",
			Chat: &telegram.Chat{ID: 1001},
			From: &telegram.User{ID: 555},
		},
	}
	status := map[int64]string{555: "delete"}
	err := HandlerDelete(bot, &update, status, dbTest)
	if err != nil {
		t.Errorf("%v", err)
	}

	if len(bot.SentMessages) != 1 {
		t.Errorf("ожидалась длина 1, а пришла %v", len(bot.SentMessages))
	}

	msg, ok := bot.SentMessages[0].(telegram.MessageConfig)
	if !ok {
		t.Errorf("не удалось провести приведение типа")
	}
	expected := "трата удалена"
	fmt.Printf("%q\n", msg.Text)
	fmt.Printf("%q\n", expected)
	if msg.Text != expected {
		t.Errorf("ожидалось %v, а пришло %v", expected, msg.Text)
	}
}

func InsertTestsForTestHandlerExpensesAllTime(db *sql.DB) (int, time.Time, int, time.Time, int, time.Time, error) {
	var id1, id2, id3 int
	var cr1, cr2, cr3 time.Time
	row1 := db.QueryRow(`INSERT INTO expensesTG (amount, category, comment) VALUES (100, 'food', 'shavuha') RETURNING id, created_at`)
	err1 := row1.Scan(&id1, &cr1)
	if err1 != nil {
		return 0, time.Time{}, 0, time.Time{}, 0, time.Time{}, fmt.Errorf("ошибка при сканировании айди1 и времени1")
	}

	row2 := db.QueryRow(`INSERT INTO expensesTG (amount, category, comment) VALUES (200, 'hells kitchen', 'get up Matty') RETURNING id, created_at`)
	err2 := row2.Scan(&id2, &cr2)
	if err2 != nil {
		return 0, time.Time{}, 0, time.Time{}, 0, time.Time{}, fmt.Errorf("ошибка при сканировании айди2 и времени2")
	}

	row3 := db.QueryRow(`INSERT INTO expensesTG (amount, category, comment) VALUES (300, 'hhh', 'hhh') RETURNING id, created_at`)
	err3 := row3.Scan(&id3, &cr3)
	if err3 != nil {
		return 0, time.Time{}, 0, time.Time{}, 0, time.Time{}, fmt.Errorf("ошибка при сканировании айди3 и времени3")
	}
	return id1, cr1, id2, cr2, id3, cr3, nil
}

func TestHandlerExpensesAllTime(t *testing.T) {
	dbTest, errDB := setUpTestDB()
	if errDB != nil {
		t.Errorf("%v", errDB)
	}
	id1, cr1, id2, cr2, id3, cr3, err := InsertTestsForTestHandlerExpensesAllTime(dbTest)
	if err != nil {
		t.Errorf("%v", err)
	}

	bot := mock.NewFakeBot()
	update := telegram.Update{
		CallbackQuery: &telegram.CallbackQuery{
			From: &telegram.User{ID: 555},
			Message: &telegram.Message{
				Chat: &telegram.Chat{ID: 1001},
			},
		},
	}

	err1 := HandlerExpensesAllTime(bot, &update, dbTest)
	if err1 != nil {
		t.Errorf("%v", err1)
	}

	msg, ok := bot.SentMessages[0].(telegram.MessageConfig)
	if !ok {
		t.Errorf("не удалось привести")
	}
	//fmt.Println(msg.Text, id1, cr1, id, cr2, id3, cr3)
	expected := fmt.Sprintf("Траты за все время\n\nТрата №1\nID - %v\nСумма - 100.00\nКатегория - food\nКоммент - shavuha\nВремя - %v\n\nТрата №2\nID - %v\nСумма - 200.00\nКатегория - hells kitchen\nКоммент - get up Matty\nВремя - %v\n\nТрата №3\nID - %v\nСумма - 300.00\nКатегория - hhh\nКоммент - hhh\nВремя - %v\n\n", id1, cr1.Format("2006-01-02 15:04:05"), id2, cr2.Format("2006-01-02 15:04:05"), id3, cr3.Format("2006-01-02 15:04:05"))
	if expected != msg.Text {
		t.Errorf("ожидалось %v, а пришло %v", expected, msg.Text)
	}
}

//"2006-01-02 15:04:05"
