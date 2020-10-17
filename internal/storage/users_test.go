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

	t.Log("Given the need to test user functionality,")

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
		t.Log("And we try to create user:")

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
	}

}
