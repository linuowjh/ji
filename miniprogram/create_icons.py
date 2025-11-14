#!/usr/bin/env python3
"""
创建TabBar占位图标
使用PIL/Pillow库创建简单的纯色占位图标
"""

try:
    from PIL import Image, ImageDraw, ImageFont
except ImportError:
    print("❌ 错误: 未找到 Pillow 库")
    print("请安装: pip3 install Pillow")
    exit(1)

import os

# 配置
SIZE = 81
GRAY = "#7A7E83"
GREEN = "#3cc51f"
OUTPUT_DIR = "miniprogram/images"

# TabBar图标配置
tabbar_icons = [
    ("home", "首页"),
    ("memorial", "纪念"),
    ("family", "家族"),
    ("profile", "我的")
]

# 页面图标配置（只需要一个版本，不需要active状态）
page_icons = [
    ("create", "创建", GREEN),
    ("list", "列表", GREEN),
    ("family-icon", "家族", GREEN),
    ("guide", "引导", GREEN),
    ("empty", "空", GRAY),
    ("empty-family", "空", GRAY),
    ("default-avatar", "头像", GRAY)
]

def create_icon(filename, color, text=""):
    """创建一个简单的图标"""
    # 创建图像
    img = Image.new('RGB', (SIZE, SIZE), color)
    draw = ImageDraw.Draw(img)
    
    # 如果有文字，添加文字（可选）
    if text:
        try:
            # 尝试使用系统字体
            font = ImageFont.truetype("/System/Library/Fonts/PingFang.ttc", 20)
        except:
            # 如果找不到字体，使用默认字体
            font = ImageFont.load_default()
        
        # 计算文字位置（居中）
        bbox = draw.textbbox((0, 0), text, font=font)
        text_width = bbox[2] - bbox[0]
        text_height = bbox[3] - bbox[1]
        x = (SIZE - text_width) // 2
        y = (SIZE - text_height) // 2
        
        # 绘制文字
        draw.text((x, y), text, fill="white", font=font)
    
    # 保存图像
    img.save(filename)
    print(f"✅ 创建: {filename}")

def main():
    print("🎨 开始创建小程序图标资源...")
    
    # 创建输出目录
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    
    print("\n📱 创建TabBar图标...")
    # 创建TabBar图标（需要两个状态）
    for icon_name, text in tabbar_icons:
        # 未选中状态（灰色）
        gray_file = os.path.join(OUTPUT_DIR, f"{icon_name}.png")
        create_icon(gray_file, GRAY, text)
        
        # 选中状态（绿色）
        active_file = os.path.join(OUTPUT_DIR, f"{icon_name}-active.png")
        create_icon(active_file, GREEN, text)
    
    print("\n🖼️  创建页面图标...")
    # 创建页面图标（只需要一个版本）
    for icon_name, text, color in page_icons:
        icon_file = os.path.join(OUTPUT_DIR, f"{icon_name}.png")
        create_icon(icon_file, color, text)
    
    print("\n✅ 所有图标创建完成！")
    print(f"📁 图标位置: {OUTPUT_DIR}/")
    print(f"📊 统计: TabBar图标 8个 + 页面图标 7个 = 共 15个")
    print("\n⚠️  注意: 这些是简单的占位图标，建议后续替换为实际设计的图标")
    print("\n📝 下一步:")
    print("1. 在微信开发者工具中重新编译")
    print("2. 所有图片资源应该可以正常加载了")
    print("3. 后续可以替换为更精美的图标设计")

if __name__ == "__main__":
    main()
