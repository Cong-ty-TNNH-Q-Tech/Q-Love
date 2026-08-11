# Use-case document: Q-Love - Super App Hẹn Hò & Mạng Xã Hội Giải Trí

**Version:** 1.0  
**Date:** June 2026  
**Status:** Draft  
**Language:** Vietnamese (Tiếng Việt)  

## Table of Contents
1. [Actor (Tác nhân hệ thống)](#1-actor)
   - 1.5. [Thuật ngữ & Khái niệm cốt lõi (Glossary)](#15-thuật-ngữ--khái-niệm-cốt-lõi-glossary)
2. [Mô tả phân cấp các ca sử dụng theo lộ trình phát hành (Phases)](#2-mô-tả-phân-cấp-các-ca-sử-dụng-theo-lộ-trình-phát-hành-phases)
3. [Chi tiết đặc tả các ca sử dụng cốt lõi (Phase 1, 2, 3)](#3-chi-tiết-đặc-tả-các-ca-sử-dụng-cốt-lõi)
4. [Cẩm nang gọi vốn dựa trên tài liệu BA kỹ thuật](#4-cẩm-nang-gọi-vốn-dựa-trên-tài-liệu-ba-kỹ-thuật)

---

## 1. Actor
- **Người dùng Độc thân (User):** Tác nhân chính sử dụng app để kết bạn, quẹt thẻ, gửi Locket, cày game.
- **Thành viên Clan (Clan Member):** Người dùng thuộc một bang hội cụ thể, tham gia đua top địa bàn.
- **Bồi thẩm đoàn (Jury User):** Người dùng ngẫu nhiên được hệ thống chọn để phân xử các vụ kiện tụng tình cảm.
- **Nhà đầu tư ảo (Trader):** Người dùng tham gia mua/bán cổ phiếu profile trong tính năng sàn chứng khoán.
- **Hệ thống AI (AI Engine):** Tác nhân tự động phân tích hội thoại, xử lý mờ ảnh và điều phối bot.
- **Quản trị viên (Admin):** Kiểm duyệt nội dung, xử lý tranh chấp cấp cao, quản lý hệ thống voucher.

---

## 1.5. Thuật ngữ & Khái niệm cốt lõi (Glossary)
- **Xu ảo (Token):** Đơn vị tiền tệ lưu hành nội bộ trong app (mua qua In-App Purchase). Xu **không thể** quy đổi ngược ra tiền mặt VNĐ (chỉ dùng đổi Voucher dịch vụ) để tuân thủ chính sách App Store/Google Play.
- **Streak:** Chuỗi ngày tương tác liên tục giữa 2 người dùng. Nếu bỏ lỡ 24h, Streak sẽ bị reset về 0.
- **Locket:** Tính năng chụp ảnh và đẩy trực tiếp lên màn hình chính (Widget) của đối phương.
- **Clan:** Bang hội do người dùng tự tạo để đua top, yêu cầu số lượng thành viên tối thiểu.

---

## 2. Mô tả phân cấp các ca sử dụng theo lộ trình phát hành (Phases)

### GIAI ĐOẠN 1: NỀN MÓNG VIRAL (PHASE 1 - MVP CORE)
- **UC-P1-001:** Đăng ký và Quẹt thẻ tìm đối tượng theo Hệ Tâm Linh.
- **UC-P1-002:** Gửi ảnh Blind Locket & Đồng bộ Widget màn hình khóa.
- **UC-P1-003:** Thành lập và Đua Top Bang hội (Chiến Trường Cupid).
- **UC-P1-004:** Đâm đơn kiện hành vi Ghosting (Tòa Án Tình Yêu).

### GIAI ĐOẠN 2: GIỮ CHÂN & DOANH THU (PHASE 2 - UPDATE 1.5)
- **UC-P2-005:** Chăm sóc và Khôi phục Đảo Tình Yêu 3D (Streak Gamification).
- **UC-P2-006:** Đánh chiếm địa bàn qua định vị GPS (Bản Đồ Chiếm Thành).
- **UC-P2-007:** Thiết lập Khế Ước Tài Chính chống bùng kèo (Cọc tiền đi Date).
- **UC-P2-008:** Trò chuyện cùng Trợ lý Cánh Gió "Mỏ Hỗn" (AI Wingman).

### GIAI ĐOẠN 3: THAO TÚNG TỐI THƯỢNG (PHASE 3 - UPDATE 2.0)
- **UC-P3-009:** Đánh giá và Tra cứu CV Tình Trường (Ex-Rating).
- **UC-P3-010:** Mua bán cổ phần Profile (Sàn Chứng Khoán Độc Thân).
- **UC-P3-011:** Sắp xếp ghép đôi hỗn loạn hàng tuần (Đêm Săn Mồi "The Purge").

---

## 3. Chi tiết đặc tả các ca sử dụng cốt lõi
*(Lưu ý: Tài liệu Draft hiện tại đang tập trung đặc tả 5 Use-case phức tạp và quan trọng nhất. Các Use-case còn lại sẽ được bổ sung chi tiết trong các phiên bản cập nhật tiếp theo).*

### NHÓM 1: PHASE 1 - NỀN MÓNG VIRAL & AN TOÀN STORE

#### 1. UC-P1-002: Gửi ảnh Blind Locket & Đồng bộ Widget
- **Mô tả ngắn:** Sử dụng thuật toán che mờ ảnh real-time gửi qua Widget màn hình chính để kích thích sự tò mò và duy trì chuỗi tương tác (Streak) của cặp đôi.
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Người dùng chụp ảnh từ Widget Q-Love ngoài màn hình chính. | 2. Hệ thống tiếp nhận, kiểm tra điểm thân thiết (Streak Score) hiện tại của cặp đôi trong database. | - ID Người gửi<br>- ID Người nhận<br>- File ảnh thô |
| | 3. Nếu Streak Score < 10, hệ thống sử dụng thư viện đồ họa để băm nát ảnh bằng bộ lọc Gaussian Blur ở mức 90%. Nếu Streak Score >= 10, giảm mức độ mờ tương ứng (VD: 50% hoặc 0%). Lưu file ảnh đã xử lý vào Cloud Storage (R2/MinIO). | - URL ảnh mờ (blur_image_url) |
| | 4. Hệ thống đẩy Payload thông báo qua FCM/APNs để Widget trên màn hình người nhận tự động render ảnh mờ. | - Payload thông báo |
| 5. Người nhận bấm vào Widget để vào app, nhập tin nhắn tương tác. | 6. Hệ thống ghi nhận tin nhắn, cộng điểm Streak và tự động giảm tỷ lệ mờ (Blur) của các ảnh tiếp theo. | - Nội dung chat<br>- Mức độ mờ mới (Blur level) |

- **Luồng ngoại lệ:** 
  - *Lỗi kết nối mạng khi tải ảnh:* Hệ thống lưu tạm ảnh vào Local Cache của thiết bị và tự động Retry khi có mạng lại.
  - *Hệ thống phát hiện ảnh nhạy cảm:* AI kiểm duyệt hình ảnh quét qua ảnh gốc, nếu dính tỷ lệ khỏa thân > 30% sẽ tự động hủy lệnh gửi và cảnh báo khóa tài khoản người gửi.
- **Yêu cầu đặc biệt:** Tốc độ đẩy ảnh từ Widget người gửi lên Cloud Storage và sync sang Widget người nhận phải dưới 3 giây. Ảnh lưu trên Widget bắt buộc phải qua server xử lý mờ, không xử lý ở Client để tránh việc User bypass code dịch ngược.
- **Tiền điều kiện:** Cả hai user đã match nhau, đã cấp quyền định vị background và quyền truy cập Widget trên iOS/Android.
- **Điều kiện sau:** Ảnh mờ được hiển thị thành công ngoài màn hình khóa, điểm Streak được kích hoạt chu kỳ 24 giờ.
- **Điểm mở rộng:** Kết nối với Apple Watch để đồng bộ rung phản hồi (Haptic feedback) theo nhịp tim khi hai người cùng mở Widget một lúc.

#### 2. UC-P1-004: Đâm đơn kiện hành vi Ghosting (Tòa Án Tình Yêu)
- **Mô tả ngắn:** Use-case cho phép người dùng trừng phạt những kẻ im lặng bùng kèo bằng cách đưa ra xét xử ẩn danh trước cộng đồng, tạo cơ chế "game hóa" luật pháp tình cảm nhằm giữ uy tín cho app.
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Người dùng chọn nút "Đâm Đơn Kiện" trong hộp thoại chat bị bỏ rơi. | 2. Hệ thống quét database để xác thực điều kiện (im lặng > 48 tiếng, và `highest_streak_score` > 5 ngày). | - ID Phòng chat<br>- Kiểm tra bảng `matches` |
| 3. Người dùng chọn lý do kiện (Ví dụ: Trapboy/Trapgirl, Đột ngột biến mất, Nói lời cay đắng) và nhấn xác nhận. | 4. Hệ thống ẩn danh thông tin cá nhân, chụp 5 block chat cuối cùng và tạo một bản ghi vụ kiện mới. | - Lý do kiện*<br>- Block chat ẩn danh |
| | 5. Hệ thống phân phối vụ kiện này vào hàng đợi (Redis Queue) của mục "Hóng Drama" thuộc 50 người dùng ngẫu nhiên thuộc các Clan khác. | - ID Vụ kiện (case_id) |
| 6. Các Bồi thẩm đoàn vào đọc tình huống và bấm vote "Có tội" hoặc "Vô tội". | 7. Sau 12 tiếng, hệ thống chốt số lượng vote. Nếu tỷ lệ "Có tội" > 65%, hệ thống thực thi lệnh trừng phạt lên tài khoản bị kiện. | - Lượt vote (bool) |

- **Luồng ngoại lệ:** 
  - *Hòa giải ngoài tòa:* Trong thời gian chờ xét xử, nếu người bị kiện tặng quà tạ lỗi (bằng xu ảo) và người kiện chấp nhận, hệ thống tự động đình chỉ vụ kiện.
  - *Không đủ số lượng bồi thẩm đoàn:* Nếu sau 12 tiếng không đủ 50 người vote, hệ thống tự động đẩy vụ kiện lên tab "Hot Drama" và thưởng thêm xu ảo cho ai vào phân xử để kích thích số lượng vote. Thời gian vote được gia hạn thêm 12 tiếng.
- **Yêu cầu đặc biệt:** Hệ thống phải tự động che mờ toàn bộ Avatar, Tên thật, Số điện thoại xuất hiện trong đoạn chat bằng regex để bảo mật thông tin cá nhân, tránh vi phạm chính sách bêu rếu của Apple App Store.
- **Tiền điều kiện:** User bị ghost thỏa mãn điều kiện thời gian im lặng hệ thống quy định (> 48h).
- **Điều kiện sau:** Kẻ tội phạm bị gắn huy hiệu "Ghost thủ" trên profile công khai, bị bóp thuật toán hiển thị quẹt thẻ (giảm 80% tần suất xuất hiện) trong 3 ngày kế tiếp.
- **Điểm mở rộng:** Cho phép tội phạm nạp xu để mua "Giấy giảm án" xóa huy hiệu ngay lập tức (Kênh tạo doanh thu trực tiếp).

### NHÓM 2: PHASE 2 - GIỮ CHÂN & DOANH THU (O2O & GAMIFICATION)

#### 3. UC-P2-006: Đánh chiếm địa bàn bằng Locket (Bản Đồ Chiếm Thành)
- **Mô tả ngắn:** Biến bản đồ thực tế thành một chiến trường. Các Clan hoặc cặp đôi sử dụng ảnh chụp Locket có đính kèm GPS định vị tại các địa điểm hot để tranh giành quyền "Cắm cờ chủ quyền".
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Thành viên Clan đến địa điểm thực tế (Ví dụ: Phố đi bộ), mở map Q-Love và bấm "Check-in Chiếm Thành". | 2. Hệ thống kích hoạt camera, bắt buộc chụp ảnh trực tiếp (không cho chọn từ thư viện) kèm tọa độ GPS real-time. | - Tọa độ GPS thực tế<br>- File ảnh thực tế |
| 3. Người dùng chọn "Đóng góp điểm cho Clan" và nhấn đăng bài. | 4. Hệ thống đối chiếu tọa độ với bán kính của Landmark trong DB. Nếu khớp, cộng điểm tích lũy của địa điểm đó cho Clan của User. | - ID Clan<br>- ID Landmark |
| | 5. Hệ thống tính tổng điểm (Mỗi lượt Check-in hợp lệ = 10 điểm). Nếu Clan của User vượt Clan cũ, hệ thống cập nhật Flag trên Map. Bảng xếp hạng reset vào 0h00 Thứ Hai hàng tuần. | - Bảng xếp hạng Landmark Địa bàn |

- **Luồng ngoại lệ:** 
  - *Fake GPS (Giả lập vị trí):* Hệ thống tích hợp bộ quét thư viện Mock Location trên Android/iOS. Nếu phát hiện user dùng phần mềm đổi tọa độ, hệ thống báo lỗi đỏ, hủy bài đăng và cấm check-in trong 7 ngày.
- **Yêu cầu đặc biệt:** Bản đồ sử dụng Mapbox API hoặc Google Maps SDK tối ưu hóa bộ nhớ đệm (Caching) để khi render hàng nghìn cờ chủ quyền không gây giật lag hoặc nóng máy thiết bị.
- **Tiền điều kiện:** Người dùng bật định vị ở chế độ "High Accuracy" và là thành viên của một Clan hợp lệ.
- **Điều kiện sau:** Avatar của Clan được hiển thị đè lên bản đồ vệ tinh của địa điểm đó, toàn bộ user đi qua khu vực đó đều nhìn thấy.
- **Điểm mở rộng:** Không có.

#### 4. UC-P2-007: Thiết lập Khế Ước Tài Chính chống bùng kèo (Cọc tiền đi Date)
- **Mô tả ngắn:** Cho phép hai người dùng khóa một lượng xu ảo (token) làm tiền bảo chứng cho cuộc hẹn ngoài đời. Nếu ai bùng kèo, tiền cọc sẽ chuyển sang ví đối phương.
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Tại giao diện chat, User A chọn "Lên lịch hẹn hò" và nhập số xu muốn cọc (Ví dụ: 100 xu). | 2. Hệ thống kiểm tra số dư ví ảo của User A. Nếu đủ, thực hiện lệnh đóng băng (Hold) số xu đó. | - Số xu cọc*<br>- Thời gian hẹn*<br>- Địa điểm hẹn (Landmark ID)* |
| 3. User B nhận được lời mời, bấm "Xác nhận khế ước" và hệ thống tự động khóa số xu tương ứng trong ví User B. | 4. Hệ thống cập nhật trạng thái bản ghi `dating_contracts` sang `active` và chuyển `balance → hold_balance` trong bảng `user_wallets` bằng một DB Transaction (Serializable Isolation) duy nhất. | - Trạng thái: `active` (Đang chờ) |
| 5. Đến giờ hẹn, hai người gặp nhau, một người bật mã Dynamic QR Code (sinh từ thuật toán `TOTP` dưới Client, đổi 30s/lần), người kia quét mã. | 6. Hệ thống xác thực mã TOTP thành công. Giải phóng lệnh đóng băng, hoàn lại xu cho cả hai. | - Mã TOTP Code<br>- Tọa độ GPS (Optional) |

- **Luồng ngoại lệ:** 
  - *Một bên không đến (Bùng kèo):* Quá giờ hẹn 30 phút mà không quét mã QR, hệ thống kiểm tra: Nếu User vắng mặt có gói **Q-Love Premium** và chưa dùng quyền miễn trừ trong tháng -> Hủy lịch hẹn, không trừ cọc, hoàn xu cho cả hai. Nếu không có hoặc đã hết quyền -> Tịch thu 100 xu chuyển thẳng vào ví người đến đúng giờ (trừ 10% phí vận hành).
- **Yêu cầu đặc biệt:** Việc nạp tiền mua xu phải qua In-App Purchase (IAP) của Apple/Google. Việc giải ngân đền bù chỉ diễn ra bằng xu nội bộ dùng để đổi Voucher (Highlands, CGV), không cho phép rút thành tiền mặt trực tiếp để tránh bị App Store quét lỗi "Ứng dụng cờ bạc, cá cược".
- **Tiền điều kiện:** Hai bên đã chat tối thiểu 20 câu và số dư tài khoản đủ để cọc.
- **Điều kiện sau:** Lịch hẹn được lưu vào lịch hệ thống, ví của cả hai bị trừ tạm thời số tiền bảo chứng.
- **Điểm mở rộng:** Liên kết trực tiếp API với các chuỗi quán cafe lớn để tự động đặt bàn và thanh toán hóa đơn bằng chính số xu cọc.

### NHÓM 3: PHASE 3 - THAO TÚNG TỐI THƯỢNG & DOANH THU ĐỘT PHÁ

#### 5. UC-P3-010: Mua bán cổ phần Profile (Sàn Chứng Khoán Độc Thân)
- **Mô tả ngắn:** Biến các profile cá nhân thành các mã cổ phiếu tài chính ảo. Giá trị mã tăng/giảm dựa trên các chỉ số hoạt động, tương tác và mức độ săn đón của cộng đồng đối với profile đó.
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Người dùng (Trader) vào sàn, tìm kiếm mã cổ phiếu của một Profile (Ví dụ: Mã $NVA của Nguyễn Văn A). | 2. Hệ thống hiển thị biểu đồ hình nến giá trị của mã $NVA cùng tập chỉ số: Lượt match, Điểm Clan đóng góp, Tổng số chuỗi Locket đang giữ. | - Mã cổ phiếu (stock_code) |
| 3. Người dùng nhập số lượng cổ phần muốn mua và nhấn "Đặt lệnh mua". | 4. Hệ thống kiểm tra thanh khoản và số lượng "Cổ phiếu tự do lưu hành" của profile đó. Thực hiện khớp lệnh qua cơ chế Order-matching engine ảo. | - Số lượng mua*<br>- Giá đặt mua* |
| | 5. Cập nhật bảng danh sách cổ đông của $NVA. Khi giá trị $NVA tăng, Trader có thể đặt lệnh "Bán" để ăn chênh lệch xu. | - Sổ lệnh (Order Book) |

- **Luồng ngoại lệ:** 
  - *Mã cổ phiếu "Hủy niêm yết" (User hủy tài khoản hoặc kết đôi thành công thoát ế):* Hệ thống kích hoạt trạng thái "Mừng Đám Cưới", toàn bộ người nắm giữ cổ phần mã này vào thời điểm hủy niêm yết sẽ được hệ thống mua lại với mức giá trần cố định + thưởng thêm Quà tặng từ quỹ của App làm tiền mừng.
- **Yêu cầu đặc biệt:** Thuật toán tính giá cổ phiếu phải cập nhật liên tục 5 phút/lần dựa trên công thức: `Giá = 100 Xu (Giá sàn IPO) + (Số lượt Match mới * 0.4) + (Lượt gửi Locket * 0.3) + (Số thành viên Clan upvote * 0.3) - (Số đơn kiện từ Tòa Án * 0.5)`. Cần có giới hạn biên độ tăng/giảm (trần/sàn) trong ngày để chống bơm thổi (Pump & Dump).
- **Tiền điều kiện:** Người dùng đạt Level 5 trong app mới được mở tính năng Sàn giao dịch.
- **Điều kiện sau:** Xu được luân chuyển giữa các ví người dùng, tạo dòng tiền lưu thông cực lớn trong hệ sinh thái app.
- **Điểm mở rộng:** Không có.

### NHÓM 4: CÁC USE-CASE BỔ SUNG (PHASE 1, 2, 3)

#### 6. UC-P1-001: Đăng ký và Quẹt thẻ tìm đối tượng theo Hệ Tâm Linh
- **Mô tả ngắn:** Cho phép người dùng mới đăng ký bằng số điện thoại (OTP), thiết lập hồ sơ cá nhân và bắt đầu quẹt thẻ tìm đối tượng với bộ lọc độc đáo dựa trên Hệ Tâm Linh (Cung hoàng đạo, Thần số học).
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Người dùng nhập số điện thoại và nhấn "Gửi mã OTP". | 2. Hệ thống gọi API của **ESMS.vn** để gửi mã OTP 6 chữ số về SĐT. Mã có hiệu lực trong 120 giây. | - Số điện thoại<br>- OTP Token (lưu Redis, TTL=120s) |
| 3. Người dùng nhập mã OTP. | 4. Hệ thống xác thực OTP với Redis. Nếu khớp, tạo bản ghi `users` mới và cấp JWT Access Token (15 phút) + Refresh Token (30 ngày). | - JWT Access Token<br>- Refresh Token |
| 5. Người dùng điền hồ sơ: Tên, ngày sinh, giới tính, ảnh đại diện, cung hoàng đạo/thần số học ưa thích. | 6. Hệ thống lưu profile, tính Hệ Tâm Linh tương thích, kích hoạt `user_wallets` với số dư 0. | - Profile data<br>- `wallet` khởi tạo |
| 7. Người dùng vào màn hình Quẹt thẻ, chọn bộ lọc "Hệ Tâm Linh". | 8. Thuật toán chỉ trả về profile có độ tương thích tâm linh > 70% trong bán kính 50km. Người dùng vuốt Phải (Like) hoặc Trái (Pass). | - Danh sách profile filtered |
| 9. Cả hai người đều vuốt Phải nhau. | 10. Hệ thống tạo bản ghi `matches`, hiển thị popup "It's a Match!" và mở khóa Chatroom. | - `match_id`<br>- Chatroom mở |

- **Luồng ngoại lệ:**
  - *OTP hết hạn hoặc sai:* Hệ thống báo lỗi, cho phép gửi lại OTP sau 60 giây (tối đa 5 lần/ngày/SĐT).
  - *Số điện thoại đã tồn tại:* Hệ thống chuyển sang luồng Đăng nhập thay vì Đăng ký.
- **Yêu cầu đặc biệt:** JWT Refresh Token phải được lưu trữ an toàn (Secure Storage trên thiết bị), không lưu trong LocalStorage. Backend phải rotate Refresh Token mỗi lần sử dụng (One-time use).
- **Tiền điều kiện:** Không có (luồng mở đầu của hệ thống).
- **Điều kiện sau:** Tài khoản và ví Xu được kích hoạt, người dùng có thể bắt đầu quẹt thẻ.
- **Điểm mở rộng:** Đăng nhập bằng Google/Apple ID ở Phase 2.

---

#### 7. UC-P1-003: Thành lập và Quản lý Bang hội (Clan)
- **Mô tả ngắn:** Cho phép người dùng đủ điều kiện tạo Clan mới, mời thành viên và tham gia vào cuộc chiến đua điểm địa bàn trên bản đồ.
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Người dùng (Level ≥ 3, Số dư ≥ 500 Xu) vào mục "Bang hội" và nhấn "Tạo Bang hội mới". | 2. Hệ thống kiểm tra điều kiện Level và số dư ví. Nếu thỏa mãn, hiển thị form nhập thông tin. | - Level user<br>- Số dư ví |
| 3. Người dùng nhập Tên bang, slogan, upload logo và nhấn "Xác nhận tạo". | 4. Hệ thống trừ 500 Xu, tạo bản ghi `clans` mới, thêm người tạo vào `clan_members` với `role = leader`. | - Tên bang (unique)<br>- Logo URL<br>- `leader_id` |
| 5. Bang chủ vào tab "Quản lý thành viên" và nhấn "Mời bạn bè". | 6. Hệ thống tìm kiếm user theo tên/SĐT và gửi lời mời tham gia vào In-App Notification của người được mời. | - Danh sách user được mời |
| 7. Người được mời nhận thông báo và nhấn "Chấp nhận". | 8. Hệ thống thêm bản ghi vào `clan_members` với `role = member`. Tự động cập nhật số thành viên hiển thị. | - `clan_id`, `user_id` |

- **Luồng ngoại lệ:**
  - *Tên bang đã tồn tại:* Hệ thống báo lỗi `Tên bang hội đã có người dùng. Vui lòng chọn tên khác.`
  - *Không đủ điều kiện:* Nút "Tạo Bang hội" hiển thị greyed-out kèm tooltip giải thích điều kiện còn thiếu.
  - *Bang chủ rời bang:* Hệ thống yêu cầu chỉ định Bang chủ mới trước khi rời. Nếu là thành viên cuối, Clan tự động giải tán.
- **Yêu cầu đặc biệt:** Tên Clan phải qua bộ lọc từ ngữ thô tục (Profanity Filter) trước khi lưu vào DB.
- **Tiền điều kiện:** Người dùng đạt Level 3 và có đủ 500 Xu trong ví.
- **Điều kiện sau:** Clan được tạo thành công, người tạo trở thành Bang chủ, 500 Xu bị trừ.
- **Điểm mở rộng:** Cho phép doanh nghiệp (Quán café, nhà hàng) tạo "Clan địa điểm" được xác nhận (Verified) để thu hút người dùng check-in.

---

#### 8. UC-P2-005: Chăm sóc và Khôi phục Đảo Tình Yêu 3D (Streak Gamification)
- **Mô tả ngắn:** Mỗi cặp đôi Match có một "Đảo Tình Yêu" 3D chung. Đảo nâng cấp khi Streak dài, héo úa khi đứt Streak, và có thể khôi phục bằng Xu ảo — tạo động lực gắn kết liên tục.
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Người dùng A và B tương tác (nhắn tin/gửi Locket) trong ngày. | 2. Hệ thống ghi nhận `last_interaction_at`, tăng `streak_score` thêm 1 điểm. | - `match_id`<br>- `streak_score` mới |
| | 3. Cron Job chạy lúc 00:00 mỗi ngày kiểm tra `streak_score` đạt milestone. Nếu đạt 7 ngày, 30 ngày — tự động nâng cấp `island_level` của match đó. | - `island_level` tăng<br>- Push notification "Đảo lên level!" |
| 4. Người dùng A và B không tương tác quá 24h. | 5. Cron Job phát hiện `last_interaction_at` > 24h, reset `streak_score = 0`, đánh dấu `island_level` sang trạng thái "Héo úa" (chuyển màu/hiệu ứng). | - Push notification "Đảo của bạn đang héo 🥀" |
| 6. Người dùng A vào app, thấy Đảo héo, nhấn "Hồi sinh đảo" và thanh toán 50 Xu. | 7. Hệ thống trừ 50 Xu, khôi phục `island_level` về cấp trước đó (không reset streak, chỉ khôi phục đảo). | - 50 Xu bị trừ<br>- Đảo hồi sinh |

- **Luồng ngoại lệ:**
  - *Không đủ Xu để hồi sinh:* Hệ thống hiển thị cửa hàng nạp Xu (In-App Purchase).
  - *Cả hai đều offline > 7 ngày:* Đảo Tình Yêu biến mất hoàn toàn. Người dùng cần nộp 200 Xu để tái thiết từ đầu (Level 1).
- **Yêu cầu đặc biệt:** Hiệu ứng 3D của Đảo phải render mượt ≥ 60fps bằng Flutter 3D (hoặc Rive animation) mà không gây nóng máy.
- **Tiền điều kiện:** Hai người dùng đã Match nhau.
- **Điều kiện sau:** `island_level` trong bảng `matches` được cập nhật, hiệu ứng hiển thị đúng trạng thái.
- **Điểm mở rộng:** Mở khóa ảnh/filter Locket độc quyền khi đảo đạt Level cao.

---

#### 9. UC-P2-008: Trò chuyện cùng Trợ lý Cánh Gió "Mỏ Hỗn" (AI Wingman)
- **Mô tả ngắn:** Khi cuộc trò chuyện bị tắt ngấm hoặc người dùng không biết nói gì, AI Wingman đọc lịch sử chat gần nhất và gợi ý 3 tin nhắn phản hồi với tone khác nhau (Hài hước, Thả thính, Thẳng thắn).
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Người dùng A đang trong chatroom với B. Chat bị ngừng > 30 phút hoặc A xóa tin nhắn 3 lần liên tiếp. | 2. Hệ thống hiển thị icon trợ lý AI "Mỏ Hỗn" nhấp nháy gợi ý ở góc dưới màn hình chat. | - Trigger: `last_message_at` > 30 phút |
| 3. Người dùng A bấm vào icon Trợ lý AI. | 4. Hệ thống lấy **5 tin nhắn gần nhất** từ `chat_messages`, ẩn danh hóa tên và gọi **OpenAI/Claude API** với prompt được thiết kế sẵn để sinh 3 gợi ý. | - 5 tin nhắn gần nhất (raw text) |
| | 5. Hệ thống trả về giao diện popup 3 gợi ý: `[💬 Hài hước] [🎯 Thả thính] [🤝 Thẳng thắn]`. | - 3 suggested messages |
| 6. Người dùng A chọn 1 trong 3 gợi ý hoặc bấm "Thử lại". | 7. Tin nhắn được đẩy vào ô soạn thảo để A có thể chỉnh sửa trước khi gửi (AI không tự gửi). | - Selected suggestion |

- **Luồng ngoại lệ:**
  - *API OpenAI/Claude timeout (> 5 giây):* Hệ thống hiển thị thông báo lỗi "Trợ lý đang bận, thử lại sau" và không tính phí lượt sử dụng.
  - *Chat chưa đủ lịch sử (< 5 tin nhắn):* AI Wingman bị vô hiệu hóa, hiển thị tooltip "Chat thêm một chút để Trợ lý hiểu bạn hơn nhé!".
- **Yêu cầu đặc biệt:** Prompt gửi sang OpenAI phải được ẩn danh hoàn toàn (không gửi tên thật, SĐT). Toàn bộ dữ liệu chat phải xử lý trong RAM của Backend, không log ra file.
- **Tiền điều kiện:** Hai người đã Match và có ít nhất 5 tin nhắn trong chatroom.
- **Điều kiện sau:** Gợi ý AI được hiển thị để người dùng lựa chọn.
- **Điểm mở rộng:** Tính năng "Phân tích nhân cách" — AI đọc 50 tin nhắn gần nhất và đưa ra nhận xét "Đối phương có vẻ đang thích bạn/chán bạn".

---

#### 10. UC-P3-009: Đánh giá và Tra cứu CV Tình Trường (Ex-Rating)
- **Mô tả ngắn:** Sau khi Unmatch, người dùng được đánh giá nhau ẩn danh bằng sao và tag. Bất kỳ ai cũng có thể tra cứu "CV Tình Trường" của một profile bằng cách trả phí Xu để xem điểm số và tag tổng hợp từ những người yêu/match cũ.
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Người dùng A bấm "Unmatch" với người dùng B (đã chat > 50 tin hoặc đã gặp mặt O2O). | 2. Hệ thống hiển thị màn hình "Đánh giá tình trường" bắt buộc trước khi Unmatch hoàn tất. | - `match_id`<br>- Điều kiện: > 50 tin nhắn |
| 3. Người dùng A chấm 1-5 sao và chọn các tag mô tả (VD: `#Lịch_sự`, `#RedFlag`, `#Ghoster`, `#Vui_vẻ`). | 4. Hệ thống lưu bản ghi vào `ex_ratings` với `reviewer_id = A (ẩn danh)`, `target_id = B`. Unmatch được hoàn tất. | - `rating_score`, `tags[]` |
| 5. Người dùng C vào xem profile của B và bấm "Tra cứu CV Tình Trường". | 6. Hệ thống kiểm tra số dư ví của C. Nếu đủ 50 Xu, trừ 50 Xu và hiển thị: Điểm trung bình, tần suất xuất hiện của từng tag — **không hiển thị danh tính người đánh giá**. | - 50 Xu trừ từ ví C<br>- Dữ liệu aggregate (avg score + tag cloud) |

- **Luồng ngoại lệ:**
  - *Người dùng B report rằng bị đánh giá oan:* B có thể nộp đơn kháng cáo (appeal). Admin xem xét và có quyền xóa review vi phạm.
  - *C không đủ 50 Xu:* Hệ thống hiển thị màn hình "Nạp Xu" thay vì dữ liệu CV.
- **Yêu cầu đặc biệt:** Danh tính người đánh giá phải được ẩn danh tuyệt đối (không suy ra được từ thông tin công khai). Hệ thống phải có cơ chế chống đánh giá trả thù (Revenge Rating) — AI phân tích nội dung tag, nếu phát hiện mẫu nhất quán từ một người dùng thì đánh dấu để Admin xét duyệt.
- **Tiền điều kiện:** Hai người dùng đã Match và đủ điều kiện (chat > 50 tin hoặc đã gặp mặt O2O).
- **Điều kiện sau:** Điểm rating được lưu, góp vào điểm trung bình tổng hợp trên CV Tình Trường của B.
- **Điểm mở rộng:** Cho phép bình chọn "Hữu ích" (Like) lên các tag anonymous để tăng độ tin cậy.

---

#### 11. UC-P3-011: Sắp xếp ghép đôi hỗn loạn hàng tuần (Đêm Săn Mồi - The Purge)
- **Mô tả ngắn:** Sự kiện có thời hạn (tối Thứ 6, 22h00-00h00) ghép đôi người dùng hoàn toàn ẩn danh, không giới hạn bán kính. Đây là cơ chế tạo nội dung viral định kỳ và kéo user quay lại mỗi tuần.
- **Luồng cơ bản:**

| Hành động của tác nhân | Phản ứng của hệ thống | Dữ liệu |
| :--- | :--- | :--- |
| 1. Đến 22h00 tối Thứ 6, hệ thống chạy Cron Job kích hoạt sự kiện "The Purge". | 2. Hệ thống hiển thị banner đếm ngược trên màn hình chính của tất cả user đang online. | - Thời gian còn lại |
| 3. Người dùng bấm "Tham gia ngay". | 4. Hệ thống thêm user vào hàng đợi (Redis Queue). Sau đó, **Matchmaking Worker** chạy độc lập sẽ xử lý thuật toán ghép đôi, giảm tải cho Server chính. Chỉ hiển thị "Đối tượng ẩn danh 🎭". | - Hàng đợi ghép đôi (Redis) |
| 5. Hai người bắt đầu chat ẩn danh trong 10 phút. | 6. Sau 10 phút, hệ thống hiển thị nút "Lộ diện" cho cả 2. Chỉ khi **cả 2 cùng bấm "Lộ diện"** trong 30 giây, Avatar và Profile thật mới được hiển thị. | - Timer 10 phút<br>- Consent của cả 2 bên |
| 7. Nếu cả 2 cùng bấm "Lộ diện". | 8. Hệ thống reveal Profile và hỏi "Bạn có muốn Match chính thức không?". Nếu cả 2 đồng ý, tạo `matches` record như bình thường. | - Match mới (optional) |

- **Luồng ngoại lệ:**
  - *Một bên không bấm "Lộ diện":* Profile ẩn danh vĩnh viễn. Cuộc trò chuyện kết thúc sau 30 giây chờ, không lưu lịch sử.
  - *Nội dung xúc phạm trong chat ẩn danh:* AI quét realtime, nếu phát hiện vi phạm thì mute tài khoản khỏi The Purge trong 4 tuần.
  - *Không đủ người tham gia (< 10 người):* Hệ thống mở rộng bán kính ghép đôi ra toàn quốc.
- **Yêu cầu đặc biệt:** Trong thời gian The Purge, WebSocket Hub phải chịu tải đột biến lớn. Hệ thống cần Auto-scaling (K8s HPA) được kích hoạt trước sự kiện 15 phút.
- **Tiền điều kiện:** Tài khoản đã xác thực, không đang bị ban.
- **Điều kiện sau:** Cuộc trò chuyện kết thúc sạch (không lưu log ẩn danh), hoặc Match chính thức được tạo nếu cả 2 đồng ý lộ diện.
- **Điểm mở rộng:** Phiên bản "The Purge VIP" cho Premium user — ghép đôi theo Hệ Tâm Linh thay vì hoàn toàn ngẫu nhiên.

---

## 4. Cẩm nang gọi vốn (BRD)
Chi tiết về Mô hình Kinh Doanh (Monetization), Mục tiêu (Objectives) và Quản trị rủi ro đã được đóng gói thành Cẩm nang gọi vốn. Vui lòng xem tại: [Business Requirements Document (brd.md)](./brd.md)
