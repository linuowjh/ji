#!/bin/bash

# 创建TabBar占位图标的脚本
# 需要安装 ImageMagick: brew install imagemagick (macOS)

echo "🎨 开始创建TabBar占位图标..."

# 检查 ImageMagick 是否安装
if ! command -v convert &> /dev/null; then
    echo "❌ 错误: 未找到 ImageMagick"
    echo "请先安装 ImageMagick:"
    echo "  macOS: brew install imagemagick"
    echo "  Ubuntu: sudo apt-get install imagemagick"
    exit 1
fi

# 创建 images 目录（如果不存在）
mkdir -p miniprogram/images

# 定义颜色
GRAY="#7A7E83"
GREEN="#3cc51f"

# 创建灰色占位图标（未选中状态）
echo "📝 创建未选中状态图标..."
convert -size 81x81 xc:$GRAY miniprogram/images/home.png
convert -size 81x81 xc:$GRAY miniprogram/images/memorial.png
convert -size 81x81 xc:$GRAY miniprogram/images/family.png
convert -size 81x81 xc:$GRAY miniprogram/images/profile.png

# 创建绿色占位图标（选中状态）
echo "📝 创建选中状态图标..."
convert -size 81x81 xc:$GREEN miniprogram/images/home-active.png
convert -size 81x81 xc:$GREEN miniprogram/images/memorial-active.png
convert -size 81x81 xc:$GREEN miniprogram/images/family-active.png
convert -size 81x81 xc:$GREEN miniprogram/images/profile-active.png

echo "✅ 占位图标创建完成！"
echo ""
echo "📁 图标位置: miniprogram/images/"
echo "📋 已创建的文件:"
ls -lh miniprogram/images/*.png
echo ""
echo "⚠️  注意: 这些是纯色占位图标，建议后续替换为实际设计的图标"
