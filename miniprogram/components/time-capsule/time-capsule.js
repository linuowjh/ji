// components/time-capsule/time-capsule.js
Component({
  properties: {
    memorialId: {
      type: String,
      value: ''
    }
  },

  data: {
    messages: [],
    loading: false
  },

  lifetimes: {
    attached() {
      this.loadMessages()
    }
  },

  methods: {
    // 加载时光信箱消息
    loadMessages() {
      this.setData({ loading: true })

      const app = getApp()
      wx.request({
        url: `${app.globalData.apiBase}/api/v1/worship/messages`,
        method: 'GET',
        header: {
          'Authorization': `Bearer ${app.globalData.token}`
        },
        data: {
          memorialId: this.properties.memorialId
        },
        success: res => {
          if (res.data.code === 0) {
            this.setData({
              messages: res.data.data || [],
              loading: false
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

    // 播放语音消息
    playVoice(e) {
      const url = e.currentTarget.dataset.url
      const innerAudioContext = wx.createInnerAudioContext()
      innerAudioContext.src = url
      innerAudioContext.play()

      wx.showToast({
        title: '正在播放...',
        icon: 'none'
      })

      innerAudioContext.onEnded(() => {
        wx.hideToast()
      })
    },

    // 播放视频消息
    playVideo(e) {
      const url = e.currentTarget.dataset.url
      wx.navigateTo({
        url: `/pages/common/video-player/video-player?url=${encodeURIComponent(url)}`
      })
    }
  }
})
