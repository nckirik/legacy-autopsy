package work

import (
	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

// ReadSet binds every record an invocation observed to the fingerprint at read
// time. It is the invocation's authoritative dependency ledger and is never
// authored by the executor.
type ReadSet map[Ref]string

func (r ReadSet) clone() ReadSet {
	out := make(ReadSet, len(r))
	for ref, fingerprint := range r {
		out[ref] = fingerprint
	}
	return out
}

// Caps is the protocol-aware capability surface for one invocation. Reads extend
// the read set automatically; there is no raw SQL access.
type Caps struct {
	db      string
	doc     *cdl.EIRDoc
	reads   ReadSet
	fetches int
}

// NewCaps opens an invocation-scoped capability surface.
func NewCaps(db string, doc *cdl.EIRDoc) *Caps {
	return &Caps{db: db, doc: doc, reads: ReadSet{}}
}

// ReadSet returns a copy of the observed-record ledger.
func (c *Caps) ReadSet() ReadSet { return c.reads.clone() }

// FetchCount reports how many payload reads this invocation performed.
func (c *Caps) FetchCount() int { return c.fetches }

// Get reads one record payload and binds its fingerprint into the read set.
func (c *Caps) Get(ref Ref) ([]byte, error) {
	payload, err := store.Get(c.db, c.doc, ref.FieldID, ref.RecordID)
	if err != nil {
		return nil, err
	}
	fingerprint, err := store.FingerprintOf(c.db, c.doc, ref.FieldID, ref.RecordID)
	if err != nil {
		return nil, err
	}
	c.reads[ref] = fingerprint
	c.fetches++
	return payload, nil
}

// List returns sorted record IDs without extending the read set; IDs alone are
// not record content.
func (c *Caps) List(fieldID string) ([]string, error) {
	return store.List(c.db, c.doc, fieldID)
}

// Search returns sorted record IDs whose canonical payload contains the exact
// needle, without extending the read set.
func (c *Caps) Search(fieldID, needle string) ([]string, error) {
	return store.Search(c.db, c.doc, fieldID, needle)
}
