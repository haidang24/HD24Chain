# HD24Chain - Nền Tảng Blockchain Doanh Nghiệp
copyright : haidang24
##  Giới Thiệu

HD24Chain là nền tảng blockchain cấp doanh nghiệp, sẵn sàng cho sản xuất, được xây dựng trên nền tảng Go Ethereum với cơ chế đồng thuận POVA (Proof of Validator Authority) tùy chỉnh. HD24Chain cung cấp hiệu suất cao, bảo mật mạnh mẽ và khả năng mở rộng cho các ứng dụng doanh nghiệp.

##  Tính Năng Chính

-  **Đồng Thuận POVA**: Cơ chế đồng thuận tùy chỉnh với luân chuyển validator xác định
-  **Hiệu Suất Cao**: Thời gian block 20 giây, tối ưu hóa bộ nhớ và CPU
-  **Bảo Mật Doanh Nghiệp**: Xác thực đa lớp, mã hóa end-to-end
-  **Tương Thích Ethereum**: Hỗ trợ EVM đầy đủ, smart contract Solidity
-  **API Đa Dạng**: HTTP RPC, WebSocket, IPC
-  **Quản Lý Validator**: Hệ thống quản lý validator linh hoạt
-  **Monitoring**: Theo dõi hiệu suất và trạng thái real-time

##  Thông Số Kỹ Thuật

| Tính Năng | Thông Số |
|-----------|----------|
| **Đồng Thuận** | POVA (Proof of Validator Authority) |
| **Thời Gian Block** | 20 giây (có thể cấu hình) |
| **Chain ID** | 99999 (có thể tùy chỉnh) |
| **Gas Limit** | 16,777,216 gas mỗi block |
| **Validators** | 1-10 (có thể cấu hình) |
| **Hỗ Trợ API** | HTTP RPC, WebSocket, IPC |
| **Nền Tảng** | Windows Server 2019/2022, Windows 10/11 |
| **Database** | Pebble (LevelDB compatible) |
| **Network Protocol** | P2P Ethereum-compatible |

##  Hướng Dẫn Cài Đặt và Chạy

### **Yêu Cầu Hệ Thống**

#### Phần Cứng
- **CPU**: 4+ cores (khuyến nghị: 8+ cores)
- **RAM**: 8GB minimum (khuyến nghị: 16GB+)
- **Ổ cứng**: 100GB+ SSD (khuyến nghị: 500GB+ NVMe)
- **Mạng**: Kết nối internet ổn định

#### Phần Mềm
- **Windows**: 10/11 hoặc Server 2019/2022
- **PowerShell**: 5.1 hoặc cao hơn
- **Go**: 1.20+ (để build từ source nếu cần)

### **Bước 1: Tạo Cấu Trúc Thư Mục**

```powershell
# Tạo thư mục chính
New-Item -ItemType Directory -Path "C:\HD24Chain" -Force
New-Item -ItemType Directory -Path "C:\HD24Chain\data" -Force
New-Item -ItemType Directory -Path "C:\HD24Chain\logs" -Force
New-Item -ItemType Directory -Path "C:\HD24Chain\config" -Force

# Copy file thực thi
Copy-Item "mychain.exe" "C:\HD24Chain\" -Force
```

### **Bước 2: Tạo File Genesis**

Tạo file `C:\HD24Chain\genesis.json` với nội dung sau:

```json
{
    "config": {
        "chainId": 99999,
        "homesteadBlock": 0,
        "eip150Block": 0,
        "eip155Block": 0,
        "eip158Block": 0,
        "byzantiumBlock": 0,
        "constantinopleBlock": 0,
        "petersburgBlock": 0,
        "istanbulBlock": 0,
        "berlinBlock": 0,
        "londonBlock": 0,
        "pova": {
            "validators": [
                "0x[ĐỊA_CHỈ_VALIDATOR_SẼ_ĐƯỢC_TẠO]"
            ],
            "period": 20
        }
    },
    "difficulty": "0x1",
    "gasLimit": "0x1000000",
    "alloc": {
        "0x[ĐỊA_CHỈ_VALIDATOR_SẼ_ĐƯỢC_TẠO]": {
            "balance": "0x3635c9adc5dea00000"
        }
    },
    "coinbase": "0x[ĐỊA_CHỈ_VALIDATOR_SẼ_ĐƯỢC_TẠO]",
    "extraData": "0x48443234436861696e202d2050726f64756374696f6e",
    "nonce": "0x1234567890abcdef",
    "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp": "0x00"
}
```

### **Bước 3: Tạo File Password**

```powershell
# Tạo file password cho validator
"validator1pass" | Out-File -FilePath "C:\HD24Chain\validator1.pass" -Encoding ASCII
```

### **Bước 4: Tạo Tài Khoản Validator**

```powershell
# Tạo tài khoản validator
C:\HD24Chain\mychain.exe --datadir C:\HD24Chain\data account new --password C:\HD24Chain\validator1.pass
```

**Lưu ý**: Ghi lại địa chỉ validator được tạo (ví dụ: `0x0C3e517A6E924d33155281106F8B2d731914D066`)

### **Bước 5: Cập Nhật File Genesis**

Cập nhật file `C:\HD24Chain\genesis.json` với địa chỉ validator thực tế.

### **Bước 6: Khởi Tạo Blockchain**

```powershell
# Xóa dữ liệu cũ nếu có
Remove-Item -Recurse -Force "C:\HD24Chain\data\*" -ErrorAction SilentlyContinue

# Khởi tạo blockchain với genesis mới
C:\HD24Chain\mychain.exe --datadir C:\HD24Chain\data init C:\HD24Chain\genesis.json
```

### **Bước 7: Tạo Lại Tài Khoản Validator**

```powershell
# Tạo lại tài khoản validator sau khi khởi tạo blockchain
C:\HD24Chain\mychain.exe --datadir C:\HD24Chain\data account new --password C:\HD24Chain\validator1.pass
```

### **Bước 8: Chạy Node**

#### Cách 1: Sử dụng Script (Khuyến nghị)
```powershell
# Chạy script khởi động node
C:\HD24Chain\start_node.ps1
```

#### Cách 2: Lệnh Trực Tiếp
```powershell
# Chạy node với validator đã unlock
C:\HD24Chain\mychain.exe --datadir C:\HD24Chain\data --mine --unlock 0x[ĐỊA_CHỈ_VALIDATOR] --password C:\HD24Chain\validator1.pass --miner.etherbase=0x[ĐỊA_CHỈ_VALIDATOR] --http --http.addr=0.0.0.0 --http.port=8545 --http.api=eth,net,web3,personal,miner,admin,debug --allow-insecure-unlock --networkid=99999 --verbosity=3
```

##  Kiểm Tra và Monitoring

### **Kiểm Tra Node**

#### Sử dụng Script Test
```powershell
# Kiểm tra node có hoạt động không
C:\HD24Chain\test_node.ps1
```

#### Kiểm Tra Thủ Công

**Kiểm tra Block Number:**
```powershell
$body = @{
    jsonrpc = "2.0"
    method = "eth_blockNumber"
    params = @()
    id = 1
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8545" -Method POST -ContentType "application/json" -Body $body
```

**Kiểm tra Network ID:**
```powershell
$body = @{
    jsonrpc = "2.0"
    method = "net_version"
    params = @()
    id = 1
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8545" -Method POST -ContentType "application/json" -Body $body
```

##  Cấu Hình Nâng Cao

### **Cấu Hình Genesis**

```json
{
  "config": {
        "chainId": 99999,
        "pova": {
            "validators": [
                "0x[ĐỊA_CHỈ_VALIDATOR_1]",
                "0x[ĐỊA_CHỈ_VALIDATOR_2]",
                "0x[ĐỊA_CHỈ_VALIDATOR_3]"
            ],
            "period": 20
        }
  }
}
```

### **Tùy Chọn Dòng Lệnh**

#### Node Cơ Bản
```powershell
C:\HD24Chain\mychain.exe --datadir C:\HD24Chain\data --mine --http --http.addr=0.0.0.0 --http.port=8545 --networkid=99999
```

#### Với Validators Đã Unlock
```powershell
C:\HD24Chain\mychain.exe --datadir C:\HD24Chain\data --mine --unlock 0x[ĐỊA_CHỈ1],0x[ĐỊA_CHỈ2] --password C:\HD24Chain\validator1.pass --password C:\HD24Chain\validator2.pass --miner.etherbase=0x[ĐỊA_CHỈ1] --http --http.addr=0.0.0.0 --http.port=8545 --http.api=eth,net,web3,personal,miner,admin,debug --allow-insecure-unlock --networkid=99999 --verbosity=3
```

#### Với WebSocket
```powershell
C:\HD24Chain\mychain.exe --datadir C:\HD24Chain\data --mine --unlock 0x[ĐỊA_CHỈ_VALIDATOR] --password C:\HD24Chain\validator1.pass --miner.etherbase=0x[ĐỊA_CHỈ_VALIDATOR] --http --http.addr=0.0.0.0 --http.port=8545 --http.api=eth,net,web3,personal,miner,admin,debug --ws --ws.addr=0.0.0.0 --ws.port=8546 --ws.api=eth,net,web3,personal,miner,admin,debug --allow-insecure-unlock --networkid=99999 --verbosity=3
```

##  Benchmark Hiệu Suất

- **Thông Lượng Giao Dịch**: 150+ TPS (chuyển khoản đơn giản)
- **Gọi Smart Contract**: 100+ TPS
- **Sử Dụng Bộ Nhớ**: 2-4 GB (thông thường)
- **Lưu Trữ**: 10-50 GB (tùy thuộc vào sử dụng)
- **Phản Hồi API**: <100ms (thông thường)
- **Thời Gian Block**: 20 giây (có thể cấu hình)

##  Tính Năng Bảo Mật

- **Tích Hợp Ví Phần Cứng**: Hỗ trợ Ledger, Trezor
- **Ví Đa Chữ Ký**: Bảo mật nâng cao
- **Cấu Hình Thủ Công**: Kiểm soát hoàn toàn cài đặt bảo mật
- **Lưu Trữ Mã Hóa**: Mã hóa cơ sở dữ liệu
- **Ghi Log Kiểm Toán**: Dấu vết giao dịch hoàn chỉnh
- **Xác Thực Validator**: Hệ thống xác thực validator nghiêm ngặt

##  Troubleshooting

### **Node Không Khởi Động**
1. Kiểm tra file `mychain.exe` có tồn tại không
2. Kiểm tra thư mục `data` có được khởi tạo không
3. Kiểm tra file `genesis.json` có đúng format không
4. Kiểm tra quyền truy cập thư mục

### **Không Kết Nối Được API**
1. Kiểm tra node có đang chạy không
2. Kiểm tra port 8545 có bị block không
3. Kiểm tra firewall settings
4. Kiểm tra địa chỉ IP binding

### **Lỗi Validator**
1. Kiểm tra địa chỉ validator có đúng không
2. Kiểm tra file password có tồn tại không
3. Kiểm tra validator có được unlock không
4. Kiểm tra cấu hình genesis có đúng không

##  Thông Tin Liên Hệ
- **GitHub**: https://github.com/haidang24 


##  Tài Liệu Thêm

- **[POVA Algorithm](POVA_ALGORITHM.md)** - Chi tiết thuật toán POVA
- **[Deployment Guide](DEPLOYMENT.md)** - Hướng dẫn triển khai 

---

**Sẵn sàng triển khai blockchain doanh nghiệp của bạn? Liên hệ đội ngũ bán hàng ngay hôm nay!**
* Gmail: haidangattt@gmail.com
