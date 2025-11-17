// pages/memorial/create/create.js
const app = getApp()

Page({
  data: {
    mode: 'create', // create or edit
    memorialId: '',
    formData: {
      deceasedName: '',
      birthDate: '',
      deathDate: '',
      biography: '',
      avatarUrl: '',
      themeStyle: 'traditional',
      tombstoneStyle: 'marble',
      epitaph: '',
      privacyLevel: 1
    },
    themeOptions: [
      { value: 'traditional', label: '中式传统' },
      { value: 'simple', label: '简约素雅' },
      { value: 'nature', label: '自然清新' }
    ],
    tombstoneOptions: [
      { value: 'marble', label: '大理石' },
      { value: 'granite', label: '花岗岩' },
      { value: 'wood', label: '木质' }
    ],
    privacyOptions: [
      { value: 1, label: '家族可见' },
      { value: 2, label: '私密' }
    ],
    themeIndex: 0,
    tombstoneIndex: 0,
    themeLabel: '',
    tombstoneLabel: ''
  },

  onLoad(options) {
    if (options.id && options.mode === 'edit') {
      this.setData({
        mode: 'edit',
        memorialId: options.id
      })
      this.loadMemorial()
    }
  },

  // 加载纪念馆信息（编辑模式）
  loadMemorial() {
    wx.showLoading({ title: '加载中...' })
    
    wx.request({
      url: `${app.globalData.apiBase}/api/v1/memorials/${this.data.memorialId}`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        wx.hideLoading()
        if (res.data.code === 0) {
          // 只更新formData中的字段，保留选项数组
          const memorialData = res.data.data
          
          // 计算主题和墓碑样式的索引和标签
          const themeIndex = this.getThemeIndex(memorialData.themeStyle || 'traditional')
          const tombstoneIndex = this.getTombstoneIndex(memorialData.tombstoneStyle || 'marble')
          const themeLabel = this.getThemeLabel(memorialData.themeStyle || 'traditional')
          const tombstoneLabel = this.getTombstoneLabel(memorialData.tombstoneStyle || 'marble')
          
          this.setData({
            'formData.deceasedName': memorialData.deceasedName || '',
            'formData.birthDate': memorialData.birthDate || '',
            'formData.deathDate': memorialData.deathDate || '',
            'formData.biography': memorialData.biography || '',
            'formData.avatarUrl': memorialData.avatarUrl || '',
            'formData.themeStyle': memorialData.themeStyle || 'traditional',
            'formData.tombstoneStyle': memorialData.tombstoneStyle || 'marble',
            'formData.epitaph': memorialData.epitaph || '',
            'formData.privacyLevel': memorialData.privacyLevel || 1,
            'themeIndex': themeIndex,
            'tombstoneIndex': tombstoneIndex,
            'themeLabel': themeLabel,
            'tombstoneLabel': tombstoneLabel
          })
          wx.setNavigationBarTitle({
            title: '编辑纪念馆'
          })
        }
      },
      fail: () => {
        wx.hideLoading()
        wx.showToast({
          title: '加载失败',
          icon: 'none'
        })
      }
    })
  },

  // 输入逝者姓名
  inputName(e) {
    this.setData({
      'formData.deceasedName': e.detail.value
    })
  },

  // 选择出生日期
  selectBirthDate(e) {
    this.setData({
      'formData.birthDate': e.detail.value
    })
  },

  // 选择逝世日期
  selectDeathDate(e) {
    this.setData({
      'formData.deathDate': e.detail.value
    })
  },

  // 输入生平简介
  inputBiography(e) {
    this.setData({
      'formData.biography': e.detail.value
    })
  },

  // 输入墓志铭
  inputEpitaph(e) {
    this.setData({
      'formData.epitaph': e.detail.value
    })
  },

  // 选择主题风格
  selectTheme(e) {
    const selectedIndex = parseInt(e.detail.value)
    const selectedOption = this.data.themeOptions[selectedIndex]
    this.setData({
      'formData.themeStyle': selectedOption.value,
      'themeIndex': selectedIndex,
      'themeLabel': selectedOption.label
    })
  },

  // 选择墓碑样式
  selectTombstone(e) {
    const selectedIndex = parseInt(e.detail.value)
    const selectedOption = this.data.tombstoneOptions[selectedIndex]
    this.setData({
      'formData.tombstoneStyle': selectedOption.value,
      'tombstoneIndex': selectedIndex,
      'tombstoneLabel': selectedOption.label
    })
  },

  // 选择隐私级别
  selectPrivacy(e) {
    const selectedIndex = parseInt(e.detail.value)
    const selectedOption = this.data.privacyOptions[selectedIndex]
    this.setData({
      'formData.privacyLevel': selectedOption.value
    })
  },

  // 根据主题值获取索引
  getThemeIndex(themeValue) {
    const themeOptions = this.data.themeOptions
    for (let i = 0; i < themeOptions.length; i++) {
      if (themeOptions[i].value === themeValue) {
        return i
      }
    }
    return 0
  },

  // 根据主题值获取标签
  getThemeLabel(themeValue) {
    const themeOptions = this.data.themeOptions
    for (let i = 0; i < themeOptions.length; i++) {
      if (themeOptions[i].value === themeValue) {
        return themeOptions[i].label
      }
    }
    return ''
  },

  // 根据墓碑样式值获取索引
  getTombstoneIndex(tombstoneValue) {
    const tombstoneOptions = this.data.tombstoneOptions
    for (let i = 0; i < tombstoneOptions.length; i++) {
      if (tombstoneOptions[i].value === tombstoneValue) {
        return i
      }
    }
    return 0
  },

  // 根据墓碑样式值获取标签
  getTombstoneLabel(tombstoneValue) {
    const tombstoneOptions = this.data.tombstoneOptions
    for (let i = 0; i < tombstoneOptions.length; i++) {
      if (tombstoneOptions[i].value === tombstoneValue) {
        return tombstoneOptions[i].label
      }
    }
    return ''
  },

  // 上传照片
  uploadAvatar() {
    wx.chooseImage({
      count: 1,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: res => {
        const tempFilePath = res.tempFilePaths[0]
        this.uploadFile(tempFilePath)
      }
    })
  },

  // 上传文件到服务器
  uploadFile(filePath) {
    wx.showLoading({ title: '上传中...' })
    
    wx.uploadFile({
      url: `${app.globalData.apiBase}/api/v1/media/upload`,
      filePath: filePath,
      name: 'file',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: res => {
        wx.hideLoading()
        const data = JSON.parse(res.data)
        if (data.code === 0) {
          // 将相对路径转换为完整 URL
          let avatarUrl = data.data.url
          if (avatarUrl && !avatarUrl.startsWith('http')) {
            avatarUrl = `${app.globalData.apiBase}${avatarUrl}`
          }
          this.setData({
            'formData.avatarUrl': avatarUrl
          })
          wx.showToast({
            title: '上传成功',
            icon: 'success'
          })
        } else {
          wx.showToast({
            title: data.message || '上传失败',
            icon: 'none'
          })
        }
      },
      fail: () => {
        wx.hideLoading()
        wx.showToast({
          title: '上传失败',
          icon: 'none'
        })
      }
    })
  },

  // 提交表单
  submitForm() {
    // 验证必填字段
    if (!this.data.formData.deceasedName) {
      wx.showToast({
        title: '请输入逝者姓名',
        icon: 'none'
      })
      return
    }

    if (!this.data.formData.birthDate || !this.data.formData.deathDate) {
      wx.showToast({
        title: '请选择生卒日期',
        icon: 'none'
      })
      return
    }

    wx.showLoading({ title: '提交中...' })

    const url = this.data.mode === 'create' 
      ? `${app.globalData.apiBase}/api/v1/memorials`
      : `${app.globalData.apiBase}/api/v1/memorials/${this.data.memorialId}`
    
    const method = this.data.mode === 'create' ? 'POST' : 'PUT'

    wx.request({
      url: url,
      method: method,
      header: {
        'Authorization': `Bearer ${app.globalData.token}`,
        'Content-Type': 'application/json'
      },
      data: this.data.formData,
      success: res => {
        wx.hideLoading()
        if (res.data.code === 0) {
          wx.showToast({
            title: this.data.mode === 'create' ? '创建成功' : '更新成功',
            icon: 'success'
          })
          
          setTimeout(() => {
            if (this.data.mode === 'create') {
              wx.redirectTo({
                url: `/pages/memorial/detail/detail?id=${res.data.data.id}`
              })
            } else {
              wx.navigateBack()
            }
          }, 1500)
        } else {
          wx.showToast({
            title: res.data.message || '操作失败',
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
