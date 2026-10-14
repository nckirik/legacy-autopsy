package handoff

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/artifact"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/handbook"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/synthesis"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// CandidateManifestArtifactType is the generic artifact type of the manifest.
const CandidateManifestArtifactType = "CANDIDATE-PAYLOAD-MANIFEST"

// CandidateManifestDefaultPath is the default package-relative artifact path.
const CandidateManifestDefaultPath = "exit-e/CANDIDATE-PAYLOAD-MANIFEST.md"

// StoreMemberPath is the containing package member of every store record under
// the declared canonical storage binding.
const StoreMemberPath = "store/legacy-autopsy.sql"

// ManifestOptions are the runtime-owned manifest directives.
type ManifestOptions struct {
	ProtocolVersion    string
	ReportPath         string
	ReportArtifactPath string
	HandbookDir        string
	HandbookPrefix     string
	ArtifactPath       string
}

// CandidateManifest is one emitted candidate payload manifest artifact.
type CandidateManifest struct {
	Payload          []byte
	Fingerprint      string
	InstanceBinding  string
	SemanticBindings int
	DECBindings      int
}

type semanticBinding struct {
	recordType  string
	recordID    string
	version     string
	fingerprint string
	dependency  string
	path        string
}

// BuildCandidatePayloadManifest emits the immutable step-2 candidate payload
// manifest referencing the earlier candidate report one-way.
func BuildCandidatePayloadManifest(db string, doc *cdl.EIRDoc, opts ManifestOptions) (CandidateManifest, error) {
	var manifest CandidateManifest
	if opts.ArtifactPath == "" {
		opts.ArtifactPath = CandidateManifestDefaultPath
	}
	if opts.ReportArtifactPath == "" {
		opts.ReportArtifactPath = CandidateReportDefaultPath
	}
	if opts.HandbookPrefix == "" {
		opts.HandbookPrefix = "handbook/sections"
	}
	reportBytes, err := os.ReadFile(opts.ReportPath)
	if err != nil {
		return manifest, fmt.Errorf("handoff: candidate payload manifest requires the emitted candidate report: %w", err)
	}
	reportFingerprint, err := artifact.PayloadFingerprint(CandidateReportArtifactType, opts.ReportArtifactPath, reportBytes)
	if err != nil {
		return manifest, err
	}
	edges, err := recordDependencies(db, doc)
	if err != nil {
		return manifest, err
	}
	var bindings []semanticBinding
	var catalogs []synthesis.Catalog
	for _, catalog := range synthesis.Catalogs {
		catalogs = append(catalogs, catalog)
	}
	sort.Slice(catalogs, func(i, j int) bool { return catalogs[i].IDType < catalogs[j].IDType })
	for _, catalog := range catalogs {
		ids, err := store.List(db, doc, catalog.FieldID)
		if err != nil {
			return manifest, err
		}
		for _, id := range ids {
			version, err := store.VersionOf(db, doc, catalog.FieldID, id)
			if err != nil {
				return manifest, err
			}
			fingerprint, err := store.FingerprintOf(db, doc, catalog.FieldID, id)
			if err != nil {
				return manifest, err
			}
			bindings = append(bindings, semanticBinding{
				recordType:  catalog.IDType,
				recordID:    id,
				version:     fmt.Sprintf("%d", version),
				fingerprint: fingerprint,
				dependency:  artifact.DependencySetFingerprint(catalog.IDType, id, edges[id]),
				path:        StoreMemberPath,
			})
		}
	}
	profileNames, _, _, err := work.FieldContract(doc, "PERSONA-PROFILE")
	if err != nil {
		return manifest, err
	}
	profileIDs, err := store.List(db, doc, "PERSONA-PROFILE")
	if err != nil {
		return manifest, err
	}
	for _, id := range profileIDs {
		payload, err := store.Get(db, doc, "PERSONA-PROFILE", id)
		if err != nil {
			return manifest, err
		}
		values, err := work.PayloadText(profileNames, payload)
		if err != nil {
			return manifest, err
		}
		fingerprint, err := capabilities.SemanticContentFingerprint(payload, "semantic-content-fingerprint")
		if err != nil {
			return manifest, err
		}
		version, err := store.VersionOf(db, doc, "PERSONA-PROFILE", id)
		if err != nil {
			return manifest, err
		}
		tuples, err := profileDependencies(db, doc, values)
		if err != nil {
			return manifest, err
		}
		bindings = append(bindings, semanticBinding{
			recordType:  "PRF",
			recordID:    id,
			version:     fmt.Sprintf("%d", version),
			fingerprint: fingerprint,
			dependency:  artifact.DependencySetFingerprint("PRF", id, tuples),
			path:        StoreMemberPath,
		})
	}
	if opts.HandbookDir != "" {
		sections, err := handbook.LoadMintedSections(opts.HandbookDir)
		if err != nil {
			return manifest, err
		}
		for _, section := range sections {
			tuples := make([][4]string, 0, len(section.Upstream))
			for _, ref := range section.Upstream {
				tuples = append(tuples, [4]string{ref.ID, fmt.Sprintf("version %d %s", ref.Version, ref.Fingerprint), "SEMANTIC-CONTENT", "None"})
			}
			bindings = append(bindings, semanticBinding{
				recordType:  "HBK",
				recordID:    section.ID,
				version:     fmt.Sprintf("%d", section.Version),
				fingerprint: section.Fingerprint,
				dependency:  artifact.DependencySetFingerprint("HBK", section.ID, tuples),
				path:        opts.HandbookPrefix + "/" + section.Path,
			})
		}
	}
	sort.Slice(bindings, func(i, j int) bool {
		if bindings[i].recordType != bindings[j].recordType {
			return bindings[i].recordType < bindings[j].recordType
		}
		return bindings[i].recordID < bindings[j].recordID
	})
	decBindings, err := decContentBindings(db, doc)
	if err != nil {
		return manifest, err
	}
	version := strings.TrimSpace(opts.ProtocolVersion) + " / " + artifact.ArtifactSchemaVersion
	var b strings.Builder
	b.WriteString("<!-- ARTIFACT-PAYLOAD: " + CandidateManifestArtifactType + " -->\n\n")
	b.WriteString("# Candidate Payload Manifest\n\n")
	fmt.Fprintf(&b, "- **ARTIFACT-TYPE:** %s\n", CandidateManifestArtifactType)
	fmt.Fprintf(&b, "- **ARTIFACT-PATH:** %s\n", opts.ArtifactPath)
	fmt.Fprintf(&b, "- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** %s\n", version)
	fmt.Fprintf(&b, "- **Exit E Candidate Report Fingerprint:** %s\n", reportFingerprint)
	b.WriteString("- **Semantic Bindings:**\n")
	b.WriteString("  | Record Type | Typed Record ID | Semantic Record Version | Semantic-Content Fingerprint | Dependency-Set Fingerprint | Normalized Path |\n")
	b.WriteString("  | --- | --- | --- | --- | --- | --- |\n")
	for _, binding := range bindings {
		fmt.Fprintf(&b, "  | %s | %s | %s | %s | %s | %s |\n",
			binding.recordType, binding.recordID, binding.version, binding.fingerprint, binding.dependency, binding.path)
	}
	b.WriteString("- **DEC Content Bindings:**\n")
	b.WriteString("  | DEC ID | Decision-Content Version | Decision-Content Fingerprint |\n")
	b.WriteString("  | --- | --- | --- |\n")
	for _, dec := range decBindings {
		fmt.Fprintf(&b, "  | %s | %s | %s |\n", dec[0], dec[1], dec[2])
	}
	fmt.Fprintf(&b, "- **Semantic Binding Count:** %d\n", len(bindings))
	fmt.Fprintf(&b, "- **DEC Content Binding Count:** %d\n", len(decBindings))
	b.WriteString("- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING\n\n")
	b.WriteString("<!-- END-ARTIFACT-PAYLOAD: " + CandidateManifestArtifactType + " -->\n\n")
	b.WriteString("<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: " + CandidateManifestArtifactType + " -->\n\n")
	b.WriteString("- **Signature / Signatories / Timestamps:** None\n")
	b.WriteString("- **Envelope Audit:** not attached at the candidate stage\n")
	b.WriteString("<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: " + CandidateManifestArtifactType + " -->\n")
	payload := []byte(b.String())
	fingerprint, err := artifact.PayloadFingerprint(CandidateManifestArtifactType, opts.ArtifactPath, payload)
	if err != nil {
		return manifest, err
	}
	final := []byte(strings.Replace(string(payload),
		"- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING", "- **ARTIFACT-PAYLOAD-FINGERPRINT:** "+fingerprint, 1))
	recomputed, err := artifact.PayloadFingerprint(CandidateManifestArtifactType, opts.ArtifactPath, final)
	if err != nil {
		return manifest, err
	}
	if recomputed != fingerprint {
		return manifest, fmt.Errorf("handoff: candidate manifest fingerprint does not recompute")
	}
	if err := artifact.ValidatePayload(doc, CandidateManifestArtifactType, final); err != nil {
		return manifest, err
	}
	instance, err := artifact.InstanceBinding(CandidateManifestArtifactType, opts.ArtifactPath, fingerprint)
	if err != nil {
		return manifest, err
	}
	manifest.Payload = final
	manifest.Fingerprint = fingerprint
	manifest.InstanceBinding = instance
	manifest.SemanticBindings = len(bindings)
	manifest.DECBindings = len(decBindings)
	return manifest, nil
}

// recordDependencies groups derived-from edge tuples by derived record ID.
func recordDependencies(db string, doc *cdl.EIRDoc) (map[string][][4]string, error) {
	names, _, _, err := work.FieldContract(doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return nil, err
	}
	edges := map[string][][4]string{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "DERIVED-FROM-EDGE", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		derived := strings.Fields(values["derived-record-id"])
		if len(derived) == 0 {
			continue
		}
		tuple := [4]string{
			strings.Fields(values["source-record-id"])[0],
			values["source-version"],
			values["binding-kind"],
			strings.Fields(values["claim-id"])[0],
		}
		edges[derived[0]] = append(edges[derived[0]], tuple)
	}
	return edges, nil
}

// typedTokenPattern matches a well-formed typed protocol ID.
var typedTokenPattern = regexp.MustCompile(`^[A-Z][A-Z0-9-]*-[A-F0-9]{12,64}$`)

// profileDependencies resolves a profile's declared dependency and claim
// bindings into exact tuples. Well-formed typed IDs must resolve; free-form
// prose tokens are not authoritative references and are skipped.
func profileDependencies(db string, doc *cdl.EIRDoc, values map[string]string) ([][4]string, error) {
	var tuples [][4]string
	seen := map[string]bool{}
	addReference := func(token string) error {
		token = strings.Trim(strings.TrimSpace(token), "[]`")
		if token == "" || strings.EqualFold(token, "None") || seen[token] {
			return nil
		}
		if idx := strings.IndexAny(token, " "); idx > 0 {
			token = token[:idx]
		}
		if !typedTokenPattern.MatchString(token) {
			if field, record, ok := strings.Cut(token, "/"); ok {
				if _, _, _, err := work.FieldContract(doc, field); err == nil {
					version, err := store.VersionOf(db, doc, field, record)
					if err != nil {
						return fmt.Errorf("profile dependency %s does not resolve", token)
					}
					fingerprint, err := store.FingerprintOf(db, doc, field, record)
					if err != nil {
						return err
					}
					seen[token] = true
					tuples = append(tuples, [4]string{record, fmt.Sprintf("version %d %s", version, fingerprint), "REFERENCE", "None"})
					return nil
				}
			}
			return nil
		}
		seen[token] = true
		if strings.HasPrefix(token, "CLM-") {
			claimNames, _, _, err := work.FieldContract(doc, "CLAIM")
			if err != nil {
				return err
			}
			payload, err := store.Get(db, doc, "CLAIM", token)
			if err != nil {
				return fmt.Errorf("profile dependency %s does not resolve", token)
			}
			claim, err := work.PayloadText(claimNames, payload)
			if err != nil {
				return err
			}
			subject := strings.Fields(claim["subject-id-record-version"])
			source := "None"
			if len(subject) > 0 {
				source = subject[0]
			}
			tuples = append(tuples, [4]string{source, claim["source-fingerprint-snapshot"], "CLAIM", token})
			return nil
		}
		fieldID, ok := synthesis.FieldForID(token)
		if !ok {
			return fmt.Errorf("profile dependency %s has no declared field", token)
		}
		version, err := store.VersionOf(db, doc, fieldID, token)
		if err != nil {
			return fmt.Errorf("profile dependency %s does not resolve", token)
		}
		fingerprint, err := store.FingerprintOf(db, doc, fieldID, token)
		if err != nil {
			return err
		}
		tuples = append(tuples, [4]string{token, fmt.Sprintf("version %d %s", version, fingerprint), "REFERENCE", "None"})
		return nil
	}
	for _, field := range []string{"derived-from", "material-claims"} {
		for _, token := range strings.Split(values[field], ",") {
			for _, piece := range strings.Split(token, ";") {
				if err := addReference(piece); err != nil {
					return nil, err
				}
			}
		}
	}
	return tuples, nil
}

// decContentBindings reads committed DECISION records in bounded form.
func decContentBindings(db string, doc *cdl.EIRDoc) ([][3]string, error) {
	names, _, _, err := work.FieldContract(doc, "DECISION")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "DECISION")
	if err != nil {
		return nil, err
	}
	var rows [][3]string
	for _, id := range ids {
		payload, err := store.Get(db, doc, "DECISION", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		rows = append(rows, [3]string{id, values["decision-content-version"], values["decision-content-fingerprint"]})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i][0] < rows[j][0] })
	return rows, nil
}
