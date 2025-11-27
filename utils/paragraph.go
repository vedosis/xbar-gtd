package utils

import "strings"

func WrapLines(s string, maxLength uint16) []string {
	if maxLength <= 0 {
		return []string{s}
	}

	var result []string

	paragraphs := strings.Split(s, "\n\n")

	for paragraphIdx, paragraph := range paragraphs {
		words := strings.Fields(paragraph)
		currentLine := ""

		for _, word := range words {
			wordLength := uint16(len(word))
			if wordLength > maxLength {
				if len(currentLine) > 0 {
					result = append(result, currentLine)
					currentLine = ""
				}

				for wordLength > maxLength {
					result = append(result, word[:maxLength])
					word = word[maxLength:]
				}

				if wordLength > 0 {
					currentLine = word
				}
				continue
			}

			if len(currentLine) == 0 {
				currentLine = word
			} else if len(currentLine)+1+int(wordLength) <= int(maxLength) {
				currentLine += " " + word
			} else {
				result = append(result, currentLine)
				currentLine = word
			}
		}

		if len(currentLine) > 0 {
			result = append(result, currentLine)
		}

		if paragraphIdx < len(paragraphs)-1 {
			result = append(result, "")
		}
	}

	return result
}
