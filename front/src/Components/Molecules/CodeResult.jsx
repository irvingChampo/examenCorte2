import "../../Styles/Components/CodeResult.css";

export default function CodeResult({ originalCode, analysisResult }) {
  return (
    <div className="code-result">
      <div className="code-block">
        <h3>Código</h3>
        <pre>{originalCode}</pre>
      </div>
      <div className="result-block">
        <h3>Resultado</h3>
        <pre>{analysisResult}</pre>
      </div>
    </div>
  );
}
