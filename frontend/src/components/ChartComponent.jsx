import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';

export default function ChartComponent({ logs }) {
  // 1. Hitung jumlah tiap status
  const statusCount = { Menunggu: 0, Diproses: 0, Selesai: 0 };
  logs.forEach(log => {
    if (statusCount[log.status] !== undefined) {
      statusCount[log.status]++;
    }
  });

  const total = logs.length;

  // 2. Format data untuk Recharts (hanya tampilkan yang angkanya lebih dari 0 agar rapi)
  const data = Object.keys(statusCount)
    .map(key => ({
      name: key,
      value: statusCount[key]
    }))
    .filter(item => item.value > 0);

  // 3. Warna statis sesuai identitas status
  const COLORS = {
    Menunggu: '#f59e0b', // Amber
    Diproses: '#6366f1', // Indigo
    Selesai: '#10b981'   // Emerald
  };

  // 4. Fungsi untuk menggambar teks (Angka & Persentase) di atas irisan donat
  const renderCustomLabel = ({ cx, cy, midAngle, innerRadius, outerRadius, percent, value }) => {
    const RADIAN = Math.PI / 180;
    const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
    const x = cx + radius * Math.cos(-midAngle * RADIAN);
    const y = cy + radius * Math.sin(-midAngle * RADIAN);

    return (
      <text x={x} y={y} fill="white" textAnchor="middle" dominantBaseline="central" className="text-[11px] font-bold drop-shadow-md">
        {value} ({(percent * 100).toFixed(0)}%)
      </text>
    );
  };

  return (
    <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 h-full flex flex-col">
      <h3 className="text-base font-bold text-slate-800 mb-4">Distribusi Penanganan</h3>
      
      <div className="flex-grow w-full h-[220px] relative">
        {total === 0 ? (
          // Empty State jika tidak ada laporan sama sekali
          <div className="absolute inset-0 flex items-center justify-center text-slate-400 text-sm font-medium">
            Belum ada data laporan
          </div>
        ) : (
          <>
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={data}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={renderCustomLabel} /* Memanggil fungsi label persentase */
                  innerRadius={65}
                  outerRadius={95}
                  paddingAngle={3}
                  dataKey="value"
                  stroke="none"
                >
                  {data.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[entry.name]} />
                  ))}
                </Pie>
                <Tooltip
                  formatter={(value, name) => [`${value} Laporan (${((value / total) * 100).toFixed(1)}%)`, name]}
                  contentStyle={{ borderRadius: '12px', border: 'none', boxShadow: '0 10px 15px -3px rgb(0 0 0 / 0.1)' }}
                  itemStyle={{ color: '#1e293b', fontWeight: 'bold' }}
                />
                <Legend iconType="circle" wrapperStyle={{ fontSize: '12px', fontWeight: '500', paddingTop: '10px' }} />
              </PieChart>
            </ResponsiveContainer>
            
            {/* Teks Angka Total di Tengah Lingkaran Donat */}
            <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none pb-8">
              <span className="text-3xl font-black text-slate-800 leading-none">{total}</span>
              <span className="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-1">Total</span>
            </div>
          </>
        )}
      </div>
    </div>
  );
}