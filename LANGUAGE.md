# YAP Language Guide

This guide covers the syntax and features of the YAP language. For the complete formal specification, see [GRAMMAR.md](GRAMMAR.md).

---

## Quick Start

YAP uses a YAML-inspired syntax with indentation-based structure. Every statement starts with a dash (`-`) followed by a keyword and colon.

```yaml
- print: "Hello, World!"
```

---

## Keywords

| Keyword | Description                   |
|---------|-------------------------------|
| `print` | Output a value to the console |
| `set`   | Assign values to variables    |
| `if`    | Conditional statement         |
| `then`  | True branch of if statement   |
| `else`  | False branch of if statement  |
| `True`  | Boolean literal (true)        |
| `False` | Boolean literal (false)       |

---

## Statements

### Print

Output a value or expression to the console:

```yaml
- print: "hello world"
- print: 42
- print: myVariable
- print: 5 + 10
- print: x * 2
- print: True
- print: x > 5
```

### Set

Assign values to one or more variables. The set block uses indentation:

```yaml
- set:
  - name: "Alice"
  - count: 100
  - total: 10 + 20
  - doubled: count * 2
  - isValid: True
  - isGreater: count > 50
```

### If/Then/Else

Conditionally execute statements based on a boolean expression:

```yaml
- if: x > 5
  then:
    - print: "x is big"
  else:
    - print: "x is small"
```

The `else` block is optional:

```yaml
- if: x > 0
  then:
    - print: "x is positive"
```

Nested if statements are supported:

```yaml
- if: x > 5
  then:
    - print: "x is big"
    - if: x < 20
      then:
        - print: "but not too big"
      else:
        - print: "it's huge!"
  else:
    - print: "x is small"
```

---

## Values

### Strings

Text enclosed in double quotes:

```yaml
"hello world"
"YAP is fun!"
""
```

### Numbers

Integer literals (no decimals yet):

```yaml
42
0
12345
```

### Booleans

Truth values (case-sensitive):

```yaml
True
False
```

### Variables

Names starting with a letter or underscore:

```yaml
myVariable
count
_private
MAX_VALUE
```

---

## Operators

### Arithmetic

| Operator | Description                              |
|----------|------------------------------------------|
| `+`      | Addition (numbers) or concatenation (strings) |
| `-`      | Subtraction                              |
| `*`      | Multiplication                           |
| `/`      | Division (integer)                       |

```yaml
- set:
  - sum: 10 + 5
  - diff: 20 - 8
  - product: 4 * 3
  - quotient: 15 / 3
  - greeting: "Hello" + " " + "World"
```

### Comparison

| Operator | Description         |
|----------|---------------------|
| `>`      | Greater than        |
| `<`      | Less than           |
| `>=`     | Greater or equal    |
| `<=`     | Less or equal       |
| `==`     | Equal               |
| `!=`     | Not equal           |

```yaml
- set:
  - isGreater: 10 > 5
  - isEqual: x == y
  - notEqual: a != b
```

**Note:** All operators are evaluated left-to-right with equal precedence.

---

## Comments

Comments start with `//` and continue to the end of the line:

```yaml
// This is a full-line comment
- print: "hello"  // This is an inline comment

- set:
  - x: 10
  // - y: 5   <- This line is commented out
  - z: 20
```

---

## Indentation

YAP uses **spaces only** for indentation (tabs are not allowed). Indentation defines block structure:

```yaml
- set:
  - x: 10      // indented 2 spaces
  - y: 20      // same level as x

- if: x > 5
  then:        // indented 2 spaces
    - print: "big"   // indented 4 spaces
  else:
    - print: "small"
```

---

## Tokens Reference

| Token        | Description                              |
|--------------|------------------------------------------|
| `DASH`       | `-` starts a statement                   |
| `COLON`      | `:` separates keyword/name from value    |
| `STRING`     | Text in double quotes (`"hello"`)        |
| `NUMERICAL`  | Integer literals (`42`)                  |
| `IDENTIFIER` | Variable names (`myVar`, `count`)        |
| `OPERATOR`   | `+`, `-`, `*`, `/`, `>`, `<`, etc.       |
| `COMMENT`    | `//` starts a comment (ignored)          |
| `INDENT`     | Increase in indentation                  |
| `DEDENT`     | Decrease in indentation                  |
| `NEWLINE`    | End of a line                            |

---

## Further Reading

- [GRAMMAR.md](GRAMMAR.md) — Complete formal language specification

