import { useState, useEffect } from 'react';

interface Transaction {
  id: string;
  userId: string;
  userName: string;
  type: string;
  amount: number;
  createdAt: string;
}

export function TransactionListPage() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);

  useEffect(() => {
    // Mock API Call
    setTransactions([
      { id: 'txn-1001', userId: 'usr-001', userName: 'Nguyễn Văn A', type: 'deposit', amount: 500, createdAt: '2026-08-28 10:00:00' },
      { id: 'txn-1002', userId: 'usr-002', userName: 'Trần Thị B', type: 'penalty', amount: -50, createdAt: '2026-08-28 11:30:00' },
      { id: 'txn-1003', userId: 'usr-003', userName: 'Lê Văn C', type: 'contract_hold', amount: -100, createdAt: '2026-08-28 14:15:00' },
      { id: 'txn-1004', userId: 'usr-001', userName: 'Nguyễn Văn A', type: 'voucher_exchange', amount: -200, createdAt: '2026-08-28 16:45:00' },
    ]);
  }, []);

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'deposit': return 'Nạp Xu';
      case 'penalty': return 'Phạt Tòa Án';
      case 'contract_hold': return 'Cọc Hẹn Hò';
      case 'voucher_exchange': return 'Đổi Voucher';
      default: return type;
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'deposit': return 'bg-green-100 text-green-800';
      case 'penalty': return 'bg-red-100 text-red-800';
      case 'contract_hold': return 'bg-yellow-100 text-yellow-800';
      case 'voucher_exchange': return 'bg-purple-100 text-purple-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Lịch sử Dòng Tiền</h1>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Mã Giao Dịch</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Người Dùng</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Loại</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Biến động Xu</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Thời gian</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {transactions.map((txn) => (
              <tr key={txn.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 font-mono">
                  {txn.id}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">{txn.userName}</div>
                  <div className="text-xs text-gray-500">{txn.userId}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getTypeColor(txn.type)}`}>
                    {getTypeLabel(txn.type)}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-bold">
                  <span className={txn.amount > 0 ? 'text-green-600' : 'text-red-600'}>
                    {txn.amount > 0 ? '+' : ''}{txn.amount} Xu
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {txn.createdAt}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
