package datahandlers

import (
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

type Key struct {
	ID   *uuid.UUID `json:"id" validate:"required"`
	Name string     `json:"name" validate:"required"`
	Key  int32      `json:"key" validate:"required"`
}

var onceKeysHandler sync.Once
var instanceKeysHandler *KeysHandler

type KeysHandler struct {
	keyListMutex sync.Mutex
	keyList      []Key
}

func GetKeysHandlerInstance() *KeysHandler {
	onceKeysHandler.Do(func() {
		instanceKeysHandler = &KeysHandler{}
	})

	return instanceKeysHandler
}

func (mh *KeysHandler) Add(key Key) bool {
	slog.Debug("Adding key: ", "key", key)

	mh.keyListMutex.Lock()
	defer mh.keyListMutex.Unlock()

	return mh.add(key)
}

func (mh *KeysHandler) AddAll(keys []Key) int {
	slog.Debug("Adding key list: ", "key", keys)

	mh.keyListMutex.Lock()
	defer mh.keyListMutex.Unlock()

	return mh.addAll(keys)
}

func (mh *KeysHandler) Remove(key Key) bool {
	slog.Debug("Removing key: ", "key", key)

	mh.keyListMutex.Lock()
	defer mh.keyListMutex.Unlock()

	for i, a := range mh.keyList {
		if a.Name == key.Name {
			mh.keyList = append(mh.keyList[:i], mh.keyList[i+1:]...)
			return true
		}
	}
	return false
}

func (mh *KeysHandler) Clear() {
	slog.Debug("Clearing key list")

	mh.keyListMutex.Lock()
	defer mh.keyListMutex.Unlock()

	mh.keyList = []Key{}
}

func (mh *KeysHandler) GetKeyList() []Key {
	mh.keyListMutex.Lock()
	defer mh.keyListMutex.Unlock()

	return mh.keyList
}

func (mh *KeysHandler) GetByID(id *uuid.UUID) *Key {
	mh.keyListMutex.Lock()
	defer mh.keyListMutex.Unlock()

	if id == nil {
		slog.Warn("ID is nil")
		return nil
	}
	for _, a := range mh.keyList {
		if a.ID == id {
			return &a
		}
	}
	return nil
}

func (mh *KeysHandler) GetByIndex(index int) *Key {
	// mh.keyListMutex.Lock()
	// defer mh.keyListMutex.Unlock()

	if index >= 0 && index < len(mh.keyList) {
		return &mh.keyList[index]
	}
	return nil
}

func (mh *KeysHandler) RemoveByIndex(index int) bool {
	mh.keyListMutex.Lock()
	defer mh.keyListMutex.Unlock()

	if index >= 0 && index < len(mh.keyList) {
		mh.keyList = append(mh.keyList[:index], mh.keyList[index+1:]...)
		return true
	}
	return false
}

func (mh *KeysHandler) Size() int {
	mh.keyListMutex.Lock()
	defer mh.keyListMutex.Unlock()

	return len(mh.keyList)
}

func (mh *KeysHandler) IsEmpty() bool {
	mh.keyListMutex.Lock()
	defer mh.keyListMutex.Unlock()

	return len(mh.keyList) == 0
}

func (mh *KeysHandler) ContainsID(id *uuid.UUID) bool {
	if id == nil {
		slog.Warn("ID is nil")
		return false
	}
	mh.keyListMutex.Lock()
	defer mh.keyListMutex.Unlock()

	return mh.containsID(id)
}

func (mh *KeysHandler) containsID(uuid *uuid.UUID) bool {
	if uuid == nil {
		slog.Warn("ID is nil")
		return false
	}
	for _, a := range mh.keyList {
		if a.ID == uuid {
			return true
		}
	}
	return false
}

func (mh *KeysHandler) add(key Key) bool {
	if key.ID == nil {
		id := uuid.New()
		key.ID = &id
	}

	if !mh.containsID(key.ID) {
		mh.keyList = append(mh.keyList, key)
		return true
	}
	return false
}

func (mh *KeysHandler) addAll(Keys []Key) int {
	addedCnt := 0
	for _, Key := range Keys {
		if mh.add(Key) {
			addedCnt++
		}
	}
	return addedCnt
}

func (kh *KeysHandler) GenerateKeyList() {
	slog.Debug("Generating key list")
	kh.Clear()

	kh.keyListMutex.Lock()
	defer kh.keyListMutex.Unlock()

	for i := 32; i < 127; i++ {
		char := rune(i) // Replace 65 with the ASCII code you want to start from
		kh.add(Key{
			ID:   nil,
			Name: string(char),
			Key:  char,
		})
	}
	slog.Debug("Generated key list: ", "list", len(kh.keyList))
}
