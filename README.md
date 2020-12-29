- Author: Jay Herron
- Version: 0.0.0
- License: Academic Free License version 3.0

This package is a Go implementation of the Haystack API as defined in the [Project Haystack Documentation](https://project-haystack.org/doc).
It is based on the [haystack-java](https://github.com/skyfoundry/haystack-java) package, with adjustments for compatibility with the Go API
and conventions. Currently, it implements the following:

- A Haystack Client
- The Haystack type system
- Zinc encoding and decoding
- JSON encoding

## How To Use
This package can be used by importing `gitlab.com/NeedleInAJayStack/haystack` (as is the norm in Go). Here is an example that uses the client:

    package main

	import (
		"fmt"

		"gitlab.com/NeedleInAJayStack/haystack"
	)

	func main() {
		client := haystack.NewClient(
			// INSERT YOUR URL AND CREDENTIALS
			"http://server/haystack",
			"username",
			"password",
		)
		openErr := client.Open()
		if openErr != nil {
			fmt.Println(openErr)
		}
		sites, readErr := client.Read("site")
		if readErr != nil {
			fmt.Println(readErr)
		}
		fmt.Println(sites.ToZinc())
	}

## Future Efforts
These are enhancement ideas, in no particular order.

- Add Client ops:
    - watchSub
    - watchUnsub
    - watchPoll
- Add JSON/Hayson unmarshalling support
- Add Haystack server support
- Optimize pointer usage

## Contributing
Contributions are absolutely welcome! To contribute, please create a branch, commit your changes, and make a pull request.