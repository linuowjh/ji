# 云念小程序

「云念」是一个网上扫墓祭奠小程序，为用户提供突破时空限制的线上祭奠服务。

## 🚀 快速开始

### ✅ 可以直接在微信开发者工具中打开！

**步骤**：
1. 下载并安装[微信开发者工具](https://developers.weixin.qq.com/miniprogram/dev/devtools/download.html)
2. 打开工具，点击"导入项目"
3. 选择 `miniprogram` 文件夹
4. AppID选择"测试号"（或填写你的AppID）
5. 点击"导入"即可开始开发

**项目已包含所有必需文件**：
- ✅ `project.config.json` - 项目配置
- ✅ `app.json` - 小程序配置  
- ✅ `app.js` - 入口文件
- ✅ 所有8个页面完整
- ✅ 所有4个组件完整

⚠️ **需要补充的资源**（不影响打开和运行）：
- TabBar图标（8个PNG文件，放在 `images/` 目录）
- 音效文件（可选，放在 `sounds/` 目录）

详细说明请查看：[开发指南.md](./开发指南.md)

## 项目结构

```
miniprogram/
├── pages/              # 页面目录
│   ├── index/         # 首页
│   ├── memorial/      # 纪念馆相关页面
│   │   ├── list/     # 纪念馆列表
│   │   ├── detail/   # 纪念馆详情
│   │   ├── create/   # 创建/编辑纪念馆
│   │   └── worship/  # 祭扫页面
│   ├── family/        # 家族圈相关页面
│   │   ├── list/     # 家族圈列表
│   │   └── detail/   # 家族圈详情
│   └── profile/       # 个人中心
├── components/        # 组件目录
│   ├── worship-panel/    # 祭扫操作面板
│   ├── media-uploader/   # 媒体上传组件
│   ├── prayer-wall/      # 祈福墙组件
│   └── time-capsule/     # 时光信箱组件
├── utils/             # 工具函数
│   ├── api.js        # API请求封装
│   ├── animation.js  # 动画工具
│   ├── performance.js # 性能优化工具
│   └── sound.js      # 音效管理
├── images/            # 图片资源（需要添加）
├── sounds/            # 音效资源（需要添加）
├── app.js            # 小程序入口
├── app.json          # 小程序配置
├── app.wxss          # 全局样式
└── sitemap.json      # 站点地图配置
```

## 功能特性

### 核心功能
- ✅ 用户登录认证（微信登录）
- ✅ 纪念馆创建和管理
- ✅ 线上祭扫（献花、点烛、上香、供品、祈福）
- ✅ 家族圈功能
- ✅ 个人中心管理

### 交互组件
- ✅ 祭扫操作面板（支持多种祭扫方式）
- ✅ 媒体上传组件（图片、视频、语音）
- ✅ 祈福墙（展示和点赞祈福留言）
- ✅ 时光信箱（语音和视频留言）

### 用户体验优化
- ✅ 流畅的动画效果（淡入淡出、缩放、滑动等）
- ✅ 庄重的音效设计
- ✅ 性能优化（图片懒加载、缓存管理、防抖节流）
- ✅ 响应式设计

## 开发指南

### 环境要求
- 微信开发者工具
- Node.js 14+
- 后端API服务

### 配置说明

1. **修改API地址**
   在 `app.js` 中修改 `apiBase` 为实际的后端API地址：
   ```javascript
   globalData: {
     apiBase: 'https://your-api-domain.com'
   }
   ```

2. **添加图片资源**
   在 `images/` 目录下添加以下图片：
   - home.png / home-active.png（首页图标）
   - memorial.png / memorial-active.png（纪念馆图标）
   - family.png / family-active.png（家族圈图标）
   - profile.png / profile-active.png（个人中心图标）
   - create.png（创建图标）
   - list.png（列表图标）
   - family-icon.png（家族图标）
   - guide.png（引导图）
   - empty.png（空状态图）
   - empty-family.png（家族圈空状态图）
   - default-avatar.png（默认头像）

3. **添加音效资源**
   在 `sounds/` 目录下添加以下音效文件：
   - flower.mp3（献花音效）
   - candle.mp3（点烛音效）
   - incense.mp3（上香音效）
   - tribute.mp3（供品音效）
   - prayer.mp3（祈福音效）
   - bell.mp3（钟声音效）
   - success.mp3（成功音效）
   - click.mp3（点击音效）

### 开发流程

1. 使用微信开发者工具打开项目
2. 配置AppID（测试可使用测试号）
3. 启动后端API服务
4. 在开发者工具中编译运行

### 组件使用示例

#### 祭扫操作面板
```xml
<worship-panel 
  type="flower" 
  memorial-id="{{memorialId}}"
  bind:submit="handleWorshipSubmit"
  bind:recordVoice="handleRecordVoice"
  bind:recordVideo="handleRecordVideo">
</worship-panel>
```

#### 媒体上传组件
```xml
<media-uploader 
  type="image" 
  max-count="9"
  bind:upload="handleUpload"
  bind:delete="handleDelete">
</media-uploader>
```

#### 祈福墙组件
```xml
<prayer-wall memorial-id="{{memorialId}}"></prayer-wall>
```

#### 时光信箱组件
```xml
<time-capsule memorial-id="{{memorialId}}"></time-capsule>
```

## 设计理念

### 视觉设计
- 采用庄重、素雅的配色方案
- 使用渐变背景营造氛围感
- 圆角卡片设计提升亲和力
- 合理的留白和间距

### 交互设计
- 流畅的页面切换动画
- 即时的操作反馈
- 清晰的视觉层级
- 符合直觉的操作流程

### 性能优化
- 图片懒加载减少首屏加载时间
- 数据缓存减少网络请求
- 防抖节流优化用户操作
- 组件化开发提高代码复用

## API接口

详细的API接口文档请参考后端项目的 `docs/` 目录。

主要接口包括：
- 用户认证：`/api/v1/auth/*`
- 纪念馆管理：`/api/v1/memorials/*`
- 祭扫功能：`/api/v1/worship/*`
- 家族圈：`/api/v1/families/*`
- 媒体上传：`/api/v1/media/*`

## 注意事项

1. **隐私保护**
   - 严格遵守微信小程序隐私政策
   - 用户数据加密传输和存储
   - 提供隐私设置选项

2. **内容审核**
   - 所有用户上传内容需经过审核
   - 敏感词过滤
   - 违规内容处理机制

3. **性能监控**
   - 定期检查页面性能
   - 优化图片和视频资源
   - 监控API响应时间

4. **兼容性**
   - 支持微信基础库 2.0+
   - 适配不同屏幕尺寸
   - 处理网络异常情况

## 后续优化

- [ ] 添加更多纪念馆主题模板
- [ ] 实现线上追思会功能
- [ ] 增加家族谱系图展示
- [ ] 支持老照片AI修复
- [ ] 添加更多祭扫仪式选项
- [ ] 优化离线缓存策略
- [ ] 增加数据统计和分析

## 技术支持

如有问题，请联系技术支持团队。
