// pages/family/create/create.js
const app = getApp()

Page({
  data: {
    formData: {
      name: '',
      description: ''
    }
  },

  // 输入家族名称
  inputName(e) {
    this.setData({
      'formData.name': e.detail.value
    })
  },

  // 输入家族描述
  inputDescription(e) {
    this.setData({
      'formData.description': e.detail.value
    })
  },

  // 提交表单
  submitForm() {
    // 验证必填字段
    if (!this.data.formData.name) {
      wx.showToast({
        title: '请输入家族名称',
        icon: 'none'
      })
      return
    }

    wx.showLoading({ title: '创建中...' })

    wx.request({
      url: `${app.globalData.apiBase}/api/v1/families`,
      method: 'POST',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`,
        'Content-Type': 'application/json'
      },
      data: this.data.formData,
      success: res => {
        wx.hideLoading()
        if (res.data.code === 0) {
          wx.showToast({
            title: '创建成功',
            icon: 'success'
          })
          
          setTimeout(() => {
            wx.navigateBack()
          }, 1500)
        } else {
          wx.showToast({
            title: res.data.message || '创建失败',
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
  }
})
