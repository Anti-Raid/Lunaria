package lunaria

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Node represents a parsed XML node
type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",chardata"`
	Nodes   []Node     `xml:",any"`
}

// Handler is a function that processes a specific XML tag
type Handler func(node Node, compiler *Compiler) (string, error)

// Compiler manages the compilation process
type Compiler struct {
	handlers map[string]Handler
	indent   int
}

// NewCompiler creates a new compiler instance
func NewCompiler() *Compiler {
	c := &Compiler{
		handlers: make(map[string]Handler),
		indent:   0,
	}

	// Register built-in handlers
	c.registerBuiltins()
	return c
}

// Register adds a custom handler for a specific XML tag
func (c *Compiler) Register(tag string, handler Handler) {
	c.handlers[tag] = handler
}

// getIndent returns the current indentation string
func (c *Compiler) getIndent() string {
	return strings.Repeat("    ", c.indent)
}

// compileNode processes a single XML node
func (c *Compiler) compileNode(node Node) (string, error) {
	// Skip text nodes that are just whitespace
	if node.XMLName.Local == "" {
		content := strings.TrimSpace(node.Content)
		if content == "" {
			return "", nil
		}
		// This is unexpected content - might be an error
		return "", fmt.Errorf("unexpected text content: %s", content)
	}

	// Look up handler for this tag
	handler, exists := c.handlers[node.XMLName.Local]
	if !exists {
		return "", fmt.Errorf("unknown tag: %s", node.XMLName.Local)
	}

	return handler(node, c)
}

// CompileFromString compiles an XML string using this compiler instance
func (c *Compiler) CompileFromString(s string) (string, error) {
	var root Node
	if err := xml.Unmarshal([]byte(s), &root); err != nil {
		return "", fmt.Errorf("XML parse error: %w", err)
	}

	// Handle root script tag
	if root.XMLName.Local == "script" {
		var results []string
		for _, child := range root.Nodes {
			code, err := c.compileNode(child)
			if err != nil {
				return "", err
			}
			if code != "" {
				results = append(results, code)
			}
		}
		return strings.Join(results, "\n"), nil
	}

	// Single command
	return c.compileNode(root)
}

// CompileFromReader compiles XML from an io.Reader using this compiler instance
func (c *Compiler) CompileFromReader(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return c.CompileFromString(string(data))
}

// Package-level convenience functions using default compiler
var defaultCompiler = NewCompiler()

// Compile compiles XML bytes to Luau code using the default compiler
func Compile(b []byte) (string, error) {
	return CompileString(string(b))
}

// CompileString compiles an XML string to Luau code using the default compiler
func CompileString(s string) (string, error) {
	return defaultCompiler.CompileFromString(s)
}

// CompileReader compiles XML from an io.Reader to Luau code using the default compiler
func CompileReader(r io.Reader) (string, error) {
	return defaultCompiler.CompileFromReader(r)
}

// Register adds a handler to the default compiler
func Register(tag string, handler Handler) {
	defaultCompiler.Register(tag, handler)
}
