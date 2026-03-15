package model

// SourceType represents the type of a source (from API type codes).
type SourceType int

const (
	SourceTypeUnknown          SourceType = 0
	SourceTypeGoogleDocs       SourceType = 1
	SourceTypeGoogleSlides     SourceType = 2
	SourceTypePDF              SourceType = 3
	SourceTypePastedText       SourceType = 4
	SourceTypeWebPage          SourceType = 5
	SourceTypeMarkdown         SourceType = 8
	SourceTypeYouTube          SourceType = 9
	SourceTypeMedia            SourceType = 10
	SourceTypeDocx             SourceType = 11
	SourceTypeImage            SourceType = 13
	SourceTypeGoogleSheets     SourceType = 14
	SourceTypeCSV              SourceType = 16
)

func (s SourceType) String() string {
	switch s {
	case SourceTypeGoogleDocs:
		return "google_docs"
	case SourceTypeGoogleSlides:
		return "google_slides"
	case SourceTypePDF:
		return "pdf"
	case SourceTypePastedText:
		return "pasted_text"
	case SourceTypeWebPage:
		return "web_page"
	case SourceTypeMarkdown:
		return "markdown"
	case SourceTypeYouTube:
		return "youtube"
	case SourceTypeMedia:
		return "media"
	case SourceTypeDocx:
		return "docx"
	case SourceTypeImage:
		return "image"
	case SourceTypeGoogleSheets:
		return "google_sheets"
	case SourceTypeCSV:
		return "csv"
	default:
		return "unknown"
	}
}

// ArtifactType represents the type of an artifact (maps to ArtifactTypeCode).
type ArtifactType int

const (
	ArtifactTypeUnknown     ArtifactType = 0
	ArtifactTypeAudio       ArtifactType = 1
	ArtifactTypeReport      ArtifactType = 2
	ArtifactTypeVideo       ArtifactType = 3
	ArtifactTypeQuiz        ArtifactType = 4
	ArtifactTypeMindMap     ArtifactType = 5
	ArtifactTypeInfographic ArtifactType = 7
	ArtifactTypeSlideDeck   ArtifactType = 8
	ArtifactTypeDataTable   ArtifactType = 9
)

func (a ArtifactType) String() string {
	switch a {
	case ArtifactTypeAudio:
		return "audio"
	case ArtifactTypeReport:
		return "report"
	case ArtifactTypeVideo:
		return "video"
	case ArtifactTypeQuiz:
		return "quiz"
	case ArtifactTypeMindMap:
		return "mind_map"
	case ArtifactTypeInfographic:
		return "infographic"
	case ArtifactTypeSlideDeck:
		return "slide_deck"
	case ArtifactTypeDataTable:
		return "data_table"
	default:
		return "unknown"
	}
}

// SharePermission represents sharing permission levels.
type SharePermission int

const (
	SharePermissionNone   SharePermission = 0
	SharePermissionViewer SharePermission = 3
	SharePermissionEditor SharePermission = 2
	SharePermissionOwner  SharePermission = 1
)

func (s SharePermission) String() string {
	switch s {
	case SharePermissionOwner:
		return "owner"
	case SharePermissionEditor:
		return "editor"
	case SharePermissionViewer:
		return "viewer"
	default:
		return "none"
	}
}
