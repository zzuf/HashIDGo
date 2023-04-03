package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//go:embed prototypes.json
var prototypesJsonStr string

type Mode struct {
	John     string `json:"john"`
	Hashcat  int    `json:"hashcat"`
	Extended bool   `json:"extended"`
	Name     string `json:"name"`
}

type Rule struct {
	Regex string `json:"regex"`
	Modes []Mode `json:"modes"`
}

func patternBuilder(base string, start int, end int) string {
	retVal := ""
	if start == -1 {
		for i := 0; i < (end / 1000); i++ {
			retVal += base + "{1000}"
		}
		if end%1000 != 0 {
			retVal += base + fmt.Sprintf("{%d}", end%1000)
		}
	} else {
		retVal = base + fmt.Sprintf("{%d}", start)
		for i := 0; i < (end / 1000); i++ {
			retVal += base + "{0,1000}"
		}
		if end%1000 != 0 {
			retVal += base + fmt.Sprintf("{0,%d}", end%1000)
		}
	}
	return retVal
}

func replaceSyntax(pattern string) string {
	r1 := regexp.MustCompile(`\{[0-9]{4}\}|\{[0-9]{1,3},[0-9]{4}\}`)
	r2 := regexp.MustCompile(`(\[[a-zA-Z0-9\-]+\])\{([0-9]{4})\}`)
	r3 := regexp.MustCompile(`(\[[a-zA-Z0-9\-]+\])\{([0-9]{1,3},[0-9]{4})\}`)
	matched := r1.MatchString(pattern)
	if matched {
		res := r2.FindAllStringSubmatch(pattern, -1)
		if len(res) > 0 {
			end, _ := strconv.Atoi(res[0][2])
			if end == 1000 {
				return pattern
			}
			return strings.Replace(pattern, res[0][0], patternBuilder(res[0][1], -1, end), 1)
		}
		res = r3.FindAllStringSubmatch(pattern, -1)
		resSplit := strings.Split(res[0][2], ",")
		start, _ := strconv.Atoi(resSplit[0])
		end, _ := strconv.Atoi(resSplit[1])
		if end == 1000 {
			return pattern
		}
		return strings.Replace(pattern, res[0][0], patternBuilder(res[0][1], start, end), 1)
	} else {
		return pattern
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: ./hashidgo <Hash>")
		return
	}
	input := []byte(os.Args[1])

	jsonStr := strings.Replace(prototypesJsonStr, "\"john\": null", "\"john\": \"\"", -1)
	jsonStr = strings.Replace(jsonStr, "\"hashcat\": null", "\"hashcat\": -1", -1)

	var rules []Rule
	err := json.Unmarshal([]byte(jsonStr), &rules)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, rule := range rules {
		replaced := replaceSyntax(rule.Regex)
		re := regexp.MustCompile("(?i)" + replaced)
		if re.Match(input) {
			for _, mode := range rule.Modes {
				fmt.Printf("john: %s, hashcat: %v, extended: %v, name: %s\n",
					mode.John, mode.Hashcat, mode.Extended, mode.Name)
			}
			return
		}
	}

	fmt.Println("No matching rule found")
}
