// pages/family/list/list.js
const app = getApp()

Page({
  data: {
    families: [],
    loading: false
  },

  onLoad() {
    this.loadFamilies()
  },

  onShow() {
    this.loadFamilies()
  },

  // 加载家族圈列表
  loadFamilies() {
    this.setData({ loading: true })

    wx.request({
      url: `${app.globalData.apiBase}/api/v1/families`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        if (res.data.code === 0) {
          this.setData({
            families: res.data.data || [],
            loading: false
          })
        } else {
          wx.showToast({
            title: res.data.message || '加载失败',
            icon: 'none'
          })
          this.setData({ loading: false })
        }
      },
      fail: () => {
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        })
        this.setData({ loading: false })
      }
    })
  },

  // 创建家族圈
  createFamily() {
    wx.navigateTo({
      url: '/pages/family/create/create'
    })
  },

  // 加入家族圈
  joinFamily() {
    wx.showModal({
      title: '加入家族圈',
      editable: true,
      placeholderText: '请输入邀请码',
      success: res => {
        if (res.confirm && res.content) {
          this.doJoinFamily(res.content)
        }
      }
    })
  },

  // 执行加入家族圈
  doJoinFamily(inviteCode) {
    wx.showLoading({ title: '加入中...' })

    wx.request({
      url: `${app.globalData.apiBase}/api/v1/families/join`,
      method: 'POST',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`,
        'Content-Type': 'application/json'
      },
      data: { inviteCode },
      success: res => {
        wx.hideLoading()
        if (res.data.code === 0) {
          wx.showToast({
            title: '加入成功',
            icon: 'success'
          })
          this.loadFamilies()
        } else {
          wx.showToast({
            title: res.data.message || '加入失败',
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

  // 前往家族圈详情
  goToDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/family/detail/detail?id=${id}`
    })
  }
})
