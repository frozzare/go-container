# container

Simple container with dependency injection.

View the [docs](http://godoc.org/github.com/frozzare/go-container).

## Installation

```
$ go get github.com/frozzare/go-container
```

## Example

```go
package main

import (
	"fmt"

	"github.com/frozzare/go-container"
)

type User struct {
	Name string
}

func main() {
	c := container.Instance()

	c.Bind("*main.User", &User{"Fredrik"})
	c.Bind("Username", func(text string, u *User) string {
		return fmt.Sprintf("%s %s", text, u.Name)
	})

	v, _ := c.Make("Username", "Hello")

	fmt.Println(v.(string))
}
```

## License

 MIT © [Fredrik Forsmo](https://github.com/frozzare)
