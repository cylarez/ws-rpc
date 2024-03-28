package main

import (
	"log"
	"os"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type messageFromProto struct {
	desc      string
	fromProto string
}

func kebabCase(str string) string {
	k := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	k = matchAllCap.ReplaceAllString(k, "${1}-${2}")
	return strings.ToLower(k)
}

func cleanChar(str string) string {
	return strings.ReplaceAll(str, ".", "")
}

func getMainPck(msg string) string {
	lastDot := strings.LastIndex(msg, ".")
	if lastDot == -1 {
		return ""
	}
	return msg[0:lastDot]
}

func findMessageDesc(msgName, protoFile string) (string, string) {
	return findDesc("message", msgName, protoFile)
}

func findDesc(t, msg, protoFile string) (string, string) {
	if msg == "Empty" || (len(msg) > 7 && msg[0:7] == "google.") || protoFile == "" {
		return "", ""
	}
	// Check if message belongs to another package (start with main package name)
	if strings.HasPrefix(msg, mainPckName) {
		split := strings.Split(msg, ".")
		b, err := os.ReadFile(protoPath + split[2] + ".proto")
		if err != nil {
			log.Fatal(err)
		}
		res, _ := findDesc(t, split[len(split)-1], string(b))
		return res, string(b)
	}
	i := strings.Index(protoFile, t+" "+msg)
	if i == -1 {
		if t != "enum" {
			return findDesc("enum", msg, protoFile)
		}
		return "", ""
	}
	end := strings.Index(protoFile[i:], "}")
	return protoFile[i : i+end+1], protoFile
}

func getRpcOptions(route, protoFile string) string {
	protoFile = strings.ReplaceAll(protoFile, " ", "")
	i := strings.Index(protoFile, "rpc"+route+"(")
	protoFile = protoFile[i:]
	end := strings.Index(protoFile, "}")
	return protoFile[0 : end+1]
}

func getScopes(options string) (scopes []string) {
	for _, ct := range clientTypes {
		pos := strings.Index(options, ct)
		if pos != -1 {
			scopes = append(scopes, ct)
		}
	}
	return
}

func isPushEvent(options string) bool {
	pos := strings.Index(options, "pushEvent=true")
	return pos != -1
}

func isPublic(options string) bool {
	hasPublic := strings.Index(options, "public=true")
	if hasPublic != -1 {
		return true
	}
	hasLogin := strings.Index(options, "login=true")
	return hasLogin != -1
}

func isDebug(options string) bool {
	options = strings.ReplaceAll(options, " ", "")
	pos := strings.Index(options, "debug=true")
	return pos != -1
}

// Find message references to generate package imports
func findAllRefs(desc ...messageFromProto) (results map[string]string) {
	var subRefs []messageFromProto
	results = map[string]string{}
	messages := findMessages(desc...)
	for _, m := range messages {
		d, from := findMessageDesc(m.desc, m.fromProto)
		if d != "" {
			results[m.desc] = d
			subRefs = append(subRefs, messageFromProto{desc: d, fromProto: from})
		}
	}
	if len(subRefs) > 0 {
		subResults := findAllRefs(subRefs...)
		for k, v := range subResults {
			results[k] = v
		}
	}
	return
}

// Find usage of messages from outside current proto module
func findMessages(contents ...messageFromProto) (messages []messageFromProto) {
	for _, content := range contents {
		lines := strings.Split(content.desc, "\n")
		for _, l := range lines[1:] {
			l = strings.TrimSpace(l)
			end := strings.Index(l, " ")
			if end != -1 {
				head := l[0:end]
				if head == "repeated" {
					newLine := l[9:]
					end = strings.Index(newLine, " ")
					head = newLine[0:end]
				}
				if head == "map<string," {
					newLine := l[12:]
					end = strings.Index(newLine, " ")
					head = newLine[0 : end-1]
				}
				custom := true
				for _, kt := range knownTypes {
					if head == kt {
						custom = false
						break
					}
				}
				if custom {
					messages = append(messages, messageFromProto{desc: head, fromProto: content.fromProto})
				}
			}
		}
	}
	return
}

// Remove package from message name
func removePackage(s string) string {
	kv := strings.Split(s, ".")
	return kv[len(kv)-1]
}

// Flag message when request/response is using an Empty message
func detectEmpty(s string, useEmpty *bool) string {
	if s == "Empty" {
		*useEmpty = true
	}
	return s
}
