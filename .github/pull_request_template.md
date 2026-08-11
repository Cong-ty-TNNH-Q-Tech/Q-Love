## 📋 Mô tả thay đổi

<!-- Tóm tắt ngắn gọn những gì PR này làm -->

Closes #<!-- số Issue -->

---

## 🏷️ Loại thay đổi

- [ ] 🐛 Bug fix (sửa lỗi, không breaking change)
- [ ] ✨ Feature mới (thêm tính năng)
- [ ] 💥 Breaking change (thay đổi làm hỏng tính năng hiện có)
- [ ] 📝 Tài liệu (chỉ thay đổi docs)
- [ ] ♻️ Refactor (không đổi behavior)
- [ ] 🔧 Config / CI/CD

---

## ✅ Checklist trước khi Submit

### Chung
- [ ] Tôi đã đọc [AGENTS.md](../.agents/AGENTS.md)
- [ ] Không sử dụng `print()` — dùng `logging` chuẩn
- [ ] Mọi dữ liệu lõi đều đã cân nhắc cơ chế Soft Delete (`deleted_at`)

### Backend (backend/server - Golang Fiber)
- [ ] Code tuân thủ Modular Monolith (tách biệt `api/`, `middleware/`, `models/`, `repository/`, `services/`)
- [ ] Giao dịch về Xu/Ví ảo phải dùng Transaction Isolation Level `SERIALIZABLE`
- [ ] Đã cập nhật `docs/api.yaml` nếu thêm/sửa endpoint
- [ ] Tầng Service giao tiếp với Repository qua Interface

### Frontend Mobile App (Flutter)
- [ ] UI/UX bám sát specs trong `docs/ui_ux.md` (Dark-first, Gen-Z aesthetics)
- [ ] Quản lý trạng thái tách biệt khỏi Widget (BLoC / Provider)
- [ ] Token lưu ở RAM, Refresh Token lưu ở secure storage

### Frontend Admin (React/Vite)
- [ ] Tách State Management khỏi UI
- [ ] Call API đúng chuẩn `/admin/v1/`

### Tài liệu
- [ ] Cập nhật `CHANGELOG.md` theo định dạng `Keep a Changelog`
- [ ] Cập nhật ERD Diagram (`docs/erd.md`) nếu có thay đổi Database

---

## 🧪 Cách test

<!-- Mô tả cách reviewer test PR này -->

```bash
# Ví dụ Backend:
cd backend/server
go run cmd/main.go
```

---

## 📸 Screenshots (nếu thay đổi UI)

<!-- Paste ảnh trước/sau ở đây -->
