# Docker操作脚本 (PowerShell版本)
# 使用方法: .\start.ps1 [命令]
# 可用命令: build, up, down, restart, logs, status, clean

param(
    [string]$Command = "help"
)

# 颜色定义
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"
$Blue = "Blue"
$White = "White"

# 打印带颜色的消息
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Red
}

function Write-Header {
    param([string]$Message)
    Write-Host "=================================" -ForegroundColor $Blue
    Write-Host "  Go Web后端项目 Docker 管理" -ForegroundColor $Blue
    Write-Host "=================================" -ForegroundColor $Blue
}

# 检查Docker是否运行
function Test-Docker {
    try {
        docker info | Out-Null
        return $true
    }
    catch {
        return $false
    }
}

# 检查docker-compose是否可用
function Test-DockerCompose {
    try {
        docker-compose version | Out-Null
        return $true
    }
    catch {
        try {
            docker compose version | Out-Null
            return $true
        }
        catch {
            return $false
        }
    }
}

# 构建镜像
function Build-Image {
    Write-Info "构建Docker镜像..."
    docker-compose build --no-cache
    Write-Info "镜像构建完成！"
}

# 启动服务
function Start-Services {
    Write-Info "启动所有服务..."
    docker-compose up -d
    Write-Info "服务启动完成！"
    Write-Info "应用将在 http://localhost:8080 上运行"
    Write-Info "MySQL端口: 3306"
    Write-Info "Redis端口: 6379"
}

# 停止服务
function Stop-Services {
    Write-Info "停止所有服务..."
    docker-compose down
    Write-Info "服务已停止！"
}

# 重启服务
function Restart-Services {
    Write-Info "重启所有服务..."
    docker-compose restart
    Write-Info "服务重启完成！"
}

# 查看日志
function Show-Logs {
    Write-Info "查看服务日志..."
    docker-compose logs -f
}

# 查看服务状态
function Show-Status {
    Write-Info "查看服务状态..."
    docker-compose ps
    Write-Host ""
    Write-Info "查看服务健康状态..."
    docker-compose ps --format "table {{.Name}}`t{{.Status}}`t{{.Ports}}"
}

# 清理资源
function Clean-Resources {
    Write-Warning "这将删除所有容器、镜像和数据卷，确定继续吗？(y/N)"
    $response = Read-Host
    if ($response -match "^[yY]$|^[yY][eE][sS]$") {
        Write-Info "清理所有Docker资源..."
        docker-compose down -v --rmi all
        docker system prune -f
        Write-Info "清理完成！"
    }
    else {
        Write-Info "取消清理操作"
    }
}

# 进入应用容器
function Enter-Shell {
    Write-Info "进入应用容器..."
    docker-compose exec app sh
}

# 查看应用日志
function Show-AppLogs {
    Write-Info "查看应用日志..."
    docker-compose logs -f app
}

# 查看数据库日志
function Show-DbLogs {
    Write-Info "查看数据库日志..."
    docker-compose logs -f mysql
}

# 查看Redis日志
function Show-RedisLogs {
    Write-Info "查看Redis日志..."
    docker-compose logs -f redis
}

# 显示帮助信息
function Show-Help {
    Write-Header
    Write-Host "可用命令:" -ForegroundColor $White
    Write-Host "  build      - 构建Docker镜像"
    Write-Host "  up         - 启动所有服务"
    Write-Host "  down       - 停止所有服务"
    Write-Host "  restart    - 重启所有服务"
    Write-Host "  logs       - 查看所有服务日志"
    Write-Host "  status     - 查看服务状态"
    Write-Host "  shell      - 进入应用容器"
    Write-Host "  app_logs   - 查看应用日志"
    Write-Host "  db_logs    - 查看数据库日志"
    Write-Host "  redis_logs - 查看Redis日志"
    Write-Host "  clean      - 清理所有Docker资源"
    Write-Host "  help       - 显示此帮助信息"
    Write-Host ""
    Write-Host "示例:" -ForegroundColor $White
    Write-Host "  .\start.ps1 build    # 构建镜像"
    Write-Host "  .\start.ps1 up       # 启动服务"
    Write-Host "  .\start.ps1 status   # 查看状态"
}

# 主函数
function Main {
    Write-Header
    
    # 检查Docker环境
    if (-not (Test-Docker)) {
        Write-Error "Docker未运行，请先启动Docker"
        exit 1
    }
    
    if (-not (Test-DockerCompose)) {
        Write-Error "docker-compose未安装或不可用"
        exit 1
    }
    
    switch ($Command) {
        "build" { Build-Image }
        "up" { Start-Services }
        "down" { Stop-Services }
        "restart" { Restart-Services }
        "logs" { Show-Logs }
        "status" { Show-Status }
        "shell" { Enter-Shell }
        "app_logs" { Show-AppLogs }
        "db_logs" { Show-DbLogs }
        "redis_logs" { Show-RedisLogs }
        "clean" { Clean-Resources }
        default { Show-Help }
    }
}

# 执行主函数
Main
