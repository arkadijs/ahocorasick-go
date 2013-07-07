// Package ahocorasick implements AhoCorack string matching algorithm
// http://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_string_matching_algorithm
package ahocorasick

import (
    "container/list"
    "errors"
    // "fmt"
)

// The root objects that represents compiled state of searched terms
type Tree struct {
    prepared bool
    root     *state
}

// Allocates new empty Tree object
func New() *Tree {
    return &Tree{prepared: false, root: newState(0)}
}

// state
type state struct {
    depth int
    edges edgeList
    fail  *state
    out   *list.List
}

func newState(depth int) *state {
    return &state{depth: depth, edges: newEdgeList(depth), fail: nil, out: list.New()}
}
func (state *state) addOutput(output string) {
    state.out.PushFront(output)
}
func (state *state) extendAll(term []byte) *state {
    s := state
    for i := 0; i < len(term); i++ {
        b := term[i]
        s2 := s.edges.get(b)
        if s2 != nil {
            s = s2
        } else {
            s = s.extend(b)
        }
    }
    return s
}
func (state *state) extend(b byte) *state {
    s := state.edges.get(b)
    if s != nil {
        return s
    }
    nextState := newState(state.depth + 1)
    state.edges.put(b, nextState)
    return nextState
}
func (state *state) size() int {
    res := 1
    for _, k := range state.edges.keys() {
        //fmt.Printf("iter %v\n", k)
        res += state.edges.get(k).size()
    }
    return res
}

// edges
type edgeList interface {
    get(ch byte) *state
    put(ch byte, state *state)
    keys() []byte
}

const threshold = 3

func newEdgeList(depth int) edgeList {
    if depth < threshold {
        return &denseEdgeList{}
    }
    return &sparseEdgeList{make(map[byte]*state)}
}

// dense edge list uses array
type denseEdgeList struct {
    entries [256]*state
}

func (edges *denseEdgeList) get(ch byte) *state {
    return edges.entries[ch]
}
func (edges *denseEdgeList) put(ch byte, state *state) {
    edges.entries[ch] = state
}
func (edges *denseEdgeList) keys() []byte {
    var keys []byte
    for i := 0; i < 256; i++ {
        if edges.entries[i] != nil {
            keys = append(keys, byte(i))
        }
    }
    return keys
}

// sparse edge list uses map
type sparseEdgeList struct {
    entries map[byte]*state
}

func (edges *sparseEdgeList) get(ch byte) *state {
    return edges.entries[ch]
}
func (edges *sparseEdgeList) put(ch byte, state *state) {
    edges.entries[ch] = state
}
func (edges *sparseEdgeList) keys() []byte {
    keys := make([]byte, len(edges.entries))
    i := 0
    for k := range edges.entries {
        keys[i] = k
        i++
    }
    return keys
}

// search result 
type searchResult struct {
    lastMatchedState *state
    bytes            []byte
    lastIndex        int
}

func (res *searchResult) out() *list.List {
    return res.lastMatchedState.out
}

// Adds search term to the Tree object
func (tree *Tree) Add(term string) error {
    if tree.prepared {
        return errors.New("ahocorasick: tree already prepared - can't Add() search terms after Search() is called")
    }
    tree.root.extendAll([]byte(term)).addOutput(term)
    return nil
}

// Initializes the fail transitions of all states
func (tree *Tree) prepare() {
    if !tree.prepared {
        tree.prepared = true
        q := list.New()
        for i := 0; i < 256; i++ {
            state := tree.root.edges.get(byte(i))
            if state != nil {
                state.fail = tree.root
                q.PushBack(state)
            } else {
                tree.root.edges.put(byte(i), tree.root)
            }
        }
        for e := q.Front(); e != nil; e = e.Next() {
            state := e.Value.(*state)
            keys := state.edges.keys()
            // fmt.Printf("keys %v\n", keys)
            for _, a := range keys {
                // fmt.Printf("state %v\n", state)
                r := state
                s := r.edges.get(a)
                // fmt.Printf("pushing %v %v\n", a, s)
                q.PushBack(s)
                r = r.fail
                for r.edges.get(a) == nil {
                    r = r.fail
                }
                r = r.edges.get(a)
                s.fail = r
                s.out.PushBackList(r.out)
            }
        }
    }
}

// starts search from the root
func (tree *Tree) startSearch(content string) *searchResult {
    tree.prepare()
    res := searchResult{lastMatchedState: tree.root, bytes: []byte(content), lastIndex: 0}
    return res.continueSearch()
}

// continues walking the tree for more matches
func (lastResult *searchResult) continueSearch() *searchResult {
    bytes := lastResult.bytes
    last := lastResult.lastMatchedState
    for i := lastResult.lastIndex; i < len(bytes); i++ {
        b := bytes[i]
        var state2 *state = nil
        for state2 == nil {
            state2 = last.edges.get(b)
            last = last.fail
        }
        last = state2
        if last.out.Front() != nil {
            return &searchResult{lastMatchedState: last, bytes: bytes, lastIndex: i + 1}
        }
    }
    return nil
}

/*
  Starts search of all Tree's terms in the content.
  Returns channel the found terms could be read from.
  Usage:
    tree := ahocorasick.New()
    tree.Add("moo")
    tree.Add("one")
    for term := range tree.Search("one moon ago") {
        fmt.Printf("found %v\n", term)
    }

  In case you don't need the results, please close the channel -- or search goroutine will stuck:
    ch := tree.Search("...")
    _, found := <-ch
    fmt.Printf("found? %v\n", found)
    close (ch)
*/
func (tree *Tree) Search(content string) <-chan string {
    c := make(chan string)
    go func() {
        res := tree.startSearch(content)
        for res != nil {
            for e := res.out().Front(); e != nil; e = e.Next() {
                c <- e.Value.(string)
            }
            res = res.continueSearch()
        }
        close(c)
    }()
    return c
}
