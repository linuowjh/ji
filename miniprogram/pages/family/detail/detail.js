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
        console.log('Members response:', res.data)
        if (res.data.code === 0) {
          const rawMembers = res.data.data?.list || res.data.data || []
          const members = rawMembers.map(member => {
            // 格式化加入时间
            let joinedAt = member.joinedAt || member.joined_at || ''
            if (joinedAt) {
              try {
                const date = new Date(joinedAt)
                joinedAt = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
              } catch (e) {
                console.error('Date parse error:', e)
              }
            }
            return {
              ...member,
              joinedAt: joinedAt
            }
          })
          console.log('Formatted members:', members)
          this.setData({ members })
        } else {
          console.error('Load members error:', res.data.message)
          wx.showToast({
            title: res.data.message || '加载成员失败',
            icon: 'none'
          })
        }
      },
      fail: err => {
        console.error('Load members failed:', err)
        wx.showToast({
          title: '加载成员失败',
          icon: 'none'
        })
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
        console.log('Activities response:', res.data)
        if (res.data.code === 0) {
          const rawActivities = res.data.data?.list || res.data.data || []
          
          // 活动类型映射
          const activityTypeMap = {
            'worship': '进行了祭扫',
            'join': '加入了家族圈',
            'create_memorial': '创建了纪念馆',
            'create_story': '发布了家族故事'
          }
          
          const activities = rawActivities.map(activity => {
            // 处理用户信息
            const userName = activity.user?.nickname || '未知用户'
            const userAvatar = activity.user?.avatarUrl || '/images/default-avatar.png'
            
            // 处理活动类型
            const activityType = activity.activityType || activity.activity_type || 'unknown'
            const actionText = activityTypeMap[activityType] || activityType
            
            // 处理时间
            let createdAt = activity.createdAt || activity.created_at || activity.timestamp || ''
            if (createdAt) {
              try {
                const date = new Date(createdAt)
                const now = new Date()
                const diff = now - date
                const minutes = Math.floor(diff / 60000)
                const hours = Math.floor(diff / 3600000)
                const days = Math.floor(diff / 86400000)
                
                if (minutes < 1) {
                  createdAt = '刚刚'
                } else if (minutes < 60) {
                  createdAt = `${minutes}分钟前`
                } else if (hours < 24) {
                  createdAt = `${hours}小时前`
                } else if (days < 7) {
                  createdAt = `${days}天前`
                } else {
                  createdAt = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
                }
              } catch (e) {
                console.error('Date parse error:', e)
              }
            }
            
            return {
              id: activity.id,
              userName: userName,
              userAvatar: userAvatar,
              actionText: actionText,
              createdAt: createdAt
            }
          })
          
          console.log('Formatted activities:', activities)
          this.setData({ activities })
        }
      },
      fail: err => {
        console.error('Load activities failed:', err)
        wx.showToast({
          title: '加载动态失败',
          icon: 'none'
        })
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
