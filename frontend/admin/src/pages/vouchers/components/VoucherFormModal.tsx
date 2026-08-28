import React, { useState, useEffect } from 'react';

interface Voucher {
  id: string;
  name: string;
  description: string;
  cost: number;
  quantity: number;
  status: 'active' | 'inactive';
}

interface Props {
  voucher: Voucher | null;
  onClose: () => void;
  onSave: (v: Voucher) => void;
}

const VoucherFormModal: React.FC<Props> = ({ voucher, onClose, onSave }) => {
  const [formData, setFormData] = useState<Voucher>({
    id: '',
    name: '',
    description: '',
    cost: 0,
    quantity: 0,
    status: 'active',
  });

  useEffect(() => {
    if (voucher) {
      setFormData(voucher);
    }
  }, [voucher]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(formData);
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-md p-6">
        <h2 className="text-xl font-bold mb-4">{voucher ? 'Chỉnh sửa Voucher' : 'Thêm mới Voucher'}</h2>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">Tên Voucher</label>
            <input 
              type="text" 
              required
              className="w-full border border-gray-300 rounded-md px-3 py-2"
              value={formData.name}
              onChange={(e) => setFormData({...formData, name: e.target.value})}
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">Mô tả</label>
            <textarea 
              required
              className="w-full border border-gray-300 rounded-md px-3 py-2"
              value={formData.description}
              onChange={(e) => setFormData({...formData, description: e.target.value})}
            />
          </div>
          <div className="flex space-x-4 mb-4">
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-700 mb-1">Giá Xu</label>
              <input 
                type="number" 
                required min="0"
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                value={formData.cost}
                onChange={(e) => setFormData({...formData, cost: Number(e.target.value)})}
              />
            </div>
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-700 mb-1">Số lượng</label>
              <input 
                type="number" 
                required min="0"
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                value={formData.quantity}
                onChange={(e) => setFormData({...formData, quantity: Number(e.target.value)})}
              />
            </div>
          </div>
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-1">Trạng thái</label>
            <select 
              className="w-full border border-gray-300 rounded-md px-3 py-2"
              value={formData.status}
              onChange={(e) => setFormData({...formData, status: e.target.value as any})}
            >
              <option value="active">Hoạt động</option>
              <option value="inactive">Hết hàng</option>
            </select>
          </div>
          
          <div className="flex justify-end space-x-3">
            <button type="button" onClick={onClose} className="px-4 py-2 text-gray-600 bg-gray-100 rounded-md hover:bg-gray-200">Hủy</button>
            <button type="submit" className="px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700">Lưu</button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default VoucherFormModal;
