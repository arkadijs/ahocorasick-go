//nolint:errcheck // intentional
package ahocorasick

import (
	"fmt"
	"testing"
)

func TestStateSingleState(t *testing.T) {
	s := newState(0)
	if 1 != s.size() {
		t.Logf("size = %d\n", s.size())
		t.Fail()
	}
}

func TestStateSingleStateSparse(t *testing.T) {
	s := newState(10)
	if 1 != s.size() {
		t.Logf("size = %d\n", s.size())
		t.Fail()
	}
}

func TestStateExtend(t *testing.T) {
	s := newState(0)
	s2 := s.extend(23)

	if s == s2 || s2 == nil || 2 != s.size() {
		t.Logf("edge = %v\n", s.edges.get(23))
		t.Logf("keys = %v\n", s.edges.keys())
		t.Logf("size = %d\n", s.size())
		t.Fail()
	}
}

func TestStateExtendSparse(t *testing.T) {
	s := newState(10)
	s2 := s.extend(23)

	if s == s2 || s2 == nil || 2 != s.size() {
		t.Logf("edge = %v\n", s.edges.get(23))
		t.Logf("keys = %v\n", s.edges.keys())
		t.Logf("size = %d\n", s.size())
		t.Fail()
	}
}

func TestStateExtendAll(t *testing.T) {
	s := newState(0)
	s2 := s.extendAll([]byte("hello world!"))

	if s == s2 || s2 == nil || 13 != s.size() {
		t.Logf("size = %d\n", s.size())
		t.Fail()
	}
}

func TestStateExtendAllTwice(t *testing.T) {
	s := newState(0)
	s2 := s.extendAll([]byte("hello world!"))
	s3 := s.extendAll([]byte("hello world!"))

	if s == s2 || s2 == nil || s2 != s3 || 13 != s.size() {
		t.Logf("size = %d\n", s.size())
		t.Fail()
	}
}

func TestStateExtendAllTwiceSparse(t *testing.T) {
	s := newState(10)
	s2 := s.extendAll([]byte("hello world!"))
	s3 := s.extendAll([]byte("hello world!"))

	if s == s2 || s2 == nil || s2 != s3 || 13 != s.size() {
		t.Logf("size = %d\n", s.size())
		t.Fail()
	}
}

func TestStateExtendMany(t *testing.T) {
	s := newState(0)
	for i := 0; i < 256; i++ {
		s.extend(byte(i))
	}

	if 257 != s.size() {
		t.Logf("size = %d\n", s.size())
		t.Fail()
	}
}

func TestStateExtendManySparse(t *testing.T) {
	s := newState(10)
	for i := 0; i < 256; i++ {
		s.extend(byte(i))
	}

	if 257 != s.size() {
		t.Logf("size = %d\n", s.size())
		t.Fail()
	}
}

func TestTree1(t *testing.T) {
	tree := New()
	tree.Add("hello")
	tree.Add("hi")
	tree.prepare()
	s0 := tree.root
	s1 := s0.edges.get('h')
	s2 := s1.edges.get('e')
	s3 := s2.edges.get('l')
	s4 := s3.edges.get('l')
	s5 := s4.edges.get('o')
	s6 := s1.edges.get('i')

	if s6 == nil {
		t.Log("s6 is nil")
		t.Fail()
	}

	if s0 != s1.fail || s0 != s2.fail || s0 != s3.fail || s0 != s4.fail || s0 != s5.fail || s0 != s6.fail {
		t.Logf("s0 %v", s0)
		t.Logf("s1 fail %v", s1.fail)
		t.Logf("s2 fail %v", s2.fail)
		t.Logf("s3 fail %v", s3.fail)
		t.Logf("s4 fail %v", s4.fail)
		t.Logf("s5 fail %v", s5.fail)
		t.Logf("s6 fail %v", s6.fail)
		t.Fail()
	}

	if l0 := s0.out.Len(); 0 != l0 {
		t.Logf("s0 out len %v", l0)
		t.Fail()
	}
	if l1 := s1.out.Len(); 0 != l1 {
		t.Logf("s1 out len %v", l1)
		t.Fail()
	}
	if l2 := s2.out.Len(); 0 != l2 {
		t.Logf("s2 out len %v", l2)
		t.Fail()
	}
	if l3 := s3.out.Len(); 0 != l3 {
		t.Logf("s3 out len %v", l3)
		t.Fail()
	}
	if l4 := s4.out.Len(); 0 != l4 {
		t.Logf("s4 out len %v", l4)
		t.Fail()
	}
	if l5 := s5.out.Len(); 1 != l5 {
		t.Logf("s5 out len %v", l5)
		t.Fail()
	}
	if l6 := s6.out.Len(); 1 != l6 {
		t.Logf("s out len %v", l6)
		t.Fail()
	}
}

func TestTree2(t *testing.T) {
	tree := New()
	tree.Add("he")
	tree.Add("she")
	tree.Add("his")
	tree.Add("hers")

	if s := tree.root.size(); 10 != s {
		t.Logf("tree size %v", s)
		t.Fail()
	}

	tree.prepare()
	s0 := tree.root
	s1 := s0.edges.get('h')
	s2 := s1.edges.get('e')
	s3 := s0.edges.get('s')
	s4 := s3.edges.get('h')
	s5 := s4.edges.get('e')
	s6 := s1.edges.get('i')

	if s6 == nil {
		t.Log("s6 is nil")
		t.Fail()

		return
	}

	s7 := s6.edges.get('s')
	s8 := s2.edges.get('r')
	s9 := s8.edges.get('s')

	if s0 != s1.fail || s0 != s2.fail || s0 != s3.fail || s0 != s6.fail || s0 != s8.fail ||
		s1 != s4.fail || s2 != s5.fail || s3 != s7.fail || s3 != s9.fail {
		t.Logf("s0 %v", s0)
		t.Logf("s1 %v", s1)
		t.Logf("s2 %v", s2)
		t.Logf("s3 %v", s3)
		t.Logf("s1 fail %v", s1.fail)
		t.Logf("s2 fail %v", s2.fail)
		t.Logf("s3 fail %v", s3.fail)
		t.Logf("s4 fail %v", s4.fail)
		t.Logf("s5 fail %v", s5.fail)
		t.Logf("s6 fail %v", s6.fail)
		t.Logf("s7 fail %v", s7.fail)
		t.Logf("s8 fail %v", s8.fail)
		t.Logf("s9 fail %v", s9.fail)
		t.Fail()
	}

	if l0 := s0.out.Len(); 0 != l0 {
		t.Logf("s0 out len %v", l0)
		t.Fail()
	}
	if l1 := s1.out.Len(); 0 != l1 {
		t.Logf("s1 out len %v", l1)
		t.Fail()
	}
	if l2 := s2.out.Len(); 1 != l2 {
		t.Logf("s2 out len %v", l2)
		t.Fail()
	}
	if l3 := s3.out.Len(); 0 != l3 {
		t.Logf("s3 out len %v", l3)
		t.Fail()
	}
	if l4 := s4.out.Len(); 0 != l4 {
		t.Logf("s4 out len %v", l4)
		t.Fail()
	}
	if l5 := s5.out.Len(); 2 != l5 {
		t.Logf("s5 out len %v", l5)
		t.Fail()
	}
	if l6 := s6.out.Len(); 0 != l6 {
		t.Logf("s6 out len %v", l6)
		t.Fail()
	}
	if l7 := s7.out.Len(); 1 != l7 {
		t.Logf("s7 out len %v", l7)
		t.Fail()
	}
	if l8 := s8.out.Len(); 0 != l8 {
		t.Logf("s8 out len %v", l8)
		t.Fail()
	}
	if l9 := s9.out.Len(); 1 != l9 {
		t.Logf("s9 out len %v", l9)
		t.Fail()
	}
}

func TestSearchSingle(t *testing.T) {
	tree := New()
	tree.Add("apple")
	res := tree.startSearch("washington cut the apple tree")
	if l := res.lastMatchedState.out.Len(); 1 != l {
		t.Logf("search result out len %v", l)
		t.Fail()
	}
	if out := res.out().Front().Value.(string); "apple" != out {
		t.Logf("search result out %v", out)
		t.Fail()
	}
	if 24 != res.lastIndex {
		t.Logf("search result last index %v", res.lastIndex)
		t.Fail()
	}
	if cont := res.continueSearch(); nil != cont {
		t.Logf("search continued %v", cont)
		t.Fail()
	}
}

func TestSearchAdjacent(t *testing.T) {
	tree := New()
	tree.Add("john")
	tree.Add("jane")
	res := tree.startSearch("johnjane")
	cont := res.continueSearch()
	if cont2 := cont.continueSearch(); nil != cont2 {
		t.Logf("search continued %v", cont2)
		t.Fail()
	}
}

func TestSearchEmpty(t *testing.T) {
	tree := New()
	tree.Add("zip")
	tree.Add("zap")
	res := tree.startSearch("")
	if nil != res {
		t.Logf("search results %v", res)
		t.Fail()
	}
}

func TestMultipleOutputs(t *testing.T) {
	tree := New()
	tree.Add("z")
	tree.Add("zz")
	tree.Add("zzz")
	res := tree.startSearch("zzz")
	if 1 != res.lastIndex {
		t.Logf("search result last index %v", res.lastIndex)
		t.Fail()
	}
	if out := res.out().Front().Value.(string); "z" != out {
		t.Logf("search result out %v", out)
		t.Fail()
	}

	res = res.continueSearch()
	if 2 != res.lastIndex {
		t.Logf("search result last index %v", res.lastIndex)
		t.Fail()
	}
	if out := res.out().Front().Value.(string); "zz" != out {
		t.Logf("search result out %v", out)
		t.Fail()
	}

	res = res.continueSearch()
	if 3 != res.lastIndex {
		t.Logf("search result last index %v", res.lastIndex)
		t.Fail()
	}
	if out := res.out().Front().Value.(string); "zzz" != out {
		t.Logf("search result out %v", out)
		t.Fail()
	}

	res = res.continueSearch()
	if nil != res {
		t.Logf("search results %v", res)
		t.Fail()
	}
}

func TestChannel(t *testing.T) {
	tree := New()
	tree.Add("moo")
	tree.Add("one")
	tree.Add("on")
	tree.Add("ne")
	ch := tree.Search("one moon ago")

	if on := <-ch; "on" != on {
		t.Logf("expected 'on' got '%v'", on)
		t.Fail()
	}
	if one := <-ch; "one" != one {
		t.Logf("expected 'one' got '%v'", one)
		t.Fail()
	}
	if ne := <-ch; "ne" != ne {
		t.Logf("expected 'ne' got '%v'", ne)
		t.Fail()
	}
	if moo := <-ch; "moo" != moo {
		t.Logf("expected 'moo' got '%v'", moo)
		t.Fail()
	}
	if on2 := <-ch; "on" != on2 {
		t.Logf("expected 'on' got '%v'", on2)
		t.Fail()
	}
	if x, ok := <-ch; ok {
		t.Logf("expected nothing got '%v'", x)
		t.Fail()
	}
}

func TestText(t *testing.T) {
	text := "The ga3 mutant of Arabidopsis is a gibberellin-responsive dwarf. We present data showing that the ga3-1 mutant is deficient in ent-kaurene oxidase activity, the first cytochrome P450-mediated step in the gibberellin biosynthetic pathway. By using a combination of conventional map-based cloning and random sequencing we identified a putative cytochrome P450 gene mapping to the same location as GA3. Relative to the progenitor line, two ga3 mutant alleles contained single base changes generating in-frame stop codons in the predicted amino acid sequence of the P450. A genomic clone spanning the P450 locus complemented the ga3-2 mutant. The deduced GA3 protein defines an additional class of cytochrome P450 enzymes. The GA3 gene was expressed in all tissues examined, RNA abundance being highest in inflorescence tissue."
	terms := []string{"microsome",
		"cytochrome",
		"cytochrome P450 activity",
		"gibberellic acid biosynthesis",
		"GA3",
		"cytochrome P450",
		"oxygen binding",
		"AT5G25900.1",
		"protein",
		"RNA",
		"gibberellin",
		"Arabidopsis",
		"ent-kaurene oxidase activity",
		"inflorescence",
		"tissue"}
	tree := New()
	for _, term := range terms {
		tree.Add(term)
	}
	found := make(map[string]int)
	for f := range tree.Search(text) {
		found[f] = 1
	}
	if 10 != len(found) {
		t.Log("more found than expected:\n")
		for f := range found {
			t.Logf("=> '%v'", f)
		}
		t.Fail()
	}
}

func TestUnicode(t *testing.T) {
	tree := New()
	tree.Add("iPhone")
	tree.Add("youtube.com/watch")
	tree.Add("profil")
	tree.Add("Заходи, сюда")

	p1 := <-tree.Search("What type of iPhone do you have?")
	p2 := <-tree.Search("Hi, take a look here: https://www.youtube.com/watch?v=RC_6skf1-t")
	p3 := <-tree.Search("Salut est-ce que tu peux aimer ma photo de profil?")
	p4 := <-tree.Search("Привет! У нас сегодня акция. Заходи, сюда узнаешь больше!")
	p5 := <-tree.Search("Yanıtını Beğen Hediye Yolluyor Dene Gör :)")

	expected := "p1: iPhone p2: youtube.com/watch p3: profil p4: Заходи, сюда p5: "
	if s := fmt.Sprintf("p1: %v p2: %v p3: %v p4: %v p5: %v", p1, p2, p3, p4, p5); s != expected {
		t.Logf("expected:\n => '%v'\ngot:\n => '%v'", expected, s)
		t.Fail()
	}
}
