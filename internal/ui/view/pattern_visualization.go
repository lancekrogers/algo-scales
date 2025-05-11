package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// PatternVisualization provides ASCII/Unicode art representations of algorithm patterns
type PatternVisualization struct {
	// Stores the pattern-specific visualizations
	visualizations map[string]func(data string, width int) string
}

// NewPatternVisualization creates a new pattern visualization
func NewPatternVisualization() *PatternVisualization {
	pv := &PatternVisualization{
		visualizations: make(map[string]func(data string, width int) string),
	}

	// Register visualizations for each pattern
	pv.visualizations["sliding-window"] = pv.visualizeSlidingWindow
	pv.visualizations["two-pointers"] = pv.visualizeTwoPointers
	pv.visualizations["fast-slow-pointers"] = pv.visualizeFastSlow
	pv.visualizations["hash-map"] = pv.visualizeHashMap
	pv.visualizations["binary-search"] = pv.visualizeBinarySearch
	pv.visualizations["dfs"] = pv.visualizeDFS
	pv.visualizations["bfs"] = pv.visualizeBFS
	pv.visualizations["dynamic-programming"] = pv.visualizeDP
	pv.visualizations["greedy"] = pv.visualizeGreedy
	pv.visualizations["union-find"] = pv.visualizeUnionFind
	pv.visualizations["heap"] = pv.visualizeHeap

	return pv
}

// VisualizePattern creates a visualization for a specific algorithm pattern
func (pv *PatternVisualization) VisualizePattern(pattern, data string, width int) string {
	// Get the visualization function for this pattern
	visualizer, ok := pv.visualizations[pattern]
	if !ok {
		// Return a generic visualization if pattern not found
		return pv.visualizeGeneric(data, width)
	}

	// Apply the pattern-specific visualization
	return visualizer(data, width)
}

// visualizeSlidingWindow shows a sliding window visualization
func (pv *PatternVisualization) visualizeSlidingWindow(data string, width int) string {
	scale := MusicScales["sliding-window"]
	
	// Parse the data (expects a comma-separated list of values)
	elements := parseDataElements(data)
	if len(elements) == 0 {
		elements = []string{"1", "3", "7", "9", "10", "11"} // Default example
	}
	
	// Create the array visualization
	arrayViz := createArrayVisualization(elements, width)
	
	// Add the sliding window
	windowStart := 1
	windowEnd := 3
	if windowEnd >= len(elements) {
		windowEnd = len(elements) - 1
	}
	
	// Calculate window position and width
	windowWidth := 0
	for i := windowStart; i <= windowEnd; i++ {
		windowWidth += len(elements[i]) + 2 // +2 for the spacing
	}
	
	// Create the window indicator line
	windowStyle := lipgloss.NewStyle().Foreground(scale.PrimaryColor)
	windowLine := strings.Repeat(" ", calculatePrefixWidth(elements, windowStart))
	windowLine += "┌" + strings.Repeat("─", windowWidth-2) + "┐"
	
	// Apply style
	styledWindowLine := windowStyle.Render(windowLine)
	
	// Combine the visualization
	return styledWindowLine + "\n" + arrayViz
}

// visualizeTwoPointers shows a two pointers visualization
func (pv *PatternVisualization) visualizeTwoPointers(data string, width int) string {
	scale := MusicScales["two-pointers"]
	
	// Parse the data
	elements := parseDataElements(data)
	if len(elements) == 0 {
		elements = []string{"1", "3", "7", "9", "10", "11"} // Default example
	}
	
	// Create the array visualization
	arrayViz := createArrayVisualization(elements, width)
	
	// Add pointers at the beginning and end
	pointerStyle := lipgloss.NewStyle().Foreground(scale.PrimaryColor)
	
	// Left pointer at position 0
	leftPointerPos := 0
	leftPointerOffset := calculatePrefixWidth(elements, leftPointerPos) + 1 // +1 to center
	
	// Right pointer at the end
	rightPointerPos := len(elements) - 1
	rightPointerOffset := calculatePrefixWidth(elements, rightPointerPos) + 1
	
	// Create the pointer line
	pointerLine := strings.Repeat(" ", leftPointerOffset) + "▼"
	pointerLine += strings.Repeat(" ", rightPointerOffset-leftPointerOffset-1) + "▼"
	
	// Apply style
	styledPointerLine := pointerStyle.Render(pointerLine)
	
	// Combine the visualization
	return styledPointerLine + "\n" + arrayViz
}

// visualizeFastSlow shows a fast/slow pointer visualization
func (pv *PatternVisualization) visualizeFastSlow(data string, width int) string {
	scale := MusicScales["fast-slow-pointers"]
	
	// Create a linked list visualization
	// For simplicity, we'll use a linear representation
	nodes := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	
	// Create the linked list
	listViz := ""
	for i, node := range nodes {
		nodeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(scale.SecondaryColor).
			Padding(0, 1).
			Bold(true)
		
		styledNode := nodeStyle.Render(node)
		
		// Add arrow except for the last node
		if i < len(nodes)-1 {
			arrow := lipgloss.NewStyle().
				Foreground(scale.PrimaryColor).
				Render(" → ")
			listViz += styledNode + arrow
		} else {
			listViz += styledNode
		}
	}
	
	// Add pointers
	slowPos := 1 // B
	fastPos := 3 // D
	
	slowPointer := lipgloss.NewStyle().
		Foreground(scale.PrimaryColor).
		Render("↑ slow")
	
	fastPointer := lipgloss.NewStyle().
		Foreground(scale.AccentColor).
		Render("↑ fast")
	
	// Calculate positions
	slowOffset := slowPos * 4 + 1 // Each node is 3 chars + arrow (4 total), +1 to center
	fastOffset := fastPos * 4 + 1
	
	// Create the pointer line
	pointerLine := strings.Repeat(" ", slowOffset) + slowPointer
	pointerLine += strings.Repeat(" ", fastOffset-slowOffset-len(slowPointer)) + fastPointer
	
	// Combine the visualization
	return listViz + "\n" + pointerLine
}

// visualizeHashMap shows a hash map visualization
func (pv *PatternVisualization) visualizeHashMap(data string, width int) string {
	scale := MusicScales["hash-map"]
	
	// Create a simple hash table visualization
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Background(scale.PrimaryColor).
		Padding(0, 1).
		Bold(true)
	
	keyStyle := lipgloss.NewStyle().
		Foreground(scale.SecondaryColor).
		Bold(true)
		
	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))
		
	// Table header
	table := headerStyle.Render(" Key ") + " │ " + headerStyle.Render(" Value ") + "\n"
	table += strings.Repeat("─", 20) + "\n"
	
	// Sample data
	entries := []struct {
		key   string
		value string
	}{
		{"apple", "5"},
		{"banana", "3"},
		{"orange", "2"},
		{"grape", "8"},
	}
	
	// Build table rows
	for _, entry := range entries {
		table += keyStyle.Render(entry.key) + 
				 strings.Repeat(" ", 8-len(entry.key)) + "│ " + 
				 valueStyle.Render(entry.value) + "\n"
	}
	
	return table
}

// visualizeBinarySearch shows a binary search visualization
func (pv *PatternVisualization) visualizeBinarySearch(data string, width int) string {
	scale := MusicScales["binary-search"]
	
	// Parse the data
	elements := parseDataElements(data)
	if len(elements) == 0 {
		elements = []string{"1", "3", "7", "9", "10", "11", "15", "19", "23"} // Default sorted example
	}
	
	// Create the array visualization
	arrayViz := createArrayVisualization(elements, width)
	
	// Add pointers for lo, mid, hi
	lo := 0
	hi := len(elements) - 1
	mid := (lo + hi) / 2
	
	// Calculate positions
	loOffset := calculatePrefixWidth(elements, lo) + 1
	midOffset := calculatePrefixWidth(elements, mid) + 1
	hiOffset := calculatePrefixWidth(elements, hi) + 1
	
	// Style for pointers
	loStyle := lipgloss.NewStyle().Foreground(scale.PrimaryColor)
	midStyle := lipgloss.NewStyle().Foreground(scale.SecondaryColor).Bold(true)
	hiStyle := lipgloss.NewStyle().Foreground(scale.PrimaryColor)
	
	// Create the pointer line
	pointerLine := strings.Repeat(" ", loOffset) + loStyle.Render("▼")
	pointerLine += strings.Repeat(" ", midOffset-loOffset-1) + midStyle.Render("▼")
	pointerLine += strings.Repeat(" ", hiOffset-midOffset-1) + hiStyle.Render("▼")
	
	// Create the label line
	labelLine := strings.Repeat(" ", loOffset) + loStyle.Render("lo")
	labelLine += strings.Repeat(" ", midOffset-loOffset-2) + midStyle.Render("mid")
	labelLine += strings.Repeat(" ", hiOffset-midOffset-3) + hiStyle.Render("hi")
	
	// Combine the visualization
	return pointerLine + "\n" + arrayViz + "\n" + labelLine
}

// visualizeDFS shows a DFS visualization
func (pv *PatternVisualization) visualizeDFS(data string, width int) string {
	scale := MusicScales["dfs"]
	
	// Simple tree visualization
	tree := lipgloss.NewStyle().Foreground(scale.PrimaryColor).Render(`
    1
   / \
  2   3
 / \   \
4   5   6
    `)[1:] // Trim the first newline
	
	// Add traversal order
	traversal := lipgloss.NewStyle().
		Foreground(scale.SecondaryColor).
		Bold(true).
		Render("DFS Traversal: 1→2→4→5→3→6")
	
	return tree + "\n" + traversal
}

// visualizeBFS shows a BFS visualization
func (pv *PatternVisualization) visualizeBFS(data string, width int) string {
	scale := MusicScales["bfs"]
	
	// Simple tree visualization
	tree := lipgloss.NewStyle().Foreground(scale.PrimaryColor).Render(`
    1
   / \
  2   3
 / \   \
4   5   6
    `)[1:] // Trim the first newline
	
	// Add traversal order
	traversal := lipgloss.NewStyle().
		Foreground(scale.SecondaryColor).
		Bold(true).
		Render("BFS Traversal: 1→2→3→4→5→6")
	
	return tree + "\n" + traversal
}

// visualizeDP shows a dynamic programming visualization
func (pv *PatternVisualization) visualizeDP(data string, width int) string {
	scale := MusicScales["dynamic-programming"]
	
	// Create a simple DP table (e.g., for fibonacci)
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Background(scale.PrimaryColor).
		Padding(0, 1).
		Bold(true)
	
	cellStyle := lipgloss.NewStyle().
		Foreground(scale.SecondaryColor).
		Bold(true)
		
	// Table header
	table := headerStyle.Render(" n ") + " │ " + headerStyle.Render(" F(n) ") + "\n"
	table += strings.Repeat("─", 15) + "\n"
	
	// Sample data for fibonacci
	fibs := []int{0, 1, 1, 2, 3, 5, 8, 13, 21}
	
	// Build table rows
	for i, fib := range fibs {
		table += fmt.Sprintf(" %d │ ", i) + cellStyle.Render(fmt.Sprintf("%d", fib)) + "\n"
	}
	
	return table
}

// visualizeGreedy shows a greedy algorithm visualization
func (pv *PatternVisualization) visualizeGreedy(data string, width int) string {
	scale := MusicScales["greedy"]
	
	// Example: coin change problem with greedy approach
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Background(scale.PrimaryColor).
		Padding(0, 1).
		Bold(true)
	
	// Coins available
	coins := []int{25, 10, 5, 1}
	target := 43
	
	// Table showing greedy choice at each step
	table := headerStyle.Render(" Step ") + " │ " + 
			 headerStyle.Render(" Coin ") + " │ " + 
			 headerStyle.Render(" Remaining ") + "\n"
	table += strings.Repeat("─", 30) + "\n"
	
	// Simulate greedy algorithm
	remaining := target
	step := 1
	
	for _, coin := range coins {
		for remaining >= coin {
			table += fmt.Sprintf("  %d   │  %2d   │    %2d     \n", step, coin, remaining-coin)
			remaining -= coin
			step++
			if remaining == 0 {
				break
			}
		}
		if remaining == 0 {
			break
		}
	}
	
	// Add header
	header := lipgloss.NewStyle().
		Foreground(scale.SecondaryColor).
		Bold(true).
		Render(fmt.Sprintf("Making change for %d cents:", target))
	
	return header + "\n\n" + table
}

// visualizeUnionFind shows a union-find visualization
func (pv *PatternVisualization) visualizeUnionFind(data string, width int) string {
	scale := MusicScales["union-find"]
	
	// Create a simple visualization of connected components
	setStyle := func(id int) lipgloss.Style {
		colors := []lipgloss.Color{
			scale.PrimaryColor,
			scale.SecondaryColor,
			scale.AccentColor,
			lipgloss.Color("#2ecc71"),
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(colors[id % len(colors)]).
			Padding(0, 1).
			Bold(true)
	}
	
	// Elements with their set IDs
	elements := []struct {
		value string
		setID int
	}{
		{"A", 0}, {"B", 0}, {"C", 0},
		{"D", 1}, {"E", 1},
		{"F", 2}, {"G", 2}, {"H", 2},
		{"I", 3},
	}
	
	// Create the visualization (each row is a set)
	setMap := make(map[int][]string)
	for _, elem := range elements {
		setMap[elem.setID] = append(setMap[elem.setID], elem.value)
	}
	
	viz := lipgloss.NewStyle().
		Foreground(scale.PrimaryColor).
		Bold(true).
		Render("Union-Find Sets:") + "\n\n"
	
	for id, members := range setMap {
		setViz := "Set " + fmt.Sprint(id) + ": "
		for i, member := range members {
			setViz += setStyle(id).Render(member)
			if i < len(members)-1 {
				setViz += " "
			}
		}
		viz += setViz + "\n"
	}
	
	return viz
}

// visualizeHeap shows a heap/priority queue visualization
func (pv *PatternVisualization) visualizeHeap(data string, width int) string {
	scale := MusicScales["heap"]
	
	// Simple max-heap visualization
	heap := lipgloss.NewStyle().Foreground(scale.PrimaryColor).Render(`
      9
     / \
    7   8
   / \ / \
  5  6 3  2
    `)[1:] // Trim the first newline
	
	// Add heap property description
	description := lipgloss.NewStyle().
		Foreground(scale.SecondaryColor).
		Render("Max Heap: Each parent is greater than its children")
	
	return heap + "\n" + description
}

// visualizeGeneric provides a generic algorithm visualization
func (pv *PatternVisualization) visualizeGeneric(data string, width int) string {
	elements := parseDataElements(data)
	if len(elements) == 0 {
		elements = []string{"1", "3", "7", "9", "10", "11"} // Default example
	}
	
	return createArrayVisualization(elements, width)
}

// Helper functions

// parseDataElements parses a comma-separated list of values
func parseDataElements(data string) []string {
	if data == "" {
		return nil
	}
	
	// Split by comma and trim spaces
	elements := strings.Split(data, ",")
	for i := range elements {
		elements[i] = strings.TrimSpace(elements[i])
	}
	
	return elements
}

// createArrayVisualization creates a visualization of an array
func createArrayVisualization(elements []string, width int) string {
	// Create array brackets
	result := "["
	
	// Add elements with spacing
	for i, elem := range elements {
		result += elem
		if i < len(elements)-1 {
			result += ", "
		}
	}
	
	result += "]"
	return result
}

// calculatePrefixWidth calculates the width of the array visualization up to a specific index
func calculatePrefixWidth(elements []string, index int) int {
	if index >= len(elements) {
		index = len(elements) - 1
	}
	
	width := 1 // Starting [
	
	for i := 0; i < index; i++ {
		width += len(elements[i]) + 2 // element + ", "
	}
	
	return width
}