export default function FileInput({ onChange, accept }) {
  return (
    <input
      type="file"
      onChange={onChange}
      accept={accept}
      className="custom-file-input"
    />
  );
}
