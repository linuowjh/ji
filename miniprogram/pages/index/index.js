// pages/index/index.js
const app = getApp()

Page({
  data: {
    recentMemorials: [],
    upcomingReminders: [],
    hasMemorial: false
  },

  onLoad() {
    this.checkLogin()
  },

  onShow() {
    if (app.globalData.token) {
      this.loadData()
    }
  },

  // 检查登录状态
  checkLogin() {
    if (!app.globalData.token) {
      wx.redirectTo({
        url: '/pages/login/login'
      })
    } else {
      this.loadData()
    }
  },

  // 加载数据
  loadData() {
    this.getRecentMemorials()
    this.getUpcomingReminders()
  },

  // 获取最近访问的纪念馆
  getRecentMemorials() {
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/memorials/recent`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        if (res.data.code === 0) {
          this.setData({
            recentMemorials: res.data.data || [],
            hasMemorial: res.data.data && res.data.data.length > 0
          })
        }
      }
    })
  },

  // 获取即将到来的纪念日
  getUpcomingReminders() {
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/users/reminders/upcoming`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        if (res.data.code === 0) {
          this.setData({
            upcomingReminders: res.data.data || []
          })
        }
      }
    })
  },

  // 创建纪念馆
  createMemorial() {
    if (!app.globalData.token) {
      this.checkLogin()
      return
    }
    wx.navigateTo({
      url: '/pages/memorial/create/create'
    })
  },

  // 前往纪念馆列表
  goToMemorialList() {
    wx.switchTab({
      url: '/pages/memorial/list/list'
    })
  },

  // 前往家族圈
  goToFamily() {
    wx.switchTab({
      url: '/pages/family/list/list'
    })
  },

  // 前往纪念馆详情
  goToMemorialDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/memorial/detail/detail?id=${id}`
    })
  }
})
