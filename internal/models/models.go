package models

type ZipInfo struct {
	Filename    string     `json:"filename"`
	ArchiveSize float64    `json:"archive_size"`
	TotalSize   float64    `json:"total_size"`
	TotalFiles  float64    `json:"total_files"`
	Files       []FileMeta `json:"files"`
}

type FileMeta struct {
	FilePath string  `json:"file_path"`
	Size     float64 `json:"size"`
	MimeType string  `json:"meme_type"`
}
