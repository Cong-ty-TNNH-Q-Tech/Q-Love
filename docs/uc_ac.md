# Use Case & Acceptance Criteria (UC-AC) Document
**Project Name:** Q-Love - Super App Hẹn Hò & Mạng Xã Hội Giải Trí  
**Version:** 1.0  
**Date:** June 2026  
**Language:** Vietnamese (Tiếng Việt)  

Tài liệu này đặc tả các Tiêu chí chấp nhận (Acceptance Criteria - AC) cho 5 Use-case cốt lõi nhất của hệ thống, phục vụ cho đội ngũ QA/Tester viết Test Case và đội ngũ Dev bám sát yêu cầu nghiệm thu. Cấu trúc được viết theo chuẩn **Given - When - Then** (BDD).

---

## 1. UC-P1-002: Gửi ảnh Blind Locket & Đồng bộ Widget

### AC1: Gửi ảnh thành công khi Streak Score < 10 (Chế độ làm mờ)
- **Given** Người dùng A và B đã Match và có Streak Score = 5 (nhỏ hơn 10). Cả 2 đã cài Widget Q-Love ra màn hình chính.
- **When** Người dùng A chụp một bức ảnh an toàn (không chứa nội dung nhạy cảm) và nhấn "Gửi qua Locket".
- **Then** Hệ thống băm mờ bức ảnh bằng filter Gaussian Blur 90% trước khi lưu trữ.
- **And** Trong vòng 3 giây, Widget trên màn hình của người dùng B tự động cập nhật hiển thị bức ảnh đã làm mờ đó kèm notification.

### AC2: Gửi ảnh thành công khi Streak Score >= 10 (Chế độ rõ nét)
- **Given** Người dùng A và B đã Match và có Streak Score = 12 (lớn hơn hoặc bằng 10).
- **When** Người dùng A chụp ảnh và gửi qua Locket.
- **Then** Hệ thống giảm mức độ mờ xuống 0% (hoặc mức tương ứng theo cấu hình) và đẩy ảnh lên Widget của B trong dưới 3 giây.

### AC3: Hệ thống phát hiện ảnh nhạy cảm (NSFW)
- **Given** Người dùng A đang cố gắng gửi một bức ảnh có tỷ lệ da thịt/khỏa thân > 30%.
- **When** Người dùng A nhấn "Gửi qua Locket".
- **Then** AI kiểm duyệt của hệ thống quét và chặn ngay lệnh gửi.
- **And** Hệ thống hiển thị popup cảnh báo lỗi cho người dùng A và ghi log cảnh cáo vào database (Nếu vi phạm 3 lần sẽ auto-ban tài khoản).

---

## 2. UC-P1-004: Đâm đơn kiện hành vi Ghosting (Tòa Án Tình Yêu)

### AC1: Khởi kiện thành công (Thỏa mãn điều kiện)
- **Given** Người dùng B đã im lặng (không reply tin nhắn) quá 48 tiếng, và Streak cũ của cả hai > 5 ngày.
- **When** Người dùng A mở đoạn chat và chọn "Đâm đơn kiện" -> Chọn lý do "Trapboy/Trapgirl".
- **Then** Hệ thống tạo thành công vụ kiện (Case ID).
- **And** Hệ thống tự động che mờ Tên, Avatar, SĐT trong 5 block chat gần nhất bằng regex và đẩy vụ kiện vào Redis Queue chờ Bồi thẩm đoàn.

### AC2: Trừng phạt khi Bồi thẩm đoàn kết án "Có tội"
- **Given** Vụ kiện đã kết thúc thời gian 12 tiếng.
- **When** Hệ thống đếm tổng số vote đạt trên 50 người, trong đó tỷ lệ vote "Có tội" > 65%.
- **Then** Hệ thống tự động gắn huy hiệu "Ghost thủ" lên profile của người dùng B.
- **And** Thuật toán hiển thị của người dùng B bị giảm 80% (Shadowban) trong 3 ngày kế tiếp.

### AC3: Hòa giải ngoài tòa thành công
- **Given** Vụ kiện đang trong thời gian 12 tiếng chờ vote.
- **When** Người dùng B gửi một Vật phẩm "Quà tạ lỗi" (được mua bằng Xu ảo) cho người dùng A và A bấm "Chấp nhận".
- **Then** Hệ thống tự động đình chỉ vụ kiện, gỡ bài khỏi Tòa Án. Bồi thẩm đoàn sẽ nhận được thông báo "Vụ kiện đã được hòa giải".

---

## 3. UC-P2-006: Đánh chiếm địa bàn qua định vị GPS

### AC1: Check-in đóng góp điểm thành công
- **Given** Người dùng A thuộc Clan "Độc Cô Cầu Bại" đang đứng tại "Phố đi bộ Nguyễn Huệ".
- **When** A mở app, chụp ảnh bằng camera trực tiếp và bấm "Check-in Chiếm Thành".
- **Then** Hệ thống đối chiếu GPS hợp lệ (trong bán kính Landmark).
- **And** Hệ thống cộng thêm 10 điểm vào quỹ điểm tuần của Clan "Độc Cô Cầu Bại" tại Landmark này.

### AC2: Cập nhật Cờ chủ quyền trên bản đồ
- **Given** Tổng điểm tuần của Clan "Độc Cô Cầu Bại" vừa vượt qua Clan "Hội F.A" (đang giữ cờ).
- **When** Hệ thống chạy lệnh trigger tính điểm.
- **Then** Avatar/Cờ của Clan "Độc Cô Cầu Bại" ngay lập tức đè lên Landmark "Phố đi bộ Nguyễn Huệ" trên bản đồ của toàn bộ user.

### AC3: Phát hiện gian lận Fake GPS
- **Given** Người dùng A sử dụng phần mềm Mock Location (Fake GPS) để giả lập vị trí tại Phố đi bộ.
- **When** A bấm "Check-in Chiếm Thành".
- **Then** Bộ quét của hệ thống phát hiện Mock Location, trả về báo lỗi đỏ "Phát hiện giả lập vị trí".
- **And** Nút Check-in của user bị vô hiệu hóa trong 7 ngày.

---

## 4. UC-P2-007: Khế Ước Tài Chính (Cọc tiền đi Date)

### AC1: Thiết lập Khế ước thành công
- **Given** User A và User B có số dư Ví ảo >= 100 Xu. Đã chat tối thiểu 20 câu.
- **When** User A gửi lời mời cọc 100 Xu, User B bấm "Xác nhận khế ước".
- **Then** Trạng thái cuộc hẹn chuyển thành "Đang chờ". 100 Xu trong ví của mỗi người bị chuyển sang trạng thái "Đóng băng" (Hold).

### AC2: Gặp mặt thành công (Quét Dynamic QR TOTP)
- **Given** Khế ước đang ở trạng thái "Đang chờ".
- **When** Đến giờ hẹn, User A quét mã Dynamic QR (sinh bằng thuật toán TOTP dưới Client, đổi 30s/lần) trên máy của User B.
- **Then** Hệ thống xác thực mã TOTP thành công, trạng thái chuyển thành "Hoàn thành".
- **And** Lệnh đóng băng được giải phóng, 100 Xu được hoàn trả về Ví của cả 2 user.

### AC3: User B bùng kèo (Không có gói Premium)
- **Given** Đã quá giờ hẹn 30 phút, mã QR không được quét. User B không có gói Q-Love Premium.
- **When** Hệ thống quét Cron Job xử lý quá giờ.
- **Then** Cuộc hẹn bị Hủy.
- **And** 100 Xu của User B bị tịch thu. Chuyển 90 Xu vào ví User A và 10 Xu vào tài khoản quỹ hệ thống (Phí vận hành).

### AC4: User B bùng kèo (Có gói Premium)
- **Given** Đã quá giờ hẹn 30 phút, mã QR không được quét. User B đang có gói Q-Love Premium và chưa dùng quyền miễn trừ trong tháng này.
- **When** Hệ thống quét Cron Job xử lý quá giờ.
- **Then** Cuộc hẹn bị Hủy.
- **And** Hệ thống không tịch thu Xu của B, hoàn lại 100 Xu bị đóng băng cho cả A và B. Ghi nhận B đã dùng hết 1 lượt miễn trừ của tháng.

---

## 5. UC-P3-010: Chợ Thẻ Bài Profile Profile

### AC1: Tính giá Giá Khởi Điểm ban đầu
- **Given** Một user đạt Level 5 và lần đầu mở tính năng Chợ Thẻ Bài Profile.
- **When** Hệ thống khởi tạo mã Thẻ Bài cho profile này (VD: #NVA).
- **Then** Giá trị khởi điểm (Giá Khởi Điểm) của mã #NVA được set mặc định là 100 Xu ảo.

### AC2: Khớp lệnh Mua thành công
- **Given** Collector có 500 Xu trong ví. Mã #NVA đang giao dịch ở giá 120 Xu.
- **When** Collector đặt lệnh mua 2 Thẻ Bài #NVA với giá thị trường (Market price).
- **Then** Hệ thống trừ 240 Xu của Collector (Cộng thêm 2% phí giao dịch nếu có).
- **And** Collector được thêm vào danh sách Người Sưu Tầm của mã #NVA, sở hữu 2 Thẻ Bài.

### AC3: Cập nhật giá theo công thức động (5 phút/lần)
- **Given** Profile #NVA vừa nhận được 10 lượt Match mới và 1 lượt gửi Locket trong 5 phút qua.
- **When** Hệ thống chạy Cron Job tính giá định kỳ.
- **Then** Giá mới của #NVA được cộng thêm: `(10 * 0.4) + (1 * 0.3) = 4.3 Xu`.
- **And** Biểu đồ hình nến (Candlestick chart) của mã #NVA trên sàn tự động vẽ thêm 1 nhịp tăng giá.

---

## 6. UC-P1-001: Đăng ký và Quẹt thẻ tìm đối tượng

### AC1: Lọc đối tượng theo "Hệ Tâm Linh"
- **Given** Người dùng đang ở màn hình Quẹt thẻ chính.
- **When** Chọn bộ lọc "Hệ Tâm Linh" (Ví dụ: Cung hoàng đạo, Thần số học).
- **Then** Thuật toán chỉ hiển thị các Profile có độ tương thích tâm linh > 70%.

### AC2: Match thành công
- **Given** Người dùng A vuốt phải (Like) profile của Người dùng B.
- **When** Người dùng B cũng vuốt phải profile của Người dùng A.
- **Then** Hệ thống hiển thị popup "It's a Match!" cho cả hai.
- **And** Mở khóa phòng chat (Chatroom) để hai người có thể nhắn tin.

---

## 7. UC-P1-003: Thành lập Bang hội (Clan)

### AC1: Tạo Clan thành công
- **Given** Người dùng đạt Level 3 và có >= 500 Xu trong ví ảo.
- **When** Nhấn "Tạo Bang hội", nhập tên và logo hợp lệ.
- **Then** Hệ thống trừ 500 Xu, tạo Clan mới và cấp quyền Bang chủ cho người dùng.

### AC2: Lỗi không đủ điều kiện tạo Clan
- **Given** Người dùng chưa đạt Level 3 hoặc số dư ví < 500 Xu.
- **When** Người dùng cố gắng nhấn "Tạo Bang hội".
- **Then** Nút bấm bị vô hiệu hóa (Greyed out) hoặc hiển thị popup báo lỗi "Bạn cần đạt Level 3 và có tối thiểu 500 Xu".

### AC3: Tiếp củi Lửa Trại thành công
- **Given** Lửa Trại của Bang hội đang có 2 thành viên tương tác trong ngày.
- **When** Thành viên thứ 3 tiến hành Check-in GPS hoặc gửi Locket vào nhóm.
- **Then** Hệ thống cập nhật `daily_active_members` = 3.
- **And** `campfire_streak` tăng thêm 1 ngày, gửi thông báo "Lửa trại đang cháy rực" cho toàn bộ bang hội.

### AC4: Lửa Trại tắt do thiếu tương tác
- **Given** Bang hội chỉ có 2 thành viên tương tác trong ngày.
- **When** Đồng hồ điểm 00:00, Cronjob quét hệ thống.
- **Then** `last_campfire_at` được xác định là quá 24h.
- **And** Hệ thống reset `campfire_streak` về 0 và gửi thông báo "Lửa trại đã tắt".

### AC5: Nhận Buff x1.5 điểm
- **Given** Bang hội đã duy trì Lửa Trại liên tục đạt `campfire_streak` >= 7.
- **When** Một thành viên thực hiện Check-in chiếm địa bàn thành công (được 10 điểm).
- **Then** Điểm cộng vào `weekly_score` của Bang hội được tự động nhân 1.5 lần (thành 15 điểm).

---

## 8. UC-P2-005: Đảo Tình Yêu 3D (Streak Gamification)

### AC1: Nâng cấp Đảo khi giữ Streak
- **Given** Người dùng A và B có Streak liên tục đạt mốc 7 ngày.
- **When** Hệ thống chạy Cron Job kiểm tra mốc Streak mỗi đêm.
- **Then** Đảo Tình Yêu 3D chung của 2 người tự động nâng cấp (Ví dụ: từ Lều cỏ lên Nhà gỗ).

### AC2: Đảo héo úa khi đứt Streak
- **Given** Người dùng A và B không tương tác (nhắn tin/Locket) quá 24h.
- **When** Hệ thống phát hiện đứt Streak.
- **Then** Streak bị reset về 0, Đảo Tình Yêu 3D chuyển sang trạng thái "Héo úa" (chờ dùng Xu để khôi phục).

---

## 9. UC-P2-008: Trợ lý Cánh Gió "Mỏ Hỗn" (AI Wingman)

### AC1: Gợi ý tin nhắn tự động
- **Given** Trong một đoạn chat, người dùng A gõ tin nhắn nhưng xóa đi viết lại 3 lần hoặc chat bị dừng quá lâu.
- **When** Người dùng A bấm vào icon "Trợ lý AI".
- **Then** AI tự động đọc 5 dòng chat gần nhất và đưa ra 3 gợi ý tin nhắn phản hồi (sarcastic, hài hước, thả thính).

### AC2: Bói bài Tarot mỗi sáng (Daily Ice-breaker)
- **Given** Người dùng A và B đã match và đang trò chuyện.
- **When** Đồng hồ điểm 08:00 sáng mỗi ngày.
- **Then** AI tự động sinh một quẻ bói Tarot/Chiêm tinh (dựa trên Cung hoàng đạo của A và B).
- **And** Gửi thẳng kết quả vào chatroom để gợi ý chủ đề nói chuyện.

---

## 10. UC-P3-009: Đánh giá và Tra cứu CV Tình Trường (Ex-Rating)

### AC1: Chấm điểm đối phương (Rating)
- **Given** Người dùng A và B đã từng match và chat trên 50 câu (hoặc đã gặp mặt O2O).
- **When** Người dùng A Unmatch B, hệ thống hiển thị bảng "Đánh giá tình trường".
- **Then** Người dùng A có thể rate B theo thang điểm 1-5 sao và để lại tag (Ví dụ: #RedFlag, #Lịch_sự).

### AC2: Tra cứu ẩn danh CV Tình trường
- **Given** Người dùng A đang xem profile của C (chưa Match) và muốn tra cứu lịch sử của C.
- **When** A bấm "Tra cứu CV Tình trường" và thanh toán 50 Xu.
- **Then** Hệ thống hiển thị điểm số Rating trung bình và các Tag ẩn danh mà người yêu cũ/match cũ đã để lại cho C.

---

## 11. UC-P3-011: Đêm săn mồi (The Purge)

### AC1: Tham gia ghép đôi hỗn loạn (Matchmaking Worker)
- **Given** Sự kiện "Đêm Săn Mồi" diễn ra vào 22h00 tối Thứ 6.
- **When** Người dùng chọn "Tham gia ngay".
- **Then** Hàng đợi Redis nhận request, sau đó **Matchmaking Worker** xử lý ghép đôi ẩn danh (không hiện Avatar/Tên) ngẫu nhiên.
- **And** Sau 10 phút trò chuyện, nếu cả 2 cùng bấm "Lộ diện", Avatar và Profile thật mới được hiển thị.

---

## 12. AC Bổ Sung: Giới hạn tần suất gửi Locket (Rate Limit)

### AC1: Rate limit gửi ảnh Locket
- **Given** Người dùng A đã gửi **10 ảnh Locket** cho cùng một người dùng B trong vòng 1 giờ.
- **When** Người dùng A cố gắng gửi ảnh Locket thứ 11 trong cùng khung giờ đó.
- **Then** Hệ thống từ chối request, trả về lỗi `429 Too Many Requests`.
- **And** Hiển thị thông báo cho A: *"Bạn đã gửi quá nhiều Locket trong giờ này. Hãy thư giãn và thử lại sau [thời gian reset]."*

### AC2: Giới hạn không áp dụng cho gói Premium
- **Given** Người dùng A đang có gói **Q-Love Premium** và đã gửi 10 ảnh Locket trong 1 giờ.
- **When** Người dùng A gửi ảnh thứ 11.
- **Then** Hệ thống cho phép gửi bình thường, không áp dụng rate limit.

---

## 13. AC Bổ Sung: Circuit Breaker Chợ Thẻ Bài Profile (Chống Pump & Dump)

### AC1: Tạm dừng giao dịch khi giá biến động bất thường
- **Given** Mã Thẻ Bài `#NVA` đang giao dịch ở 200 Xu.
- **When** Trong vòng **5 phút**, giá `#NVA` tăng hoặc giảm hơn **30%** (lên 260 Xu hoặc xuống 140 Xu).
- **Then** Hệ thống tự động kích hoạt **Circuit Breaker**: tạm dừng toàn bộ lệnh mua/bán của mã `#NVA` trong **15 phút**.
- **And** Hiển thị banner trên trang Thẻ Bài: *"Giao dịch tạm dừng do biến động bất thường. Tiếp tục lúc [thời gian]."*

### AC2: Thông báo cho Người Sưu Tầm hiện hữu
- **Given** Circuit Breaker vừa được kích hoạt cho mã `#NVA`.
- **When** Hệ thống thực thi lệnh tạm dừng.
- **Then** Toàn bộ user đang nắm giữ Thẻ Bài `#NVA` nhận được **push notification** thông báo tình trạng tạm dừng.
- **And** Các lệnh đang chờ khớp (Pending Orders) bị **hủy tự động**, Xu được hoàn trả vào ví.

---

## 14. AC Bổ Sung: Rút đơn kiện (Tòa Án Tình Yêu)

### AC1: Người kiện rút đơn trước khi đủ 50 vote
- **Given** Vụ kiện do Người dùng A khởi tạo chống B đang trong thời gian bỏ phiếu và chưa đạt đủ 50 phiếu.
- **When** Người dùng A vào mục "Vụ kiện của tôi" và nhấn "Rút đơn kiện".
- **Then** Hệ thống **hủy vụ án**, gỡ bài khỏi tab "Hóng Drama".
- **And** Bồi thẩm đoàn nhận thông báo *"Vụ kiện đã được nguyên đơn rút lại."*
- **And** Không có hình phạt hoặc phần thưởng nào được áp dụng cho cả hai bên.

### AC2: Không thể rút đơn sau khi đã đủ 50 vote và có kết quả
- **Given** Vụ kiện đã nhận đủ 50 phiếu và đã có kết quả (guilty/not_guilty).
- **When** Người dùng A cố gắng rút đơn.
- **Then** Nút "Rút đơn" bị vô hiệu hóa (greyed-out), hệ thống hiển thị thông báo *"Vụ án đã có phán quyết và không thể rút lại."*

---

## 15. UC-P2-012: Bảng Truy Nã Hẹn Hò (Bounty Hunter Mode)

### AC1: Đăng lệnh truy nã
- **Given** Người dùng A có đủ Xu trong ví.
- **When** A đăng một lệnh truy nã hẹn hò với tiền thưởng X Xu.
- **Then** Lệnh được hiển thị lên bảng chung, hệ thống đóng băng X Xu trong ví của A.

### AC2: Nộp đơn và Hoàn thành
- **Given** Lệnh truy nã của A đang "open".
- **When** B nộp đơn ứng tuyển, A chọn B làm người thắng cuộc (matched).
- **And** A và B gặp mặt ngoài đời, quét mã QR thành công.
- **Then** Trạng thái lệnh chuyển thành "completed", hệ thống chuyển X Xu từ ví A (đang đóng băng) sang ví B.

---

## 16. UC-P3-013: Đấu Giá Đặc Quyền (Top-Tier Blind Auction)

### AC1: Mở phiên đấu giá mù
- **Given** Là ngày cuối tháng.
- **When** Hệ thống tự động chọn ra Top 5 user có giá trị Thẻ Bài cao nhất.
- **Then** Khởi tạo phiên đấu giá (blind auction) cho từng Top-Tier trong 24h.

### AC2: Đấu giá thành công và nhận đặc quyền
- **Given** Phiên đấu giá của Top-Tier C kết thúc.
- **When** Người dùng A là người trả giá cao nhất (Winner).
- **Then** Hệ thống khóa chat của C, C chỉ được phép nhắn tin với A trong vòng 24h tiếp theo.
- **And** Số Xu trúng thầu được chia 50% cho C và 50% bị hệ thống thu hồi (burn).

---

## 17. UC-P3-014: Tường Thành Phong Sát (Wall of Shame)

### AC1: Bị đưa lên Tường thành
- **Given** Người dùng A vừa bị Tòa Án phán quyết "Có Tội".
- **When** Phán quyết được lưu vào hệ thống.
- **Then** Thông tin của A tự động xuất hiện trên Tường Thành Phong Sát với thời hạn 24h.

### AC2: Ném cà chua
- **Given** Người dùng B đang xem Tường Thành Phong Sát và có đủ Xu.
- **When** B nhấn nút "Ném Cà Chua" vào A.
- **Then** Hệ thống trừ 1 Xu của B.
- **And** Biến đếm `tomatoes_thrown` của A tăng lên 1, đồng thời có hiệu ứng animation cà chua đập vào avatar.

---

## 18. UC-P2-015: Vibe Check Nửa Đêm

### AC1: Lấy bài hát Spotify
- **Given** Là 23:00 đêm, người dùng A mở app.
- **When** A vào tab Vibe Check.
- **Then** Hệ thống gọi API Spotify để lấy bài hát A đang nghe. Nếu A không nghe gì, yêu cầu A chọn 1 bài hát thủ công.

### AC2: Match qua bài hát
- **Given** A đang nghe bài "Lối Nhỏ", B cũng đang nghe bài "Lối Nhỏ" (hoặc B lướt thấy và bấm Like).
- **When** Cả hai cùng Like bài hát của nhau.
- **Then** Hệ thống tạo Match ẩn danh và mở chatroom có thời hạn 60 phút.

---

## 19. UC-P1-016: Nghề Cò Mối / Bà Nguyệt

### AC1: Nhận hoa hồng Cò mối
- **Given** User C (Wingman) đã gửi link ghép đôi cho A và B.
- **When** A và B đồng ý Match, sau đó gặp mặt thành công (quét QR).
- **Then** Hệ thống ghi nhận trạng thái `rewarded` cho C.
- **And** C nhận được 10% tổng giá trị Thẻ Bài của A và B (Quy đổi ra Xu).

---

## 20. UC-P3-017: PK Cướp Đoạt Thẻ Bài (The Steal)

### AC1: Mua thẻ Đạo Tặc
- **Given** A muốn cướp Thẻ `#NVA` từ B.
- **When** A vào cửa hàng mua Thẻ Đạo Tặc (1000 Xu).
- **Then** Hệ thống trừ 1000 Xu và cấp vật phẩm cho A.

### AC2: Kích hoạt PK và Thua cuộc
- **Given** A dùng Thẻ Đạo Tặc nhắm vào B.
- **When** Cả hai vào phòng PK minigame 10 giây.
- **And** B thắng minigame.
- **Then** A mất Thẻ Đạo Tặc. B giữ được Thẻ `#NVA` và nhận được 500 Xu (trích từ giá trị Thẻ Đạo Tặc của A).
