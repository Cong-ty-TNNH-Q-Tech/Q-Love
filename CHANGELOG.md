# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - Sắp tới

### Added
- [frontend/app] Tích hợp tính năng UI Quẹt Thẻ (Tinder-like swipe) & Bộ lọc Tâm Linh (Issue #14).
- [backend/server/api] Thêm endpoint `/ai/suggest` để Trợ lý Mỏ Hỗn sinh gợi ý tin nhắn.
- [frontend/app] Tạo thư mục native (android/ và ios/) cho Flutter app, cấu hình Bundle ID, Permissions và Deep Links (Issue #111).
- [backend/server/api] Thêm API Unmatch (Issue #42).
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
- [backend/server/api] Thêm API đánh giá Ex-Rating ẩn danh sau Unmatch (`POST /ex-ratings`).
- [backend/server/api] Thêm API tra cứu CV Tình Trường (Ex-Rating) tốn 50 Xu (`GET /users/:user_id/ex-rating`).
- [backend/server/api] Thêm endpoint `/ai/suggest` để Trợ lý Mỏ Hỗn sinh gợi ý tin nhắn.
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
- [backend/server/cmd] Refactor `main.go` để loại bỏ global state `app`, chuyển sang DI để hỗận Unit Test độc lập.
- [backend/server/Dockerfile] Cập nhật base image lên `golang:alpine` để tương thích với `go 1.25.0`.
- [docs/readme] Cập nhật toàn bộ README theo chuẩn Investor Ready và thêm logo mới.
- [backend/server] Bổ sung file `.env.example` và thiết lập hệ thống database migrations bằng `golang-migrate` (Issue #120): Thêm Makefile commands và 2 file migration khởi tạo schema.
- [backend/server] Bổ sung header bản quyền AGPLv3 cho các file còn thiếu và thêm pre-commit hook (Issue #115).
- [backend/server] Triển khai `UploadFile` lên Cloudflare R2 cho LocketService thay vì dùng URL ảo (Issue #112).
- [backend/server] Refactor error handling trong `ShameHandler` và `WingmanHandler` (Issue #118): Sử dụng sentinel errors (`ErrInsufficientBalance`, `ErrReferralNotFound`, v.v.) thay vì so sánh string literals.

### Fixed
- [frontend/app] Khắc phục `ShameWallBloc` dùng mock data (Issue #119): Thêm `ShameRepository`, tách `ShameModel`, gọi API thực sự và thêm empty state UI.
- [backend/server/services] Sửa lỗi logic `WingmanService.AcceptReferral` (Issue #117): Thực sự tạo `Match` trong database, cập nhật trạng thái `referral`, và phân phối hoa hồng một cách chính xác thay vì chỉ trả về dummy `MatchID`.
- [backend/server/security] Implement actual NSFW detection using AWS Rekognition in `NSFWService` (Issue #116): Fixes the mock implementation to perform real image content moderation.
- [backend/server/services] Sửa lỗi kiểu dữ liệu LocketRateLimiter (`uint64` sang `int64`) gây panic khi so sánh (Issue #110).
- [backend/server/architecture] Khắc phục vi phạm DIP và Repository pattern ở `AuctionHandler` và `AuctionService` (Issue #113): Thêm `GetActiveAuctions` vào `AuctionService`, tạo `ChatLockRepository`, loại bỏ sử dụng trực tiếp `auctionRepo` trong Handler và `*gorm.DB` trong Service.
- [backend/server] Loại bỏ hardcoded `JWT_SECRET` trong `jwt_middleware.go` và yêu cầu tải từ Environment Variable, nếu thiếu sẽ panic khi khởi động (Issue #107).

- [backend/server] Khắc phục vi phạm audit từ issue 74 (thiếu bản quyền, cập nhật CHANGELOG).
- [backend/server/api] Sửa lỗi thiếu `JWTMiddleware` cho các API `/wingman` và `/upload`.
- [backend/server/api] Sửa lỗi sử dụng dummy UUID thay cho JWT context tại `WingmanHandler`.
- [backend/server/tests] Sửa lỗi sai argument count của `go-sqlmock` trong `wingman_service_test.go`.
- [backend/server/api] Sửa lỗi bảo mật auth bypass trong WebSocket chat (Issue #109): Xác thực user claim từ JWT payload trực tiếp trong socket handler thay vì phụ thuộc vào query parameter không an toàn.
