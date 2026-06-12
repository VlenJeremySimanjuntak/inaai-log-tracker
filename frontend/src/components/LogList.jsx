import { useState } from 'react';
import StatusBadge from './StatusBadge';

// PERHATIKAN: Di baris ini kita menggunakan "onActionClick", bukan "onUpdateStatus" lagi
export default function LogList({ logs, onActionClick, isAdmin }) {
  // --- LOGIKA PAGINATION ---
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 5; 

  const countMenunggu = logs.filter(l => l.status === 'Menunggu').length;
  const countDiproses = logs.filter(l => l.status === 'Diproses').length;
  const countSelesai = logs.filter(l => l.status === 'Selesai').length;

  const totalPages = Math.max(1, Math.ceil(logs.length / itemsPerPage));
  const startIndex = (currentPage - 1) * itemsPerPage;
  const currentLogs = logs.slice(startIndex, startIndex + itemsPerPage);

  const goToNext = () => setCurrentPage((p) => Math.min(p + 1, totalPages));
  const goToPrev = () => setCurrentPage((p) => Math.max(p - 1, 1));

  return (
    <div className="space-y-4">
      
      {/* --- REKAP STATUS --- */}
      <div className="flex gap-3">
        <div className="bg-white border border-slate-200 px-4 py-2 rounded-xl shadow-sm text-sm font-semibold flex items-center gap-2 text-slate-700">
          Total Laporan: <span className="bg-slate-100 px-2 py-0.5 rounded-md">{logs.length}</span>
        </div>
        <div className="bg-amber-50 border border-amber-200 px-4 py-2 rounded-xl text-sm font-semibold flex items-center gap-2 text-amber-700">
          Menunggu: <span className="bg-amber-100 px-2 py-0.5 rounded-md">{countMenunggu}</span>
        </div>
        <div className="bg-indigo-50 border border-indigo-200 px-4 py-2 rounded-xl text-sm font-semibold flex items-center gap-2 text-indigo-700 hidden sm:flex">
          Diproses: <span className="bg-indigo-100 px-2 py-0.5 rounded-md">{countDiproses}</span>
        </div>
        <div className="bg-emerald-50 border border-emerald-200 px-4 py-2 rounded-xl text-sm font-semibold flex items-center gap-2 text-emerald-700 hidden sm:flex">
          Selesai: <span className="bg-emerald-100 px-2 py-0.5 rounded-md">{countSelesai}</span>
        </div>
      </div>

      {/* --- TABEL UTAMA --- */}
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-slate-50">
              <tr>
                <th className="w-20 px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">ID</th>
                <th className="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Judul Gangguan</th>
                <th className="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Kategori</th>
                <th className="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Status</th>
                
                {isAdmin && (
                  <th className="px-6 py-4 text-left text-xs font-bold text-slate-800 uppercase tracking-wider">Aksi Admin</th>
                )}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 bg-white">
              {currentLogs.map((log) => (
                <tr key={log.id} className="hover:bg-slate-50 transition-colors">
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-500 font-medium">#{log.id}</td>
                  <td className="px-6 py-4 text-sm text-slate-900">
                    <div className="font-bold">{log.title}</div>
                    <div className="text-xs text-slate-500 truncate max-w-[200px]">{log.description}</div>
                  </td>
                  
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-500">
                    {log.category_id === 1 ? 'Jaringan' : log.category_id === 2 ? 'Server' : 'Aplikasi'}
                  </td>
                  
                  <td className="px-6 py-4 whitespace-nowrap">
                    <StatusBadge status={log.status} />
                  </td>

                  {/* TOMBOL TINJAU DETAIL */}
                  {isAdmin && (
                    <td className="px-6 py-4 whitespace-nowrap">
                      <button
                        onClick={() => onActionClick(log)}
                        className="bg-white border border-slate-300 text-slate-700 hover:bg-slate-50 hover:border-slate-400 px-4 py-2 rounded-lg text-xs font-bold transition-all flex items-center gap-2 shadow-sm"
                      >
                        🔍 Tinjau Detail
                      </button>
                    </td>
                  )}
                </tr>
              ))}
              
              {currentLogs.length === 0 && (
                <tr>
                  <td colSpan={isAdmin ? 5 : 4} className="px-6 py-10 text-center text-slate-500">
                    Tidak ada laporan gangguan saat ini.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* --- KONTROL PAGINATION --- */}
        {logs.length > 0 && (
          <div className="bg-slate-50 border-t border-slate-200 px-6 py-3 flex items-center justify-between">
            <span className="text-sm text-slate-500">
              Menampilkan <span className="font-bold text-slate-700">{startIndex + 1}</span> hingga <span className="font-bold text-slate-700">{Math.min(startIndex + itemsPerPage, logs.length)}</span> dari <span className="font-bold text-slate-700">{logs.length}</span> entri
            </span>
            
            <div className="flex gap-2">
              <button
                onClick={goToPrev}
                disabled={currentPage === 1}
                className="px-3 py-1.5 border border-slate-300 text-slate-600 rounded-lg text-sm font-semibold hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
              >
                Sebelumnya
              </button>
              <button
                onClick={goToNext}
                disabled={currentPage === totalPages}
                className="px-3 py-1.5 border border-slate-300 text-slate-600 rounded-lg text-sm font-semibold hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
              >
                Selanjutnya
              </button>
            </div>
          </div>
        )}
      </div>

    </div>
  );
}