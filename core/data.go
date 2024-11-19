package core

const Version = "v1.03"

type FipConfig struct {
	Watch struct {
		Include []string          `yaml:"include"`
		Exclude []string          `yaml:"exclude"`
		Type    []string          `yaml:"type"`
		Mode    []string          `yaml:"mode"`
		Release map[string]string `yaml:"release"`
	}
	Cgroup struct {
		MaxCpuUsage    string `yaml:"max_cpu_usage"`
		MaxMemoryUsage string `yaml:"max_memory_usage"`
	}
}

const ConfigTemplate = `# 监控配置
watch:
  # 包含路径
  # ./          监听当前文件目录及其所有子目录
  # /data/sec   监听/data/sec目录及其所有子目录
  include:
    - ./

  # 排除路径
  # /data/sec/.git    忽略/data/sec/.git目录及其所有子目录
  exclude:

  # 文件类型
  # .*    所有文件
  # .go   后缀为.go的文件
  type:
    - .*

  # 监听模式
  # create  创建文件事件
  # write   写入文件事件
  # remove  删除文件事件
  # rename  重命名文件事件
  # chmod   修改文件权限事件
  mode:
    - create
    - write
    - remove
    - rename
    - chmod

  # 放行时间
  # start  开始时间
  # end    结束时间
  # mode   放行模式（quiet-不产生告警；tag-添加放行标签）
  release:
    start: "00:00:00"
    end: "00:00:00"
    mode: tag

# 资源配置
cgroup:
  # 最大CPU使用率（单位 %）
  max_cpu_usage: 10%
  # 最大内存占用值（可选单位：K、M、G）
  max_memory_usage: 500M`
