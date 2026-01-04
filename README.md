## YAP (YAML as Programming)

Have you ever read YAML and thought to yourself - "I wish I could write my full code in that sick markdown language."?

**Well, now you can!**

### Status

🚧 **Early Development** - YAP is in its infancy. 

### Currently Supported

- **Print statement**: Output text to the console
(multiple lines NOT supported)

```yaml
- print: "hello world"
```

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
- print: "hello world"
```

Then run it:

```bash
yap run hello.yap
```

Output:
```
hello world
```

### License

See [LICENSE](LICENSE) for details.
