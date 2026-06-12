export default function StatusBadge({ status }) {
  const styles = {
    Menunggu: 'bg-amber-50 text-amber-700 border border-amber-200',
    Diproses: 'bg-indigo-50 text-indigo-700 border border-indigo-200',
    Selesai: 'bg-emerald-50 text-emerald-700 border border-emerald-200',
  };

  return (
    <span className={`px-3 py-1 rounded-md text-xs font-bold uppercase tracking-wider ${styles[status] || 'bg-slate-100 text-slate-500'}`}>
      {status}
    </span>
  );
}