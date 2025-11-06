package repository

import "time"

// もともとopenapi.ymlに書いていて自動生成されていた構造体群
// openapi.ymlで自動生成されなくなったが、DB:minkan_statesの初期化でのみ使うため転記

// Edge defines model for Edge.
type Edge struct {
	Id string `json:"id"`

	// Source 接続元ノードID
	Source string `json:"source"`

	// Target 接続先ノードID
	Target string `json:"target"`
	Type   string `json:"type"`
}

// KanbanCardRef defines model for KanbanCardRef.
type KanbanCardRef struct {
	NodeId string `json:"nodeId"`
	PjId   string `json:"pjId"`
}

// KanbanColumns 各カラムごとのカード参照配列
type KanbanColumns struct {
	Backlog []KanbanCardRef `json:"backlog"`
	Doing   []KanbanCardRef `json:"doing"`
	Done    []KanbanCardRef `json:"done"`
	Todo    []KanbanCardRef `json:"todo"`
}

// KanbanIndex 本来 pjId -> nodeIdのSetだが、Setがないのでstring[]で代替
type KanbanIndex map[string][]string

// Minkan mindmap,kanban,作業中プロジェクトIDの状態
type Minkan struct {
	// CurrentPjId 作業中マインドマッププロジェクトのID
	CurrentPjId string `json:"currentPjId"`

	// KanbanColumns 各カラムごとのカード参照配列
	KanbanColumns KanbanColumns `json:"kanbanColumns"`

	// KanbanIndex 本来 pjId -> nodeIdのSetだが、Setがないのでstring[]で代替
	KanbanIndex KanbanIndex `json:"kanbanIndex"`

	// Projects pjID -> Project のマップ
	Projects Projects `json:"projects"`
}

// Node defines model for Node.
type Node struct {
	Data     NodeData `json:"data"`
	Id       string   `json:"id"`
	Position struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	} `json:"position"`
	Type string `json:"type"`
}

// NodeComment defines model for NodeComment.
type NodeComment struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Id        string    `json:"id"`
}

// NodeData defines model for NodeData.
type NodeData struct {
	Comments []NodeComment `json:"comments"`
	IsDone   bool          `json:"isDone"`
	Label    string        `json:"label"`
	ParentId *string       `json:"parentId"`
}

// Project defines model for Project.
type Project struct {
	CreatedAt time.Time `json:"createdAt"`
	Edges     []Edge    `json:"edges"`
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Nodes     []Node    `json:"nodes"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Projects pjID -> Project のマップ
type Projects map[string]Project
