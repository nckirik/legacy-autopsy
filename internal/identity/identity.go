// Package identity implements the bootstrap subset of Protocol §4.1:
// normalized relative paths and deterministic typed IDs.
package identity

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
)

var allowedTypes = set("CMP", "CLM", "ER", "REL", "BR", "UC", "IF", "DEP", "CFG", "SCHED", "NFR", "SEC", "SM", "DR", "CAP", "PRV", "FLT", "MOD", "PRF", "HBK", "INV", "AUTH", "STATE", "DEC", "CON", "FRT", "SWP", "CLU", "SHR", "DSC", "EXP", "GAP", "CNF")
var deferredSpecializedTypes = set("COV", "CND")

var (
	personaPrefixPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]{1,7}$`)
	hbkAnchorPattern     = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)
var sourceKinds = set("FILE", "SYM", "ROUTE", "RPC", "JOB", "QUERY", "TABLE", "VIEW", "DR", "CONSTRAINT", "GENERATED", "CFG", "WORKFLOW", "NODE", "EDGE", "PAGE", "WIDGET", "ACTION", "BINDING", "QUEUE", "EVENT", "STORAGE", "IFACE", "OTHER")

func set(values ...string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

type Result struct {
	ID           string
	CanonicalKey string
}

func Generate(recordType, sourceKind, namespace, ownerCoordinate, discriminator string) (Result, error) {
	recordType = strings.ToUpper(strings.TrimSpace(recordType))
	keyType := recordType
	if recordType == "SRC" {
		sourceKind = strings.ToUpper(strings.TrimSpace(sourceKind))
		if _, ok := sourceKinds[sourceKind]; !ok {
			return Result{}, fmt.Errorf("unsupported SRC kind %q", sourceKind)
		}
		keyType = "SRC-" + sourceKind
	} else if _, deferred := deferredSpecializedTypes[recordType]; deferred {
		return Result{}, fmt.Errorf("protocol ID type %q requires specialized coordinates not implemented by this bootstrap", recordType)
	} else if _, ok := allowedTypes[recordType]; !ok {
		return Result{}, fmt.Errorf("unsupported protocol ID type %q", recordType)
	}
	switch recordType {
	case "PRF":
		if discriminator != "persona-profile" {
			return Result{}, fmt.Errorf("PRF semantic discriminator must be the literal persona-profile")
		}
		if !personaPrefixPattern.MatchString(strings.TrimSpace(ownerCoordinate)) {
			return Result{}, fmt.Errorf("PRF owner coordinate must be a persona prefix matching [A-Z][A-Z0-9]{1,7}")
		}
		ownerCoordinate = strings.TrimSpace(ownerCoordinate)
	case "HBK":
		if discriminator != "handbook-section" {
			return Result{}, fmt.Errorf("HBK semantic discriminator must be the literal handbook-section")
		}
		path, anchor, found := strings.Cut(ownerCoordinate, "#")
		if !found {
			return Result{}, fmt.Errorf("HBK owner coordinate must be normalized-path#stable-anchor")
		}
		normalized, err := NormalizeRelativePath(path)
		if err != nil {
			return Result{}, fmt.Errorf("HBK path: %w", err)
		}
		if !hbkAnchorPattern.MatchString(anchor) {
			return Result{}, fmt.Errorf("HBK stable section anchor %q must match [a-z0-9]+(?:-[a-z0-9]+)*", anchor)
		}
		ownerCoordinate = normalized + "#" + anchor
	}
	namespace = normalizeCoordinate(namespace)
	if recordType == "SRC" && sourceKind == "FILE" {
		var err error
		ownerCoordinate, err = NormalizeRelativePath(ownerCoordinate)
		if err != nil {
			return Result{}, fmt.Errorf("normalize SRC FILE owner: %w", err)
		}
	} else {
		ownerCoordinate = normalizeCoordinate(ownerCoordinate)
	}
	discriminator = normalizeCoordinate(discriminator)
	if namespace == "" || ownerCoordinate == "" || discriminator == "" {
		return Result{}, fmt.Errorf("namespace, owner coordinate, and discriminator are required")
	}
	if strings.Contains(namespace, "|") || strings.Contains(ownerCoordinate, "|") || strings.Contains(discriminator, "|") {
		return Result{}, fmt.Errorf("canonical-key components must not contain the delimiter '|'")
	}
	canonicalKey := strings.Join([]string{keyType, namespace, ownerCoordinate, discriminator}, "|")
	id, err := capabilities.TypedID(keyType, namespace, ownerCoordinate, discriminator)
	if err != nil {
		return Result{}, err
	}
	return Result{ID: id, CanonicalKey: canonicalKey}, nil
}

func normalizeCoordinate(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

// NormalizeRelativePath delegates to the declared path.normalize capability.
func NormalizeRelativePath(value string) (string, error) {
	return capabilities.NormalizeRelativePath(value)
}
