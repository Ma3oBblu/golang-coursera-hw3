package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	seenBrowsers := make(map[string]struct{})
	var foundUsers string

	scanner := bufio.NewScanner(file)
	i := 0
	user := make(map[string]interface{})
	for scanner.Scan() {
		i++
		err = json.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			panic(err)
		}

		var isAndroid, isMSIE bool
		browsers := user["browsers"].([]interface{})
		for _, browserRaw := range browsers {
			browser := browserRaw.(string)
			if strings.Contains(browser, "Android") {
				isAndroid = true
				seenBrowsers[browser] = struct{}{}
			} else if strings.Contains(browser, "MSIE") {
				isMSIE = true
				seenBrowsers[browser] = struct{}{}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		//log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.ReplaceAll(user["email"].(string), "@", " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i-1, user["name"], email)
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
