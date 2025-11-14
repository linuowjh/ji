// utils/performance.js

/**
 * 图片懒加载
 */
function lazyLoadImage(selector, callback) {
  const observer = wx.createIntersectionObserver()
  
  observer.relativeToViewport({ bottom: 100 }).observe(selector, (res) => {
    if (res.intersectionRatio > 0) {
      if (callback) callback(res)
      observer.disconnect()
    }
  })
  
  return observer
}

/**
 * 防抖函数
 */
function debounce(fn, delay = 300) {
  let timer = null
  return function(...args) {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      fn.apply(this, args)
    }, delay)
  }
}

/**
 * 节流函数
 */
function throttle(fn, delay = 300) {
  let lastTime = 0
  return function(...args) {
    const now = Date.now()
    if (now - lastTime >= delay) {
      fn.apply(this, args)
      lastTime = now
    }
  }
}

/**
 * 图片压缩
 */
function compressImage(src, quality = 0.8) {
  return new Promise((resolve, reject) => {
    wx.compressImage({
      src: src,
      quality: quality * 100,
      success: res => resolve(res.tempFilePath),
      fail: reject
    })
  })
}

/**
 * 预加载图片
 */
function preloadImages(urls) {
  return Promise.all(
    urls.map(url => {
      return new Promise((resolve, reject) => {
        wx.getImageInfo({
          src: url,
          success: resolve,
          fail: reject
        })
      })
    })
  )
}

/**
 * 缓存数据
 */
function cacheData(key, data, expire = 3600000) {
  const cacheItem = {
    data: data,
    timestamp: Date.now(),
    expire: expire
  }
  wx.setStorageSync(key, cacheItem)
}

/**
 * 获取缓存数据
 */
function getCachedData(key) {
  try {
    const cacheItem = wx.getStorageSync(key)
    if (!cacheItem) return null
    
    const now = Date.now()
    if (now - cacheItem.timestamp > cacheItem.expire) {
      wx.removeStorageSync(key)
      return null
    }
    
    return cacheItem.data
  } catch (e) {
    return null
  }
}

/**
 * 清除过期缓存
 */
function clearExpiredCache() {
  try {
    const info = wx.getStorageInfoSync()
    const keys = info.keys
    
    keys.forEach(key => {
      const cacheItem = wx.getStorageSync(key)
      if (cacheItem && cacheItem.timestamp && cacheItem.expire) {
        const now = Date.now()
        if (now - cacheItem.timestamp > cacheItem.expire) {
          wx.removeStorageSync(key)
        }
      }
    })
  } catch (e) {
    console.error('Clear cache error:', e)
  }
}

/**
 * 页面性能监控
 */
function monitorPerformance(pageName) {
  const performance = wx.getPerformance()
  const observer = performance.createObserver((entryList) => {
    const entries = entryList.getEntries()
    entries.forEach(entry => {
      console.log(`[Performance] ${pageName}:`, {
        name: entry.name,
        duration: entry.duration,
        startTime: entry.startTime
      })
    })
  })
  
  observer.observe({ entryTypes: ['render', 'script', 'navigation'] })
  
  return observer
}

module.exports = {
  lazyLoadImage,
  debounce,
  throttle,
  compressImage,
  preloadImages,
  cacheData,
  getCachedData,
  clearExpiredCache,
  monitorPerformance
}
