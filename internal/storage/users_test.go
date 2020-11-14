package storage_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/igomonov88/nimbler_writer/internal/storage"
	"github.com/igomonov88/nimbler_writer/internal/tests"
)

func TestUser(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	t.Log("Given the need to test user functionality:")

	u := struct {
		email    string
		name     string
		password string
	}{
		"igor@gmail.com",
		"igor",
		"qwerty",
	}

	// Operate with user
	{
		nu, err := storage.CreateUser(context.Background(), db, u.name, u.email, u.password)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to add new user to storage: %s", tests.Failed, err)
		}

		if cmp.Diff(nu.Name, u.name) != "" {
			t.Fatalf("\t%s\tCreated user should have same name as was provided.", tests.Failed)
		}

		if cmp.Diff(nu.Email, u.email) != "" {
			t.Fatalf("\t%s\tCreated user should have same email as was provided.", tests.Failed)
		}

		if len(nu.Password) == 0 {
			t.Fatalf("\t%s\tCreated user should have same email as was provided.", tests.Failed)
		}

		t.Logf("\t%s\tShould be able to add new user to storage.", tests.Success)

		_, err = storage.Authenticate(context.Background(), db, u.email, u.password)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to authenticate user: %s", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to authenticate user.", tests.Success)

		nui := struct {
			name, email string
		}{
			"maris",
			"maris@gmai.com",
		}

		err = storage.UpdateUserInfo(context.Background(), db, nu.ID, nui.name, nui.email)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to update user info: %s", tests.Failed, err)
		}

		ru, err := storage.RetrieveUser(context.Background(), db, nu.ID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve user: %s", tests.Failed, err)
		}

		if cmp.Diff(ru.Name, nui.name) != "" || cmp.Diff(ru.Email, nui.email) != "" {
			t.Fatalf("\t%s\tShould update user info: %s", tests.Failed, err)
		}
		t.Logf("\t%s\tShould update user info storage.", tests.Success)

		err = storage.UpdateUsersPassword(context.Background(), db, nu.ID, "qazwsxedc")
		if err != nil {
			t.Fatalf("\t%s\tShould update users password: %s", tests.Failed, err)
		}
		t.Logf("\t%s\tShould update user users password.", tests.Success)

		exist, err := storage.DoesUserEmailExist(context.Background(), db, nui.email)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to check email exist: %s", tests.Failed, err)
		}
		if !exist {
			t.Fatalf("\t%s\tShould exist email for previosly created user: %s", tests.Failed, nui.email)
		}
		t.Logf("\t%s\tShould be able to check email exist.", tests.Success)

		err = storage.DeleteUser(context.Background(), db, ru.ID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to delete user with userID: %s", tests.Failed, err)
		}

		_, err = storage.RetrieveUser(context.Background(), db, ru.ID)
		if err != storage.ErrNotFound {
			t.Fatalf("\t%s\tShould not return deleted user: %s", tests.Failed, err)
		}

		t.Logf("\t%s\tShould be able to delete user with userID.", tests.Success)
	}
}

func TestNegativeScenariosForUser(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	t.Log("Given the need to test user negative scenarios:")

	u := struct {
		email    string
		name     string
		password string
	}{
		"igor@gmail.com",
		"igor",
		"qwerty",
	}
	{
		_, err := storage.CreateUser(context.Background(), db, u.name, u.email, u.password)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to add new user to storage: %s", tests.Failed, err)
		}

		_, err = storage.Authenticate(context.Background(), db, u.email, "123456")
		if err == nil {
			t.Fatalf("\t%s\tShould return %s, when user provide not valid password: %s", tests.Failed, storage.ErrAuthenticationFailure, err)
		}
		if err != storage.ErrAuthenticationFailure {
			t.Fatalf("\t%s\tShould return %s, when user provide not valid password: %s", tests.Failed, storage.ErrAuthenticationFailure, err)
		}
		t.Logf("\t%s\tShould return %s, when user provide not valid password", tests.Success, storage.ErrEmailAlreadyExist)

		_, err = storage.Authenticate(context.Background(), db, "qmail1@gmail.com", "qwerty")
		if err == nil {
			t.Fatalf("\t%s\tShould return %s, when user provide not valid password: %s", tests.Failed, storage.ErrAuthenticationFailure, err)
		}
		if err != storage.ErrNotFound {
			t.Fatalf("\t%s\tShould return %s, when user provide not existing email: %s", tests.Failed, storage.ErrNotFound, err)
		}
		t.Logf("\t%s\tShould return %s, when user provide not existing email.", tests.Success, storage.ErrNotFound)

		_, err = storage.CreateUser(context.Background(), db, u.name, u.email, u.password)
		if err != storage.ErrEmailAlreadyExist {
			t.Fatalf("\t%s\tShould return %s, when user with such email already exist: %s", tests.Failed, storage.ErrEmailAlreadyExist, err)
		}
		t.Logf("\t%s\tShould return %s when user with such email already exist.", tests.Success, storage.ErrEmailAlreadyExist)

		// try to retrieve user with not valid uuid as user_id
		_, err = storage.RetrieveUser(context.Background(), db, "qwerasdzxc")
		if err != storage.ErrInvalidUserID {
			t.Fatalf("\t%s\tShould return %s while using invalid uuid when retrieve user: %s", tests.Failed, storage.ErrInvalidUserID, err)
		}
		t.Logf("\t%s\tShould return %s while using invalid uuid when retrieve user.", tests.Success, storage.ErrInvalidUserID)

		// try to retrieve user with not existing user_id
		_, err = storage.RetrieveUser(context.Background(), db, uuid.New().String())
		if err != storage.ErrNotFound {
			t.Fatalf("\t%s\tShould return %s while using not existing uuid: %s", tests.Failed, storage.ErrNotFound, err)
		}
		t.Logf("\t%s\tShould return %s while using not existing uuid.", tests.Success, storage.ErrNotFound)

		// try to update user's password with not valid uuid as user_id
		err = storage.UpdateUsersPassword(context.Background(), db, "qweasdzxc", "qwerty")
		if err != storage.ErrInvalidUserID {
			t.Fatalf("\t%s\tShould return %s while using invalid uuid when update password: %s", tests.Failed, storage.ErrInvalidUserID, err)
		}
		t.Logf("\t%s\tShould return %s while using invalid uuid when update password.", tests.Success, storage.ErrInvalidUserID)

		// try to delete user with not valid uuid as user_id
		err = storage.DeleteUser(context.Background(),db, "qweqweqwe")
		if err != storage.ErrInvalidUserID {
			t.Fatalf("\t%s\tShould return %s while using invalid uuid when deleting user: %s", tests.Failed, storage.ErrInvalidUserID, err)
		}
		t.Logf("\t%s\tShould return %s while using invalid uuid when delete user.", tests.Success, storage.ErrInvalidUserID)

		err = storage.UpdateUserInfo(context.Background(), db, "qweqweqwe", "qweqwe", "qweqwe")
		if err != storage.ErrInvalidUserID {
			t.Fatalf("\t%s\tShould return %s while using invalid uuid when updating user info: %s", tests.Failed, storage.ErrInvalidUserID, err)
		}
		t.Logf("\t%s\tShould return %s while using invalid uuid when updating user info.", tests.Success, storage.ErrInvalidUserID)
	}
}
