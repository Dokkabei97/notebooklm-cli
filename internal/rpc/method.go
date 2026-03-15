package rpc

// RPC method IDs for Google NotebookLM batchexecute protocol.
// Reverse-engineered from notebooklm-py (teng-lin/notebooklm-py).

const (
	// Notebook operations
	MethodListNotebooks  = "wXbhsf"
	MethodCreateNotebook = "CCqFvf"
	MethodGetNotebook    = "rLM1Ne"
	MethodRenameNotebook = "s0tc2d"
	MethodDeleteNotebook = "WWINqb"

	// Source operations
	MethodAddSource            = "izAoDd"
	MethodAddSourceFile        = "o4cbdc" // Register uploaded file as source
	MethodDeleteSource         = "tGMBJ"
	MethodGetSource            = "hizoJc"
	MethodRefreshSource        = "FLmJqe"
	MethodCheckSourceFreshness = "yR9Yof"
	MethodUpdateSource         = "b7Wfje"
	MethodDiscoverSources      = "qXyaNe"

	// Summary and query
	MethodSummarize         = "VfAZjd"
	MethodGetSourceGuide    = "tr032e"
	MethodGetSuggestedReports = "ciyUvf"

	// Artifact operations
	MethodCreateArtifact    = "R7cb6c"
	MethodListArtifacts     = "gArtLc"
	MethodDeleteArtifact    = "V5N4be"
	MethodRenameArtifact    = "rc3d8d"
	MethodExportArtifact    = "Krh3pd"
	MethodShareArtifact     = "RGP97b"
	MethodGetInteractiveHTML = "v9rmvd"
	MethodReviseSlide       = "KmcKPe"

	// Research
	MethodStartFastResearch = "Ljjv0c"
	MethodStartDeepResearch = "QA9ei"
	MethodPollResearch      = "e3bVqc"
	MethodImportResearch    = "LBwxtb"

	// Note and mind map operations
	MethodGenerateMindMap      = "yyryJe"
	MethodCreateNote           = "CYK0Xb"
	MethodGetNotesAndMindMaps  = "cFji9"
	MethodUpdateNote           = "cYAfTb"
	MethodDeleteNote           = "AH0mwd"

	// Conversation
	MethodGetLastConversationID = "hPTbtc"
	MethodGetConversationTurns  = "khqZz"

	// Sharing operations
	MethodShareNotebook = "QDyure"
	MethodGetShareStatus = "JFMDGd"

	// Additional operations
	MethodRemoveRecentlyViewed = "fejl7e"
	MethodGetUserSettings      = "ZwVcOc"
	MethodSetUserSettings      = "hT54vc"
)

// API endpoint URLs
const (
	BatchExecuteURL = "https://notebooklm.google.com/_/LabsTailwindUi/data/batchexecute"
	QueryURL        = "https://notebooklm.google.com/_/LabsTailwindUi/data/google.internal.labs.tailwind.orchestration.v1.LabsTailwindOrchestrationService/GenerateFreeFormStreamed"
	UploadURL       = "https://notebooklm.google.com/upload/_/"
)

// ArtifactTypeCode represents artifact type codes used in RPC calls.
type ArtifactTypeCode int

const (
	ArtifactCodeAudio      ArtifactTypeCode = 1
	ArtifactCodeReport     ArtifactTypeCode = 2
	ArtifactCodeVideo      ArtifactTypeCode = 3
	ArtifactCodeQuiz       ArtifactTypeCode = 4
	ArtifactCodeMindMap    ArtifactTypeCode = 5
	ArtifactCodeInfographic ArtifactTypeCode = 7
	ArtifactCodeSlideDeck  ArtifactTypeCode = 8
	ArtifactCodeDataTable  ArtifactTypeCode = 9
)

// ArtifactStatusCode represents artifact processing status.
type ArtifactStatusCode int

const (
	ArtifactStatusProcessing ArtifactStatusCode = 1
	ArtifactStatusPending    ArtifactStatusCode = 2
	ArtifactStatusCompleted  ArtifactStatusCode = 3
	ArtifactStatusFailed     ArtifactStatusCode = 4
)

func (s ArtifactStatusCode) String() string {
	switch s {
	case ArtifactStatusProcessing:
		return "in_progress"
	case ArtifactStatusPending:
		return "pending"
	case ArtifactStatusCompleted:
		return "completed"
	case ArtifactStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// SourceStatusCode represents source processing status.
type SourceStatusCode int

const (
	SourceStatusProcessing SourceStatusCode = 1
	SourceStatusReady      SourceStatusCode = 2
	SourceStatusError      SourceStatusCode = 3
	SourceStatusPreparing  SourceStatusCode = 5
)

func (s SourceStatusCode) String() string {
	switch s {
	case SourceStatusProcessing:
		return "processing"
	case SourceStatusReady:
		return "ready"
	case SourceStatusError:
		return "error"
	case SourceStatusPreparing:
		return "preparing"
	default:
		return "unknown"
	}
}
