/*
	reboxed - the toybox server emulator
	Copyright (C) 2024  patapancakes <patapancakes@pagefault.games>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package browser

import (
	"fmt"
	"html/template"
	"net/http"
	"reboxed/db"
	"reboxed/utils"
	"strconv"
)

type BrowserTemplateData struct {
	InGame   bool
	Category string
	PageNum  int
	Packages []utils.Package
	PrevLink string
	NextLink string
}

const itemsPerPage = 50

var (
	categories = map[string]string{
		"entities": "entity",
		"weapons":  "weapon",
		"props":    "prop",
		"saves":    "savemap",
		"maps":     "map",
	}
	t, _ = template.New("Browser").Parse(tmpl)
)

func Handle(w http.ResponseWriter, r *http.Request) {
	category, ok := categories[r.PathValue("category")]
	if !ok {
		http.Error(w, "unknown category", http.StatusNotFound)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	list, err := db.FetchPackageListByTypePaged(category, (page-1)*itemsPerPage, itemsPerPage)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch package list: %s", err))
		return
	}

	prev := fmt.Sprintf("?page=%d", page-1)
	if page <= 1 {
		prev = "#"
	}

	next := fmt.Sprintf("?page=%d", page+1)
	if len(list) < itemsPerPage {
		next = "#"
	}

	err = t.Execute(w, BrowserTemplateData{
		InGame:   r.Header.Get("GMOD_VERSION") != "",
		Category: category,
		PageNum:  page,
		Packages: list,
		PrevLink: prev,
		NextLink: next,
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
