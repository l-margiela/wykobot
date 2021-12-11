package dedup

import (
	"fmt"
	"io"
	"strconv"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

type Deduper interface {
	io.Closer

	// Sent checks if given entry has already been sent.
	Sent(int) bool

	// SetSent marks the entry as sent and remembers it for the given duration.
	SetSent(int, time.Duration) error
}

type Badger struct {
	db *badger.DB
}

func NewBadger(filepath string) (Badger, error) {
	db, err := badger.Open(badger.DefaultOptions(filepath))
	if err != nil {
		return Badger{}, err
	}

	return Badger{db}, nil
}

func (b Badger) Sent(entryID int) bool {
	if err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(strconv.Itoa(entryID)))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return false
	}
	return true
}

func (b Badger) SetSent(entryID int, ttl time.Duration) error {
	txn := b.db.NewTransaction(true)

	entry := badger.NewEntry([]byte(strconv.Itoa(entryID)), []byte("true")).WithTTL(ttl)
	if err := txn.SetEntry(entry); err != nil {
		fmt.Printf("err setting entry for %d: %s\n", entryID, err)
		return err
	}

	fmt.Printf("commit %d\n", entryID)
	return txn.Commit()
}

func (b Badger) Close() error {
	return b.db.Close()
}
