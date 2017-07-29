package lex

import "fmt"

type Item struct {
	pos int

	Typ itemType // The type of this item.
	Val string   // The value of this item.
}

func (i Item) String() string {
	switch {
	case i.Typ == ItemEOF:
		return "EOF"
	case i.Typ == ItemError:
		return i.Val
	}

	return fmt.Sprintf("{%s: %q}", i.Typ, i.Val)
}

// itemType identifies the type of lex items.
type itemType int

const eof = -1

const (
	ItemError itemType = iota // error occurred; value is text of error
	ItemEOF
	ItemText
	ItemRawString
	ItemLeftDelim
	ItemRightDelim
	ItemLeftParen
	ItemRightParen
	ItemParameter
	ItemScript
	ItemRubyBlock
	ItemNumber
)

func (t itemType) String() string {
	switch t {
	case ItemError:
		return "itemError"
	case ItemEOF:
		return "itemEOF"
	case ItemText:
		return "itemText"
	case ItemRawString:
		return "itemRawString"
	case ItemLeftDelim:
		return "itemLeftDelim"
	case ItemRightDelim:
		return "itemRightDelim"
	case ItemLeftParen:
		return "itemLeftParen"
	case ItemRightParen:
		return "itemRightParen"
	case ItemParameter:
		return "itemParameter"
	case ItemScript:
		return "itemScript"
	case ItemRubyBlock:
		return "itemRubyBlock"
	case ItemNumber:
		return "itemNumber"
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

	items chan Item // channel of scanned items
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
