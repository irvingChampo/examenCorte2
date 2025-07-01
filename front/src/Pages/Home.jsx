import { useState } from "react";
import FileUploader from "../Components/Molecules/FileUploader";
import CodeInput from "../Components/Molecules/CodeInput";
import TokenAnalysisTable from "../Components/Molecules/CodeAnalysisTable";
import "../Styles/App.css";

export default function Home() {
  const [code, setCode] = useState(`def main():
    edad = 22
    escuela = "upchiapas"
    if edad > 18:
        print("Mayor de edad")
    if escuela.lower() == "upchiapas":
        print("Bienvenido a UPChiapas")

if __name__ == "__main__":
    main()`);
  const [result, setResult] = useState("");
  const [isAnalyzing, setIsAnalyzing] = useState(false);

  const handleAnalyze = async () => {
    if (!code.trim()) {
      alert("❗ No hay código para analizar.");
      return;
    }

    setIsAnalyzing(true);
    try {
      const response = await fetch("http://localhost:8080/analyze", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error("❌ Error del servidor:", errorText);
        alert("Error al analizar código:\n" + errorText);
        return;
      }

      const data = await response.text();
      setResult(data);
    } catch (error) {
      console.error("❌ Error al analizar:", error);
      alert("Error de conexión con el servidor.");
    } finally {
      setIsAnalyzing(false);
    }
  };

  return (
    <div className="container">
      <div className="left-panel">
        <h2 className="panel-title">Analizador de Código Python</h2>
        <FileUploader onFileUpload={setCode} />
        <CodeInput 
          code={code} 
          setCode={setCode} 
          onAnalyze={handleAnalyze} 
          isAnalyzing={isAnalyzing}
        />
      </div>
      
      <div className="right-panel">
        <h2 className="panel-title">Resultados del Análisis</h2>
        <div className="results-container">
          <TokenAnalysisTable analysisResult={result} />
        </div>
      </div>
    </div>
  );
}