package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Node represents a file or directory in the tree structure
type Node struct {
	Name     string
	IsDir    bool
	Children []*Node
}

func ParseTree(input string) (*Node, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var root *Node
	var stack []*Node
	lastDepth := -1

	// Handle root node
	if scanner.Scan() {
		line := scanner.Text()
		rootName := strings.TrimSuffix(line, "/")
		root = &Node{Name: rootName, IsDir: true}
		stack = append(stack, root)
		lastDepth = 0
	}

	// Use precise regex to match tree structure format
	dirRegex := regexp.MustCompile(`^([│\s]*)([├└]── )?(.+?)/?$`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := dirRegex.FindStringSubmatch(line)
		if len(matches) != 4 {
			continue
		}

		indent := matches[1]
		name := matches[3]
		depth := strings.Count(indent, "│") + strings.Count(indent, "    ")
		isDir := strings.HasSuffix(line, "/")

		// Adjust stack based on current depth
		if depth <= lastDepth {
			stack = stack[:depth+1]
		}
		lastDepth = depth

		// Create new node
		node := &Node{
			Name:  name,
			IsDir: isDir,
		}

		// Add to parent's children list
		parent := stack[depth]
		parent.Children = append(parent.Children, node)

		// If directory, add to stack for tracking
		if isDir {
			stack = append(stack, node)
		}
	}

	return root, nil
}

func BuildFileStructure(node *Node, basePath string) error {
	path := filepath.Join(basePath, node.Name)

	// Create directory or file based on node type
	if node.IsDir {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}

		// Recursively process child nodes
		for _, child := range node.Children {
			if err := BuildFileStructure(child, path); err != nil {
				return err
			}
		}
	} else {
		// Ensure parent directory exists before creating file
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", path, err)
		}

		// Create empty file
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
		f.Close()
	}

	return nil
}

func main() {
	treeInput := `project/
├── README.md
├── src/
│   ├── index.ts
│   └── utils/
│       └── date.ts
└── .gitignore`

	root, err := ParseTree(treeInput)
	if err != nil {
		fmt.Println("Error parsing tree:", err)
		return
	}

	basePath := "."
	if err := BuildFileStructure(root, basePath); err != nil {
		fmt.Printf("Error building file structure: %v\n", err)
		return
	}

	fmt.Printf("File structure successfully created in %s\n", filepath.Join(basePath, root.Name))
}
