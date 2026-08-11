# System Architecture & Tech Stack Document
**Project Name:** Q-Love - Super App Hẹn Hò & Mạng Xã Hội Giải Trí  
**Version:** 1.0  
**Date:** June 2026  
**Language:** Vietnamese (Tiếng Việt)  

Tài liệu này đặc tả quy hoạch công nghệ (Tech Stack) và Kiến trúc hệ thống để đáp ứng các yêu cầu kinh doanh khắt khe (Real-time, Gamification, Widget) đã định nghĩa tại [brd.md](file:///c:/Coding%20Space/Projects/Q-Love/docs/brd.md) và [ba.md](file:///c:/Coding%20Space/Projects/Q-Love/docs/ba.md).

---

## 1. Ứng dụng Di động (Mobile Client)
Sử dụng phương pháp phát triển đa nền tảng (Cross-platform) để tối ưu thời gian ra mắt (Time-to-market), kết hợp với Native Code cho các tính năng đặc thù.

- **Core Framework:** **Flutter**
  - *Lý do:* Q-Love chứa rất nhiều Animation phức tạp (Quẹt thẻ mượt mà, Đảo tình yêu 3D, Hiệu ứng Locket). Engine đồ họa (Skia/Impeller) của Flutter vẽ UI trực tiếp ở tốc độ 60-120fps, mang lại trải nghiệm mượt mà tiệm cận Native tuyệt đối. Ngoài ra, hệ thống cầu nối (Platform Channels) sang code Native (Swift/Kotlin) để làm tính năng Widget Locket cũng cực kỳ ổn định và tài liệu rõ ràng.
- **Widget System (Cốt lõi):**
  - **iOS:** Viết bằng `Swift` sử dụng `WidgetKit`.
  - **Android:** Viết bằng `Kotlin` sử dụng `AppWidgetProvider`.
  - *Lý do:* Flutter/RN không hỗ trợ tạo Widget trực tiếp. Bắt buộc phải viết Native Code và dùng Method Channels để đồng bộ dữ liệu với App chính.
- **Bản đồ & Định vị (Map & GPS):** `Mapbox SDK`.
  - *Lý do:* Tốc độ render hàng ngàn cờ/điểm Check-in của Clan tốt hơn, chi phí API rẻ hơn Google Maps và cho phép custom giao diện bản đồ linh hoạt.

## 2. Hệ thống Máy chủ (Backend & APIs)
Để tối ưu nguồn lực cho giai đoạn MVP (Phase 1), toàn bộ kiến trúc Backend sẽ được gom về một ngôn ngữ duy nhất thay vì dùng quá nhiều loại Tech Stack làm phân mảnh đội Dev.

- **Ngôn ngữ & Framework độc tôn:** **Golang (Go)** kết hợp framework **`Fiber`** *(đã chốt, lý do: Fiber nhanh hơn Gin ~20% trên benchmark, API syntax gần giống Express.js giúp dev onboard nhanh hơn, phù hợp cho workload Real-time cao)*.
  - *Lý do:* Golang là "vua" trong việc xử lý đồng thời (Concurrency) bằng Goroutines. Nó gánh vác cực tốt các kết nối Socket liên tục (Real-time Chat, Matchmaking), tốn ít RAM và chạy rất mượt trên môi trường Cloud.
- **Cách Golang xử lý toàn bộ bài toán của Q-Love:**
  - *Real-time & Core:* Xử lý giao dịch Sàn Chứng Khoán, Cọc Khế ước và Chat thông qua WebSockets.
  - *Admin CMS:* Dùng Golang viết luôn các API CRUD quản trị (chạy rất nhanh và an toàn với tính định kiểu chặt chẽ).
  - *AI & Xử lý ảnh (Thay thế Python):* Thay vì đẻ thêm service Python, Golang sẽ gọi trực tiếp API của OpenAI/Claude cho tính năng Wingman. Với bài toán làm mờ ảnh (Blur) và NSFW, Golang dùng thư viện `bimg` (dựa trên libvips C++) để xử lý ảnh siêu tốc.
- **Giao tiếp Real-time:** `WebSockets` tiêu chuẩn hoặc `gRPC`.

## 3. Cơ sở dữ liệu (Database Layer)
- **Primary Database (Relational):** **PostgreSQL**
  - *Nhiệm vụ:* Lưu trữ thông tin User, Lịch sử giao dịch Xu (ACID), Cổ phiếu Profile.
  - *Plugin bắt buộc:* **PostGIS** để tối ưu hóa truy vấn tọa độ GPS (Tìm người quanh đây, tính bán kính Check-in).
- **In-memory Cache & Queue:** **Redis**
  - *Nhiệm vụ:* Quản lý hàng đợi (Queue) cho Tòa Án Tình Yêu, lưu trữ trạng thái Online/Offline, đếm Streak, và Cache giá Cổ phiếu.

## 4. Hạ tầng & Triển khai (Cloud Infrastructure)
- **Cloud Provider:** **Amazon Web Services (AWS)**.
- **Media Storage:** **Cloudflare R2** hoặc **MinIO** (S3-Compatible).
  - *Lý do:* App mảng hẹn hò/Locket tiêu tốn cực kỳ nhiều băng thông tải ảnh (Egress Bandwidth). AWS S3 thu phí băng thông đầu ra rất đắt.
  - **Cloudflare R2:** Miễn phí 100% băng thông đầu ra (Zero Egress Fee), tích hợp sẵn CDN cực nhanh của Cloudflare. Đây là giải pháp *hoàn hảo* cho tính năng Widget Locket.
  - **MinIO:** Phù hợp nếu tự host (On-premise / VPS) để làm chủ hoàn toàn hạ tầng và tiết kiệm chi phí, tốc độ đọc/ghi nhị phân cực kỳ khủng.
- **Push Notification:** **Firebase Cloud Messaging (FCM)** & **APNs (Apple Push Notification service)**.
  - *Lưu ý:* Phải dùng "Silent Push Notification" để đánh thức (Wake up) ứng dụng ngầm và ra lệnh cho Widget cập nhật ảnh mới.
- **Orchestration & Scaling:** **Docker** và **Kubernetes (K8s)**.
  - *Lý do:* Đảm bảo tính sẵn sàng (High Availability) 99.9%, tự động Scale up (mở rộng server) khi sự kiện "Đêm Săn Mồi" diễn ra.

## 5. Cấu trúc Nhân sự đề xuất (Team Structure)
Để ra mắt Phase 1 (MVP) trong vòng 3 - 4 tháng, đội hình tối thiểu cần có:
1. **Project Manager / BA (1 người):** Theo dõi tiến độ, bám sát các tài liệu BA/BRD.
2. **Mobile Developer (2 người):** Chuyên Flutter, trong đó ít nhất 1 người có khả năng can thiệp Native code (Swift/Kotlin) để làm Widget Locket.
3. **Backend Developer (2 người):** Chuyên Golang. Có kinh nghiệm làm việc với kiến trúc Microservices/Modular Monolith, WebSockets và xử lý ảnh cơ bản.
4. **DevOps / Cloud Engineer (0.5 - 1 người):** Thiết lập hạ tầng AWS, CI/CD, Kubernetes ban đầu.
5. **QA / Tester (1 người):** Viết Test script dựa trên tài liệu `uc_ac.md`.

---

## 6. CI/CD Pipeline

Toàn bộ quy trình tích hợp và triển khai liên tục (CI/CD) được tự động hóa bằng **GitHub Actions**.

```text
[Push / PR to main]
       │
       ▼
[1. CI: Test & Build]
   - go test ./... (Unit + Integration tests)
   - flutter test
   - Lint (golangci-lint, flutter analyze)
       │
       ▼
[2. Docker Build & Push]
   - docker build -t q-love-api:sha-{commit}
   - Push to Amazon ECR
       │
       ▼
[3. Deploy to Staging (Auto)]
   - kubectl set image deployment/q-love-api ... (K8s Rolling Update)
   - Smoke test tự động sau deploy
       │
       ▼
[4. Deploy to Production (Manual Approval)]
   - PM / Tech Lead approve trên GitHub Actions UI
   - Blue/Green Deployment để đảm bảo Zero Downtime
```

**Quy tắc:**
- Branch `main` → Deploy tự động lên **Staging**.
- Tag `v*.*.*` (VD: `v1.2.0`) → Trigger deploy lên **Production** sau khi được approve.

---

## 7. Monitoring & Observability Stack

Hệ thống Real-time với giao dịch Xu ảo phức tạp bắt buộc phải có đầy đủ 3 trụ cột Observability: **Logs, Metrics, Traces**.

| Lớp | Công cụ | Mục đích |
| :--- | :--- | :--- |
| **Error Tracking** | **Sentry** (SDK tích hợp trong Go + Flutter) | Bắt và báo động ngay lập tức khi có crash (App) hoặc panic (Backend). Tích hợp alert vào Slack. |
| **Metrics & Dashboards** | **Prometheus + Grafana** | Theo dõi: Request latency, WebSocket connection count, Worker queue size, Xu transaction volume theo thời gian thực. |
| **Centralized Logging** | **Grafana Loki** | Thu thập log từ tất cả pods K8s vào một nơi. Giúp debug nhanh khi có sự cố (VD: trace luồng giao dịch Xu bị lỗi theo `trace_id`). |
| **Uptime Monitoring** | **Uptime Robot** (Free tier) | Ping endpoint `/health` mỗi 5 phút, SMS alert khi downtime. Đảm bảo SLA 99.9%. |

**Cảnh báo (Alerting Rules) quan trọng:**
- API latency P99 > 500ms → Alert PagerDuty (On-call dev)
- WebSocket active connections > 10,000 → Trigger K8s HPA scale-out
- `wallet_transactions` error rate > 0.1% → Alert khẩn cấp (Liên quan đến tiền)
- Sentry error rate tăng đột biến > 5x baseline → Alert Slack channel #incidents
