package analyzer

import (
	"fmt"
	"regexp"
	"strings"
)

type AnalysisError struct {
	Line    int
	Message string
}

type Token struct {
	Type  string
	Value string
}

func isNumber(str string) bool {
	ok, _ := regexp.MatchString(`^\d+$`, str)
	return ok
}

func AnalyzeCode(code string) string {
	lines := strings.Split(code, "\n")
	var tokens []Token
	var errors []AnalysisError
	var funcDefined bool
	variables := make(map[string]string)

	// ---------- ANÁLISIS LÉXICO ----------
	reserved := map[string]bool{
		"def": true, "if": true, "else": true, "return": true, "print": true,
	}

	operators := map[string]bool{
		"=": true, "*": true, "-": true, "<=": true, "(": true, ")": true, ":": true, ",": true,
	}

	reString := regexp.MustCompile(`"(.*?)"`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extraer todas las cadenas y agregarlas como tokens
		matches := reString.FindAllString(line, -1)
		for _, str := range matches {
			tokens = append(tokens, Token{"STRING", str})
			line = strings.Replace(line, str, "", 1)
		}

		// Tokenizar por espacio
		words := strings.Fields(line)
		for _, word := range words {
			if reserved[word] {
				tokens = append(tokens, Token{"KEYWORD", word})
			} else if operators[word] {
				tokens = append(tokens, Token{"OPERATOR", word})
			} else if isNumber(word) {
				tokens = append(tokens, Token{"NUMBER", word})
			} else if regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(word) {
				tokens = append(tokens, Token{"IDENTIFIER", word})
			} else if word != "" {
				tokens = append(tokens, Token{"UNKNOWN", word})
			}
		}
	}

	// ---------- ANÁLISIS SINTÁCTICO ----------
	for i, line := range lines {
		if strings.HasPrefix(line, "def factorial(") {
			funcDefined = true
			if !strings.Contains(line, "):") {
				errors = append(errors, AnalysisError{i + 1, "Definición de función inválida"})
			}
		}
		if strings.Contains(line, "if") && !strings.Contains(line, ":") {
			errors = append(errors, AnalysisError{i + 1, "Condición 'if' sin ':' final"})
		}
		if strings.Contains(line, "else") && !strings.Contains(line, ":") {
			errors = append(errors, AnalysisError{i + 1, "'else' sin ':' final"})
		}
		if strings.Contains(line, "return") && !strings.Contains(line, "factorial(") && !strings.Contains(line, "*") && !strings.Contains(line, "1") {
			errors = append(errors, AnalysisError{i + 1, "Return inválido"})
		}
	}

	if !funcDefined {
		errors = append(errors, AnalysisError{1, "Falta definición de la función factorial"})
	}

	// ---------- ANÁLISIS SEMÁNTICO ----------
	for idx, line := range lines {
		if strings.Contains(line, "x =") {
			variables["x"] = "int"
		}
		if strings.Contains(line, "factorial(x)") && !funcDefined {
			errors = append(errors, AnalysisError{idx + 1, "Uso de factorial() antes de ser definida"})
		}
		if strings.Contains(line, "print") {
			if !strings.Contains(line, "(") || !strings.Contains(line, ")") {
				errors = append(errors, AnalysisError{idx + 1, "Uso de print() sin paréntesis válidos"})
			}
		}
	}

	// ---------- RESULTADOS ----------
	var sb strings.Builder
	sb.WriteString("🧪 Tokens:\n")
	for _, tok := range tokens {
		sb.WriteString(fmt.Sprintf("  [%s] → %s\n", tok.Type, tok.Value))
	}
	sb.WriteString("\n")

	if len(errors) > 0 {
		sb.WriteString("❌ Errores encontrados:\n")
		for _, e := range errors {
			sb.WriteString(fmt.Sprintf("  Línea %d: %s\n", e.Line, e.Message))
		}
	} else {
		sb.WriteString("✅ No se encontraron errores sintácticos ni semánticos.\n")
	}

	return sb.String()
}
