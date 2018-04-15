package models

type Mountain struct {
  ID          int64   `json:"id"`
  Webid       string  `json:"web_id"`
  Name        string  `json:"title"`
  Lat         string  `json:"latitude"`
  Longt       string  `json:"longtitude"`
  Height      string  `json:"height"`
  Linktype    string  `json:"type_link"`
}
