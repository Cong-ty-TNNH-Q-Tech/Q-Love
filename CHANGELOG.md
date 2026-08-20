# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - Sắp tới

### Added
- [backend/server/ai] Tích hợp AI xử lý làm mờ ảnh (Gaussian Blur) dựa trên số nguyên âm bằng Go `image/draw` (Issue #21).
- [frontend/app] Hoàn thiện cấu hình App Store Optimization (ASO): Thêm App Icon, Splash Screen, Deep Links (`qlove://match`), Localization (en, vi) và tính năng In-App Review (Issue #121).
- [frontend/app] Locket Widget iOS/Android (Issue #15).
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
- [backend/server/api] Thêm endpoint gửi Locket (Issue #24).
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
- [backend/server/api] Thêm API Minigame Steal (Issue #58) hỗ trợ InitSteal và SubmitStealResult với Transaction `SERIALIZABLE`.

### Changed
- [backend/server/cmd] Refactor `main.go` để loại bỏ global state `app`, chuyển sang DI để hỗ trợ Unit Test độc lập.
- [backend/server/Dockerfile] Cập nhật base image lên `golang:alpine` để tương thích với `go 1.25.0`.
- [docs/readme] Cập nhật toàn bộ README theo chuẩn Investor Ready và thêm logo mới.

### Fixed
- [frontend/app] Khắc phục `ShameWallBloc` dùng mock data (Issue #119): Thêm `ShameRepository`, tách `ShameModel`, gọi API thực sự và thêm empty state UI.
- [backend/server/architecture] Khắc phục vi phạm DIP và Repository pattern ở `AuctionHandler` và `AuctionService` (Issue #113): Thêm `GetActiveAuctions` vào `AuctionService`, tạo `ChatLockRepository`, loại bỏ sử dụng trực tiếp `auctionRepo` trong Handler và `*gorm.DB` trong Service.
- [backend/server] Bổ sung file `.env.example` và thiết lập hệ thống database migrations bằng `golang-migrate` (Issue #120): Thêm Makefile commands và 2 file migration khởi tạo schema.
- [backend/server] Loại bỏ hardcoded `JWT_SECRET` trong `jwt_middleware.go` và yêu cầu tải từ Environment Variable, nếu thiếu sẽ panic khi khởi động (Issue #107).
- [backend/server] Bổ sung header bản quyền AGPLv3 cho các file còn thiếu và thêm pre-commit hook (Issue #115).
- [backend/server] Triển khai `UploadFile` lên Cloudflare R2 cho LocketService thay vì dùng URL ảo (Issue #112).
- [backend/server] Refactor error handling trong `ShameHandler` và `WingmanHandler` (Issue #118): Sử dụng sentinel errors (`ErrInsufficientBalance`, `ErrReferralNotFound`, v.v.) thay vì so sánh string literals.
- [backend/server] Khắc phục vi phạm audit từ issue 74 (thiếu bản quyền, cập nhật CHANGELOG).
- [backend/server/api] Sửa lỗi thiếu `JWTMiddleware` cho các API `/wingman` và `/upload`.
- [backend/server/api] Sửa lỗi sử dụng dummy UUID thay cho JWT context tại `WingmanHandler`.
- [backend/server/tests] Sửa lỗi sai argument count của `go-sqlmock` trong `wingman_service_test.go`.
