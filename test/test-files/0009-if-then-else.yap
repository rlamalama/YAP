- set:
  - x: 10

- if: x > 5
  then:
    - print: "x is big"
    - if: x < 20
      then: 
        - print: "well not that big"
      else: 
        - print: "it must be huge!"
  else:
    - print: "x is small"
