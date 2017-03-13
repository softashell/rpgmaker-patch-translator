package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
)

func getOnlyText(text string) string {
	items := parseText(text)

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

		// TODO: Avoid hardcoding character limit
		if len(getOnlyText(l)) > 42 {
			line := ""
			s := strings.Split(l, " ")
			for i := range s {
				// TODO: Ignore variables and functions when calculating len
				if len(line)+len(s[i]) > 42 {
					log.Infof("Split! %q from %q", line, l)

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
				log.Infof("Split! %q from %q", line, l)

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
