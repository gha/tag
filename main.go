package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	tagPrefix    = "v"
	tagSeparator = "."
)

var validTagTypes = map[string]bool{
	"major": true,
	"minor": true,
	"patch": true,
}

func main() {
	cmdArgs := os.Args[1:]
	if len(cmdArgs) != 1 || !validTagTypes[cmdArgs[0]] {
		log.Fatal("Tag type required (major, minor or patch)")
	}

	tagType := cmdArgs[0]

	var out bytes.Buffer
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	currentTag := strings.TrimSpace(out.String())
	tagParts := strings.Split(strings.TrimLeft(currentTag, tagPrefix), tagSeparator)

	if len(tagParts) < 1 || len(tagParts) > 3 {
		log.Fatal("Invalid tag structure in repository: ", currentTag)
	}

	majVersion := tagParts[0]
	minVersion := "0"
	patVersion := "0"

	if len(tagParts) > 1 {
		minVersion = tagParts[1]
	}
	if len(tagParts) > 2 {
		patVersion = tagParts[2]
	}

	iMajV, err := strconv.Atoi(majVersion)
	if err != nil {
		log.Fatal(err)
	}

	iMinV, err := strconv.Atoi(minVersion)
	if err != nil {
		log.Fatal(err)
	}

	iPatV, err := strconv.Atoi(patVersion)
	if err != nil {
		log.Fatal(err)
	}

	switch tagType {
	case "major":
		iMajV++
		iMinV = 0
		iPatV = 0
	case "minor":
		iMinV++
		iPatV = 0
	case "patch":
		iPatV++
	default:
		log.Fatal("Invalid tag type")
	}

	majVersion = strconv.Itoa(iMajV)
	minVersion = strconv.Itoa(iMinV)
	patVersion = strconv.Itoa(iPatV)
	newTag := tagPrefix + strings.Join([]string{majVersion, minVersion, patVersion}, tagSeparator)

	var tagCmdOut bytes.Buffer
	tagCmd := exec.Command("git", "tag", newTag)
	tagCmd.Stdout = &tagCmdOut
	tagCmd.Stderr = &tagCmdOut
	if err := tagCmd.Run(); err != nil {
		log.Fatal(err)
	}

	if tagCmdOut.String() != "" {
		log.Fatal("Unexpected tag output: ", tagCmdOut.String())
	}

	fmt.Println("Old Tag:", currentTag)
	fmt.Println("New Tag:", newTag)
}
