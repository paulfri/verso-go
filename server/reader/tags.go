package reader

import "net/http"

type Tag struct {
	ID     string `json:"id"`
	SortID string `json:"sortid"` // e.g. A0000001
}

type TagList struct {
	Tags []Tag `json:"tags"`
}

func (c *ReaderController) TagList(w http.ResponseWriter, req *http.Request) {
	c.Container.Render.JSON(w, http.StatusOK, TagList{
		Tags: []Tag{},
	})
}

// {
//     "tags": [
//         {
//             "id": "user/1/state/com.google/starred",
//             "sortid": "A0000001"

//         },
//         {
//             "id": "user/1/states/com.google/broadcast",
//             "sortid": "A0000002"

//	        },
//	        {
//	            "id": "user/1/label/Tech",
//	            "sortid": "A0000003"
//	        },
//	    ]
//	}
