### game2048
* 基于go的后端2048游戏

### 环境
* ubuntu 系统
* 安装好 docker
* make 命令
* golang 1.8+

### 本地运行
* git clone git@github.com:yahaa/game2048.git
* cd game2048
* make
* docker run -d -p 8080:8080 game2048

### api使用说明

*  GET /game2048      //获取当前状态

        // 必须带Authorization：value,value取POST返回的Authorization的值

        //Gmap   当前网格状态
        //Score  当前得分
        //Size   当前游戏网格数 n*n
        {
            "Gmap": [
                [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    2
                ],
                [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                [
                    0,
                    0,
                    0,
                    2,
                    0,
                    0,
                    0,
                    0
                ]
            ],
            "Score": 0,
            "Size": 8
        }


* POST /game2048    //开始游戏


        //不需要带 Authorization，这里相当于登录
        //必须带body,如下
        {
            "username":"zihuaaa",   //username用于标识用户
            "size":8                //游戏的格子数 n*n
        }

        //返回数据格式同GET


* PUT /game2048


        //往四个方向玩游戏
        //必须带有 Authorization,value看 GET
        //必须带body,输入如下
        dir="up"||"left"||"right"||"down"
        {
            "dir":"up"
        }

        //返回数据格式同GET


* DELETE /game2048 退出游戏
