<template>
  <div
    :class="['ai-btn', radius && 'radius', size, type, disabled && 'disabled', disabled2 && 'disabled2', disabled3 && 'disabled3']"
    :style="btnStyle"
    @click="onClicked">
      <text
      :class="['btn-text',(disabled || disabled2 || disabled3) ? 'text-disabled' : `text-${type}`, `text-${size}`]"
      :style="textStyle" @click="onClicked">{{text}}</text>
      <slot></slot>
  </div>
</template>
<script>
export default {
  name: 'ai-button',
  props: {
    text: {
      type: String,
      default: '确认',
    },
    type: {
      type: String,
      default: 'normal',
    },
    size: {
      type: String,
      default: 'medium',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabled2: {
      type: Boolean,
      default: false,
    },
    disabled3: {
      type: Boolean,
      default: false,
    },
    radius: {
      type: Boolean,
      default: true,
    },
    btnStyle: {
      type: Object,
      default() {
        return {}
      },
    },
    textStyle: {
      type: Object,
      default() {
        return {}
      },
    },
  },
  data() {
    return {
    }
  },
  methods: {
    onClicked(e) {
      const {
        type,
        disabled,
        disabled2,
        disabled3,
      } = this
      if (disabled || disabled2 || disabled3) {
        return
      }

      this.$emit('click', { e, type, disabled })
    },
  },
}
</script>

<style lang="scss">
@import '../style/default-variables';

.smaller{
  width: $ai-btn-smaller--width;
  height: $ai-btn-smaller--height;
}
.small{
  width: $ai-btn-small--width;
  height: $ai-btn-small--height;
}
.medium{
  width: $ai-btn-medium--width;
  height: $ai-btn-medium--height;
}
.big{
  width: $ai-btn-big--width;
  height: $ai-btn-big--height;
}
.bigger{
  width: $ai-btn-bigger--width;
  height: $ai-btn-bigger--height;
}
.large{
  width: $ai-btn-large--width;
  height: $ai-btn-large--height;
}
.text-bigger {
  font-size: $ai-fs32;
  font-weight: 600;
}
.radius {
  border-radius: $ai-btn--border-radius;
}
.ai-btn{
  justify-content: center;
  align-items: center;
}
.normal{
  background-color: $ai-btn-nomal--background-color;
  background-image: none;
  border-color: $ai-btn-nomal--border-color;
  border-width: $ai-btn-nomal--border-width;
}
.normal-line{
  background-color: $color-white;
  background-image: none;
  border-color: $ai-color-new-gray-border;
  border-width: $ai-btn-nomal--border-width;
}
.text-normal-line {
  color: $ai-color-new-gray-border;
}
.green {
  background-color: $ai-color-green-btn;
  background-image: none;
  border-color: $ai-color-green-btn;
  border-width: $ai-color-green-btn;
}
.green:active{
  background-image: none;
  background-color: $ai-color-green-btn-active;
}
.text-green{
  color: $color-white;
}
// background-image 不需要时需要清空
.primary{
  background-color: $ai-btn-primary--background-color;
  background-image: none;
  border-width: $ai-btn-primary--border-width;
}
.primary:active{
  background-image: none;
}

.primary-line{
  background-color: $ai-btn-primary-line--background-color;
  background-image: $ai-btn-primary-line--background-image;
  border-width: $ai-btn-primary-line--border-width;
}
.primary-line:active{
  background-image: $ai-btn-primary-line--background-image-active;
}

.text-primary{
  color: $ai-btn-primary--color;
}
.text-primary-line{
  color: $ai-btn-primary-line--color;
}
.clicked{
  background-color: $ai-btn-clicked--background-color;
  background-image: none;
  border-color: $ai-btn-clicked--border-color;
  border-width: $ai-btn-clicked--border-width;
}
.text-clicked{
  color: $ai-btn-clicked--color;
}

.disabled{
  background-color: $ai-btn-disabled--background-color;
  background-image: none;
  border-width: 0;
}
.disabled:active{
  background-image: none;
}
.disabled2{
  background-color: $ai-btn-disabled2--background-color;
  background-image: none;
  border-width: 0;
}
.disabled2:active{
  background-image: none;
}
.disabled3{
  background-color: $ai-color-green-btn-50;
  background-image: none;
  border-width: 0;
}
.disabled3:active{
  background-image: none;
}
.text-disabled{
  color: $ai-btn-disabled--color;
}

.btn-text{
  text-align: center;
  color: $ai-btn--color;
  font-size: $ai-btn--font-size;
}
.green2 {
  background-color: $ai-color-light-green-btn;
  background-image: none;
  border-color: $ai-color-light-green-btn;
  border-width: $ai-color-light-green-btn;
}
.text-green2{
  color: $color-white;
  font-size: $ai-fs22;
}
.green3 {
  background-color: $ai-color-green-btn-50;
  background-image: none;
  border-color: $ai-color-light-green-btn;
  border-width: $ai-color-light-green-btn;
}
.text-green3{
  color: $color-white;
  font-size: $ai-fs22;
}
.text-large {
  font-size: $ai-fs32;
}
.clicked2 {
  background-color: $ai-color-light-gray-border;
  background-image: none;
  border-color: $ai-color-light-gray-border;
  border-width: $ai-color-light-gray-border;
}
.text-clicked2{
  color: $color-black;
  font-size: $ai-fs22;
}

</style>
