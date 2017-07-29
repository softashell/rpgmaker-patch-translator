package lex

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	log "github.com/Sirupsen/logrus"
)

const (
	slashCharacters = "abcdefghijklmnopqrstuvxzwyABCDEFGHIJKLMNOPQRSTUVXZWY0123456789[]{}()\\/<>!|$^."
	rawCharacters   = "\u3000\t\n・･！？。…「」『』()（）/\"“”[]【】<>〈〉：:*＊_＿#$%="
	// Skipping these might not actually be such a good idea since in some cases translator will lack context
	ignoredCharacters = "abcdefghijklmnopqrstuvxzwyABCDEFGHIJKLMNOPQRSTUVXZWYａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｘｚｗｙＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＸＺＷＹ.,!?" // + " "
	numbers           = "0123456789０１２３４５６７８９"
	numberEndings     = "つ十百千万"
	numberAdditions   = "%％階"
)

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

func (l *lexer) nextItem() Item {
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

// emit passes an Item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- Item{l.start, t, l.input[l.start:l.pos]}

	l.start = l.pos
}

func (l *lexer) emitBefore(t itemType) {
	l.backup(1)

	if l.pos > l.start {
		l.emit(t)

		l.ignore()
	}
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
	l.items <- Item{l.start, ItemError, fmt.Sprintf(format, args...)}

	return nil
}

// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan Item),
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
		case r == '%' && l.peek(1) == 's':
			l.emitBefore(ItemText)

			l.next()

			if l.accept("s") {
				l.emit(ItemRawString)
			}
		case (r == 'i' && strings.HasPrefix(l.input[l.pos:], "f(")) ||
			(r == 'e' && strings.HasPrefix(l.input[l.pos:], "n(")):
			l.emitBefore(ItemText)

			return lexScript
		case r == '\\' || r == '@':
			l.emitBefore(ItemText)

			return lexScript
		case r == '#' && l.peek(1) == '{':
			l.emitBefore(ItemText)
			return lexRubyBlock
		case strings.ContainsRune(numbers, r):
			l.emitBefore(ItemText)
			return lexNumber
		case strings.ContainsRune(rawCharacters, r) || unicode.IsSymbol(r):
			l.emitBefore(ItemText)
			l.next()
			l.emit(ItemRawString)

			return lexText
		case r == '-' && l.peek(1) == '-':
			l.emitBefore(ItemText)
			l.acceptRun("-")
			l.emit(ItemRawString)
		case strings.ContainsRune(ignoredCharacters, r):
			l.emitBefore(ItemText)

			if !strings.Contains(l.input[l.pos:], "if(") && !strings.Contains(l.input[l.pos:], "en(") {
				l.acceptRun(ignoredCharacters + " ")
			} else {
				l.acceptRun(ignoredCharacters)
			}

			l.emit(ItemRawString)

			return lexText
		}

	}

	if l.pos > l.start {
		l.emit(ItemText)

		l.ignore()
	}

	l.emit(ItemEOF)

	return nil

}

func lexScript(l *lexer) stateFn {
	log.Debugf("lexScript %q", l.input[l.pos:])

Loop:
	for {
		switch l.next() {
		case eof:
			log.Warn("Script not terminated properly %q", l.input[l.start:])
			break Loop
		case '(':
			l.emitBefore(ItemScript)

			l.next()

			l.emit(ItemLeftParen)

			l.parenDepth++

			return lexInsideAction
		case '[':
			l.backup(1)
			l.emit(ItemScript)

			return lexLeftDelim
		case '\\':
			l.acceptRun(slashCharacters)

			break Loop
		case '\n':
			l.acceptRun("[0123456789]")

			break Loop
		case '@':
			l.acceptRun("0123456789-")

			break Loop

		default:
			log.Debug(string(l.mark))
		}
	}

	l.emit(ItemScript)

	return lexText
}

func lexLeftDelim(l *lexer) stateFn {
	l.next()

	log.Debug("leftDelim: ", string(l.mark))

	l.emit(ItemLeftDelim)

	return lexInsideAction
}

func lexRightDelim(l *lexer) stateFn {
	l.next()

	log.Debug("rightDelim: ", string(l.mark))

	l.emit(ItemRightDelim)

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
	case r == eof:
		return l.errorf("unclosed action")
	case r == '(':
		l.emitBefore(ItemParameter)

		l.next()

		l.emit(ItemLeftParen)

		l.parenDepth++
	case r == ')':
		l.emitBefore(ItemParameter)

		l.next()

		l.emit(ItemRightParen)

		l.parenDepth--

		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren %#U", r)
		}

		if l.parenDepth == 0 {
			return lexText
		}
	case r == '[':
		l.emitBefore(ItemParameter)

		return lexLeftDelim
	case r == ']':
		l.emitBefore(ItemParameter)

		return lexRightDelim
	case r == '%' && l.accept("s"):
		l.emit(ItemParameter)
	case r == '"':
		l.emit(ItemParameter)

		return lexText
	}

	return lexInsideAction
}

func lexRubyBlock(l *lexer) stateFn {
	log.Debugf("lexRubyBlock %q", l.input[l.pos:])

	opened := 0

Loop:
	for {
		switch l.next() {
		case eof:
			log.Warn("Ruby block not terminated properly %q", l.input[l.start:])
			break Loop
		case '#':
			log.Debug("Starting ruby block")
		case '{':
			log.Debug("Opening brackets")
			opened++
		case '}':
			log.Debug("Closing brackets")
			opened--
			if opened <= 0 {
				log.Debug("Ending ruby block")
				break Loop
			}
		default:
			log.Debug(string(l.mark))
		}
	}

	l.emit(ItemRubyBlock)

	return lexText
}

func lexNumber(l *lexer) stateFn {
	log.Debugf("lexNumber %q", l.input[l.pos:])

Loop:
	for {
		switch r := l.next(); {
		case r == eof:
			break Loop
		case strings.ContainsRune(numbers, r):
			l.acceptRun(numbers)
		case strings.ContainsRune(numberEndings, r) || strings.ContainsRune(numberAdditions, r):
			break Loop // Shouldn't have more than one of these
		default:
			l.backup(1)
			break Loop
		}
	}

	l.emit(ItemNumber)

	return lexText
}
