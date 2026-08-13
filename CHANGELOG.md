# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - Sắp tới

### Added
- [backend/server/api] Thêm endpoint `/api/v1/upload/presigned-url` cho chức năng Upload Avatar (tích hợp Cloudflare R2).
- [backend/server/api] Thêm endpoint Nghề Cò Mối (Wingman) với đầy đủ CRUD và xử lý hoa hồng Ví Ảo (Transaction SERIALIZABLE).
- [backend/server/models] Thêm model `wingman_referrals` và `wallet_transactions`.
- [backend/server/services] Thêm `WingmanService` với logic Referral và hoa hồng (Commission).

### Changed
- [backend/server/Dockerfile] Cập nhật base image lên `golang:alpine` để tương thích với `go 1.25.0`.
- [docs] Cập nhật `README.md` theo định hướng Professional Executive Summary cho nhà đầu tư.

### Fixed
- [backend/server/tests] Sửa lỗi sai argument count của `go-sqlmock` trong `wingman_service_test.go`.
