package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	data := `Hello
this 
is
streaming
for
testing
and
know
about 
how
can
stream
work`

	reader := strings.NewReader(data)

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
