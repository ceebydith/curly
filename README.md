# curly

curly is a simple Go package that provides a flexible way to format strings using dynamic values. It allows you to define placeholders within a string and replace them with values from predefined maps or custom maps you provide.

## Features

- Dynamic date and time formatting
- Automatic detection of application and current working directories
- Customizable placeholders with static or dynamic values
- Easy integration into your Go projects
- Error handling for unresolved placeholders

## Installation

To install Curly, use `go get`:

```bash
go get github.com/ceebydith/curly
```

## Usage
Here's a basic example of how to use Curly:
```go
package main

import (
    "fmt"
    "github.com/ceebydith/curly"
)

func main() {
    str := "Year: {yyyy}, Month: {mm}, Day: {dd}, App Directory: {appdir}"
    result, err := curly.Format(str)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println(result)
    }
}
```
You can also provide custom maps to override the default values:
```go
customMaps := map[string]any{"yyyy": "2025"}
str := "Year: {yyyy}"
result, err := curly.Format(str, customMaps)
if err != nil {
        fmt.Println("Error:", err)
} else {
        fmt.Println(result)
}
```

## Functions
`Format`
```go
func Format(str string, maps ...map[string]any) (string, error)
```
Replaces placeholders in the input string with their corresponding values from the provided maps. Returns the formatted string and an error if unresolved placeholders remain.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/ceebydith/curly/blob/master/LICENSE) file for details.

## Acknowledgements
Special thanks to the Go community and contributors who made this project possible.


Feel free to customize the content to better fit your project's specific details, such as replacing placeholder URLs and user information. Let me know if there's anything else you'd like to add or adjust!