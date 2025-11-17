// pages/profile/profile.js
const app = getApp()

Page({
  data: {
    userInfo: null,
    avatarUrl: '/images/default-avatar.png',
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
      this.setData({ 
        userInfo: app.globalData.userInfo,
        avatarUrl: app.globalData.userInfo.avatarUrl || '/images/default-avatar.png'
      })
    } else if (app.globalData.token) {
      // 有 token 但没有用户信息，尝试获取
      app.getUserInfo()
        .then(userInfo => {
          this.setData({ 
            userInfo,
            avatarUrl: userInfo.avatarUrl || '/images/default-avatar.png'
          })
        })
        .catch(() => {
          // 获取失败，可能 token 过期，清除 token
          app.globalData.token = null
          wx.removeStorageSync('token')
          this.setData({ 
            userInfo: null,
            avatarUrl: '/images/default-avatar.png'
          })
        })
    } else {
      // 没有登录
      this.setData({ 
        userInfo: null,
        avatarUrl: '/images/default-avatar.png'
      })
    }
  },

  // 执行登录
  doLogin() {
    wx.showLoading({ title: '登录中...' })
    app.login()
      .then(userInfo => {
        wx.hideLoading()
        this.setData({ 
          userInfo,
          avatarUrl: userInfo.avatarUrl || '/images/default-avatar.png'
        })
        this.loadStats()
        wx.showToast({
          title: '登录成功',
          icon: 'success'
        })
      })
      .catch(err => {
        wx.hideLoading()
        wx.showToast({
          title: err || '登录失败',
          icon: 'none'
        })
      })
  },

  // 头像加载失败处理
  onAvatarError() {
    this.setData({
      avatarUrl: '/images/default-avatar.png'
    })
  },

  // 加载统计数据
  loadStats() {
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/users/statistics`,
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
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
    // TODO: 实现隐私设置页面
    // wx.navigateTo({
    //   url: '/pages/profile/privacy/privacy'
    // })
  },

  // 前往帮助中心
  goToHelp() {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
    // TODO: 实现帮助中心页面
    // wx.navigateTo({
    //   url: '/pages/profile/help/help'
    // })
  },

  // 前往关于我们
  goToAbout() {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
    // TODO: 实现关于我们页面
    // wx.navigateTo({
    //   url: '/pages/profile/about/about'
    // })
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
