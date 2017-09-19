package shellexec

import (
	"regexp"
	"runtime"
	"strings"
)

var isWindows = runtime.GOOS == "windows"

func PrepareCommand(cmd string) string {
	ret := ""
	lineEnding := regexp.MustCompile("\r\n|\n")
	continueStart := regexp.MustCompile(`^\s*&&\s*`)
	continueEnd := regexp.MustCompile(`[\\]\s*$`)
	lines := lineEnding.Split(cmd, -1)

	if isWindows {
		for _, line := range lines {
			if continueEnd.MatchString(line) {
				line = continueEnd.ReplaceAllString(line, "")
			} else {
				line += " && "
			}
			ret += strings.TrimSpace(line) + " "
		}
		ret = strings.TrimSpace(ret)
		if strings.HasSuffix(ret, " &&") {
			ret = ret[0 : len(ret)-3]
		}
	} else {
		isContinuing := false
		for i, line := range lines {
			if i == 0 {
				isContinuing = false
				if !continueEnd.MatchString(line) {
					line += " \\"
					isContinuing = true
				}
			} else {
				if isContinuing && !continueStart.MatchString(line) {
					line = " && " + line
				}
				isContinuing = false
				if !continueEnd.MatchString(line) {
					line += " \\"
					isContinuing = true
				}
			}
			line += "\n"
			ret += line
		}
		if len(ret) > 3 {
			ret = ret[0 : len(ret)-3]
		}
	}
	return ret
}
