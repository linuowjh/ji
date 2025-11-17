// pages/memorial/worship/worship.js
const app = getApp()

Page({
  data: {
    memorialId: '',
    memorial: null,
    activeTab: 'flower',
    flowerType: 'chrysanthemum',
    candleType: 'red',
    candleDuration: 60,
    incenseCount: 3,
    incenseType: 'sandalwood',
    tributeType: 'fruit',
    prayerText: ''
  },

  onLoad(options) {
    if (options.id) {
      this.setData({ memorialId: options.id })
      this.loadMemorial()
    }
  },

  // 加载纪念馆信息
  loadMemorial() {
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/memorials/${this.data.memorialId}`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        if (res.data.code === 0) {
          this.setData({ memorial: res.data.data })
        }
      }
    })
  },

  // 切换标签
  switchTab(e) {
    const tab = e.currentTarget.dataset.tab
    this.setData({ activeTab: tab })
  },

  // 选择花卉类型
  selectFlower(e) {
    this.setData({ flowerType: e.currentTarget.dataset.type })
  },

  // 选择香柱数量
  selectIncense(e) {
    this.setData({ incenseCount: e.currentTarget.dataset.count })
  },

  // 选择供品类型
  selectTribute(e) {
    this.setData({ tributeType: e.currentTarget.dataset.type })
  },

  // 输入祈福语
  inputPrayer(e) {
    this.setData({ prayerText: e.detail.value })
  },

  // 献花
  offerFlower() {
    this.performWorship('flower', {
      flowerType: this.data.flowerType,
      quantity: 1  // 默认数量为1
    })
  },

  // 点烛
  lightCandle() {
    this.performWorship('candle', {
      candleType: this.data.candleType,
      duration: this.data.candleDuration
    })
  },

  // 上香
  offerIncense() {
    this.performWorship('incense', {
      incenseCount: this.data.incenseCount,
      incenseType: this.data.incenseType
    })
  },

  // 供奉供品
  offerTribute() {
    this.performWorship('tribute', {
      tributeType: this.data.tributeType,
      items: ['水果']  // 默认供品项目
    })
  },

  // 祈福
  sendPrayer() {
    if (!this.data.prayerText.trim()) {
      wx.showToast({
        title: '请输入祈福语',
        icon: 'none'
      })
      return
    }

    this.performWorship('prayer', {
      content: this.data.prayerText
    })
  },

  // 执行祭扫操作
  performWorship(type, data) {
    wx.showLoading({ title: '提交中...' })

    // 构建正确的URL
    const typeMap = {
      'flower': 'flowers',
      'candle': 'candles',
      'incense': 'incense',
      'tribute': 'tributes',
      'prayer': 'prayers'
    }
    const endpoint = typeMap[type] || type

    wx.request({
      url: `${app.globalData.apiBase}/api/v1/worship/memorials/${this.data.memorialId}/${endpoint}`,
      method: 'POST',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`,
        'Content-Type': 'application/json'
      },
      data: data,
      success: res => {
        wx.hideLoading()
        if (res.data.code === 0) {
          wx.showToast({
            title: '祭扫成功',
            icon: 'success'
          })
          
          // 播放音效
          this.playSound(type)
          
          // 显示动画
          this.showAnimation(type)
          
          // 清空祈福语
          if (type === 'prayer') {
            this.setData({ prayerText: '' })
          }
          
          // 刷新纪念馆数据（更新祭扫次数）
          this.loadMemorial()
        } else {
          wx.showToast({
            title: res.data.message || '操作失败',
            icon: 'none'
          })
        }
      },
      fail: () => {
        wx.hideLoading()
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        })
      }
    })
  },

  // 播放音效
  playSound(type) {
    const soundMap = {
      flower: '/sounds/flower.mp3',
      candle: '/sounds/candle.mp3',
      incense: '/sounds/incense.mp3',
      tribute: '/sounds/tribute.mp3',
      prayer: '/sounds/prayer.mp3'
    }
    
    const innerAudioContext = wx.createInnerAudioContext()
    innerAudioContext.src = soundMap[type]
    innerAudioContext.play()
  },

  // 显示动画
  showAnimation(type) {
    // 这里可以添加动画效果
    console.log(`Show animation for ${type}`)
  },

  // 录制语音留言
  recordVoice() {
    wx.navigateTo({
      url: `/pages/memorial/message/message?id=${this.data.memorialId}&type=voice`
    })
  },

  // 录制视频留言
  recordVideo() {
    wx.navigateTo({
      url: `/pages/memorial/message/message?id=${this.data.memorialId}&type=video`
    })
  }
})
