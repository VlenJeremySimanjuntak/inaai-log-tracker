import React, { useState, useEffect } from 'react';
import axios from 'axios';

function App() {
  const [logs, setLogs] = useState([]);

  // Fungsi untuk mengambil data dari backend kita
  const fetchLogs = async () => {
    try {
      const res = await axios.get('http://localhost:8081/api/logs');
      setLogs(res.data);
    } catch (err) {
      console.error("Gagal ambil data:", err);
    }
  };

  // Efek polling setiap 5 detik
  useEffect(() => {
    fetchLogs(); // Panggil pertama kali

    const interval = setInterval(fetchLogs, 5000); // Polling setiap 5 detik
    return () => clearInterval(interval); // Bersihkan interval saat komponen mati
  }, []);

  return (
    <div style={{ padding: '20px' }}>
      <h1>Dashboard Gangguan InaAI</h1>
      <table border="1" cellPadding="10">
        <thead>
          <tr>
            <th>ID</th>
            <th>Judul</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {logs.map((log) => (
            <tr key={log.id}>
              <td>{log.id}</td>
              <td>{log.title}</td>
              <td>{log.status}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default App;