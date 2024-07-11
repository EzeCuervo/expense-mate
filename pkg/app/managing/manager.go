package managing

// Service just holds all the managing use cases
type Service struct {
	// UserCreator     UsersCreator
	CategoryDeleter       CategoryDeleter
	CategoryCreator       CategoryCreator
	CategoryUpdater       CategoryUpdater
	CategoryBudgetCreator CategoryBudgetCreator
	TelegramCommander     TelegramCommander
	RuleManager           RuleManager
	UserManager           UserManager
}

// NewService is the interctor for all Managing Use cases
func NewService(cd CategoryDeleter, cc CategoryCreator, cu CategoryUpdater, cbc CategoryBudgetCreator, tc TelegramCommander, rm RuleManager, um UserManager) Service {
	return Service{cd, cc, cu, cbc, tc, rm, um}
}
