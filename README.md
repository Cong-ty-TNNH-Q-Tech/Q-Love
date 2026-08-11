# Q-Love: Super App Hẹn Hò & Mạng Xã Hội Giải Trí

Chào mừng đến với kho lưu trữ mã nguồn chính thức của **Q-Love**. Dự án này được phát triển bởi Công ty TNHH Q-Tech và vận hành theo mô hình kiến trúc Modular Monolith hiện đại.

## Cấu trúc mã nguồn (Monorepo)

- `/backend/server`: API Backend viết bằng Golang (Framework Fiber). Tương tác với PostgreSQL, PostGIS và Redis.
- `/frontend/app`: Ứng dụng di động đa nền tảng viết bằng Flutter, áp dụng kiến trúc Feature-First.
- `/frontend/admin`: Web Dashboard dành cho quản trị viên, viết bằng React (Vite + TailwindCSS + Zustand).
- `/docs`: Toàn bộ tài liệu nghiệp vụ (BA, BRD), thiết kế hệ thống (Architecture, ERD) và đặc tả giao diện (UI/UX).
- `/docker-compose.yml`: Cấu trúc hạ tầng Local (PostgreSQL + PostGIS, Redis).

## Bắt đầu nhanh (Quick Start)

Yêu cầu môi trường tối thiểu:
- Go 1.22+
- Node.js 20+
- Flutter SDK (stable channel)
- Docker Desktop

Vui lòng tham khảo tài liệu [Tech Stack](docs/tech_stack.md) để biết thêm chi tiết về cách khởi chạy dự án.

---
*Bản quyền © 2026 Q-Tech Team. Áp dụng giấy phép AGPL-3.0.*