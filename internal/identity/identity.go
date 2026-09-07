// Package identity implements the bootstrap subset of Protocol §4.1:
// normalized relative paths and deterministic typed IDs.
package identity

import (
	"crypto/sha256"
	"fmt"
	"path"
	"strings"
)

var allowedTypes = set("CMP", "CLM", "ER", "REL", "BR", "UC", "IF", "DEP", "CFG", "SCHED", "NFR", "SEC", "SM", "DR", "CAP", "PRV", "FLT", "MOD", "INV", "AUTH", "STATE", "DEC", "CON", "FRT", "SWP", "CLU", "SHR", "DSC", "EXP", "GAP", "CNF")
var deferredSpecializedTypes = set("PRF", "HBK", "COV", "CND")
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
	sum := sha256.Sum256([]byte(canonicalKey))
	hash := strings.ToUpper(fmt.Sprintf("%x", sum))[:12]
	return Result{ID: keyType + "-" + hash, CanonicalKey: canonicalKey}, nil
}

func normalizeCoordinate(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func NormalizeRelativePath(value string) (string, error) {
	value = strings.ReplaceAll(strings.TrimSpace(value), "\\", "/")
	if value == "" {
		return "", fmt.Errorf("relative path is empty")
	}
	if strings.HasPrefix(value, "/") || (len(value) >= 2 && value[1] == ':') {
		return "", fmt.Errorf("path must be relative: %q", value)
	}
	parts := strings.Split(value, "/")
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		switch part {
		case "", ".":
			continue
		case "..":
			return "", fmt.Errorf("path traversal is forbidden: %q", value)
		default:
			cleaned = append(cleaned, part)
		}
	}
	if len(cleaned) == 0 {
		return "", fmt.Errorf("relative path resolves to root")
	}
	return path.Join(cleaned...), nil
}
