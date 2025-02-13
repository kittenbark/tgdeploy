package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

const (
	dstDefault        = "README.md"
	srcDefault        = "template/README.md"
	composeDefault    = "compose.yml"
	dockerfileDefault = "Dockerfile"
)

func main() {
	dst := at(os.Args, 1, dstDefault)
	if slices.Contains([]string{"-h", "-help", "--help", "help"}, strings.ToLower(dst)) {
		fmt.Printf("Usage: %s [<destination> | %s] [<source> | %s] [<compose> | %s] [<dockerfile> | %s]\n",
			os.Args[0], dstDefault, srcDefault, composeDefault, dockerfileDefault)
		os.Exit(0)
	}
	src := at(os.Args, 2, srcDefault)

	srcFile, err := os.ReadFile(src)
	if err != nil {
		log.Fatalf("failed to read source file: %v", err)
	}

	compose := at(os.Args, 3, composeDefault)
	composeFile, err := os.ReadFile(compose)
	if err != nil {
		log.Fatalf("failed to read compose file: %v", err)
	}

	dockerfile := at(os.Args, 4, dockerfileDefault)
	dockerfileFile, err := os.ReadFile(dockerfile)
	if err != nil {
		log.Fatalf("failed to read dockerfile file: %v", err)
	}

	if err := os.WriteFile(dst, build(string(srcFile), string(composeFile), string(dockerfileFile)), os.ModePerm); err != nil {
		log.Fatalf("failed to write to destination file: %v", err)
	}
}

func build(source string, compose string, dockerfile string) []byte {
	replacer := strings.NewReplacer(
		"{{compose.yml}}", strings.TrimSpace(compose),
		"{{Dockerfile}}", strings.TrimSpace(dockerfile),
	)

	return []byte(replacer.Replace(source))
}

func at[T any](list []T, index int, or T) T {
	if index < 0 || index >= len(list) {
		return or
	}
	return list[index]
}
