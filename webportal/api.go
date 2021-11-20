package webportal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//"encoding/json"

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

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println()
		c, err := r.Cookie("token")

		//If the request doesn't contain the authentication token handle the request without the authentication id in the context
		if err != nil {
			next.ServeHTTP(w, r)
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

func RunAPI(dbtype uint8, addr string, dbconnection string, filestoragetype string) error {
	r := mux.NewRouter()
	db, err := persistence.GetDataBaseHandler(dbtype, dbconnection)
	fh, err := persistence.GetFileHandler(filestoragetype, "")
	if err != nil {
		log.Fatal(err)
	}

	r.Path("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		template.Homepage("Hi there!", "My name is Nikola, I'm a full-stack web developer from Belgrade, Serbia. Welcome to my website, here you can find all sorts of information about me, projects I've worked on, and technologies I'm currently interested in!", w)
	})

	r.PathPrefix("/login").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		template.HandleLogin("", w)
	})

	r.PathPrefix("/login").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		username := req.FormValue("username")
		password := req.FormValue("password")

		out, err := functionality.Login(username, password)
		if err == functionality.ErrUsernamePasswordNotFound {
			w.WriteHeader(http.StatusBadRequest)
			template.HandleLogin("Wrong username or password", w)
			//http.Redirect(w, req, "/", http.StatusBadRequest)
			return
		} else if err == functionality.ErrCreatingJWT {
			w.WriteHeader(http.StatusInternalServerError)
			template.HandleLogin("Error occured, try again", w)
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

	// r.PathPrefix("/blog/edit/{postID}").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// 	defer req.Body.Close()
	// 	vars := mux.Vars(req)
	// 	postID := vars["postID"]
	// 	ctx := req.Context()
	// 	post, err := db.GetPost(ctx, postID)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusNotFound)
	// 		return
	// 	}
	// 	template.HandleEditPost(post, w)
	// })))

	// r.PathPrefix("/blog/edit/{postID}").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// 	defer req.Body.Close()
	// 	vars := mux.Vars(req)
	// 	post := models.Post{}
	// 	if err := json.NewDecoder(req.Body).Decode(post); err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}

	// 	ctx := req.Context()
	// 	post.ID = vars["postID"]
	// 	newPost, err := db.UpdatePost(ctx, post)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	template.HandleEditPost(newPost, w)
	// })))

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
		ctx := req.Context()
		posts, err := db.GetPosts(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		template.HandleEditPost(posts, w)
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
			post.Thumbnail = fileName

			ctx := req.Context()
			if selectedPostId == "" || selectedPostId == "None" {
				err := db.AddPost(ctx, post)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				post.ID = selectedPostId
				err := db.ReplacePost(ctx, post)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			posts, err = db.GetPosts(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(post.Title)
			fmt.Println(post.Content)

			template.HandleEditPost(posts, w)
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

			posts, err := db.GetPosts(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Printf("Deleted post: $s\n", selectedPostId)
			template.HandleEditPost(posts, w)
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

			posts, err := db.GetPosts(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Printf("Published post: $s\n", selectedPostId)
			template.HandleEditPost(posts, w)
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
		projects, err := db.GetProjects(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		template.HandleEditProject(projects, w)
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

			//2. retrieve file from posted form-data
			file, handler, err := req.FormFile("Thumbnail")
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

			selectedProjectId := req.FormValue("SelectedProject")
			project.Title = req.FormValue("Title")
			project.Description = req.FormValue("Description")
			project.Link = req.FormValue("Link")
			project.Thumbnail = fileName

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

			template.HandleEditProject(projects, w)
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
			template.HandleEditProject(projects, w)
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
		template.HandleEditKnowledgeTimeline(knowledgeTimeline, w)
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

			template.HandleEditKnowledgeTimeline(events, w)
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
			template.HandleEditKnowledgeTimeline(events, w)
		}

	})))

	r.PathPrefix("/projects").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//posts, err := db.GetPosts()
		//var projects []models.Project
		// projects = []models.Project{
		// 	models.Project{
		// 		Title:       "Bookings app",
		// 		Description: "An app where you can book shit",
		// 	},
		// 	models.Project{
		// 		Title:       "Animator",
		// 		Description: "An app where you can animate shit",
		// 	},
		// 	models.Project{
		// 		Title:       "Social network",
		// 		Description: "A social network that helps you connect with other creators, share your work and make money doing it.",
		// 	},
		// 	models.Project{
		// 		Title:       "Discord clone",
		// 		Description: "An alternative to the trainwreck that is Discord.",
		// 	},
		// 	models.Project{
		// 		Title:       "Video player",
		// 		Description: "Web based video player",
		// 	},
		// 	models.Project{
		// 		Title:       "Game engine",
		// 		Description: "Game engine for creating 2d video games",
		// 	},
		// 	models.Project{
		// 		Title:       "Bookings app",
		// 		Description: "An app where you can book shit",
		// 	},
		// 	models.Project{
		// 		Title:       "Animator",
		// 		Description: "An app where you can animate shit",
		// 	},
		// }

		ctx := req.Context()
		projects, err := db.GetProjects(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		template.HandleShowcase(projects, w)
	})

	r.PathPrefix("/blog/{postID}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		ctx := req.Context()
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
		template.HandlePost(links, post, w)
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

		posts, err = db.GetPosts(req.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		template.HandleBlog(posts, w)
	})

	r.PathPrefix("/knowledgetimeline").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		var timeline []models.TimelineEvent
		timeline, err = db.GetKnowledgeTimelineEvents(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		template.HandleKnowledgeTimeline(timeline, w)
	})

	r.PathPrefix("/content/images/{imageID}").Methods("GET").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		//ctx := req.Context()
		imageBytes, err := fh.GetFile(vars["imageID"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image")
		w.Write(imageBytes)
	})))
	r.PathPrefix("/content/images").Methods("POST").Handler(Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

	r.PathPrefix("/content/").Handler(http.StripPrefix("/content/", http.FileServer(http.Dir("./webportal/content"))))

	return http.ListenAndServe(addr, r)
}
