package analyzer

import (
	"fmt"
	"regexp"
	"strings"
)

type AnalysisError struct {
	Line    int
	Type    string
	Message string
}

type Token struct {
	Type  string
	Value string
	Line  int
}

func isNumber(str string) bool {
	ok, _ := regexp.MatchString(`^\d+$`, str)
	return ok
}

func isString(str string) bool {
	return (strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"")) ||
		(strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'"))
}

func isIdentifier(str string) bool {
	ok, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, str)
	return ok
}

func AnalyzeCode(code string) string {
	lines := strings.Split(code, "\n")
	var tokens []Token
	var errors []AnalysisError

	// Variables para análisis semántico
	variables := make(map[string]string)
	functions := make(map[string]bool)

	// Identificadores especiales válidos
	specialIdentifiers := map[string]bool{
		"__name__": true,
		"__main__": true,
	}

	// ---------- ANÁLISIS LÉXICO ----------
	reserved := map[string]bool{
		"def": true, "if": true, "else": true, "elif": true, "return": true,
		"print": true, "import": true, "from": true, "as": true, "class": true,
		"try": true, "except": true, "finally": true, "with": true, "for": true,
		"while": true, "break": true, "continue": true, "pass": true, "lambda": true,
		"and": true, "or": true, "not": true, "in": true, "is": true, "None": true,
		"True": true, "False": true, "__name__": true, "__main__": true,
	}

	operators := map[string]bool{
		"=": true, "==": true, "!=": true, "<": true, ">": true, "<=": true, ">=": true,
		"+": true, "-": true, "*": true, "/": true, "%": true, "**": true, "//": true,
		"(": true, ")": true, "[": true, "]": true, "{": true, "}": true,
		":": true, ",": true, ".": true, ";": true,
	}

	// Análisis línea por línea
	for lineNum, line := range lines {
		originalLine := line
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Detectar strings y agregarlos como tokens
		reString := regexp.MustCompile(`"([^"]*)"`).FindAllString(line, -1)
		for _, str := range reString {
			tokens = append(tokens, Token{"STRING", str, lineNum + 1})
			line = strings.Replace(line, str, " ", 1)
		}

		reString2 := regexp.MustCompile(`'([^']*)'`).FindAllString(originalLine, -1)
		for _, str := range reString2 {
			tokens = append(tokens, Token{"STRING", str, lineNum + 1})
			line = strings.Replace(line, str, " ", 1)
		}

		// Detectar operadores de dos caracteres primero
		for op := range operators {
			if len(op) == 2 && strings.Contains(line, op) {
				tokens = append(tokens, Token{"OPERATOR", op, lineNum + 1})
				line = strings.Replace(line, op, " ", -1)
			}
		}

		// Tokenizar por caracteres especiales y espacios
		words := strings.FieldsFunc(line, func(c rune) bool {
			char := string(c)
			return c == ' ' || c == '\t' || operators[char]
		})

		// Agregar operadores de un caracter encontrados
		for _, char := range line {
			charStr := string(char)
			if operators[charStr] && len(charStr) == 1 {
				tokens = append(tokens, Token{"OPERATOR", charStr, lineNum + 1})
			}
		}

		// Procesar palabras
		for _, word := range words {
			word = strings.TrimSpace(word)
			if word == "" {
				continue
			}

			if reserved[word] {
				tokens = append(tokens, Token{"KEYWORD", word, lineNum + 1})
			} else if isNumber(word) {
				tokens = append(tokens, Token{"NUMBER", word, lineNum + 1})
			} else if isIdentifier(word) {
				tokens = append(tokens, Token{"IDENTIFIER", word, lineNum + 1})
			} else {
				tokens = append(tokens, Token{"UNKNOWN", word, lineNum + 1})
				errors = append(errors, AnalysisError{lineNum + 1, "LEXICAL", fmt.Sprintf("Token desconocido: '%s'", word)})
			}
		}
	}

	// ---------- ANÁLISIS SINTÁCTICO ----------
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "def ") {
			if !strings.Contains(line, "(") || !strings.Contains(line, ")") || !strings.HasSuffix(line, ":") {
				errors = append(errors, AnalysisError{i + 1, "SYNTAX", "Definición de función inválida - debe tener paréntesis y ':' al final"})
			} else {
				funcMatch := regexp.MustCompile(`def\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`).FindStringSubmatch(line)
				if len(funcMatch) > 1 {
					functions[funcMatch[1]] = true
				}
			}
		}

		if strings.HasPrefix(line, "if ") {
			if !strings.HasSuffix(line, ":") {
				errors = append(errors, AnalysisError{i + 1, "SYNTAX", "Declaración 'if' debe terminar con ':'"})
			}
		}

		if strings.Contains(line, "print(") {
			if !strings.Contains(line, ")") {
				errors = append(errors, AnalysisError{i + 1, "SYNTAX", "Llamada a print() incompleta - falta ')'"})
			}
		}

		if strings.HasPrefix(line, "def ") || strings.HasPrefix(line, "if ") || strings.HasPrefix(line, "else:") {
			for j := i + 1; j < len(lines); j++ {
				nextLine := lines[j]
				if strings.TrimSpace(nextLine) == "" {
					continue
				}
				if !strings.HasPrefix(nextLine, "    ") && !strings.HasPrefix(nextLine, "\t") {
					errors = append(errors, AnalysisError{j + 1, "SYNTAX", "Falta indentación después de ':'"})
				}
				break
			}
		}

		// Validar asignaciones y registrar variables
		if strings.Contains(line, "=") && !strings.Contains(line, "==") && !strings.Contains(line, "!=") && !strings.Contains(line, "<=") && !strings.Contains(line, ">=") {
			parts := strings.Split(line, "=")
			if len(parts) >= 2 {
				varName := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				if isIdentifier(varName) {
					if isNumber(value) {
						variables[varName] = "int"
					} else if isString(value) {
						if varName == "edad" {
							errors = append(errors, AnalysisError{i + 1, "SEMANTIC", fmt.Sprintf("Variable '%s' debe ser numérica, no string", varName)})
						}
						variables[varName] = "string"
					} else if value == "True" || value == "False" {
						variables[varName] = "bool"
					} else if value == "None" {
						variables[varName] = "none"
					} else if isIdentifier(value) {
						if otherType, exists := variables[value]; exists {
							variables[varName] = otherType
						} else {
							variables[varName] = "unknown"
							errors = append(errors, AnalysisError{i + 1, "SEMANTIC", fmt.Sprintf("Variable '%s' asignada con variable '%s' no declarada", varName, value)})
						}
					} else {
						variables[varName] = "unknown"
					}

					// ✅ Aquí agregamos la validación para escuela
					if varName == "escuela" && variables[varName] != "string" {
						errors = append(errors, AnalysisError{i + 1, "SEMANTIC", fmt.Sprintf("Variable '%s' debe ser de tipo string", varName)})
					}
				}
			}
		}
	}

	// ---------- ANÁLISIS SEMÁNTICO adicional: verificar uso de variables declaradas ----------
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Verificar print
		printMatch := regexp.MustCompile(`print\s*\(\s*([^)]+)\s*\)`).FindStringSubmatch(line)
		if len(printMatch) > 1 {
			content := strings.TrimSpace(printMatch[1])
			if !isString(content) && isIdentifier(content) {
				if _, exists := variables[content]; !exists {
					if !specialIdentifiers[content] {
						errors = append(errors, AnalysisError{i + 1, "SEMANTIC", fmt.Sprintf("Variable '%s' usada sin declarar", content)})
					}
				}
			}
		}

		// Comparaciones
		allOps := []string{"==", "!=", ">", "<", ">=", "<="}
		for _, op := range allOps {
			if strings.Contains(line, op) {
				parts := strings.Split(line, op)
				if len(parts) == 2 {
					left := strings.TrimSpace(parts[0])
					left = strings.TrimPrefix(left, "if")
					left = strings.TrimSpace(left)
					if isIdentifier(left) {
						if _, exists := variables[left]; !exists && !specialIdentifiers[left] {
							errors = append(errors, AnalysisError{i + 1, "SEMANTIC", fmt.Sprintf("Variable '%s' usada en comparación sin declarar", left)})
						}
					}
				}
			}
		}

		// Métodos en string
		methodMatch := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\s*\.\s*lower\s*\(\s*\)`).FindStringSubmatch(line)
		if len(methodMatch) > 1 {
			varName := methodMatch[1]
			if _, exists := variables[varName]; !exists && !specialIdentifiers[varName] {
				errors = append(errors, AnalysisError{i + 1, "SEMANTIC", fmt.Sprintf("Variable '%s' usada sin declarar", varName)})
			}
		}
	}

	// ---------- GENERAR RESULTADOS ----------
	var sb strings.Builder

	// Contar tokens
	tokenCount := make(map[string]int)
	for _, token := range tokens {
		tokenCount[token.Type]++
	}

	categories := map[string]string{
		"KEYWORD":    "Palabras reservadas",
		"IDENTIFIER": "Identificadores",
		"OPERATOR":   "Operadores",
		"NUMBER":     "Números",
		"STRING":     "Cadenas",
	}

	for tokenType, categoryName := range categories {
		if tokenCount[tokenType] > 0 {
			sb.WriteString(fmt.Sprintf("%s (%d):\n", categoryName, tokenCount[tokenType]))
			for _, token := range tokens {
				if token.Type == tokenType {
					sb.WriteString(fmt.Sprintf("%s\n", token.Value))
				}
			}
			sb.WriteString("\n")
		}
	}

	// Mostrar errores
	lexicalErrors := []AnalysisError{}
	syntaxErrors := []AnalysisError{}
	semanticErrors := []AnalysisError{}

	for _, err := range errors {
		switch err.Type {
		case "LEXICAL":
			lexicalErrors = append(lexicalErrors, err)
		case "SYNTAX":
			syntaxErrors = append(syntaxErrors, err)
		case "SEMANTIC":
			semanticErrors = append(semanticErrors, err)
		}
	}

	if len(lexicalErrors) > 0 {
		sb.WriteString("⚠️ Errores Léxicos:\n")
		for _, err := range lexicalErrors {
			sb.WriteString(fmt.Sprintf("❌ '%s' → Error: %s\n", err.Message, err.Message))
		}
		sb.WriteString("\n")
	}

	if len(syntaxErrors) > 0 {
		sb.WriteString("🧱 Errores de Sintaxis:\n")
		for _, err := range syntaxErrors {
			sb.WriteString(fmt.Sprintf("- Línea %d: %s\n", err.Line, err.Message))
		}
		sb.WriteString("\n")
	}

	if len(semanticErrors) > 0 {
		sb.WriteString("🧠 Errores Semánticos:\n")
		for _, err := range semanticErrors {
			sb.WriteString(fmt.Sprintf("- Línea %d: %s\n", err.Line, err.Message))
		}
		sb.WriteString("\n")
	}

	if len(errors) == 0 {
		sb.WriteString("✅ Análisis completado sin errores.\n")
		sb.WriteString("✅ Código Python sintácticamente correcto.\n")
		sb.WriteString("✅ Análisis semántico exitoso.\n\n")
	}

	sb.WriteString("===RESUMEN===\n")
	for tokenType, categoryName := range categories {
		if tokenCount[tokenType] > 0 {
			sb.WriteString(fmt.Sprintf("%s: %d\n", categoryName, tokenCount[tokenType]))
		}
	}

	return sb.String()
}
