# ai-button

### 示例

```html
      <ai-button slot="button" type="normal" text="关注" :btnStyle="{width: '156px'}" />

      <ai-button slot="button" type="primary" text="关注" :btnStyle="{width: '156px'}" />

      <ai-button slot="button" type="clicked" text="已关注" :btnStyle="{width: '156px'}"  disabled/>


```

### 可配置的参数

| 参数 | 说明 | 类型 | 默认值 | 必填 |
| :--- | :--- | :--- | :--- | :--- |
| text | 展现的文字 | String | 确认 | 是 |
| type | 类型：primary-line,primary，normal,clicked | String | normal | - |
| size | 按钮大小：small,medium,big,large | String | medium | - |
| radius | 圆角 | Boolean | true | 否 |
| disabled | 是否禁用 | Boolean | false | - |
| btnStyle | 按钮的样式对象 | Object | {} | 是 |
| textStyle | 字体的样式对象 | Object | {} | 是 |

> btnStyle > disabled > type
### 回调事件

> @click=""
