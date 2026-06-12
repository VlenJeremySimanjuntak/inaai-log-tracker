export default function SummaryCard({ summary }) {
  if (!summary) return (
    <div className="bg-white rounded-2xl border border-slate-200 p-6 shadow-sm h-full flex flex-col justify-center animate-pulse">
      <div className="h-6 bg-slate-200 rounded w-1/3 mb-4"></div>
      <div className="h-4 bg-slate-200 rounded w-full mb-2"></div>
      <div className="h-4 bg-slate-200 rounded w-5/6"></div>
    </div>
  );

  return (
    <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 lg:p-8 h-full flex flex-col relative overflow-hidden">
      {/* Aksen background dekoratif halus di pojok */}
      <div className="absolute top-0 right-0 w-32 h-32 bg-indigo-50 rounded-bl-full -z-10 opacity-60"></div>
      
      <div className="flex items-center gap-3 mb-4 border-b border-slate-100 pb-4">
        <div className="bg-indigo-100 p-2.5 rounded-xl text-indigo-600">
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
        </div>
        <div>
          <h3 className="font-bold text-lg text-slate-800 leading-tight">Analisis AI (Gemini)</h3>
          <p className="text-xs text-slate-500 mt-0.5">Ringkasan otomatis dari seluruh laporan aktif</p>
        </div>
      </div>
      
      <div className="flex-grow flex items-start mt-2">
        <p className="text-slate-700 leading-relaxed text-sm sm:text-base font-medium">
          {summary.summary_text}
        </p>
      </div>
      
      <div className="mt-6 pt-4 border-t border-slate-100 flex flex-col sm:flex-row justify-between items-start sm:items-center text-xs text-slate-400 gap-2">
        <span className="bg-slate-50 px-2 py-1 rounded-md border border-slate-100">ID Teranalisis: {summary.log_ids_analyzed || '-'}</span>
        <span>Pembaruan Terakhir: {new Date(summary.created_at).toLocaleString('id-ID')}</span>
      </div>
    </div>
  );
}