# Curly

Curly is a versatile Go package designed for formatting and parsing text based on custom identifiers and expressions. It provides a flexible framework for working with text modifications, number calculations, and more, making it an essential tool for various text processing tasks.

## Features

- **Format Text**: Apply a series of formatters to the text and modify it according to custom rules.
- **Parse Text**: Extract data from text based on predefined expressions and parsers.
- **Number Calculations**: Perform mathematical operations on formatted strings.
- **String Modifications**: Modify strings based on specified expressions.
- **More...**: Explore more features and functionalities in the test package [`curly_test.go`](https://github.com/ceebydith/curly/blob/main/curly_test.go).

## Installation

To install the package, use the following command:

```sh
go get github.com/ceebydith/curly
```

## Usage

### Formatting Text

Use the `Format` function to apply a series of formatters to a text string:

```go
package main

import (
    "fmt"
    "github.com/ceebydith/curly"
)

func main() {
    text := "Today's date is {yyyy}-{mm}-{dd}"
    formatted, err := curly.Format(text)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Formatted Text:", formatted)
}
```

### Parsing Text

Use the `Parse` function to extract data from text based on provided expressions and parsers:

```go
package main

import (
    "fmt"
    "github.com/ceebydith/curly"
)

func main() {
    text := "The price is 123.45 dollars"
    expression := "{price:num}"
    result, err := curly.Parse(text, expression)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Parsed Result:", result)
}
```

### Number Calculations

Use the `NumberCalculate` function to evaluate mathematical expressions after formatting:

```go
package main

import (
    "fmt"
    "github.com/ceebydith/curly"
)

func main() {
    expression := "2 * (3 + 4)"
    result, err := curly.NumberCalculate(expression)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Calculation Result:", result)
}
```

### String Modifications

Use the `StringModify` function to apply modifications to text based on specified expressions:

```go
package main

import (
    "fmt"
    "github.com/ceebydith/curly"
)

func main() {
    text := " hello world "
    expressions := "pre(Greetings, )|post(!)"
    result, err := curly.StringModify(text, expressions)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Modified String:", result)
}
```

## Documentation

For more detailed documentation, visit the [pkg.go.dev](https://pkg.go.dev/github.com/ceebydith/curly) page.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request. For major changes, please open an issue first to discuss what you would like to change.
Fork the repository
Create your feature branch (`git checkout -b feature/your-feature`)
Commit your changes (`git commit -m 'Add some feature'`)
Push to the branch (`git push origin feature/your-feature`)
Open a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/ceebydith/curly/blob/main/LICENSE) file for details.
