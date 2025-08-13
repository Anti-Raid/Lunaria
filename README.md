# Lunaria
Lunaria is an XML-based scripting language that compiles to Luau. It lets you author structured scripts in XML and transpile them to Luau for use in Roblox projects or any Luau runtime.

> ✨ Why Lunaria? You get schema-able, tool-friendly XML for authoring, with clean Luau generation for execution.

## Features

- XML-powered DSL → write logic as XML, great for pipelines, tools, or designers.
- Transpiles to Luau → output idiomatic Luau ready for Roblox Studio or any Luau VM.
- Pluggable commands → register custom tags and code generators in Go.
- Safe by design → control what functions and libraries are exposed.

## Install

```bash
go get github.com/Anti-Raid/Lunaria
```

## Quick Start
### Example XML script

```xml
<script>
  <set var="name" local="true">"AntiRaid"</set>
  <print>Hello, {{name}}!</print>
</script>
```

### Compile in Go

```go
package main

import (
    "fmt"
    "github.com/Anti-Raid/Lunaria"
)

func main() {
    luau, err := lunaria.CompileString(`<script><print>Hello World</print></script>`)
    if err != nil { panic(err) }
    fmt.Println(luau)
}
```

### Core XML Spec
```xml
<set var="x" local="true|false">EXPR</set> → local x = EXPR

<print>TEXT {{var}}</print> → print(...) with interpolation

<if test="EXPR">...</if> → conditional

<for var="i" from="A" to="B">...</for> → numeric loop

<call name="FN">...</call> → function call

<raw>...</raw> → pass-through Luau
```

### Go API
```xml
func Compile(b []byte) (string, error)
func CompileString(s string) (string, error)
func CompileReader(r io.Reader) (string, error)

type Handler func(node Node) (string, error)
func Register(tag string, h Handler)
```