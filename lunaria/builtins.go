package lunaria

import (
	"fmt"
	"strings"
)

// registerBuiltins adds all built-in Lunaria commands to the compiler
func (c *Compiler) registerBuiltins() {
	c.registerVariableCommands()
	c.registerControlFlowCommands()
	c.registerFunctionCommands()
	c.registerDataCommands()
	c.registerIOCommands()
	c.registerUtilityCommands()
}

// registerVariableCommands registers variable-related commands
func (c *Compiler) registerVariableCommands() {
	// <set> command
	c.Register("set", func(node Node, compiler *Compiler) (string, error) {
		varName := GetAttr(node, "var")
		if varName == "" {
			return "", fmt.Errorf("set command requires 'var' attribute")
		}

		if !IsValidIdentifier(varName) {
			return "", fmt.Errorf("invalid variable name: %s", varName)
		}

		isLocal := GetBoolAttr(node, "local")
		value := strings.TrimSpace(node.Content)

		if value == "" {
			return "", fmt.Errorf("set command requires a value")
		}

		prefix := ""
		if isLocal {
			prefix = "local "
		}

		return fmt.Sprintf("%s%s%s = %s", compiler.getIndent(), prefix, varName, value), nil
	})
}

// registerControlFlowCommands registers control flow commands
func (c *Compiler) registerControlFlowCommands() {
	// <if> command
	c.Register("if", func(node Node, compiler *Compiler) (string, error) {
		test := GetAttr(node, "test")
		if test == "" {
			return "", fmt.Errorf("if command requires 'test' attribute")
		}

		result := fmt.Sprintf("%sif %s then\n", compiler.getIndent(), test)

		compiler.indent++
		for _, child := range node.Nodes {
			childCode, err := compiler.compileNode(child)
			if err != nil {
				return "", err
			}
			if childCode != "" {
				result += childCode + "\n"
			}
		}
		compiler.indent--

		result += compiler.getIndent() + "end"
		return result, nil
	})

	// <elseif> command (used within if blocks)
	c.Register("elseif", func(node Node, compiler *Compiler) (string, error) {
		test := GetAttr(node, "test")
		if test == "" {
			return "", fmt.Errorf("elseif command requires 'test' attribute")
		}

		result := fmt.Sprintf("%selseif %s then\n", compiler.getIndent(), test)

		compiler.indent++
		for _, child := range node.Nodes {
			childCode, err := compiler.compileNode(child)
			if err != nil {
				return "", err
			}
			if childCode != "" {
				result += childCode + "\n"
			}
		}
		compiler.indent--

		return result, nil
	})

	// <else> command (used within if blocks)
	c.Register("else", func(node Node, compiler *Compiler) (string, error) {
		result := fmt.Sprintf("%selse\n", compiler.getIndent())

		compiler.indent++
		for _, child := range node.Nodes {
			childCode, err := compiler.compileNode(child)
			if err != nil {
				return "", err
			}
			if childCode != "" {
				result += childCode + "\n"
			}
		}
		compiler.indent--

		return result, nil
	})

	// <for> command
	c.Register("for", func(node Node, compiler *Compiler) (string, error) {
		varName := GetAttr(node, "var")
		from := GetAttr(node, "from")
		to := GetAttr(node, "to")
		step := GetAttrWithDefault(node, "step", "1")

		if varName == "" {
			return "", fmt.Errorf("for command requires 'var' attribute")
		}

		if !IsValidIdentifier(varName) {
			return "", fmt.Errorf("invalid variable name: %s", varName)
		}

		var result string
		if from != "" && to != "" {
			// Numeric for loop
			if step != "1" {
				result = fmt.Sprintf("%sfor %s = %s, %s, %s do\n", compiler.getIndent(), varName, from, to, step)
			} else {
				result = fmt.Sprintf("%sfor %s = %s, %s do\n", compiler.getIndent(), varName, from, to)
			}
		} else {
			// Generic for loop (for k, v in pairs(...))
			iterator := GetAttr(node, "in")
			if iterator == "" {
				return "", fmt.Errorf("for command requires either 'from'/'to' or 'in' attributes")
			}
			result = fmt.Sprintf("%sfor %s in %s do\n", compiler.getIndent(), varName, iterator)
		}

		compiler.indent++
		for _, child := range node.Nodes {
			childCode, err := compiler.compileNode(child)
			if err != nil {
				return "", err
			}
			if childCode != "" {
				result += childCode + "\n"
			}
		}
		compiler.indent--

		result += compiler.getIndent() + "end"
		return result, nil
	})

	// <while> command
	c.Register("while", func(node Node, compiler *Compiler) (string, error) {
		test := GetAttr(node, "test")
		if test == "" {
			return "", fmt.Errorf("while command requires 'test' attribute")
		}

		result := fmt.Sprintf("%swhile %s do\n", compiler.getIndent(), test)

		compiler.indent++
		for _, child := range node.Nodes {
			childCode, err := compiler.compileNode(child)
			if err != nil {
				return "", err
			}
			if childCode != "" {
				result += childCode + "\n"
			}
		}
		compiler.indent--

		result += compiler.getIndent() + "end"
		return result, nil
	})

	// <repeat> command
	c.Register("repeat", func(node Node, compiler *Compiler) (string, error) {
		until := GetAttr(node, "until")
		if until == "" {
			return "", fmt.Errorf("repeat command requires 'until' attribute")
		}

		result := fmt.Sprintf("%srepeat\n", compiler.getIndent())

		compiler.indent++
		for _, child := range node.Nodes {
			childCode, err := compiler.compileNode(child)
			if err != nil {
				return "", err
			}
			if childCode != "" {
				result += childCode + "\n"
			}
		}
		compiler.indent--

		result += fmt.Sprintf("%suntil %s", compiler.getIndent(), until)
		return result, nil
	})

	// <break> command
	c.Register("break", func(node Node, compiler *Compiler) (string, error) {
		return compiler.getIndent() + "break", nil
	})
}

// registerFunctionCommands registers function-related commands
func (c *Compiler) registerFunctionCommands() {
	// <function> command
	c.Register("function", func(node Node, compiler *Compiler) (string, error) {
		name := GetAttr(node, "name")
		params := GetAttrWithDefault(node, "params", "")
		isLocal := GetBoolAttr(node, "local")

		if name == "" {
			return "", fmt.Errorf("function command requires 'name' attribute")
		}

		if !IsValidIdentifier(name) {
			return "", fmt.Errorf("invalid function name: %s", name)
		}

		prefix := ""
		if isLocal {
			prefix = "local "
		}

		result := fmt.Sprintf("%s%sfunction %s(%s)\n", compiler.getIndent(), prefix, name, params)

		compiler.indent++
		for _, child := range node.Nodes {
			childCode, err := compiler.compileNode(child)
			if err != nil {
				return "", err
			}
			if childCode != "" {
				result += childCode + "\n"
			}
		}
		compiler.indent--

		result += compiler.getIndent() + "end"
		return result, nil
	})

	// <call> command
	c.Register("call", func(node Node, compiler *Compiler) (string, error) {
		name := GetAttr(node, "name")
		if name == "" {
			return "", fmt.Errorf("call command requires 'name' attribute")
		}

		args := []string{}
		content := strings.TrimSpace(node.Content)
		if content != "" {
			args = append(args, content)
		}

		// Process child nodes as arguments
		for _, child := range node.Nodes {
			if child.XMLName.Local == "arg" {
				argValue := strings.TrimSpace(child.Content)
				if argValue != "" {
					args = append(args, argValue)
				}
			}
		}

		argsStr := JoinWithCommas(args)
		return fmt.Sprintf("%s%s(%s)", compiler.getIndent(), name, argsStr), nil
	})

	// <return> command
	c.Register("return", func(node Node, compiler *Compiler) (string, error) {
		content := strings.TrimSpace(node.Content)
		if content == "" {
			return compiler.getIndent() + "return", nil
		}
		return fmt.Sprintf("%sreturn %s", compiler.getIndent(), content), nil
	})

	// <arg> command (used within call blocks)
	c.Register("arg", func(node Node, compiler *Compiler) (string, error) {
		// Args are processed by the parent call command
		return "", nil
	})
}

// registerDataCommands registers data structure commands
func (c *Compiler) registerDataCommands() {
	// <table> command
	c.Register("table", func(node Node, compiler *Compiler) (string, error) {
		varName := GetAttr(node, "var")
		isLocal := GetBoolAttr(node, "local")

		prefix := ""
		if isLocal {
			prefix = "local "
		}

		if varName != "" {
			if !IsValidIdentifier(varName) {
				return "", fmt.Errorf("invalid variable name: %s", varName)
			}

			result := fmt.Sprintf("%s%s%s = {\n", compiler.getIndent(), prefix, varName)

			compiler.indent++
			for _, child := range node.Nodes {
				if child.XMLName.Local == "entry" {
					key := GetAttr(child, "key")
					value := strings.TrimSpace(child.Content)
					if key != "" && value != "" {
						if IsValidIdentifier(key) {
							result += fmt.Sprintf("%s%s = %s,\n", compiler.getIndent(), key, value)
						} else {
							result += fmt.Sprintf("%s[%s] = %s,\n", compiler.getIndent(), WrapInQuotes(key), value)
						}
					}
				}
			}
			compiler.indent--

			result += compiler.getIndent() + "}"
			return result, nil
		}

		// Inline table
		result := "{\n"
		compiler.indent++
		for _, child := range node.Nodes {
			if child.XMLName.Local == "entry" {
				key := GetAttr(child, "key")
				value := strings.TrimSpace(child.Content)
				if key != "" && value != "" {
					if IsValidIdentifier(key) {
						result += fmt.Sprintf("%s%s = %s,\n", compiler.getIndent(), key, value)
					} else {
						result += fmt.Sprintf("%s[%s] = %s,\n", compiler.getIndent(), WrapInQuotes(key), value)
					}
				}
			}
		}
		compiler.indent--
		result += compiler.getIndent() + "}"

		return result, nil
	})

	// <entry> command (used within table blocks)
	c.Register("entry", func(node Node, compiler *Compiler) (string, error) {
		// Entries are processed by the parent table command
		return "", nil
	})

	// <array> command for creating arrays
	c.Register("array", func(node Node, compiler *Compiler) (string, error) {
		varName := GetAttr(node, "var")
		isLocal := GetBoolAttr(node, "local")

		prefix := ""
		if isLocal {
			prefix = "local "
		}

		values := []string{}
		content := strings.TrimSpace(node.Content)
		if content != "" {
			values = append(values, content)
		}

		// Process child nodes as array items
		for _, child := range node.Nodes {
			if child.XMLName.Local == "item" {
				itemValue := strings.TrimSpace(child.Content)
				if itemValue != "" {
					values = append(values, itemValue)
				}
			}
		}

		arrayContent := JoinWithCommas(values)

		if varName != "" {
			if !IsValidIdentifier(varName) {
				return "", fmt.Errorf("invalid variable name: %s", varName)
			}
			return fmt.Sprintf("%s%s%s = {%s}", compiler.getIndent(), prefix, varName, arrayContent), nil
		}

		return fmt.Sprintf("{%s}", arrayContent), nil
	})

	// <item> command (used within array blocks)
	c.Register("item", func(node Node, compiler *Compiler) (string, error) {
		// Items are processed by the parent array command
		return "", nil
	})
}

// registerIOCommands registers input/output commands
func (c *Compiler) registerIOCommands() {
	// <print> command
	c.Register("print", func(node Node, compiler *Compiler) (string, error) {
		content := strings.TrimSpace(node.Content)
		if content == "" {
			return "", fmt.Errorf("print command requires content")
		}

		// Handle interpolation
		if strings.Contains(content, "{{") {
			interpolated := Interpolate(content)
			return fmt.Sprintf("%sprint(\"%s\")", compiler.getIndent(), interpolated), nil
		}

		return fmt.Sprintf("%sprint(%s)", compiler.getIndent(), content), nil
	})

	// <warn> command
	c.Register("warn", func(node Node, compiler *Compiler) (string, error) {
		content := strings.TrimSpace(node.Content)
		if content == "" {
			return "", fmt.Errorf("warn command requires content")
		}

		// Handle interpolation
		if strings.Contains(content, "{{") {
			interpolated := Interpolate(content)
			return fmt.Sprintf("%swarn(\"%s\")", compiler.getIndent(), interpolated), nil
		}

		return fmt.Sprintf("%swarn(%s)", compiler.getIndent(), content), nil
	})

	// <error> command
	c.Register("error", func(node Node, compiler *Compiler) (string, error) {
		content := strings.TrimSpace(node.Content)
		if content == "" {
			return "", fmt.Errorf("error command requires content")
		}

		level := GetAttrWithDefault(node, "level", "1")

		// Handle interpolation
		if strings.Contains(content, "{{") {
			interpolated := Interpolate(content)
			return fmt.Sprintf("%serror(\"%s\", %s)", compiler.getIndent(), interpolated, level), nil
		}

		return fmt.Sprintf("%serror(%s, %s)", compiler.getIndent(), content, level), nil
	})
}

// registerUtilityCommands registers utility commands
func (c *Compiler) registerUtilityCommands() {
	// <raw> command - pass-through Luau
	c.Register("raw", func(node Node, compiler *Compiler) (string, error) {
		content := strings.TrimSpace(node.Content)
		if content == "" {
			return "", nil
		}

		// Apply current indentation to each line
		return IndentLines(content, compiler.getIndent()), nil
	})

	// <comment> command
	c.Register("comment", func(node Node, compiler *Compiler) (string, error) {
		content := strings.TrimSpace(node.Content)
		if content == "" {
			return "", nil
		}

		comment := FormatComment(content)
		return IndentLines(comment, compiler.getIndent()), nil
	})

	// <assert> command
	c.Register("assert", func(node Node, compiler *Compiler) (string, error) {
		condition := GetAttr(node, "test")
		if condition == "" {
			return "", fmt.Errorf("assert command requires 'test' attribute")
		}

		message := strings.TrimSpace(node.Content)
		if message != "" {
			return fmt.Sprintf("%sassert(%s, %s)", compiler.getIndent(), condition, WrapInQuotes(message)), nil
		}

		return fmt.Sprintf("%sassert(%s)", compiler.getIndent(), condition), nil
	})

	// <typeof> command
	c.Register("typeof", func(node Node, compiler *Compiler) (string, error) {
		varName := GetAttr(node, "var")
		value := strings.TrimSpace(node.Content)

		if varName == "" && value == "" {
			return "", fmt.Errorf("typeof command requires either 'var' attribute or content")
		}

		if varName != "" {
			if !IsValidIdentifier(varName) {
				return "", fmt.Errorf("invalid variable name: %s", varName)
			}

			isLocal := GetBoolAttr(node, "local")
			expr := value
			if expr == "" {
				return "", fmt.Errorf("typeof command with 'var' requires content")
			}

			prefix := ""
			if isLocal {
				prefix = "local "
			}

			return fmt.Sprintf("%s%s%s = typeof(%s)", compiler.getIndent(), prefix, varName, expr), nil
		}

		// Return typeof expression directly
		return fmt.Sprintf("typeof(%s)", value), nil
	})
}
