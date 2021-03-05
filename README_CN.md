# DERO Golang Pool
DERO Golang Pool目前还处于alpha阶段，欢迎大家对这个池子软件提出反馈/修正。虽然很多池子都使用Redis或其他数据库，但我还是选择了DERO，使用了最新发布的Graviton数据库，该数据库由DERO开发团队构建并支持。关于这个选择，请参见[后端数据库选择](#backend-database-choices)。请阅读README和代码，我很乐意通过问题和其他方式协助解决，请随时联系我。我目前使用这个代码库的 "实时 "池是[pool.dero.net](https://pool.dero.network)。

## 使用这个代码库的池子
* [pool.dero.net](https://pool.dero.network)
* [dero.xmrminers.club](https://dero.xmrminers.club/)
* [lethash.it/dero](https://letshash.it/dero)

## 特点
* 在Golang中开发
* 利用Graviton作为后台，由deroproject核心团队构建和支持。
* 内置http/https服务器，用于网页用户界面。
*采矿硬件监控，跟踪工人是否生病。
* 追踪接受、拒绝和阻止的统计信息。
* 守护进程故障转移，利用多个守护进程（上行），池子将从第一个活着的节点获得工作，同时监控其余节点进行备份。
* 通过使用多线程进行并发的共享处理
* 支持直接向交易所或钱包发送挖矿奖励。
* 允许使用综合地址(dERi)和paymentIDs。
* JSON中的API(http/https)，便于与网络前端整合。
* 利用函数和开关来支持挖矿算法，这样你就可以轻松地从config.json中修改所需的挖矿算法，并且只需在几个地方更新代码。
* 支持固定难度，每个端口都有最低难度设置。
* 支持不同难度的最大跳跃灵活性和自定义设置。
* 支持联营和单独采矿
* PROP付款计划
竞彩足球比分直播现场 *轻量级网页，内置基本的资金池统计，但模板用来做为起步阶段的运行。
* 允许矿工设置捐款到一些定义的捐款地址。只需在连接时将%5（0-100，默认为0）添加到钱包地址/用户名中，每个提交的份额就会有一定比例的捐赠。

###未来功能
* (未来)PPLNS和可能的其他集合计划支持。
* (未来)管理功能/go文件，用于修改/报告付款等。
* (未来)更多的功能集在前端，如管理页面等。

## ＃＃目录
1. [要求](#要求) 
1. 下载与安装](#下载--安装) 
1. 配置](#配置)
1. 建立/启动池子](#buildstart-the-pool)
1. 1. [托管api](#host-the-api)
1. 主持前端](#host-the-frontend)
1. SSL for API and frontend](#ssl-for-api-and-frontend)
1. 后台数据库选择](#backend-database-choices)
1. 捐款](#捐款)
1. [学分](#学分)

###要求
* Coin daemon（找到硬币的repo并从源头构建最新版本）。
    * [Derosuite](https://github.com/deroproject/derosuite/releases/latest)
* [Golang](https://golang.org/dl/)
    * 所有代码都是在Windows和Ubuntu 18.04上用Go v1.13.6构建和测试的。

**不要以root身份运行池**：创建一个没有ssh权限的新用户，以避免安全问题。
```bash
sudo adduser --disabled-password --disabled-login your-user（禁用密码）。
```
要用该用户登录: 
```
sudo su - your-user
```

###下载和安装

* 支持获得的repos。
```bash
去获取github.com/deroproject/derosuite/...
go get github.com/deroproject/graviton/...
去找github.com/gorilla/mux/...
```

* 获取项目回购。
```bash
去获取github.com/Nelbert442/dero-golang-pool。
```

###配置

将您选择的`config_example.json`文件复制到`config.json`中，然后查看每个选项并更改任何选项以符合您的首选设置。

每个字段的解释。

```javascript
{
    /*  将显示在前端的矿工连接到的池主机。
                重要：该值用于定义存储的Graviton树名。一旦你设置了这个值，如果你改变了它
        		你会看到数据的 "损失"，因为你现在存储在一个新的树上，而旧的数据在旧的树上。
        		一旦你达到gravitonMaxSnapshots，TREE和所有的k/v对都会被迁移，而不是整个DB。TODO将来要把所有的东西都搬走。
	*/
	"poolHost": "your_pool_host_name",

	/* 区块链探索器，即explorer.dero.io。 */
	"blockchainExplorer": "https://explorer.dero.io/block/{id}",

	/* 交易探索器，即 explorer.dero.io。 */
	"transactionExplorer": "https://explorer.dero.io/tx/{id}",

    /*  矿池地址 */
	"address": "<pool_DERO_Address>",

	/*	矿池捐赠地址 - 在登录地址上连接%<#>的矿工将捐赠一些%的股份到定义的捐赠地址。默认值为0（无捐赠 */
	"donationAddress": "<pool_donation_DERO_Address>",

	/*	Pool donation description - 用于描述捐款的去向，并显示在前端(#getting_started页面) */
	"donationDescription": "Thank you for supporting our mining pool!",

    /*  True：不用担心验证矿工股份[处理速度更快，但有可能是错误的algo]，False。用内置的derosuite函数验证矿工份额 */
	"bypassShareValidation": false,

    /*  生成层的线程数。 */
	"threads": 1,

    /*  定义矿池使用的算法。在 miner.go 中有对这个开关的引用。 */
	"algo": "astrobwt",

	/* 	定义钱币名称 */
	"coin": "DERO",

	/* 定义DERO的基数，小数点后12位。 */
	"coinUnits": 1000000000000,

	/* 定义前端显示的小数点。 */
	"coinDecimalPlaces": 4,

	/* 定义平均找到一个区块的难度目标（以秒为单位）。 */
	"coinDifficultyTarget": 27,

	/* 用于定义在passThru哈希之前，一行要提交多少个有效的股份[可信]。 */
	"trustedSharesCount": 30,

    /*  定义上游（守护进程）getblocktemplate的刷新频率。
                DERO区块链速度很快，运行在27秒的区块时间上。最佳实践是至少每秒钟更新一次挖矿作业。
                比特币池也是每10秒更新一次矿工作业，BTC区块时间是10分钟 -Captain [03/08/2020] 。
                BTC上10分钟的区块时间10秒更新的例子。~10/600 * 27 = 0.45 [你可以随心所欲的低，它只是增加了守护进程的查询] 。
	*/
	"blockRefreshInterval": "450ms",

	"hashrateExpiration": "3h",		// 工人统计的TTL，通常应该等于API部分的大哈希特窗口。注意：使用 "0s "表示无限的过期时间。

	"storeMinerStatsInterval": "5s",	// 多久运行一次WriteMinerStats()来同步当前所有矿工的矿工图和DB。[不要把这个值用毫秒来表示，至少要留>=1s，2s更好] 。

	/*
		定义在迁移到新的DB之前，对实时DB进行多少次快照（提交）。这个值会直接影响DB随时间增长的大小。
        		经测试，Pool在100k以上的提交范围内工作良好，但是当进入上范围时，DB的大小远远超过了25-40GB，因此实现了迁移。
        		当命中这个值时，pooldb目录将被重命名为pooldb_bak，并提供一个新的pooldb目录。如果需要的话，你有pooldb_bak目录用于任何查询或卸载到冷存储。这个值是可变的，取决于池主的舒适度和存储的可用性。
	*/
	"gravitonMaxSnapshots": 5000,

	/*
		定义一个进程（读或写）等待数据库迁移完成的时间长度。理论上进程不会重叠，但是随着规模的增长和提交数量的增长（因为是线性处理），迁移时间可能会增加超过一些首选时间，并与其他操作重叠。这可以保护这些操作不指向DB的错误内存位置，也可以让进程有效地完成，并在开始任何广泛的事情之前暂停，直到DB迁移完成。
        		这个值可以做得很小或很大，最好不要大于250ms左右，因为它会一直循环，直到不能运行为止，但可能不应该低于25ms左右，这样就不会给进程带来垃圾（不过应该不会引起其他问题，只是比正常的日志说明等待DB的时间要大一些而已
	*/
	"gravitonMigrateWait": "100ms",

	"upstreamCheckInterval": "5s",  // 多久轮询一次上游(守护进程)的成功连接？

	/*
		要轮询新工作的守护进程节点列表。Pool将从第一个活着的节点获得工作，并且
        		在后台检查失败的守护进程，作为备份。当前池的块模板
        		始终缓存在RAM中，所以即使守护进程被切换，块模板仍然存在（除非新的块/工作）。
	*/
	"upstream": [
		{
			"enabled": true,        // 将守护进程的启用设置为 "true"、"utilized "或 "false"、"not utilized"。
			"name": "Derod",        // 为守护进程连接设置名称
			"host": "127.0.0.1",    // 设置到达守护进程的地址
			"port": 30306,          // 设置附加到主机的端口
			"timeout": "10s"        // 设置守护进程连接的超时值
		},
		{
			"enabled": false,
			"name": "Remote Derod",
			"host": "derodaemon.nelbert442.com",
			"port": 20206,
			"timeout": "10s"
		}
	],

	"stratum": {
		"paymentId": {
			"addressSeparator": "+",	// 定义从矿工登录到解析paymentID的分隔符。
		},
		"fixedDiff": {
			"addressSeparator": "."		// 定义从矿工登录到解析固定难度的分离器。
		},
		"workerID": {
			"addressSeparator": "@"		// 定义从矿工登录到解析workerID的分隔符。
		},
		"donatePercent": {
			"addressSeparator": "%"		// 定义从矿工登录到解析捐赠百分比的分离器（提交的股份中捐赠到矿池捐赠地址的百分比）。
		},
		"soloMining": {
			"enabled": true,			// 定义是否启用solo挖矿。通过将此设置为false，即使矿工用适当的solo~连接，他们的ID也不会包括solo~。
			"addressSeparator": "~"		// 定义从矿工登录到解析soloMining时使用的分隔符。
		},

		"timeout": "15m",           // See SetDeadline - https://golang.org/pkg/net/
		"healthCheck": true,		// 如果redis不可用，将错误回复给miner而不是job。 (https://github.com/sammy007/monero-stratum)
		"maxFails": 100,			// 在此数目的redis失败之后，将池标记为有病 (https://github.com/sammy007/monero-stratum)

		"listen": [
			{
				"host": "0.0.0.0",  		// 绑定地址
				"port": 1111,       		// Port for mining apps to connect to
				"diff": 1000,       		// 困难矿工设置在此端口上。 TODO：varDiff并将diff设置为开始diff
				"minDiff": 500,				// 设置每个端口可用于固定（可能用于varDiff [future]）的最小难度
				"maxConn": 32768,    		// 该端口上的最大连接数
				"desc": "Low end hardware"	// 端口配置说明
			},
			{
				"host": "0.0.0.0",
				"port": 3333,
				"diff": 2500,
				"minDiff": 500,
				"maxConn": 32768,
				"desc": "Mid range hardware"
			},
			{
				"host": "0.0.0.0",
				"port": 5555,
				"diff": 5000,
				"minDiff": 500,
				"maxConn": 32768,
				"desc": "High end hardware"
			}
		],

		"varDiff": {
			"enabled": false,		// 将varDiff设置为true，非固定式diff矿工的可变难度或false，默认设置为以上难度配置或固定难度
			"minDiff": 100,			// 设置varDiff的最低难度
			"maxDiff": 1000000,		// 设置varDiff的最大难度
			"targetTime": 20,		// 尝试每这么多秒钟获得1个份额
			"retargetTime": 120,	// 检查我们是否应该每隔这么多秒重新定位一次
			"variancePercent": 30,	// 留出时间将目标百分比更改为目标，而无需重新定位
			"maxJump": 50			// 在一次重新定位中限制差异百分比的增加/减少
		}
	},

	"api": {
		"enabled": true,				// 将api设置为true，自托管api或false（不托管）
		"listen": "0.0.0.0:8082",		// 设置api的绑定地址和端口[注意：poolAddr / api / *（在api.go中定义的统计信息，块等）]
		"statsCollectInterval": "5s",	// 设置统计信息收集的运行间隔
		"hashrateWindow": "10m",		// 每家矿工的快速哈希特估计窗口，从其股份中可以看出
		"payments": 30,					// 前端显示的最大支付次数
		"blocks": 50,					// 在前端显示的最大块数
		"ssl": false,					// 为 api 启用 SSL
		"sslListen": "0.0.0.0:9092",	// 为 SSL api 设置绑定地址和端口
		"certFile": "fullchain.cer",	// 设置完整的证书文件链。包括证书，链和 ca。位于同一目录作为 exe 文件。未来可以使用 filepath 包。		
		"keyFile": "cert.key"			// 为 cert 文件设置密钥文件。位于 exe 文件的同一目录中。 TODO Future 可以使用 filepath 包。
	},

	"unlocker": {
		"enabled": true,			// 设置块 unlocker 为 true，used，或 false，not used
		"poolFee": 0.1,				// 设置游泳池费用。这将从块奖励(支付给游泳池添加器)
		"depth": 60,				// 为块解锁设置深度。此值与核心基本块深度进行比较以进行验证
		"interval": "5m"			// 设置间隔来检查块解锁。检查得越快，进程就会变得越嘈杂/忙碌。
	},

	"payments": {
		"enabled": false,			// 将支付设置为真实的、利用的或假的，未利用的
		"interval": "10m",			// 在此时间间隔内运行支付
		"mixin": 8,					// 为事务定义 mixin
		"maxAddresses": 2,			// 定义发送一个 TX 的最大地址数[通常保持较低的安全性，但1-5就足够了]
		"minPayment": 100,			// 定义最低付款额(uint64) ，即: 1 DERO = 1000000000000
		"walletHost": "127.0.0.1",	// 定义钱包守护进程的主机
		"walletPort": "30309"		// 定义钱包守护进程的端口[ deromainnet 默认为20209，Testnet 默认为30309]
	},

	"website": {
		"enabled": true,			// 设置网站启用为true, utilized, 或false, not utilized。
		"port": "8080",				// 设置网站要绑定的端口。
		"ssl": false,				// 为网站启用SSL。
		"sslPort": "9090",			// 设置SSL站点的绑定端口。
		"certFile": "fullchain.cer",// 设置全链证书文件。"fullchain.cer",//设置全链证书文件。包括证书、链和ca，与exe文件位于同一目录下。TODO 未来可以使用文件路径包。
		"keyFile": "cert.key"		// 设置完整的链式证书文件。"cert.key": // 为证书文件设置密钥文件。与exe文件位于同一目录下。TODO未来可以使用filepath包。
	},

	"poolcharts": {
		"updateInterval": "60s",	// 设置“池”图表的更新间隔
		"hashrate": {
			"enabled": true,		// 将此图表数据类型的存储设置为true / false
			"maximumPeriod": 86400	// 设置图表在前端显示的最长时间。 该值以秒表示，并且应等于或大于/被updateInterval整除。
		},
		"miners": {
			"enabled": true,
			"maximumPeriod": 86400
		},
		"workers": {
			"enabled": true,
			"maximumPeriod": 86400
		},
		"difficulty": {
			"enabled": true,
			"maximumPeriod": 604800
		}
	}
}
```

###建造/启动泳池

每次运行基础。

```bash
go run main.go
```

或建立。

```bash
go build main.go
```

注意：logs/和pooldb/目录是在工作目录下创建的。如果你要配置systemd运行或运行应用程序本身时，请记住这一点。

如果你打算用systemd运行，你可以利用类似下面的配置。

```
[Unit]
Description=dero-golang-pool
After=network.target

[Service]
Type=simple
Restart=on-failure
RestartSec=10
SyslogIdentifier=dero-golang-pool
ExecStart=/pathtoyourdir/yourexecutable
WorkingDirectory=/pathtoyourdir

[Install]
WantedBy=multi-user.target
```

### Host the api

一旦`config.json`将 "api". "enabled "设置为 "true"，它将默认在本地 :8082（或定义的任何端口）上进行监听。你定义的地址和端口需要更新，并反映在`config.js`中，以便前端加载数据到它。你可以使用下面的一个例子来拉取内容，或者直接在浏览器中轮询它。

API示例： * ".../api

* ".../api/stats" Example:

```json
{"blocksTotal":18,"candidates":null,"candidatesTotal":0,"config":{"algo":"astrobwt","blockchainExplorer":"http://127.0.0.1:8081/block/{id}","coin":"DERO","coinDecimalPlaces":4,"coinDifficultyTarget":27,"coinUnits":1000000000000,"fixedDiffAddressSeparator":".","payIDAddressSeparator":"+","paymentInterval":30,"paymentMinimum":10000000000,"paymentMixin":8,"poolFee":0.1,"poolHost":"127.0.0.1","ports":[{"diff":1000,"minDiff":500,"host":"0.0.0.0","port":1111,"maxConn":32768},{"diff":2500,"minDiff":500,"host":"0.0.0.0","port":3333,"maxConn":32768},{"diff":5000,"minDiff":500,"host":"0.0.0.0","port":5555,"maxConn":32768}],"transactionExplorer":"http://127.0.0.1:8081/tx/{id}","unlockDepth":5,"unlockInterval":10,"version":"1.0.0","workIDAddressSeparator":"@"},"immature":[{"Hash":"770efbc1377ca0f1818ac9e01b0f697bd461e716160b24826b6b96931ac392d2","Address":"dEToUEe...8gVNr","Height":1017,"Orphan":false,"Timestamp":1600807603,"Difficulty":22254,"TotalShares":29975,"Reward":2351321493449,"Solo":false},{"Hash":"efca19034b80b48366f984a2bdb81647e786481a1528942d406412b219109f6a","Address":"dEToUEe...8gVNr","Height":1014,"Orphan":false,"Timestamp":1600807420,"Difficulty":21816,"TotalShares":2000,"Reward":2345322388119,"Solo":false},{"Hash":"c3d54ee8d3c7919e0f426ec964516efa33f5d00b4608536c47e389329677425d","Address":"dEToUEe...8gVNr","Height":1016,"Orphan":false,"Timestamp":1600807598,"Difficulty":22254,"TotalShares":27780,"Reward":2345321791672,"Solo":false},{"Hash":"5ba9184f441c125fd67549d1aeecc8a1d1d664d51e1caf62b0357089492a1ee3","Address":"dEToUEe...8gVNr","Height":1013,"Orphan":false,"Timestamp":1600807411,"Difficulty":21600,"TotalShares":2000,"Reward":2345322686342,"Solo":false},{"Hash":"1c3bfe247f02f44c60301bfa54f85fa7e18f1604320ee8f2a775dea66567d128","Address":"dEToUEe...8gVNr","Height":1015,"Orphan":false,"Timestamp":1600807439,"Difficulty":22034,"TotalShares":5000,"Reward":2349822089896,"Solo":false}],"immatureTotal":5,"lastblock":{"Difficulty":"22254","Height":1017,"Timestamp":1600807598,"Reward":2351321493449,"Hash":"770efbc1377ca0f1818ac9e01b0f697bd461e716160b24826b6b96931ac392d2"},"matured":[{"Hash":"339ad336c07e86913f388fb45fc3d03dc03ef9ae7cdd82e98e7ee0d97c470f79","Address":"dEToUEe...8gVNr","Height":1000,"Orphan":false,"Timestamp":1600806375,"Difficulty":21600,"TotalShares":13000,"Reward":2354326563247,"Solo":false},{"Hash":"b2cbf4b90d36a10521092ea3bd8d20d0a29676b190492bb715b188fec17b0130","Address":"dEToUEe...8gVNr","Height":1007,"Orphan":false,"Timestamp":1600807040,"Difficulty":21600,"TotalShares":0,"Reward":2349824475682,"Solo":false},{"Hash":"4454bf01932bc8ae601e8aee345a294e8fde99790e05b71a481b7c4eec4bd084","Address":"dEToUEe...8gVNr","Height":1008,"Orphan":false,"Timestamp":1600807153,"Difficulty":21600,"TotalShares":0,"Reward":2349824177459,"Solo":false},{"Hash":"aadf5246f36cc098b341bf6c694dd08d6ca6969b0784d91c82f3cb3791812652","Address":"dEToUEe...8gVNr","Height":1011,"Orphan":false,"Timestamp":1600807224,"Difficulty":21600,"TotalShares":12000,"Reward":2349823282789,"Solo":false},{"Hash":"dfa60fede87c7c4e7d351c54b87e46c3239209ae10d6db58050a27a9b147457d","Address":"dEToUEe...8gVNr","Height":1012,"Orphan":false,"Timestamp":1600807401,"Difficulty":21600,"TotalShares":5000,"Reward":2354322984565,"Solo":false},{"Hash":"a6eccb0be31558bed06a8add669fe7846d388410e09bb37e8a29c1d5ab992f3e","Address":"dEToUEe...8gVNr","Height":1003,"Orphan":false,"Timestamp":1600806585,"Difficulty":21600,"TotalShares":10500,"Reward":2345325668576,"Solo":false},{"Hash":"f79af5914e15373fa998819cfacc7d74ffe18bb315787572c7fbbe1bb93aaed4","Address":"dEToUEe...8gVNr","Height":1004,"Orphan":false,"Timestamp":1600806855,"Difficulty":21600,"TotalShares":43500,"Reward":2345325370353,"Solo":false},{"Hash":"da99e1f3600508708a38f48959210ca9de914ab524aaa153882fa04c3873811a","Address":"dEToUEe...8gVNr","Height":1010,"Orphan":false,"Timestamp":1600807222,"Difficulty":21600,"TotalShares":0,"Reward":2349823581012,"Solo":false},{"Hash":"3fe81b154a9f4a07fce72d621fbaf169e457baf918be8d092a9b735a2159ce73","Address":"dEToUEe...8gVNr","Height":1002,"Orphan":false,"Timestamp":1600806516,"Difficulty":21600,"TotalShares":11500,"Reward":2345325966800,"Solo":false},{"Hash":"1068ccc0d92c1d49d375a675018154c29b5404bbb297b0f2da329154efe9e832","Address":"dEToUEe...8gVNr","Height":1006,"Orphan":false,"Timestamp":1600807020,"Difficulty":21600,"TotalShares":11250,"Reward":2345324773905,"Solo":false},{"Hash":"98310319fd9e80d97742e4e906a8b594f5423122b6a133511c672aaedfa29277","Address":"dEToUEe...8gVNr","Height":1001,"Orphan":false,"Timestamp":1600806383,"Difficulty":21600,"TotalShares":0,"Reward":2345326265023,"Solo":false},{"Hash":"e5fbce21b8003876d249ff2b050c474c44bc54dbfc7069d1845100d6b55cae42","Address":"dEToUEe...8gVNr","Height":1009,"Orphan":false,"Timestamp":1600807188,"Difficulty":21600,"TotalShares":0,"Reward":2349823879235,"Solo":false},{"Hash":"38984e8ac3ccd2c1ebc4eba781d38a4ecc76d461c80731d6c81ad94265e9d8e4","Address":"dEToUEe...8gVNr","Height":1005,"Orphan":false,"Timestamp":1600806908,"Difficulty":21600,"TotalShares":13500,"Reward":2345325072129,"Solo":false}],"maturedTotal":13,"miners":[{"LastBeat":1600807678,"StartedAt":1600807391,"ValidShares":36,"InvalidShares":0,"StaleShares":0,"Accepts":6,"Rejects":0,"RoundShares":29975,"Hashrate":151,"Offline":false,"Id":"dEToUEe3q57XoqLgbuDE7DUmoB6byMtNBWtz85DmLAHAC8wSpetw4ggLVE4nB3KRMRhnFdxRT3fnh9geaAMmGrhP2UDY18gVNr","Address":"dEToUEe3q57XoqLgbuDE7DUmoB6byMtNBWtz85DmLAHAC8wSpetw4ggLVE4nB3KRMRhnFdxRT3fnh9geaAMmGrhP2UDY18gVNr","IsSolo":false}],"now":1600807685,"payments":[{"Hash":"205e4ac6547a784eb94cba28f50f4a26595f3335ae28a8d3d39dccdf6e0fae10","Timestamp":1600807021,"Payees":1,"Mixin":8,"Amount":2345326265023},{"Hash":"88621a2fee06d0c2d97b8bf5137ed26d22789ec5602263bcad9505c32f9caaf1","Timestamp":1600807202,"Payees":1,"Mixin":8,"Amount":2342980044983},{"Hash":"c24bedcaa513204d5663028821559379544754132d515030c68cf75f76a9eb70","Timestamp":1600807263,"Payees":1,"Mixin":8,"Amount":2342979449131},{"Hash":"2616b795413d6207da75aff72c1b66fd17af3cb7f99fca06bd073c60bd398088","Timestamp":1600807627,"Payees":1,"Mixin":8,"Amount":4699442121086},{"Hash":"e64c7bed69b3dfd2aa02100e9790dfa3e4904c63f59bd5067e4d0f71dbbb4b19","Timestamp":1600806931,"Payees":1,"Mixin":8,"Amount":2351972236684},{"Hash":"186615582db0e54b2e21c23f715d82ccc8b686e3aaeb243486a805517def5872","Timestamp":1600807051,"Payees":1,"Mixin":8,"Amount":2342980640833},{"Hash":"969334e0cd6e40947d9d016509965c7e52ef66e17ed650e700d29285f9c6824d","Timestamp":1600807172,"Payees":1,"Mixin":8,"Amount":2342980342907},{"Hash":"c2f3413e0579de5bba9bd10e810586d051f7a4b4e37e1f316278f15daf5e52ca","Timestamp":1600807233,"Payees":1,"Mixin":8,"Amount":2342979747057},{"Hash":"3eaa0b54c80b7856b46226d927cf114a7abbcbeb8a947cb7d9769590c9abbc24","Timestamp":1600807417,"Payees":1,"Mixin":8,"Amount":2349824475682},{"Hash":"b4e24d9a16ab1a3ae7c9254f43e660b3e697330925d933601c289fecc75f1e8e","Timestamp":1600807447,"Payees":1,"Mixin":8,"Amount":4699647460247}],"poolHashrate":151,"soloHashrate":0,"totalMinersPaid":1,"totalPayments":10,"totalPoolMiners":1,"totalSoloMiners":0}
```

* ".../api/accounts?address=<yourwalletaddress>" Example:

```json
{"address":"dEToUEe3q57XoqLgbuDE7DUmoB6byMtNBWtz85DmLAHAC8wSpetw4ggLVE4nB3KRMRhnFdxRT3fnh9geaAMmGrhP2UDY18gVNr","miners":[{"LastBeat":1603719621,"StartedAt":1603719611,"ValidShares":3,"InvalidShares":0,"StaleShares":0,"Accepts":0,"Rejects":0,"LastRoundShares":0,"RoundShares":4000,"Hashrate":0,"Offline":true,"Id":"dEToUEe3q57XoqLgbuDE7DUmoB6byMtNBWtz85DmLAHAC8wSpetw4ggLVE4nB3KRMRhnFdxRT3fnh9geaAMmGrhP2UDY18gVNr","Address":"dEToUEe3q57XoqLgbuDE7DUmoB6byMtNBWtz85DmLAHAC8wSpetw4ggLVE4nB3KRMRhnFdxRT3fnh9geaAMmGrhP2UDY18gVNr","IsSolo":false},{"LastBeat":1603719643,"StartedAt":1603719633,"ValidShares":1,"InvalidShares":0,"StaleShares":0,"Accepts":0,"Rejects":0,"LastRoundShares":0,"RoundShares":0,"Hashrate":0,"Offline":true,"Id":"solo~dEToUEe3q57XoqLgbuDE7DUmoB6byMtNBWtz85DmLAHAC8wSpetw4ggLVE4nB3KRMRhnFdxRT3fnh9geaAMmGrhP2UDY18gVNr","Address":"dEToUEe3q57XoqLgbuDE7DUmoB6byMtNBWtz85DmLAHAC8wSpetw4ggLVE4nB3KRMRhnFdxRT3fnh9geaAMmGrhP2UDY18gVNr","IsSolo":true}],"payments":[{"Hash":"fae0a899fac54452f90bc4a0c883705fd3ebc17193d169345b3b0476ab5ab48f","Timestamp":1603719241,"Payees":1,"Mixin":8,"Amount":2344919251485},{"Hash":"54656d899b0764639302f19ff6a56985d939b51e3f5748325d04154fadc1ac83","Timestamp":1603719152,"Payees":1,"Mixin":8,"Amount":2344919549085},{"Hash":"7f6a32ab4d95b527cf0b6b3f9a5f4ce52ef2d136d3910459d60ae6a3ad943425","Timestamp":1603718732,"Payees":1,"Mixin":8,"Amount":2340424346685},{"Hash":"0a98cc001b1a677c31c6ac2747b41ba86722b43ef9118299760c8bf80e16cd55","Timestamp":1603718341,"Payees":1,"Mixin":8,"Amount":2350914144285},{"Hash":"10a9632e96d50584ed575e4176393ca30057299e17139fdb16ddb9b702a6c6f4","Timestamp":1603717846,"Payees":1,"Mixin":8,"Amount":2344920441886},{"Hash":"b88604f42dede0d2427c63cbc4bff7d908a36d3fffe2a4080c49d2482686b741","Timestamp":1603717696,"Payees":1,"Mixin":8,"Amount":2344920739487},{"Hash":"0997ecd4ba65e042ed8942769ca57c3facbeccad2ade681de19f780ec05e2843","Timestamp":1603717635,"Payees":1,"Mixin":8,"Amount":2344921037087},{"Hash":"485e602aa179abcc39e14afe1c41aeee5716ee5ccf0ab2a66be9027ed4e820f1","Timestamp":1603717125,"Payees":1,"Mixin":8,"Amount":2344921334688},],"poolHashrate":0,"soloHashrate":0,"totalPayments":196,"totalPoolMiners":0,"totalSoloMiners":0}
```

### ＃＃主持前台

一旦`config.json`将 "website". "enabled "设置为true，它将默认在本地的:8080(或任何一个定义的端口)进行监听。它将像静态网页一样利用标准的js/html/css文件，并与#4中的API集成。

`website.go`是运行器，它只是在定义的端口上启动监听服务，然后在/website/pages中提供内容，可以随意修改文件夹结构，只要确保更新website.go即可。

`config.js`是javascript配置文档，指向API的网址。请务必更新api变量指向上面定义的正确api url。

```
var api = "http://127.0.0.1:8082/api";
```

![DERO Pool Home](images/home.PNG?raw=true "DERO Pool Home") 
![DERO Pool Getting Started](images/gettingstarted.PNG?raw=true "DERO Pool Getting Started")
![DERO Pool Worker Stats](images/workerstats.PNG?raw=true "DERO Pool Worker Stats")
![DERO Pool Blocks](images/poolBlock.PNG?raw=true "DERO Pool Blocks")
![DERO Pool Pay](images/poolpayment.PNG?raw=true "DERO Pool Pay")

### API和前端的SSL

在 "config.json "中，API和网站都有SSL部分，可以利用这些部分来定义你的密钥/证书文件和端口，以便在本地运行网站或api。cert和key文件要放置在与你的构建包相同的目录中。cert文件应该包括cert、chain和ca，而key文件则是你的私钥。将这些文件存储在本地，并确保正确地.gitignore它们，默认有*.key和*.cer，但是不同的人使用不同的文件类型。

### 后台数据库选择

由于这个池子是以[DERO](https://github.com/deroproject/)为基础开发的，所以我决定效仿，将这个池子配置为利用新发布的[Graviton](https://github.com/deroproject/graviton)数据库作为池子的后端。我会介绍一些目前存在的弊端和我已经实现的变通方法，但是这个池子可以很容易地被重新配置为其他DB类型，比如：redis、boltdb、badgerdb、Graviton等。Redis在之前的提交历史中被使用过，你也可以在其他池子上看到Redis的实现，参见Credis部分。

为了让Graviton能够持续发展，我实现了以下一些变通方法。

* Graviton在每次提交数据时都会进行快照。这意味着每次你向它写入数据时，都会生成一个新的快照，作为版本控制历史记录，可以在任何时间点重新引用，并在任何时候与当前或两个不同的快照/树进行比较。

在早期的改编中，我意识到在很短的时间内，我的提交量就达到了100,000+，这使得数据库的大小急剧膨胀。这在一定程度上是由于我最初实施的重提交性质，以及其他一些片子，在时间上进行了一些优化。虽然 "重点 "是保留这个快照的历史备份，但我的实现就是不需要这个要求。为了不直接报废，回到boltdb或者redis，我继续往前走，决定保留X个数量的提交历史和单个备份，以备我以后想把备份推送到某个云/冷存储上，保留一段时间。

为了做到这一点，我定义了一些gravitonMaxSnapshots，我在每次读/写DB时都会检查（低ms检查），直到达到该值（或超过该值），然后抓取所有的k/v对，对当前pooldb目录进行重命名为pooldb_bak，然后提供一个新的pooldb存储，并将k/v对放入其中，提交后继续进行。在这期间，有一个g.migrating属性（设置为0（未迁移）或1（迁移）），在每次读/写时都会对其进行检查。如果db正在迁移，而一些读/写操作试图利用它，那么这个过程将 "等待 "gravitonMigrateWait的时间（比如100ms左右），并不断循环，直到这个过程打开。由于这种情况发生在所有的读写过程中，所以到目前为止，还没有出现进程中途卡住的测试问题，因为进程的提交发生在尾部，而不是沿途。

随着时间的推移，可能看起来Graviton并不合适，然而我并没有因此而放弃，因为我喜欢它的功能，目录的可移植性（可以复制/粘贴实时数据而不会损坏），以及其他潜在的未来功能集。对于每个人来说，欢迎任何人使用这个repo来实现他们喜欢的任何形式的DB。我曾想过保留一个历史记录，这样你就可以很容易地在使用redis或graviton或其他方面进行切换，然而这对于alpha阶段来说似乎有点过于雄心勃勃，也许以后会有更多的东西：)

###捐款

自从我第一次开始为DERO（testnet和mainnet都一样）托管矿池以来，我已经把它作为一个长期的传统，只要我个人能够维持，我就会以0费操作的方式免费托管它们。几年下来，我真的很享受在这些矿池上工作的每一分钟，这个版本是我第一次可以说我知道矿池的来龙去脉，而不是仅仅利用另一个矿池代码库，修改一下就可以了。在[Credits](#credits)下面的一些功能集和想法的帮助下，打好了基础，并进行了大量的重写/修改，形成了今天这个泳池软件。希望您能和我一样喜欢它，在此祝愿技术尽可能的进步。非常感谢您的关注，并随时欢迎您的反馈。

```
dERopdjpGmr2DEwQJdRrKc8M6obca9NQu2EaC2fNe3RNHonYcCfqmjGF7NBEHoB8dpLXWhnjdW7dugFTVhofuKTb4sfzmyBSAj
```

Credits
---------

* [sammy007](https://github.com/sammy007) - Developer on [monero-stratum](https://github.com/sammy007/monero-stratum) 为这个项目的地层建设奠定了基础。
* [JKKGBE](https://github.com/JKKGBE) - Developer on [open-zcash-pool](https://github.com/JKKGBE/open-zcash-pool) which is forked from [sammy007](https://github.com/sammy007) project [open-ethereum-pool](https://github.com/sammy007/open-ethereum-pool) for some additional ideas/thoughts throughout dev when REDIS was utilized, but later migrated to Graviton from scratch with other implementations.
* [Graviton](https://github.com/deroproject/graviton) - 在这个项目中用于后端数据存储的 Graviton DB。
* [Derosuite](https://github.com/deroproject/derosuite) - Derosuite (DERO) ，这是一种加密货币，这个数据池最初是为其建立的，并为其聚焦。