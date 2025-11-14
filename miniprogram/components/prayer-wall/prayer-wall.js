// components/prayer-wall/prayer-wall.js
Component({
  properties: {
    memorialId: {
      type: String,
      value: ''
    }
  },

  data: {
    prayers: [],
    loading: false,
    hasMore: true,
    page: 1
  },

  lifetimes: {
    attached() {
      this.loadPrayers()
    }
  },

  methods: {
    // 加载祈福列表
    loadPrayers() {
      if (this.data.loading || !this.data.hasMore) return

      this.setData({ loading: true })

      const app = getApp()
      wx.request({
        url: `${app.globalData.apiBase}/api/v1/worship/prayers`,
        method: 'GET',
        header: {
          'Authorization': `Bearer ${app.globalData.token}`
        },
        data: {
          memorialId: this.properties.memorialId,
          page: this.data.page,
          pageSize: 10
        },
        success: res => {
          if (res.data.code === 0) {
            const newPrayers = res.data.data || []
            this.setData({
              prayers: [...this.data.prayers, ...newPrayers],
              loading: false,
              hasMore: newPrayers.length === 10,
              page: this.data.page + 1
            })
          } else {
            this.setData({ loading: false })
          }
        },
        fail: () => {
          this.setData({ loading: false })
        }
      })
    },

    // 加载更多
    loadMore() {
      this.loadPrayers()
    },

    // 点赞祈福
    likePrayer(e) {
      const id = e.currentTarget.dataset.id
      const app = getApp()

      wx.request({
        url: `${app.globalData.apiBase}/api/v1/worship/prayers/${id}/like`,
        method: 'POST',
        header: {
          'Authorization': `Bearer ${app.globalData.token}`
        },
        success: res => {
          if (res.data.code === 0) {
            // 更新点赞状态
            const prayers = this.data.prayers.map(prayer => {
              if (prayer.id === id) {
                return {
                  ...prayer,
                  liked: !prayer.liked,
                  likeCount: prayer.liked ? prayer.likeCount - 1 : prayer.likeCount + 1
                }
              }
              return prayer
            })
            this.setData({ prayers })
          }
        }
      })
    }
  }
})
