package debug

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Names list of which name has which color
var names map[string]int = make(map[string]int)

// List of available colors
var colors = []int{6, 2, 3, 4, 5, 1}

// The default color
var defaultColor = 6

// Previous milliseconds
var prevTime int64 = 0

// The type that is returned in Debug function.
type printType func(format string, a ...interface{})

// Get string with terminal color.
func getColorString(firstColor int, lastColor int, message string) string {
	if !useColors() {
		return message
	}

	return fmt.Sprintf("\u001b[%d%dm%s\u001b[0m", firstColor, lastColor, message)
}

// Print debug message with namespace and message.
func printDebug(namespace string, format string, a ...interface{}) {
	color := names[namespace]
	namespace = getColorString(9, color, namespace)
	message := getColorString(9, 0, fmt.Sprintf(format, a...))

	if timeFormat, ok := useMS(); !ok {
		now := time.Now()

		if timeFormat == "utc" {
			now = time.Now().UTC()
		}

		fmt.Println(fmt.Sprintf("%s %s %s", now, namespace, message))
	} else {
		ms := getColorString(3, color, fmt.Sprintf("+%dms", getMs()))
		fmt.Println(fmt.Sprintf("%s %s %s", namespace, message, ms))
	}

	// No support for logging debug to file, feel free to contribute!
}

// Check if string is in slice.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Check if namespace is allowed to print the debug message or not.
func checkNamespaceStatus(namespace string) bool {
	value := os.Getenv("DEBUG")

	if len(value) == 0 {
		return false
	}

	values := strings.Split(value, ",")

	if len(values) == 1 && values[0] == "*" {
		return true
	}

	if stringInSlice(namespace, values) {
		return true
	}

	star := false

	for _, name := range values {
		ns := strings.Split(namespace, ":")

		if name == "*" && !star {
			star = true
		}

		if strings.Contains(name, ns[0]) {
			if string(name[0]) == "-" && (string(name[1:]) == namespace || string(name[len(name)-1]) == "*") {
				return false
			} else if string(name[len(name)-1]) == "*" {
				parent := name[0 : len(name)-2]
				return strings.Split(parent, ":")[0] == parent
			} else {
				return false
			}
		}
	}

	return star
}

// Get milliseconds between debug messages.
func getMs() int64 {
	curr := time.Now().UnixNano() % 1e6 / 1e3
	ms := curr

	if prevTime == 0 {
		ms = 0
	} else {
		ms -= prevTime
	}

	prevTime = curr

	return ms
}

// Check if we should use colors or not.
func useColors() bool {
	value := os.Getenv("DEBUG_COLORS")
	return value != "0" && value != "no" && value != "false" && value != "disabled"
}

// Check if we should print milliseconds or time.
func useMS() (string, bool) {
	value := os.Getenv("DEBUG_TIME")

	if len(value) == 0 || value == "ms" {
		return "ms", true
	}

	return value, false
}

// Create a new namespace to debug from.
func Debug(namespace string) printType {
	enabled := checkNamespaceStatus(namespace)

	if !enabled {
		return func(format string, a ...interface{}) {}
	}

	if _, ok := names[namespace]; !ok {
		color := defaultColor

		// We are out of colors, default to green.
		if len(colors) != 0 {
			color = colors[0]
			colors = colors[1:]
		}

		names[namespace] = color
	}

	return func(format string, a ...interface{}) {
		printDebug(namespace, format, a...)
	}
}
