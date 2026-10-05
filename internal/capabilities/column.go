package capabilities

import "strings"

// ColumnName maps a protocol field label to the declared column identifier
// convention: lowercase alphanumeric runs joined with hyphens
// ("Verification Vectors" -> "verification-vectors"). It is the deterministic
// inverse used when FIELD schemas were derived from normative record labels, so
// ingestion can match payload labels to declared columns without a lookup table.
func ColumnName(label string) string {
	var b strings.Builder
	pendingHyphen := false
	for _, r := range strings.ToLower(label) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			if pendingHyphen && b.Len() > 0 {
				b.WriteByte('-')
			}
			pendingHyphen = false
			b.WriteRune(r)
		default:
			pendingHyphen = true
		}
	}
	return b.String()
}
