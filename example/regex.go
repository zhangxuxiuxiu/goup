package util

import (
	"fmt"
	"regexp"
	"testing"
)

// ref:https://blog.logrocket.com/deep-dive-regular-expressions-golang/
func TestRegex(*testing.T) {
	{ //MatchString
		inputText := "I love new york city"
		match, err := regexp.MatchString("[A-z]ork", inputText)
		if err == nil {
			fmt.Println("Match:", match)
		} else {
			fmt.Println("Error:", err)
		}
	}

	{ //FindStringIndex
		pattern := regexp.MustCompile("H[a-z]{4}|[A-z]ork")
		welcomeMessage := "Hello guys, welcome to new york city"
		firstMatchIndex := pattern.FindStringIndex(welcomeMessage)
		fmt.Println("First matched index", firstMatchIndex[0], "-", firstMatchIndex[1],
			"=", welcomeMessage[firstMatchIndex[0]:firstMatchIndex[1]])
	}

	{ //FindString
		pattern := regexp.MustCompile("H[a-z]{4}|[A-z]ork")
		welcomeMessage := "Hello guys, welcome to new york city"
		firstMatchSubstring := pattern.FindString(welcomeMessage)
		fmt.Println("First matched substring:", firstMatchSubstring)
	}

	{ //FindAllString
		pattern := regexp.MustCompile("H[a-z]{4}|[A-z]ork")
		welcomeMessage := "Hello guys, welcome to new york city"
		allSubstringMatches := pattern.FindAllString(welcomeMessage, -1)
		fmt.Println(allSubstringMatches)
	}

	{ //FindAllStringIndex
		pattern := regexp.MustCompile("H[a-z]{4}|[A-z]ork")
		welcomeMessage := "Hello guys, welcome to new york city"
		allIndices := pattern.FindAllStringIndex(welcomeMessage, -1)
		fmt.Println(allIndices)
		for i, idx := range allIndices {
			fmt.Println("Index", i, "=", idx[0], "-",
				idx[1], "=", welcomeMessage[idx[0]:idx[1]])
		}
	}

	{ //MatchString
		pattern, compileErr := regexp.Compile("[A-z]ork")
		correctAnswer := "Yes, I love new york city"
		question := "Do you love new york city?"
		wrongAnswer := "I love dogs"
		if compileErr == nil {
			fmt.Println("Question:", pattern.MatchString(question))
			fmt.Println("Correct Answer:", pattern.MatchString(correctAnswer))
			fmt.Println("Wrong Answer:", pattern.MatchString(wrongAnswer))
		} else {
			fmt.Println("Error:", compileErr)
		}
	}

	{ //Split
		pattern := regexp.MustCompile("guys|york")
		welcomeMessage := "Hello guys, welcome to new york city"
		split := pattern.Split(welcomeMessage, -1)
		for _, s := range split {
			if s != "" {
				fmt.Println("Substring:", s)
			}
		}
		/*
			Substring: Hello
			Substring: , welcome to new
			Substring:  city
		*/
	}

	{ //FindStringSubmatch
		pattern := regexp.MustCompile("welcome ([A-z]*) new ([A-z]*) city")
		welcomeMessage := "Hello guys, welcome to new york city"
		subStr := pattern.FindStringSubmatch(welcomeMessage)
		fmt.Println(subStr)

		for _, s := range subStr {
			fmt.Println("Match:", s)
		}
		/*
			Match: welcome to new york city
			Match: to
			Match: york
		*/
	}

	{ //SubexpNames
		pattern := regexp.MustCompile("welcome (?P<val1>[A-z]*) new (?P<val2>[A-z]*) city")
		replaced := pattern.SubexpNames()
		replacedCnt := pattern.NumSubexp()

		fmt.Println(replaced, replacedCnt)
		/*
			[ val1 val2] 2
		*/
	}

	{ //ReplaceAllString
		pattern := regexp.MustCompile("welcome (?P<val1>[A-z]*) new (?P<val2>[A-z]*) city")
		welcomeMessage := "Hello guys, welcome to new york city"
		template := "(value 1: ${val1}, value 2: ${val2})"
		replaced := pattern.ReplaceAllString(welcomeMessage, template)
		fmt.Println(replaced)
		/*
			Hello guys, (value 1: to, value 2: york)
		*/
	}

	{ //ReplaceAllLiteralString
		pattern := regexp.MustCompile("welcome (?P<val1>[A-z]*) new (?P<val2>[A-z]*) city")
		welcomeMessage := "Hello guys, welcome to new york city"
		template := "(value 1: ${val1}, value 2: ${val2})"
		replaced := pattern.ReplaceAllLiteralString(welcomeMessage, template)
		fmt.Println(replaced)
		/*
			Hello guys, (value 1: ${val1}, value 2: ${val2})
		*/
	}

	{ //ReplaceAllStringFunc
		pattern := regexp.MustCompile("welcome ([A-z]*) new ([A-z]*) city")
		welcomeMessage := "Hello guys, welcome to new york city"
		replaced := pattern.ReplaceAllStringFunc(welcomeMessage, func(s string) string {
			return "here is the replacement content for the matched string."
		})
		fmt.Println(replaced)
		/*
			Hello guys, here is the replacement content for the matched string
		*/
	}
}
