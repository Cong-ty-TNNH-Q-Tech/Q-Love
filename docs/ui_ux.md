# UI/UX Design Specification
**Project Name:** Q-Love - Super App Hẹn Hò & Mạng Xã Hội Giải Trí  
**Version:** 1.0  
**Date:** August 2026  
**Audience:** UI/UX Designer, Flutter Developer, QA  

> Tài liệu này là cầu nối giữa Business Requirements ([brd.md](./brd.md)) và Implementation ([architecture.md](./architecture.md)). Mọi quyết định UI phải bám sát Use-case trong [ba.md](./ba.md) và Acceptance Criteria trong [uc_ac.md](./uc_ac.md).

---

## Table of Contents
1. [Design System (Hệ thống thiết kế)](#1-design-system)
2. [Navigation Architecture (Kiến trúc điều hướng)](#2-navigation-architecture)
3. [Screen Inventory (Danh sách màn hình)](#3-screen-inventory)
4. [Key User Flows (Luồng màn hình chi tiết)](#4-key-user-flows)
5. [Component Specs (Đặc tả Components)](#5-component-specs)
6. [Micro-animation Guidelines](#6-micro-animation-guidelines)
7. [Accessibility & Edge Cases](#7-accessibility--edge-cases)

---

## 1. Design System

### 1.1. Triết lý thiết kế (Design Philosophy)
Q-Love nhắm đến Gen Z và Millennials — thế hệ thích những thứ **táo bạo, vui vẻ và hơi hỗn**. Design language phải thể hiện điều này:

- **Bold & Playful:** Không ngại màu mạnh, font lớn, icon có tính cách
- **Dark-first:** Nền tối là mặc định (giảm mỏi mắt khi dùng app buổi tối - đúng với hành vi target user)
- **Gamified:** Mọi tương tác phải cho cảm giác có phần thưởng (animation, confetti, haptic)
- **Trustworthy in Finance:** Các màn hình liên quan đến Xu/Giao dịch phải toát lên sự tin cậy (clean layout, rõ số liệu)

---

### 1.2. Color Palette

#### Primary Colors (Màu chính)
| Token | Hex | HSL | Mục đích |
|:---|:---|:---|:---|
| `--color-primary` | `#FF3D6B` | `hsl(348, 100%, 62%)` | CTA chính, Like button, Accent |
| `--color-primary-dark` | `#C4003A` | `hsl(348, 100%, 38%)` | Pressed state |
| `--color-primary-glow` | `#FF3D6B33` | Primary + 20% opacity | Glow effect, Shadow |

#### Secondary Colors
| Token | Hex | HSL | Mục đích |
|:---|:---|:---|:---|
| `--color-gold` | `#FFB830` | `hsl(42, 100%, 59%)` | Premium badge, Xu currency icon |
| `--color-gold-dark` | `#E09B00` | `hsl(42, 100%, 44%)` | Pressed state của gold elements |
| `--color-teal` | `#00D4AA` | `hsl(168, 100%, 42%)` | Stock price up, Success state |
| `--color-red-stock` | `#FF4757` | `hsl(355, 100%, 64%)` | Stock price down, Error |

#### Background & Surface
| Token | Hex | Mục đích |
|:---|:---|:---|
| `--bg-base` | `#0D0D14` | Nền tối toàn app (Root background) |
| `--bg-surface` | `#1A1A2E` | Card, Modal, Bottom sheet |
| `--bg-surface-2` | `#16213E` | Input field, List item |
| `--bg-surface-elevated` | `#22224A` | Hover, Selected state |

#### Text Colors
| Token | Hex | Mục đích |
|:---|:---|:---|
| `--text-primary` | `#FFFFFF` | Tiêu đề, Text chính |
| `--text-secondary` | `#A0A3BD` | Placeholder, Label phụ |
| `--text-disabled` | `#4A4A6A` | Disabled state |
| `--text-inverse` | `#0D0D14` | Text trên nền sáng (Gold button) |

#### Semantic Colors
| Token | Hex | Mục đích |
|:---|:---|:---|
| `--color-success` | `#00D4AA` | Thành công |
| `--color-warning` | `#FFB830` | Cảnh báo |
| `--color-error` | `#FF4757` | Lỗi, Pass button |
| `--color-info` | `#3D91FF` | Thông tin |

---

### 1.3. Typography

**Font chính:** `Inter` (Google Fonts) — Modern, dễ đọc, hỗ trợ tiếng Việt tốt.  
**Font phụ (Display):** `Clash Display` — Dùng cho Hero text, Tên tính năng lớn.

| Token | Font | Weight | Size | Line Height | Mục đích |
|:---|:---|:---|:---|:---|:---|
| `--text-display-xl` | Clash Display | 700 | 40sp | 48sp | Hero screen title |
| `--text-display-lg` | Clash Display | 700 | 32sp | 40sp | Feature name, Match popup |
| `--text-headline` | Inter | 700 | 24sp | 32sp | Section header |
| `--text-title` | Inter | 600 | 18sp | 26sp | Card title, User name |
| `--text-body-lg` | Inter | 400 | 16sp | 24sp | Body text chính |
| `--text-body` | Inter | 400 | 14sp | 22sp | Body text phụ |
| `--text-caption` | Inter | 400 | 12sp | 18sp | Label, Timestamp |
| `--text-label` | Inter | 600 | 11sp | 16sp | Button text, Tag |

---

### 1.4. Spacing System

Sử dụng base-8 grid system:

| Token | Value | Mục đích |
|:---|:---|:---|
| `--space-xs` | 4dp | Icon padding, Tight spacing |
| `--space-sm` | 8dp | Khoảng cách nhỏ trong component |
| `--space-md` | 16dp | Padding mặc định card |
| `--space-lg` | 24dp | Section spacing |
| `--space-xl` | 32dp | Screen padding top/bottom |
| `--space-2xl` | 48dp | Hero section |

**Screen horizontal padding (Edge margin):** `20dp` cố định cho toàn app.

---

### 1.5. Border Radius

| Token | Value | Mục đích |
|:---|:---|:---|
| `--radius-sm` | 8dp | Input, Tag chip |
| `--radius-md` | 16dp | Card, Bottom sheet |
| `--radius-lg` | 24dp | Profile card, Action sheet |
| `--radius-xl` | 32dp | Large modal |
| `--radius-full` | 9999dp | Avatar, Badge, Pill button |

---

### 1.6. Elevation & Shadow

```
Level 1 (Card):     box-shadow: 0 4px 16px rgba(255,61,107, 0.08)
Level 2 (Modal):    box-shadow: 0 8px 32px rgba(0,0,0, 0.48)
Level 3 (Popup):    box-shadow: 0 16px 48px rgba(0,0,0, 0.64)
Glow Primary:       box-shadow: 0 0 20px rgba(255,61,107, 0.40)
Glow Gold:          box-shadow: 0 0 20px rgba(255,184,48, 0.40)
```

---

### 1.7. Icon System

- **Icon Library:** `Phosphor Icons` (Có cả Regular và Bold weight, hỗ trợ tốt cho dark theme)
- **Icon Size Standards:**
  - Navigation bar: `24dp`
  - In-content action: `20dp`
  - Inline with text: `16dp`
  - Hero/Feature icon: `32-48dp`

---

## 2. Navigation Architecture

### 2.1. App Structure

```
App Root
├── Auth Stack (Unauthenticated)
│   ├── SplashScreen
│   ├── OnboardingScreen (3 slides)
│   ├── PhoneInputScreen
│   └── OtpVerificationScreen
│
└── Main Shell (Authenticated) ← Bottom Navigation
    ├── Tab 1: Discover (Quẹt thẻ)
    │   ├── SwipeScreen
    │   ├── MatchPopupOverlay
    │   └── FilterSheet
    ├── Tab 2: Map (Bản đồ Clan)
    │   ├── MapScreen
    │   ├── LandmarkDetailSheet
    │   └── CheckInScreen
    ├── Tab 3: Messages (Chat)
    │   ├── MatchListScreen
    │   └── ChatScreen
    │       ├── LocketCameraSheet
    │       ├── AIWingmanSheet
    │       └── ContractSheet
    ├── Tab 4: Market (Sàn CK + Court)
    │   ├── MarketTabsScreen (Stock / Court)
    │   ├── StockDetailScreen
    │   └── CourtCaseDetailScreen
    └── Tab 5: Profile
        ├── ProfileScreen
        ├── WalletScreen
        ├── ClanScreen
        └── SettingsScreen
```

### 2.2. Bottom Navigation Bar Spec

| Tab | Icon | Label | Badge |
|:---|:---|:---|:---|
| Discover | `ph:cards` | Khám phá | Không |
| Map | `ph:map-pin` | Bản đồ | Điểm mới (số) |
| Messages | `ph:chat-circle-dots` | Tin nhắn | Unread count |
| Market | `ph:chart-line-up` | Thị trường | Hot case (🔥) |
| Profile | `ph:user-circle` | Hồ sơ | Notification dot |

**Active State:** Tab đang chọn hiển thị màu `--color-primary` với pill indicator phía trên.  
**Height:** `72dp` + safe area bottom (iPhone notch).

---

## 3. Screen Inventory

### Phase 1 — MVP (16 màn hình)

| ID | Tên màn hình | Tab | Priority |
|:---|:---|:---|:---:|
| S-001 | Splash Screen | — | P0 |
| S-002 | Onboarding (3 slides) | — | P0 |
| S-003 | Phone Input | — | P0 |
| S-004 | OTP Verification | — | P0 |
| S-005 | Profile Setup (Step 1/3: Info) | — | P0 |
| S-006 | Profile Setup (Step 2/3: Photos) | — | P0 |
| S-007 | Profile Setup (Step 3/3: Spiritual Filter) | — | P0 |
| S-008 | Swipe / Discover | Discover | P0 |
| S-009 | Match Popup Overlay | Discover | P0 |
| S-010 | Match List | Messages | P0 |
| S-011 | Chat Screen | Messages | P0 |
| S-012 | Locket Camera | Messages | P0 |
| S-013 | Court Cases List (Hóng Drama) | Market | P0 |
| S-014 | Court Case Detail | Market | P0 |
| S-015 | My Profile | Profile | P0 |
| S-016 | Wallet Screen | Profile | P0 |

### Phase 2 (10 màn hình)

| ID | Tên màn hình | Tab | Priority |
|:---|:---|:---|:---:|
| S-017 | Map Screen (Bản đồ Clan) | Map | P1 |
| S-018 | Check-in Camera | Map | P1 |
| S-019 | Clan List & Leaderboard | Map | P1 |
| S-020 | Clan Detail | Map | P1 |
| S-021 | Create Clan | Profile | P1 |
| S-022 | Island (Đảo Tình Yêu) | Chat | P1 |
| S-023 | Dating Contract Setup | Chat | P1 |
| S-024 | QR Code Show/Scan | Chat | P1 |
| S-025 | AI Wingman Suggestions | Chat | P1 |
| S-026 | Notification Center | — (Global) | P1 |

### Phase 3 (8 màn hình)

| ID | Tên màn hình | Tab | Priority |
|:---|:---|:---|:---:|
| S-027 | Stock Market (Sàn CK) | Market | P2 |
| S-028 | Stock Detail (Biểu đồ nến) | Market | P2 |
| S-029 | Trade Order Sheet | Market | P2 |
| S-030 | Ex-Rating (Đánh giá sau Unmatch) | Chat | P2 |
| S-031 | CV Tình Trường (Kết quả tra cứu) | Profile | P2 |
| S-032 | The Purge Lobby | Discover | P2 |
| S-033 | Anonymous Chat (The Purge) | Discover | P2 |
| S-034 | Reveal Confirmation | Discover | P2 |

---

## 4. Key User Flows

### 4.1. Flow: Đăng ký & Setup Profile (S-003 → S-008)

```
[S-003] Phone Input
    │
    │  User nhập SĐT → Tap "Gửi mã OTP"
    │  Button state: Loading spinner 2s
    ▼
[S-004] OTP Verification
    │  6-ô input tự động focus next
    │  Countdown "Gửi lại sau 60s"
    │  Khi điền đủ 6 số → Tự động submit (không cần nhấn nút)
    ▼
[S-005] Profile Setup — Step 1/3 (Personal Info)
    │  Progress bar 33%
    │  Fields: Tên, Ngày sinh (Date picker), Giới tính (3 chips)
    │  Bio (optional, max 150 ký tự)
    ▼
[S-006] Profile Setup — Step 2/3 (Photos)
    │  Progress bar 66%
    │  Grid 2x3 slot ảnh (Tối thiểu 1, tối đa 6 ảnh)
    │  Slot đầu tiên = Ảnh đại diện chính (có crown indicator)
    │  Drag-and-drop để sắp xếp lại
    ▼
[S-007] Profile Setup — Step 3/3 (Spiritual Filter)
    │  Progress bar 100%
    │  Cung hoàng đạo: 12 icon zodiac dạng grid
    │  Thần số học: Input ngày sinh tự tính Life Path Number
    │  Preference: "Tôi muốn hẹn hò với ai tương thích > 70%"
    │  Toggle switch
    ▼
[S-008] Swipe Screen (Main App)
    │  Confetti animation + "Chào mừng đến Q-Love! 🎉"
    │  Toast 3s rồi tắt
```

**Empty State (S-003):** Khi chưa điền SĐT — nút "Gửi mã" bị disabled (opacity 40%).

---

### 4.2. Flow: Gửi Blind Locket (S-011 → S-012)

```
[S-011] Chat Screen
    │  User tap icon Camera ở thanh tool bottom
    ▼
[S-012] Locket Camera Screen
    │  ┌─────────────────────────────┐
    │  │  [X]              [⚡ Streak: 7]│
    │  │                             │
    │  │    [Camera Viewfinder]      │
    │  │                             │
    │  │  Gaussian Blur Preview      │
    │  │  (Preview real-time theo    │
    │  │   streak hiện tại)          │
    │  │                             │
    │  │     [ 🔴 CHỤP & GỬI ]      │
    │  └─────────────────────────────┘
    │
    ├─ [NSFW Detected] → Error Toast "Ảnh không phù hợp 🚫"
    │                    Camera reset, không gửi
    │
    └─ [Success] → 
        Ripple animation từ nút chụp lan ra toàn màn hình
        Toast: "Đã gửi Locket 📸 ~2.5s"
        Auto-dismiss camera, quay về Chat
```

**Blur Preview Logic:**
- Streak < 10: Viewfinder hiển thị preview với Gaussian Blur 90% realtime
- Streak 10-29: Blur 50%
- Streak ≥ 30: Blur 0% (rõ nét)

**Rate Limit State (429):** Toast màu vàng "Bạn đã gửi 10 Locket trong giờ này. Reset lúc [HH:MM] ⏳"

---

### 4.3. Flow: Khế Ước Tài Chính (S-011 → S-023 → S-024)

```
[S-011] Chat Screen
    │  Trong ≥ 20 tin nhắn → Icon 💍 Contract xuất hiện ở toolbar
    │  (Icon ẩn nếu < 20 tin, thay bằng lock icon với tooltip)
    │
    │  User tap 💍
    ▼
[S-023] Dating Contract Setup Sheet (Bottom Sheet, 80% height)
    │  ┌──────────────────────────────┐
    │  │  ━━━━━━  (Drag handle)       │
    │  │  🤝 Tạo Khế Ước Hẹn Hò      │
    │  │                              │
    │  │  Số Xu cọc mỗi bên:         │
    │  │  [  - ] [ 100 Xu ] [ + ]    │
    │  │  Số dư của bạn: 350 Xu ✓    │
    │  │                              │
    │  │  📅 Thời gian hẹn:          │
    │  │  [ Chọn ngày giờ... ]       │
    │  │                              │
    │  │  📍 Ghi chú địa điểm:       │
    │  │  [ Café The Coffee House... ]│
    │  │                              │
    │  │  ⚠️ 100 Xu sẽ bị đóng băng  │
    │  │  cho đến khi hẹn hoàn thành  │
    │  │                              │
    │  │  [    Gửi lời mời Khế Ước  ]│
    │  └──────────────────────────────┘
    │
    ├─ [Insufficient Balance] → Cột số dư chuyển màu đỏ
    │                           "Bạn chỉ còn X Xu. Cần X thêm."
    │                           Button "Gửi" bị disabled
    │
    └─ [Success] → 
        Chat bubble đặc biệt xuất hiện: 
        "💍 Khế ước đã được gửi. 100 Xu đã bị đóng băng."
        B nhận notification và thấy trong chat bubble [Xác nhận / Từ chối]

[Khi đến giờ hẹn, B bật màn hình QR]
    ▼
[S-024] QR Code Show/Scan Screen
    │  ┌──────────────────────────────┐
    │  │  [B] Màn hình hiển thị QR   │
    │  │                              │
    │  │  [QR Code 256x256]           │
    │  │  Rotate mỗi 30 giây         │
    │  │  ⏱️ Còn: 28s                │
    │  │                              │
    │  │  ⛔ KHÔNG CHỤP MÀN HÌNH    │
    │  │  (Screenshot block active)   │
    │  └──────────────────────────────┘
    │
    │  [A] Màn hình Camera Quét QR
    │  └─ Quét thành công → 
    │      Confetti 🎉 + Toast "Gặp mặt thành công! 100 Xu đã được hoàn lại"
    │      Trạng thái Khế Ước → "Hoàn thành ✅"
```

---

### 4.4. Flow: Tòa Án Tình Yêu (S-011 → S-014)

```
[S-011] Chat Screen (A bị ghost)
    │  Hiển thị Ghost Indicator nếu B im lặng > 24h:
    │  "👻 B đã im lặng 2 ngày rồi..."
    │  Button đỏ: [Đâm đơn kiện]
    │
    │  User tap [Đâm đơn kiện]
    ▼
[Bottom Sheet] Chọn lý do kiện
    │  • 👻 Ghosting (im lặng đột ngột)
    │  • 🪤 Trapboy / Trapgirl (câu like rồi biến)
    │  • 💸 Bùng kèo (không đến đúng giờ)
    │  • 😤 Nói lời xúc phạm
    │
    │  [Xác nhận đâm đơn →]
    ▼
[Modal Confirmation]
    │  "Bạn đang chuẩn bị đưa [***] ra Tòa án tình yêu."
    │  "Thông tin sẽ được ẩn danh hoàn toàn."
    │  [Hủy] [Xác nhận kiện]
    ▼
[S-013] Court Cases List (Tab Market)
    │  Vụ kiện mới xuất hiện với badge "Mới 🔥"
    ▼
[S-014] Court Case Detail
    │  ┌──────────────────────────────┐
    │  │  🏛️ TÒA ÁN TÌNH YÊU        │
    │  │  ─────────────────────────  │
    │  │  Nguyên đơn: Ng*** (ẩn danh)│
    │  │  Bị đơn:     Tr*** (ẩn danh)│
    │  │  Lý do:      Ghosting        │
    │  │                              │
    │  │  📜 Bằng chứng (5 tin nhắn) │
    │  │  ┌────────────────────────┐ │
    │  │  │ [***]: "Chiều nay đi   │ │
    │  │  │  cà phê không?"        │ │
    │  │  │ [Ng***]: "OK !"        │ │
    │  │  │ [***]: ...im lặng...   │ │
    │  │  └────────────────────────┘ │
    │  │                              │
    │  │  ⏱️ Còn 8 tiếng để vote     │
    │  │  48 / 50 phiếu biểu quyết   │
    │  │  [████████░░] 76% Có tội    │
    │  │                              │
    │  │  [  Có tội 😡  ] [ Vô tội 🕊️]│
    │  └──────────────────────────────┘
```

**Verdict State:**  
- `guilty`: Banner đỏ "⚖️ Phán quyết: CÓ TỘI. Bị đơn bị Ghost Badge."
- `not_guilty`: Banner xanh "⚖️ Phán quyết: VÔ TỘI. Vụ án đã kết thúc."

---

### 4.5. Flow: Sàn Chứng Khoán (S-027 → S-028 → S-029)

```
[S-027] Stock Market Screen
    │  Header: "📈 Sàn Chứng Khoán Độc Thân"
    │
    │  Tabs: [Trending 🔥] [Đang tăng ↑] [Đang giảm ↓] [Danh mục]
    │
    │  List item mỗi cổ phiếu:
    │  ┌────────────────────────────────┐
    │  │ [Avatar] $NVA  Nguyễn Văn A   │
    │  │          147.5 Xu  +4.3 (+3%) │ ← Màu teal nếu tăng
    │  │          [Mini Sparkline Chart]│
    │  └────────────────────────────────┘
    │
    │  [Circuit Breaker Active Banner]
    │  ⚡ "$ABC đang tạm dừng giao dịch — Còn 12:34"
    │  (Màu vàng, sticky top)
    ▼
[S-028] Stock Detail Screen
    │  Header: "$NVA — Nguyễn Văn A"
    │  Avatar + Level badge
    │
    │  [Candlestick Chart - chiếm 40% màn hình]
    │  Time selector: [1H] [1D] [1W] [1M]
    │
    │  Chỉ số:
    │  ┌─────────────┬──────────────┐
    │  │ Giá hiện tại│ 147.5 Xu     │
    │  │ Biến động   │ +4.3 (+3%)  │
    │  │ Lượt Match  │ 12 (hôm nay)│
    │  │ Clan Score  │ 85 điểm     │
    │  │ Cổ đông     │ 23 người    │
    │  └─────────────┴──────────────┘
    │
    │  [  🛒 MUA  ] [  📤 BÁN  ]
    ▼
[S-029] Trade Order Bottom Sheet
    │  "Mua cổ phiếu $NVA"
    │  Giá hiện tại: 147.5 Xu
    │
    │  Số lượng: [ - ] [ 2 ] [ + ]
    │
    │  Chi tiết:
    │  Tiền vốn:   295.0 Xu
    │  Phí (2%):   5.9 Xu
    │  Tổng:       300.9 Xu
    │
    │  Số dư:      350.0 Xu ✓
    │
    │  [    Xác nhận đặt lệnh Mua    ]
```

---

### 4.6. Flow: The Purge (S-032 → S-033 → S-034)

```
[Push Notification 21:45 Thứ 6]
"🎭 The Purge bắt đầu sau 15 phút! Bạn sẵn sàng chưa?"

[S-032] The Purge Lobby
    │  Full-screen dark background
    │  Glitch effect animation trên title "THE PURGE"
    │  Countdown: 00:14:32
    │  "👁️ 1,247 người đang chờ ghép đôi"
    │
    │  [   ⚡ THAM GIA NGAY   ]
    │
    │  Sau 22:00 → Lobby chuyển sang "Đang tìm đối tượng..."
    │  Loading animation dạng radar pulse
    ▼
[S-033] Anonymous Chat Screen
    │  Hoàn toàn khác chat thường:
    │  - Background: Tối hơn (#080810)
    │  - Avatar: 🎭 Mask icon thay cho ảnh thật
    │  - Tên: "Đối tượng ẩn danh"
    │
    │  Timer countdown to góc trên: "🕐 09:32"
    │
    │  Khi còn 1 phút: Timer đổi màu đỏ, rung nhẹ
    │
    │  Khi hết giờ:
    │  → Màn hình mờ dần
    │  → Nút "Lộ diện" xuất hiện ở giữa (30s timeout)
    ▼
[S-034] Reveal Confirmation Screen
    │  "Bạn có muốn lộ diện với đối tượng này không?"
    │
    │  [Ảnh mờ blur 100%] [Ảnh mờ blur 100%]
    │  "Bạn"              "Họ"
    │
    │  Timer: 28s (countdown)
    │
    │  [  👁️ LỘ DIỆN  ] [  Thôi, lần sau  ]
    │
    │  ├─ [Cả 2 cùng bấm] → 
    │  │   Blur ảnh fade out từ từ (1.5s)
    │  │   "✨ Match thật sự!" → Modal hỏi Match chính thức
    │  │
    │  └─ [1 bên không bấm] →
    │       "Lần này không duyên. Thử lại tuần sau!" 
    │       Màn hình fade to black
```

---

## 5. Component Specs

### 5.1. Profile Swipe Card (Quẹt thẻ)

```
┌─────────────────────────────┐ ← border-radius: 24dp
│                             │   width: screen - 32dp
│     [Profile Photo]         │   height: 70vh
│     (Full bleed, object-fit │
│      cover)                 │
│                             │
│  ╔═══════════════════════╗  │ ← Gradient overlay from bottom
│  ║ Nguyễn Văn A, 25 ♂   ║  │   Linear: transparent → #0D0D14 80%
│  ║ 📍 2.3km              ║  │
│  ║ ♐ Nhân Mã • 86% 💫   ║  │ ← Spiritual compatibility badge
│  ║ "Thích cà phê và..."  ║  │
│  ╚═══════════════════════╝  │
└─────────────────────────────┘

[ ✕ PASS ]        [ ♥ LIKE ]
  Màu #FF4757      Màu #FF3D6B
  Size: 64dp       Size: 72dp (lớn hơn để encourage like)
  border-radius: full
  Glow khi hover/press
```

**Swipe Physics:**
- Ngưỡng Like: kéo sang phải > 30% màn hình
- Ngưỡng Pass: kéo sang trái > 30% màn hình
- Card rotation: `rotate(swipeX * 0.1deg)` khi kéo
- Like overlay: Heart icon đỏ fade in khi kéo phải (opacity tỷ lệ với khoảng cách)
- Pass overlay: ✕ icon xanh fade in khi kéo trái

---

### 5.2. Locket Widget (iOS WidgetKit / Android AppWidget)

```
Small Widget (2x2):
┌──────────────────┐
│ [Ảnh mờ]        │ ← Gaussian Blur 90% (nếu Streak < 10)
│                  │
│ Q-Love    🔥 7  │ ← Logo + Streak số
└──────────────────┘

Medium Widget (2x4):
┌──────────────────────────────────┐
│ [Ảnh mờ/rõ tùy Streak]         │
│                                  │
│ 💌 Nguyễn Văn A gửi cho bạn     │
│ Q-Love                    🔥 7  │
└──────────────────────────────────┘
```

**Widget States:**
| State | Hiển thị |
|:---|:---|
| `no_match` | Placeholder "Chưa có Match. Quẹt thẻ ngay!" |
| `waiting` | "Đang chờ Locket từ [Tên]..." |
| `received_blurred` | Ảnh mờ + "Nhấn để xem và phản hồi" |
| `received_clear` | Ảnh rõ (Streak ≥ 30) |
| `error` | "Không tải được ảnh. Nhấn để thử lại" |

---

### 5.3. Xu / Token Display

```
Trong Wallet Screen:
┌────────────────────────────────┐
│  💰 Ví Xu của bạn              │
│  ┌──────────────┬────────────┐ │
│  │  Khả dụng   │ Đóng băng  │ │
│  │  350 Xu 🟡  │ 100 Xu 🔒  │ │
│  └──────────────┴────────────┘ │
│  [  + Nạp Xu  ]               │
└────────────────────────────────┘

Inline (trong app):
🟡 350 Xu  ← Icon coin vàng + số + text "Xu"
```

**Màu sắc Xu:**
- Đủ tiền để thực hiện hành động: `--color-gold`
- Thiếu tiền (warning): `--color-warning` + shake animation
- Balance 0: `--color-error`

---

### 5.4. Streak Fire Badge

```
🔥 7  ← Text label
```

| Streak | Icon | Màu | Animation |
|:---|:---|:---|:---|
| 0 | 💀 | `#4A4A6A` | None |
| 1-6 | 🔥 | `#FF6B35` | None |
| 7-29 | 🔥 | `#FF3D6B` | Pulse nhẹ (1s interval) |
| 30-89 | 🔥🔥 | `#FF3D6B` | Pulse nhanh hơn |
| 90+ | 🔥🔥🔥 | Gradient vàng-đỏ | Flame flicker animation |

---

### 5.5. Ghost Badge (Huy hiệu trên Profile)

```
Hiển thị trên Profile Card:
┌──────────────────┐
│  👻 GHOST THỦ   │ ← Nền đen mờ, viền đỏ
│  Còn 2 ngày      │
└──────────────────┘
```

Badge xuất hiện ở góc trên-phải ảnh đại diện, dạng chip với countdown text.

---

### 5.6. Đảo Tình Yêu 3D — State Machine

```
Level 1: Lều cỏ      (Streak 0-6)   — Màu xanh lá nhạt
Level 2: Nhà gỗ      (Streak 7-29)  — Màu nâu ấm
Level 3: Nhà gạch    (Streak 30-59) — Màu cam đất
Level 4: Biệt thự    (Streak 60-89) — Màu tím hoàng gia
Level 5: Lâu đài     (Streak 90+)   — Màu vàng + particle effect
```

**Trạng thái đặc biệt:**
- `withered` (Héo úa): Desaturate toàn bộ đảo + mưa rơi animation + màu xám
- `rebuilding` (Đang tái thiết sau khi mua): 3 giây grow animation từ Level 1

---

## 6. Micro-animation Guidelines

### 6.1. Nguyên tắc chung
- **Duration:** Tất cả transition dưới **350ms** để không làm chậm UX
- **Easing:** 
  - Vào màn hình: `cubic-bezier(0.22, 1, 0.36, 1)` (ease out quart)
  - Ra màn hình: `cubic-bezier(0.76, 0, 0.24, 1)` (ease in quart)
  - Spring: `spring(mass: 1, stiffness: 200, damping: 20)` cho card swipe
- **Tool:** `Rive` cho animation phức tạp (Đảo, Locket blur effect), Flutter `AnimationController` cho micro-interaction đơn giản

---

### 6.2. Animation Inventory

| Trigger | Animation | Duration | Easing |
|:---|:---|:---|:---|
| Match thành công | Confetti burst + Heartbeat popup scale 0→1 | 600ms | Spring |
| Swipe Like | Card bay ra phải + fade + Heart trail | 300ms | Ease out |
| Swipe Pass | Card bay ra trái + fade + X trail | 250ms | Ease out |
| Gửi Locket | Ripple từ nút lan ra màn hình | 400ms | Ease out |
| Streak tăng | Fire icon scale 1→1.3→1 + particle | 200ms | Spring |
| Nâng cấp Đảo | Glow pulse + Island grow + Stars | 1500ms | Ease in-out |
| Đảo héo úa | Color desaturate + wither particle | 1000ms | Ease in |
| Đổ Xu (chi) | Số đếm ngược + shake + màu đỏ nhấp | 400ms | Ease out |
| Nhận Xu (thu) | Số đếm tăng + coin flip animation | 500ms | Spring |
| Circuit Breaker | Banner slide down + pulse warning | 300ms | Ease out |
| The Purge timer < 1 phút | Timer đổi đỏ + shake 2px | Loop | Ease in-out |
| Verdict guilty | Gavel drop sound cue + screen flash đỏ | 800ms | Ease out |
| QR timeout | Ô QR blur dần + countdown spin | 500ms | Linear |

---

### 6.3. Haptic Feedback Map

| Sự kiện | Haptic Pattern |
|:---|:---|
| Swipe Like | `HeavyImpact` |
| Swipe Pass | `LightImpact` |
| Match | `SuccessNotification` (double pulse) |
| Gửi Locket | `MediumImpact` |
| Cảnh báo thiếu Xu | `ErrorNotification` |
| Vote tại Tòa Án | `SelectionClick` |
| QR Scan thành công | `SuccessNotification` |
| The Purge Lộ diện | `HeavyImpact` + `SuccessNotification` |

---

## 7. Accessibility & Edge Cases

### 7.1. Empty States

| Màn hình | Empty State Content |
|:---|:---|
| Swipe Feed | Illustration chú chim nhỏ cô đơn + "Bạn đã quẹt hết người rồi 😅 Thử mở rộng bán kính?" + [Mở rộng bán kính] |
| Match List | Illustration + "Chưa có Match nào. Bắt đầu quẹt thẻ ngay!" + [Khám phá] |
| Court Cases | "Hôm nay Q-Love yên bình 🕊️ Không có vụ kiện nào đang chờ." |
| Wallet Transactions | "Chưa có giao dịch nào. Hãy mua Xu để bắt đầu!" |
| Clan Leaderboard | "Bang hội của bạn chưa check-in tuần này. Dẫn đầu bảng xếp hạng nào!" |

### 7.2. Error States

| Lỗi | UI |
|:---|:---|
| Mất kết nối mạng | Banner sticky trên cùng: "📡 Mất kết nối. Đang thử lại..." |
| API 500 | Toast đỏ "Có lỗi xảy ra. Thử lại sau." + Retry button |
| OTP sai | Các ô input rung + viền đỏ + text "Mã OTP không chính xác" |
| Fake GPS | Full-screen error: "🚫 Phát hiện giả lập vị trí. Tính năng Check-in bị khóa 7 ngày." |

### 7.3. Loading States

- **Skeleton Screen:** Dùng cho Feed, Match List, Court List (không dùng spinner)
- **Shimmer Color:** `#1A1A2E` → `#22224A` (phù hợp dark theme)
- **Pull-to-Refresh:** Custom animation: trái tim đập + Q logo xoay

### 7.4. Permission Requests

| Permission | Màn hình yêu cầu | Context |
|:---|:---|:---|
| Camera | Trước S-006 (Setup ảnh) | "Cần quyền Camera để upload ảnh đại diện" |
| Location | Trước S-008 (Swipe) | "Cần vị trí để hiển thị người gần bạn" |
| Notification | Sau Match đầu tiên | "Bật thông báo để không bỏ lỡ Locket mới!" |
| Widget (iOS) | Sau khi Match + gửi Locket đầu tiên | "Thêm Q-Love Widget để nhận Locket ngay trên màn hình chính!" |

> **Rule:** Không yêu cầu nhiều permission cùng lúc. Mỗi permission request phải xuất hiện đúng lúc user thực sự cần tính năng đó (Contextual permission).

### 7.5. Screen Size Support

| Device | Chú ý |
|:---|:---|
| iPhone SE (375px) | Swipe card height giảm xuống 60vh, font scale xuống 0.9 |
| iPhone 15 Pro Max (430px) | Layout chuẩn |
| Android nhỏ (360px) | Bottom nav label ẩn, chỉ icon |
| Tablet (≥768px) | Màn hình Swipe layout 2 cột (card + profile detail) |

---

*Tài liệu này nên được đồng bộ với file Figma của team Design. Mọi thay đổi lớn cần được review bởi cả Design Lead và Tech Lead trước khi implement.*
