# Business Requirements Document (BRD)
**Project Name:** Q-Love - Super App Hẹn Hò & Mạng Xã Hội Giải Trí  
**Version:** 1.0  
**Date:** June 2026  
**Status:** Draft  
**Language:** Vietnamese (Tiếng Việt)  

---

## Table of Contents
1. [Tổng quan dự án (Executive Summary)](#1-tổng-quan-dự-án-executive-summary)
2. [Mục tiêu kinh doanh (Business Objectives)](#2-mục-tiêu-kinh-doanh-business-objectives)
3. [Đối tượng khách hàng mục tiêu (Target Audience)](#3-đối-tượng-khách-hàng-mục-tiêu-target-audience)
4. [Phạm vi dự án (Project Scope)](#4-phạm-vi-dự-án-project-scope)
5. [Tổng quan yêu cầu nghiệp vụ (Business Requirements)](#5-tổng-quan-yêu-cầu-nghiệp-vụ-business-requirements)
6. [Mô hình doanh thu (Monetization Strategy)](#6-mô-hình-doanh-thu-monetization-strategy)
7. [Yêu cầu phi chức năng (Non-Functional Requirements)](#7-yêu-cầu-phi-chức-năng-non-functional-requirements)
8. [Rủi ro & Biện pháp giảm thiểu (Risks & Mitigations)](#8-rủi-ro--biện-pháp-giảm-thiểu-risks--mitigations)

---

## 1. Tổng quan dự án (Executive Summary)
Q-Love là một dự án "Super App" kết hợp giữa mô hình ứng dụng hẹn hò truyền thống (Dating) và mạng xã hội giải trí đa nền tảng (Social Entertainment). 

**Vấn đề (Pain Points) của thị trường hiện tại:**
- Người dùng các app hẹn hò truyền thống (Tinder, Bumble) thường xuyên gặp tình trạng **Ghosting** (nhắn tin đột ngột biến mất) hoặc **Flaking / Bùng kèo** (hẹn nhưng không đến).
- Trải nghiệm quẹt thẻ dần trở nên nhàm chán, thiếu sự tương tác sâu và thiếu động lực gắn kết lâu dài.

**Giải pháp của Q-Love:**
Tích hợp mạnh mẽ các cơ chế **Gamification** (Game hóa) và **Virtual Economy** (Nền kinh tế ảo) vào quá trình hẹn hò. Các tính năng nổi bật bao gồm "Tòa án tình yêu" trừng phạt kẻ bùng kèo, "Sàn chứng khoán Profile" giúp người dùng đầu tư vào sự nổi tiếng của người khác, và "Khế ước tài chính" để cọc tiền đảm bảo các cuộc hẹn O2O (Online to Offline).

---

## 2. Mục tiêu kinh doanh (Business Objectives)
- **User Acquisition:** Đạt 1 triệu người dùng đăng ký (Registered Users) và 300,000 người dùng hoạt động hàng tháng (MAU) trong 6 tháng đầu ra mắt.
  - *Đo lường:* Amplitude (event: `app_open`, `registration_completed`). Báo cáo hàng tuần. Chủ sở hữu: PM.
- **Retention (Giữ chân):** Tăng tỷ lệ giữ chân ngày 1 (D1 Retention) lên 45% và ngày 30 (D30 Retention) lên 25% nhờ vào cơ chế chuỗi tương tác (Streak) và Đảo Tình Yêu.
  - *Đo lường:* Amplitude Retention Analysis (Cohort). Báo cáo hàng tháng. Chủ sở hữu: Product.
- **Monetization (Doanh thu):** Đạt ARPU (Doanh thu trung bình trên mỗi người dùng) tối thiểu $2.5/tháng thông qua hệ thống In-App Purchase và mua bán vật phẩm nội bộ.
  - *Đo lường:* RevenueCat (IAP analytics). Báo cáo hàng tháng. Chủ sở hữu: Business.
- **Thương hiệu:** Trở thành Top 3 ứng dụng Hẹn Hò & Giải Trí được tải nhiều nhất trên App Store / Google Play tại Việt Nam sau 1 năm.
  - *Đo lường:* App Annie / data.ai rank tracking. Theo dõi hàng tuần.

---

## 3. Đối tượng khách hàng mục tiêu (Target Audience)
- **Độ tuổi:** Gen Z và Millennials (18 - 30 tuổi).
- **Vị trí địa lý:** Tập trung ban đầu tại các thành phố lớn (Hà Nội, TP.HCM, Đà Nẵng) trước khi mở rộng.
- **Đặc điểm hành vi:**
  - Thích trải nghiệm các tính năng mới lạ, trending (Locket, đầu tư tài chính ảo).
  - Có nhu cầu kết bạn, hẹn hò nhưng e ngại các rủi ro, lừa đảo hoặc sự nhàm chán của app cũ.
  - Sẵn sàng chi trả (Micro-transactions) cho các dịch vụ tăng cường trải nghiệm cá nhân hoặc trừng phạt/tương tác với người khác.

---

## 4. Phạm vi dự án (Project Scope)

### 4.1. In-Scope (Trong phạm vi triển khai)
Các tính năng được quy hoạch theo 3 giai đoạn (Phases) tương ứng với tài liệu Đặc tả Use-case [ba.md](./ba.md):

**Phase 1 (MVP Core - Nền móng Viral):**
- Đăng ký và Quẹt thẻ tìm đối tượng (Matchmaking cơ bản).
- Gửi ảnh Blind Locket & Đồng bộ Widget màn hình khóa.
- Thành lập Bang hội (Clan) và Đua top.
- Tòa Án Tình Yêu (xử lý hành vi Ghosting).

**Phase 2 (Giữ chân & Doanh thu O2O):**
- Đảo Tình Yêu 3D (Streak Gamification).
- Đánh chiếm địa bàn qua định vị GPS.
- Khế Ước Tài Chính chống bùng kèo bằng QR Code.
- Trợ lý AI (Wingman).

**Phase 3 (Kinh tế ảo & Thao túng):**
- Đánh giá và Tra cứu CV Tình Trường.
- Sàn Chứng Khoán Độc Thân (Mua bán cổ phần profile).
- Đêm săn mồi (The Purge).

**Hệ thống chung (Phát triển xuyên suốt):**
- Hệ thống Nạp/Tiêu Xu ảo (Token) liên kết Voucher (Highlands, CGV,...).
- Hệ thống quản trị (Admin/CMS Dashboard): Kiểm duyệt nội dung, xử lý khiếu nại (appeal) từ Tòa Án, đối soát Voucher và quản trị kinh tế ảo.

### 4.2. Out-of-Scope (Nằm ngoài phạm vi giai đoạn đầu)
- Tính năng Livestream (Sẽ được cân nhắc ở Phase 4).
- Tính năng gọi điện Video/Voice Call (Sử dụng rất nhiều băng thông, tạm hoãn để tối ưu chi phí hạ tầng).
- Cổng rút tiền thật trực tiếp từ Xu ảo ra VNĐ (Tránh vi phạm quy định chống rửa tiền/cờ bạc của Apple/Google).

---

## 5. Tổng quan yêu cầu nghiệp vụ (Business Requirements)
Dựa trên tài liệu BA Use-case, các yêu cầu nghiệp vụ chính được thiết kế để giải quyết 3 vòng lặp cốt lõi (Core Loops) của hệ thống:

1. **Vòng lặp Viral & Tương tác (Acquisition & Engagement):** 
   - Sử dụng Blind Locket Widget để ép người dùng phải mở app hàng ngày để xem ảnh (Nuôi Streak).
   - "Đâm đơn kiện" kẻ Ghosting để tạo Drama ẩn danh, thu hút người dùng khác vào đọc (Crowdsourcing content).
2. **Vòng lặp Offline (O2O & Gamification):** 
   - Kéo người dùng từ nhà ra đường thông qua "Đánh chiếm địa bàn Clan".
   - Khuyến khích gặp mặt thực tế an toàn bằng "Khế ước tài chính chống bùng kèo".
3. **Vòng lặp Kinh tế ảo (Virtual Economy & Monetization):** 
   - Sử dụng Sàn Chứng Khoán Profile để tạo tính thanh khoản và lý do để nạp tiền, giữ xu trong hệ thống thay vì tiêu xài hết.

---

## 6. Mô hình doanh thu (Monetization Strategy)
Q-Love sử dụng mô hình "Freemium + Micro-transactions" kết hợp kinh tế ảo nội bộ:

- **Gói Đăng ký định kỳ (Subscription - Q-Love Premium):** Nguồn thu cốt lõi ổn định hàng tháng. Cung cấp các đặc quyền VIP: Quẹt không giới hạn, xem ai đã Like mình, miễn phí 1 lần hủy kèo/tháng không bị trừ cọc, v.v.
- **Bán Xu Ảo (Tokens):** Thông qua cổng In-App Purchase (IAP). Xu là đơn vị tiền tệ duy nhất dùng cho mọi giao dịch nhỏ lẻ (Micro-transactions).
- **Bán Vật Phẩm Đặc Biệt:**
  - *Giấy giảm án / Quà tạ lỗi:* Tội phạm bị đưa ra Tòa Án Tình Yêu phải nạp tiền mua giấy giảm án để xóa huy hiệu "Ghost thủ".
  - *Bộ lọc ảnh Locket VIP:* Các hiệu ứng làm mờ đặc biệt hoặc khung ảnh đẹp.
- **Thu phí giao dịch (Transaction Fee - "Thuế"):**
  - Thu 10% phí quản lý khi có người bị tịch thu cọc hẹn hò (chuyển 90% cho nạn nhân, 10% chảy vào quỹ hệ thống).
  - Thu phí giao dịch 2% cho mỗi lệnh mua/bán thành công trên Sàn Chứng Khoán Profile.
- **Doanh thu B2B (Tương lai):** Bán các gói "Cờ chủ quyền địa điểm" cho các quán Cafe, Trà sữa để kích cầu người dùng tới Check-in.

---

## 7. Yêu cầu phi chức năng (Non-Functional Requirements)
- **Hiệu năng & Tính mở rộng (Performance & Scalability):** 
  - Ảnh Locket đẩy qua Widget phải được server xử lý làm mờ và gửi đến thiết bị đích trong dưới `3 giây`.
  - Bản đồ GPS có khả năng render `10,000` cờ địa điểm cùng lúc mà không gây crash app.
  - Hệ thống phải có khả năng Tự động mở rộng (Auto-scaling) để xử lý lượng truy cập tăng vọt (Spike traffic) trong các sự kiện đua top hoặc khi có vụ kiện Hot.
- **Tính sẵn sàng (High Availability):** Đảm bảo Uptime đạt mức `99.9%` để duy trì trải nghiệm Widget liên tục (không bị mất kết nối thời gian thực).
- **Bảo mật & Quyền riêng tư (Security & Privacy):**
  - Che mờ (Redact) 100% các dữ liệu PII (Avatar, tên, SĐT) trong đoạn chat bằng regex trước khi đưa vụ kiện ra Tòa Án Tình Yêu.
  - Sử dụng hệ thống phát hiện "Mock Location" để khóa ngay lập tức các tài khoản dùng phần mềm fake GPS.
- **Pháp lý (Compliance):** 
  - Thiết kế quy trình "Cọc - Đền bù" khéo léo, chỉ cho phép đổi Xu ra Voucher dịch vụ thực tế, nghiêm cấm quy đổi ra tiền mặt để tránh vi phạm chính sách Cờ Bạc (Gambling) của App Store và luật pháp sở tại.

---

## 8. Rủi ro & Biện pháp giảm thiểu (Risks & Mitigations)
| Rủi ro | Mức độ | Biện pháp giảm thiểu |
| :--- | :--- | :--- |
| **Bị Apple/Google Store từ chối do tính năng Sàn Chứng Khoán** | Cao | Đóng gói tính năng này dưới dạng "Minigame sưu tầm thẻ bài Profile", không dùng các từ ngữ như Trading, Stock, Invest trong UI/UX public. |
| **Người dùng lạm dụng Tòa Án Tình Yêu để spam, bôi nhọ** | Trung bình | Yêu cầu phải có Streak tối thiểu 5 ngày mới được kiện. Áp dụng AI quét nội dung ngôn từ thù địch trước khi đưa vụ kiện ra công khai. |
| **Gian lận khi quét QR Code Khế Ước (chụp gửi qua mạng)** | Thấp | Bắt buộc sử dụng Dynamic QR Code (đổi mã mỗi 30 giây) ngay trong app, vô hiệu hóa tính năng chụp màn hình (Screenshot) tại màn hình QR. |
| **Rủi ro Pháp lý tại thị trường Việt Nam (Giấy phép)** | Cao | App có tính năng Mạng xã hội và giao dịch Xu ảo. Bắt buộc phải có cố vấn pháp lý ngay từ Phase 1 để xin cấp phép Giấy phép Mạng Xã Hội và Thương mại điện tử (nếu cần). |
