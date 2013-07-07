### Aho-Corasick string matching algorithm

This is an implementation of [Aho-Corasick string matching algorithm][ac] in Google Go.
Based on original BSD licensed [implementation][orig] by Danny Yoo from UC Berkeley.

BSD license.

Start using it by importing the dependency:
```go
import "github.com/arkadijs/ahocorasick-go"
```
then install the package with `go` command:
```shell
$ go get
```
The name of the package is `ahocorasick`.

#### API

```go
type Tree struct {
    // contains filtered or unexported fields
}
```
The `Tree` is a root objects that represents compiled state of searched terms.

```go
func New() *Tree
```
Allocates new empty `Tree` object.

```go
func (tree *Tree) Add(term string) error
```
Adds search `term` to the `Tree` object.

```go
func (tree *Tree) Search(content string) <-chan string
```
Starts search of all `Tree` terms in the `content`. Returns Go channel the found terms could be read from.

```go
	tree := ahocorasick.New()
	tree.Add("moo")
	tree.Add("one")
	for term := range tree.Search("one moon ago") {
	    fmt.Printf("found %v\n", term)
	}
```

In case you don't need the results, please close the channel -- or search goroutine may stuck:
```go
	ch := tree.Search("...")
	_, found := <-ch
	fmt.Printf("found? %v\n", found)
	close(ch)
```

[ac]: http://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_string_matching_algorithm
[orig]: https://hkn.eecs.berkeley.edu/~dyoo/java/index.html
