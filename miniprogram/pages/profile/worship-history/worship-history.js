// pages/profile/worship-history/worship-history.js
const app = getApp()

Page({
  data: {
    records: [],
    loading: false,
    hasMore: true,
    page: 1,
    pageSize: 20
  },

  onLoad() {
    this.loadRecords()
  },

  // 加载祭扫记录
  loadRecords() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    wx.request({
      url: `${app.globalData.apiBase}/api/v1/worship/user/history`,
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
          const data = res.data.data
          const rawRecords = data.records || []
          
          // 转换数据格式，添加友好的显示字段
          const worshipTypeMap = {
            'flower': '献了鲜花',
            'candle': '点燃了蜡烛',
            'incense': '敬献了香火',
            'tribute': '供奉了供品',
            'prayer': '送上了祈福',
            'message': '留下了留言'
          }
          
          const formattedRecords = rawRecords.map(record => {
            // 处理纪念馆信息
            const memorialName = record.memorial?.deceasedName || '未知纪念馆'
            const memorialId = record.memorial?.id || record.memorialId || ''
            
            // 处理祭扫类型
            const worshipTypeText = worshipTypeMap[record.worshipType] || record.worshipType
            
            // 处理时间
            let createdAt = record.createdAt || ''
            if (createdAt) {
              try {
                const date = new Date(createdAt)
                createdAt = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
              } catch (e) {
                console.error('Date parse error:', e)
              }
            }
            
            // 处理内容 - 尝试解析JSON格式的content
            let content = ''
            if (record.content) {
              try {
                const contentObj = JSON.parse(record.content)
                // 根据不同的祭扫类型提取内容
                if (contentObj.message) {
                  content = contentObj.message
                } else if (contentObj.count) {
                  content = `数量: ${contentObj.count}`
                } else if (contentObj.type) {
                  content = `类型: ${contentObj.type}`
                }
              } catch (e) {
                // 如果不是JSON，直接使用原始内容
                content = record.content
              }
            }
            
            return {
              id: record.id,
              memorialId: memorialId,
              memorialName: memorialName,
              worshipType: record.worshipType,
              worshipTypeText: worshipTypeText,
              createdAt: createdAt,
              content: content
            }
          })
          
          this.setData({
            records: [...this.data.records, ...formattedRecords],
            loading: false,
            hasMore: formattedRecords.length === this.data.pageSize,
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
      records: [],
      page: 1,
      hasMore: true
    })
    this.loadRecords()
    setTimeout(() => {
      wx.stopPullDownRefresh()
    }, 1000)
  },

  // 上拉加载更多
  onReachBottom() {
    this.loadRecords()
  },

  // 前往纪念馆详情
  goToMemorial(e) {
    const memorialId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/memorial/detail/detail?id=${memorialId}`
    })
  }
})
