// pages/family/detail/detail.js
const app = getApp()

Page({
  data: {
    familyId: '',
    family: null,
    members: [],
    activities: [],
    activeTab: 'activity'
  },

  onLoad(options) {
    if (options.id) {
      this.setData({ familyId: options.id })
      this.loadFamilyDetail()
      this.loadMembers()
      this.loadActivities()
    }
  },

  // 加载家族圈详情
  loadFamilyDetail() {
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/families/${this.data.familyId}`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        if (res.data.code === 0) {
          this.setData({ family: res.data.data })
        }
      }
    })
  },

  // 加载成员列表
  loadMembers() {
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/families/${this.data.familyId}/members`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        if (res.data.code === 0) {
          this.setData({ members: res.data.data || [] })
        }
      }
    })
  },

  // 加载动态列表
  loadActivities() {
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/families/${this.data.familyId}/activities`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        if (res.data.code === 0) {
          this.setData({ activities: res.data.data || [] })
        }
      }
    })
  },

  // 切换标签
  switchTab(e) {
    const tab = e.currentTarget.dataset.tab
    this.setData({ activeTab: tab })
  },

  // 邀请成员
  inviteMember() {
    wx.showModal({
      title: '邀请成员',
      content: `邀请码：${this.data.family.inviteCode}`,
      confirmText: '复制',
      success: res => {
        if (res.confirm) {
          wx.setClipboardData({
            data: this.data.family.inviteCode,
            success: () => {
              wx.showToast({
                title: '已复制邀请码',
                icon: 'success'
              })
            }
          })
        }
      }
    })
  },

  // 前往纪念馆详情
  goToMemorial(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/memorial/detail/detail?id=${id}`
    })
  }
})
