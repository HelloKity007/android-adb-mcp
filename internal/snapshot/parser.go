package snapshot

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// UINode represents a node in the Android UI hierarchy XML
type UINode struct {
	XMLName        xml.Name   `xml:"node"`
	Class          string     `xml:"class,attr"`
	Text           string     `xml:"text,attr"`
	ResourceID     string     `xml:"resource-id,attr"`
	ContentDesc    string     `xml:"content-desc,attr"`
	Bounds         string     `xml:"bounds,attr"`
	Checkable      string     `xml:"checkable,attr"`
	Checked        string     `xml:"checked,attr"`
	Clickable      string     `xml:"clickable,attr"`
	Enabled        string     `xml:"enabled,attr"`
	Focusable      string     `xml:"focusable,attr"`
	Focused        string     `xml:"focused,attr"`
	Scrollable     string     `xml:"scrollable,attr"`
	LongClickable  string     `xml:"long-clickable,attr"`
	Password       string     `xml:"password,attr"`
	Selected       string     `xml:"selected,attr"`
	Hint           string     `xml:"hint,attr"`
	Nodes          []UINode   `xml:"node"`
}

// UIHierarchy represents the root of the Android UI hierarchy
type UIHierarchy struct {
	XMLName xml.Name `xml:"hierarchy"`
	Nodes   []UINode `xml:"node"`
}

// boundsRegex matches Android bounds format [x1,y1][x2,y2]
var boundsRegex = regexp.MustCompile(`\[(\d+),(\d+)\]\[(\d+),(\d+)\]`)

// ExtractBounds parses the bounds string and returns coordinates
func ExtractBounds(bounds string) (x1, y1, x2, y2 int, err error) {
	matches := boundsRegex.FindStringSubmatch(bounds)
	if matches == nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid bounds format: %s", bounds)
	}

	x1, _ = strconv.Atoi(matches[1])
	y1, _ = strconv.Atoi(matches[2])
	x2, _ = strconv.Atoi(matches[3])
	y2, _ = strconv.Atoi(matches[4])

	return x1, y1, x2, y2, nil
}

// GetCenter calculates the center point of a bounding box
func GetCenter(x1, y1, x2, y2 int) (int, int) {
	return (x1 + x2) / 2, (y1 + y2) / 2
}

// IsInteractive checks if a UI node is interactive
func IsInteractive(node UINode) bool {
	// Check state attributes
	if node.Focusable == "true" ||
		node.Clickable == "true" ||
		node.LongClickable == "true" ||
		node.Checkable == "true" ||
		node.Scrollable == "true" ||
		node.Selected == "true" ||
		node.Password == "true" {
		return true
	}

	// Check class against interactive classes whitelist
	for _, class := range INTERACTIVE_CLASSES {
		if node.Class == class {
			return true
		}
	}

	return false
}

// GetElementName resolves the name of a UI element
func GetElementName(node UINode) string {
	// Primary: content-desc or text
	name := node.ContentDesc
	if name == "" {
		name = node.Text
	}

	if name != "" {
		return name
	}

	// Fallback: collect text from children
	var texts []string
	var fallbackTexts []string

	var collectText func(n UINode, isRoot bool)
	collectText = func(n UINode, isRoot bool) {
		isActionable := !isRoot && (
			n.Clickable == "true" ||
			n.LongClickable == "true" ||
			n.Checkable == "true" ||
			n.Scrollable == "true")

		val := n.Text
		if val == "" {
			val = n.ContentDesc
		}
		if val == "" {
			val = n.Hint
		}

		if isActionable {
			if val != "" {
				fallbackTexts = append(fallbackTexts, val)
			}
			return
		}

		if val != "" {
			texts = append(texts, val)
		}

		for _, child := range n.Nodes {
			collectText(child, false)
		}
	}

	collectText(node, true)

	// Use primary texts if found, otherwise use fallback texts
	finalTexts := texts
	if len(finalTexts) == 0 {
		finalTexts = fallbackTexts
	}

	return strings.Join(finalTexts, " ")
}

// ParseUIHierarchy parses XML data into a UIHierarchy structure
func ParseUIHierarchy(xmlData string) (*UIHierarchy, error) {
	var hierarchy UIHierarchy
	err := xml.Unmarshal([]byte(xmlData), &hierarchy)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UI hierarchy XML: %w", err)
	}
	return &hierarchy, nil
}

// GetInteractiveElements extracts interactive elements from UI hierarchy
func GetInteractiveElements(hierarchy *UIHierarchy) []ElementNode {
	var elements []ElementNode

	var processNode func(node UINode)
	processNode = func(node UINode) {
		if node.Enabled != "true" {
			return
		}

		if IsInteractive(node) {
			x1, y1, x2, y2, err := ExtractBounds(node.Bounds)
			if err != nil {
				return
			}

			name := GetElementName(node)
			if name == "" {
				return
			}

			cx, cy := GetCenter(x1, y1, x2, y2)

			// Extract short resource ID
			resourceID := node.ResourceID
			if idx := strings.LastIndex(resourceID, "/"); idx >= 0 {
				resourceID = resourceID[idx+1:]
			}

			elements = append(elements, ElementNode{
				Name:      name,
				ClassName: node.Class,
				Coordinates: CenterCord{
					X: cx,
					Y: cy,
				},
				BoundingBox: BoundingBox{
					X1: x1,
					Y1: y1,
					X2: x2,
					Y2: y2,
				},
				ResourceID: resourceID,
			})
		}

		// Process children
		for _, child := range node.Nodes {
			processNode(child)
		}
	}

	// Process all root nodes
	for _, node := range hierarchy.Nodes {
		processNode(node)
	}

	return elements
}

// GetTreeState parses XML data and returns the tree state
func GetTreeState(xmlData string) (*TreeState, error) {
	hierarchy, err := ParseUIHierarchy(xmlData)
	if err != nil {
		return nil, err
	}

	elements := GetInteractiveElements(hierarchy)
	return &TreeState{
		InteractiveElements: elements,
	}, nil
}
