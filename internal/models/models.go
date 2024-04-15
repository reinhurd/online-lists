package models

import "time"

type YADISKList struct {
	Items  []YDItem `json:"items"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

type YDUploadResponse struct {
	OperationID string `json:"operation_id"`
	Href        string `json:"href"`
	Method      string `json:"method"`
	Templated   bool   `json:"templated"`
}

type YDItem struct {
	AntivirusStatus string `json:"antivirus_status"`
	Size            int    `json:"size"`
	CommentIds      struct {
		PrivateResource string `json:"private_resource"`
		PublicResource  string `json:"public_resource"`
	} `json:"comment_ids"`
	Name string `json:"name"`
	Exif struct {
	} `json:"exif"`
	Created    time.Time `json:"created"`
	ResourceID string    `json:"resource_id"`
	Modified   time.Time `json:"modified"`
	MimeType   string    `json:"mime_type"`
	Sizes      []struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"sizes,omitempty"`
	File           string    `json:"file"`
	MediaType      string    `json:"media_type"`
	Preview        string    `json:"preview,omitempty"`
	Path           string    `json:"path"`
	Sha256         string    `json:"sha256"`
	Type           string    `json:"type"`
	Md5            string    `json:"md5"`
	Revision       int64     `json:"revision"`
	PublicKey      string    `json:"public_key,omitempty"`
	PublicURL      string    `json:"public_url,omitempty"`
	PhotosliceTime time.Time `json:"photoslice_time,omitempty"`
}
