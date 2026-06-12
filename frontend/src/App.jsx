  import { useEffect, useState } from 'react';
  import { fetchLogs, createLog, updateStatus, fetchSummary } from './api';
  import Layout from './components/Layout';
  import LogForm from './components/LogForm';
  import LogList from './components/LogList';
  import ChartComponent from './components/ChartComponent';
  import SummaryCard from './components/SummaryCard';
  import TicketModal from './components/TicketModal';

  function App() {
    const [logs, setLogs] = useState([]);
    const [summary, setSummary] = useState(null);
    const [loading, setLoading] = useState(true);
    
    // State untuk menyimpan role yang sedang aktif (Default: Teknisi)
    const [role, setRole] = useState('Teknisi'); 
    const [isFormOpen, setIsFormOpen] = useState(false);
    const [adminTab, setAdminTab] = useState('kendali');
    const [selectedTicket, setSelectedTicket] = useState(null);

    const loadData = async () => {
      try {
        const [logsRes, summaryRes] = await Promise.all([fetchLogs(), fetchSummary()]);
        setLogs(logsRes.data);
        setSummary(summaryRes.data);
      } catch (err) {
        console.error(err);
        alert('Gagal memuat data. Pastikan backend Go berjalan.');
      } finally {
        setLoading(false);
      }
    };

    useEffect(() => {
      loadData();
      // Simulasi Short-polling setiap 10 detik agar data real-time
      const interval = setInterval(loadData, 10000);
      return () => clearInterval(interval);
    }, []);

    const handleCreate = async (formData) => {
      try {
        await createLog(formData);
        await loadData(); 
        alert("Laporan berhasil dikirim!");
        setIsFormOpen(false); // <--- Tambahkan baris ini untuk menutup form otomatis
      } catch (err) {
        alert('Gagal menambah laporan');
      }
    };

    const handleUpdateStatus = async (id, newStatus) => {
      try {
        await updateStatus(id, newStatus);
        await loadData();
      } catch (err) {
        alert('Gagal mengupdate status, coba lagi.');
      }
    };

    if (loading) return <Layout role={role} setRole={setRole}><div className="text-center py-20 text-blue-600 font-bold animate-pulse">Memuat Database...</div></Layout>;

    return (
      <Layout role={role} setRole={setRole}>
        
      {/* ========================================= */}
        {/* V I E W   T E K N I S I                   */}
        {/* ========================================= */}
        {role === 'Teknisi' && (
          <div className="max-w-5xl mx-auto">
            
            {/* SAKLAR: Jika Form Tertutup, tampilkan Tabel */}
            {!isFormOpen ? (
              <div className="space-y-6 animate-in fade-in zoom-in-95 duration-300">
                {/* Header Riwayat & Tombol Buat Laporan */}
                <div className="bg-slate-800 text-white px-6 py-6 rounded-2xl flex flex-col sm:flex-row justify-between items-start sm:items-center shadow-md gap-4">
                  <div>
                    <h2 className="text-xl font-bold">Riwayat Laporan Anda</h2>
                    <p className="text-sm text-slate-300 mt-1">Pemantauan status laporan secara real-time</p>
                  </div>
                  <button
                    onClick={() => setIsFormOpen(true)}
                    className="bg-indigo-500 hover:bg-indigo-400 text-white px-6 py-2.5 rounded-xl font-bold text-sm flex items-center gap-2 transition-all shadow-lg"
                  >
                    <span className="text-lg leading-none">+</span> Buat Tiket Baru
                  </button>
                </div>

                {/* Tabel */}
                <LogList logs={logs} isAdmin={false} />
              </div>
            ) : (
              /* SAKLAR: Jika Form Terbuka, tampilkan Form saja (Terpusat) */
              <div className="max-w-2xl mx-auto animate-in fade-in slide-in-from-bottom-4 duration-300">
                <LogForm 
                  onSubmit={handleCreate} 
                  onCancel={() => setIsFormOpen(false)} // Mengirim sinyal untuk menutup form
                />
              </div>
            )}
          </div>
        )}

        {/* ========================================= */}
        {/* V I E W   A D M I N                       */}
        {/* ========================================= */}
        {role === 'Admin' && (
          <div className="max-w-6xl mx-auto animate-in fade-in duration-300">
            
            {/* Menu Tab Navigasi Admin */}
            <div className="flex gap-2 border-b border-slate-200 mb-6 bg-white p-1 rounded-t-2xl">
              <button 
                onClick={() => setAdminTab('kendali')} 
                className={`flex-1 py-3 px-4 font-bold text-sm sm:text-base rounded-t-xl transition-all border-b-2 ${adminTab === 'kendali' ? 'border-indigo-600 text-indigo-700 bg-indigo-50/50' : 'border-transparent text-slate-500 hover:text-slate-700 hover:bg-slate-50'}`}
              >
                📋 Pusat Kendali Laporan
              </button>
              <button 
                onClick={() => setAdminTab('analitik')} 
                className={`flex-1 py-3 px-4 font-bold text-sm sm:text-base rounded-t-xl transition-all border-b-2 ${adminTab === 'analitik' ? 'border-indigo-600 text-indigo-700 bg-indigo-50/50' : 'border-transparent text-slate-500 hover:text-slate-700 hover:bg-slate-50'}`}
              >
                📈 Analitik & AI Summary
              </button>
            </div>

            {/* KONTEN TAB: PUSAT KENDALI */}
            {adminTab === 'kendali' && (
              <div className="animate-in slide-in-from-left-4 duration-300">
                <LogList 
                  logs={logs} 
                  onActionClick={(ticket) => setSelectedTicket(ticket)} 
                  isAdmin={true} 
                />
              </div>
            )}

            {/* KONTEN TAB: ANALITIK & AI */}
            {adminTab === 'analitik' && (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 items-stretch animate-in slide-in-from-right-4 duration-300">
                <ChartComponent logs={logs} />
                <SummaryCard summary={summary} />
              </div>
            )}

            {/* Komponen Modal Pop-up (Hanya muncul jika selectedTicket tidak null) */}
            <TicketModal 
              ticket={selectedTicket} 
              onClose={() => setSelectedTicket(null)} 
              onUpdateStatus={handleUpdateStatus} 
            />
            
          </div>
        )}

      </Layout>
    );
  }

  export default App;