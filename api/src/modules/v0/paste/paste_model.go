//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description: --
package paste

import (
	"../common"
	"errors"
	"time"
)

type Note struct {
	ID      string    `json:"_id" valid:"-"`
	UserID  string    `json:"user_id,omitempty" valid:"-"`
	Note    string    `json:"note" valid:"required"`
	ViewCnt int       `json:"view_cnt,omitempty" valid:"-"`
	TTL     int       `json:"ttl,omitempty" valid:"-"`
	Created time.Time `json:"created" valid:"-"`
	Updated time.Time `json:"updated" valid:"-"`
}

type Model struct {
}

func NewModel() *Model {
	return new(Model)
}

//NewNote - shortcut for populating note struct
func NewNote(body string) *Note {
	return &Note{
		ID:      common.NewUUID(),
		UserID:  "-",
		Note:    body,
		TTL:     3600,
		Created: time.Now(),
		Updated: time.Now(),
	}
}

func InitNoteDefaults(n *Note) {
	n.ID = common.NewUUID()
	n.TTL = 3600
	n.Created = time.Now()
	n.Updated = time.Now()
}

//StoreEntry - create new entry
func (m *Model) StoreEntry(n *Note) (*Note, error) {
	rc := common.GetRedisClient()
	nb, err := common.Encode(n)
	if err != nil {
		return nil, errors.New("Could not encode entry")
	}

	if !rc.SetKey(n.ID, nb, n.TTL) {
		return nil, errors.New("Could not store entry")
	}

	return n, nil
}

func (m *Model) GetEntry(id string) (*Note, error) {
	rc := common.GetRedisClient()
	nb := rc.GetKey(id)

	if len(nb) == 0 {
		return nil, errors.New("Entry not found")
	}

	var n *Note
	if err := common.Decode(nb, n); err != nil {
		return nil, errors.New("Serialization error")
	}

	return n, nil
}
