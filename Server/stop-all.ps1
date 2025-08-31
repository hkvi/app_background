# 停止所有服务的PowerShell脚本

Write-Host "正在停止所有服务..." -ForegroundColor Yellow

# 停止Docker Compose服务
Write-Host "正在停止Docker服务..." -ForegroundColor Blue
docker-compose down

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 所有服务已成功停止!" -ForegroundColor Green
} else {
    Write-Host "⚠️  停止服务时出现问题" -ForegroundColor Yellow
}

# 清理悬空镜像（可选）
$cleanup = Read-Host "是否清理Docker悬空镜像? (y/n)"
if ($cleanup -eq "y" -or $cleanup -eq "Y") {
    Write-Host "正在清理悬空镜像..." -ForegroundColor Blue
    docker image prune -f
    Write-Host "清理完成!" -ForegroundColor Green
}
