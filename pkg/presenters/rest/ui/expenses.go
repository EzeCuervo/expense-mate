package ui

import (
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const DEFAULT_DAYS_FROM_PARAM = "190"
const DEFAULT_DAYS_TO_PARAM = "0" // Now
const DEFAULT_PSIZE_PARAM = "30"
const DEFAULT_PNUM_PARAM = "0"

func DeleteExpense(ed tracking.ExpenseDeleter) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := tracking.DeleteExpenseReq{
			IDS: []string{c.Params("id")},
		}
		resp, err := ed.Delete(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Deletion Error",
				"Msg":   err.Error(),
			})
		}
		if len(resp.FailedDeletes) > 0 {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Could not delete expense",
				"Msg":   "Expense could not be deleted.",
			})
		}
		// c.Append("Hx-Trigger", "reloadExpensesTable")
		c.Append("Hx-Trigger", "deletedExpense-"+c.Params("id"))
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Deleted",
			"Msg":   "Expenses deleted successfully.",
		})
	}
}

func DeleteExpenses(ed tracking.ExpenseDeleter) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		payload := struct {
			IDs []string `json:"ids"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			panic("Form parsing error")
		}
		req := tracking.DeleteExpenseReq{
			IDS: payload.IDs,
		}
		resp, err := ed.Delete(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Deletion Error",
				"Msg":   err.Error(),
			})
		}
		if len(resp.FailedDeletes) > 0 {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Partial Deletion",
				"Msg":   "Some expenses could not be deleted.",
			})
		}
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Deleted",
			"Msg":   "Expenses deleted successfully.",
		})
	}
}

func CreateExpense(eu tracking.ExpenseCreator) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		payload := struct {
			Date       string  `form:"date"`
			Shop       string  `form:"shop"`
			Product    string  `form:"product"`
			CategoryID string  `form:"category"`
			Amount     float64 `form:"amount"`
		}{}
		selectedUsers := slices.DeleteFunc(strings.Split(c.FormValue("users"), ","), func(s string) bool { return s == "" })
		if err := c.BodyParser(&payload); err != nil {
			panic("Form parsing error")
		}
		inputLayout := "2006-01-02"
		parsedDate, err := time.Parse(inputLayout, payload.Date)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Wrong Date",
				"Msg":   "Error parsing date",
			})
		}
		req := tracking.CreateExpenseReq{
			Product:    payload.Product,
			Amount:     payload.Amount,
			Shop:       payload.Shop,
			Date:       parsedDate,
			UsersID:    selectedUsers,
			CategoryID: payload.CategoryID,
		}
		_, err = eu.Create(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Can't create expense",
				"Msg":   err,
			})
		}
		c.Append("Hx-Trigger", "reloadExpensesTable")
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   "Expense created.",
		})
	}
}

func LoadExpenseRow(eq querying.ExpenseQuerier, cq querying.CategoryQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		respExpense, err := eq.GetByID(c.Params("id"))
		if err != nil {
			panic("Implement error")
		}
		// c.Append("Hx-Trigger", fmt.Sprintf("reloadRow-%s", c.Params("id")))
		return c.Render("sections/expenses/row", fiber.Map{
			"Expense": respExpense.Expenses[0],
		})
	}
}

func LoadAddExpensesRow(cq querying.CategoryQuerier, mu managing.UserManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		respCategories, err := cq.Query()
		if err != nil {
			panic("Implement error")
		}
		respUsers, err := mu.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "User error",
				"Msg":   "Could not load users.",
			})
		}
		return c.Render("sections/expenses/rowAdd", fiber.Map{
			"Categories": respCategories.Categories,
			"NoUserID":   querying.NoUserID,
			"Users":      respUsers.Users,
		})
	}
}

func EditExpense(eq querying.ExpenseQuerier, eu tracking.ExpenseUpdater) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		respExpense, err := eq.GetByID(c.Params("id"))
		if err != nil {
			panic("Implement error")
		}
		// Payload to unmarshal te form
		payload := struct {
			Date       string  `form:"date"`
			Shop       string  `form:"shop"`
			Product    string  `form:"product"`
			CategoryID string  `form:"category"`
			Amount     float64 `form:"amount"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Form",
				"Msg":   "Error parsing form",
			})
		}
		selectedUsers := slices.DeleteFunc(strings.Split(c.FormValue("users"), ","), func(s string) bool { return s == "" })
		inputLayout := "2006-01-02"
		parsedDate, err := time.Parse(inputLayout, payload.Date)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Wrong Date",
				"Msg":   "Error parsing date. Please use YYYY-MM-DD format.",
			})
		}
		req := tracking.UpdateExpenseReq{
			Amount:     payload.Amount,
			UsersID:    selectedUsers,
			CategoryID: payload.CategoryID,
			Date:       parsedDate,
			ExpenseID:  respExpense.Expenses[0].ID,
			Product:    payload.Product,
			Shop:       payload.Shop,
		}
		_, err = eu.Update(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Creation error",
				"Msg":   "Error updating expense",
			})
		}
		c.Append("Hx-Trigger", "reloadRow-"+c.Params("id"))
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   "Expense updated.",
		})
	}
}

func LoadExpenseEditRow(eq querying.ExpenseQuerier, cq querying.CategoryQuerier, mu managing.UserManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			return c.Render("main", fiber.Map{
				"ExpensesTrigger": "revealed",
			})
		}
		respCategories, err := cq.Query()
		if err != nil {
			panic("Implement error")
		}
		respExpense, err := eq.GetByID(c.Params("id"))
		if err != nil {
			fmt.Println(err)
			panic("Implement error")
		}
		respUsers, err := mu.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "User error",
				"Msg":   "Could not load users.",
			})
		}
		return c.Render("sections/expenses/rowEdit", fiber.Map{
			"Expense":    respExpense.Expenses[0],
			"Categories": respCategories.Categories,
			"NoUserID":   querying.NoUserID,
			"Users":      respUsers.Users,
		})
	}

}
func LoadExpenseFilter(cq querying.CategoryQuerier, mu managing.UserManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		respCategories, err := cq.Query()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Expense filter",
				"Msg":   "Could not load the expense filter.",
			})
		}
		respUsers, err := mu.List()
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "User error",
				"Msg":   "Could not load users.",
			})
		}
		return c.Render("sections/expenses/filter", fiber.Map{
			"Categories": respCategories.Categories,
			"NoUserID":   querying.NoUserID,
			"Users":      respUsers.Users,
		})
	}
}

// LoadExpensesTable rendersn the Expenses section
func LoadExpensesTable(eq querying.ExpenseQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		pageNum, err := strconv.Atoi(c.Query("page_num", DEFAULT_PNUM_PARAM))
		if err != nil {
			panic("Atoi parse error")
		}
		pageSize, err := strconv.Atoi(c.Query("page_size", DEFAULT_PSIZE_PARAM))
		if err != nil {
			panic("Atoi parse error")
		}
		fromDate, err := time.Parse("2006-01-02", c.Query("from-date", time.Time{}.Format("2006-01-02")))
		if err != nil {
			panic("Date parse error")
		}
		toDate, err := time.Parse("2006-01-02", c.Query("to-date", time.Time{}.Format("2006-01-02")))
		if err != nil {
			panic("Date parse error")
		}
		min_amount, err := strconv.Atoi(c.Query("min_amount", "0"))
		if err != nil {
			panic("Atoi parse error")
		}
		max_amount, err := strconv.Atoi(c.Query("max_amount", "0"))
		if err != nil {
			panic("Atoi parse error")
		}
		selectedUsers := slices.DeleteFunc(strings.Split(c.Query("users"), ","), func(s string) bool { return s == "" })
		categories := slices.DeleteFunc(strings.Split(c.Query("categories"), ","), func(s string) bool { return s == "" })
		req := querying.ExpenseQuerierReq{
			Page:        uint(pageNum),
			MaxPageSize: uint(pageSize),
			ExpenseFilter: querying.ExpenseQuerierFilter{
				ByCategoryID: categories,
				ByUsers:      selectedUsers,
				ByShop:       c.Query("shop"),
				ByProduct:    c.Query("product"),
				ByAmount:     [2]uint{uint(min_amount), uint(max_amount)},
				ByTime:       [2]time.Time{fromDate, toDate},
			},
		}
		fmt.Println(req.ExpenseFilter.ByUsers)
		resp, err := eq.Query(req)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Table error",
				"Msg":   "Error loading expenses table.",
			})
		}
		return c.Render("sections/expenses/table", fiber.Map{
			"Expenses":      resp.Expenses,
			"CurrentPage":   resp.Page,     // Add this line
			"NextPage":      resp.Page + 1, // Add this line
			"PrevPage":      resp.Page - 1, // Add this line
			"TotalPages":    uint(math.Ceil(float64(resp.ExpensesCount / req.MaxPageSize))),
			"PageSize":      resp.PageSize,
			"ExpensesCount": resp.ExpensesCount,
		})
	}
}
