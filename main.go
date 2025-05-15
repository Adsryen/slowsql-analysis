package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const helpText = `慢查询日志分析工具 v1.0

用法: 
    ./slowsql-analysis -f <慢查询日志路径> [-port <端口>] [-startTime <开始时间>] [-endTime <结束时间>]

参数:
    -f          慢查询日志文件路径
    -port       Web服务端口，设置后可通过浏览器访问报告
    -startTime  开始时间 (可选，格式: yyyy-mm-dd HH:mm:ss)
    -endTime    结束时间 (可选，格式: yyyy-mm-dd HH:mm:ss)

示例:
    1. 基本分析:
       ./slowsql-analysis -f /var/log/mysql-slow.log

    2. 启动Web服务:
       ./slowsql-analysis -f /var/log/mysql-slow.log -port 6033

    3. 指定时间范围:
       ./slowsql-analysis -f /var/log/mysql-slow.log -startTime="2024-04-16 00:00:00" -endTime="2024-04-16 23:59:59"

    4. 完整功能:
       ./slowsql-analysis -f /var/log/mysql-slow.log -port 6033 -startTime="2024-04-16 00:00:00" -endTime="2024-04-16 23:59:59"

输出:
    生成的报告文件格式: <原始日志名>-<生成时间>.html
    如果指定了端口，可以通过浏览器访问: http://<IP>:<端口>/<报告文件名>`

func init() {
	flag.Usage = func() {
		printColoredInfo("blue", helpText)
	}
}

var logAddress = flag.String("f", "", "慢查询日志文件路径")
var startTime = flag.String("startTime", "", "分析开始时间 (格式: yyyy-mm-dd HH:mm:ss)")
var endTime = flag.String("endTime", "", "分析结束时间 (格式: yyyy-mm-dd HH:mm:ss)")
var port = flag.Int("port", 0, "Web服务端口，设置后可通过浏览器访问报告")

// 打印分隔线
func printDivider() {
	fmt.Println("========================================")
}

// 打印带颜色的信息
func printColoredInfo(color string, format string, args ...interface{}) {
	colorCode := map[string]string{
		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"blue":   "\033[34m",
		"reset":  "\033[0m",
	}
	fmt.Printf(colorCode[color]+format+colorCode["reset"]+"\n", args...)
}

func hasDuplicate(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

type SlowSqlInfoSliceDecrement []SlowSqlInfo

func (s SlowSqlInfoSliceDecrement) Len() int { return len(s) }

func (s SlowSqlInfoSliceDecrement) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SlowSqlInfoSliceDecrement) Less(i, j int) bool { return s[i].Time95 > s[j].Time95 }

func getBaseFileName(logPath string) string {
	// 获取文件名（不含路径）
	fileName := logPath[strings.LastIndex(logPath, "/")+1:]
	// 移除.log后缀
	baseName := strings.TrimSuffix(fileName, ".log")
	return baseName
}

// 启动Web服务
func startWebServer(port int, fileName string) {
	// 创建一个文件服务器，提供当前目录下的文件访问
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// 获取本机IP地址
	addrs, err := getLocalIPs()
	if err != nil {
		printColoredInfo("red", "获取本机IP地址失败: %s", err.Error())
		return
	}

	// 打印访问链接
	printColoredInfo("green", "\n报告可通过以下地址访问:")
	for _, addr := range addrs {
		printColoredInfo("blue", "http://%s:%d/%s", addr, port, fileName)
	}
	printColoredInfo("yellow", "按 Ctrl+C 停止Web服务\n")

	// 启动HTTP服务
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			printColoredInfo("red", "启动Web服务失败: %s", err.Error())
			os.Exit(1)
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

// 获取本机IP地址
func getLocalIPs() ([]string, error) {
	cmd := exec.Command("hostname", "-I")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	// 分割并清理IP地址
	ips := strings.Fields(string(output))
	return ips, nil
}

func main() {
	execStartTime := time.Now()
	
	flag.Parse()

	if *logAddress == "" {
		printColoredInfo("blue", "使用方法: ./slowsql-analysis -f <慢查询日志路径> [-port <端口>]")
		printColoredInfo("blue", "示例: ./slowsql-analysis -f /var/log/mysql-slow.log -port 6033")
		printColoredInfo("yellow", "请输入慢查询日志文件路径: ")
		
		// 读取用户输入
		var input string
		fmt.Scanln(&input)
		
		if input == "" {
			os.Exit(1)
		}
		
		// 更新日志文件路径
		*logAddress = input
	}

	printDivider()
	printColoredInfo("blue", "开始分析慢查询日志...")
	printColoredInfo("blue", "日志文件: %s", *logAddress)
	if *startTime != "" && *endTime != "" {
		printColoredInfo("blue", "分析时间范围: %s 至 %s", *startTime, *endTime)
	}
	printDivider()

	var ptCmd string
	if *startTime == "" || *endTime == "" {
		ptCmd = fmt.Sprintf("./cmd/pt-query-digest  %s --output json  --noversion-check --progress time,1 --charset=utf8mb4 >mysql_slow.json", *logAddress)
	} else {
		ptCmd = fmt.Sprintf("./cmd/pt-query-digest  %s --output json  --noversion-check --set-vars time_zone='+8:00' --progress time,1 --charset=utf8mb4 --since='%s' --until='%s' >mysql_slow.json", *logAddress, *startTime, *endTime)
	}

	printColoredInfo("yellow", "正在执行日志分析...")
	cmd := exec.Command("/bin/bash", "-c", ptCmd)
	if err := cmd.Run(); err != nil {
		printColoredInfo("red", "分析过程出错: %s", err.Error())
		os.Exit(1)
	}

	// 生成基于原始日志文件名的输出文件名
	baseFileName := getBaseFileName(*logAddress)
	currentTime := time.Now().Format("2006-01-02-15-04")
	fileName := fmt.Sprintf("%s-%s.html", baseFileName, currentTime)

	printColoredInfo("yellow", "正在生成分析报告: %s", fileName)

	newFile, err := os.Create(fileName)
	if err != nil {
		printColoredInfo("red", "创建报告文件失败: %s", err.Error())
		os.Exit(1)
	}

	file, err := os.Open("mysql_slow.json")
	if err != nil {
		printColoredInfo("red", "打开分析结果失败: %s", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	var report Report
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&report)
	if err != nil {
		printColoredInfo("red", "解析JSON数据失败: %s", err.Error())
		os.Exit(1)
	}

	var slowSqlInfos []SlowSqlInfo
	allSqlInfo := report.Classes

	printColoredInfo("yellow", "正在处理查询信息...")
	for _, sqlInfo := range allSqlInfo {
		var allTables []string
		var slowSqlInfo SlowSqlInfo
		for _, slowTable := range sqlInfo.Tables {
			s := strings.Split(slowTable.Create, ".")
			re := regexp.MustCompile("`([^`]+)`")
			match := re.FindStringSubmatch(s[1])
			tableName := match[1]
			flag := hasDuplicate(allTables, tableName)
			if !flag {
				allTables = append(allTables, tableName)
			}
		}
		slowSqlInfo.RowsSum = sqlInfo.Metrics.RowsExamined.Sum
		slowSqlInfo.RowsMax = sqlInfo.Metrics.RowsExamined.Max
		slowSqlInfo.LengthSum = sqlInfo.Metrics.QueryLength.Sum
		slowSqlInfo.LengthMax = sqlInfo.Metrics.QueryLength.Max
		slowSqlInfo.TimeMax = sqlInfo.Metrics.QueryTime.Max
		slowSqlInfo.TimeMin = sqlInfo.Metrics.QueryTime.Min
		slowSqlInfo.Time95 = sqlInfo.Metrics.QueryTime.Pct95
		slowSqlInfo.TimeMedian = sqlInfo.Metrics.QueryTime.Median
		slowSqlInfo.RowSendMax = sqlInfo.Metrics.RowsSent.Max
		slowSqlInfo.QueryDb = sqlInfo.Metrics.Db.Value
		slowSqlInfo.QueryCount = sqlInfo.QueryCount
		slowSqlInfo.Sql = sqlInfo.Example.Query
		slowSqlInfo.QueryTables = allTables
		slowSqlInfo.Id = sqlInfo.Checksum
		slowSqlInfo.User = sqlInfo.Metrics.User.Value
		slowSqlInfo.Host = sqlInfo.Metrics.Host.Value
		slowSqlInfo.LockTimeMax = sqlInfo.Metrics.LockTime.Max
		slowSqlInfo.LockTimeMin = sqlInfo.Metrics.LockTime.Min
		slowSqlInfo.LockTime95 = sqlInfo.Metrics.LockTime.Pct95
		slowSqlInfo.QueryId = sqlInfo.Example.Id
		slowSqlInfo.Timestamp = sqlInfo.Example.Ts
		slowSqlInfos = append(slowSqlInfos, slowSqlInfo)
	}

	sort.Sort(SlowSqlInfoSliceDecrement(slowSqlInfos))

	// 添加自定义模板函数
	funcMap := template.FuncMap{
		"float64": func(s string) float64 {
			f, _ := strconv.ParseFloat(s, 64)
			return f
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"int64": func(i int64) int64 {
			return i
		},
		"formatTime": func(s string) string {
			t, _ := strconv.ParseFloat(s, 64)
			switch {
			case t >= 1:
				return fmt.Sprintf("%.2fs", t)
			case t >= 0.001:
				return fmt.Sprintf("%.0fms", t*1000)
			default:
				return fmt.Sprintf("%.0fμs", t*1000000)
			}
		},
	}

	printColoredInfo("yellow", "正在生成HTML报告...")
	
	// 使用自定义函数创建新模板
	tmpl, err := template.New("template.html").Funcs(funcMap).ParseFiles("./template/template.html")
	if err != nil {
		printColoredInfo("red", "创建HTML模板失败: %s", err.Error())
		os.Exit(1)
	}

	err = tmpl.Execute(newFile, slowSqlInfos)
	if err != nil {
		printColoredInfo("red", "生成HTML报告失败: %s", err.Error())
		os.Exit(1)
	}

	// 输出统计信息
	printDivider()
	printColoredInfo("green", "分析完成!")
	printColoredInfo("blue", "统计信息:")
	printColoredInfo("blue", "- 总分析SQL数: %d", len(slowSqlInfos))
	printColoredInfo("blue", "- 分析耗时: %.2f秒", time.Since(execStartTime).Seconds())
	printColoredInfo("blue", "- 报告文件: %s", fileName)
	printDivider()

	// 在生成报告后，如果指定了端口，启动Web服务
	if *port > 0 {
		startWebServer(*port, fileName)
	} else {
		printColoredInfo("blue", "\n提示: 使用 -port 参数可启动Web服务访问报告")
		printColoredInfo("blue", "示例: ./slowsql-analysis -f %s -port 6033\n", *logAddress)
	}
}

type Report struct {
	Global struct {
		UniqueQueryCount int `json:"unique_query_count"`
		Files            []struct {
			Name string `json:"name"`
			Size int    `json:"size"`
		} `json:"files"`
		QueryCount int `json:"query_count"`
		Metrics    struct {
			QueryLength struct {
				Sum    string `json:"sum"`
				Stddev string `json:"stddev"`
				Min    string `json:"min"`
				Avg    string `json:"avg"`
				Median string `json:"median"`
				Max    string `json:"max"`
				Pct95  string `json:"pct_95"`
			} `json:"Query_length"`
			LockTime struct {
				Pct95  string `json:"pct_95"`
				Max    string `json:"max"`
				Median string `json:"median"`
				Avg    string `json:"avg"`
				Min    string `json:"min"`
				Stddev string `json:"stddev"`
				Sum    string `json:"sum"`
			} `json:"Lock_time"`
			RowsExamined struct {
				Avg    string `json:"avg"`
				Min    string `json:"min"`
				Median string `json:"median"`
				Max    string `json:"max"`
				Pct95  string `json:"pct_95"`
				Sum    string `json:"sum"`
				Stddev string `json:"stddev"`
			} `json:"Rows_examined"`
			RowsSent struct {
				Max    string `json:"max"`
				Pct95  string `json:"pct_95"`
				Avg    string `json:"avg"`
				Min    string `json:"min"`
				Median string `json:"median"`
				Sum    string `json:"sum"`
				Stddev string `json:"stddev"`
			} `json:"Rows_sent"`
			QueryTime struct {
				Median string `json:"median"`
				Min    string `json:"min"`
				Avg    string `json:"avg"`
				Pct95  string `json:"pct_95"`
				Max    string `json:"max"`
				Stddev string `json:"stddev"`
				Sum    string `json:"sum"`
			} `json:"Query_time"`
		} `json:"metrics"`
	} `json:"global"`
	Classes []struct {
		Distillate string `json:"distillate"`
		Example    struct {
			QueryTime string `json:"Query_time"`
			Query     string `json:"query"`
			Ts        string `json:"ts"`
			AsSelect  string `json:"as_select,omitempty"`
			Id        string `json:"Id,omitempty"`
		} `json:"example"`
		Histograms struct {
			QueryTime []int `json:"Query_time"`
		} `json:"histograms"`
		Fingerprint string `json:"fingerprint"`
		Metrics     struct {
			LockTime struct {
				Pct    string `json:"pct"`
				Stddev string `json:"stddev"`
				Sum    string `json:"sum"`
				Pct95  string `json:"pct_95"`
				Max    string `json:"max"`
				Median string `json:"median"`
				Avg    string `json:"avg"`
				Min    string `json:"min"`
			} `json:"Lock_time"`
			QueryLength struct {
				Pct    string `json:"pct"`
				Stddev string `json:"stddev"`
				Sum    string `json:"sum"`
				Pct95  string `json:"pct_95"`
				Max    string `json:"max"`
				Median string `json:"median"`
				Avg    string `json:"avg"`
				Min    string `json:"min"`
			} `json:"Query_length"`
			RowsSent struct {
				Max    string `json:"max"`
				Pct95  string `json:"pct_95"`
				Avg    string `json:"avg"`
				Min    string `json:"min"`
				Median string `json:"median"`
				Pct    string `json:"pct"`
				Sum    string `json:"sum"`
				Stddev string `json:"stddev"`
			} `json:"Rows_sent"`
			User struct {
				Value string `json:"value"`
			} `json:"user"`
			Db struct {
				Value string `json:"value"`
			} `json:"db,omitempty"`
			RowsExamined struct {
				Median string `json:"median"`
				Min    string `json:"min"`
				Avg    string `json:"avg"`
				Pct95  string `json:"pct_95"`
				Max    string `json:"max"`
				Stddev string `json:"stddev"`
				Sum    string `json:"sum"`
				Pct    string `json:"pct"`
			} `json:"Rows_examined"`
			Host struct {
				Value string `json:"value"`
			} `json:"host"`
			QueryTime struct {
				Avg    string `json:"avg"`
				Min    string `json:"min"`
				Median string `json:"median"`
				Max    string `json:"max"`
				Pct95  string `json:"pct_95"`
				Sum    string `json:"sum"`
				Stddev string `json:"stddev"`
				Pct    string `json:"pct"`
			} `json:"Query_time"`
		} `json:"metrics"`
		TsMin      string `json:"ts_min"`
		Attribute  string `json:"attribute"`
		TsMax      string `json:"ts_max"`
		Checksum   string `json:"checksum"`
		QueryCount int    `json:"query_count"`
		Tables     []struct {
			Status string `json:"status"`
			Create string `json:"create"`
		} `json:"tables,omitempty"`
	} `json:"classes"`
}

type SlowSqlInfo struct {
	Id          string
	RowsSum     string
	RowsMax     string
	LengthSum   string
	LengthMax   string
	TimeMax     string
	TimeMin     string
	Time95      string
	TimeMedian  string
	RowSendMax  string
	QueryDb     string
	QueryCount  int
	QueryTables []string
	Sql         string
	User        string
	Host        string
	LockTimeMax string
	LockTimeMin string
	LockTime95  string
	QueryId     string
	Timestamp   string
}