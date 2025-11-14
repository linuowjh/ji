// pages/memorial/detail/detail.js
const app = getApp()

Page({
  data: {
    memorialId: '',
    memorial: null,
    recentWorship: [],
    loading: true
  },

  onLoad(options) {
    if (options.id) {
      this.setData({ memorialId: options.id })
      this.loadMemorialDetail()
      this.loadRecentWorship()
    }
  },

  // 加载纪念馆详情
  loadMemorialDetail() {
    wx.showLoading({ title: '加载中...' })
    
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/memorials/${this.data.memorialId}`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        wx.hideLoading()
        if (res.data.code === 0) {
          this.setData({
            memorial: res.data.data,
            loading: false
          })
        } else {
          wx.showToast({
            title: res.data.message || '加载失败',
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

  // 加载最近祭扫记录
  loadRecentWorship() {
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/worship/records`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      data: {
        memorialId: this.data.memorialId,
        page: 1,
        pageSize: 5
      },
      success: res => {
        if (res.data.code === 0) {
          this.setData({
            recentWorship: res.data.data || []
          })
        }
      }
    })
  },

  // 前往祭扫页面
  goToWorship() {
    wx.navigateTo({
      url: `/pages/memorial/worship/worship?id=${this.data.memorialId}`
    })
  },

  // 编辑纪念馆
  editMemorial() {
    wx.navigateTo({
      url: `/pages/memorial/create/create?id=${this.data.memorialId}&mode=edit`
    })
  },

  // 分享纪念馆
  onShareAppMessage() {
    return {
      title: `${this.data.memorial.deceasedName}的纪念馆`,
      path: `/pages/memorial/detail/detail?id=${this.data.memorialId}`,
      imageUrl: this.data.memorial.avatarUrl
    }
  },

  // 预览图片
  previewImage(e) {
    const url = e.currentTarget.dataset.url
    wx.previewImage({
      current: url,
      urls: [url]
    })
  }
})
