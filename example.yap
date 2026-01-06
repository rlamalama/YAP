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