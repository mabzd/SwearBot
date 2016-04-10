package dictmatch

import (
	"testing"
)

func TestNewDict(t *testing.T) {
	dict := NewDict()
	if dict == nil {
		t.Fatal("Dict is nil")
	}
}

func TestMatchEmpty(t *testing.T) {
	dict := NewDict()
	assertNotMatch(t, dict, "")
}

func TestMatchNotExisting(t *testing.T) {
	dict := NewDict()
	assertNotMatch(t, dict, "abc")
}

func TestMatch(t *testing.T) {
	dict := NewDict()
	assertAddEntry(t, dict, "abc")

	assertNotMatch(t, dict, "a")
	assertNotMatch(t, dict, "ab")
	assertNotMatch(t, dict, "abcd")
	assertMatch(t, dict, "abc", "abc")
}

func TestMatchUnicode(t *testing.T) {
	dict := NewDict()
	assertAddEntry(t, dict, "ταБЬℓσ")

	assertNotMatch(t, dict, "ταБЬℓ")
	assertNotMatch(t, dict, "ταБЬℓ*")
	assertMatch(t, dict, "ταБЬℓσ", "ταБЬℓσ")
}

func TestMatchWildcard(t *testing.T) {
	dict := NewDict()
	assertAddEntry(t, dict, "abc*")

	assertNotMatch(t, dict, "a")
	assertNotMatch(t, dict, "ab")
	assertMatch(t, dict, "abc", "abc")
	assertMatch(t, dict, "abcd", "abc")
	assertMatch(t, dict, "abcdefghijkl", "abc")
}

func TestMatchUnicodeWildcard(t *testing.T) {
	dict := NewDict()
	assertAddEntry(t, dict, "ταБЬℓσ*")

	assertNotMatch(t, dict, "ταБЬℓ")
	assertMatch(t, dict, "ταБЬℓσ", "ταБЬℓσ")
	assertMatch(t, dict, "ταБЬℓσ*", "ταБЬℓσ")
	assertMatch(t, dict, "ταБЬℓσБЬasd", "ταБЬℓσ")
}

func TestMatchMultiple(t *testing.T) {
	dict := NewDict()
	assertAddEntry(t, dict, "a")
	assertAddEntry(t, dict, "aa")
	assertAddEntry(t, dict, "ab")
	assertAddEntry(t, dict, "ac*")

	assertMatch(t, dict, "a", "a")
	assertMatch(t, dict, "aa", "aa")
	assertMatch(t, dict, "ab", "ab")
	assertMatch(t, dict, "ac", "ac")
	assertMatch(t, dict, "acc", "ac")
	assertNotMatch(t, dict, "aaa")
	assertNotMatch(t, dict, "abb")
	assertNotMatch(t, dict, "abc")
}

func TestAddDuplicatedEntry(t *testing.T) {
	dict := NewDict()
	assertAddEntry(t, dict, "abc")
	assertAddEntryError(t, dict, "abc", WordExistErr)
}

func TestAddWordOverlappedByWildcard(t *testing.T) {
	dict := NewDict()
	assertAddEntry(t, dict, "abc*")
	assertAddEntryError(t, dict, "abcdef", WordOverlappedByWildcardErr)
}

func TestAddWildcardOverlappedByWord(t *testing.T) {
	dict := NewDict()
	assertAddEntry(t, dict, "abcdef")
	assertAddEntryError(t, dict, "abc*", WildcardOverlappedByWordErr)
}

func TestAddDuplicatedWildcardRoot(t *testing.T) {
	dict := NewDict()
	assertAddEntry(t, dict, "abc")
	assertAddEntryError(t, dict, "abc*", WildcardRootExistErr)
}

func TestAddDuplicatedWildcard(t *testing.T) {
	dict := NewDict()
	assertAddEntry(t, dict, "abc*")
	assertAddEntryError(t, dict, "abc*", WordOverlappedByWildcardErr)
}

func TestInvalidWildcardPlacement(t *testing.T) {
	dict := NewDict()
	assertAddEntryError(t, dict, "**a", InvalidWildardPlacementErr)
	assertAddEntryError(t, dict, "*a", InvalidWildardPlacementErr)
	assertAddEntryError(t, dict, "a*b", InvalidWildardPlacementErr)
	assertAddEntryError(t, dict, "a**", InvalidWildardPlacementErr)
	assertAddEntryError(t, dict, "**", InvalidWildardPlacementErr)
}

func assertAddEntry(t *testing.T, dict *Dict, word string) {
	err := dict.AddEntry(word)
	if err != nil {
		t.Fatal(err)
	}
}

func assertAddEntryError(t *testing.T, dict *Dict, word string, errType int) {
	err := dict.AddEntry(word)
	if err == nil {
		t.Fatalf("Adding entry '%s' should yield error type %d (no error found)", word, errType)
	}
	if err.ErrType != errType {
		t.Fatalf(
			"Adding entry '%s' should yield error type %d (error %d found instead: '%s')",
			word,
			errType,
			err.ErrType,
			err.Desc)
	}
}

func assertNotMatch(t *testing.T, dict *Dict, word string) {
	success, _ := dict.Match(word)
	if success {
		t.Fatalf("'%s' should not be matched", word)
	}
}

func assertMatch(t *testing.T, dict *Dict, word string, expectedMatch string) {
	success, match := dict.Match(word)
	if !success {
		t.Fatalf("'%s' should be matched", word)
	}
	if match != expectedMatch {
		t.Fatalf("Actual match '%s' is not equal to expected match '%s'", match, expectedMatch)
	}
}
