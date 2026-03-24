package ui

import "fmt"

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
)

// Symbols
const (
	SymInfo    = "ℹ"
	SymSuccess = "✔"
	SymError   = "✘"
	SymWait    = "➜"
	SymPackage = "📦"
)

func LogInfo(format string, a ...interface{}) {
	fmt.Printf(ColorBlue+SymWait+" "+format+ColorReset+"\n", a...)
}

func LogSuccess(format string, a ...interface{}) {
	fmt.Printf(ColorGreen+SymSuccess+" "+format+ColorReset+"\n", a...)
}

func LogError(format string, a ...interface{}) {
	fmt.Printf(ColorRed+SymError+" "+format+ColorReset+"\n", a...)
}

func LogWarning(format string, a ...interface{}) {
	fmt.Printf(ColorYellow+SymInfo+" "+format+ColorReset+"\n", a...)
}
