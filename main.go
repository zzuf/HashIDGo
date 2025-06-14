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

//go:embed hashID/prototypes.json
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

func buildPattern(patternBase string, min int, max int) string {
	result := ""
	if min == -1 {
		for i := 0; i < (max / 1000); i++ {
			result += patternBase + "{1000}"
		}
		if max%1000 != 0 {
			result += patternBase + fmt.Sprintf("{%d}", max%1000)
		}
	} else {
		result = patternBase + fmt.Sprintf("{%d}", min)
		for i := 0; i < (max / 1000); i++ {
			result += patternBase + "{0,1000}"
		}
		if max%1000 != 0 {
			result += patternBase + fmt.Sprintf("{0,%d}", max%1000)
		}
	}
	return result
}

func expandPatternSyntax(pattern string) string {
	r1 := regexp.MustCompile(`\{[0-9]{4}\}|\{[0-9]{1,3},[0-9]{4}\}`)
	r2 := regexp.MustCompile(`(\[[a-zA-Z0-9\-]+\])\{([0-9]{4})\}`)
	r3 := regexp.MustCompile(`(\[[a-zA-Z0-9\-]+\])\{([0-9]{1,3},[0-9]{4})\}`)
	matched := r1.MatchString(pattern)
	if matched {
		res := r2.FindAllStringSubmatch(pattern, -1)
		if len(res) > 0 {
			max, _ := strconv.Atoi(res[0][2])
			if max == 1000 {
				return pattern
			}
			return strings.Replace(pattern, res[0][0], buildPattern(res[0][1], -1, max), 1)
		}
		res = r3.FindAllStringSubmatch(pattern, -1)
		resSplit := strings.Split(res[0][2], ",")
		min, _ := strconv.Atoi(resSplit[0])
		max, _ := strconv.Atoi(resSplit[1])
		if max == 1000 {
			return pattern
		}
		return strings.Replace(pattern, res[0][0], buildPattern(res[0][1], min, max), 1)
	} else {
		return pattern
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: ./hashidgo <Hash>")
		return
	}
	hashInput := []byte(os.Args[1])

	normalizedJsonStr := strings.Replace(prototypesJsonStr, "\"john\": null", "\"john\": \"\"", -1)
	normalizedJsonStr = strings.Replace(normalizedJsonStr, "\"hashcat\": null", "\"hashcat\": -1", -1)

	var ruleList []Rule
	err := json.Unmarshal([]byte(normalizedJsonStr), &ruleList)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, ruleItem := range ruleList {
		expandedRegex := expandPatternSyntax(ruleItem.Regex)
		regex := regexp.MustCompile("(?i)" + expandedRegex)
		if regex.Match(hashInput) {
			for _, modeItem := range ruleItem.Modes {
				fmt.Printf("john: %s, hashcat: %v, extended: %v, name: %s\n",
					modeItem.John, modeItem.Hashcat, modeItem.Extended, modeItem.Name)
			}
			return
		}
	}

	fmt.Println("No matching rule found")
}
