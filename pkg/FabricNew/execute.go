package FabricNew

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"get.porter.sh/porter/pkg/exec/builder"
	"gopkg.in/yaml.v2"
)

type Dashes struct {
	Long  string
	Short string
}

var DefaultFlagDashes = Dashes{
	Long:  "--",
	Short: "-",
}

type HasCustomDashes interface {
	GetDashes() Dashes
}

func (m *Mixin) loadAction(ctx context.Context) (*Action, error) {
	var action Action
	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &action)
		return &action, err
	})
	return &action, err
}

func (m *Mixin) Execute(ctx context.Context) error {
	fmt.Println("Entering Execute method..")
	action, err := m.loadAction(ctx)
	if err != nil {
		return err
	}

	//fmt.Println(action)
	//step := action.Steps[0]

	/**Get Steps from action.
	1. Currently we are suppporting single step.
	2. Each Fabric artifact is to be mentioned as separate step.
		Eg. Workspace as step 1, Warehouse as Step 2
	**/
	//steps := action.GetSteps()

	/*
		1. Executes command.
		2. In Fabric Case it would FabricClient (.NET Cli program that has Fabric Deployment Library loaded in process.)
		**/
	//command := steps[0].GetCommand()

	/**
		Fetching Arguments from command.
		In Fabric there are two sections
		1. service: core.workspace / core.items (similar to namespaces)
		2. Command: create (actual commands)
	**/
	// var cmd *exec.Cmd
	//Method : 1
	// arguments := steps[0].GetArguments()
	// flags := steps[0].GetFlags()

	// strCmd := command + " --action " + arguments[1] + " "
	// for _, chunk := range flags {
	// 	flagCmd := "--" + chunk.Name + " " + chunk.Values[0] + " "
	// 	strCmd = strCmd + flagCmd
	// }
	// fmt.Println(strCmd)
	// cmd = exec.Command(strCmd)
	// output, err := cmd.CombinedOutput()
	// fmt.Println(output)
	//iterate through flags and append them to command
	fmt.Println(action.Steps[0].GetArguments())
	_, err = builder.ExecuteSingleStepAction(ctx, m.RuntimeConfig, action)
	//fmt.Println(result)
	//return m.handleOutputs(ctx, step.Outputs)

	//Method 2:
	// args := make([]string, len(arguments), 1+len(arguments)+len(flags)*2)
	// // group := []string{"group", arguments[0]}
	// copy(args, arguments)
	// // operation := []string{"--action", arguments[1]}
	// // copy(args, operation)
	// flagsSlice := splitCommand(flags.ToSlice(builder.Dashes(DefaultFlagDashes)))
	// args = append(args, flagsSlice...)
	// cmd := m.Context.NewCommand(ctx, steps[0].GetCommand(), args...)
	// //cmd := exec.Command(command, args...)
	// fmt.Println(cmd)
	// output, err := cmd.CombinedOutput()

	// m.Out.Write([]byte(action.Steps[0].Outputs[0].JsonPath))
	// m.Out.Write([]byte(action.Steps[0].Outputs[0].Name))
	// fmt.Println((action.Steps[0].Outputs[0]))
	// fmt.Println(string(output))
	fmt.Println(err)
	return err
}

var whitespace = string([]rune{space, newline, tab})

const (
	space       = rune(' ')
	newline     = rune('\n')
	tab         = rune('\t')
	backslash   = rune('\\')
	doubleQuote = rune('"')
	singleQuote = rune('\'')
)

func splitCommand(slice []string) []string {
	expandedSlice := make([]string, 0, len(slice))
	for _, chunk := range slice {
		chunkettes := findWords(chunk)
		expandedSlice = append(expandedSlice, chunkettes...)
	}

	return expandedSlice
}
func findWords(input string) []string {
	words := make([]string, 0, 1)
	next := input
	for len(next) > 0 {
		word, remainder, err := findNextWord(next)
		if err != nil {
			return []string{input}
		}
		next = remainder
		words = append(words, word)
	}

	return words
}

func findNextWord(input string) (string, string, error) {
	var buf bytes.Buffer

	// Remove leading whitespace before starting
	input = strings.TrimLeft(input, whitespace)

	var escaped bool
	var wordStart, wordStop int
	var closingQuote rune

	for i, r := range input {
		// Prevent escaped characters from matching below
		if escaped {
			r = -1
			escaped = false
		}

		switch r {
		case backslash:
			// Escape the next character
			escaped = true
			continue
		case closingQuote:
			wordStop = i
			closingQuote = 0 // Reset looking for a closing quote
		case singleQuote, doubleQuote:
			// Seek to the closing quote only
			if closingQuote != 0 {
				continue
			}

			wordStart = 1    // Skip opening quote
			closingQuote = r // Seek to the same closing quote
		case space, tab, newline:
			// Seek to the closing quote only
			if closingQuote != 0 {
				continue
			}

			wordStart = 0
			wordStop = i
		}

		// Found the end of a word
		if wordStop > 0 {
			_, err := buf.WriteString(input[wordStart:wordStop])
			if err != nil {
				return "", input, errors.New("error writing to buffer")
			}
			return buf.String(), input[wordStop+1:], nil
		}
	}

	if closingQuote != 0 {
		return "", "", errors.New("unmatched quote found")
	}

	// Hit the end of input, flush the remainder
	_, err := buf.WriteString(input)
	if err != nil {
		return "", input, errors.New("error writing to buffer")
	}

	return buf.String(), "", nil
}
