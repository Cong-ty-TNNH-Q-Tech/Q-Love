import React from 'react';
interface Voucher {
  id: string;
  name: string;
}

interface Props {
  voucher: Voucher;
  onClose: () => void;
}

const MOCK_HISTORY = [
  { id: '1', userName: 'Alex123', date: '2026-08-28 14:30' },
  { id: '2', userName: 'JennyLove', date: '2026-08-27 09:15' },
  { id: '3', userName: 'CuongKenn', date: '2026-08-25 18:45' },
];

const VoucherHistoryModal: React.FC<Props> = ({ voucher, onClose }) => {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-lg p-6 max-h-[80vh] flex flex-col">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold">Lịch sử đổi: {voucher.name}</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div className="overflow-y-auto flex-1">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Người dùng</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Thời gian đổi</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {MOCK_HISTORY.map((h) => (
                <tr key={h.id}>
                  <td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-gray-900">{h.userName}</td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{h.date}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default VoucherHistoryModal;
