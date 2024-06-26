package ui

import (
	"expenses-app/pkg/app/querying"
	"expenses-app/pkg/app/tracking"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const DEFAULT_DAYS_FROM_PARAM = "190"
const DEFAULT_DAYS_TO_PARAM = "0" // Now
const DEFAULT_PSIZE_PARAM = "35"
const DEFAULT_PNUM_PARAM = "0"

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
			panic("Form parsing error")
		}
		fmt.Println(payload)
		inputLayout := "2006-01-02"
		parsedDate, err := time.Parse(inputLayout, payload.Date)
		if err != nil {
			return c.Render("alerts/toastErr", fiber.Map{
				"Title": "Wrong Date",
				"Msg":   "Error parsing date",
			})

		}
		req := tracking.UpdateExpenseReq{
			Amount:     payload.Amount,
			CategoryID: payload.CategoryID,
			Date:       parsedDate,
			ExpenseID:  respExpense.Expenses[c.Params("id")].ID,
			People:     respExpense.Expenses[c.Params("id")].People,
			Product:    payload.Product,
			Shop:       payload.Shop,
		}
		_, err = eu.Update(req)
		if err != nil {
			panic("Update Error")
		}
		c.Append("Hx-Trigger", "reloadRow")
		return c.Render("alerts/toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   "Expense updated.",
		})
	}
}

func LoadTrackingSection() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			// c.Append("hx-trigger", "newPair")  // Not working :(
			return c.Render("main", fiber.Map{
				"ExpensesTrigger": "revealed",
			})
		}
		return c.Render("sections/tracking/index", fiber.Map{})
	}
}

func LoadExpenseRow(eq querying.ExpenseQuerier, cq querying.CategoryQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		respExpense, err := eq.GetByID(c.Params("id"))
		if err != nil {
			panic("Implement error")
		}
		return c.Render("sections/tracking/row", fiber.Map{
			"Expense": respExpense.Expenses[c.Params("id")],
		})
	}
}

func LoadExpenseEditRow(eq querying.ExpenseQuerier, cq querying.CategoryQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			// c.Append("hx-trigger", "newPair")  // Not working :(
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
			panic("Implement error")
		}
		fmt.Println(respCategories)
		fmt.Println(respExpense)
		return c.Render("sections/tracking/rowEdit", fiber.Map{
			"Expense":    respExpense.Expenses[c.Params("id")],
			"Categories": respCategories.Categories,
		})
	}
}

// LoadTrackingTable rendersn the Expenses section
func LoadTrackingTable(eq querying.ExpenseQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			fmt.Println("No HX-Request refreshing with revealed")
			// c.Append("hx-trigger", "newPair")  // Not working :(
			return c.Render("main", fiber.Map{
				"ExpensesTrigger": "revealed",
			})
		}
		pageNum, err := strconv.Atoi(c.Query("page_num", DEFAULT_PNUM_PARAM))
		if err != nil {
			panic("Implement me")
		}
		pageSize, err := strconv.Atoi(c.Query("page_size", DEFAULT_PSIZE_PARAM))
		if err != nil {
			panic("Implement me")
		}
		req := querying.ExpenseQuerierReq{
			// From:        time.Now().Add(-1 * time.Hour * 24 * time.Duration(daysFrom)),
			// To:          time.Now().Add(-1 * time.Hour * 24 * time.Duration(daysTo)),
			Page:        uint(pageNum),
			MaxPageSize: uint(pageSize),
		}
		resp, err := eq.Query(req)
		if err != nil {
			panic("Implement error UI")
		}
		return c.Render("sections/tracking/table", fiber.Map{
			"Expenses":    resp.Expenses,
			"CurrentPage": req.Page,      // Add this line
			"NextPage":    resp.Page + 1, // Add this line
			"PrevPage":    resp.Page - 1, // Add this line
			"PageSize":    resp.PageSize,
		})
	}
}