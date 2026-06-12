import { useState } from 'react';

export default function LogForm({ onSubmit, onCancel }) {
  const [form, setForm] = useState({ user_id: 1, category_id: 1, title: '', description: '' });
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value });

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!form.title.trim()) return alert('Judul gangguan wajib diisi');
    setLoading(true);
    await onSubmit(form);
    setLoading(false);
    setForm({ ...form, title: '', description: '' });
  };

  return (
    <div className="bg-white rounded-2xl border border-slate-200 p-8 shadow-sm">
      <div className="mb-6">
        <h2 className="text-xl font-bold text-slate-800">Buat Laporan Baru</h2>
        <p className="text-sm text-slate-500 mt-1">Laporkan gangguan operasional di lapangan.</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-5">
        <div>
          <label className="block text-sm font-semibold text-slate-700 mb-1.5">Judul Gangguan</label>
          <input
            name="title"
            placeholder="Maksimal 50 karakter..."
            value={form.title}
            onChange={handleChange}
            className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-sm focus:bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all"
          />
        </div>

        {/* --- BAGIAN LABEL YANG DIPERBAIKI --- */}
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="w-full sm:w-2/3">
            <label className="block text-xs font-bold text-slate-500 mb-1.5 uppercase tracking-wide">Kategori Gangguan</label>
            <select
              name="category_id"
              value={form.category_id}
              onChange={handleChange}
              className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-sm text-slate-700 focus:bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500"
            >
              <option value="1">🌐 Jaringan & Internet</option>
              <option value="2">🖥️ Server & Database</option>
              <option value="3">📱 Aplikasi Internal</option>
            </select>
          </div>
          
          <div className="w-full sm:w-1/3">
            <label className="block text-xs font-bold text-slate-500 mb-1.5 uppercase tracking-wide">Nama Teknisi</label>
            <select
              name="user_id"
              value={form.user_id}
              onChange={handleChange}
              className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-sm text-slate-700 focus:bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500"
            >
              <option value="1">Vlen</option>
              <option value="2">Gilberd</option>
            </select>
          </div>
        </div>
        {/* ------------------------------------ */}

        <div>
          <label className="block text-sm font-semibold text-slate-700 mb-1.5">Deskripsi Detail</label>
          <textarea
            name="description"
            placeholder="Jelaskan kronologi, lokasi, atau dampak dari gangguan ini..."
            rows="4"
            value={form.description}
            onChange={handleChange}
            className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-sm focus:bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all resize-none"
          />
        </div>

        <div className="flex gap-3 pt-2">
          {onCancel && (
            <button
              type="button"
              onClick={onCancel}
              className="w-1/3 bg-white border border-slate-300 text-slate-700 font-semibold py-3 rounded-xl hover:bg-slate-50 transition-all text-sm"
            >
              Kembali
            </button>
          )}
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-indigo-600 text-white font-semibold py-3 rounded-xl hover:bg-indigo-700 active:scale-[0.98] transition-all disabled:opacity-50 flex justify-center items-center text-sm shadow-md shadow-indigo-200"
          >
            {loading ? 'Memproses...' : 'Kirim Tiket Laporan'}
          </button>
        </div>
      </form>
    </div>
  );
}