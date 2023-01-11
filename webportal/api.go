package webportal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	//"encoding/json"
	"github.com/ip2location/ip2location-go/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/doublen987/Projects/MySite/server/functionality"
	"github.com/doublen987/Projects/MySite/server/persistence"
	"github.com/doublen987/Projects/MySite/server/persistence/models"
	"github.com/doublen987/Projects/MySite/server/webportal/template"
	"github.com/gorilla/mux"
)

type Suggestion struct {
	ID    string `bson:"_id",json:"id"`
	Title string `bson:"title",json:"title"`
}

type Image struct {
	FileName string `bson:"filename",json:"filename"`
	Bytes    []byte `bson:"bytes",json:"bytes"`
}

type ImageResponse struct {
	Location string `json:"location"`
}

type SectionRequest struct {
	SelectedPageID string        `json:"selectedPageID"`
	Sections       []interface{} `json:"sections"`
}

type TextSectionRequest struct {
	Header  string `json:"header"`
	Content string `json:"content"`
}

type StackSectionRequest struct {
	Header      string `json:"header"`
	TagSections string `json:"content"`
}

type ModelSectionRequest struct {
	Filename string `json:"filename"`
	Bytes    string `json:"bytes"`
}

type ImageSectionRequest struct {
	Filename string `json:"filename"`
	Bytes    string `json:"bytes"`
}

type SettingsReq struct {
	models.Settings
	Bytes    []byte `json:"bytes"`
	FileName string `json:"filename"`
}

type key string

const KeyAuthUserID key = "auth_user_id"

// type authenticationMiddleware struct {
// 	tokenUsers map[string]string
// }

// func (amw *authenticationMiddleware) Populate() {
// 	amw.tokenUsers["00000000"] = "user0"
// 	amw.tokenUsers["aaaaaaaa"] = "userA"
// 	amw.tokenUsers["05f717e5"] = "randomUser"
// 	amw.tokenUsers["deadbeef"] = "user0"
// }

func IPLoggerMiddleware(database persistence.DBHandler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB1.IPV6.BIN")
			if err != nil {
				fmt.Println(err)
				return
			}

			ips := strings.Split(r.RemoteAddr, ":")

			results, err := db.Get_all(ips[0])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("country_long: %s\n", results.Country_short)

			ctx := r.Context()

			visit := models.Visit{
				URL:     r.URL.Path,
				Country: strings.Trim(strings.ReplaceAll(results.Country_long, " ", ""), "."),
			}
			err = database.AddVisit(ctx, visit)
			if err != nil {
				fmt.Println(err)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")

		//If the request doesn't contain the authentication token handle the request without the authentication id in the context
		if err != nil {
			//next.ServeHTTP(w, r)
			//template.HandleLogin("", w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tknStr := c.Value

		if user, err := functionality.Authenticate(tknStr); err == nil {
			log.Printf("Authenticated user %s\n", user.Username)
			ctx := r.Context()
			ctx = context.WithValue(ctx, KeyAuthUserID, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func RunAPI(dbtype uint8, endpoint string, cert string, key string, tlsendpoint string, dbconnection string, databaseName, filestoragetype string) (chan error, chan error) {
	rootrouter := mux.NewRouter()
	db, err := persistence.GetDataBaseHandler(dbtype, dbconnection, databaseName)
	fh, err := persistence.GetFileHandler(filestoragetype, "")
	localfh, err := persistence.GetFileHandler("filesystem", "")
	settings, err := db.GetSettings(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	if users, err := db.GetUsers(context.Background()); len(users) == 0 {
		err = db.AddUser(context.Background(), models.User{Username: "admin", Password: "admin", Description: ""})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}

	r := rootrouter.PathPrefix("").Subrouter()

	r.Use(IPLoggerMiddleware(db))

	r.Path("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//template.Homepage("Hi there!", "My name is Nikola, I'm a full-stack web developer from Belgrade, Serbia. Welcome to my website, here you can find all sorts of information about me, projects I've worked on, and technologies I'm currently interested in!", w)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "homepage", true)
		homepage, err := db.GetPage(ctx, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		fmt.Println(homepage)
		page := models.Page2{}
		page.ID = homepage.ID
		page.Homepage = homepage.Homepage
		page.Name = homepage.Name
		for _, section := range homepage.Sections {
			bla, ok := section.(map[string]interface{})
			fmt.Println(ok)
			fmt.Println(bla)
			fmt.Println(section)

			if bla["type"] == "stack" {
				fmt.Println(reflect.TypeOf(bla["tagssections"]))
				header, _ := bla["name"].(string)
				newSection := models.StackSection{
					Name: header,
				}
				la, _ := bla["tagssections"].(primitive.A)
				for _, a := range la {
					fmt.Println(reflect.TypeOf(a))
					tagSection := a.(map[string]interface{})
					name := tagSection["name"].(string)
					fmt.Println(reflect.TypeOf(tagSection["tags"]))
					tags := tagSection["tags"].(primitive.A)
					newtags := []models.Tag{}
					for _, tag := range tags {
						tagMap := tag.(map[string]interface{})
						tagid := tagMap["ID"].(string)
						tagname := tagMap["name"].(string)
						tagthumbnail := tagMap["thumbnail"].(string)
						newtags = append(newtags, models.Tag{
							ID:        tagid,
							Name:      tagname,
							Thumbnail: tagthumbnail,
						})
					}
					newSection.TagSections = append(newSection.TagSections, models.TagSection{
						Name: name,
						Tags: newtags,
					})
				}
				page.Sections = append(page.Sections, newSection)
			}
			if bla["type"] == "text" {
				header, _ := bla["header"].(string)
				content, _ := bla["content"].(string)

				page.Sections = append(page.Sections, models.TextSection{
					Header:  header,
					Content: content,
				})
			}
			if bla["type"] == "image" {
				filename, _ := bla["filename"].(string)

				page.Sections = append(page.Sections, models.ImageSection{
					Image: filename,
				})
			}
			if bla["type"] == "3dmodel" {
				filename, _ := bla["filename"].(string)

				page.Sections = append(page.Sections, models.ModelSection{
					FileName: filename,
				})
			}
		}
		template.Homepage(settings, page, w)
	})

	r.PathPrefix("/login").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		template.HandleLogin(settings, "", w)
	})

	r.PathPrefix("/login").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		username := req.FormValue("username")
		password := req.FormValue("password")
		fmt.Println("Username: " + username + " Password: " + password)
		authenticated, err := db.Authenticate(ctx, username, password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			template.HandleLogin(settings, "Server error", w)
			return
		}

		if !authenticated {
			w.WriteHeader(http.StatusInternalServerError)
			template.HandleLogin(settings, "Server error", w)
			return
		}

		out, err := functionality.GenerateJWTToken(username, password)
		if err == functionality.ErrUsernamePasswordNotFound {
			w.WriteHeader(http.StatusBadRequest)
			template.HandleLogin(settings, "Wrong username or password", w)
			//http.Redirect(w, req, "/", http.StatusBadRequest)
			return
		} else if err == functionality.ErrCreatingJWT {
			w.WriteHeader(http.StatusInternalServerError)
			template.HandleLogin(settings, "Error occured, try again", w)
			//http.Redirect(w, req, "/", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   out.Token,
			Expires: out.ExpiresAt,
		})

		http.Redirect(w, req, "/", http.StatusFound)
	})

	r.PathPrefix("/logout").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		http.SetCookie(w, &http.Cookie{
			Name:   "token",
			Value:  "",
			MaxAge: -1,
		})

		http.Redirect(w, req, "/", http.StatusFound)
	})

	r.PathPrefix("/dashboard").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// defer req.Body.Close()
		// vars := mux.Vars(req)
		// postID := vars["postID"]
		// ctx := req.Context()
		// post, err := db.GetPost(ctx, postID)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }
		// ctx := req.Context()
		// users, err := db.GetUsers(ctx)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		ctx := req.Context()
		settings, err := db.GetSettings(ctx)

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// socialLinks := map[string]string{
		// 	"github":   "https://github.com/doublen987",
		// 	"linkedin": "https://www.linkedin.com/in/nikola-nesovic-24214219a/",
		// 	"email":    "doublen987@gmail.com",
		// }
		template.HandleEditSettings(settings, w)
	})))

	r.PathPrefix("/dashboard").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	})))

	r.PathPrefix("/stats").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		visits, err := db.GetVisits(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(visits)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})))

	r.PathPrefix("/settings").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		settings, err := db.GetSettings(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(settings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})))

	r.PathPrefix("/settings").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		settingsreq := SettingsReq{}
		err := json.NewDecoder(req.Body).Decode(&settingsreq)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(settingsreq.Bytes) != 0 {
			newFilename, err := fh.AddFile(settingsreq.Bytes, settingsreq.FileName)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			settingsreq.Logo = newFilename
		}

		ctx := req.Context()
		err = db.UpdateSettings(ctx, settingsreq.Settings)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		settings, err = db.GetSettings(ctx)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})))

	r.PathPrefix("/pages").Methods("GET").Handler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		pages, err := db.GetPages(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		json.NewEncoder(w).Encode(pages)
	}))

	r.PathPrefix("/homepage/edit").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// defer req.Body.Close()
		// vars := mux.Vars(req)
		// postID := vars["postID"]
		// ctx := req.Context()
		// post, err := db.GetPost(ctx, postID)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }
		// ctx := req.Context()
		// posts, err := db.GetPosts(ctx)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// tags := []models.Tag{}
		// tags = append(tags, models.Tag{
		// 	Name:      "MongoDB",
		// 	Thumbnail: "",
		// })
		tags, err := db.GetTags(req.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		sections := []models.PageSection{}
		sections = append(sections,
			models.TextSection{
				ID:      "bla",
				Header:  "About me",
				Content: "Hi, my name is Nikola and I am a software developer",
			},
			models.TextSection{
				ID:      "bla2",
				Header:  "Qualifications",
				Content: "Hi, my name is Nikola and I am a software developer",
			},
		)

		pages, err := db.GetPages(req.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		template.HandleEditHomePage(settings, "Hi, there", pages, sections, tags, w)
	})))

	r.PathPrefix("/homepage/edit").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		page := models.Page{}

		err := json.NewDecoder(req.Body).Decode(&page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//fmt.Println(reflect.TypeOf(page.Sections[0]["type"]))

		for index, section := range page.Sections {
			sectionMap := section.(map[string]interface{})
			st, _ := sectionMap["type"].(string)
			fmt.Println(st)
			switch st {
			case "3dmodel":
				{
					filename, _ := sectionMap["filename"].(string)
					bytesStr, _ := sectionMap["bytes"].(string)

					bytes, err := base64.StdEncoding.DecodeString(bytesStr)
					if err != nil {
						fmt.Println(err)
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					bla := make(map[string]string)
					if len(bytes) != 0 {
						newFileName, err := fh.AddFile(bytes, filename)
						if err != nil {
							fmt.Println(err)
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						bla["filename"] = newFileName
					} else {
						bla["filename"] = filename
					}

					bla["type"] = st
					page.Sections[index] = bla
				}
				break
			case "text":
				{
					header, _ := sectionMap["header"].(string)
					content, _ := sectionMap["content"].(string)
					fmt.Println(header)
					fmt.Println(content)
				}
				break
			case "stack":
				{
					name, _ := sectionMap["name"].(string)
					itagssections, _ := sectionMap["tagssections"].([]interface{})
					tagssections := []models.TagSection{}
					for _, itag := range itagssections {
						bla := itag.(map[string]interface{})

						tags := []models.Tag{}

						tagsSlice := bla["tags"].([]interface{})
						for _, tagInSlice := range tagsSlice {
							tagsMap := tagInSlice.(map[string]interface{})
							tags = append(tags, models.Tag{
								ID:        tagsMap["ID"].(string),
								Name:      tagsMap["name"].(string),
								Thumbnail: tagsMap["thumbnail"].(string),
							})
						}
						tagssections = append(tagssections, models.TagSection{
							Name: bla["name"].(string),
							Tags: tags,
						})

					}
					fmt.Println(name)
					fmt.Println(len(tagssections))

				}
				break
			case "image":
				{
					filename, _ := sectionMap["filename"].(string)
					bytesStr, _ := sectionMap["bytes"].(string)

					bytes, err := base64.StdEncoding.DecodeString(bytesStr)
					if err != nil {
						fmt.Println(err)
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					bla := make(map[string]string)
					if len(bytes) != 0 {
						newFileName, err := fh.AddFile(bytes, filename)
						if err != nil {
							fmt.Println(err)
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						bla["filename"] = newFileName
					} else {
						bla["filename"] = filename
					}

					bla["type"] = st
					page.Sections[index] = bla

					fmt.Println(filename)
					fmt.Println(len(bytes))
				}
				break
			default:
				{

				}
			}
		}

		fmt.Println(page.Sections...)

		ctx := req.Context()
		if page.ID == "" || page.ID == "None" {
			err := db.AddPage(ctx, page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := db.UpdatePage(ctx, page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		tags, err := db.GetTags(req.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		pages, err := db.GetPages(req.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		template.HandleEditHomePage(settings, "", pages, []models.PageSection{}, tags, w)

	})))

	r.PathPrefix("/homepage/edit").Methods("DELETE").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		page := models.Page{}

		err := json.NewDecoder(req.Body).Decode(&page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(page)

		ctx := req.Context()
		if page.ID == "" || page.ID == "None" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			err := db.RemovePage(ctx, page.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	})))

	r.PathPrefix("/tags/edit").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// defer req.Body.Close()
		// vars := mux.Vars(req)
		// postID := vars["postID"]
		// ctx := req.Context()
		// post, err := db.GetPost(ctx, postID)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }
		ctx := req.Context()
		tags := []models.Tag{}

		tags, err := db.GetTags(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		template.HandleEditTag(settings, tags, w)
	})))

	r.PathPrefix("/tags/edit").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		method := req.FormValue("Send")

		if method == "POST" {
			tag := models.Tag{}
			// if err := json.NewDecoder(req.Body).Decode(post); err != nil {
			// 	http.Error(w, err.Error(), http.StatusBadRequest)
			// 	return
			// }

			//1. parse input, type multipart/form-data
			//ParseMultipartForm parses a request body as multipart/form-data. The whole request body is parsed and up to a
			//total of maxMemory bytes of its file parts are stored in memory, with the remainder stored on disk in temporary
			//files. ParseMultipartForm calls ParseForm if necessary. After one call to ParseMultipartForm, subsequent calls
			//have no effect.
			req.ParseMultipartForm(10 << 20)

			//2. retrieve file from posted form-data
			var fileName string = ""
			file, handler, err := req.FormFile("Thumbnail")

			if err != nil {
				fmt.Println("Error Retrieving thumbnail from request:")
				fmt.Println(err)
				fileName = req.FormValue("ThumbnailName")
			} else {
				defer file.Close()
				fmt.Printf("Uploaded File: %+v\n", handler.Filename)
				fmt.Printf("File Size: %+v\n", handler.Size)
				fmt.Printf("MIME Header: %+v\n", handler.Header)
				fileBytes, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println("Error Reading file bytes:")
					fmt.Println(err)
				} else {
					fileName, err = fh.AddFile(fileBytes, handler.Filename)
					if err != nil {
						fmt.Println("Error Adding File:")
						fmt.Println(err)
					}
				}

			}

			selectedTagId := req.FormValue("SelectedTag")
			tag.Name = req.FormValue("Name")
			tag.Description = req.FormValue("Description")
			tag.Content = req.FormValue("Content")
			tag.Thumbnail = fileName

			if req.FormValue("ThumbnailStretched") == "true" {
				tag.ThumbnailStretched = true
			} else {
				tag.ThumbnailStretched = false
			}

			ctx := req.Context()
			if selectedTagId == "" || selectedTagId == "None" {
				err := db.AddTag(ctx, tag)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				tag.ID = selectedTagId
				err := db.UpdateTag(ctx, tag)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			tags, err := db.GetTags(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(tag.Name)
			fmt.Println(tag.Content)

			template.HandleEditTag(settings, tags, w)
		}
		if method == "DELETE" {
			selectedTagId := req.FormValue("SelectedTag")
			fmt.Printf("Deleting post: $s\n", selectedTagId)

			ctx := req.Context()
			if selectedTagId != "" && selectedTagId != "None" {
				err := db.RemoveTag(ctx, selectedTagId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			tags, err := db.GetTags(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Printf("Deleted tag: $s\n", selectedTagId)
			template.HandleEditTag(settings, tags, w)
		}

	})))

	r.PathPrefix("/tags").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		tags, err := db.GetTags(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(tags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})))

	r.PathPrefix("/users/edit").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// defer req.Body.Close()
		// vars := mux.Vars(req)
		// postID := vars["postID"]
		// ctx := req.Context()
		// post, err := db.GetPost(ctx, postID)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }
		ctx := req.Context()
		users, err := db.GetUsers(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		template.HandleEditUsers(settings, users, w)
	})))

	r.PathPrefix("/users/edit").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		method := req.FormValue("Send")
		ctx := req.Context()

		if method == "POST" {
			user := models.User{}
			users := []models.User{}
			// if err := json.NewDecoder(req.Body).Decode(post); err != nil {
			// 	http.Error(w, err.Error(), http.StatusBadRequest)
			// 	return
			// }

			//1. parse input, type multipart/form-data
			//ParseMultipartForm parses a request body as multipart/form-data. The whole request body is parsed and up to a
			//total of maxMemory bytes of its file parts are stored in memory, with the remainder stored on disk in temporary
			//files. ParseMultipartForm calls ParseForm if necessary. After one call to ParseMultipartForm, subsequent calls
			//have no effect.
			req.ParseMultipartForm(10 << 20)

			//2. retrieve file from posted form-data
			var fileName string = ""
			file, handler, err := req.FormFile("Thumbnail")

			oldPassword := req.FormValue("OldPassword")
			selectedUserId := req.FormValue("SelectedUser")
			user.Username = req.FormValue("Username")
			user.Password = req.FormValue("Password")
			user.Description = req.FormValue("Description")

			if err != nil {
				fmt.Println("Error Retrieving thumbnail from request:")
				fmt.Println(err)
				fileName = req.FormValue("ThumbnailName")
			} else {
				defer file.Close()
				fmt.Printf("Uploaded File: %+v\n", handler.Filename)
				fmt.Printf("File Size: %+v\n", handler.Size)
				fmt.Printf("MIME Header: %+v\n", handler.Header)
				fileBytes, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println("Error Reading file bytes:")
					fmt.Println(err)
				} else {
					fileName, err = fh.AddFile(fileBytes, handler.Filename)
					if err != nil {
						fmt.Println("Error Adding File:")
						fmt.Println(err)
					}
				}
			}

			user.Thumbnail = fileName

			if selectedUserId == "" || selectedUserId == "None" {
				err := db.AddUser(ctx, user)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				if user.Password != "" {
					if authenticated, err := db.Authenticate(ctx, user.Username, oldPassword); authenticated != true || err != nil {
						http.Error(w, "Error confirming your old password", http.StatusBadRequest)
						return
					}
				}

				user.ID = selectedUserId
				err := db.UpdateUser(ctx, user)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			users, err = db.GetUsers(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(user.Username)
			fmt.Println(user.Password)

			template.HandleEditUsers(settings, users, w)
		}
		if method == "DELETE" {
			selectedUserId := req.FormValue("SelectedUser")
			fmt.Printf("Deleting user: $s\n", selectedUserId)

			ctx := req.Context()
			if selectedUserId != "" && selectedUserId != "None" {
				err := db.RemoveUser(ctx, selectedUserId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			users, err := db.GetUsers(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Printf("Deleted user: $s\n", selectedUserId)
			template.HandleEditUsers(settings, users, w)
		}

	})))

	r.PathPrefix("/blog/edit").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// defer req.Body.Close()
		// vars := mux.Vars(req)
		// postID := vars["postID"]
		// ctx := req.Context()
		// post, err := db.GetPost(ctx, postID)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }
		tags, err := db.GetTags(req.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		ctx := req.Context()
		posts, err := db.GetPosts(ctx)
		fmt.Println(posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		template.HandleEditPost(settings, posts, tags, w)
	})))

	r.PathPrefix("/blog/edit").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		method := req.FormValue("Send")

		if method == "POST" {
			post := models.Post{}
			posts := []models.Post{}
			// if err := json.NewDecoder(req.Body).Decode(post); err != nil {
			// 	http.Error(w, err.Error(), http.StatusBadRequest)
			// 	return
			// }

			//1. parse input, type multipart/form-data
			//ParseMultipartForm parses a request body as multipart/form-data. The whole request body is parsed and up to a
			//total of maxMemory bytes of its file parts are stored in memory, with the remainder stored on disk in temporary
			//files. ParseMultipartForm calls ParseForm if necessary. After one call to ParseMultipartForm, subsequent calls
			//have no effect.
			req.ParseMultipartForm(10 << 20)

			tags := []models.Tag{}
			noMoreTags := false
			for i := 0; !noMoreTags; i++ {
				tag := req.FormValue("tag-" + strconv.FormatInt(int64(i), 10))

				if tag == "" {
					noMoreTags = true
					continue
				}
				tags = append(tags, models.Tag{ID: tag})
			}
			post.Tags = tags

			//2. retrieve file from posted form-data
			var fileName string = ""
			file, handler, err := req.FormFile("Thumbnail")

			if err != nil {
				fmt.Println("Error Retrieving thumbnail from request:")
				fmt.Println(err)
				fileName = req.FormValue("ThumbnailName")
			} else {
				defer file.Close()
				fmt.Printf("Uploaded File: %+v\n", handler.Filename)
				fmt.Printf("File Size: %+v\n", handler.Size)
				fmt.Printf("MIME Header: %+v\n", handler.Header)
				fileBytes, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println("Error Reading file bytes:")
					fmt.Println(err)
				} else {
					fileName, err = fh.AddFile(fileBytes, handler.Filename)
					if err != nil {
						fmt.Println("Error Adding File:")
						fmt.Println(err)
					}
				}

			}

			selectedPostId := req.FormValue("SelectedPost")
			post.Title = req.FormValue("Title")
			post.Description = req.FormValue("Description")
			post.Content = req.FormValue("Content")
			if req.FormValue("Hidden") == "true" {
				post.Hidden = true
			} else {
				post.Hidden = false
			}
			post.Thumbnail = fileName

			if req.FormValue("ThumbnailStretched") == "true" {
				post.ThumbnailStretched = true
			} else {
				post.ThumbnailStretched = false
			}

			ctx := req.Context()
			if selectedPostId == "" || selectedPostId == "None" {
				err := db.AddPost(ctx, post)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				post.ID = selectedPostId
				_, err := db.UpdatePost(ctx, post)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			alltags, err := db.GetTags(req.Context())
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			posts, err = db.GetPosts(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(post.Title)
			fmt.Println(post.Content)

			template.HandleEditPost(settings, posts, alltags, w)
		}
		if method == "DELETE" {
			selectedPostId := req.FormValue("SelectedPost")
			fmt.Printf("Deleting post: $s\n", selectedPostId)

			ctx := req.Context()
			if selectedPostId != "" && selectedPostId != "None" {
				err := db.RemovePost(ctx, selectedPostId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			tags, err := db.GetTags(req.Context())
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			posts, err := db.GetPosts(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Printf("Deleted post: $s\n", selectedPostId)
			template.HandleEditPost(settings, posts, tags, w)
		}
		if method == "PUBLISH" {
			selectedPostId := req.FormValue("SelectedPost")
			fmt.Printf("Publishing post: $s\n", selectedPostId)

			ctx := req.Context()
			if selectedPostId != "" && selectedPostId != "None" {
				err := db.PublishPost(ctx, selectedPostId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			tags, err := db.GetTags(req.Context())
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			posts, err := db.GetPosts(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Printf("Published post: $s\n", selectedPostId)
			template.HandleEditPost(settings, posts, tags, w)
		}

	})))

	r.PathPrefix("/projects/edit").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// defer req.Body.Close()
		// vars := mux.Vars(req)
		// postID := vars["postID"]
		// ctx := req.Context()
		// post, err := db.GetPost(ctx, postID)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }
		ctx := req.Context()
		tags, err := db.GetTags(req.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		projects, err := db.GetProjects(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		template.HandleEditProject(settings, projects, tags, w)
	})))

	r.PathPrefix("/projects/edit").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		method := req.FormValue("Send")

		if method == "POST" {
			project := models.Project{}
			projects := []models.Project{}
			// if err := json.NewDecoder(req.Body).Decode(post); err != nil {
			// 	http.Error(w, err.Error(), http.StatusBadRequest)
			// 	return
			// }

			//1. parse input, type multipart/form-data
			//ParseMultipartForm parses a request body as multipart/form-data. The whole request body is parsed and up to a
			//total of maxMemory bytes of its file parts are stored in memory, with the remainder stored on disk in temporary
			//files. ParseMultipartForm calls ParseForm if necessary. After one call to ParseMultipartForm, subsequent calls
			//have no effect.
			req.ParseMultipartForm(10 << 20)

			projecttags := []models.Tag{}
			noMoreTags := false
			for i := 0; !noMoreTags; i++ {
				tag := req.FormValue("tag-" + strconv.FormatInt(int64(i), 10))

				if tag == "" {
					noMoreTags = true
					continue
				}
				projecttags = append(projecttags, models.Tag{ID: tag})
			}
			project.Tags = projecttags

			//2. retrieve file from posted form-data
			var fileName string = ""
			file, handler, err := req.FormFile("Thumbnail")
			if err != nil {
				fmt.Println("Error Retrieving file from form-data")
				fileName = req.FormValue("ThumbnailName")
			} else {
				defer file.Close()
				fmt.Printf("Uploaded File: %+v\n", handler.Filename)
				fmt.Printf("File Size: %+v\n", handler.Size)
				fmt.Printf("MIME Header: %+v\n", handler.Header)
				fileBytes, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println("Error Reading file bytes:")
					fmt.Println(err)
				} else {
					fileName, err = fh.AddFile(fileBytes, handler.Filename)
					if err != nil {
						fmt.Println("Error Adding File:")
						fmt.Println(err)
					}
				}
			}

			selectedProjectId := req.FormValue("SelectedProject")
			project.Title = req.FormValue("Title")
			project.Description = req.FormValue("Description")
			project.Link = req.FormValue("Link")
			project.Thumbnail = fileName
			fmt.Println(req.FormValue("ThumbnailStretched"))
			if req.FormValue("ThumbnailStretched") == "true" {
				project.ThumbnailStretched = true
			} else {
				project.ThumbnailStretched = false
			}

			ctx := req.Context()
			if selectedProjectId == "" || selectedProjectId == "None" {
				err := db.AddProject(ctx, project)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				project.ID = selectedProjectId
				_, err := db.UpdateProject(ctx, project)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			projects, err = db.GetProjects(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(project.Title)
			fmt.Println(project.Description)

			tags, err := db.GetTags(req.Context())
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			template.HandleEditProject(settings, projects, tags, w)
		}
		if method == "DELETE" {
			selectedProjectId := req.FormValue("SelectedProject")
			fmt.Printf("Deleting project: $s\n", selectedProjectId)

			ctx := req.Context()
			if selectedProjectId != "" && selectedProjectId != "None" {
				err := db.RemoveProject(ctx, selectedProjectId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			projects, err := db.GetProjects(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Printf("Deleted project: $s\n", selectedProjectId)

			tags, err := db.GetTags(req.Context())
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			template.HandleEditProject(settings, projects, tags, w)
		}

	})))

	r.PathPrefix("/knowledgetimeline/edit").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// defer req.Body.Close()
		// vars := mux.Vars(req)
		// postID := vars["postID"]
		// ctx := req.Context()
		// post, err := db.GetPost(ctx, postID)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }
		ctx := req.Context()
		knowledgeTimeline, err := db.GetKnowledgeTimelineEvents(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		template.HandleEditKnowledgeTimeline(settings, knowledgeTimeline, w)
	})))

	r.PathPrefix("/knowledgetimeline/edit").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		method := req.FormValue("Send")

		if method == "POST" {
			event := models.TimelineEvent{}
			events := []models.TimelineEvent{}
			// if err := json.NewDecoder(req.Body).Decode(post); err != nil {
			// 	http.Error(w, err.Error(), http.StatusBadRequest)
			// 	return
			// }

			//1. parse input, type multipart/form-data
			//ParseMultipartForm parses a request body as multipart/form-data. The whole request body is parsed and up to a
			//total of maxMemory bytes of its file parts are stored in memory, with the remainder stored on disk in temporary
			//files. ParseMultipartForm calls ParseForm if necessary. After one call to ParseMultipartForm, subsequent calls
			//have no effect.
			req.ParseMultipartForm(10 << 20)

			//2. retrieve file from posted form-data
			file, handler, err := req.FormFile("Image")
			if err != nil {
				fmt.Println("Error Retrieving file from form-data")
				fmt.Println(err)
				return
			}
			defer file.Close()
			fmt.Printf("Uploaded File: %+v\n", handler.Filename)
			fmt.Printf("File Size: %+v\n", handler.Size)
			fmt.Printf("MIME Header: %+v\n", handler.Header)
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}

			fileName, err := fh.AddFile(fileBytes, handler.Filename)
			if err != nil {
				fmt.Println(err)
			}

			selectedEventId := req.FormValue("SelectedEvent")
			event.Title = req.FormValue("Title")
			event.Description = req.FormValue("Description")
			event.Image = fileName

			ctx := req.Context()
			if selectedEventId == "" || selectedEventId == "None" {
				err := db.AddKnowledgeTimelineEvent(ctx, event)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				event.ID = selectedEventId
				_, err := db.UpdateKnowledgeTimelineEvent(ctx, event)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			events, err = db.GetKnowledgeTimelineEvents(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(event.Title)
			fmt.Println(event.Description)

			template.HandleEditKnowledgeTimeline(settings, events, w)
		}
		if method == "DELETE" {
			selectedEventId := req.FormValue("SelectedEvent")
			fmt.Printf("Deleting knowledge timeline event: $s\n", selectedEventId)

			ctx := req.Context()
			if selectedEventId != "" && selectedEventId != "None" {
				err := db.RemoveKnowledgeTimelineEvent(ctx, selectedEventId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			events, err := db.GetKnowledgeTimelineEvents(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Printf("Deleted knowledge timeline event: $s\n", selectedEventId)
			template.HandleEditKnowledgeTimeline(settings, events, w)
		}

	})))

	r.PathPrefix("/projects").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		projects, err := db.GetProjects(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		fmt.Println(projects)
		template.HandleShowcase(settings, projects, w)
	})

	r.PathPrefix("/blog/{postID}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		ctx := req.Context()

		ctx = context.WithValue(ctx, "published", true)
		ctx = context.WithValue(ctx, "hidden", false)
		post, err := db.GetPost(ctx, vars["postID"])
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		links, err := db.GetLinks(ctx)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		template.HandlePost(settings, links, post, w)
	})

	r.PathPrefix("/blog").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//posts, err := db.GetPosts()
		var posts []models.Post
		// posts = []models.Post{
		// 	models.Post{
		// 		Title:   "WebAssembly",
		// 		Content: "WebAssembly enables us to run native C++ code in the browser",
		// 		Image:   "",
		// 	},
		// 	models.Post{
		// 		Title:   "React",
		// 		Content: "React is a Javascript framework that allows us to utilise functional programing in our apps",
		// 		Image:   "",
		// 	},
		// 	models.Post{
		// 		Title:   "Creating filters with Pixi.js",
		// 		Content: "Pixi filters allow us to generate all kinds of awesome effects.",
		// 		Image:   "",
		// 	},
		// }

		newCtx := req.Context()

		newCtx = context.WithValue(newCtx, "published", true)
		newCtx = context.WithValue(newCtx, "hidden", false)

		if searchTerm := req.URL.Query().Get("search"); searchTerm != "" {
			newCtx = context.WithValue(newCtx, "search-term", searchTerm)
		} else {
			newCtx = context.WithValue(newCtx, "search-term", "")
		}

		var currentPage int = 1

		if page := req.URL.Query().Get("page"); page != "" {
			newCtx = context.WithValue(newCtx, "page", page)
			currentPage, err = strconv.Atoi(page)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			newCtx = context.WithValue(newCtx, "page", "1")
		}

		postsCount, err := db.GetPostsCount(newCtx)
		if err != nil {
			fmt.Println(err)
			return
		}

		pageBlockStart := (currentPage/5)*5 + 1
		numOfPages := (postsCount / 10) + 1
		if postsCount%10 == 0 {
			numOfPages--
		}

		if currentPage < 1 {
			http.Redirect(w, req, "/blog?page=1", http.StatusSeeOther)
			return
		}
		if currentPage > int(numOfPages) && int(numOfPages) != 0 {
			http.Redirect(w, req, fmt.Sprintf("/blog?page=%d", numOfPages), http.StatusSeeOther)
			return
		}

		posts, err = db.GetPosts(newCtx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		template.HandleBlog(settings, posts, currentPage, pageBlockStart, int(numOfPages), w)
	})

	r.PathPrefix("/knowledgetimeline").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		var timeline []models.TimelineEvent
		timeline, err = db.GetKnowledgeTimelineEvents(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		template.HandleKnowledgeTimeline(settings, timeline, w)
	})

	//r.PathPrefix("/content/").Handler(http.StripPrefix("/content/", http.FileServer(http.Dir("./webportal/content"))))

	subr := rootrouter.PathPrefix("/content").Subrouter()

	subr.PathPrefix("/images/{imageID}").Methods("GET").Handler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		//ctx := req.Context()

		storagetype := req.URL.Query().Get("storagetype")

		var imageBytes []byte

		if storagetype == "filesystem" {
			imageBytes, err = localfh.GetFile(vars["imageID"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		} else {

			imageBytes, err = fh.GetFile(vars["imageID"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}

		s := strings.Split(vars["imageID"], ".")
		if len(s) < 1 {
			http.Error(w, "Invalid file name", http.StatusNotFound)
			return
		}

		fileextension := s[len(s)-1]

		switch fileextension {
		case "svg":
			w.Header().Set("content-type", "image/svg+xml")
			break
		default:
			fmt.Println(fileextension)
			w.Header().Set("content-type", "image")
		}

		w.Write(imageBytes)
		w.WriteHeader(http.StatusOK)
		return
	}))

	subr.PathPrefix("/images").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//vars := mux.Vars(req)
		//ctx := req.Context()

		image := Image{}
		err := json.NewDecoder(req.Body).Decode(&image)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		newFilename, err := fh.AddFile(image.Bytes, image.FileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)

		imageResponse := &ImageResponse{
			Location: newFilename,
		}

		json.NewEncoder(w).Encode(imageResponse)
	})))

	subr.PathPrefix("/").Handler(http.StripPrefix("/content/", http.FileServer(http.Dir("./webportal/content"))))

	httpErrChan := make(chan error)
	httpIsErrChan := make(chan error)

	fmt.Println("Cert.pem location: " + cert)
	fmt.Println("Key.pem location: " + key)

	go func() { httpIsErrChan <- http.ListenAndServeTLS(tlsendpoint, cert, key, rootrouter) }()
	go func() { httpErrChan <- http.ListenAndServe(endpoint, rootrouter) }()

	return httpErrChan, httpIsErrChan
}
