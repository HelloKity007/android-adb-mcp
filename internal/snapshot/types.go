package snapshot

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

// ElementNode represents a single interactive UI element for annotation
type ElementNode struct {
	Name        string      `json:"name"`
	ClassName   string      `json:"class_name"`
	Coordinates CenterCord  `json:"coordinates"`
	BoundingBox BoundingBox `json:"bounding_box"`
	ResourceID  string      `json:"resource_id"`
}
