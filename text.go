package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
)

func getOnlyText(text string) string {
	items, err := parseText(text)
	if err != nil {
		log.Errorf("%s\ntext: %q", err, text)
		log.Error(spew.Sdump(items))

		panic(err)
	}

	var out string

	for _, item := range items {
		switch item.typ {
		case itemText, itemRawString:
			out += item.val
		}
	}

	return out
}

func breakLines(text string) string {
	var out string

	for _, l := range strings.Split(text, "\n") {
		// Remove any extra trailing new lines
		l = strings.TrimRight(l, "\n")
		if len(l) < 1 {
			continue
		}

		// TODO: Avoid hardcoding character limit
		if len(getOnlyText(l)) > 42 {
			line := ""
			s := strings.Split(l, " ")
			for i := range s {
				// TODO: Ignore variables and functions when calculating len
				if len(line)+len(s[i]) > 42 {
					log.Debugf("Split! %q from %q", line, l)

					if !strings.HasSuffix(line, "\n") {
						line += "\n"
					}

					out += line

					line = ""
				}

				line += s[i] + " "
			}

			line = strings.TrimRight(line, " ")
			if len(line) > 0 {
				log.Debugf("Split! %q from %q", line, l)

				if !strings.HasSuffix(line, "\n") {
					line += "\n"
				}

				out += line
			}
		} else {
			out += l

			if !strings.HasSuffix(out, "\n") {
				out += "\n"
			}
		}
	}

	return out
}
