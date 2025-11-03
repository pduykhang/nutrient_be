# Demo: Ví dụ về xung đột tên khi không dùng .PHONY

## Tình huống vấn đề

Khi có file/folder trùng tên với target trong Makefile và KHÔNG khai báo `.PHONY`:

### Trường hợp 1: File "build" tồn tại

```bash
# Tạo file có tên trùng với target
$ echo "fake file" > build

# Chạy make build
$ make build
make: 'build' is up to date.  # ❌ Make nghĩ file "build" đã up-to-date, không chạy command!
```

### Trường hợp 2: Folder "build" tồn tại

```bash
# Tạo folder có tên trùng với target
$ mkdir build

# Chạy make build
$ make build
make: 'build' is up to date.  # ❌ Make kiểm tra folder mới, không chạy command!
```

## Tại sao Make nghĩ "up to date"?

Make có logic:
1. Kiểm tra target (ví dụ: `build`)
2. Nếu TỒN TẠI file/folder trùng tên → Make nghĩ target đã được "tạo ra"
3. So sánh timestamp với dependencies
4. Nếu target mới hơn → "up to date", không chạy command

## Giải pháp: Dùng .PHONY

```makefile
# Khai báo target là PHONY (không tạo ra file)
.PHONY: build test run

build:
	go build -o bin/app ./cmd

test:
	go test ./...
```

### Bây giờ Make sẽ:
1. Thấy `.PHONY: build` → Biết đây là target ảo (không tạo file)
2. KHÔNG kiểm tra file/folder trùng tên
3. LUÔN chạy command trong target
4. Tăng hiệu suất (không cần kiểm tra timestamp)

## Ví dụ thực tế

```bash
# Terminal 1: Tạo file trùng tên
$ touch build test lint run

# Terminal 2: Chạy make (KHÔNG dùng .PHONY)
$ make build
make: 'build' is up to date.  # ❌ Không chạy!
$ make test  
make: 'test' is up to date.    # ❌ Không chạy!

# Có .PHONY trong Makefile
$ make build
go build -o bin/nutrient-api cmd/api/main.go  # ✅ Chạy ngay!
$ make test
go test -v -race -coverprofile=coverage.out ./...  # ✅ Chạy ngay!
```

## Best Practices

### ✅ ĐÚNG: Khai báo tất cả target không tạo file
```makefile
.PHONY: build run test clean docker-up docker-down migrate lint help

build:
	go build -o bin/app ./cmd

run:
	go run ./cmd

test:
	go test ./...

clean:
	rm -rf bin/
```

### ❌ SAI: Không khai báo .PHONY
```makefile
# Nếu có file "build" trong project → xung đột!
build:
	go build -o bin/app ./cmd
```

## Khi nào KHÔNG cần .PHONY?

Khi target TẠO RA file:
```makefile
# binary là file output, KHÔNG cần .PHONY
binary: main.go
	go build -o binary main.go

# Chạy: make binary
# Kết quả: Tạo file "binary"
```

## Tóm tắt

1. `.PHONY` đánh dấu target "ảo" (không tạo file)
2. Giúp tránh xung đột khi có file/folder trùng tên
3. Tăng hiệu suất (không kiểm tra timestamp)
4. Nên dùng cho tất cả target "command" (build, run, test, clean, etc.)
5. Không cần cho target tạo file thực (compiler output)

## Command để test

```bash
# Tạo file trùng tên với target
$ touch build test lint run

# Chạy make (Makefile hiện tại có .PHONY nên vẫn OK)
$ make build    # ✅ Chạy bình thường
$ make test    # ✅ Chạy bình thường
$ make lint    # ✅ Chạy bình thường

# Xóa file test
$ rm -f build test lint run
```
