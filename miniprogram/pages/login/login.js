// pages/login/login.js
const app = getApp()

Page({
  data: {
    hasPhone: false,
    phoneCode: null,
    avatarUrl: '',
    nickname: ''
  },

  onLoad() {
    // 检查是否已登录
    if (app.globalData.token) {
      this.redirectToHome()
    }
  },

  // 选择头像
  onChooseAvatar(e) {
    console.log('选择头像:', e.detail.avatarUrl);
    this.setData({
      avatarUrl: e.detail.avatarUrl,
    });
    console.log('当前状态:', {
      avatarUrl: this.data.avatarUrl,
      nickname: this.data.nickname,
      hasPhone: this.data.hasPhone,
    });
  },

  // 输入昵称
  onNicknameInput(e) {
    this.setData({
      nickname: e.detail.value,
    });
    console.log('当前状态:', {
      avatarUrl: this.data.avatarUrl,
      nickname: this.data.nickname,
      hasPhone: this.data.hasPhone,
    });
  },

  // 获取手机号
  getPhoneNumber(e) {
    console.log('获取手机号回调:', e)
    if (e.detail.code) {
      this.setData({
        hasPhone: true,
        phoneCode: e.detail.code
      })
      wx.showToast({
        title: '手机号授权成功',
        icon: 'success'
      })
    } else {
      console.error('获取手机号失败:', e.detail.errMsg)
      
      // 开发环境模拟授权成功（仅用于测试）
      if (e.detail.errMsg.includes('no permission')) {
        wx.showModal({
          title: '开发环境提示',
          content: '当前为开发环境，手机号授权功能需要小程序认证。是否使用模拟数据继续测试？',
          confirmText: '继续测试',
          cancelText: '取消',
          success: res => {
            if (res.confirm) {
              // 使用模拟的手机号code
              this.setData({
                hasPhone: true,
                phoneCode: 'mock_phone_code_for_dev'
              })
              wx.showToast({
                title: '已启用测试模式',
                icon: 'success'
              })
            }
          }
        })
      } else {
        wx.showToast({
          title: '授权失败',
          icon: 'none'
        })
      }
    }
  },

  // 完成登录
  completeLogin() {
    if (!this.data.nickname) {
      wx.showToast({
        title: '请输入昵称',
        icon: 'none'
      })
      return
    }

    if (!this.data.avatarUrl) {
      wx.showToast({
        title: '请选择头像',
        icon: 'none'
      })
      return
    }

    if (!this.data.hasPhone) {
      wx.showToast({
        title: '请先授权手机号',
        icon: 'none'
      })
      return
    }

    wx.showLoading({ title: '登录中...' })

    // 构建用户信息
    const userInfo = {
      nickName: this.data.nickname,
      avatarUrl: this.data.avatarUrl
    }

    // 先进行基础登录
    app.login(userInfo)
      .then(user => {
        console.log('登录成功:', user)
        
        // 更新手机号（必填）
        return app.updatePhone(this.data.phoneCode)
          .then(() => {
            wx.hideLoading()
            wx.showToast({
              title: '登录成功',
              icon: 'success'
            })
            setTimeout(() => {
              this.redirectToHome()
            }, 1500)
          })
          .catch(err => {
            console.error('更新手机号失败:', err)
            wx.hideLoading()
            wx.showModal({
              title: '手机号授权失败',
              content: err || '请重新授权手机号',
              showCancel: false,
              success: () => {
                // 清除手机号状态，让用户重新授权
                this.setData({
                  hasPhone: false,
                  phoneCode: null
                })
              }
            })
          })
      })
      .catch(err => {
        wx.hideLoading()
        console.error('登录失败:', err)
        wx.showModal({
          title: '登录失败',
          content: err || '请稍后重试',
          showCancel: false
        })
      })
  },

  // 跳转到首页
  redirectToHome() {
    wx.switchTab({
      url: '/pages/index/index'
    })
  }
})
