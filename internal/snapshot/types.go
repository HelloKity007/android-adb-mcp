package snapshot

import "fmt"

// CenterCord represents the center point of a UI element
type CenterCord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// BoundingBox represents the rectangular bounds of a UI element
type BoundingBox struct {
	X1 int `json:"x1"`
	Y1 int `json:"y1"`
	X2 int `json:"x2"`
	Y2 int `json:"y2"`
}

// ElementNode represents a single interactive UI element
type ElementNode struct {
	Name        string      `json:"name"`
	ClassName   string      `json:"class_name"`
	Coordinates CenterCord  `json:"coordinates"`
	BoundingBox BoundingBox `json:"bounding_box"`
	ResourceID  string      `json:"resource_id"`
}

// TreeState represents the full parsed UI state
type TreeState struct {
	InteractiveElements []ElementNode `json:"interactive_elements"`
}

// ToString returns a formatted table representation of the tree state
func (ts *TreeState) ToString() string {
	if len(ts.InteractiveElements) == 0 {
		return "No interactive elements found"
	}

	result := "Label | Name | ResourceId | Class | Coordinates\n"
	result += "------|------|------------|-------|------------\n"

	for i, elem := range ts.InteractiveElements {
		name := elem.Name
		if len(name) > 30 {
			name = name[:27] + "..."
		}
		result += fmt.Sprintf("%d | %s | %s | %s | (%d,%d)\n",
			i, name, elem.ResourceID, elem.ClassName,
			elem.Coordinates.X, elem.Coordinates.Y)
	}

	return result
}
