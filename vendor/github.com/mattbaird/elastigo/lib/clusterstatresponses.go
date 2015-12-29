package elastigo

type NodeStatsResponse struct {
	ClusterName string `json:"cluster_name"`
	Nodes       map[string]NodeStatsNodeResponse
}

type NodeStatsNodeResponse struct {
	Name             string                                     `json:"name"`
	Timestamp        int64                                      `json:"timestamp"`
	TransportAddress string                                     `json:"transport_address"`
	Hostname         string                                     `json:"hostname"`
	Host             string                                     `json:"host"`
	IP               []string                                   `json:"ip"`
	Attributes       NodeStatsNodeAttributes                    `json:"attributes"`
	Indices          NodeStatsIndicesResponse                   `json:"indices"`
	OS               NodeStatsOSResponse                        `json:"os"`
	Process          NodeStatsProcessResponse                   `json:"process"`
	JVM              NodeStatsJVMResponse                       `json:"jvm"`
	Network          NodeStatsNetworkResponse                   `json:"network"`
	FS               NodeStatsFSResponse                        `json:"fs"`
	ThreadPool       map[string]NodeStatsThreadPoolPoolResponse `json:"thread_pool"`
	Transport        NodeStatsTransportResponse                 `json:"transport"`
	FieldDataBreaker NodeStatsFieldDataBreakerResponse          `json:"fielddata_breaker"`
}

type NodeStatsNodeAttributes struct {
	Data   string `json:"data"`
	Client string `json:"client"`
}
type NodeStatsNetworkResponse struct {
	TCP NodeStatsTCPResponse `json:"tcp"`
}

type NodeStatsFieldDataBreakerResponse struct {
	MaximumSizeInBytes   int64   `json:"maximum_size_in_bytes"`
	MaximumSize          string  `json:"maximum_size"`
	EstimatedSizeInBytes int64   `json:"estimated_size_in_bytes"`
	EstimatedSize        string  `json:"estimated_size"`
	Overhead             float64 `json:"overhead"`
	Tripped              int64   `json:"tripped"`
}
type NodeStatsTransportResponse struct {
	ServerOpen int64 `json:"server_open"`
	RxCount    int64 `json:"rx_count"`
	RxSize     int64 `json:"rx_size_in_bytes"`
	TxCount    int64 `json:"tx_count"`
	TxSize     int64 `json:"tx_size_in_bytes"`
}

type NodeStatsThreadPoolPoolResponse struct {
	Threads   int64 `json:"threads"`
	Queue     int64 `json:"queue"`
	Active    int64 `json:"active"`
	Rejected  int64 `json:"rejected"`
	Largest   int64 `json:"largest"`
	Completed int64 `json:"completed"`
}

type NodeStatsTCPResponse struct {
	ActiveOpens  int64 `json:"active_opens"`
	PassiveOpens int64 `json:"passive_opens"`
	CurrEstab    int64 `json:"curr_estab"`
	InSegs       int64 `json:"in_segs"`
	OutSegs      int64 `json:"out_segs"`
	RetransSegs  int64 `json:"retrans_segs"`
	EstabResets  int64 `json:"estab_resets"`
	AttemptFails int64 `json:"attempt_fails"`
	InErrs       int64 `json:"in_errs"`
	OutRsts      int64 `json:"out_rsts"`
}

type NodeStatsIndicesResponse struct {
	Docs        NodeStatsIndicesDocsResponse        `json:"docs"`
	Store       NodeStatsIndicesStoreResponse       `json:"store"`
	Indexing    NodeStatsIndicesIndexingResponse    `json:"indexing"`
	Get         NodeStatsIndicesGetResponse         `json:"get"`
	Search      NodeStatsIndicesSearchResponse      `json:"search"`
	Merges      NodeStatsIndicesMergesResponse      `json:"merges"`
	Refresh     NodeStatsIndicesRefreshResponse     `json:"refresh"`
	Flush       NodeStatsIndicesFlushResponse       `json:"flush"`
	Warmer      NodeStatsIndicesWarmerResponse      `json:"warmer"`
	FilterCache NodeStatsIndicesFilterCacheResponse `json:"filter_cache"`
	IdCache     NodeStatsIndicesIdCacheResponse     `json:"id_cache"`
	FieldData   NodeStatsIndicesFieldDataResponse   `json:"fielddata"`
	Percolate   NodeStatsIndicesPercolateResponse   `json:"percolate"`
	Completion  NodeStatsIndicesCompletionResponse  `json:"completion"`
	Segments    NodeStatsIndicesSegmentsResponse    `json:"segments"`
	Translog    NodeStatsIndicesTranslogResponse    `json:"translog"`
	Suggest     NodeStatsIndicesSuggestResponse     `json:"suggest"`
}

type NodeStatsIndicesDocsResponse struct {
	Count   int64 `json:"count"`
	Deleted int64 `json:"deleted"`
}

type NodeStatsIndicesStoreResponse struct {
	Size         int64 `json:"size_in_bytes"`
	ThrottleTime int64 `json:"throttle_time_in_millis"`
}

type NodeStatsIndicesIndexingResponse struct {
	IndexTotal    int64 `json:"index_total"`
	IndexTime     int64 `json:"index_time_in_millis"`
	IndexCurrent  int64 `json:"index_current"`
	DeleteTotal   int64 `json:"delete_total"`
	DeleteTime    int64 `json:"delete_time_in_millis"`
	DeleteCurrent int64 `json:"delete_current"`
}

type NodeStatsIndicesGetResponse struct {
	Total        int64 `json:"total"`
	Time         int64 `json:"time_in_millis"`
	ExistsTotal  int64 `json:"exists_total"`
	ExistsTime   int64 `json:"exists_time_in_millis"`
	MissingTotal int64 `json:"missing_total"`
	MissingTime  int64 `json:"missing_time_in_millis"`
	Current      int64 `json:"current"`
}

type NodeStatsIndicesSearchResponse struct {
	OpenContext  int64 `json:"open_contexts"`
	QueryTotal   int64 `json:"query_total"`
	QueryTime    int64 `json:"query_time_in_millis"`
	QueryCurrent int64 `json:"query_current"`
	FetchTotal   int64 `json:"fetch_total"`
	FetchTime    int64 `json:"fetch_time_in_millis"`
	FetchCurrent int64 `json:"fetch_current"`
}
type NodeStatsIndicesMergesResponse struct {
	Current            int64 `json:"current"`
	CurrentDocs        int64 `json:"current_docs"`
	CurrentSizeInBytes int64 `json:"current_size_in_bytes"`
	Total              int64 `json:"total"`
	TotalTimeInMs      int64 `json:"total_time_in_millis"`
	TotalDocs          int64 `json:"total_docs"`
	TotalSizeInBytes   int64 `json:"total_size_in_bytes"`
}
type NodeStatsIndicesRefreshResponse struct {
	Total         int64 `json:"total"`
	TotalTimeInMs int64 `json:"total_time_in_millis"`
}
type NodeStatsIndicesFlushResponse struct {
	Total         int64 `json:"total"`
	TotalTimeInMs int64 `json:"total_time_in_millis"`
}
type NodeStatsIndicesWarmerResponse struct {
	Current       int64 `json:"current"`
	Total         int64 `json:"total"`
	TotalTimeInMs int64 `json:"total_time_in_millis"`
}
type NodeStatsIndicesFilterCacheResponse struct {
	MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
	Evictions         int64 `json:"evictions"`
}
type NodeStatsIndicesIdCacheResponse struct {
	MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
}
type NodeStatsIndicesFieldDataResponse struct {
	MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
	Evictions         int64 `json:"evictions"`
}
type NodeStatsIndicesPercolateResponse struct {
	Total             int64  `json:"total"`
	TimeInMs          int64  `json:"time_in_millis"`
	Current           int64  `json:"current"`
	MemorySizeInBytes int64  `json:"memory_size_in_bytes"`
	MemorySize        string `json:"memory_size"`
	Queries           int64  `json:"queries"`
}
type NodeStatsIndicesCompletionResponse struct {
	SizeInBytes int64 `json:"size_in_bytes"`
}
type NodeStatsIndicesSegmentsResponse struct {
	Count                    int64 `json:"count"`
	MemoryInBytes            int64 `json:"memory_in_bytes"`
	IndexWriterMemoryInBytes int64 `json:"index_writer_memory_in_bytes"`
	VersionMapMemoryInBytes  int64 `json:"version_map_memory_in_bytes"`
}
type NodeStatsIndicesTranslogResponse struct {
	Operations  int64 `json:"operations"`
	SizeInBytes int64 `json:"size_in_bytes"`
}
type NodeStatsIndicesSuggestResponse struct {
	Total    int64 `json:"total"`
	TimeInMs int64 `json:"time_in_millis"`
	Current  int64 `json:"current"`
}
type NodeStatsOSResponse struct {
	Timestamp int64                   `json:"timestamp"`
	Uptime    int64                   `json:"uptime_in_millis"`
	LoadAvg   []float64               `json:"load_average"`
	CPU       NodeStatsOSCPUResponse  `json:"cpu"`
	Mem       NodeStatsOSMemResponse  `json:"mem"`
	Swap      NodeStatsOSSwapResponse `json:"swap"`
}

type NodeStatsOSMemResponse struct {
	Free       int64 `json:"free_in_bytes"`
	Used       int64 `json:"used_in_bytes"`
	ActualFree int64 `json:"actual_free_in_bytes"`
	ActualUsed int64 `json:"actual_used_in_bytes"`
}

type NodeStatsOSSwapResponse struct {
	Used int64 `json:"used_in_bytes"`
	Free int64 `json:"free_in_bytes"`
}

type NodeStatsOSCPUResponse struct {
	Sys   int64 `json:"sys"`
	User  int64 `json:"user"`
	Idle  int64 `json:"idle"`
	Steal int64 `json:"stolen"`
}

type NodeStatsProcessResponse struct {
	Timestamp int64                       `json:"timestamp"`
	OpenFD    int64                       `json:"open_file_descriptors"`
	CPU       NodeStatsProcessCPUResponse `json:"cpu"`
	Memory    NodeStatsProcessMemResponse `json:"mem"`
}

type NodeStatsProcessMemResponse struct {
	Resident     int64 `json:"resident_in_bytes"`
	Share        int64 `json:"share_in_bytes"`
	TotalVirtual int64 `json:"total_virtual_in_bytes"`
}

type NodeStatsProcessCPUResponse struct {
	Percent int64 `json:"percent"`
	Sys     int64 `json:"sys_in_millis"`
	User    int64 `json:"user_in_millis"`
	Total   int64 `json:"total_in_millis"`
}

type NodeStatsJVMResponse struct {
	Timestame   int64                                      `json:"timestamp"`
	UptimeInMs  int64                                      `json:"uptime_in_millis"`
	Mem         NodeStatsJVMMemResponse                    `json:"mem"`
	Threads     NodeStatsJVMThreadsResponse                `json:"threads"`
	GC          NodeStatsJVMGCResponse                     `json:"gc"`
	BufferPools map[string]NodeStatsJVMBufferPoolsResponse `json:"buffer_pools"`
}

type NodeStatsJVMMemResponse struct {
	HeapUsedInBytes         int64                                   `json:"heap_used_in_bytes"`
	HeapUsedPercent         int64                                   `json:"heap_used_percent"`
	HeapCommitedInBytes     int64                                   `json:"heap_commited_in_bytes"`
	HeapMaxInBytes          int64                                   `json:"heap_max_in_bytes"`
	NonHeapUsedInBytes      int64                                   `json:"non_heap_used_in_bytes"`
	NonHeapCommittedInBytes int64                                   `json:"non_heap_committed_in_bytes"`
	Pools                   map[string]NodeStatsJVMMemPoolsResponse `json:"pools"`
}
type NodeStatsJVMMemPoolsResponse struct {
	UsedInBytes     int64 `json:"used_in_bytes"`
	MaxInBytes      int64 `json:"max_in_bytes"`
	PeakUsedInBytes int64 `json:"peak_used_in_bytes"`
	PeakMaxInBytes  int64 `json:"peak_max_in_bytes"`
}
type NodeStatsJVMThreadsResponse struct {
	Count     int64 `json:"count"`
	PeakCount int64 `json:"peak_count"`
}
type NodeStatsJVMGCResponse struct {
	Collectors map[string]NodeStatsJVMGCCollectorsAgeResponse `json:"collectors"`
}
type NodeStatsJVMGCCollectorsAgeResponse struct {
	Count    int64 `json:"collection_count"`
	TimeInMs int64 `json:"collection_time_in_millis"`
}
type NodeStatsJVMBufferPoolsResponse struct {
	Count                int64 `json:"count"`
	UsedInBytes          int64 `json:"used_in_bytes"`
	TotalCapacityInBytes int64 `json:"total_capacity_in_bytes"`
}
type NodeStatsHTTPResponse struct {
	CurrentOpen int64 `json:"current_open"`
	TotalOpen   int64 `json:"total_open"`
}

type NodeStatsFSResponse struct {
	Timestamp int64                     `json:"timestamp"`
	Data      []NodeStatsFSDataResponse `json:"data"`
}

type NodeStatsFSDataResponse struct {
	Path          string `json:"path"`
	Mount         string `json:"mount"`
	Device        string `json:"dev"`
	Total         int64  `json:"total_in_bytes"`
	Free          int64  `json:"free_in_bytes"`
	Available     int64  `json:"available_in_bytes"`
	DiskReads     int64  `json:"disk_reads"`
	DiskWrites    int64  `json:"disk_writes"`
	DiskReadSize  int64  `json:"disk_read_size_in_bytes"`
	DiskWriteSize int64  `json:"disk_write_size_in_bytes"`
}
