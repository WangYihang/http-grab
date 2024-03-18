package option

type Option struct {
	InputFilePath  string `long:"input" description:"input file path" required:"true"`
	OutputFilePath string `long:"output" description:"output file path" required:"true"`
	StatusFilePath string `long:"status" description:"status file path" required:"true" default:"-"`

	NumWorkers               int   `long:"num-workers" description:"number of workers" default:"32"`
	NumShards                int64 `long:"num-shards" description:"number of shards" default:"1"`
	Shard                    int64 `long:"shard" description:"shard" default:"0"`
	MaxTries                 int   `long:"max-tries" description:"max tries" default:"2"`
	MaxRuntimePerTaskSeconds int   `long:"max-runtime-per-task-seconds" description:"max runtime per task seconds" default:"16"`

	Port int    `long:"port" description:"port" default:"80"`
	Path string `long:"path" description:"path" default:"index.html"`
	Host string `long:"host" description:"http host header, leave it blank to use the IP address" default:""`
}
