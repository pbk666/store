//功能可能并不是很完善
//增加的接口如下
1.商品点赞接口
接口描述:用于给指定的商品点赞，提高其热度。

请求方式
POST /product/praise

请求参数:
参数名	类型	  必填	 说明
id	  string	 是	 商品的 ID

响应参数:
参数名	类型	说明
info	string	响应状态信息
status	int	响应状态码
message	string	具体信息
成功响应示例:
{
    "info": "success",
    "status": 10000,
    "message": "点赞成功"
}


2.热门商品推荐接口
接口描述:获取热度最高的商品列表。
请求方式：GET /product/heat
请求参数：无

响应参数:
参数名	类型	说明
info	string	响应状态信息
status	int	响应状态码
products	array	热门商品列表

成功响应示例;
{
    "status": 10000,
    "info": "success",
    "products": 
        {
            "ID": 1,
            "Name": "商品名称",
            "Description": "商品描述",
            "Type": "商品类型",
            "Price": 99.99,
            "Heat": 500
        },
        {
            "ID": 2,
            "Name": "商品名称",
            "Description": "商品描述"
            "Type": "商品类型",
            "Price": 199.99,
            "Heat": 450
        }
    ]
}

3.商品评分接口
接口描述：用于给指定的商品评分，评分范围为 1 到 5 分，并刷新商品的平均评分。
请求方式：POST /product/rating
请求参数：
参数名	类型	必填	说明
productID	int	是	商品的 ID
userID	int	是	用户的 ID
score	int	是	评分（1-5）

响应参数：
参数名	类型	说明
info	string	响应状态信息
status	int	响应状态码

成功响应示例：
json
{
    "info": "success",
    "status": 10000
}
