//тестирование слоя репозитория

package repository

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/user/internal/entities"
	"github.com/stretchr/testify/assert"
)

// Тесты репозитрия с моком вместо RelationalDataBase для тестирования
// Особого смысла в данном конкретном случае для таких тестов нет, поскольку всё,
// Что делает слой репозитория, это вызывает интерфейс relDB и логирует ошибку в случае
// Её возникновения. Эти тесты следует рассматривать как нечто для будущего)

//Мок для интерфейса RelationalDataBase
type MockRelationDB struct{}

func (m MockRelationDB) AddNewUser(ctx context.Context, user *entities.UserInfo, code string) (int, error) {
	if user.Usrname == "" {
		return 0, errors.New("Empty Username")
	}
	return 0, nil
}

func (m MockRelationDB) GetUserById(ctx context.Context, id int) (*entities.UserInfo, error) {
	if id < 0 {
		return nil, errors.New("incorrect Id")
	}
	return nil, nil
}

func (m MockRelationDB) GetUserByEmail(ctx context.Context, email string) (*entities.UserInfo, error) {
	if email == "" {
		return nil, errors.New("Empty Email")
	}
	return nil, nil
}

func (m MockRelationDB) VerifyCode(ctx context.Context, userId int, code string) (bool, error) {

	codeFromDB := "code"

	if userId < 0 {
		return false, errors.New("incorrect Id")
	}
	if code == "" {
		return false, errors.New("Empty Code")
	}
	if code != codeFromDB {
		return false, nil
	}
	return true, nil
}

func (m MockRelationDB) UpdateUser(ctx context.Context, userId int, user *entities.UserInfo) error {
	if user.Usrname == "" {
		return errors.New("Empty Username")
	}
	return nil
}

func TestUserRepository_AddNewUser_CorrectCreditionals(t *testing.T) {

	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	repo := NewUserRepository(MockRelationDB{}, logger)

	userId, err := repo.AddNewUser(context.Background(), &entities.UserInfo{Usrname: "Andrew"}, "code")
	assert.NoError(t, err)
	assert.Equal(t, 0, userId)
}

func TestUserRepository_AddNewUser_IncorrectCreditionals(t *testing.T) {
	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	repo := NewUserRepository(MockRelationDB{}, logger)

	userId, err := repo.AddNewUser(context.Background(), &entities.UserInfo{Usrname: ""}, "code")
	assert.Error(t, err)
	assert.Equal(t, 0, userId)
}

func TestUserRepository_GetUserById_CorrectCreditionals(t *testing.T) {
	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	repo := NewUserRepository(MockRelationDB{}, logger)
	user, err := repo.GetUserById(context.Background(), 1)
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetUserById_IncorrectCreditionals(t *testing.T) {
	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	repo := NewUserRepository(MockRelationDB{}, logger)
	user, err := repo.GetUserById(context.Background(), -1)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetUserByEmail_CorrectCreditionals(t *testing.T) {
	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	repo := NewUserRepository(MockRelationDB{}, logger)

	user, err := repo.GetUserByEmail(context.Background(), "test@test.com")
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetUserByEmail_IncorrectCreditionals(t *testing.T) {
	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	repo := NewUserRepository(MockRelationDB{}, logger)

	user, err := repo.GetUserByEmail(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_UpdateUser_CorrectCreditionals(t *testing.T) {
	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	repo := NewUserRepository(MockRelationDB{}, logger)
	err := repo.UpdateUser(context.Background(), 1, &entities.UserInfo{Usrname: "Andrew"})
	assert.NoError(t, err)
}

func TestUserRepository_UpdateUser_IncorrectCreditionals(t *testing.T) {
	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	repo := NewUserRepository(MockRelationDB{}, logger)
	err := repo.UpdateUser(context.Background(), 1, &entities.UserInfo{Usrname: ""})
	assert.Error(t, err)
}

func TestUserRepository_VerifyCode_IncorrectCode(t *testing.T) {
	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	repo := NewUserRepository(MockRelationDB{}, logger)
	isVErified, err := repo.VerifyCode(context.Background(), 1, "123456")
	assert.NoError(t, err)
	assert.False(t, isVErified)
}

func TestUserRepository_VerifyCode_CorrectCode(t *testing.T) {
	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	repo := NewUserRepository(MockRelationDB{}, logger)
	isVErified, err := repo.VerifyCode(context.Background(), 1, "code")
	assert.NoError(t, err)
	assert.True(t, isVErified)
}

func TestUserRepository_VerifyCode_IncorrectUserId(t *testing.T) {
	var logger *slog.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	repo := NewUserRepository(MockRelationDB{}, logger)
	isVErified, err := repo.VerifyCode(context.Background(), -1, "code")
	assert.Error(t, err)
	assert.False(t, isVErified)
}
