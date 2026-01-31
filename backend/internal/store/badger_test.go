package store

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tobi/contracts/backend/internal/model"
)

const testUser = "test-user"

func newTestStore(t *testing.T) *BadgerStore {
	t.Helper()
	s, err := NewBadgerStore(t.TempDir(), slog.Default())
	if err != nil {
		t.Fatalf("NewBadgerStore: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func makeCategory(name string) model.Category {
	now := time.Now().UTC()
	return model.Category{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func makeContract(categoryID uuid.UUID, name string) model.Contract {
	now := time.Now().UTC()
	return model.Contract{
		ID:         uuid.New(),
		CategoryID: categoryID,
		Name:       name,
		StartDate:  "2025-01-01",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Category CRUD

func TestCreateAndGetCategory(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	cat := makeCategory("Insurance")

	if err := s.CreateCategory(ctx, testUser, cat); err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}

	got, err := s.GetCategory(ctx, testUser, cat.ID)
	if err != nil {
		t.Fatalf("GetCategory: %v", err)
	}
	if got.Name != cat.Name {
		t.Errorf("Name = %q, want %q", got.Name, cat.Name)
	}
}

func TestListCategories(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	cats, err := s.ListCategories(ctx, testUser)
	if err != nil {
		t.Fatalf("ListCategories: %v", err)
	}
	if len(cats) != 0 {
		t.Fatalf("expected empty list, got %d", len(cats))
	}

	for _, name := range []string{"A", "B", "C"} {
		if err := s.CreateCategory(ctx, testUser, makeCategory(name)); err != nil {
			t.Fatalf("CreateCategory(%s): %v", name, err)
		}
	}

	cats, err = s.ListCategories(ctx, testUser)
	if err != nil {
		t.Fatalf("ListCategories: %v", err)
	}
	if len(cats) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(cats))
	}
}

func TestUpdateCategory(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	cat := makeCategory("Old")

	if err := s.CreateCategory(ctx, testUser, cat); err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}

	cat.Name = "New"
	cat.UpdatedAt = time.Now().UTC()
	if err := s.UpdateCategory(ctx, testUser, cat); err != nil {
		t.Fatalf("UpdateCategory: %v", err)
	}

	got, err := s.GetCategory(ctx, testUser, cat.ID)
	if err != nil {
		t.Fatalf("GetCategory: %v", err)
	}
	if got.Name != "New" {
		t.Errorf("Name = %q, want %q", got.Name, "New")
	}
}

func TestDeleteCategory(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	cat := makeCategory("ToDelete")

	if err := s.CreateCategory(ctx, testUser, cat); err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}
	if err := s.DeleteCategory(ctx, testUser, cat.ID); err != nil {
		t.Fatalf("DeleteCategory: %v", err)
	}

	_, err := s.GetCategory(ctx, testUser, cat.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetCategory_NotFound(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetCategory(context.Background(), testUser, uuid.New())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateCategory_NotFound(t *testing.T) {
	s := newTestStore(t)
	err := s.UpdateCategory(context.Background(), testUser, makeCategory("Ghost"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteCategory_NotFound(t *testing.T) {
	s := newTestStore(t)
	err := s.DeleteCategory(context.Background(), testUser, uuid.New())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// Contract CRUD

func TestCreateAndGetContract(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	cat := makeCategory("Cat")
	s.CreateCategory(ctx, testUser, cat)

	con := makeContract(cat.ID, "Phone Plan")
	if err := s.CreateContract(ctx, testUser, con); err != nil {
		t.Fatalf("CreateContract: %v", err)
	}

	got, err := s.GetContract(ctx, testUser, con.ID)
	if err != nil {
		t.Fatalf("GetContract: %v", err)
	}
	if got.Name != con.Name {
		t.Errorf("Name = %q, want %q", got.Name, con.Name)
	}
	if got.CategoryID != cat.ID {
		t.Errorf("CategoryID = %s, want %s", got.CategoryID, cat.ID)
	}
}

func TestListContracts(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	cat := makeCategory("Cat")
	s.CreateCategory(ctx, testUser, cat)

	for _, name := range []string{"A", "B"} {
		s.CreateContract(ctx, testUser, makeContract(cat.ID, name))
	}

	all, err := s.ListContracts(ctx, testUser)
	if err != nil {
		t.Fatalf("ListContracts: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 contracts, got %d", len(all))
	}
}

func TestListContractsByCategory(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	cat1 := makeCategory("Cat1")
	cat2 := makeCategory("Cat2")
	s.CreateCategory(ctx, testUser, cat1)
	s.CreateCategory(ctx, testUser, cat2)

	s.CreateContract(ctx, testUser, makeContract(cat1.ID, "C1"))
	s.CreateContract(ctx, testUser, makeContract(cat1.ID, "C2"))
	s.CreateContract(ctx, testUser, makeContract(cat2.ID, "C3"))

	list, err := s.ListContractsByCategory(ctx, testUser, cat1.ID)
	if err != nil {
		t.Fatalf("ListContractsByCategory: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 contracts for cat1, got %d", len(list))
	}

	list, err = s.ListContractsByCategory(ctx, testUser, cat2.ID)
	if err != nil {
		t.Fatalf("ListContractsByCategory: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 contract for cat2, got %d", len(list))
	}
}

func TestUpdateContract(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	cat := makeCategory("Cat")
	s.CreateCategory(ctx, testUser, cat)

	con := makeContract(cat.ID, "Old")
	s.CreateContract(ctx, testUser, con)

	con.Name = "New"
	con.UpdatedAt = time.Now().UTC()
	if err := s.UpdateContract(ctx, testUser, con); err != nil {
		t.Fatalf("UpdateContract: %v", err)
	}

	got, err := s.GetContract(ctx, testUser, con.ID)
	if err != nil {
		t.Fatalf("GetContract: %v", err)
	}
	if got.Name != "New" {
		t.Errorf("Name = %q, want %q", got.Name, "New")
	}
}

func TestDeleteContract(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	cat := makeCategory("Cat")
	s.CreateCategory(ctx, testUser, cat)

	con := makeContract(cat.ID, "ToDelete")
	s.CreateContract(ctx, testUser, con)

	if err := s.DeleteContract(ctx, testUser, con.ID); err != nil {
		t.Fatalf("DeleteContract: %v", err)
	}

	_, err := s.GetContract(ctx, testUser, con.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetContract_NotFound(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetContract(context.Background(), testUser, uuid.New())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateContract_NotFound(t *testing.T) {
	s := newTestStore(t)
	cat := makeCategory("Cat")
	s.CreateCategory(context.Background(), testUser, cat)
	err := s.UpdateContract(context.Background(), testUser, makeContract(cat.ID, "Ghost"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteContract_NotFound(t *testing.T) {
	s := newTestStore(t)
	err := s.DeleteContract(context.Background(), testUser, uuid.New())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// Cascade delete

func TestDeleteCategory_CascadesContracts(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	cat := makeCategory("Cat")
	s.CreateCategory(ctx, testUser, cat)

	con1 := makeContract(cat.ID, "C1")
	con2 := makeContract(cat.ID, "C2")
	s.CreateContract(ctx, testUser, con1)
	s.CreateContract(ctx, testUser, con2)

	if err := s.DeleteCategory(ctx, testUser, cat.ID); err != nil {
		t.Fatalf("DeleteCategory: %v", err)
	}

	// Both contracts should be gone
	for _, id := range []uuid.UUID{con1.ID, con2.ID} {
		_, err := s.GetContract(ctx, testUser, id)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("contract %s: expected ErrNotFound, got %v", id, err)
		}
	}

	// Index should be clean — listing by deleted category returns empty
	list, err := s.ListContractsByCategory(ctx, testUser, cat.ID)
	if err != nil {
		t.Fatalf("ListContractsByCategory: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 contracts after cascade, got %d", len(list))
	}
}

// Index consistency

func TestUpdateContract_CategoryChange_UpdatesIndex(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	cat1 := makeCategory("Cat1")
	cat2 := makeCategory("Cat2")
	s.CreateCategory(ctx, testUser, cat1)
	s.CreateCategory(ctx, testUser, cat2)

	con := makeContract(cat1.ID, "Moveable")
	s.CreateContract(ctx, testUser, con)

	// Move contract to cat2
	con.CategoryID = cat2.ID
	con.UpdatedAt = time.Now().UTC()
	if err := s.UpdateContract(ctx, testUser, con); err != nil {
		t.Fatalf("UpdateContract: %v", err)
	}

	// cat1 should have 0, cat2 should have 1
	list1, _ := s.ListContractsByCategory(ctx, testUser, cat1.ID)
	list2, _ := s.ListContractsByCategory(ctx, testUser, cat2.ID)

	if len(list1) != 0 {
		t.Errorf("cat1 should have 0 contracts, got %d", len(list1))
	}
	if len(list2) != 1 {
		t.Errorf("cat2 should have 1 contract, got %d", len(list2))
	}
}

func TestDeleteContract_CleansIndex(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	cat := makeCategory("Cat")
	s.CreateCategory(ctx, testUser, cat)

	con := makeContract(cat.ID, "C")
	s.CreateContract(ctx, testUser, con)
	s.DeleteContract(ctx, testUser, con.ID)

	list, _ := s.ListContractsByCategory(ctx, testUser, cat.ID)
	if len(list) != 0 {
		t.Errorf("expected 0 contracts after delete, got %d", len(list))
	}
}

// User auth

func TestCreateAndGetUserByEmail(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	user := model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$fakehashfortest",
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	got, err := s.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}
	if got.ID != user.ID {
		t.Errorf("ID = %s, want %s", got.ID, user.ID)
	}
	if got.Email != user.Email {
		t.Errorf("Email = %q, want %q", got.Email, user.Email)
	}
	if got.PasswordHash != user.PasswordHash {
		t.Errorf("PasswordHash = %q, want %q", got.PasswordHash, user.PasswordHash)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	user := model.User{
		ID:           uuid.New(),
		Email:        "dup@example.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now().UTC(),
	}
	if err := s.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	user.ID = uuid.New()
	err := s.CreateUser(ctx, user)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetUserByEmail(context.Background(), "nobody@test.com")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// User isolation

func TestUserIsolation(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	cat := makeCategory("UserA-Cat")
	s.CreateCategory(ctx, "user-a", cat)

	cats, _ := s.ListCategories(ctx, "user-b")
	if len(cats) != 0 {
		t.Errorf("user-b should see 0 categories, got %d", len(cats))
	}

	_, err := s.GetCategory(ctx, "user-b", cat.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("user-b should get ErrNotFound, got %v", err)
	}
}
