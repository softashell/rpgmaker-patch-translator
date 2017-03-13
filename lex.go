package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	log "github.com/Sirupsen/logrus"
)

const eof = -1

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])

	l.mark = r
	l.width = w

	l.pos += l.width

	return r
}

func (l *lexer) nextItem() item {
	item := <-l.items

	l.lastPos = item.pos

	return item
}

func (l *lexer) backup(pos int) {
	for i := 0; i < pos; i++ {
		l.pos -= l.width
		if l.pos < 0 {
			l.pos = 0
		}

		r, w := utf8.DecodeRuneInString(l.input[l.pos:])
		l.mark = r
		l.width = w
	}
}

func (l *lexer) peek(locs int) rune {
	pos := saveLexerPosition(l)

	var r rune

	x := 0

	for x < locs {
		l.next()

		if x == locs-1 {
			r = l.mark
		}

		x++
	}

	pos.restore(l)

	return r
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}

	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}

	l.backup(1)

	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}

	l.backup(1)
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}

	return nil
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) drain() {
	for range l.items {
	}
}

// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}

	go l.run()

	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}

	close(l.items)
}

func lexText(l *lexer) stateFn {
	log.Debugf("lexText %q", l.input[l.pos:])

	l.width = 0

Loop:
	for {
		//switch l.next() {
		switch r := l.next(); {
		case r == eof:
			break Loop
		case r == '"':
			if l.pos > l.start {
				l.backup(1)

				if len(l.input[l.start:l.pos]) > 0 {
					l.emit(itemText)

					l.ignore()
				}
			}

			l.next()

			return lexInsideAction
		case r == '%':
			if l.pos > l.start {
				l.backup(1)

				if len(l.input[l.start:l.pos]) > 0 {
					l.emit(itemText)

					l.ignore()
				}
			}

			l.next()

			if r = l.next(); r != eof && r != '\n' {
				switch r {
				case 's':
					l.emit(itemRawString)
				}
			}
		case r == 'i' && strings.HasPrefix(l.input[l.pos:], "f("):
			l.backup(1)

			if l.pos > l.start {
				l.emit(itemText)
			}

			l.ignore()

			return lexScript
		case r == 'e' && strings.HasPrefix(l.input[l.pos:], "n("):
			l.backup(1)

			if l.pos > l.start {
				l.emit(itemText)
			}

			l.ignore()

			return lexScript
		case r == '\\':
			if r = l.peek(1); r != eof && r != '\n' && r == '\\' {
				l.backup(1)

				if l.pos > l.start {
					l.emit(itemText)
				}

				l.ignore()

				return lexScript
			}
		case strings.ContainsRune("\u3000。…【】」「\n()", r) || unicode.IsSymbol(r):
			if l.pos > l.start {
				l.backup(1)

				if len(l.input[l.start:l.pos]) > 0 {
					l.emit(itemText)

					l.ignore()
				}

			}

			l.next()

			l.emit(itemRawString)

			return lexText
		}
	}

	if l.pos > l.start {
		l.emit(itemText)
	}

	l.ignore()

	l.emit(itemEOF)

	return nil

}

func lexScript(l *lexer) stateFn {
	if l.pos < len(l.input) {
		log.Debugf("lexScript %q", l.input[l.pos:])
	} else {
		log.Debugf("lexScript out of bounds %d/%d", l.pos, len(l.input))
	}

Loop:
	for {
		switch l.next() {
		case eof, '\n':
			return l.errorf("unterminated script")
		case '(':
			l.emit(itemLeftParen)

			l.parenDepth++

			return lexInsideAction
		case '[':
			l.backup(1)
			l.emit(itemScript)
			return lexLeftDelim
		case '\\':
			if r := l.peek(1); r != eof && r != '\n' {
				if strings.ContainsRune(">lrt{}$G.|^", r) {
					l.next()

					log.Debug("Found escaped ", string(r))
					break Loop
				}
			}

			fallthrough
		default:
			if l.pos < len(l.input) {
				log.Debug(string(l.input[l.pos]))
			}
		}
	}

	l.emit(itemScript)

	return lexText
}

func lexLeftDelim(l *lexer) stateFn {
	log.Debug("leftDelim: ", l.input[l.pos:l.pos+len("[")])

	l.pos += len("[")

	l.emit(itemLeftDelim)
	l.ignore()
	//l.parenDepth = 0

	return lexInsideAction
}

func lexRightDelim(l *lexer) stateFn {
	log.Debug("rightDelim: ", l.input[l.pos:l.pos+len("]")])

	l.pos += len("]")

	l.emit(itemRightDelim)

	log.Debug("Paren depth: ", l.parenDepth)

	if l.parenDepth == 0 {
		return lexText
	}

	return lexInsideAction
}

// lexInsideAction scans the elements inside action delimiters.
func lexInsideAction(l *lexer) stateFn {
	log.Debugf("lexInsideAction %q", l.input[l.pos:])

	switch r := l.next(); {
	case r == eof || isEndOfLine(r):
		return l.errorf("unclosed action")
	case r == '(':
		l.emit(itemLeftParen)
		l.parenDepth++
	case r == ')':
		l.emit(itemRightParen)
		l.parenDepth--

		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren %#U", r)
		}

		if l.parenDepth == 0 {
			return lexText
		}
	case r == '[':
		l.backup(1)

		if l.pos > l.start {
			l.emit(itemParameter)
		}

		l.ignore()

		return lexLeftDelim
	case r == ']':
		l.backup(1)

		if l.pos > l.start {
			l.emit(itemParameter)
		}

		l.ignore()

		return lexRightDelim
	case r == '%':
		if r = l.next(); r != eof && r != '\n' {
			switch r {
			case 's':
				l.emit(itemParameter)
			}
		}
	case r == '"':
		l.emit(itemParameter)

		return lexText
	}

	return lexInsideAction
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}
