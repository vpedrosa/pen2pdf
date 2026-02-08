package domain

import "sort"

// FontRef represents a unique font family + weight + style combination.
type FontRef struct {
	Family string
	Weight string
	Style  string
}

// CollectFontRefs walks the document tree and returns all unique font references.
func CollectFontRefs(doc *Document) []FontRef {
	seen := make(map[FontRef]bool)
	var refs []FontRef

	for _, child := range doc.Children {
		collectFromNode(child, seen, &refs)
	}
	return refs
}

// CollectFontFamilies returns a sorted, deduplicated list of font family
// names used in the document.
func CollectFontFamilies(doc *Document) []string {
	refs := CollectFontRefs(doc)
	familySet := make(map[string]bool, len(refs))
	for _, r := range refs {
		familySet[r.Family] = true
	}
	families := make([]string, 0, len(familySet))
	for f := range familySet {
		families = append(families, f)
	}
	sort.Strings(families)
	return families
}

func collectFromNode(node Node, seen map[FontRef]bool, refs *[]FontRef) {
	switch n := node.(type) {
	case *Text:
		if n.FontFamily != "" {
			ref := FontRef{
				Family: n.FontFamily,
				Weight: n.FontWeight,
				Style:  n.FontStyle,
			}
			if !seen[ref] {
				seen[ref] = true
				*refs = append(*refs, ref)
			}
		}
	case *Frame:
		for _, child := range n.Children {
			collectFromNode(child, seen, refs)
		}
	}
}
