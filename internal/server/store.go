package server

import (
	"time"

	badger "github.com/dgraph-io/badger/v2"
	jsoniter "github.com/json-iterator/go"
)

type StorageKey byte

const (
	StorageFileGroupKey StorageKey = 'g'
	StorageFileKey      StorageKey = 'f'
	StorageTextKey      StorageKey = 't'
	StorageLinkKey      StorageKey = 'l'
)

var jsCfg = jsoniter.Config{
	EscapeHTML:                    true,
	SortMapKeys:                   false,
	ValidateJsonRawMessage:        false,
	MarshalFloatWith6Digits:       true,
	ObjectFieldMustBeSimpleString: false,
}.Froze()

// UserMeta flags

type UMField byte

const (
	BurnAfterRead UMField = 1 << iota
	Hidden
	// ...
)

func (u UMField) Set(flag UMField) UMField    { return u | flag }
func (u UMField) Clear(flag UMField) UMField  { return u &^ flag }
func (u UMField) Toggle(flag UMField) UMField { return u ^ flag }
func (u UMField) Has(flag UMField) bool       { return u&flag != 0 }

func getAllForType(db *badger.DB, sk StorageKey) ([][]byte, error) {
	total := make([][]byte, 0, 20)

	pfx := []byte{byte(sk)}
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = true
	opts.Prefix = pfx
	err := db.View(func(tx *badger.Txn) error {
		it := tx.NewIterator(opts)
		defer it.Close()
		for it.Seek(pfx); it.ValidForPrefix(pfx); it.Next() {
			u := UMField(it.Item().UserMeta())
			if u.Has(Hidden) {
				continue
			}

			if err := it.Item().Value(func(v []byte) error {
				var buf []byte
				var err error
				if u.Has(BurnAfterRead) {
					buf, err = censorPreventBurn(sk, v)
					if err != nil {
						return err
					}
				} else {
					buf = make([]byte, len(v))
					copy(buf, v)
				}
				total = append(total, buf)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	})

	return total, err
}

func makeMeta(c CommonFields) UMField {
	u := UMField(0)
	if c.BurnAfterRead {
		u = u.Set(BurnAfterRead)
	}
	return u
}

func writeType(db *badger.DB, sk StorageKey, id string, item interface{}, u UMField, ttl int64) error {
	buf, err := jsCfg.Marshal(item)
	if err != nil {
		return err
	}
	return writeBytes(db, sk, id, buf, u, ttl)
}

func writeBytes(db *badger.DB, sk StorageKey, id string, buf []byte, u UMField, ttl int64) error {
	key := makeKey(sk, id)
	entry := badger.NewEntry(key, buf).WithMeta(byte(u))
	if ttl > 0 {
		entry = entry.WithTTL(time.Duration(ttl) * time.Second)
	}
	return db.Update(func(tx *badger.Txn) error {
		return tx.SetEntry(entry)
	})
}

type FetchOpts struct {
	SkipBurn bool // does not burn item on read
}

var DontBurn = &FetchOpts{SkipBurn: true}

func _getOne(db *badger.DB, f *FetchOpts, sk StorageKey, id string, cb func([]byte) error) error {
	key := makeKey(sk, id)

	if f == nil {
		f = &FetchOpts{}
	}

	err := db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(cb)
		if err != nil {
			return err
		}
		if UMField(item.UserMeta()).Has(BurnAfterRead) && !f.SkipBurn {
			deleteRecord(db, sk, id)
		}
		return nil
	})

	return err
}

func getOneBytes(db *badger.DB, f *FetchOpts, sk StorageKey, id string) ([]byte, error) {
	var buf []byte
	err := _getOne(db, f, sk, id, func(b []byte) error {
		buf = make([]byte, len(b))
		copy(buf, b)
		return nil
	})
	return buf, err
}

func getOne(db *badger.DB, f *FetchOpts, sk StorageKey, id string, t interface{}) error {
	return _getOne(db, f, sk, id, func(b []byte) error {
		return jsCfg.Unmarshal(b, t)
	})
}

func getMeta(db *badger.DB, sk StorageKey, id string) (UMField, error) {
	var u UMField
	key := makeKey(sk, id)
	return u, db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(key)
		if err != nil {
			return err
		}
		u = UMField(item.UserMeta())
		return nil
	})
}

// deletes a single record
func deleteRecord(db *badger.DB, sk StorageKey, id string) error {
	key := makeKey(sk, id)
	return db.Update(func(tx *badger.Txn) error {
		return tx.Delete(key)
	})
}

// composes the indexing Key
func makeKey(sk StorageKey, id string) []byte {
	key := []byte(id)
	key = append([]byte{byte(sk)}, key...)
	return key
}

type Info struct {
	ID   string
	Meta UMField
}

func listInfo(db *badger.DB, sk StorageKey) ([]Info, error) {
	total := make([]Info, 0, 20)

	pfx := []byte{byte(sk)}
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = pfx
	err := db.View(func(tx *badger.Txn) error {
		it := tx.NewIterator(opts)
		defer it.Close()
		for it.Seek(pfx); it.ValidForPrefix(pfx); it.Next() {
			total = append(total, Info{
				ID:   string(it.Item().KeyCopy(nil)[1:]),
				Meta: UMField(it.Item().UserMeta()),
			})
		}
		return nil
	})

	return total, err
}
