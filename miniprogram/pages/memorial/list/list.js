// pages/memorial/list/list.js
const app = getApp()

Page({
  data: {
    memorials: [],
    loading: false,
    hasMore: true,
    page: 1,
    pageSize: 10
  },

  onLoad() {
    this.loadMemorials()
  },

  onShow() {
    // 刷新列表
    this.setData({
      memorials: [],
      page: 1,
      hasMore: true
    })
    this.loadMemorials()
  },

  // 加载纪念馆列表
  loadMemorials() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    wx.request({
      url: `${app.globalData.apiBase}/api/v1/memorials`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      data: {
        page: this.data.page,
        pageSize: this.data.pageSize
      },
      success: res => {
        if (res.data.code === 0) {
          const newMemorials = res.data.data || []
          this.setData({
            memorials: [...this.data.memorials, ...newMemorials],
            loading: false,
            hasMore: newMemorials.length === this.data.pageSize,
            page: this.data.page + 1
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

  // 下拉刷新
  onPullDownRefresh() {
    this.setData({
      memorials: [],
      page: 1,
      hasMore: true
    })
    this.loadMemorials()
    setTimeout(() => {
      wx.stopPullDownRefresh()
    }, 1000)
  },

  // 上拉加载更多
  onReachBottom() {
    this.loadMemorials()
  },

  // 创建纪念馆
  createMemorial() {
    wx.navigateTo({
      url: '/pages/memorial/create/create'
    })
  },

  // 前往纪念馆详情
  goToDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/memorial/detail/detail?id=${id}`
    })
  }
})
