package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"tgBotWithDataBase/models"

	_ "github.com/lib/pq"
)

func setUpTestDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password='C529@4H0OdEO%}Y' dbname=puvelka_test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД - %w", err)
	}

	errPing := db.Ping()
	if errPing != nil {
		return nil, fmt.Errorf("ошибка пинга БД - %w", errPing)
	}
	_, errDelete := db.Exec(`TRUNCATE TABLE expensesTG RESTART IDENTITY;`)
	if errDelete != nil {
		return nil, fmt.Errorf("ошибка удаления полей из БД - %w", errDelete)
	}
	return db, nil
}

func TestGetExpensesAllTime(t *testing.T) {
	db, err := setUpTestDB()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	testData1 := `INSERT INTO expensesTG (amount, category, comment) VALUES (100, 'еда', 'тест1')`
	_, errExec1 := db.Exec(testData1)
	if errExec1 != nil {
		t.Errorf("ошибка запроса к БД - %v", err)
	}

	testData2 := `INSERT INTO expensesTG (amount, category, comment) VALUES (200, 'метро', 'тест2')`
	_, errExec2 := db.Exec(testData2)
	if errExec2 != nil {
		t.Errorf("ошибка запроса к БД - %v", err)
	}

	testData3 := `INSERT INTO expensesTG (amount, category, comment) VALUES (300, 'шаурма', 'тест3')`
	_, errExec3 := db.Exec(testData3)
	if errExec3 != nil {
		t.Errorf("ошибка запроса к БД - %v", err)
	}

	expenses, errExpensesAllTime := GetExpensesAllTime(db)
	if errExpensesAllTime != nil {
		t.Errorf("%v", errExpensesAllTime)
	}

	if len(expenses) == 0 {
		t.Errorf("длина не должна быть 0, а должна быть %v", len(expenses))
	}
	expected0 := "100.00\nеда\nтест1"
	//expected1 :=
	//expected2 :=

	str0 := ""
	//str1 := ""
	//str2 := ""
	for index, exp := range expenses {
		if index == 0 {
			str0 += strconv.FormatFloat(exp.Amount, 'f', 2, 64) + "\n"
			str0 += exp.Category + "\n"
			str0 += exp.Comment
		}
	}
	fmt.Printf("Ожидается - %q\n", expected0)
	fmt.Printf("Получено - %q\n", str0)
	if expected0 != str0 {
		t.Errorf("ожидается\n%v, а получено\n%v", expected0, str0)
	}
}

func InsertForTestExpensesByCategories(t *testing.T) {
	db, err := setUpTestDB()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	insert1 := `INSERT INTO expensesTG (amount, category, comment) VALUES (100, 'шаурма', '')`
	insert2 := `INSERT INTO expensesTG (amount, category, comment) VALUES (200, 'шаурма', 'что-то')`
	insert3 := `INSERT INTO expensesTG (amount, category, comment) VALUES (300, 'метро', 'тест3')`
	insert4 := `INSERT INTO expensesTG (amount, category, comment) VALUES (600, 'метро', 'тест4 пускай будет')`

	_, err1 := db.Exec(insert1)
	if err1 != nil {
		t.Errorf("ошибка запроса на вставку в 1 попытке")
	}

	_, err2 := db.Exec(insert2)
	if err2 != nil {
		t.Errorf("ошибка запроса на вставку в 2 попытке")
	}

	_, err3 := db.Exec(insert3)
	if err3 != nil {
		t.Errorf("ошибка запроса на вставку в 3 попытке")
	}

	_, err4 := db.Exec(insert4)
	if err4 != nil {
		t.Errorf("ошибка запроса на вставку в 4 попытке")
	}
}

func TestExpensesByCategories(t *testing.T) {
	db, err := setUpTestDB()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	InsertForTestExpensesByCategories(t)

	result, errExpensesByCategories := ExpensesByCategories(db)
	if errExpensesByCategories != nil {
		t.Errorf("%v", errExpensesByCategories)
	}

	expectedAmount1 := 300
	expectedAmount2 := 900

	expectedCategory1 := "шаурма"
	expectedCategory2 := "метро"

	if result[expectedCategory1] != float64(expectedAmount1) {
		t.Errorf("ожидалась сумма - %v, а получено - %v", expectedAmount1, result[expectedCategory2])
	}

	if result[expectedCategory2] != float64(expectedAmount2) {
		t.Errorf("ожидалась сумма - %v, а получено - %v", expectedAmount1, result[expectedCategory2])
	}
}

type ExpectedExpense struct {
	ID       int     `json:"id"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Comment  string  `json:"comment"`
}

func TestAddExpense(t *testing.T) {
	db, err := setUpTestDB()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	testExpense := models.GetExpense{
		Amount:   150,
		Category: "метро",
		Comment:  "прокатился с кайфом",
	}

	res, errAdd := AddExpense(db, testExpense)
	if errAdd != nil {
		t.Errorf("%v", errAdd)
	}

	expected := ExpectedExpense{
		ID:       1,
		Amount:   150,
		Category: "метро",
		Comment:  "прокатился с кайфом",
	}

	if res.ID != expected.ID {
		t.Errorf("ожидалось %v, а пришло %v", expected.ID, res.ID)
	}

	if res.Amount != expected.Amount {
		t.Errorf("ожидалось %v, а пришло %v", expected.Amount, res.Amount)
	}

	if res.Category != expected.Category {
		t.Errorf("ожидалось %v, а пришло %v", expected.Category, res.Category)
	}

	if res.Comment != expected.Comment {
		t.Errorf("ожидалось %v, а пришло %v", expected.Comment, res.Comment)
	}
}

func InsertForDelete(t *testing.T) {
	db, err := setUpTestDB()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	_, err1 := db.Exec(`INSERT INTO expensesTG (amount, category, comment) VALUES (150, 'еда', 'ресторан') RETURNING id, created_at`)
	if err1 != nil {
		t.Error("ошибка запроса на вставку")
	}
}

func TestDeleteExpense(t *testing.T) {
	db, err := setUpTestDB()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	InsertForDelete(t)
	id := 1
	errDelete1 := DeleteExpense(db, id)
	if errDelete1 != nil {
		t.Error("не получилось удалить")
	}

	errDelete2 := DeleteExpense(db, id)
	if errDelete2 == nil {
		t.Errorf("ожидается ошибка - %v", errDelete2)
	}
}
