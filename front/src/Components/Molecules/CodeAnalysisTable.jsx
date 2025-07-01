import React from "react";
import "../../Styles/Components/CodeAnalysisTable.css";

const TokenAnalysisTable = ({ analysisResult }) => {
  if (!analysisResult) return null;

  const parseAnalysis = (result) => {
    const categories = {
      "Palabras reservadas": new Set(),
      "Identificadores": new Set(),
      "Operadores": new Set(),
      "Números": new Set(),
      "Símbolos": new Set(),
      "Cadenas": new Set(),
      "Comentarios": new Set(),
      "Errores": new Map(),
    };

    const syntaxErrors = [];
    const semanticErrors = [];
    const successMessages = [];
    const summary = {};
    const lines = result.split('\n');
    let currentCategory = '';
    let inSummary = false;

    lines.forEach(line => {
      const trimmedLine = line.trim();

      if (!trimmedLine) return;

      // Detectar resumen
      if (trimmedLine === "===RESUMEN===") {
        inSummary = true;
        return;
      }

      if (trimmedLine.startsWith("✅")) {
        successMessages.push(trimmedLine);
        return;
      }


      // Resumen
      if (inSummary) {
        const summaryMatch = trimmedLine.match(/^([A-Za-záéíóúñÑ\s]+):\s*(\d+)/);
        if (summaryMatch) {
          const key = summaryMatch[1].trim();
          const count = parseInt(summaryMatch[2]);
          summary[key] = count;
        }
        return;
      }

      // Detectar secciones de errores
      if (trimmedLine.startsWith("🧱 Errores de Sintaxis:")) {
        currentCategory = "ErroresSintacticos";
        return;
      }
      if (trimmedLine.startsWith("🧠 Errores Semánticos:")) {
        currentCategory = "ErroresSemanticos";
        return;
      }
      if (trimmedLine.startsWith("⚠️ Errores Léxicos:")) {
        currentCategory = "Errores";
        return;
      }

      if (currentCategory === "ErroresSintacticos" && trimmedLine.startsWith("- Línea")) {
        syntaxErrors.push(trimmedLine);
        return;
      }

      if (currentCategory === "ErroresSemanticos" && trimmedLine.startsWith("- Línea")) {
        semanticErrors.push(trimmedLine);
        return;
      }

      if (currentCategory === "Errores" && trimmedLine.includes('→ Error:')) {
        const [tokenPart, suggestionPart] = trimmedLine.split('→');
        const token = tokenPart.replace('❌', '').trim().replace(/^'|'$/g, '');
        const suggestion = suggestionPart.replace('Error:', '').trim();
        categories["Errores"].set(token, suggestion);
        return;
      }

      // Categoría de tokens
      const categoryMatch = trimmedLine.match(/^([A-Za-záéíóúñÑ\s]+) \(\d+\):/);
      if (categoryMatch) {
        currentCategory = categoryMatch[1].trim();
        return;
      }

      // Agregar token a categoría
      if (currentCategory && categories[currentCategory]) {
        const token = trimmedLine.replace(/^[;❌]/, '').trim();
        categories[currentCategory].add(token);
      }
    });

    return {
      ...categories,
      Errores: Array.from(categories["Errores"].entries()),
      syntaxErrors,
      semanticErrors,
      summary,
      successMessages,
    };
  };

  const analysisData = parseAnalysis(analysisResult);

  return (
    <div className="token-analysis-container">
      {/* Errores Léxicos */}
      {analysisData.Errores.length > 0 && (
        <div className="errors-section">
          <h3 className="errors-title">⚠️ Errores Léxicos</h3>
          <ul className="errors-list">
            {analysisData.Errores.map(([token, suggestion], index) => (
              <li key={index} className="error-item">
                <span className="error-token">'{token}'</span>
                <span className="error-suggestion">→ {suggestion}</span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Errores Sintácticos */}
      {analysisData.syntaxErrors.length > 0 && (
        <div className="errors-section">
          <h3 className="errors-title">🧱 Errores de Sintaxis</h3>
          <ul className="errors-list">
            {analysisData.syntaxErrors.map((err, i) => (
              <li key={i} className="error-item">{err}</li>
            ))}
          </ul>
        </div>
      )}

      {/* Errores Semánticos */}
      {analysisData.semanticErrors.length > 0 && (
        <div className="errors-section">
          <h3 className="errors-title">🧠 Errores Semánticos</h3>
          <ul className="errors-list">
            {analysisData.semanticErrors.map((err, i) => (
              <li key={i} className="error-item">{err}</li>
            ))}
          </ul>
        </div>
      )}

      {/* Mensajes de éxito (sin errores) */}
      {analysisData.successMessages.length > 0 && (
      <div className="errors-section" style={{ borderLeft: "4px solid #28a745" }}>
        <h3 className="errors-title" style={{ color: "#28a745" }}>✅ Análisis Exitoso</h3>
          <ul className="errors-list">
          {analysisData.successMessages.map((msg, i) => (
            <li key={i} className="error-item">{msg}</li>
          ))}
          </ul>
      </div>
)}


      {/* Tokens por categoría */}
      <div className="tokens-table-section">
        <h3>Tokens Encontrados</h3>
        <div className="categories-grid">
          {Object.entries(analysisData)
            .filter(([name, items]) =>
              name !== "Errores" &&
              name !== "summary" &&
              name !== "syntaxErrors" &&
              name !== "semanticErrors" &&
              items.size > 0
            )
            .map(([name, items]) => (
              <div key={name} className="category-block">
                <h4 className="category-title">
                  {name} <span className="count-badge">{items.size}</span>
                </h4>
                <ul className="token-list">
                  {Array.from(items).map((token, i) => (
                    <li key={i} className="token-item">
                      {name === "Comentarios" ? `;${token}` : token}
                    </li>
                  ))}
                </ul>
              </div>
            ))}
        </div>
      </div>

      {/* Resumen de Tokens */}
      {Object.keys(analysisData.summary).length > 0 && (
        <div className="summary-section">
          <h3>📊 Resumen de Tokens</h3>
          <div className="summary-grid">
            {Object.entries(analysisData.summary).map(([category, count], i) => (
              <div key={i} className="summary-item">
                <span className="summary-category">{category}</span>
                <span className="summary-count">{count}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default TokenAnalysisTable;
