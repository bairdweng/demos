
export const suanfa = function () {
  sums();
}
// 两数之和
function sums() {
  console.time();
  // 官方
  var twoSum1 = function (nums, target) {
    var temp = [];
    for (var i = 0; i < nums.length; i++) {
      var dif = target - nums[i];
      if (temp[dif] != undefined) {
        return [temp[dif], i];
      }
      temp[nums[i]] = i;
    }
  };

  var twoSum2 = function (nums, target) {
    let one = 0;
    let two = 0;
    for (let i = 0; i < nums.length; i++) {
      let num = nums[i];
      for (let j = i + 1; j < nums.length; j++) {
        let num2 = nums[j];
        if (num + num2 === target) {
          one = i;
          two = j;
          break;
        }
      }
    }
    return [one, two];
  };


  console.time('时间');
  let indexs2 = twoSum1([1, 2, 3, 4, 5], 4);
  console.log('两数之和', indexs2);
  console.timeEnd('时间');

  // console.time();
  // let indexs1 = twoSum1([1, 2, 3, 4, 5], 4);
  // console.log('两数之和1', indexs1);
  // console.timeEnd();


}