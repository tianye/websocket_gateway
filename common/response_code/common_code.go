package response_code

const SUCCESS_CODE = 200      //(成功）  服务器已成功处理了请求。通常，这表示服务器提供了请求的网页。
const SUCCESS_NULL_CODE = 204 //(无内容） 服务器成功处理了请求，但未返回任何内容。

const AUTHENTICATION_FAILED = 401 //身份验证失败 踢掉用户

const FAIL_CODE = 500 //失败处理
