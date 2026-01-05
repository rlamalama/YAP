## YAP (YAML as Programming)

Have you ever read YAML and thought to yourself - "I wish I could write my full code in that sick markdown language."?

**Well, now you can!**

### Status

🚧 **Early Development** - YAP is in its infancy. 

### Language Overview

YAP uses a YAML-inspired syntax with indentation-based structure.

#### Keywords

| Keyword | Description                   |
|---------|-------------------------------|
| `print` | Output a value to the console |
| `set`   | Assign values to variables    |

#### Tokens

| Token        | Description                              |
|--------------|------------------------------------------|
| `DASH`       | `-` starts a statement                   |
| `COLON`      | `:` separates keyword/name from value    |
| `STRING`     | Text in double quotes (`"hello"`)        |
| `NUMERICAL`  | Integer literals (`42`)                  |
| `IDENTIFIER` | Variable names (`myVar`, `count`)        |
| `OPERATOR`   | Arithmetic operators (`+`, `-`, `*`, `/`)|
| `INDENT`     | Increase in indentation                  |
| `DEDENT`     | Decrease in indentation                  |

#### Statements

**Print** — Output a value or expression:

```yaml
- print: "hello world"
- print: 42
- print: myVariable
- print: 5 + 10
- print: x * 2
```

**Set** — Assign variables (supports expressions):

```yaml
- set:
  - name: "Alice"
  - count: 100
  - total: 10 + 20
  - doubled: count * 2
```

#### Operators

| Operator | Description                              |
|----------|------------------------------------------|
| `+`      | Addition (numbers) or concatenation (strings) |
| `-`      | Subtraction                              |
| `*`      | Multiplication                           |
| `/`      | Division (integer)                       |

> 📖 See [GRAMMAR.md](GRAMMAR.md) for the complete language specification.

### Building

```bash
# Build for your local machine
make build-local

# Build for Linux (amd64)
make build

# Run tests
make test
```

### Running

```bash
# Run a .yap file
./bin/yap run yourfile.yap
```

### Examples

#### Basic Example

Create a file called `hello.yap`:

```yaml
- set:
  - greeting: "Hello, YAP!"
  - count: 3

- print: greeting
- print: count
- print: "Done!"
```

Run it:

```bash
yap run hello.yap
```

Output:
```
Hello, YAP!
3
Done!
```

#### Expressions Example

Create a file called `math.yap`:

```yaml
- set:
  - x: 10 + 10 - 15
  - y: x * 4
  - z: y / 5

- print: x
- print: y
- print: x * z
```

Run it:

```bash
yap run math.yap
```

Output:
```
5
20
20
```

### License

See [LICENSE](LICENSE) for details.
