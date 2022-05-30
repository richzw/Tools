
# 查看有多少个IP访问：
awk '{print $1}' log_file|sort|uniq|wc -l

# 查看每一个IP访问了多少个页面：
awk '{++S[$1]} END {for (a in S) print a,S[a]}' log_file > log.txt
sort -n -t ' ' -k 2 log.txt #配合sort进一步排序

# 将每个IP访问的页面数进行从小到大排序：
awk '{++S[$1]} END {for (a in S) print S[a],a}' log_file | sort -n

# 查看访问前十个ip地址
awk '{print $1}' |sort|uniq -c|sort -nr |head -10 access_log

# 访问次数最多的10个文件或页面
cat log_file|awk '{print $11}'|sort|uniq -c|sort -nr | head -10
cat log_file|awk '{print $11}'|sort|uniq -c|sort -nr|head -20
awk '{print $1}' log_file |sort -n -r |uniq -c | sort -n -r | head -20

# 列出最最耗时的页面(超过60秒的)的以及对应页面发生次数
cat www.access.log |awk '($NF > 60 && $7~/\.php/){print $7}'|sort -n|uniq -c|sort -nr|head -100

# Shell 分析服务器日志命令集锦 https://mp.weixin.qq.com/s/vUvYdeo5eAXR1vdOpSN0WQ

