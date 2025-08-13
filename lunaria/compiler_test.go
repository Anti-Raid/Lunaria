package lunaria

import (
	"fmt"
	"strings"
	"testing"
)

func TestBasicSet(t *testing.T) {
	xml := `<set var="x" local="true">42</set>`
	expected := `local x = 42`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestPrintWithInterpolation(t *testing.T) {
	xml := `<script>
  <set var="name" local="true">"World"</set>
  <print>Hello, {{name}}!</print>
</script>`

	expected := `local name = "World"
print("Hello, " .. tostring(name) .. "!")`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestIfStatement(t *testing.T) {
	xml := `<if test="x > 0">
  <print>"Positive"</print>
</if>`

	expected := `if x > 0 then
    print("Positive")
end`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestForLoop(t *testing.T) {
	xml := `<for var="i" from="1" to="10">
  <print>{{i}}</print>
</for>`

	expected := `for i = 1, 10 do
    print("" .. tostring(i) .. "")
end`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestForLoopWithStep(t *testing.T) {
	xml := `<for var="i" from="0" to="10" step="2">
  <print>{{i}}</print>
</for>`

	expected := `for i = 0, 10, 2 do
    print("" .. tostring(i) .. "")
end`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestGenericForLoop(t *testing.T) {
	xml := `<for var="k, v" in="pairs(table)">
  <print>{{k}}: {{v}}</print>
</for>`

	expected := `for k, v in pairs(table) do
    print("" .. tostring(k) .. ": " .. tostring(v) .. "")
end`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestWhileLoop(t *testing.T) {
	xml := `<while test="x < 10">
  <set var="x">x + 1</set>
</while>`

	expected := `while x < 10 do
    x = x + 1
end`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestFunction(t *testing.T) {
	xml := `<function name="greet" params="name" local="true">
  <print>Hello, {{name}}!</print>
  <return>"greeting sent"</return>
</function>`

	expected := `local function greet(name)
    print("Hello, " .. tostring(name) .. "!")
    return "greeting sent"
end`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestFunctionCall(t *testing.T) {
	xml := `<call name="greet">
  <arg>"Alice"</arg>
  <arg>"Bob"</arg>
</call>`

	expected := `greet("Alice", "Bob")`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestTable(t *testing.T) {
	xml := `<table var="config" local="true">
  <entry key="name">"MyApp"</entry>
  <entry key="version">1.0</entry>
  <entry key="debug">true</entry>
</table>`

	expected := `local config = {
    name = "MyApp",
    version = 1.0,
    debug = true,
}`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestArray(t *testing.T) {
	xml := `<array var="numbers" local="true">
  <item>1</item>
  <item>2</item>
  <item>3</item>
</array>`

	expected := `local numbers = {1, 2, 3}`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestRawCode(t *testing.T) {
	xml := `<raw>
local function complex()
    return math.random() * 100
end
</raw>`

	expected := `local function complex()
    return math.random() * 100
end`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestComment(t *testing.T) {
	xml := `<comment>This is a test comment</comment>`
	expected := `-- This is a test comment`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestMultiLineComment(t *testing.T) {
	xml := `<comment>This is a
multi-line
comment</comment>`

	expected := `-- This is a
-- multi-line
-- comment`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestAssert(t *testing.T) {
	xml := `<assert test="x ~= nil">Variable x must not be nil</assert>`
	expected := `assert(x ~= nil, "Variable x must not be nil")`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestComplexScript(t *testing.T) {
	xml := `<script>
  <comment>A complex example script</comment>
  <set var="numbers" local="true">{1, 2, 3, 4, 5}</set>
  
  <function name="processNumbers" params="nums" local="true">
    <set var="sum" local="true">0</set>
    <for var="i, num" in="ipairs(nums)">
      <set var="sum">sum + num</set>
      <if test="num % 2 == 0">
        <print>{{num}} is even</print>
      </if>
    </for>
    <return>sum</return>
  </function>
  
  <set var="result" local="true">processNumbers(numbers)</set>
  <print>Total sum: {{result}}</print>
</script>`

	result, err := CompileString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	// Check that it contains expected elements
	if !strings.Contains(result, "-- A complex example script") {
		t.Error("Missing comment")
	}
	if !strings.Contains(result, "local numbers = {1, 2, 3, 4, 5}") {
		t.Error("Missing numbers assignment")
	}
	if !strings.Contains(result, "local function processNumbers(nums)") {
		t.Error("Missing function declaration")
	}
	if !strings.Contains(result, "for i, num in ipairs(nums) do") {
		t.Error("Missing for loop")
	}
	if !strings.Contains(result, "print(\"Total sum: \" .. tostring(result) .. \"\")") {
		t.Error("Missing interpolated print")
	}
}

func TestCustomHandler(t *testing.T) {
	compiler := NewCompiler()

	// Register a custom log handler
	compiler.Register("log", func(node Node, c *Compiler) (string, error) {
		level := GetAttrWithDefault(node, "level", "info")
		message := strings.TrimSpace(node.Content)
		return fmt.Sprintf("%slogger.%s(%s)", c.getIndent(), level, WrapInQuotes(message)), nil
	})

	xml := `<log level="debug">Application starting</log>`
	expected := `logger.debug("Application starting")`

	result, err := compiler.CompileFromString(xml)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		name        string
		xml         string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Missing var attribute",
			xml:         `<set local="true">42</set>`,
			shouldError: true,
			errorMsg:    "set command requires 'var' attribute",
		},
		{
			name:        "Invalid variable name",
			xml:         `<set var="123invalid">42</set>`,
			shouldError: true,
			errorMsg:    "invalid variable name",
		},
		{
			name:        "Missing test attribute",
			xml:         `<if><print>"test"</print></if>`,
			shouldError: true,
			errorMsg:    "if command requires 'test' attribute",
		},
		{
			name:        "Unknown tag",
			xml:         `<unknown>content</unknown>`,
			shouldError: true,
			errorMsg:    "unknown tag: unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CompileString(tc.xml)
			if tc.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Expected error containing '%s', got: %v", tc.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkSimpleCompilation(b *testing.B) {
	xml := `<script>
  <set var="x" local="true">42</set>
  <print>Hello World</print>
</script>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompileString(xml)
		if err != nil {
			b.Fatalf("Compilation failed: %v", err)
		}
	}
}

func BenchmarkComplexCompilation(b *testing.B) {
	xml := `<script>
  <function name="fibonacci" params="n" local="true">
    <if test="n <= 1">
      <return>n</return>
    </if>
    <return>fibonacci(n-1) + fibonacci(n-2)</return>
  </function>
  
  <for var="i" from="1" to="10">
    <set var="result" local="true">fibonacci(i)</set>
    <print>Fibonacci({{i}}) = {{result}}</print>
  </for>
</script>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompileString(xml)
		if err != nil {
			b.Fatalf("Compilation failed: %v", err)
		}
	}
}
