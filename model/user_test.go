package model

import (
	"context"
	"testing"

	"myapp/testutil"

	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	return testutil.NewTestDB(t, &User{})
}

func newTestUser() *User {
	return &User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
	}
}

func TestCreate(t *testing.T) {
	repo := NewUserRepository(newTestDB(t))
	user := newTestUser()

	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if user.ID.String() == "" {
		t.Error("expected ID to be set after Create")
	}
}

func TestGetByID(t *testing.T) {
	repo := NewUserRepository(newTestDB(t))
	user := newTestUser()
	_ = repo.Create(context.Background(), user)

	t.Run("found", func(t *testing.T) {
		got, err := repo.GetByID(context.Background(), user.ID.String())
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if got.Email != user.Email {
			t.Errorf("got email %q, want %q", got.Email, user.Email)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
		if err == nil {
			t.Error("expected error for missing ID, got nil")
		}
	})
}

func TestGetByEmail(t *testing.T) {
	repo := NewUserRepository(newTestDB(t))
	user := newTestUser()
	_ = repo.Create(context.Background(), user)

	t.Run("found", func(t *testing.T) {
		got, err := repo.GetByEmail(context.Background(), user.Email)
		if err != nil {
			t.Fatalf("GetByEmail failed: %v", err)
		}
		if got.ID != user.ID {
			t.Errorf("got ID %v, want %v", got.ID, user.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByEmail(context.Background(), "nobody@example.com")
		if err == nil {
			t.Error("expected error for missing email, got nil")
		}
	})
}

func TestUpdate(t *testing.T) {
	repo := NewUserRepository(newTestDB(t))
	user := newTestUser()
	_ = repo.Create(context.Background(), user)

	user.Name = "Updated Name"
	if err := repo.Update(context.Background(), user); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, _ := repo.GetByID(context.Background(), user.ID.String())
	if got.Name != "Updated Name" {
		t.Errorf("got name %q, want %q", got.Name, "Updated Name")
	}
}

func TestDelete(t *testing.T) {
	repo := NewUserRepository(newTestDB(t))
	user := newTestUser()
	_ = repo.Create(context.Background(), user)

	t.Run("success", func(t *testing.T) {
		if err := repo.Delete(context.Background(), user.ID.String()); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
		_, err := repo.GetByID(context.Background(), user.ID.String())
		if err == nil {
			t.Error("expected error after delete, got nil")
		}
	})

	t.Run("already deleted", func(t *testing.T) {
		err := repo.Delete(context.Background(), user.ID.String())
		if err == nil {
			t.Error("expected error for already-deleted user, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(context.Background(), "00000000-0000-0000-0000-000000000000")
		if err == nil {
			t.Error("expected error for missing ID, got nil")
		}
	})
}

func TestExistsByEmail(t *testing.T) {
	repo := NewUserRepository(newTestDB(t))
	user := newTestUser()
	_ = repo.Create(context.Background(), user)

	t.Run("exists", func(t *testing.T) {
		ok, err := repo.ExistsByEmail(context.Background(), user.Email)
		if err != nil {
			t.Fatalf("ExistsByEmail failed: %v", err)
		}
		if !ok {
			t.Error("expected true for existing email, got false")
		}
	})

	t.Run("not exists", func(t *testing.T) {
		ok, err := repo.ExistsByEmail(context.Background(), "nobody@example.com")
		if err != nil {
			t.Fatalf("ExistsByEmail failed: %v", err)
		}
		if ok {
			t.Error("expected false for missing email, got true")
		}
	})

	t.Run("not exists after delete", func(t *testing.T) {
		_ = repo.Delete(context.Background(), user.ID.String())
		ok, err := repo.ExistsByEmail(context.Background(), user.Email)
		if err != nil {
			t.Fatalf("ExistsByEmail failed: %v", err)
		}
		if ok {
			t.Error("expected false for soft-deleted user, got true")
		}
	})
}
