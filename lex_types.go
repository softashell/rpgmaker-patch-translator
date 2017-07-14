package main

import "fmt"

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType // The type of this item.
	pos int
	val string // The value of this item.
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	}

	return fmt.Sprintf("{%s: %q}", i.typ, i.val)
}

// itemType identifies the type of lex items.
type itemType int

const eof = -1

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF
	itemText
	itemRawString
	itemLeftDelim
	itemRightDelim
	itemLeftParen
	itemRightParen
	itemParameter
	itemScript
	itemRubyBlock
)

func (t itemType) String() string {
	switch t {
	case itemError:
		return "itemError"
	case itemEOF:
		return "itemEOF"
	case itemText:
		return "itemText"
	case itemRawString:
		return "itemRawString"
	case itemLeftDelim:
		return "itemLeftDelim"
	case itemRightDelim:
		return "itemRightDelim"
	case itemLeftParen:
		return "itemLeftParen"
	case itemRightParen:
		return "itemRightParen"
	case itemParameter:
		return "itemParameter"
	case itemScript:
		return "itemScript"
	case itemRubyBlock:
		return "itemRubyBlock"
	default:
		panic(fmt.Sprintf("unknown item type: %d", t))
	}
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	input string  // the string being scanned
	state stateFn // the next lexing function to enter

	pos        int // current position in the input
	start      int // start position of this item
	width      int // width of last rune read from input
	lastPos    int // position of most recent item returned by nextItem
	parenDepth int

	mark rune // The current lexed rune

	items chan item // channel of scanned items
}

type lexPosition struct {
	pos   int
	start int
	width int
	mark  rune
}

func saveLexerPosition(lexState *lexer) *lexPosition {
	return &lexPosition{
		pos:   lexState.pos,
		start: lexState.start,
		width: lexState.width,
		mark:  lexState.mark,
	}
}

func (l *lexPosition) restore(lexState *lexer) {
	lexState.pos = l.pos
	lexState.start = l.start
	lexState.width = l.width
	lexState.mark = l.mark
}
