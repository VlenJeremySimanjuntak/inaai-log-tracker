    export default function Layout({ children, role, setRole }) {
    return (
        <div className="min-h-screen bg-gray-100 font-sans">
        <header className="bg-gradient-to-r from-blue-800 to-blue-600 text-white shadow-lg">
            <div className="container mx-auto px-4 py-5 flex flex-col sm:flex-row justify-between items-center">
            <h1 className="text-2xl font-bold tracking-tight">📡 Log Tracker - Operasional</h1>
            
            {/* Tombol Pilihan Role */}
            <div className="flex gap-3 mt-4 sm:mt-0 bg-black/20 p-1 rounded-full">
                <button
                onClick={() => setRole('Teknisi')}
                className={`px-5 py-1.5 rounded-full text-sm font-semibold transition-all duration-300 ${
                    role === 'Teknisi' ? 'bg-white text-blue-800 shadow-md' : 'text-white hover:bg-white/20'
                }`}
                >
                👷 Teknisi
                </button>
                <button
                onClick={() => setRole('Admin')}
                className={`px-5 py-1.5 rounded-full text-sm font-semibold transition-all duration-300 ${
                    role === 'Admin' ? 'bg-white text-blue-800 shadow-md' : 'text-white hover:bg-white/20'
                }`}
                >
                📊 Admin
                </button>
            </div>

            </div>
        </header>
        <main className="container mx-auto px-4 py-8">
            {children}
        </main>
        <footer className="text-center text-gray-500 text-sm py-4 border-t mt-8">
            © 2026 Log Tracker - Arsitektur Clean & Role-Based UI
        </footer>
        </div>
    );
    }