import React from 'react';
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import { VoucherListPage } from './pages/vouchers/VoucherListPage';

function App() {
  return (
    <BrowserRouter>
      <div className="flex h-screen bg-gray-100 font-sans">
        {/* Sidebar */}
        <aside className="w-64 bg-gray-900 text-white flex flex-col">
          <div className="p-6 border-b border-gray-800">
            <h1 className="text-2xl font-bold text-pink-500">Q-Love Admin</h1>
          </div>
          <nav className="flex-1 p-4 space-y-2">
            <Link to="/" className="block px-4 py-2 rounded hover:bg-gray-800 text-gray-300 hover:text-white">Dashboard</Link>
            <Link to="/vouchers" className="block px-4 py-2 rounded hover:bg-gray-800 text-gray-300 hover:text-white">Voucher & Vật phẩm</Link>
          </nav>
        </aside>

        {/* Main Content */}
        <main className="flex-1 overflow-y-auto">
          <Routes>
            <Route path="/" element={
              <div className="p-6">
                <h1 className="text-2xl font-bold text-gray-800 mb-4">Dashboard</h1>
                <p className="text-gray-600">Welcome to Q-Love Admin Panel.</p>
              </div>
            } />
            <Route path="/vouchers" element={<VoucherListPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
