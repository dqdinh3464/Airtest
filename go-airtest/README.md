# Go Airtest (airtestgo)

Một bản triển khai **Airtest-like** bằng Golang, tập trung vào Android qua `adb`.

## Mục tiêu

Dự án này mô phỏng một phần workflow cốt lõi của Airtest:
- Chụp màn hình thiết bị.
- Tìm ảnh mẫu (template matching) trên màn hình.
- Tự động chạm theo ảnh mẫu (`touch(template)`).
- Điều khiển thao tác cơ bản (tap theo toạ độ).

## Tính năng hiện có

- `screenshot`: chụp ảnh màn hình bằng `adb exec-out screencap -p`.
- `find`: tìm vị trí template trong màn hình hiện tại.
- `touch`: tìm template và chạm vào tâm template nếu đủ ngưỡng confidence.
- `tap`: chạm toạ độ tuyệt đối.

## Cài đặt

Yêu cầu:
- Go >= 1.21
- Android SDK platform-tools (`adb`) trong PATH
- Thiết bị Android bật USB debugging

Build:

```bash
cd go-airtest
go build -o airtestgo ./cmd/airtestgo
```

## Sử dụng

Chụp màn hình:

```bash
./airtestgo screenshot --serial <device_serial> --out screen.png
```

Tìm ảnh mẫu:

```bash
./airtestgo find --serial <device_serial> --template tpl_login_button.png
```

Tìm và chạm theo ảnh mẫu:

```bash
./airtestgo touch --serial <device_serial> --template tpl_login_button.png --threshold 0.90
```

Chạm toạ độ:

```bash
./airtestgo tap --serial <device_serial> --x 500 --y 900
```

## Kiến trúc

- `internal/android`: lớp tích hợp `adb` (screen capture + input).
- `internal/vision`: thuật toán template matching (grayscale + normalized SAD).
- `internal/runner`: ghép device + vision thành engine thao tác.
- `cmd/airtestgo`: CLI.

## Giới hạn hiện tại

- Mới hỗ trợ Android qua `adb`.
- Thuật toán matching đang tối ưu cho dễ hiểu (chưa tối ưu tốc độ cho ảnh lớn).
- Chưa có report HTML, record video, hay DSL script như Airtest đầy đủ.

## Hướng mở rộng

- Thêm thực thi test script (YAML/JSON DSL).
- Thêm đa nền tảng (Windows/iOS).
- Thêm assertion API (`assert_exists`, `assert_not_exists`).
- Thêm report dạng HTML + step log + screenshot timeline.
