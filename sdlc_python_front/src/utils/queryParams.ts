export const QueryParams = {
  /**
   * 获取url地址
   * @param name
   */
  getQueryString: function (name, url = '') {
    let reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    let searchUrl;
    let res;
    searchUrl = (url ? `?${url.split('?')[1]}` : url) || window.location.search || `?${window.location.href.split('?')[1]}`;
    res = searchUrl.substr(1).match(reg);
    if (res === null) {
      searchUrl = `?${window.location.href.split('?')[window.location.href.split('?').length - 1]}`;
    }
    res = searchUrl.substr(1).match(reg);
    if (res !== null) return decodeURIComponent(res[2]);
    return '';
  },

  /**
   * 获取url地址--NEW ,如果该方法获取不到会重新用上面的方法获取
   * @param _that
   * @param name
   * @returns {*}
   */

  localQuery: function (_that: any, name: any) {
    var value = '';
    if (!this.isEmpty(_that) &&
      !this.isEmpty(_that.props) &&
      !this.isEmpty(_that.props.location) &&
      !this.isEmpty(_that.props.location.query)) {
      value = _that.props.location.query[name];
    }
    if (this.isEmpty(value) || this.isNull(value)) {
      value = this.getQueryString(name);
    }
    return value;
  },
  /**
   * 判断是不是空的或者undefined
   * @param obj
   * @returns {boolean}
   */

  isNull: function (obj: any) {
    return obj === null || typeof obj === 'undefined' || obj === undefined || obj === 'undefined';
  },

  /**
   * 判断是不是空的字符串
   * @param obj
   * @returns {boolean}
   */

  isEmpty: function (obj: any) {
    return this.isNull(obj) || obj === '';
  },

  replaceParamVal: function (url: string, paramNames: any, replaceWiths: any) {
    let newUrl: string = url
    paramNames.forEach((item: string, index: any) => {
      const value = replaceWiths[index]
      if (newUrl.indexOf(item) === -1) {
        newUrl = newUrl + '&' + item + '=' + value
      } else {
        const re = eval('/(' + item + '=)([^&]*)/gi');
        newUrl = newUrl.replace(re, item + '=' + value);
      }
    });
    return newUrl;
  }


};

/**
* @name 判断是不是微信浏览器
*/
export const isWxBrowser = () => {
  const ua = navigator.userAgent.toLowerCase();
  const isWx = ua.indexOf('micromessenger') === -1;
  if (isWx) {
    return false;
  } else {
    return ua.indexOf('wxwork') === -1;
  }
};

/**处理卡号 每四位隔离 */
export const calcCardNumber = (str = '') => {
  return str.replace(/\s/g, '').replace(/(.{4})/g, "$1 ")
};

/**验证身份证是否正确 */
export const checkIDCard = (idcode = '') => {
  // 加权因子
  let weight_factor = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2];
  // 校验码
  let check_code = ['1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'];

  let code = idcode + "";
  let last = idcode[17];//最后一位

  let seventeen = code.substring(0, 17);

  // ISO 7064:1983.MOD 11-2
  // 判断最后一位校验码是否正确
  let arr: any = seventeen.split("");
  let len = arr.length;
  let num = 0;
  for (let i = 0; i < len; i++) {
    num = num + arr[i] * weight_factor[i];
  }

  // 获取余数
  let resisue = num % 11;
  let last_no = check_code[resisue];

  // 格式的正则
  // 正则思路
  /*
  第一位不可能是0
  第二位到第六位可以是0-9
  第七位到第十位是年份，所以七八位为19或者20
  十一位和十二位是月份，这两位是01-12之间的数值
  十三位和十四位是日期，是从01-31之间的数值
  十五，十六，十七都是数字0-9
  十八位可能是数字0-9，也可能是X
  */
  let idcard_patter = /^[1-9][0-9]{5}([1][9][0-9]{2}|[2][0][0|1][0-9])([0][1-9]|[1][0|1|2])([0][1-9]|[1|2][0-9]|[3][0|1])[0-9]{3}([0-9]|[X])$/;

  // 判断格式是否正确
  let format = idcard_patter.test(idcode);

  // 返回验证结果，校验码和格式同时正确才算是合法的身份证号码
  return last === last_no && format;
};
