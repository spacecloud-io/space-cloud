package types

import (
	"encoding/json"
)

type DbStore map[string]ColStore
type ColStore map[string]IdStore
type IdStore map[string]*Store

// RealtimeRequest is the object sent for realtime requests
type RealtimeRequest struct {
	Token   string                 `json:"token"`
	DBType  string                 `json:"dbType"`
	Project string                 `json:"project"`
	Group   string                 `json:"group"` // Group is the collection name
	Type    string                 `json:"type"`  // Can either be subscribe or unsubscribe
	ID      string                 `json:"id"`    // id is the query id
	Where   map[string]interface{} `json:"where"`
	Options *LiveQueryOptions      `json:"options"`
}

// LiveQueryOptions is used to set the options for the live query
type LiveQueryOptions struct {
	ChangesOnly bool
	SkipInitial bool `json:"skipInitial"`
}

type RealtimeParams struct {
	Find M
}

type Store struct {
	QueryOptions LiveQueryOptions
	Snapshot     []*SnapshotData
	C            chan *SubscriptionEvent
	Find         interface{}
	Options      interface{}
	Unsubscribe  func()
}

func LiveQuerySubscriptionInit(store *Store) *SubscriptionObject {
	return &SubscriptionObject{store: store}
}

type SubscriptionObject struct {
	store *Store
}

func (s *SubscriptionObject) C() chan *SubscriptionEvent {
	return s.store.C
}

func (s *SubscriptionObject) GetSnapshot() []DocumentSnapshot {
	docs := make([]DocumentSnapshot, 0)
	for _, v := range s.store.Snapshot {
		if !v.IsDeleted {
			docs = append(docs, DocumentSnapshot{doc: v})
		}
	}
	return docs
}

// Unsubscribe
func (s *SubscriptionObject) Unsubscribe() {
	s.store.Unsubscribe()
}

// DocumentSnapshot contains the data and meta info of a single document
type DocumentSnapshot struct {
	doc *SnapshotData
}

// Unmarshal parses the document and stores the value into vPtr
func (s *DocumentSnapshot) Unmarshal(vPtr interface{}) error {
	data, err := json.Marshal(s.doc.Payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, vPtr)
}

type SubscriptionEvent struct {
	err    error
	find   map[string]interface{}
	doc    interface{}
	evType string
}

func NewSubscriptionEvent(evType string, doc interface{}, find map[string]interface{}, err error) *SubscriptionEvent {
	return &SubscriptionEvent{evType: evType, doc: doc, find: find, err: err}
}

func (s *SubscriptionEvent) Unmarshal(vPtr interface{}) error {
	data, err := json.Marshal(s.doc)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, vPtr)
}

func (s *SubscriptionEvent) Type() string {
	return s.evType
}

func (s *SubscriptionEvent) Find() map[string]interface{} {
	return s.find
}

func (s *SubscriptionEvent) Err() error {
	return s.err
}

type SnapshotData struct {
	Find      map[string]interface{}
	Time      int64
	Payload   interface{}
	IsDeleted bool
}
