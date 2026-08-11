<div align="center">
  <img src="https://via.placeholder.com/800x200.png?text=Q-Love+Super+App" alt="Q-Love Banner">
  
  # Q-Love: Trải nghiệm Hẹn hò Gamification Đỉnh cao dành cho Gen Z
  
  **Định nghĩa lại giới hạn của ứng dụng hẹn hò bằng hệ sinh thái tài chính ảo, tương tác O2O và trí tuệ nhân tạo.**

  [![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](#)
  [![Flutter](https://img.shields.io/badge/Flutter-Stable-02569B?logo=flutter&logoColor=white)](#)
  [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-PostGIS-336791?logo=postgresql&logoColor=white)](#)
  [![License](https://img.shields.io/badge/License-AGPL%203.0-blue.svg)](#)
  [![Status](https://img.shields.io/badge/Status-Active_Development-success.svg)](#)

</div>

---

## 🌟 Tầm Nhìn (The Vision)

**Q-Love** không chỉ là một ứng dụng hẹn hò (Dating App) vuốt thả nhàm chán. Đây là một **Siêu ứng dụng (Super App) giải trí mạng xã hội** được thiết kế nguyên bản dành riêng cho Gen Z. Chúng tôi kết hợp cơ chế ghép đôi cốt lõi với **Gamification (Trò chơi hóa)**, **Nền kinh tế ảo (Virtual Economy)**, và **Tương tác Offline-to-Online (O2O)** để tạo ra một hệ sinh thái gây nghiện, viral tự nhiên và an toàn.

---

## 🚀 Tính Năng Đột Phá (Breakthrough Features)

Q-Love sở hữu những vũ khí tối thượng đánh bại các ứng dụng truyền thống:

*   **🃏 Chợ Thẻ Bài Profile:** Biến hồ sơ người dùng thành Thẻ Bài (Cards). Người dùng dùng Xu ảo để mua bán, đầu cơ thẻ bài của những Hot boy/Hot girl để kiếm lời. Một nền kinh tế ảo hoàn chỉnh!
*   **🍅 Tường Thành Phong Sát (Wall of Shame):** Trừng phạt những kẻ bùng kèo bằng Tòa Án Tình Yêu. Kẻ thua kiện bị bêu tên công khai 24h và bị cộng đồng "ném cà chua" tốn Xu. Drama không bao giờ tàn!
*   **⚔️ PK Cướp Thẻ (The Steal):** Dùng thẻ Đạo Tặc mua bằng Xu để kích hoạt minigame đối kháng 10 giây. Thắng sẽ cướp trắng Thẻ Bài Profile đắt giá từ tay người khác.
*   **🎧 Vibe Check Nửa Đêm:** Đúng 23:00, "Đài Phát Thanh Tình Yêu" mở cửa. Ghép đôi ẩn danh hoàn toàn dựa trên bài hát đối phương đang nghe qua API Spotify.
*   **👼 Nghề Cò Mối (Wingman):** Bạn đang ế? Không sao, hãy đi "ép duyên" bạn bè. Nếu họ thành đôi và đi chơi, bạn nhận hoa hồng Thẻ bài/Xu. Vòng lặp Viral hoàn hảo!
*   **🤖 Trợ lý "Mỏ Hỗn" (AI Wingman):** Tích hợp LLM (Claude/OpenAI) đọc ngữ cảnh chat để mớm lời mặn mòi, hoặc bói bài Tarot/Chiêm tinh mỗi sáng để phá băng (Ice-breaker).

---

## 🏗️ Kiến Trúc Hệ Thống (Architecture)

Dự án tuân thủ nghiêm ngặt mô hình **Modular Monolith** kết hợp với **Feature-First Architecture** cho tính bền vững và khả năng mở rộng (Scale).

### 🛠 Tech Stack
*   **Core Backend:** Golang (Fiber Framework), Clean Architecture.
*   **Database:** PostgreSQL với PostGIS (Xử lý truy vấn radar bán kính < 50ms).
*   **Caching / PubSub:** Redis (Quản lý trạng thái The Purge, Rate Limiting).
*   **Mobile App:** Flutter (BLoC / Riverpod), Clean Architecture.
*   **Admin Dashboard:** React.js, Vite, TailwindCSS, Zustand.
*   **Infrastructure:** Docker Compose (Local), Kubernetes (Production), Cloudflare R2 (Storage).

---

## 📂 Cấu Trúc Mã Nguồn (Monorepo Structure)

```text
Q-Love/
├── backend/
│   └── server/          # 🚀 Golang Fiber API Server
├── frontend/
│   ├── app/             # 📱 Flutter Mobile Application
│   └── admin/           # 💻 React Web Admin Dashboard
├── docs/                # 📚 Toàn bộ tài liệu nghiệp vụ & kỹ thuật
│   ├── api.yaml         # Đặc tả OpenAPI 3.0
│   ├── architecture.md  # Sơ đồ Kiến trúc & Luồng dữ liệu
│   ├── ba.md            # Phân tích Nghiệp vụ (Use-cases)
│   ├── brd.md           # Business Requirements (Mô hình kiếm tiền)
│   ├── erd.md           # Database Schema
│   └── ...
├── .github/workflows/   # ⚙️ CI/CD Pipelines (Github Actions)
└── docker-compose.yml   # 🐳 Cấu hình môi trường Local (DB, Redis)
```

---

## 🏁 Bắt Đầu Nhanh (Quick Start)

### 1. Khởi chạy Hạ tầng (Local Infrastructure)
Bạn cần cài đặt **Docker** và **Docker Compose**.
```bash
docker-compose up -d
```
*Lệnh này sẽ khởi chạy PostgreSQL (kèm PostGIS) ở port 5432 và Redis ở port 6379.*

### 2. Khởi chạy Backend API (Golang)
```bash
cd backend/server
go mod tidy
go run cmd/main.go
```
*API sẽ lắng nghe tại `http://localhost:3000`.*

### 3. Khởi chạy Mobile App (Flutter)
```bash
cd frontend/app
flutter pub get
flutter run
```

---

## 📖 Tài Liệu Tham Khảo Kỹ Thuật

Mọi thành viên tham gia dự án **BẮT BUỘC** phải đọc các tài liệu sau trước khi đóng góp mã nguồn:
1.  [Cẩm nang gọi vốn (BRD)](docs/brd.md) - Hiểu về cách dự án kiếm tiền.
2.  [Đặc tả hệ thống (Architecture)](docs/architecture.md) - Hiểu luồng dữ liệu cốt lõi.
3.  [Tiêu chuẩn Giao diện (UI/UX)](docs/ui_ux.md) - Đảm bảo thẩm mỹ Gen-Z, Dark-first.

---

<div align="center">
  <i>Được xây dựng bằng 💻 và ☕ bởi đội ngũ tinh nhuệ của <b>Q-Tech</b>.</i><br>
  Bản quyền © 2026 Q-Tech Team.
</div>