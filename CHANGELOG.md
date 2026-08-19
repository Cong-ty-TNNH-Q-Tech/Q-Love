# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - Sắp tới

### Added
- [frontend/app] UI Quẹt Thẻ & Bộ lọc Tâm Linh (Issue #14).
- [frontend/app] UI Tòa Án Tình Yêu (Vuốt để Vote Tội/Vô tội) (Issue #16).
- [frontend/app] UI Đăng nhập, Xác thực OTP và Tạo Profile (Issue #13).
- [frontend/app] UI Khế Ước Hẹn Hò & Quét QR Code (Issue #46).
- [frontend/app] Giao diện Chat 1-1 & Tích hợp WebSocket, Locket, Wingman (Issue #44).
- [frontend/app] Thêm UI/UX Luồng chia sẻ Deep Link Cò Mối & Quản lý Hoa hồng (Issue #57).
- [backend/server/api] Thêm API backend cho tính năng Vibe Check.
- [frontend/app] Tích hợp API Backend thật cho tính năng Vibe Check (thay thế Mock Data) và thêm UI Vibe Check Screen (Issue #55).
- [backend/server/api] Thêm endpoint kiểm duyệt ảnh NSFW và Locket (Issue #24).
- [frontend/app] Thêm UI/UX Minigame PK Cướp Đoạt Thẻ Bài (Glassmorphism & Haptics).
- [backend/server/api] Thêm API Đấu Giá Đặc Quyền (Blind Auction) và Cronjob khóa chat 24h (Issue #50).
- [backend/server] Thêm endpoint `/shames` và sửa transaction, join query, API docs (Issue #74).
- [backend/server/api] Thêm API tạo và quản lý Clan (Issue #25).
- [backend/server/api] Thêm API tích hợp Locket (Issue #26).
- [backend/server/api] Thêm tính năng chat realtime với Websocket và Redis (Issue #33).
- [backend/server] Tích hợp RevenueCat webhook cho IAP deposit (Issue #32).

### Changed
- [frontend/app] Sử dụng Theme cho TextStyles trong ShameWallScreen để dễ quản lý.
- [backend/server] Sử dụng Zap cho logging và Sentry cho error tracking.
- [backend/server/api] Tách handler, service, repository theo chuẩn Standard Go Project Layout.

### Fixed
- [frontend/app] Sửa lỗi theme chói mắt ở màn hình Tòa án.
- [backend/server] Sửa lỗi IAP Premium Idempotency và thời hạn tự cộng dồn.
- [backend/server] Thêm isolation level SERIALIZABLE vào API ném cà chua (Issue #90).
- [backend/server/api] Ngăn chặn fatal panic do thay đổi map đồng thời trong websocket.
- [backend/server/api] Xóa mềm (Soft Delete) người dùng không phá vỡ logic quan hệ.
- [docs] Sửa lỗi API Docs trùng lặp endpoint và lỗi font.
