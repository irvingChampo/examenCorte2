export default function TextArea({ value, onChange }) {
  return <textarea rows={10} cols={50} value={value} onChange={onChange} />;
}
