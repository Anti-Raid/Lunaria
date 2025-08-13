package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"lunaria/lunaria"
)

const (
	version = "1.0.0"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	switch os.Args[1] {
	case "-h", "--help", "help":
		showHelp()
	case "-v", "--version", "version":
		fmt.Printf("Lunaria %s\n", version)
	case "examples":
		showExamples()
	case "-":
		compileFromStdin()
	default:
		compileFromFile(os.Args[1])
	}
}

func showHelp() {
	fmt.Println("Lunaria XML-to-Luau Compiler")
	fmt.Printf("Version: %s\n\n", version)
	fmt.Println("USAGE:")
	fmt.Println("    lunaria [OPTIONS] [FILE]")
	fmt.Println()
	fmt.Println("ARGS:")
	fmt.Println("    <FILE>    XML file to compile (use '-' for stdin)")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("    -h, --help       Show this help message")
	fmt.Println("    -v, --version    Show version information")
	fmt.Println("    examples         Show usage examples")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("    lunaria script.xml    # Compile script.xml to Luau")
	fmt.Println("    lunaria -             # Read from stdin")
	fmt.Println("    cat script.xml | lunaria -")
}

func showExamples() {
	fmt.Println("Lunaria Examples")
	fmt.Println("================")
	fmt.Println()

	examples := []struct {
		title       string
		description string
		xml         string
	}{
		{
			"Basic Variables and Output",
			"Simple variable assignment and string interpolation",
			`<script>
  <set var="name" local="true">"World"</set>
  <set var="version" local="true">1.0</set>
  <print>Hello, {{name}}! Version: {{version}}</print>
</script>`,
		},
		{
			"Control Flow",
			"Conditionals and loops",
			`<script>
  <set var="max" local="true">5</set>
  <for var="i" from="1" to="max">
    <if test="i % 2 == 0">
      <print>{{i}} is even</print>
    </if>
  </for>
</script>`,
		},
		{
			"Functions and Tables",
			"Function definitions and table creation",
			`<script>
  <function name="calculateArea" params="width, height" local="true">
    <return>width * height</return>
  </function>
  
  <table var="rectangle" local="true">
    <entry key="width">10</entry>
    <entry key="height">20</entry>
    <entry key="area">calculateArea(10, 20)</entry>
  </table>
  
  <print>Area: {{rectangle.area}}</print>
</script>`,
		},
		{
			"Error Handling",
			"Assertions and error handling",
			`<script>
  <function name="divide" params="a, b" local="true">
    <assert test="b ~= 0">Cannot divide by zero</assert>
    <return>a / b</return>
  </function>
  
  <set var="result" local="true">divide(10, 2)</set>
  <print>Result: {{result}}</print>
</script>`,
		},
		{
			"Mixed Raw Code",
			"Combining XML commands with raw Luau",
			`<script>
  <comment>Custom math utilities</comment>
  <raw>
local function clamp(value, min, max)
    return math.max(min, math.min(max, value))
end
  </raw>
  
  <set var="value" local="true">clamp(15, 0, 10)</set>
  <print>Clamped value: {{value}}</print>
</script>`,
		},
	}

	for i, example := range examples {
		fmt.Printf("%d. %s\n", i+1, example.title)
		fmt.Printf("   %s\n\n", example.description)
		fmt.Printf("   XML:\n")
		printIndented(example.xml, "   ")
		fmt.Println()

		// Show compiled output
		if result, err := lunaria.CompileString(example.xml); err == nil {
			fmt.Printf("   Compiles to:\n")
			printIndented(result, "   ")
		} else {
			fmt.Printf("   Error: %v\n", err)
		}

		fmt.Println()
		fmt.Println(strings.Repeat("-", 60))
		fmt.Println()
	}
}

func printIndented(text, indent string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fmt.Printf("%s%s\n", indent, line)
		} else {
			fmt.Println()
		}
	}
}

func compileFromStdin() {
	result, err := lunaria.CompileReader(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}

func compileFromFile(filename string) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File '%s' does not exist\n", filename)
		os.Exit(1)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	result, err := lunaria.CompileReader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compilation error in %s: %v\n", filename, err)
		os.Exit(1)
	}

	// If output filename is not specified, print to stdout
	if len(os.Args) == 2 {
		fmt.Println(result)
		return
	}

	// Optional: Save to file if a third argument is provided
	if len(os.Args) >= 3 {
		outputFile := os.Args[2]
		if err := saveToFile(outputFile, result); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Compiled %s -> %s\n", filename, outputFile)
	}
}

func saveToFile(filename, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	return err
}

// Additional CLI utilities

func isXMLFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".xml" || ext == ".lunaria"
}

func getOutputFilename(inputFile string) string {
	ext := filepath.Ext(inputFile)
	base := strings.TrimSuffix(inputFile, ext)
	return base + ".lua"
}

// Advanced CLI features (can be extended)

func compileBatch(pattern string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return fmt.Errorf("no files match pattern: %s", pattern)
	}

	for _, filename := range matches {
		if !isXMLFile(filename) {
			continue
		}

		fmt.Printf("Compiling %s...", filename)

		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf(" ERROR: %v\n", err)
			continue
		}

		result, err := lunaria.CompileReader(file)
		file.Close()

		if err != nil {
			fmt.Printf(" ERROR: %v\n", err)
			continue
		}

		outputFile := getOutputFilename(filename)
		if err := saveToFile(outputFile, result); err != nil {
			fmt.Printf(" ERROR saving: %v\n", err)
			continue
		}

		fmt.Printf(" -> %s\n", outputFile)
	}

	return nil
}

// Watch mode (placeholder for future implementation)
func watchMode(filename string) error {
	// This would implement file watching and auto-compilation
	// For now, just return an error indicating it's not implemented
	return fmt.Errorf("watch mode not yet implemented")
}
