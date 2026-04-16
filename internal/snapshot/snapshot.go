package snapshot

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/HelloKity007/android-adb-mcp/internal/adb"
)

// SnapshotTool provides enhanced snapshot functionality
// combining UI tree parsing with annotated screenshots
type SnapshotTool struct {
	client    *adb.Client
	annotator *Annotator
}

// NewSnapshotTool creates a new snapshot tool
func NewSnapshotTool(client *adb.Client) *SnapshotTool {
	return &SnapshotTool{
		client:    client,
		annotator: NewAnnotator(0.7),
	}
}

// SnapshotResult contains the result of a snapshot operation
type SnapshotResult struct {
	TreeState        string `json:"tree_state"`
	ScreenshotPath   string `json:"screenshot_path,omitempty"`
	ScreenshotBase64 string `json:"screenshot_base64,omitempty"`
	AnnotatedPath    string `json:"annotated_path,omitempty"`
	AnnotatedBase64  string `json:"annotated_base64,omitempty"`
}

// TakeSnapshot takes a snapshot with UI tree parsing and optional annotated screenshot
func (s *SnapshotTool) TakeSnapshot(ctx context.Context, serial string, useVision bool, useAnnotation bool) (*SnapshotResult, error) {
	result := &SnapshotResult{}

	// Get UI hierarchy XML using adb-tui's existing function
	xmlData, err := s.client.DumpUIHierarchy(ctx, serial)
	if err != nil {
		return nil, fmt.Errorf("failed to dump UI hierarchy: %w", err)
	}

	// Parse using adb-tui's existing parser
	hierarchy, err := adb.ParseUIHierarchy(xmlData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UI hierarchy: %w", err)
	}

	// Flatten nodes
	elements := adb.FlattenNodes(hierarchy.Nodes)

	// Create tree state string
	treeState := "Label | Text | ResourceId | Class | Bounds\n"
	treeState += "------|------|------------|-------|--------\n"
	for i, elem := range elements {
		if elem.Text != "" || elem.ContentDescription != "" || elem.ResourceID != "" {
			name := elem.Text
			if name == "" {
				name = elem.ContentDescription
			}
			if len(name) > 30 {
				name = name[:27] + "..."
			}
			treeState += fmt.Sprintf("%d | %s | %s | %s | %v\n",
				i, name, elem.ResourceID, elem.Class, elem.Bounds)
		}
	}
	result.TreeState = treeState

	if useVision {
		// Take screenshot
		screenshotPath := filepath.Join(os.TempDir(), "android_snapshot.png")
		err = s.client.Screenshot(ctx, serial, screenshotPath)
		if err != nil {
			return nil, fmt.Errorf("failed to take screenshot: %w", err)
		}

		if useAnnotation && len(elements) > 0 {
			// Load screenshot
			img, err := LoadImage(screenshotPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load screenshot: %w", err)
			}

			// Convert adb-tui elements to snapshot ElementNode format
			var snapshotElements []ElementNode
			for _, elem := range elements {
				if elem.Text != "" || elem.ContentDescription != "" || elem.ResourceID != "" {
					snapshotElements = append(snapshotElements, ElementNode{
						Name:      elem.Text,
						ClassName: elem.Class,
						Coordinates: CenterCord{
							X: elem.Bounds.Left + (elem.Bounds.Right-elem.Bounds.Left)/2,
							Y: elem.Bounds.Top + (elem.Bounds.Bottom-elem.Bounds.Top)/2,
						},
						BoundingBox: BoundingBox{
							X1: elem.Bounds.Left,
							Y1: elem.Bounds.Top,
							X2: elem.Bounds.Right,
							Y2: elem.Bounds.Bottom,
						},
						ResourceID: elem.ResourceID,
					})
				}
			}

			// Annotate screenshot
			annotatedImg, err := s.annotator.AnnotateScreenshot(img, snapshotElements)
			if err != nil {
				return nil, fmt.Errorf("failed to annotate screenshot: %w", err)
			}

			// Save annotated screenshot
			annotatedPath := filepath.Join(os.TempDir(), "android_snapshot_annotated.png")
			err = SavePNG(annotatedImg, annotatedPath)
			if err != nil {
				return nil, fmt.Errorf("failed to save annotated screenshot: %w", err)
			}

			result.AnnotatedPath = annotatedPath

			// Convert to base64
			pngData, err := ImageToPNG(annotatedImg)
			if err == nil {
				result.AnnotatedBase64 = base64.StdEncoding.EncodeToString(pngData)
			}
		} else {
			result.ScreenshotPath = screenshotPath

			// Convert to base64
			file, err := os.Open(screenshotPath)
			if err == nil {
				defer file.Close()
				info, _ := file.Stat()
				buf := make([]byte, info.Size())
				file.Read(buf)
				result.ScreenshotBase64 = base64.StdEncoding.EncodeToString(buf)
			}
		}
	}

	return result, nil
}

// FindElement finds an element by various criteria using adb-tui's existing function
func (s *SnapshotTool) FindElement(ctx context.Context, serial string, text string, resourceID string) (*adb.UIElement, error) {
	xmlData, err := s.client.DumpUIHierarchy(ctx, serial)
	if err != nil {
		return nil, fmt.Errorf("failed to dump UI hierarchy: %w", err)
	}

	hierarchy, err := adb.ParseUIHierarchy(xmlData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UI hierarchy: %w", err)
	}

	elements := adb.FlattenNodes(hierarchy.Nodes)
	matches := adb.FindElements(elements, text, resourceID, "")
	if len(matches) == 0 {
		return nil, fmt.Errorf("element not found: text=%q resource_id=%q", text, resourceID)
	}

	return &matches[0], nil
}

// ClickElement finds and clicks an element
func (s *SnapshotTool) ClickElement(ctx context.Context, serial string, text string, resourceID string) error {
	elem, err := s.FindElement(ctx, serial, text, resourceID)
	if err != nil {
		return err
	}

	x, y := elem.Bounds.Center()
	return s.client.Tap(ctx, serial, x, y)
}

// TypeText finds an element and types text into it
func (s *SnapshotTool) TypeText(ctx context.Context, serial string, text string, resourceID string, inputText string) error {
	elem, err := s.FindElement(ctx, serial, text, resourceID)
	if err != nil {
		return err
	}

	x, y := elem.Bounds.Center()
	err = s.client.Tap(ctx, serial, x, y)
	if err != nil {
		return err
	}

	return s.client.Text(ctx, serial, inputText)
}

// GetFocusedApp returns the currently focused application
func (s *SnapshotTool) GetFocusedApp(ctx context.Context, serial string) (string, error) {
	return s.client.GetFocusedApp(ctx, serial)
}

// GetUIHierarchy returns the raw UI hierarchy XML
func (s *SnapshotTool) GetUIHierarchy(ctx context.Context, serial string) (string, error) {
	return s.client.DumpUIHierarchy(ctx, serial)
}
