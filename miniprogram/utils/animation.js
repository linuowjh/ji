// utils/animation.js

/**
 * 淡入动画
 */
function fadeIn(duration = 300) {
  const animation = wx.createAnimation({
    duration: duration,
    timingFunction: 'ease'
  })
  animation.opacity(1).step()
  return animation.export()
}

/**
 * 淡出动画
 */
function fadeOut(duration = 300) {
  const animation = wx.createAnimation({
    duration: duration,
    timingFunction: 'ease'
  })
  animation.opacity(0).step()
  return animation.export()
}

/**
 * 缩放动画
 */
function scale(scaleValue = 1.1, duration = 300) {
  const animation = wx.createAnimation({
    duration: duration,
    timingFunction: 'ease-in-out'
  })
  animation.scale(scaleValue).step()
  return animation.export()
}

/**
 * 旋转动画
 */
function rotate(angle = 360, duration = 1000) {
  const animation = wx.createAnimation({
    duration: duration,
    timingFunction: 'linear'
  })
  animation.rotate(angle).step()
  return animation.export()
}

/**
 * 滑入动画（从下往上）
 */
function slideUp(duration = 400) {
  const animation = wx.createAnimation({
    duration: duration,
    timingFunction: 'ease-out'
  })
  animation.translateY(0).opacity(1).step()
  return animation.export()
}

/**
 * 滑出动画（从上往下）
 */
function slideDown(duration = 400) {
  const animation = wx.createAnimation({
    duration: duration,
    timingFunction: 'ease-in'
  })
  animation.translateY(100).opacity(0).step()
  return animation.export()
}

/**
 * 弹跳动画
 */
function bounce(duration = 600) {
  const animation = wx.createAnimation({
    duration: duration,
    timingFunction: 'ease-in-out'
  })
  animation.scale(1.2).step({ duration: duration / 2 })
  animation.scale(1).step({ duration: duration / 2 })
  return animation.export()
}

/**
 * 摇晃动画
 */
function shake(duration = 500) {
  const animation = wx.createAnimation({
    duration: duration / 4,
    timingFunction: 'ease-in-out'
  })
  animation.rotate(5).step()
  animation.rotate(-5).step()
  animation.rotate(5).step()
  animation.rotate(0).step()
  return animation.export()
}

/**
 * 心跳动画
 */
function heartbeat(duration = 800) {
  const animation = wx.createAnimation({
    duration: duration / 4,
    timingFunction: 'ease-in-out'
  })
  animation.scale(1.1).step()
  animation.scale(1).step()
  animation.scale(1.1).step()
  animation.scale(1).step()
  return animation.export()
}

/**
 * 烛光闪烁动画
 */
function candleFlicker() {
  const animation = wx.createAnimation({
    duration: 2000,
    timingFunction: 'ease-in-out'
  })
  animation.opacity(1).step({ duration: 1000 })
  animation.opacity(0.8).step({ duration: 1000 })
  return animation.export()
}

module.exports = {
  fadeIn,
  fadeOut,
  scale,
  rotate,
  slideUp,
  slideDown,
  bounce,
  shake,
  heartbeat,
  candleFlicker
}
