package storage

import (
	"time"

	"github.com/google/uuid"
)

// User represents the users table
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName string    `gorm:"type:varchar(50);unique;not null" json:"first_name"`
	LastName  string    `gorm:"type:varchar(50);unique;not null" json:"last_name"`
	Password  string    `gorm:"type:varchar(255);unique;not null" json:"password"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	Accounts   []Account  `gorm:"foreignKey:UserID" json:"accounts,omitempty"`
	Categories []Category `gorm:"foreignKey:UserID" json:"categories,omitempty"`
	Budgets    []Budget   `gorm:"foreignKey:UserID" json:"budgets,omitempty"`
	Goals      []Goal     `gorm:"foreignKey:UserID" json:"goals,omitempty"`
}

// Account represents the accounts table with dual balances
type Account struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name             string    `gorm:"type:varchar(50);not null" json:"name"`
	SpendableBalance float64   `gorm:"column:spendable_balance;type:numeric(15,2);not null;default:0.00" json:"spendable_balance"`
	SavingsBalance   float64   `gorm:"column:savings_balance;type:numeric(15,2);not null;default:0.00" json:"savings_balance"`
	Currency         string    `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	CreatedAt        time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:AccountID" json:"transactions,omitempty"`
}

// Category represents the main budget categories table
type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	SubCategories []SubCategory `gorm:"foreignKey:CategoryID" json:"sub_categories,omitempty"`
}

// SubCategory represents the sub_categories table
type SubCategory struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	Name       string    `gorm:"type:varchar(50);not null" json:"name"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// Budget represents the monthly rolling budget configuration
type Budget struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	SubCategoryID     uuid.UUID `gorm:"type:uuid;not null" json:"sub_category_id"`
	BudgetMonth       time.Time `gorm:"column:budget_month;type:date;not null" json:"budget_month"` // Store as first day of month
	AllocatedAmount   float64   `gorm:"column:allocated_amount;type:numeric(15,2);not null;default:0.00" json:"allocated_amount"`
	CarriedOverAmount float64   `gorm:"column:carried_over_amount;type:numeric(15,2);not null;default:0.00" json:"carried_over_amount"`
	CurrentSpend      float64   `gorm:"column:current_spend;type:numeric(15,2);not null;default:0.00" json:"current_spend"`
	CreatedAt         time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	SubCategory *SubCategory `gorm:"foreignKey:SubCategoryID" json:"sub_category,omitempty"`
}

// Goal represents target goals linked to an account's savings pool
type Goal struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Name         string     `gorm:"type:varchar(100);not null" json:"name"`
	TargetAmount float64    `gorm:"column:target_amount;type:numeric(15,2);not null" json:"target_amount"`
	SavedAmount  float64    `gorm:"column:saved_amount;type:numeric(15,2);not null;default:0.00" json:"saved_amount"`
	IsCompleted  bool       `gorm:"column:is_completed;default:false" json:"is_completed"`
	TargetDate   *time.Time `gorm:"column:target_date;type:date" json:"target_date,omitempty"` // Pointer to handle nullable dates
	CreatedAt    time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// Transaction represents the ledger entries with standard check types
type Transaction struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	AccountID       uuid.UUID  `gorm:"type:uuid;not null" json:"account_id"`
	SubCategoryID   *uuid.UUID `gorm:"column:sub_category_id;type:uuid" json:"sub_category_id,omitempty"` // Nullable for pure transfers/credits
	Type            string     `gorm:"type:varchar(10);not null" json:"type"`                             // 'Credit' or 'Debit'
	Amount          float64    `gorm:"type:numeric(15,2);not null" json:"amount"`
	Note            string     `gorm:"type:text" json:"note,omitempty"`
	TransactionDate time.Time  `gorm:"column:transaction_date;default:CURRENT_TIMESTAMP" json:"transaction_date"`

	Account     *Account     `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	SubCategory *SubCategory `gorm:"foreignKey:SubCategoryID" json:"sub_category,omitempty"`
}

// SavingsWithdrawal tracks emergency dips into goal savings
type SavingsWithdrawal struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID   uuid.UUID `gorm:"type:uuid;not null" json:"account_id"`
	GoalID      uuid.UUID `gorm:"type:uuid;not null" json:"goal_id"`
	Amount      float64   `gorm:"type:numeric(15,2);not null" json:"amount"`
	WithdrawnAt time.Time `gorm:"column:withdrawn_at;default:CURRENT_TIMESTAMP" json:"withdrawn_at"`
	Direction   string    `gorm:"type:varchar(100);not null" json:"direction"`

	Goal *Goal `gorm:"foreignKey:GoalID" json:"goal,omitempty"`
}

// AccountTransfer tracks pure account-to-account velocity
type AccountTransfer struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FromAccountID uuid.UUID `gorm:"column:from_account_id;type:uuid;not null" json:"from_account_id"`
	ToAccountID   uuid.UUID `gorm:"column:to_account_id;type:uuid;not null" json:"to_account_id"`
	Amount        float64   `gorm:"type:numeric(15,2);not null" json:"amount"`
	TransferDate  time.Time `gorm:"column:transfer_date;default:CURRENT_TIMESTAMP" json:"transfer_date"`

	FromAccount *Account `gorm:"foreignKey:FromAccountID" json:"from_account,omitempty"`
	ToAccount   *Account `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"`
}

type BudgetTemplate struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_category_template" json:"user_id"`
	SubCategoryID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_category_template" json:"sub_category_id"`
	AllocatedAmount float64   `gorm:"column:allocated_amount;type:numeric(15,2);not null;default:0.00" json:"allocated_amount"`
	AllowCarryOver  bool      `gorm:"column:allow_carry_over;not null;default:true" json:"allow_carry_over"`
	CreatedAt       time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	SubCategory *SubCategory `gorm:"foreignKey:SubCategoryID" json:"sub_category,omitempty"`
}
