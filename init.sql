-- 1. Tabel Master: Users (Teknisi/Admin)
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    role ENUM('Teknisi', 'Admin') DEFAULT 'Teknisi',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Tabel Master: Categories (Kategori Gangguan)
CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. Tabel Transaksi: Incident Logs
CREATE TABLE IF NOT EXISTS incident_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    category_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status ENUM('Menunggu', 'Diproses', 'Selesai') DEFAULT 'Menunggu',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. Tabel Analisa: AI Summaries (Hasil Rangkuman Massal Gemini)
CREATE TABLE IF NOT EXISTS ai_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    summary_text TEXT NOT NULL,
    log_ids_analyzed VARCHAR(255) NOT NULL, -- Menyimpan id gabungan, contoh: "1,2,3"
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Data Awal / Dummy untuk testing
INSERT INTO users (name, role) VALUES ('Vlen Jeremy', 'Admin'), ('Gilberd', 'Teknisi');
INSERT INTO categories (name) VALUES ('Jaringan & Internet'), ('Server & Database'), ('Aplikasi Internal');
INSERT INTO incident_logs (user_id, category_id, title, description, status) VALUES 
(2, 1, 'Koneksi WiFi Drop', 'Access point lantai 2 Gedung Kuliah Utama mati total tidak memancarkan sinyal.', 'Menunggu'),
(2, 2, 'Koneksi MySQL Lambat', 'Terjadi lonjakan query dan memory spikes pada server MySQL utama.', 'Menunggu');