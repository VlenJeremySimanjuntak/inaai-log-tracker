import StatusBadge from './StatusBadge';

export default function TicketModal({ ticket, onClose, onUpdateStatus }) {
  if (!ticket) return null;

  // Helper agar Kategori terlihat rapi
  const getCategoryName = (id) => {
    if (id === 1) return '🌐 Jaringan & Internet';
    if (id === 2) return '🖥️ Server & Database';
    return '📱 Aplikasi Internal';
  };

  return (
    <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4 sm:p-0 animate-in fade-in duration-200">
      
      {/* Container Utama Modal */}
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg overflow-hidden animate-in zoom-in-95 duration-200 border border-slate-200 flex flex-col max-h-[90vh]">

        {/* --- HEADER GELAP --- */}
        <div className="bg-slate-800 px-6 py-4 flex justify-between items-center text-white shrink-0">
          <div className="flex items-center gap-3">
            <div className="bg-slate-700/50 p-2 rounded-lg">
              <svg className="w-5 h-5 text-indigo-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
            </div>
            <div>
              <h3 className="font-bold text-base sm:text-lg leading-tight">Detail Laporan</h3>
              <p className="text-slate-400 text-xs font-mono mt-0.5">ID: #{ticket.id}</p>
            </div>
          </div>
          <button 
            onClick={onClose} 
            className="text-slate-400 hover:text-white bg-slate-700/50 hover:bg-slate-700 w-8 h-8 flex items-center justify-center rounded-full transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>

        {/* --- KONTEN TENGAH --- */}
        <div className="p-6 overflow-y-auto flex-1 space-y-6 bg-slate-50/50">

          {/* KOTAK INFORMASI TERPADU (Judul, Kategori, Status) */}
          <div className="bg-white border border-slate-200 rounded-xl shadow-sm overflow-hidden">
            {/* Bagian Atas: Judul */}
            <div className="p-5">
              <p className="text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1.5 flex items-center gap-1.5">
                <span className="w-1.5 h-1.5 rounded-full bg-indigo-500"></span> Judul Gangguan
              </p>
              <h4 className="text-xl font-bold text-slate-800 leading-snug">{ticket.title}</h4>
            </div>
            
            {/* Bagian Bawah: Grid Kategori & Status */}
            <div className="grid grid-cols-2 divide-x divide-slate-100 border-t border-slate-100 bg-slate-50/80">
              <div className="p-4">
                <p className="text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-2">Kategori</p>
                <p className="text-sm text-slate-700 font-bold">{getCategoryName(ticket.category_id)}</p>
              </div>
              <div className="p-4 flex flex-col justify-center">
                <p className="text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-2">Status Saat Ini</p>
                <div className="self-start"><StatusBadge status={ticket.status} /></div>
              </div>
            </div>
          </div>

          {/* KOTAK DESKRIPSI */}
          <div>
            <p className="text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2 flex items-center gap-1.5 px-1">
              <span className="w-1.5 h-1.5 rounded-full bg-slate-400"></span> Deskripsi Lapangan
            </p>
            <div className="bg-white border border-slate-200 rounded-xl p-5 text-sm text-slate-700 leading-relaxed min-h-[120px] whitespace-pre-wrap shadow-sm">
              {ticket.description}
            </div>
          </div>

        </div>

        {/* --- FOOTER / AKSI --- */}
        <div className="bg-white border-t border-slate-200 px-6 py-4 flex flex-col-reverse sm:flex-row gap-3 shrink-0">
          
          <button 
            onClick={onClose} 
            className="w-full sm:w-auto px-6 py-2.5 bg-white border border-slate-300 hover:bg-slate-50 text-slate-700 rounded-xl text-sm font-bold transition-all focus:ring-2 focus:ring-slate-200"
          >
            Tutup
          </button>

          <div className="flex-1 flex gap-3">
            {ticket.status === 'Menunggu' && (
              <button
                onClick={() => { onUpdateStatus(ticket.id, 'Diproses'); onClose(); }}
                className="w-full bg-indigo-600 hover:bg-indigo-700 text-white py-2.5 rounded-xl text-sm font-bold transition-all shadow-md shadow-indigo-200 flex justify-center items-center gap-2"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                Mulai Proses Laporan
              </button>
            )}
            
            {ticket.status === 'Diproses' && (
              <button
                onClick={() => { onUpdateStatus(ticket.id, 'Selesai'); onClose(); }}
                className="w-full bg-emerald-600 hover:bg-emerald-700 text-white py-2.5 rounded-xl text-sm font-bold transition-all shadow-md shadow-emerald-200 flex justify-center items-center gap-2"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7"></path></svg>
                Tandai Selesai
              </button>
            )}
            
            {ticket.status === 'Selesai' && (
              <div className="w-full bg-slate-100 border border-slate-200 text-slate-400 py-2.5 rounded-xl text-sm font-bold flex justify-center items-center gap-2 cursor-not-allowed">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                Laporan Telah Tuntas
              </div>
            )}
          </div>

        </div>

      </div>
    </div>
  );
}