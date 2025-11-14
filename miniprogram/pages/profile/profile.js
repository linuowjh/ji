// pages/profile/profile.js
const app = getApp()

Page({
  data: {
    userInfo: null,
    stats: {
      memorialCount: 0,
      worshipCount: 0,
      familyCount: 0
    }
  },

  onLoad() {
    this.loadUserInfo()
    this.loadStats()
  },

  onShow() {
    this.loadUserInfo()
    this.loadStats()
  },

  // 加载用户信息
  loadUserInfo() {
    if (app.globalData.userInfo) {
      this.setData({ userInfo: app.globalData.userInfo })
    } else {
      app.getUserInfo()
        .then(userInfo => {
          this.setData({ userInfo })
        })
        .catch(() => {
          wx.showToast({
            title: '获取用户信息失败',
            icon: 'none'
          })
        })
    }
  },

  // 加载统计数据
  loadStats() {
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/users/stats`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        if (res.data.code === 0) {
          this.setData({ stats: res.data.data })
        }
      }
    })
  },

  // 前往我的纪念馆
  goToMyMemorials() {
    wx.switchTab({
      url: '/pages/memorial/list/list'
    })
  },

  // 前往祭扫记录
  goToWorshipHistory() {
    wx.navigateTo({
      url: '/pages/profile/worship-history/worship-history'
    })
  },

  // 前往家族圈
  goToFamilies() {
    wx.switchTab({
      url: '/pages/family/list/list'
    })
  },

  // 前往隐私设置
  goToPrivacySettings() {
    wx.navigateTo({
      url: '/pages/profile/privacy/privacy'
    })
  },

  // 前往帮助中心
  goToHelp() {
    wx.navigateTo({
      url: '/pages/profile/help/help'
    })
  },

  // 前往关于我们
  goToAbout() {
    wx.navigateTo({
      url: '/pages/profile/about/about'
    })
  },

  // 联系客服
  contactService() {
    wx.makePhoneCall({
      phoneNumber: '400-123-4567'
    })
  },

  // 退出登录
  logout() {
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: res => {
        if (res.confirm) {
          app.logout()
        }
      }
    })
  }
})
