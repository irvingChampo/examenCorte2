import FileInput from "../Atoms/FileInput";
import "../../Styles/Components/FileUploader.css";

export default function FileUploader({ onFileUpload }) {
  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    if (file.type !== "text/plain" && !file.name.endsWith(".txt")) {
      alert("❗ Solo se permiten archivos de texto (.txt).");
      return;
    }

    const reader = new FileReader();
    reader.onload = (event) => {
      const fileContent = event.target.result;
      onFileUpload(fileContent);
    };
    reader.onerror = () => {
      alert("❌ Error al leer el archivo.");
    };

    reader.readAsText(file);
  };

  return (
    <div className="file-uploader">
      <FileInput onChange={handleFileChange} accept=".txt,text/plain" />
    </div>
  );
}
