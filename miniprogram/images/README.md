# TabBar 图标说明

## ✅ 已创建的图标文件

当前目录包含以下15个占位图标：

### TabBar图标（8个）
1. `home.png` + `home-active.png` - 首页图标
2. `memorial.png` + `memorial-active.png` - 纪念馆图标
3. `family.png` + `family-active.png` - 家族圈图标
4. `profile.png` + `profile-active.png` - 我的图标

### 页面图标（7个）
5. `create.png` - 创建纪念馆图标
6. `list.png` - 列表图标
7. `family-icon.png` - 家族圈图标
8. `guide.png` - 引导图标
9. `empty.png` - 空状态图标
10. `empty-family.png` - 家族圈空状态图标
11. `default-avatar.png` - 默认头像

## 📐 图标规格要求

- **尺寸**：81px × 81px
- **格式**：PNG
- **背景**：透明
- **颜色**：
  - 未选中状态：灰色 `#7A7E83`
  - 选中状态：绿色 `#3cc51f`

## 🎨 快速解决方案

### 方案1：使用在线图标生成器（推荐）

1. 访问 [iconfont](https://www.iconfont.cn/) 或 [iconpark](https://iconpark.oceanengine.com/)
2. 搜索并下载以下图标：
   - home / 首页
   - memorial / 纪念碑 / 墓碑
   - family / 家庭 / 家族
   - profile / 用户 / 我的
3. 调整尺寸为 81px × 81px
4. 导出为PNG格式
5. 分别保存为灰色和绿色版本

### 方案2：使用设计工具

使用 Figma / Sketch / Photoshop 创建：
1. 创建 81px × 81px 画布
2. 绘制简单的图标
3. 导出为PNG

### 方案3：临时占位图标

如果只是想快速测试，可以：
1. 创建纯色的 81px × 81px PNG图片
2. 命名为对应的文件名
3. 放在此目录下

## 🔧 创建占位图标的命令

在 macOS/Linux 上，可以使用 ImageMagick 快速创建：

```bash
# 安装 ImageMagick（如果还没有）
# macOS: brew install imagemagick
# Ubuntu: sudo apt-get install imagemagick

# 创建灰色占位图标
convert -size 81x81 xc:#7A7E83 miniprogram/images/home.png
convert -size 81x81 xc:#7A7E83 miniprogram/images/memorial.png
convert -size 81x81 xc:#7A7E83 miniprogram/images/family.png
convert -size 81x81 xc:#7A7E83 miniprogram/images/profile.png

# 创建绿色占位图标
convert -size 81x81 xc:#3cc51f miniprogram/images/home-active.png
convert -size 81x81 xc:#3cc51f miniprogram/images/memorial-active.png
convert -size 81x81 xc:#3cc51f miniprogram/images/family-active.png
convert -size 81x81 xc:#3cc51f miniprogram/images/profile-active.png
```

## 📝 图标设计建议

### 首页图标 (home)
- 可以使用：房子、首页、主页图标
- 风格：简洁、线条清晰

### 纪念馆图标 (memorial)
- 可以使用：纪念碑、墓碑、蜡烛、花朵图标
- 风格：庄重、简约

### 家族圈图标 (family)
- 可以使用：家庭、人群、树形图标
- 风格：温馨、亲切

### 我的图标 (profile)
- 可以使用：用户、个人、头像图标
- 风格：简单、通用

## ✅ 添加图标后

将图标文件放入此目录后：
1. 重新编译小程序
2. TabBar应该正常显示
3. 错误提示消失

## 🎯 推荐资源

- [iconfont 阿里巴巴矢量图标库](https://www.iconfont.cn/)
- [IconPark 字节跳动图标库](https://iconpark.oceanengine.com/)
- [Flaticon 免费图标](https://www.flaticon.com/)
- [Icons8 图标资源](https://icons8.com/)
