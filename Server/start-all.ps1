# 启动所有服务的PowerShell脚本

Write-Host "正在启动微服务架构..." -ForegroundColor Green

# 检查Docker是否运行
$dockerRunning = docker info 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "错误: Docker未运行，请先启动Docker Desktop" -ForegroundColor Red
    exit 1
}

# 检查环境变量文件
if (-not (Test-Path ".env")) {
    if (Test-Path "env.example") {
        Write-Host "正在复制环境变量示例文件..." -ForegroundColor Yellow
        Copy-Item "env.example" ".env"
        Write-Host "请编辑 .env 文件配置您的环境变量" -ForegroundColor Yellow
    } else {
        Write-Host "警告: 未找到环境变量配置文件" -ForegroundColor Yellow
    }
}

# 构建并启动所有服务
Write-Host "正在构建并启动服务..." -ForegroundColor Blue
docker-compose up --build -d

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✅ 所有服务已成功启动!" -ForegroundColor Green
    Write-Host "`n📋 服务信息:" -ForegroundColor Cyan
    Write-Host "  🌐 网关服务: http://localhost:8080" -ForegroundColor White
    Write-Host "  🏢 业务服务: http://localhost:8081" -ForegroundColor White
    Write-Host "  🗄️  MySQL数据库: localhost:3306" -ForegroundColor White
    Write-Host "  🔧 Redis缓存: localhost:6379" -ForegroundColor White
    
    Write-Host "`n📚 API文档:" -ForegroundColor Cyan
    Write-Host "  健康检查: GET http://localhost:8080/api/health" -ForegroundColor White
    Write-Host "  用户注册: POST http://localhost:8080/api/auth/register" -ForegroundColor White
    Write-Host "  用户登录: POST http://localhost:8080/api/auth/login" -ForegroundColor White
    Write-Host "  发送短信: POST http://localhost:8080/api/auth/sms/send" -ForegroundColor White
    Write-Host "  短信登录: POST http://localhost:8080/api/auth/sms/login" -ForegroundColor White
    Write-Host "  刷新令牌: POST http://localhost:8080/api/auth/refresh" -ForegroundColor White
    
    Write-Host "`n🔧 管理命令:" -ForegroundColor Cyan
    Write-Host "  查看日志: docker-compose logs -f" -ForegroundColor White
    Write-Host "  停止服务: docker-compose down" -ForegroundColor White
    Write-Host "  重启服务: docker-compose restart" -ForegroundColor White
} else {
    Write-Host "❌ 服务启动失败!" -ForegroundColor Red
    Write-Host "请检查错误信息并重试" -ForegroundColor Yellow
}
