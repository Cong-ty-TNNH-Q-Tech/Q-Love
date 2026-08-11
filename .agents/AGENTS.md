# Hướng dẫn AI Agent làm việc với Q-Love

Chào AI Agent, khi bạn được giao nhiệm vụ viết code, phân tích hoặc gỡ lỗi trong dự án **Q-Love**, vui lòng tuân thủ tuyệt đối các quy tắc (Rules) dưới đây. Dự án này được thiết kế theo kiến trúc Modular Monolith với việc phân tách rõ ràng Frontend và Backend.

## 1. Kiến trúc Tổng quan (Must Read)
- **Monorepo**: Dự án chia làm 2 phần chính:
  - `frontend/`: Chứa các Client Apps bao gồm `app/` (Mobile App bằng Flutter) và `admin/` (Web Dashboard bằng React/Vite).
  - `backend/`: Chứa `server/` là Backend API viết bằng Golang (Fiber framework).
- **Hệ sinh thái công nghệ**: Hệ thống tích hợp PostgreSQL (kèm PostGIS), Redis (Cache & PubSub), Cloudflare R2 (Media Storage), Mapbox (Bản đồ), ESMS.vn (OTP) và Firebase Cloud Messaging (FCM). **Tuyệt đối tuân thủ kiến trúc đã định sẵn trong `docs/architecture.md`**.
- **Không tự phát minh lại bánh xe**: Sử dụng triệt để sức mạnh của các dịch vụ bên thứ 3 đã quy định.

## 2. Quy tắc cho Frontend
### 2.1. Mobile App (frontend/app)
- **Framework**: Flutter. Quản lý trạng thái bằng BLoC/Provider/Riverpod.
- **UI/UX**: Bắt buộc tuân thủ tài liệu `docs/ui_ux.md`. Áp dụng phong cách **Dark-first**, Premium, Gen-Z aesthetic, sử dụng các hiệu ứng Glassmorphism và viền phát sáng (glow) phù hợp.
- **Bảo mật**: Access Token chỉ lưu trong RAM, Refresh Token lưu an toàn trong `flutter_secure_storage`.

### 2.2. Admin Dashboard (frontend/admin)
- **Framework**: React + Vite + TailwindCSS.
- **Giao tiếp API**: Giao tiếp với backend qua route riêng `/admin/v1/` và bảo vệ bằng JWT (`role: admin`).

## 3. Quy tắc cho Backend (backend/server)
- **Framework**: Golang + Fiber.
- **Kiến trúc**: Sử dụng **Modular Monolith** theo chuẩn `golang-standards/project-layout` (chia rõ `cmd`, `internal`, `pkg`). **Cấm** viết toàn bộ logic vào file main hoặc quăng hết vào thư mục gốc. Tách biệt `Handler`, `Service`, và `Repository`.
- **Database & Xử lý Giao dịch**: PostgreSQL (PostGIS). Mọi giao dịch liên quan đến Ví ảo (Xu, Đóng băng cọc) **BẮT BUỘC** phải sử dụng Transaction của Database ở mức độ cách ly `SERIALIZABLE` để tránh Race Condition.
- **Truy vấn Địa lý**: Bắt buộc sử dụng Index `GIST(location)` của PostGIS cho các truy vấn quét bán kính (radar).
- **API Docs**: Bám sát 100% schema và error codes đã định nghĩa trong `docs/api.yaml`. Nếu có thay đổi code backend làm đổi Response/Request, **phải cập nhật file yaml**.

## 4. Trước khi Code (Cực kỳ Quan trọng)
- Hãy luôn ưu tiên dùng tool để đọc kỹ toàn bộ thư mục `docs/` (bao gồm `brd.md`, `ba.md`, `uc_ac.md`, `erd.md`, `architecture.md`, `tech_stack.md`, `ui_ux.md`, `api.yaml`) trước khi thêm tính năng mới. Không tự ý sáng tạo sai lệch so với đặc tả BA/Product.

## 5. Quy tắc Cập nhật CHANGELOG.md (Bắt buộc)
**BẮT BUỘC** cập nhật `CHANGELOG.md` sau mỗi thay đổi đáng kể.
- Thêm tính năng, sửa API, đổi Schema DB, sửa kiến trúc, sửa bug quan trọng, update Docs lớn → ✅ Bắt buộc update.
- Sửa typo, format code, refactor không đổi behavior → ❌ Không cần.

**Định dạng bắt buộc**: Luôn thêm mục mới vào section `## [Unreleased]` ở đầu file theo dạng Keep a Changelog.
```markdown
## [Unreleased] - Sắp tới
### Added
- [backend/server/api] Thêm endpoint gửi Locket.
### Changed
- [frontend/app] Đổi animation vuốt thẻ trong Discover.
### Fixed
- [docs/erd.md] Sửa lỗi thiếu khóa ngoại.
```

## 6. Nguyên tắc bổ sung (Logging & Data)
- **Logging Standards**: Sử dụng thư viện logging chuẩn (như Zap/Logrus cho Go) để log thông tin hệ thống. Tích hợp cảnh báo lỗi khẩn cấp qua Sentry đúng như thiết kế kiến trúc.
- **Soft Delete (Xóa mềm)**: Mọi dữ liệu nhạy cảm hoặc cốt lõi (Users, Khế ước, Vụ kiện tòa án, Giao dịch) bắt buộc triển khai cơ chế Soft Delete. Không dùng lệnh `DELETE` cứng trực tiếp vào CSDL để tránh mất vết kiểm toán (Audit Trail).

## 7. Tiêu Chuẩn Viết Code (SOLID & Design Patterns)
- **S (Single Responsibility):** Tách biệt API Handler (tiếp nhận HTTP), Business Logic (Service), và Data Access (Repository) ở Backend. Ở Frontend, tách State Management ra khỏi UI widget.
- **D (Dependency Inversion):** Tầng Service chỉ giao tiếp với Database thông qua Interface của Repository. Điều này giúp dễ dàng Unit Test bằng Mock.
- **Repository Pattern:** Bắt buộc dùng ở tầng Backend Data Layer để tập trung mọi logic truy xuất PostgreSQL/Redis.

## 8. Quy trình Giải quyết Issues & Triển khai Task (Bắt buộc)
- Đọc kỹ nội dung yêu cầu của user và kiểm tra các file docs liên quan.
- **Kiểm tra Dependency (Phụ thuộc)**: Nếu task hiện tại phụ thuộc vào task khác chưa hoàn thành (ví dụ: đòi code UI nhưng API backend chưa có, hoặc DB chưa thiết kế), phải DỪNG LẠI và báo cáo cho user, hoặc yêu cầu thực hiện task phụ thuộc trước. 
- **Checklist Trước khi Code**:
  1. Đã đọc tài liệu `docs/` chưa?
  2. Có làm thay đổi cấu trúc DB không? Nếu có, phải cập nhật `docs/erd.md`.
  3. Có thêm sửa API không? Nếu có, phải cập nhật `docs/api.yaml`.

## 9. Quy tắc Commit, Code Review và Pull Request
- Luôn đảm bảo code mới không phá vỡ thiết kế và architecture.
- Chạy Unit Test và Linting (`golangci-lint`, `flutter analyze`) trước khi commit.
- Phải fix hết tất cả các code review comments (nếu có công cụ CI/CD review) thì mới được báo cáo hoàn thành.
