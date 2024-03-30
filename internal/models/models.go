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
	AntivirusStatus struct {
	} `json:"antivirus_status"`
	ResourceID string `json:"resource_id"`
	Share      struct {
		IsRoot  bool   `json:"is_root"`
		IsOwned bool   `json:"is_owned"`
		Rights  string `json:"rights"`
	} `json:"share"`
	File           string    `json:"file"`
	Size           int       `json:"size"`
	PhotosliceTime time.Time `json:"photoslice_time"`
	Embedded       struct {
		Sort  string `json:"sort"`
		Items []struct {
		} `json:"items"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
		Path   string `json:"path"`
		Total  int    `json:"total"`
	} `json:"_embedded"`
	Exif struct {
		DateTime     time.Time `json:"date_time"`
		GpsLongitude struct {
		} `json:"gps_longitude"`
		GpsLatitude struct {
		} `json:"gps_latitude"`
	} `json:"exif"`
	CustomProperties struct {
	} `json:"custom_properties"`
	MediaType string    `json:"media_type"`
	Preview   string    `json:"preview"`
	Type      string    `json:"type"`
	MimeType  string    `json:"mime_type"`
	Revision  int       `json:"revision"`
	PublicURL string    `json:"public_url"`
	Path      string    `json:"path"`
	Md5       string    `json:"md5"`
	PublicKey string    `json:"public_key"`
	Sha256    string    `json:"sha256"`
	Name      string    `json:"name"`
	Created   time.Time `json:"created"`
	Sizes     []struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"sizes"`
	Modified   time.Time `json:"modified"`
	CommentIds struct {
		PrivateResource string `json:"private_resource"`
		PublicResource  string `json:"public_resource"`
	} `json:"comment_ids"`
}
