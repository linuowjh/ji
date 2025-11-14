// app.js
App({
  globalData: {
    userInfo: null,
    token: null,
    apiBase: 'http://localhost:8080' // 替换为实际API地址
  },

  onLaunch() {
    // 检查登录状态
    const token = wx.getStorageSync('token')
    if (token) {
      this.globalData.token = token
      this.getUserInfo()
    }
  },

  // 微信登录
  login() {
    return new Promise((resolve, reject) => {
      wx.login({
        success: res => {
          if (res.code) {
            // 发送 res.code 到后台换取 token
            wx.request({
              url: `${this.globalData.apiBase}/api/v1/auth/wechat-login`,
              method: 'POST',
              data: { code: res.code },
              success: response => {
                if (response.data.code === 0) {
                  const { token, userInfo } = response.data.data
                  this.globalData.token = token
                  this.globalData.userInfo = userInfo
                  wx.setStorageSync('token', token)
                  resolve(userInfo)
                } else {
                  reject(response.data.message)
                }
              },
              fail: reject
            })
          } else {
            reject('登录失败：' + res.errMsg)
          }
        },
        fail: reject
      })
    })
  },

  // 获取用户信息
  getUserInfo() {
    return new Promise((resolve, reject) => {
      wx.request({
        url: `${this.globalData.apiBase}/api/v1/users/profile`,
        method: 'GET',
        header: {
          'Authorization': `Bearer ${this.globalData.token}`
        },
        success: res => {
          if (res.data.code === 0) {
            this.globalData.userInfo = res.data.data
            resolve(res.data.data)
          } else {
            reject(res.data.message)
          }
        },
        fail: reject
      })
    })
  },

  // 退出登录
  logout() {
    this.globalData.token = null
    this.globalData.userInfo = null
    wx.removeStorageSync('token')
    wx.reLaunch({
      url: '/pages/index/index'
    })
  }
})
