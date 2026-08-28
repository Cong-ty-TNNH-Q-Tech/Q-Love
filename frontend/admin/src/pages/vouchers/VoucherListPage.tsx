import React, { useState } from 'react';
import VoucherFormModal from './components/VoucherFormModal';
import VoucherHistoryModal from './components/VoucherHistoryModal';

interface Voucher {
  id: string;
  name: string;
  description: string;
  cost: number;
  quantity: number;
  status: 'active' | 'inactive';
}

const MOCK_VOUCHERS: Voucher[] = [
  { id: '1', name: 'Highlands Coffee 50k', description: 'Voucher giảm giá 50k toàn hệ thống', cost: 1000, quantity: 50, status: 'active' },
  { id: '2', name: 'CGV Ticket', description: 'Vé xem phim 2D các ngày trong tuần', cost: 2500, quantity: 20, status: 'active' },
  { id: '3', name: 'Shopee 20k', description: 'Mã giảm giá 20k đơn từ 0đ', cost: 500, quantity: 0, status: 'inactive' },
];

export const VoucherListPage = () => {
  const [vouchers, setVouchers] = useState<Voucher[]>(MOCK_VOUCHERS);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isHistoryOpen, setIsHistoryOpen] = useState(false);
  const [selectedVoucher, setSelectedVoucher] = useState<Voucher | null>(null);

  const handleAdd = () => {
    setSelectedVoucher(null);
    setIsFormOpen(true);
  };

  const handleEdit = (v: Voucher) => {
    setSelectedVoucher(v);
    setIsFormOpen(true);
  };

  const handleDelete = (id: string) => {
    if (confirm('Bạn có chắc chắn muốn xóa voucher này?')) {
      setVouchers(vouchers.filter(v => v.id !== id));
    }
  };

  const handleViewHistory = (v: Voucher) => {
    setSelectedVoucher(v);
    setIsHistoryOpen(true);
  };

  const handleSave = (voucher: Voucher) => {
    if (selectedVoucher) {
      setVouchers(vouchers.map(v => v.id === voucher.id ? voucher : v));
    } else {
      setVouchers([...vouchers, { ...voucher, id: Math.random().toString() }]);
    }
    setIsFormOpen(false);
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Quản lý Voucher & Vật phẩm</h1>
        <button 
          onClick={handleAdd}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg shadow font-medium"
        >
          + Thêm Mới
        </button>
      </div>

      <div className="bg-white rounded-xl shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tên Voucher</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Giá Xu</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Kho</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Trạng thái</th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Hành động</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {vouchers.map((v) => (
              <tr key={v.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">{v.name}</div>
                  <div className="text-sm text-gray-500">{v.description}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-yellow-600 font-bold">
                  {v.cost.toLocaleString()} Xu
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {v.quantity}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${v.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                    {v.status === 'active' ? 'Hoạt động' : 'Hết hàng'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button onClick={() => handleViewHistory(v)} className="text-indigo-600 hover:text-indigo-900 mr-4">Lịch sử</button>
                  <button onClick={() => handleEdit(v)} className="text-blue-600 hover:text-blue-900 mr-4">Sửa</button>
                  <button onClick={() => handleDelete(v.id)} className="text-red-600 hover:text-red-900">Xóa</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {isFormOpen && (
        <VoucherFormModal 
          voucher={selectedVoucher} 
          onClose={() => setIsFormOpen(false)} 
          onSave={handleSave} 
        />
      )}

      {isHistoryOpen && selectedVoucher && (
        <VoucherHistoryModal 
          voucher={selectedVoucher} 
          onClose={() => setIsHistoryOpen(false)} 
        />
      )}
    </div>
  );
};
