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
| `INDENT`     | Increase in indentation                  |
| `DEDENT`     | Decrease in indentation                  |

#### Statements

**Print** — Output a value:

```yaml
- print: "hello world"
- print: 42
- print: myVariable
```

**Set** — Assign variables:

```yaml
- set:
  - name: "Alice"
  - count: 100
```

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

### Example

Create a file called `hello.yap`:

```yaml
- set:
  - greeting: "Hello, YAP!"
  - count: 3

- print: greeting
- print: count
- print: "Done!"
```

Then run it:

```bash
yap run hello.yap
```

Output:
```
Hello, YAP!
3
Done!
```

### License

See [LICENSE](LICENSE) for details.
