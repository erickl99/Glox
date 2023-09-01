package main

import (
	"bufio"
	"fmt"
	"os"
)

var error_occurred = false

func main() {
    if len(os.Args) > 2 {
        fmt.Println("Usage: glox [script]")
    } else if len(os.Args) == 2 {
        run_file(os.Args[1])
    } else {
        run_prompt()
    }
}

func run_file(name string) {
    bytes, err := os.ReadFile(name)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    source := string(bytes[:])
    run(source)
    // if err != nil {
    //     fmt.Fprintln(os.Stderr, err)
    //     os.Exit(65)
    // }
    if error_occurred {
        os.Exit(65)
    }
}

func run_prompt() {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("> ")
    for scanner.Scan() {
        source := scanner.Text()
        if source == "quit" {
            break
        }
        run(source)
        error_occurred = false
        fmt.Print("> ")
    }
    if err := scanner.Err(); err != nil {
        fmt.Println(err)
    }
    fmt.Println("Bye")
}

func run(source string) {
    lscanner := NewLexer(source)
    tokens := lscanner.scan_tokens()
    parser := Parser{tokens: tokens}
    expr, _ := parser.parse()
    if error_occurred {
        fmt.Println("Uh oh")
        return
    }
    fmt.Println(print(expr))
}

func line_error(line int, message string) {
    report(line, "", message)
}

func token_error(token Token, message string) {
    if token.t_type == EOF {
        report(token.line, "at end", message)
    } else {
        report(token.line, "at '" + token.lexeme + "'", message)
    }
}

func report(line int, where, message string) {
    fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s", line, where, message)
    error_occurred = true
}
