package storage_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"nimbler_writer/internal/storage"
	"nimbler_writer/internal/tests"
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

	// Create New User
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
	}
}
