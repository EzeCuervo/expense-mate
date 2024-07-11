package managing

import (
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/expense"
)

type CreateCategoryBudgetResp struct {
	ID  string
	Msg string
}

type CreateCategoryBudgetReq struct {
	CategoryID int
	Amount     float64
	Year       int
	Month      int
}

// CategoryBudgetCreator use case creates a category budget for a month/year
type CategoryBudgetCreator struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewCategoryBudgetCreator(l app.Logger, e expense.Expenses) *CategoryBudgetCreator {
	return &CategoryBudgetCreator{l, e}
}

// Create use cases function creates a new category budget
func (s *CategoryBudgetCreator) Create(req CreateCategoryBudgetReq) (*CreateCategoryBudgetResp, error) {
	
	return nil, nil
}
