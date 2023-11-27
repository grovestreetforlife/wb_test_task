package psql

import (
	"context"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"testing"
	"wb_test_task/libs/model"
)

func TestUserStorage_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Error(err)
	}
	defer mock.Close()

	storage := newUserStorage(mock)

	testCases := []struct {
		name      string
		mockInput struct {
			id string
		}
		mock           func(id string)
		expectedResult *model.User
		wantErr        bool
		errMsg         string
	}{
		{
			name:      "OK",
			mockInput: struct{ id string }{id: "test"},
			mock: func(id string) {
				mock.ExpectExec("INSERT INTO users").WithArgs(id).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedResult: &model.User{ID: "test"},
			wantErr:        false,
			errMsg:         "",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.mockInput.id)

			user, err := storage.Create(context.Background(), test.mockInput.id)

			if test.wantErr {
				assert.EqualError(t, err, test.errMsg)
				assert.Equal(t, test.expectedResult, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, user)
			}
		})
	}
}

func TestUserStorage_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Error(err)
	}
	defer mock.Close()

	storage := newUserStorage(mock)

	testCases := []struct {
		name      string
		mockInput struct {
			id string
		}
		mock           func(id string)
		expectedResult *model.User
		wantErr        bool
		errMsg         string
	}{
		{
			name:      "OK",
			mockInput: struct{ id string }{id: "test"},
			mock: func(id string) {
				rows := mock.NewRows([]string{"u.id"})
				rows.AddRow(id)
				mock.ExpectQuery("SELECT u.id FROM users u WHERE").WithArgs(id).WillReturnRows(rows)
			},
			expectedResult: &model.User{ID: "test"},
			wantErr:        false,
			errMsg:         "",
		},
		{
			name:      "User does not exists",
			mockInput: struct{ id string }{id: "test"},
			mock: func(id string) {
				rows := mock.NewRows([]string{"u.id"})
				mock.ExpectQuery("SELECT u.id FROM users u WHERE").WithArgs(id).WillReturnRows(rows)
			},
			expectedResult: &model.User{},
			wantErr:        true,
			errMsg:         "no rows in result set",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.mockInput.id)

			user, err := storage.GetByID(context.Background(), test.mockInput.id)

			if test.wantErr {
				assert.EqualError(t, err, test.errMsg)
				assert.Equal(t, test.expectedResult, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, user)
			}
		})
	}
}
