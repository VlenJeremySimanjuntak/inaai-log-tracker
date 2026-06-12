export default function SummaryCard({ summary }) {
  if (!summary) return <div className="p-6 text-slate-400">Memuat analisis AI...</div>;

  // Fungsi sederhana untuk membersihkan text tanpa library
  const formatText = (text) => {
    if (!text) return '';
    return text
      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') // Ubah **tebal** jadi bold HTML
      .split('\n') // Pisahkan berdasarkan baris baru
      .map((paragraph, index) => {
        if (paragraph.trim().startsWith('*')) {
          // Jika diawali tanda bintang, buat jadi baris list
          return <li key={index} className="ml-4 list-disc mb-1" dangerouslySetInnerHTML={{ __html: paragraph.replace('*', '').trim() }} />;
        }
        return paragraph.trim() ? <p key={index} className="mb-4" dangerouslySetInnerHTML={{ __html: paragraph }} /> : null;
      });
  };

  return (
    <div className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm flex flex-col h-full">
      <div className="flex items-center gap-3 border-b border-slate-100 pb-4 mb-4">
        <div className="bg-indigo-50 p-2 rounded-xl text-indigo-600">✨</div>
        <div>
          <h3 className="text-base font-bold text-slate-800">Analisis AI (Gemini 2.5)</h3>
          <p className="text-xs text-slate-400">Ringkasan otomatis seluruh gangguan aktif</p>
        </div>
      </div>

      {/* Konten Hasil Format */}
      <div className="flex-1 text-sm text-slate-600 leading-relaxed overflow-y-auto max-h-[350px]">
        {formatText(summary.summary_text)}
      </div>

      <div className="mt-4 pt-3 border-t border-slate-100 flex justify-between items-center text-[11px] text-slate-400">
        <span>Log Teranalisis: <span className="font-mono bg-slate-100 px-1.5 py-0.5 rounded text-slate-600">{summary.log_ids_analyzed || '-'}</span></span>
      </div>
    </div>
  );
}