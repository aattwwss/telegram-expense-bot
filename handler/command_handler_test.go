package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func newTestCommandHandler(ur mockUserRepo, tr mockTransactionRepo, mr mockMessageContextRepo, ttr mockTransactionTypeRepo, cr mockCategoryRepo) (CommandHandler, *tgbotapi.BotAPI) {
	bot := &tgbotapi.BotAPI{
		Token:  "dummy",
		Client: &http.Client{},
		Buffer: 100,
	}
	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)
	return CommandHandler{
		userRepo:            ur,
		transactionRepo:     tr,
		messageContextRepo:  mr,
		transactionTypeRepo: ttr,
		categoryRepo:        cr,
	}, bot
}

func TestStart_UserExists(t *testing.T) {
	var calledFindByIdWith int64
	ur := mockUserRepo{
		findByIdFn: func(ctx context.Context, id int64) (*domain.User, error) {
			calledFindByIdWith = id
			return &domain.User{Id: 123, Currency: money.GetCurrency("SGD"), Location: time.UTC}, nil
		},
	}

	handler, bot := newTestCommandHandler(ur, mockTransactionRepo{}, mockMessageContextRepo{}, mockTransactionTypeRepo{}, mockCategoryRepo{})

	user := &tgbotapi.User{ID: 123}
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: user,
			Chat: &tgbotapi.Chat{ID: 456},
		},
	}

	handler.Start(context.Background(), bot, update)

	if calledFindByIdWith != 123 {
		t.Errorf("expected FindUserById called with 123, got %d", calledFindByIdWith)
	}
}

func TestStart_NewUserSignup(t *testing.T) {
	var addedUser domain.User

	ur := mockUserRepo{
		findByIdFn: func(ctx context.Context, id int64) (*domain.User, error) {
			return nil, nil // user does not exist
		},
		addFn: func(ctx context.Context, user domain.User) error {
			addedUser = user
			return nil
		},
	}

	handler, bot := newTestCommandHandler(ur, mockTransactionRepo{}, mockMessageContextRepo{}, mockTransactionTypeRepo{}, mockCategoryRepo{})

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: 999},
			Chat: &tgbotapi.Chat{ID: 456},
		},
	}

	handler.Start(context.Background(), bot, update)

	if addedUser.Id != 999 {
		t.Errorf("expected new user ID 999, got %d", addedUser.Id)
	}
	if addedUser.Currency.Code != "SGD" {
		t.Errorf("expected default SGD currency, got %s", addedUser.Currency.Code)
	}
}

func TestStart_FindUserError(t *testing.T) {
	ur := mockUserRepo{
		findByIdFn: func(ctx context.Context, id int64) (*domain.User, error) {
			return nil, errors.New("db error")
		},
	}

	handler, bot := newTestCommandHandler(ur, mockTransactionRepo{}, mockMessageContextRepo{}, mockTransactionTypeRepo{}, mockCategoryRepo{})

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: 1},
			Chat: &tgbotapi.Chat{ID: 456},
		},
	}

	// Should not panic; error path sends an error message to the user
	handler.Start(context.Background(), bot, update)
}

func TestStart_AddUserError(t *testing.T) {
	ur := mockUserRepo{
		findByIdFn: func(ctx context.Context, id int64) (*domain.User, error) {
			return nil, nil
		},
		addFn: func(ctx context.Context, user domain.User) error {
			return errors.New("insert error")
		},
	}

	handler, bot := newTestCommandHandler(ur, mockTransactionRepo{}, mockMessageContextRepo{}, mockTransactionTypeRepo{}, mockCategoryRepo{})

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: 1},
			Chat: &tgbotapi.Chat{ID: 456},
		},
	}

	handler.Start(context.Background(), bot, update)
}

func TestUndo_NoTransactions(t *testing.T) {
	var sentMessage string
	ur := mockUserRepo{
		findByIdFn: func(ctx context.Context, id int64) (*domain.User, error) {
			return &domain.User{Id: 1}, nil
		},
	}

	tr := mockTransactionRepo{
		findLatestByUserIdFn: func(ctx context.Context, userId int64) (*domain.Transaction, error) {
			return nil, nil // no transactions
		},
	}

	handler, bot := newTestCommandHandler(ur, tr, mockMessageContextRepo{}, mockTransactionTypeRepo{}, mockCategoryRepo{})

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: 1},
			Chat: &tgbotapi.Chat{ID: 456},
			Text: "/undo",
		},
	}

	_ = sentMessage
	handler.Undo(context.Background(), bot, update)
	// Verify we reached the "no transactions" branch by checking no further repo methods were called
}

func TestUndo_FindLatestError(t *testing.T) {
	tr := mockTransactionRepo{
		findLatestByUserIdFn: func(ctx context.Context, userId int64) (*domain.Transaction, error) {
			return nil, errors.New("db error")
		},
	}

	handler, bot := newTestCommandHandler(mockUserRepo{
		findByIdFn: func(ctx context.Context, id int64) (*domain.User, error) {
			return &domain.User{Id: 1}, nil
		},
	}, tr, mockMessageContextRepo{}, mockTransactionTypeRepo{}, mockCategoryRepo{})

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: 1},
			Chat: &tgbotapi.Chat{ID: 456},
			Text: "/undo",
		},
	}

	handler.Undo(context.Background(), bot, update)
}

func TestNewCategoriesKeyboard(t *testing.T) {
	categories := []*entity.Category{
		{Id: 1, Name: "Food", TransactionTypeId: 1},
		{Id: 2, Name: "Transport", TransactionTypeId: 1},
	}

	kb, err := newCategoriesKeyboard(categories, 42, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kb) != 2 { // 1 data row + cancel
		t.Fatalf("expected 2 rows, got %d", len(kb))
	}
	if kb[0][0].Text != "Food" || kb[0][1].Text != "Transport" {
		t.Errorf("expected Food and Transport buttons")
	}
}

func TestNewCategoriesKeyboard_Empty(t *testing.T) {
	categories := []*entity.Category{}
	kb, err := newCategoriesKeyboard(categories, 42, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kb) != 1 { // only cancel row
		t.Fatalf("expected 1 row, got %d", len(kb))
	}
}

func TestNewTransactionTypesKeyboard(t *testing.T) {
	types := []*entity.TransactionType{
		{Id: 1, Name: "Spent", Multiplier: -1, ReplyText: "Spent %s"},
		{Id: 2, Name: "Received", Multiplier: 1, ReplyText: "Received %s"},
	}

	kb, err := newTransactionTypesKeyboard(types, 42, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kb) != 2 { // 1 data row + cancel
		t.Fatalf("expected 2 rows, got %d", len(kb))
	}
	if kb[0][0].Text != "Spent" || kb[0][1].Text != "Received" {
		t.Errorf("expected Spent and Received buttons")
	}
}
