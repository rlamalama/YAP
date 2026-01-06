# YAP (YAML as Programming)

Have you ever read YAML and thought to yourself - "I wish I could write my full code in that sick markdown language."?

**Well, now you can!**

---

## Status

🚧 **Early Development** - YAP is in its infancy.

---

## Building

```bash
# Build for your local machine
make build-local

# Build for Linux (amd64)
make build

# Run tests
make test
```

---

## Running

```bash
# Run a .yap file
./bin/yap run yourfile.yap
```

---

## Roadmap

**Core Language:**
- [x] Print statements
- [x] Variables (`set`)
- [x] String literals
- [x] Numeric literals
- [x] Boolean literals (`True`/`False`)
- [x] Arithmetic operators (`+`, `-`, `*`, `/`)
- [x] Comparison operators (`>`, `<`, `>=`, `<=`, `==`, `!=`)
- [x] Comments (`//`)
- [x] Conditional statements (`if`/`then`/`else`)
- [ ] Loops (`while`)
- [ ] Functions (`function`/`call`)

**Future:**
- [ ] Lists/Arrays
- [ ] Logical operators (`and`, `or`, `not`)
- [ ] User input
- [ ] Floating-point numbers
- [ ] String operations
- [ ] File I/O

---

## Example

Create a file called `demo.yap`:

```yaml
// YAP Demo - showcasing all language features

// Variables: strings, numbers, booleans, expressions
- set:
  - name: "YAP"
  - version: 1
  - x: 10 + 5 - 3
  - y: x * 2
  - greeting: "Hello, " + name + "!"
  - isReady: True
  - isBig: x > 10

// Print values and expressions
- print: greeting
- print: "Version:"
- print: version
- print: "x ="
- print: x
- print: "y ="
- print: y

// Conditionals with nesting
- if: isBig
  then:
    - print: "x is big!"
    - if: x < 20
      then:
        - print: "but not too big"
      else:
        - print: "wow, it's huge!"
  else:
    - print: "x is small"

// Comparisons
- print: "Is x equal to 12?"
- print: x == 12
- print: "Is y greater than or equal to 30?"
- print: y >= 30

// Final message
- print: "Done!"  // inline comment
```

Run it:

```bash
./bin/yap run demo.yap
```

Output:

```
Hello, YAP!
Version:
1
x =
12
y =
24
x is big!
but not too big
Is x equal to 12?
true
Is y greater than or equal to 30?
false
Done!
```

---

## Documentation

| Document | Description |
|----------|-------------|
| [LANGUAGE.md](LANGUAGE.md) | Language guide and syntax reference |
| [GRAMMAR.md](GRAMMAR.md) | Complete formal language specification |

---

## License

See [LICENSE](LICENSE) for details.
