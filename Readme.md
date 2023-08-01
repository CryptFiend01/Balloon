1. 获取token
   /token
   参数：invitor=xxx
   返回：{"code":0, "data":{"user_id":"xxxxx"}}
2. 重命名
   /name
   参数：user_id=xxxx&name=xxxx
   返回：{"code":0, "data":"ok"}
3. 数据信息
   /info
   参数：user_id=xxxx
   返回：{"code":0, "data":{"user_id":"xxxx","name":"xxx","score":199,"energy":10,"ad_times":0,"invite_times":0}}
4. 开始游戏
   /start
   参数：user_id=xxx
   返回：{"code":0, "data":[150,300,1000]}
5. 记录积分
   /score
   参数：user_id=xxx&score=10
   返回：{"code":0, "data":"ok"}  or {"code":0, "data":{"rank":10, "score":15,"total":100}}
6. 排行榜
   /rank
   参数：user_id=xxxx
   返回：{"code":0, "data":{"rank_info":[{"user_id":"xxx","name":"xxx","score":10,"energy":10}],"self":3,"total":1000}}
7. 获得体力
   /energy
   参数：user_id=xxx&type=ad|invite
   返回：{"code":0, "data":{"energy":3,"ad_times":1}}
8. *注意：所有请求参数使用POST方式*
