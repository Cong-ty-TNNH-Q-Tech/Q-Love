# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - Sắp tới

### Added
- [backend/server] Thêm endpoint /shames và sửa transaction, join query, API docs (Issue #74).
- [frontend/app] Tích hợp package `google_fonts` và cấu hình `textTheme` chuẩn cho `AppTheme.darkTheme`.
- [backend/server/pkg] Tích hợp thư viện logging chuẩn `Zap` và báo lỗi qua `Sentry`.
- [backend/server/pkg] Cấu hình và khởi tạo kết nối CSDL PostgreSQL sử dụng `gorm`.
- [backend/server] Thêm headers bản quyền `GNU AGPLv3` cho các file còn thiếu.
- [backend/server/api] Thêm endpoint `/api/v1/upload/presigned-url` cho chức năng Upload Avatar (tích hợp Cloudflare R2).
- [backend/server/api] Thêm endpoint Nghề Cò Mối (Wingman) với đầy đủ CRUD và xử lý hoa hồng Ví Ảo (Transaction SERIALIZABLE).
- [backend/server/models] Thêm model `wingman_referrals` và `wallet_transactions`.
- [backend/server/services] Thêm `WingmanService` với logic Referral và hoa hồng (Commission).
- [frontend/app] Khởi tạo `AppTheme` chuẩn UI/UX Gen-Z (Dark-first, Premium) và thêm headers bản quyền.
- [frontend/app] Khởi tạo kiến trúc State Management bằng BLoC (Thêm `AppBloc`, `MultiBlocProvider`).
- [backend/server/api] Thêm API Tường Thành Phong Sát (Wall of Shame) và luồng giao dịch Ném Cà Chua trừ Xu.
- [frontend/app] Thêm giao diện Tường Thành Phong Sát `ShameWallScreen` (UI Dark-first, Glassmorphism, Haptic Feedback) và `ShameWallBloc`.

### Changed
- [backend/server/cmd] Refactor `main.go` để loại bỏ global state `app`, chuyển sang DI để hỗ trợ Unit Test độc lập.
- [backend/server/Dockerfile] Cập nhật base image lên `golang:alpine` để tương thích với `go 1.25.0`.
- [docs/readme] Cập nhật toàn bộ README theo chuẩn Investor Ready và thêm logo mới.

### Fixed
- [backend/server/api] Sửa lỗi thiếu `JWTMiddleware` cho các API `/wingman` và `/upload`.
- [backend/server/api] Sửa lỗi sử dụng dummy UUID thay cho JWT context tại `WingmanHandler`.
- [backend/server/tests] Sửa lỗi sai argument count của `go-sqlmock` trong `wingman_service_test.go`.
