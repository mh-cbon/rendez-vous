# {{.Name}}

{{pkgdoc}}

# cli

#### $ {{shell "go run main.go -h" | color "sh"}}

# tests

#### $ {{shell "go test -v" | color "sh"}}

# todos

#### $ {{shell "grep --include='*go' -r todo -B 1 -A 1 -n" | color "sh"}}
