package util

const (
	nc = "\033[0m"
	red = "\033[1;31m"
	green = "\033[1;32m"
	blue = "\033[1;34m"
)	

func Green(text string) string {
	return green + text + nc
}

func Red(text string) string {
	return red + text + nc
}

func Blue(text string) string {
	return blue + text + nc
}