// components/media-uploader/media-uploader.js
const app = getApp()

Component({
  properties: {
    type: {
      type: String,
      value: 'image' // image, video, voice
    },
    maxCount: {
      type: Number,
      value: 9
    },
    maxSize: {
      type: Number,
      value: 10 * 1024 * 1024 // 10MB
    }
  },

  data: {
    fileList: [],
    uploading: false,
    uploadProgress: 0
  },

  methods: {
    // 选择文件
    chooseFile() {
      if (this.data.fileList.length >= this.properties.maxCount) {
        wx.showToast({
          title: `最多上传${this.properties.maxCount}个文件`,
          icon: 'none'
        })
        return
      }

      if (this.properties.type === 'image') {
        this.chooseImage()
      } else if (this.properties.type === 'video') {
        this.chooseVideo()
      } else if (this.properties.type === 'voice') {
        this.recordVoice()
      }
    },

    // 选择图片
    chooseImage() {
      wx.chooseImage({
        count: this.properties.maxCount - this.data.fileList.length,
        sizeType: ['compressed'],
        sourceType: ['album', 'camera'],
        success: res => {
          res.tempFilePaths.forEach(filePath => {
            this.uploadFile(filePath, 'image')
          })
        }
      })
    },

    // 选择视频
    chooseVideo() {
      wx.chooseVideo({
        sourceType: ['album', 'camera'],
        maxDuration: 60,
        camera: 'back',
        success: res => {
          if (res.size > this.properties.maxSize) {
            wx.showToast({
              title: '视频文件过大',
              icon: 'none'
            })
            return
          }
          this.uploadFile(res.tempFilePath, 'video')
        }
      })
    },

    // 录制语音
    recordVoice() {
      const recorderManager = wx.getRecorderManager()
      
      wx.showModal({
        title: '录制语音',
        content: '点击确定开始录制，最长60秒',
        success: res => {
          if (res.confirm) {
            recorderManager.start({
              duration: 60000,
              format: 'mp3'
            })

            wx.showToast({
              title: '正在录制...',
              icon: 'loading',
              duration: 60000
            })

            recorderManager.onStop(res => {
              wx.hideToast()
              this.uploadFile(res.tempFilePath, 'voice')
            })

            // 60秒后自动停止
            setTimeout(() => {
              recorderManager.stop()
            }, 60000)
          }
        }
      })
    },

    // 上传文件
    uploadFile(filePath, fileType) {
      this.setData({ uploading: true, uploadProgress: 0 })

      const uploadTask = wx.uploadFile({
        url: `${app.globalData.apiBase}/api/v1/media/upload`,
        filePath: filePath,
        name: 'file',
        formData: {
          type: fileType
        },
        header: {
          'Authorization': `Bearer ${app.globalData.token}`
        },
        success: res => {
          const data = JSON.parse(res.data)
          if (data.code === 0) {
            const newFile = {
              url: data.data.url,
              type: fileType,
              localPath: filePath
            }
            
            this.setData({
              fileList: [...this.data.fileList, newFile],
              uploading: false
            })

            this.triggerEvent('upload', { file: newFile, fileList: this.data.fileList })

            wx.showToast({
              title: '上传成功',
              icon: 'success'
            })
          } else {
            wx.showToast({
              title: data.message || '上传失败',
              icon: 'none'
            })
            this.setData({ uploading: false })
          }
        },
        fail: () => {
          wx.showToast({
            title: '上传失败',
            icon: 'none'
          })
          this.setData({ uploading: false })
        }
      })

      uploadTask.onProgressUpdate(res => {
        this.setData({ uploadProgress: res.progress })
      })
    },

    // 预览文件
    previewFile(e) {
      const index = e.currentTarget.dataset.index
      const file = this.data.fileList[index]

      if (file.type === 'image') {
        const urls = this.data.fileList
          .filter(f => f.type === 'image')
          .map(f => f.url)
        wx.previewImage({
          current: file.url,
          urls: urls
        })
      } else if (file.type === 'video') {
        wx.navigateTo({
          url: `/pages/common/video-player/video-player?url=${encodeURIComponent(file.url)}`
        })
      } else if (file.type === 'voice') {
        this.playVoice(file.url)
      }
    },

    // 播放语音
    playVoice(url) {
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

    // 删除文件
    deleteFile(e) {
      const index = e.currentTarget.dataset.index
      
      wx.showModal({
        title: '提示',
        content: '确定要删除这个文件吗？',
        success: res => {
          if (res.confirm) {
            const fileList = this.data.fileList.filter((_, i) => i !== index)
            this.setData({ fileList })
            this.triggerEvent('delete', { index, fileList })
          }
        }
      })
    }
  }
})
