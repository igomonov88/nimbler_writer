package storage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/igomonov88/nimbler_writer/internal/storage"
	"github.com/igomonov88/nimbler_writer/internal/tests"
)

func TestURLs(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	t.Log("Given the need to test url functionality:")

	u := struct {
		name, email, password string
	}{
		"maris",
		"maris@gmai.com",
		"myPassword",
	}

	nu, err := storage.CreateUser(context.Background(), db, u.name, u.email, u.password)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to add new user to storage: %s", tests.Failed, err)
	}

	// Generating urls for tests
	urls := make([]storage.Url, 30)
	for i := range urls {
		url := storage.Url{
			URLHash:     fmt.Sprintf("%v_string", i),
			UserID:      nu.ID,
			ExpiredAt:   time.Now().Add(5 * time.Millisecond),
			OriginalURL: fmt.Sprintf("original_url_%v", i),
			CustomAlias: fmt.Sprintf("custom_alias_for_url_%v", i),
		}
		urls[i] = url
	}

	t.Log("Given the need to test url functionality:")
	{
		for i := range urls {
			if err := storage.StoreURL(context.Background(), db, urls[i]); err != nil {
				t.Fatalf("\t%s\tShould be able to add new url to storage: %s", tests.Failed, err)
			}
		}
		t.Logf("\t%s\tShould be able to add new url to storage.", tests.Success)
	}

	urlsForDelete := make([]string, 0, len(urls))
	for i := range urls {
		urlsForDelete = append(urlsForDelete, urls[i].URLHash)
	}

	if err := storage.DeleteURLS(context.Background(), db, urlsForDelete); err != nil {
		t.Fatalf("\t%s\tShould be able to batch delete urls from storage: %s", tests.Failed, err)
	}

	t.Logf("\t%s\tShould be able to batch delete urls from storage.", tests.Success)

	url := storage.Url{
		URLHash:     fmt.Sprintf("%v_string", 1),
		UserID:      nu.ID,
		ExpiredAt:   time.Now().Add(5 * time.Millisecond),
		OriginalURL: fmt.Sprintf("original_url_%v", 1),
		CustomAlias: fmt.Sprintf("custom_alias_for_url_%v", 1),
	}

	if err := storage.StoreURL(context.Background(), db, url); err != nil {
		t.Fatalf("\t%s\tShould be able to add new url to storage: %s", tests.Failed, err)
	}

	exist, err := storage.DoesURLAliasExist(context.Background(), db, url.CustomAlias)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get info about custom alias from storage: %s", tests.Failed, err)
	}
	if !exist {
		t.Fatalf("\t%s\tAlias should be existed in storage: %s", tests.Failed, err)
	}

	t.Logf("\t%s\tShould be able to get info about alias existing.", tests.Success)

	if err := storage.DeleteURL(context.Background(), db, url.URLHash); err != nil {
		t.Fatalf("\t%s\tShould be able to delete url from storage: %s", tests.Failed, err)
	}

	t.Logf("\t%s\tShould be able to delete url from storage.", tests.Success)
}

func TestNegativeURL(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	t.Log("Given the need to test negative scenarios of store url functionality:")

	url := storage.Url{
		URLHash:     uuid.New().String(),
		UserID:      uuid.New().String(),
		ExpiredAt:   time.Now().Add(5 * time.Millisecond),
		OriginalURL: uuid.New().String(),
		CustomAlias: uuid.New().String(),
	}

	if err := storage.StoreURL(context.Background(), db, url); err == nil && err != storage.ErrInvalidUserID {
		t.Fatalf("\t%s\tShould got an error when trying to store url with userID not existing in database: %s", tests.Failed, err)
	}
	t.Logf("\t%s\tShould got an error when trying to store url with userID not existing in database.", tests.Success)

	if err := storage.DeleteURL(context.Background(), db, "test"); err != nil {
		t.Fatalf("\t%s\tShould not get an error when trying to delete not existing hash from database: %s", tests.Failed, err)
	}

	t.Logf("\t%s\tShould not get an error when trying to delete not existing hash from database.", tests.Success)

	if err := storage.DeleteURLS(context.Background(), db, []string{}); err != nil {
		t.Fatalf("\t%s\tShould not get an error when trying to use batch delete with empty slice or urls: %s", tests.Failed, err)
	}

	t.Logf("\t%s\tShould not get an error when trying to use batch delete with empty slice or urls.", tests.Success)

	exist, err := storage.DoesURLAliasExist(context.Background(), db, "custom")
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get info about custom alias from storage: %s", tests.Failed, err)
	}
	if exist {
		t.Fatalf("\t%s\tShould return not exist on not existing custom alias from storage.", tests.Failed)
	}

	t.Logf("\t%s\tShould return not exist on not existing custom alias from storage.", tests.Success)
}