// utils/sound.js

/**
 * 音效管理器
 */
class SoundManager {
  constructor() {
    this.audioContext = null
    this.sounds = {
      flower: '/sounds/flower.mp3',
      candle: '/sounds/candle.mp3',
      incense: '/sounds/incense.mp3',
      tribute: '/sounds/tribute.mp3',
      prayer: '/sounds/prayer.mp3',
      bell: '/sounds/bell.mp3',
      success: '/sounds/success.mp3',
      click: '/sounds/click.mp3'
    }
  }

  /**
   * 播放音效
   */
  play(soundName, volume = 1.0) {
    if (!this.sounds[soundName]) {
      console.warn(`Sound ${soundName} not found`)
      return
    }

    const innerAudioContext = wx.createInnerAudioContext()
    innerAudioContext.src = this.sounds[soundName]
    innerAudioContext.volume = volume
    innerAudioContext.play()

    innerAudioContext.onError((err) => {
      console.error('Audio play error:', err)
    })

    return innerAudioContext
  }

  /**
   * 播放背景音乐
   */
  playBackground(url, loop = true, volume = 0.5) {
    if (this.audioContext) {
      this.stopBackground()
    }

    this.audioContext = wx.createInnerAudioContext()
    this.audioContext.src = url
    this.audioContext.loop = loop
    this.audioContext.volume = volume
    this.audioContext.play()

    return this.audioContext
  }

  /**
   * 停止背景音乐
   */
  stopBackground() {
    if (this.audioContext) {
      this.audioContext.stop()
      this.audioContext.destroy()
      this.audioContext = null
    }
  }

  /**
   * 暂停背景音乐
   */
  pauseBackground() {
    if (this.audioContext) {
      this.audioContext.pause()
    }
  }

  /**
   * 恢复背景音乐
   */
  resumeBackground() {
    if (this.audioContext) {
      this.audioContext.play()
    }
  }

  /**
   * 设置音量
   */
  setVolume(volume) {
    if (this.audioContext) {
      this.audioContext.volume = volume
    }
  }
}

// 创建单例
const soundManager = new SoundManager()

module.exports = soundManager
