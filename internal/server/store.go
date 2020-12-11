package server

import (
	badger "github.com/dgraph-io/badger/v2"
	jsoniter "github.com/json-iterator/go"
)

type StorageKey byte

const (
	StorageFileKey StorageKey = 'f'
	StorageTextKey StorageKey = 't'
	StorageLinkKey StorageKey = 'l'
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

type ListCallback func(UMField) func([]byte) error

func getAllForType(db *badger.DB, sk StorageKey, cb ListCallback) error {
	pfx := []byte{byte(sk)}
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = true
	opts.Prefix = pfx
	return db.View(func(tx *badger.Txn) error {
		it := tx.NewIterator(opts)
		defer it.Close()
		for it.Seek(pfx); it.ValidForPrefix(pfx); it.Next() {
			u := UMField(it.Item().UserMeta())
			if u.Has(Hidden) {
				continue
			}
			cb2 := cb(u)
			if err := it.Item().Value(cb2); err != nil {
				return err
			}
		}
		return nil
	})
}

func makeMeta(c CommonFields) UMField {
	u := UMField(0)
	if c.BurnAfterRead {
		u = u.Set(BurnAfterRead)
	}
	return u
}

func writeType(db *badger.DB, sk StorageKey, id string, item interface{}, u UMField) error {
	buf, err := jsCfg.Marshal(item)
	if err != nil {
		return err
	}

	key := makeKey(sk, id)
	return db.Update(func(tx *badger.Txn) error {
		return tx.SetEntry(badger.NewEntry(key, buf).WithMeta(byte(u)))
	})
}

func writeBytes(db *badger.DB, sk StorageKey, id string, buf []byte, u UMField) error {
	key := makeKey(sk, id)
	return db.Update(func(tx *badger.Txn) error {
		return tx.SetEntry(badger.NewEntry(key, buf).WithMeta(byte(u)))
	})
}

func _getOne(db *badger.DB, sk StorageKey, id string, cb func([]byte) error) error {
	key := makeKey(sk, id)

	err := db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(cb)
		if err != nil {
			return err
		}
		if UMField(item.UserMeta()).Has(BurnAfterRead) {
			deleteRecord(db, sk, id)
		}
		return nil
	})

	return err
}

func getOneBytes(db *badger.DB, sk StorageKey, id string) ([]byte, error) {
	var buf []byte
	err := _getOne(db, sk, id, func(b []byte) error {
		buf = make([]byte, len(b))
		copy(buf, b)
		return nil
	})
	return buf, err
}

func getOne(db *badger.DB, sk StorageKey, id string, t interface{}) error {
	return _getOne(db, sk, id, func(b []byte) error {
		return jsCfg.Unmarshal(b, t)
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
