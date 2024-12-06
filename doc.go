// Package curly provides tools for flexible and powerful text formatting and parsing.
// It includes functionality for dynamically formatting text, parsing structured data,
// performing number calculations, and applying string modifications based on custom rules
// and identifiers.
package curly

/*
Curly is a robust Go package designed to provide flexible and powerful tools for text formatting and parsing.
With Curly, you can dynamically format text, parse structured data, perform calculations, and apply string modifications.

Key Features:

- Text Formatting: Apply a series of custom formatters to transform text dynamically.
- Data Parsing: Extract structured data from strings based on predefined expressions and parsers.
- Number Calculations: Evaluate mathematical expressions embedded in strings.
- String Modifications: Modify strings using a variety of transformation rules.

Installation:

To install the package, run the following command:

    go get github.com/ceebydith/curly

Usage Examples:

Formatting Text:

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

Parsing Text:

    package main

    import (
        "fmt"
        "github.com/ceebydith/curly"
    )

    func main() {
        text := "The price is 123.45 dollars"
        expression := "price is {price} dollars"
        result, err := curly.Parse(text, expression, curly.NewNumberParser())
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        fmt.Println("Parsed Result:", result)
    }

Number Calculations:

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

String Modifications:

    package main

    import (
        "fmt"
        "github.com/ceebydith/curly"
    )

    func main() {
        text := " hello world "
        expressions := "pre(Greetings, )|post(!)|trim"
        result, err := curly.StringModify(text, expressions)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        fmt.Println("Modified String:", result)
    }

Contributing:

Contributions are welcome! Please see the GitHub repository (https://github.com/ceebydith/curly) for details on how to contribute.

License:

This project is licensed under the MIT License. See the LICENSE file for details.
*/
