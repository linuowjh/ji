// components/worship-panel/worship-panel.js
Component({
  properties: {
    type: {
      type: String,
      value: 'flower' // flower, candle, incense, tribute, prayer
    },
    memorialId: {
      type: String,
      value: ''
    }
  },

  data: {
    selectedOption: '',
    inputValue: '',
    options: {
      flower: [
        { value: 'chrysanthemum', label: '菊花', icon: '🌼' },
        { value: 'carnation', label: '康乃馨', icon: '🌹' },
        { value: 'lily', label: '百合', icon: '🌷' }
      ],
      incense: [
        { value: 3, label: '三柱香', icon: '🪔' },
        { value: 9, label: '九柱香', icon: '🪔🪔🪔' }
      ],
      tribute: [
        { value: 'fruit', label: '水果', icon: '🍎' },
        { value: 'cake', label: '糕点', icon: '🍰' },
        { value: 'tea', label: '茶水', icon: '🍵' }
      ]
    }
  },

  lifetimes: {
    attached() {
      // 设置默认选项
      if (this.data.type === 'flower') {
        this.setData({ selectedOption: 'chrysanthemum' })
      } else if (this.data.type === 'incense') {
        this.setData({ selectedOption: 3 })
      } else if (this.data.type === 'tribute') {
        this.setData({ selectedOption: 'fruit' })
      }
    }
  },

  methods: {
    // 选择选项
    selectOption(e) {
      const value = e.currentTarget.dataset.value
      this.setData({ selectedOption: value })
    },

    // 输入内容
    onInput(e) {
      this.setData({ inputValue: e.detail.value })
    },

    // 提交祭扫
    submit() {
      const data = {
        type: this.data.type,
        memorialId: this.properties.memorialId
      }

      if (this.data.type === 'prayer') {
        if (!this.data.inputValue.trim()) {
          wx.showToast({
            title: '请输入祈福语',
            icon: 'none'
          })
          return
        }
        data.content = this.data.inputValue
      } else {
        data.option = this.data.selectedOption
      }

      this.triggerEvent('submit', data)
      
      // 清空输入
      if (this.data.type === 'prayer') {
        this.setData({ inputValue: '' })
      }
    },

    // 录制语音
    recordVoice() {
      this.triggerEvent('recordVoice')
    },

    // 录制视频
    recordVideo() {
      this.triggerEvent('recordVideo')
    }
  }
})
