@echo off
chcp 65001
set docker_username="pgman"
set /p docker_repo="请输入您要发布的Docker镜像仓库名称："
set /p docker_tag="请输入您要发布的Docker镜像标签名称（默认为latest）:"

if "%docker_tag%" == "" (
    set docker_tag=latest
)

echo 正在构建Docker镜像，请等待...
docker build -t %docker_username%/%docker_repo%:%docker_tag% .
if %errorlevel% neq 0 (
    goto :error
)

echo 正在推送Docker镜像，请等待...
docker push %docker_username%/%docker_repo%:%docker_tag%
if %errorlevel% neq 0 (
    goto :error
)

echo Docker镜像已成功推送到：%docker_username%/%docker_repo%:%docker_tag%
exit /b 0

:error
echo 发布Docker镜像时出现错误，请检查并重试。
exit /b 1