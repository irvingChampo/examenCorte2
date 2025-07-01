import "../../Styles/Components/CodeInput.css"

export default function CodeInput({ code, setCode, onAnalyze }) {
  return (
    <div className="code-input">
      <textarea
        value={code}
        onChange={(e) => setCode(e.target.value)}
        placeholder="Escribe o sube tu código aquí..."
      />
      <button onClick={onAnalyze}>Analizar</button>
    </div>
  );
}
