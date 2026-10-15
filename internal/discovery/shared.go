package discovery

import (
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// SharedEntry is one committed SHARED-ENTRY record in bounded form: the shared
// concrete entry, its participants, and its coordination owner.
type SharedEntry struct {
	ConcreteEntrySrc  string
	Participants      []string
	CoordinationOwner string
}

// LoadSharedEntries reads every committed shared-entry record, sorted by record
// ID, so discovery units can resolve the participating personas of a shared
// cluster instead of leaving the binding unbound.
func LoadSharedEntries(db string, doc *cdl.EIRDoc) ([]SharedEntry, error) {
	names, _, _, err := work.FieldContract(doc, "SHARED-ENTRY")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "SHARED-ENTRY")
	if err != nil {
		return nil, err
	}
	sort.Strings(ids)
	var entries []SharedEntry
	for _, id := range ids {
		payload, err := store.Get(db, doc, "SHARED-ENTRY", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		entry := SharedEntry{
			ConcreteEntrySrc:  strings.Join(strings.Fields(values["concrete-entry-src"]), " "),
			CoordinationOwner: strings.Join(strings.Fields(values["coordination-owner"]), " "),
		}
		entry.Participants = append(entry.Participants, personaTokens(values["participating-personas-and-invocation-conditions"])...)
		entry.Participants = append(entry.Participants, personaTokens(values["persona-invocation-pointers"])...)
		entry.Participants = dedupeSorted(entry.Participants)
		entries = append(entries, entry)
	}
	return entries, nil
}

// SharedParticipants returns the distinct persona prefixes a shared-entry set
// declares for one cluster: entries whose concrete entry source matches one of
// the cluster's concrete entry coordinates. The coordination owner comes first
// when declared; the rest are lexicographically sorted for determinism.
func SharedParticipants(entries []SharedEntry, cluster Cluster) []string {
	var coordination string
	rest := map[string]bool{}
	for _, entry := range entries {
		if !entryMatchesCluster(entry, cluster) {
			continue
		}
		if coordination == "" && entry.CoordinationOwner != "" {
			coordination = entry.CoordinationOwner
		}
		for _, participant := range entry.Participants {
			rest[participant] = true
		}
	}
	var participants []string
	if coordination != "" {
		participants = append(participants, coordination)
		delete(rest, coordination)
	}
	var sorted []string
	for participant := range rest {
		sorted = append(sorted, participant)
	}
	sort.Strings(sorted)
	return append(participants, sorted...)
}

func entryMatchesCluster(entry SharedEntry, cluster Cluster) bool {
	if entry.ConcreteEntrySrc == "" {
		return false
	}
	for _, coordinate := range cluster.Entries {
		if entry.ConcreteEntrySrc == coordinate {
			return true
		}
	}
	return false
}

// personaTokens extracts persona prefixes from a declaration value: the leading
// token of every semicolon/comma-separated segment (for example
// "AD: local condition; WEB: ...") and every `<PREFIX>-<slug>` pointer token.
func personaTokens(raw string) []string {
	var tokens []string
	for _, piece := range strings.FieldsFunc(raw, func(r rune) bool { return r == ';' || r == ',' || r == '\n' }) {
		piece = strings.TrimSpace(piece)
		if piece == "" || strings.EqualFold(piece, "None") {
			continue
		}
		head := piece
		if idx := strings.IndexAny(piece, ": "); idx > 0 {
			head = piece[:idx]
		}
		if prefix := validPersonaPrefix(head); prefix != "" {
			tokens = append(tokens, prefix)
			continue
		}
		for _, segment := range strings.Split(piece, "/") {
			segment = strings.Trim(segment, "`[]() ")
			idx := strings.Index(segment, "-")
			if idx <= 0 || idx == len(segment)-1 {
				continue
			}
			prefix := validPersonaPrefix(segment[:idx])
			if prefix == "" {
				continue
			}
			rest := segment[idx+1]
			if rest < 'a' || rest > 'z' {
				continue
			}
			tokens = append(tokens, prefix)
		}
	}
	return tokens
}

// validPersonaPrefix reports the upper-cased token when it matches the §2.3
// persona prefix grammar and is not the reserved shared marker.
func validPersonaPrefix(token string) string {
	token = strings.ToUpper(strings.TrimSpace(token))
	if token == "" || token == "SHARED" || token == "NONE" || len(token) > 8 {
		return ""
	}
	if token[0] < 'A' || token[0] > 'Z' {
		return ""
	}
	for i := 1; i < len(token); i++ {
		c := token[i]
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return ""
		}
	}
	return token
}

func dedupeSorted(values []string) []string {
	set := map[string]bool{}
	var out []string
	for _, value := range values {
		if value == "" || set[value] {
			continue
		}
		set[value] = true
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
