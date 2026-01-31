package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"github.com/tobi/contracts/backend/internal/model"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type BadgerStore struct {
	db     *badger.DB
	logger *slog.Logger
	done   chan struct{}
}

type badgerLogger struct {
	logger *slog.Logger
}

func (l *badgerLogger) Errorf(f string, v ...interface{})   { l.logger.Error(fmt.Sprintf(f, v...)) }
func (l *badgerLogger) Warningf(f string, v ...interface{}) { l.logger.Warn(fmt.Sprintf(f, v...)) }
func (l *badgerLogger) Infof(f string, v ...interface{})    { l.logger.Info(fmt.Sprintf(f, v...)) }
func (l *badgerLogger) Debugf(f string, v ...interface{})   { l.logger.Debug(fmt.Sprintf(f, v...)) }

func NewBadgerStore(path string, logger *slog.Logger) (*BadgerStore, error) {
	opts := badger.DefaultOptions(path).
		WithLogger(&badgerLogger{logger: logger.With("component", "badger")})

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("opening badger db: %w", err)
	}

	s := &BadgerStore{
		db:     db,
		logger: logger,
		done:   make(chan struct{}),
	}
	go s.runGC()
	return s, nil
}

func (s *BadgerStore) runGC() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			for s.db.RunValueLogGC(0.5) == nil {
				// keep running until no more GC needed
			}
		}
	}
}

func (s *BadgerStore) Close() error {
	close(s.done)
	return s.db.Close()
}

func (s *BadgerStore) Healthy() error {
	return s.db.View(func(txn *badger.Txn) error {
		return nil
	})
}

// User keys

func usrKey(id uuid.UUID) []byte {
	return []byte(fmt.Sprintf("usr/%s", id))
}

func usrEmailKey(email string) []byte {
	return []byte(fmt.Sprintf("usr_email/%s", email))
}

// Users

func (s *BadgerStore) CreateUser(_ context.Context, u model.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return s.db.Update(func(txn *badger.Txn) error {
		// Check email uniqueness
		if _, err := txn.Get(usrEmailKey(u.Email)); err == nil {
			return ErrConflict
		}
		if err := txn.Set(usrKey(u.ID), data); err != nil {
			return err
		}
		return txn.Set(usrEmailKey(u.Email), []byte(u.ID.String()))
	})
}

func (s *BadgerStore) GetUserByEmail(_ context.Context, email string) (model.User, error) {
	var user model.User
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(usrEmailKey(email))
		if err != nil {
			return err
		}
		var idStr string
		if err := item.Value(func(val []byte) error {
			idStr = string(val)
			return nil
		}); err != nil {
			return err
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return err
		}
		uItem, err := txn.Get(usrKey(id))
		if err != nil {
			return err
		}
		return uItem.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		return user, ErrNotFound
	}
	return user, err
}

// Key helpers

func catKey(userID string, id uuid.UUID) []byte {
	return []byte(fmt.Sprintf("u/%s/cat/%s", userID, id))
}

func catPrefix(userID string) []byte {
	return []byte(fmt.Sprintf("u/%s/cat/", userID))
}

func conKey(userID string, id uuid.UUID) []byte {
	return []byte(fmt.Sprintf("u/%s/con/%s", userID, id))
}

func conPrefix(userID string) []byte {
	return []byte(fmt.Sprintf("u/%s/con/", userID))
}

func idxCatConKey(userID string, categoryID, contractID uuid.UUID) []byte {
	return []byte(fmt.Sprintf("u/%s/idx/cat_con/%s/%s", userID, categoryID, contractID))
}

func idxCatConPrefix(userID string, categoryID uuid.UUID) []byte {
	return []byte(fmt.Sprintf("u/%s/idx/cat_con/%s/", userID, categoryID))
}

// Categories

func (s *BadgerStore) ListCategories(_ context.Context, userID string) ([]model.Category, error) {
	var categories []model.Category
	prefix := catPrefix(userID)

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var cat model.Category
			if err := it.Item().Value(func(val []byte) error {
				return json.Unmarshal(val, &cat)
			}); err != nil {
				return err
			}
			categories = append(categories, cat)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if categories == nil {
		categories = []model.Category{}
	}
	return categories, nil
}

func (s *BadgerStore) GetCategory(_ context.Context, userID string, id uuid.UUID) (model.Category, error) {
	var cat model.Category
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(catKey(userID, id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &cat)
		})
	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		return cat, ErrNotFound
	}
	return cat, err
}

func (s *BadgerStore) CreateCategory(_ context.Context, userID string, c model.Category) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(catKey(userID, c.ID), data)
	})
}

func (s *BadgerStore) UpdateCategory(_ context.Context, userID string, c model.Category) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return s.db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(catKey(userID, c.ID)); err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrNotFound
			}
			return err
		}
		return txn.Set(catKey(userID, c.ID), data)
	})
}

func (s *BadgerStore) DeleteCategory(_ context.Context, userID string, id uuid.UUID) error {
	return s.db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(catKey(userID, id)); err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrNotFound
			}
			return err
		}

		if err := txn.Delete(catKey(userID, id)); err != nil {
			return err
		}

		// Delete all contracts in this category via the index
		prefix := idxCatConPrefix(userID, id)
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		var contractIDs []uuid.UUID
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().Key()
			// Extract contract ID from the end of the index key
			conIDStr := string(key[len(prefix):])
			conID, err := uuid.Parse(conIDStr)
			if err != nil {
				continue
			}
			contractIDs = append(contractIDs, conID)
		}
		it.Close()

		for _, conID := range contractIDs {
			if err := txn.Delete(conKey(userID, conID)); err != nil {
				return err
			}
			if err := txn.Delete(idxCatConKey(userID, id, conID)); err != nil {
				return err
			}
		}

		return nil
	})
}

// Contracts

func (s *BadgerStore) ListContracts(_ context.Context, userID string) ([]model.Contract, error) {
	var contracts []model.Contract
	prefix := conPrefix(userID)

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var con model.Contract
			if err := it.Item().Value(func(val []byte) error {
				return json.Unmarshal(val, &con)
			}); err != nil {
				return err
			}
			contracts = append(contracts, con)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if contracts == nil {
		contracts = []model.Contract{}
	}
	return contracts, nil
}

func (s *BadgerStore) ListContractsByCategory(_ context.Context, userID string, categoryID uuid.UUID) ([]model.Contract, error) {
	var contracts []model.Contract
	prefix := idxCatConPrefix(userID, categoryID)

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().Key()
			conIDStr := string(key[len(prefix):])
			conID, err := uuid.Parse(conIDStr)
			if err != nil {
				continue
			}

			item, err := txn.Get(conKey(userID, conID))
			if err != nil {
				if errors.Is(err, badger.ErrKeyNotFound) {
					continue
				}
				return err
			}

			var con model.Contract
			if err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &con)
			}); err != nil {
				return err
			}
			contracts = append(contracts, con)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if contracts == nil {
		contracts = []model.Contract{}
	}
	return contracts, nil
}

func (s *BadgerStore) GetContract(_ context.Context, userID string, id uuid.UUID) (model.Contract, error) {
	var con model.Contract
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(conKey(userID, id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &con)
		})
	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		return con, ErrNotFound
	}
	return con, err
}

func (s *BadgerStore) CreateContract(_ context.Context, userID string, c model.Contract) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return s.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(conKey(userID, c.ID), data); err != nil {
			return err
		}
		return txn.Set(idxCatConKey(userID, c.CategoryID, c.ID), []byte{})
	})
}

func (s *BadgerStore) UpdateContract(_ context.Context, userID string, c model.Contract) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return s.db.Update(func(txn *badger.Txn) error {
		// Get existing to check old categoryID
		item, err := txn.Get(conKey(userID, c.ID))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrNotFound
			}
			return err
		}

		var old model.Contract
		if err := item.Value(func(val []byte) error {
			return json.Unmarshal(val, &old)
		}); err != nil {
			return err
		}

		if err := txn.Set(conKey(userID, c.ID), data); err != nil {
			return err
		}

		// Update index if category changed
		if old.CategoryID != c.CategoryID {
			if err := txn.Delete(idxCatConKey(userID, old.CategoryID, c.ID)); err != nil {
				return err
			}
			if err := txn.Set(idxCatConKey(userID, c.CategoryID, c.ID), []byte{}); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *BadgerStore) DeleteContract(_ context.Context, userID string, id uuid.UUID) error {
	return s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(conKey(userID, id))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrNotFound
			}
			return err
		}

		var con model.Contract
		if err := item.Value(func(val []byte) error {
			return json.Unmarshal(val, &con)
		}); err != nil {
			return err
		}

		if err := txn.Delete(conKey(userID, id)); err != nil {
			return err
		}
		return txn.Delete(idxCatConKey(userID, con.CategoryID, id))
	})
}
