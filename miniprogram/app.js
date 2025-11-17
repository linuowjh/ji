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

  // 微信登录（基础登录，不包含手机号）
  login(userInfo) {
    return new Promise((resolve, reject) => {
      wx.login({
        success: res => {
          if (res.code) {
            // 发送 res.code 和用户信息到后台换取 token
            wx.request({
              url: `${this.globalData.apiBase}/api/v1/auth/wechat-login`,
              method: 'POST',
              data: { 
                code: res.code,
                nickname: userInfo?.nickName || '微信用户',
                avatar: userInfo?.avatarUrl || ''
              },
              success: response => {
                if (response.data.code === 0) {
                  const { token, user } = response.data.data
                  this.globalData.token = token
                  this.globalData.userInfo = user
                  wx.setStorageSync('token', token)
                  resolve(user)
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

  // 更新手机号
  updatePhone(phoneCode) {
    return new Promise((resolve, reject) => {
      if (!this.globalData.token) {
        reject('请先登录')
        return
      }

      wx.request({
        url: `${this.globalData.apiBase}/api/v1/users/phone`,
        method: 'POST',
        header: {
          'Authorization': `Bearer ${this.globalData.token}`
        },
        data: {
          code: phoneCode
        },
        success: res => {
          if (res.data.code === 0) {
            // 更新本地用户信息
            if (this.globalData.userInfo) {
              this.globalData.userInfo.phone = res.data.data.phone
            }
            resolve(res.data.data)
          } else {
            reject(res.data.message)
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
