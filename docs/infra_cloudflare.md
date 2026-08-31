# Cấu hình Hạ tầng Cloudflare cho Q-Love

Tài liệu này hướng dẫn DevOps cách cấu hình Cloudflare để đảm bảo an ninh, hiệu suất và tương thích với ứng dụng Q-Love.

## 1. Cloudflare DNS & Proxy (Orange Cloud)

Để kích hoạt CDN và WAF của Cloudflare, toàn bộ traffic phải đi qua hệ thống proxy của họ.

- **Frontend / Admin Dashboard**: Trỏ domain (vd: `admin.qlove.vn`) thông qua Cloudflare Pages (tự động cấu hình DNS khi deploy).
- **Backend API**:
  - Tạo bản ghi `A` hoặc `CNAME` (vd: `api.qlove.vn`) trỏ về IP của máy chủ Backend.
  - Bật **Proxy status** (đám mây màu cam).

### 1.1 Cấu hình cho WebSocket Hub (wss://)

Hệ thống Chat và Matchmaking của Q-Love sử dụng WebSockets. Cloudflare hỗ trợ WebSockets mặc định trên các cổng tiêu chuẩn (như 443 cho wss://).
- **Yêu cầu:** Backend Server (chạy Fiber) phải lắng nghe WebSockets trên port được Cloudflare hỗ trợ (port 443 hoặc cấu hình port mapping trên server đích).
- Trong giao diện Cloudflare Dashboard -> **Network** -> Bật **WebSockets** (thường đã bật sẵn).
- Lưu ý không sử dụng các port non-standard cho WebSocket nếu bật Proxy Cloudflare, vì Cloudflare chỉ proxy các cổng cụ thể (80, 443, 8080, 8443,...).

## 2. Bảo mật - WAF (Web Application Firewall) & DDoS Protection

### 2.1 Mức độ bảo mật (Security Level)
- Vào Security -> Settings -> Security Level: Chọn **Medium** hoặc **High**.
- Bật **Browser Integrity Check**.

### 2.2 WAF Custom Rules (Tùy chọn)
Nếu phát hiện lưu lượng bất thường vào `/api/v1/auth` hoặc `/api/v1/ws`, hãy tạo Custom Rules (WAF -> Custom rules):
- **Rate Limiting**: Giới hạn số lượng request OTP/Login từ cùng một IP. (Ví dụ: Block IP nếu request `POST /api/v1/auth/otp` quá 10 lần trong 1 phút).
- Block các quốc gia không phải mục tiêu phát hành của Q-Love (nếu app chỉ phục vụ Việt Nam, có thể chặn các IP từ Châu Phi/Đông Âu để giảm tải).

## 3. SSL/TLS

Vì API xử lý dữ liệu nhạy cảm (Tin nhắn, Vị trí, Khế ước), phải bắt buộc dùng HTTPS.
- Vào **SSL/TLS** -> **Overview**: Đặt chế độ **Full (Strict)**.
- Đảm bảo máy chủ Backend (EC2/VPS) có cài đặt chứng chỉ SSL hợp lệ (Let's Encrypt hoặc Origin CA của Cloudflare).
- Bật **Always Use HTTPS** (Edge Certificates).

## 4. Caching

- Mặc định Cloudflare không cache nội dung động (API JSON).
- Không được tạo Page Rules để cache các endpoint API chứa thông tin người dùng (`/api/v1/user/...`).
- Các static assets từ Admin Dashboard (do deploy trên Cloudflare Pages) sẽ tự động được cache tối ưu bởi hệ thống.
