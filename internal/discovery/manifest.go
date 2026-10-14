package discovery

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// povLabels name the five isolated POVs on the reading surface.
var povLabels = map[string]string{
	"POV-1": "Frontend Workflows",
	"POV-2": "Backend Workflows",
	"POV-3": "Background Workflows",
	"POV-4": "Data Movement Workflows",
	"POV-5": "System Wiring Workflows",
}

// PersonaManifest is one persona-local POV projection.
type PersonaManifest struct {
	Directory string
	FileName  string
	Path      string
	Markdown  []byte
}

// POVManifests projects the committed shared references into persona-local POV
// manifest files (<out>/{persona-directory}/{NN}-{NAME}-POV.md). The projection
// is deterministic and non-authoritative: the canonical store stays the source
// of truth and the files carry no semantic state. A reference whose component
// is not committed, or whose persona is not registered, is skipped instead of
// guessed; the returned manifests contain only attributions the store proves.
func POVManifests(db string, doc *cdl.EIRDoc) ([]PersonaManifest, error) {
	components, err := LoadComponents(db, doc)
	if err != nil {
		return nil, err
	}
	byID := map[string]ComponentRef{}
	for _, component := range components {
		byID[component.ID] = component
	}
	directories, err := personaDirectories(db, doc)
	if err != nil {
		return nil, err
	}
	names, _, _, err := work.FieldContract(doc, "SHARED-REFERENCE")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "SHARED-REFERENCE")
	if err != nil {
		return nil, err
	}
	sort.Strings(ids)
	type key struct {
		directory string
		pov       string
	}
	grouped := map[key][]string{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "SHARED-REFERENCE", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		componentID := componentIDFromLocation(values["canonical-location"])
		if componentID == "" {
			continue
		}
		component, ok := byID[componentID]
		if !ok || component.Persona == "" {
			continue
		}
		directory, ok := directories[component.Persona]
		if !ok || directory == "" {
			continue
		}
		if povFiles[component.POV] == "" {
			continue
		}
		entry := fmt.Sprintf("- **%s** - local trigger: `%s`; track: %s; capability state: %s; canonical staging: %s; canonical location: %s",
			component.ID, component.VirtualCoordinate, orNone(component.Track), orNone(component.CapabilityState), orNone(component.StagingStatus), values["canonical-location"])
		k := key{directory: directory, pov: component.POV}
		grouped[k] = append(grouped[k], entry)
	}
	var keys []key
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].directory != keys[j].directory {
			return keys[i].directory < keys[j].directory
		}
		return keys[i].pov < keys[j].pov
	})
	var manifests []PersonaManifest
	for _, k := range keys {
		entries := grouped[k]
		sort.Strings(entries)
		var b strings.Builder
		fmt.Fprintf(&b, "# %s - %s POV manifest\n\n", k.directory, k.pov)
		b.WriteString("> Generated projection from committed shared references. The canonical store is authoritative; this navigation view carries no independent semantic state.\n\n")
		fmt.Fprintf(&b, "- **Persona Directory:** %s\n", k.directory)
		fmt.Fprintf(&b, "- **POV:** %s (%s)\n", k.pov, povLabels[k.pov])
		fmt.Fprintf(&b, "- **Shared References:** %d\n\n", len(entries))
		b.WriteString("## Shared references\n\n")
		for _, entry := range entries {
			b.WriteString(entry)
			b.WriteString("\n")
		}
		fileName := povFiles[k.pov]
		filePath := path.Join("personas", k.directory, fileName)
		manifests = append(manifests, PersonaManifest{Directory: k.directory, FileName: fileName, Path: filePath, Markdown: []byte(b.String())})
	}
	return manifests, nil
}

// personaDirectories maps persona prefixes to canonical persona-directory
// basenames from the committed 0A registry.
func personaDirectories(db string, doc *cdl.EIRDoc) (map[string]string, error) {
	names, _, _, err := work.FieldContract(doc, "PERSONA-REGISTRY")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "PERSONA-REGISTRY")
	if err != nil {
		return nil, err
	}
	directories := map[string]string{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "PERSONA-REGISTRY", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		prefix := strings.ToUpper(strings.TrimSpace(values["prefix"]))
		directory := strings.TrimSpace(values["persona-directory"])
		if prefix != "" && directory != "" {
			directories[prefix] = directory
		}
	}
	return directories, nil
}

// componentIDFromLocation extracts the CMP ID from a canonical shared location
// of the form `personas/_shared/XX-NAME-POV.md -> CMP-...`.
func componentIDFromLocation(location string) string {
	index := strings.LastIndex(location, "->")
	if index < 0 {
		return ""
	}
	id := strings.TrimSpace(location[index+2:])
	if !strings.HasPrefix(id, "CMP-") {
		return ""
	}
	return id
}
