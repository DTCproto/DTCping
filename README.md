# DTCping
批量Ping工具

----------------------------------------------
_**配置文件说明：**_
|参数|说明|
|---:|:---|
|number|遍历次数|
|csv_file_name|保存文件名|
|random_values|抽取D段地址值列表(1-254)|
|extract_flag|是否抽取D段地址进行查询(优先度高于MergeFlag)|
|merge_flag|是否合并结果到一个CSV|
|ip_segments|IP段列表|
|colo_flag|是否查询Colo信息|
|iata_src_path|优先读取本地iata对照文件|

----------------------------------------------

_**Usage of DTCping.exe:**_
|参数|说明|
|---:|:---|
|-colo|是否查询Colo信息|
| |Colo Open Flag (default false)|
|-ip string|单ip查询|
| |IP (1.0.0.1)|
|-ips string|ip段批量查询|
| |IP Sgt (1.0.0.0/24)|
|-n int|循环次数(例如：10次)|
| | Ping Number (10) (default 10)
|-name string|域名DNS反查ping所有解析IP|
| |Addr Name (cloudflare.com)|
|-path string|读取配置文件信息(<优先级最高>)|
| |config file path(./config.json)|
|-s string|保存文件名设置|
| |Save File Path (pingIpv4) (default "pingIpv4")|
|-iata string|优先读取本地iata对照文件|
| |iata src file path(default iatas.json)|
| |(default "iatas.json")|