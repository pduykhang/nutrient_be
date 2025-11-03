# Hot Reload Development với Air

Hệ thống đã được cấu hình để sử dụng [Air](https://github.com/air-verse/air) cho hot reload trong Docker container.

## Features

- ✅ **Automatic rebuild** khi code thay đổi
- ✅ **Automatic restart** sau khi build xong
- ✅ **File watching** cho tất cả `.go` files
- ✅ **Colored logs** để dễ debug
- ✅ **Excludes** tmp, vendor, .git folders
- ✅ **Docker volume mapping** cho live changes

## Quick Start

### 1. Start với Hot Reload

```bash
# Start tất cả services với hot reload
make docker-dev

# Hoặc chỉ start API service
docker-compose -f deployments/docker-compose.yml up --build api
```

### 2. View Logs

```bash
# Xem logs của API service
make docker-logs

# Hoặc
docker-compose -f deployments/docker-compose.yml logs -f api
```

### 3. Stop Services

```bash
make docker-down
```

## Cách hoạt động

1. **File Changes**: Khi bạn edit bất kỳ `.go` file nào
2. **Air Detects**: Air tự động phát hiện changes
3. **Rebuild**: Air chạy build command (`go build`)
4. **Restart**: Application tự động restart với binary mới
5. **Watch Again**: Air tiếp tục watch cho changes tiếp theo

## Cấu hình

### .air.toml

File `.air.toml` chứa cấu hình cho Air:

```toml
[build]
  cmd = "go build -o ./tmp/nutrient-api ./cmd/api/main.go ..."
  bin = "./tmp/nutrient-api"
  args_bin = ["server", "--config=configs/config.dev.yaml"]
  
  # Exclude directories không cần watch
  exclude_dir = ["tmp", "vendor", ".git", "bin"]
  
  # Include only .go files
  include_ext = [".go", ".tpl", ".tmpl", ".html"]
```

### Dockerfile.dev

Dockerfile development có:
- Air đã được cài đặt
- Source code được mount qua volume
- `.air.toml` được copy vào container

### docker-compose.yml

Volume mapping:
```yaml
volumes:
  - ..:/app              # Mount toàn bộ project
  - /app/tmp             # Exclude tmp (Air sẽ tạo trong container)
  - /app/bin             # Exclude bin
```

## Excluded Directories

Các thư mục sau sẽ **KHÔNG** trigger rebuild:
- `tmp/` - Temporary files
- `bin/` - Build output
- `vendor/` - Dependencies
- `.git/` - Git files
- `node_modules/` - Node modules (nếu có)
- `.vscode/`, `.idea/` - IDE files
- `docs/`, `deployments/`, `examples/` - Documentation
- `scripts/`, `proto/` - Scripts và proto files

## Included Files

Chỉ các file sau sẽ trigger rebuild:
- `.go` files - Go source code
- `.tpl`, `.tmpl`, `.html` - Template files

## Debugging

### Xem Air logs

```bash
# Trong container
docker exec -it nutrient-api-dev cat tmp/air.log

# Hoặc view container logs
make docker-logs
```

### Manual rebuild trong container

```bash
# Enter container
docker exec -it nutrient-api-dev sh

# Chạy build manually
go build -o ./tmp/nutrient-api ./cmd/api/main.go ./cmd/api/server.go ./cmd/api/migrate.go ./cmd/api/version.go
```

### Restart Air manually

```bash
# Restart container
docker-compose -f deployments/docker-compose.yml restart api
```

## Troubleshooting

### 1. Changes không trigger rebuild

**Problem**: Sửa code nhưng không rebuild

**Solutions**:
- Check xem file có nằm trong `exclude_dir` không
- Check xem extension có nằm trong `include_ext` không
- Xem logs: `make docker-logs`

### 2. Build errors

**Problem**: Build failed

**Solutions**:
- Check syntax errors trong code
- Xem Air logs: `docker exec -it nutrient-api-dev cat tmp/air.log`
- Rebuild container: `make docker-rebuild`

### 3. Application không restart

**Problem**: Build thành công nhưng app không restart

**Solutions**:
- Check binary path trong `.air.toml`
- Check `args_bin` có đúng không
- Restart container: `docker-compose restart api`

### 4. Volume mapping issues

**Problem**: Changes không sync vào container

**Solutions**:
- Verify volume mount: `docker inspect nutrient-api-dev | grep Mounts`
- Check file permissions
- Rebuild: `make docker-rebuild`

## Best Practices

1. **Chỉ watch những gì cần**: Exclude unnecessary directories
2. **Use tmp for builds**: Không commit tmp/ folder
3. **Monitor logs**: Watch Air logs để debug issues
4. **Test after changes**: Verify app restart và hoạt động đúng

## Performance Tips

1. **Exclude lớn directories**: Thêm vào `exclude_dir`
2. **Use polling if needed**: Set `poll = true` nếu filesystem không support events
3. **Adjust poll interval**: `poll_interval = 500` (milliseconds)

## Manual Air Usage (Without Docker)

Nếu muốn dùng Air locally:

```bash
# Install Air
go install github.com/air-verse/air@latest

# Chạy Air
air -c .air.toml

# Hoặc với custom config
air -c custom.air.toml
```

## References

- [Air GitHub](https://github.com/air-verse/air)
- [Air Documentation](https://github.com/air-verse/air#usage)
- [Air Example Config](https://github.com/air-verse/air/blob/master/air_example.toml)
